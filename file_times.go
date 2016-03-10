package main

import (
	"fmt"
	"os"
	"path/filepath"
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

	stat, err := os.Stat(path)
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
	for idx := range *times.list {
		err = (*times.list)[idx].Check()
		if err != nil {
			return
		}
	}
	return
}

func (time FileTime) Check() (err error) {
	stat, err := os.Stat(time.Path)
	switch {
	case os.IsNotExist(err):
		if time.Exists {
			return checkFailed{fmt.Sprintf("File %q is missing", time.Path)}
		}
	case err != nil:
		return err
	case !time.Exists:
		return checkFailed{fmt.Sprintf("File %q newly created", time.Path)}
	case stat.ModTime().Unix() > time.Modtime:
		return checkFailed{fmt.Sprintf("File %q is stale", time.Path)}
	}
	return nil
}

func (times *FileTimes) Marshal() string {
	return marshal(*times.list)
}

func (times *FileTimes) Unmarshal(from string) error {
	return unmarshal(from, times.list)
}
