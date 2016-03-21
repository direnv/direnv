package main

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	times := NewFileTimes()
	times.Update("file_times.go")
	if len(*times.list) != 1 {
		t.Error("Length of updated list not 1")
	}

	if !(*times.list)[0].Exists {
		t.Error("Existing file marked not existing")
	}
}

func TestFTJsons(t *testing.T) {
	ft := FileTime{"something.txt", time.Now().Unix(), true}
	marshalled, err := json.Marshal(ft)
	if err != nil {
		t.Error("Filetime failed to marshal:", err)
	}
	if bytes.NewBuffer(marshalled).String() == "{}" {
		t.Error(ft, "marshals as empty object")
	}

}

func TestRoundTrip(t *testing.T) {
	watches := NewFileTimes()
	watches.Update("file_times.go")

	rt_chk := NewFileTimes()
	rt_chk.Unmarshal(watches.Marshal())

	compareFTs(t, watches, rt_chk, "length", func(ft FileTimes) interface{} { return len(*ft.list) })
	compareFTs(t, watches, rt_chk, "first path", func(ft FileTimes) interface{} { return (*ft.list)[0].Path })
}

func compareFTs(t *testing.T, left, right FileTimes, desc string, compare func(ft FileTimes) (res interface{})) {
	lc, rc := compare(left), compare(right)
	if lc != rc {
		t.Error("Filetimes didn't round trip.",
			"Original", desc, "was:", lc,
			"RT", desc, "was:", rc)
	}
}

func TestCanonicalAdds(t *testing.T) {
	fts := NewFileTimes()
	fts.NewTime("docs/../file_times.go", 0, true)
	fts.NewTime("file_times.go", 0, true)
	if len(*fts.list) > 1 {
		t.Error("Double add of the same file")
	}
}

func TestCheckPasses(t *testing.T) {
	fts := NewFileTimes()
	fts.Update("file_times.go")
	err := fts.Check()
	if err != nil {
		t.Error("Check that should pass fails with:", err)
	}
}

func TestCheckStale(t *testing.T) {
	fts := NewFileTimes()
	fts.NewTime("file_times.go", 0, true)
	err := fts.Check()
	if err == nil {
		t.Error("Check that should fail because stale passes")
	}
}

func TestCheckAppeared(t *testing.T) {
	fts := NewFileTimes()
	fts.NewTime("file_times.go", 0, false)
	err := fts.Check()
	if err == nil {
		t.Error("Check that should fail because appeared passes")
	}
}

func TestCheckGone(t *testing.T) {
	fts := NewFileTimes()
	fts.NewTime("nosuchfileevarright.go", time.Now().Unix()+1000, true)
	err := fts.Check()
	if err == nil {
		t.Error("Check that should fail because gone passes")
	}
}
