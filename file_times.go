package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/direnv/direnv/gzenv"
)

type FileTime struct {
	Path    string
	Modtime int64
	Exists  bool
}

type FileTimes struct {
	list *[]FileTime
}

func NewFileTimes() (times FileTimes) {
	list := make([]FileTime, 0)
	times.list = &list
	return
}

func (times *FileTimes) Update(path string) (err error) {
	var modtime int64
	var exists bool

	stat, err := getLatestStat(path)
	if os.IsNotExist(err) {
		exists = false
	} else {
		exists = true
		if err != nil {
			return
		}
		modtime = stat.ModTime().Unix()
	}

	err = times.NewTime(path, modtime, exists)

	return
}

func (times *FileTimes) NewTime(path string, modtime int64, exists bool) (err error) {
	var time *FileTime

	path, err = filepath.Abs(path)
	if err != nil {
		return
	}

	path = filepath.Clean(path)

	for idx := range *(times.list) {
		if (*times.list)[idx].Path == path {
			time = &(*times.list)[idx]
			break
		}
	}
	if time == nil {
		newTimes := append(*times.list, FileTime{Path: path})
		times.list = &newTimes
		time = &((*times.list)[len(*times.list)-1])
	}

	time.Modtime = modtime
	time.Exists = exists

	return
}

type checkFailed struct {
	message string
}

func (err checkFailed) Error() string {
	return err.message
}

func (times *FileTimes) Check() (err error) {
	if len(*times.list) == 0 {
		return checkFailed{"Times list is empty"}
	}
	for idx := range *times.list {
		err = (*times.list)[idx].Check()
		if err != nil {
			return
		}
	}
	return
}

func (times *FileTimes) CheckOne(path string) (err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	for idx := range *times.list {
		if time := (*times.list)[idx]; time.Path == path {
			err = time.Check()
			return
		}
	}
	return checkFailed{fmt.Sprintf("File %q is unknown", path)}
}

func (time FileTime) Check() (err error) {
	stat, err := getLatestStat(time.Path)

	switch {
	case os.IsNotExist(err):
		if time.Exists {
			log_debug("Stat Check: %s: gone", time.Path)
			return checkFailed{fmt.Sprintf("File %q is missing (Stat)", time.Path)}
		}
	case err != nil:
		log_debug("Stat Check: %s: ERR: %v", time.Path, err)
		return err
	case !time.Exists:
		log_debug("Check: %s: appeared", time.Path)
		return checkFailed{fmt.Sprintf("File %q newly created", time.Path)}
	case stat.ModTime().Unix() != time.Modtime:
		log_debug("Check: %s: stale (stat: %v, lastcheck: %v)",
			time.Path, stat.ModTime().Unix(), time.Modtime)
		return checkFailed{fmt.Sprintf("File %q has changed", time.Path)}
	}
	log_debug("Check: %s: up to date", time.Path)
	return nil
}

func (self *FileTime) Formatted(relDir string) string {
	timeBytes, err := time.Unix(self.Modtime, 0).MarshalText()
	if err != nil {
		timeBytes = []byte("<<???>>")
	}
	path, err := filepath.Rel(relDir, self.Path)
	if err != nil {
		path = self.Path
	}
	return fmt.Sprintf("%q - %s", path, timeBytes)
}

func (times *FileTimes) Marshal() string {
	return gzenv.Marshal(*times.list)
}

func (times *FileTimes) Unmarshal(from string) error {
	return gzenv.Unmarshal(from, times.list)
}

func getLatestStat(path string) (os.FileInfo, error) {
	var lstat_modtime int64
	var stat_modtime int64

	// Check the examine-a-symlink case first:
	lstat, err := os.Lstat(path)
	if err != nil {
		log_debug("getLatestStat,Lstat: %s: error: %v", path, err)
		return nil, err
	}
	lstat_modtime = lstat.ModTime().Unix()

	stat, err := os.Stat(path)
	if err != nil {
		log_debug("getLatestStat,Stat: %s: error: %v (Lstat time: %v)",
			path, err, lstat_modtime)
		return nil, err
	}
	stat_modtime = stat.ModTime().Unix()

	if lstat_modtime > stat_modtime {
		log_debug("getLatestStat: %s: Lstat: %v, Stat: %v -> preferring Lstat",
			path, lstat_modtime, stat_modtime)
		return lstat, nil
	} else {
		log_debug("getLatestStat: %s: Lstat: %v, Stat: %v -> preferring Stat",
			path, lstat_modtime, stat_modtime)
		return stat, nil
	}

}
