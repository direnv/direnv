package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/direnv/direnv/v2/gzenv"
)

// FileTime represents a single recorded file status
type FileTime struct {
	Path    string `json:"path"`
	Modtime int64  `json:"modtime"`
	Exists  bool   `json:"exists"`
}

// FileTimes represent a record of all the known files and times
type FileTimes struct {
	list *[]FileTime
}

// NewFileTimes creates a new empty FileTimes
func NewFileTimes() (times FileTimes) {
	list := make([]FileTime, 0)
	times.list = &list
	return
}

// Update gets the latest stats on the path and updates the record.
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

// NewTime add the file on path, with modtime and exists flag to the list of known
// files.
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

// Check validates all the recorded file times
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

// CheckOne compares notes between the given path and the recorded times
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

// Check verifies that the file is good and hasn't changed
func (times FileTime) Check() (err error) {
	stat, err := getLatestStat(times.Path)

	switch {
	case os.IsNotExist(err):
		if times.Exists {
			logDebug("Stat Check: %s: gone", times.Path)
			return checkFailed{fmt.Sprintf("File %q is missing (Stat)", times.Path)}
		}
	case err != nil:
		logDebug("Stat Check: %s: ERR: %v", times.Path, err)
		return err
	case !times.Exists:
		logDebug("Check: %s: appeared", times.Path)
		return checkFailed{fmt.Sprintf("File %q newly created", times.Path)}
	case stat.ModTime().Unix() != times.Modtime:
		logDebug("Check: %s: stale (stat: %v, lastcheck: %v)",
			times.Path, stat.ModTime().Unix(), times.Modtime)
		return checkFailed{fmt.Sprintf("File %q has changed", times.Path)}
	}
	logDebug("Check: %s: up to date", times.Path)
	return nil
}

// Formatted shows the times in a user-friendly format.
func (times *FileTime) Formatted(relDir string) string {
	timeBytes, err := time.Unix(times.Modtime, 0).MarshalText()
	if err != nil {
		timeBytes = []byte("<<???>>")
	}
	path, err := filepath.Rel(relDir, times.Path)
	if err != nil {
		path = times.Path
	}
	return fmt.Sprintf("%q - %s", path, timeBytes)
}

// Marshal dumps the times into gzenv format
func (times *FileTimes) Marshal() string {
	return gzenv.Marshal(*times.list)
}

// Unmarshal loads the watches back from gzenv
func (times *FileTimes) Unmarshal(from string) error {
	return gzenv.Unmarshal(from, times.list)
}

func getLatestStat(path string) (os.FileInfo, error) {
	var lstatModTime int64
	var statModTime int64

	// Check the examine-a-symlink case first:
	lstat, err := os.Lstat(path)
	if err != nil {
		logDebug("getLatestStat,Lstat: %s: error: %v", path, err)
		return nil, err
	}
	lstatModTime = lstat.ModTime().Unix()

	stat, err := os.Stat(path)
	if err != nil {
		logDebug("getLatestStat,Stat: %s: error: %v (Lstat time: %v)",
			path, err, lstatModTime)
		return nil, err
	}
	statModTime = stat.ModTime().Unix()

	if lstatModTime > statModTime {
		logDebug("getLatestStat: %s: Lstat: %v, Stat: %v -> preferring Lstat",
			path, lstatModTime, statModTime)
		return lstat, nil
	}
	logDebug("getLatestStat: %s: Lstat: %v, Stat: %v -> preferring Stat",
		path, lstatModTime, statModTime)
	return stat, nil
}
