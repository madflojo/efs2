package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func Test_TmpFn(t *testing.T) {
	if len(TmpFn()) == 0 {
		t.Errorf("TmpFn returned an empty String")
	}
}

func Test_Parse_nofile(t *testing.T) {
	_, err := Parse("/nofile")
	if err == nil {
		t.Errorf("Parse did not return the expected error")
	}
}

func Test_Parse(t *testing.T) {
	b := []byte("RUN CMD ls -la\nRUN SCRIPT /test.sh\nPUT somefile /somefile 0644")
	f, _ := ioutil.TempFile("/tmp/", "")
	defer os.Remove(f.Name())
	_, err := f.Write(b)
	if err != nil {
		t.Errorf("Error when creating test file - %s", err)
	}
	_ = f.Close()

	tasks, err := Parse(f.Name())
	if err != nil {
		t.Errorf("Unexpected error when calling Parse - %s", err)
	}

	if len(tasks) != 3 {
		t.Errorf("Parsing error occured when calling Parse, expected 3 tasks, got %d", len(tasks))
	}

	tsk := tasks[0]
	if tsk.command.active == false || tsk.file.active == true {
		t.Errorf("Parsing error, task 1 should be a run command")
	}

	tsk = tasks[1]
	if tsk.command.active == false || tsk.file.active == false {
		t.Errorf("Parsing error, task 2 should be a run script")
	}

	tsk = tasks[2]
	if tsk.command.active == true || tsk.file.active == false {
		t.Errorf("Parsing error, task 3 should be a put")
	}
}
