package main

import (
	"fmt"
	"os"
	"time"
)

type FileTime struct {
	path    string
	modtime time.Time
	exists  bool
}

type FileTimes []FileTime

func (times FileTimes) Update(path string) (err error) {
	var time *FileTime
	for idx := range times {
		if times[idx].path == path {
			time = &times[idx]
			break
		}
	}
	if time == nil {
		times = append(times, FileTime{path: path})
		time = &times[len(times)]
	}

	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		time.exists = false
	} else {
		if err != nil {
			return
		}
		time.modtime = stat.ModTime()
	}

	return
}

type checkFailed struct {
	message string
}

func (err checkFailed) Error() string {
	return err.message
}

func (times FileTimes) Check() (err error) {
	for idx := range times {
		err = times[idx].Check()
		if err != nil {
			return
		}
	}
	return
}

func (time FileTime) Check() (err error) {
	stat, err := os.Stat(time.path)
	switch {
	case os.IsNotExist(err):
		if time.exists {
			return checkFailed{fmt.Sprintf("File %q is missing", time.path)}
		}
	case err != nil:
		return err
	case stat.ModTime().After(time.modtime):
		return checkFailed{fmt.Sprintf("File %q is stale", time.path)}
	}
	return nil
}

func (times FileTimes) Marshal() string {
	return marshal(times)
}

func (times FileTimes) Unmarshal(from string) error {
	return unmarshal(from, times)
}
