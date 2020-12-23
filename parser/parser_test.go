package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type testCase struct {
	// pass should indicated if the test is expected to pass
	pass bool
	// data contains the Efs2file contents being tested
	data []byte
	// instructions is a count of instructions that should be returned
	instructions int
	// name is the name of the test case
	name string
}

func TestTmpFn(t *testing.T) {
	if len(TmpFn()) == 0 {
		t.Errorf("TmpFn returned an empty String")
	}
}

func TestParsenofile(t *testing.T) {
	_, err := Parse("/nofile")
	if err == nil {
		t.Errorf("Parse did not return the expected error")
	}
}

func TestParsing(t *testing.T) {
  t.Parallel()
	// Define Test Cases
	var cc []testCase
	c := testCase{
		name: "Empty File",
		data: []byte(""),
		pass: true,
	}
	cc = append(cc, c)
	c = testCase{
		name:         "Base example",
		data:         []byte("RUN CMD ls -la\nRUN SCRIPT /test.sh\nPUT somefile /somefile 0644"),
		pass:         true,
		instructions: 3,
	}
	cc = append(cc, c)
	c = testCase{
		name:         "Invalid Put",
		data:         []byte("RUN ls -la\nPUT somefile nofile"),
		pass:         false,
		instructions: 1,
	}
	cc = append(cc, c)
	c = testCase{
		name:         "Invalid Mode",
		data:         []byte("RUN ls -la\nPUT somefile /somefile sevenfivefive"),
		pass:         false,
		instructions: 1,
	}
	cc = append(cc, c)
	c = testCase{
		name:         "Empty Lines",
		data:         []byte("RUN CMD ls -la\n\nRUN SCRIPT /test.sh"),
		pass:         true,
		instructions: 2,
	}
	cc = append(cc, c)
	c = testCase{
		name:         "Comments",
		data:         []byte("# This is a Comment\nRUN CMD ls -la\n"),
		pass:         true,
		instructions: 1,
	}
	cc = append(cc, c)
	c = testCase{
		name:         "MultiLine",
		data:         []byte("# This is a Comment\nRUN CMD ls -la \\ /tmp/somedir\nRUN this is a command\n"),
		pass:         true,
		instructions: 2,
	}
	cc = append(cc, c)

	// Execute Tests in a bunch of sub-tests
	for _, x := range cc {
		t.Run("Test Parsing - "+x.name, func(t *testing.T) {
      t.Parallel()
			// Create Temp File
			f, _ := ioutil.TempFile("/tmp/", "testing.*.txt")
			defer os.Remove(f.Name())
			_, err := f.Write(x.data)
			if err != nil {
				t.Errorf("Error creating test file - %s", err)
			}

			// Test Parsing
			tasks, err := Parse(f.Name())
			if err != nil && x.pass {
				t.Errorf("Unexpected error when calling Parser - %s", err)
			}
			if !x.pass && err == nil {
				t.Errorf("Unexpected success when calling Parser - tasks %+v", tasks)
			}

			if x.pass {
				if len(tasks) != x.instructions {
					t.Errorf("Parser did not return the expected number of tasks got %d, expected %d", len(tasks), x.instructions)
				}
			}
		})
	}

}

func TestMultiLineParsing(t *testing.T) {
	var cc []testCase
	c := testCase{
		name: "Slash at the end",
		pass: true,
		data: []byte("this is a multi \\\nline command"),
	}
	cc = append(cc, c)
	c = testCase{
		name: "Slash in the middle",
		pass: true,
		data: []byte("this is not a multi\\line command"),
	}
	cc = append(cc, c)
	c = testCase{
		name: "Slash with whitespace",
		pass: true,
		data: []byte("this is a multi\\ \n\tline command"),
	}
	cc = append(cc, c)
	c = testCase{
		name: "No Slash multiline",
		pass: false,
		data: []byte("this is not a multi\n\tline command"),
	}
	cc = append(cc, c)

	for _, x := range cc {
		t.Run("Test MultiLine Parser - "+x.name, func(t *testing.T) {
			cmd := fmt.Sprintf("%s", x.data)
			data := []byte("RUN " + cmd)

			// Create Temp File
			f, _ := ioutil.TempFile("/tmp/", "testing.*.txt")
			defer os.Remove(f.Name())
			_, err := f.Write(data)
			if err != nil {
				t.Errorf("Error creating test file - %s", err)
				return
			}

			// Test Parsing
			tasks, err := Parse(f.Name())
			if err != nil && x.pass {
				t.Errorf("Unexpected error when calling Parser - %s", err)
				return
			}
			if x.pass == false && err == nil {
				t.Errorf("Unexpected success when calling parser with invalid data")
				return
			}

			if len(tasks) < 1 {
				t.Errorf("Tasks count is not correct, got %d", len(tasks))
				return
			}

			if tasks[0].Command.Cmd != cmd && x.pass {
				t.Errorf("Command in task does not match base command got %#q expected %#q", tasks[0].Command.Cmd, cmd)
				return
			}
		})
	}
}

func TestParsingTasks(t *testing.T) {
	b := []byte("RUN ls -la\nRUN CMD ls -la\nRUN SCRIPT /test.sh\nPUT somefile /somefile 0644")
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

	if len(tasks) != 4 {
		t.Errorf("Parsing error occurred when calling Parse, expected 4 tasks, got %d", len(tasks))
	}

	tsk := tasks[0]
	if tsk.Command.Cmd == "" || tsk.File.Source != "" {
		t.Errorf("Parsing error, task 1 should be a run Command")
	}

	tsk = tasks[1]
	if tsk.Command.Cmd == "" || tsk.File.Source != "" {
		t.Errorf("Parsing error, task 2 should be a run Command")
	}

	tsk = tasks[2]
	if tsk.Command.Cmd == "" || tsk.File.Source == "" {
		t.Errorf("Parsing error, task 3 should be a run script")
	}

	tsk = tasks[3]
	if tsk.Command.Cmd != "" || tsk.File.Source == "" {
		t.Errorf("Parsing error, task 4 should be a put")
	}
}
