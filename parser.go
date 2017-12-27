package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Parse will return a slice of tasks built from a Efs2file
func Parse(f string) ([]*task, error) {
	var tasks []*task
	var s *bufio.Scanner

	if f == "-" {
		s = bufio.NewScanner(os.Stdin)
	} else {
		fh, err := os.Open(f)
		if err != nil {
			return nil, fmt.Errorf("Could not read Efs2file - %s", err)
		}
		defer fh.Close()
		s = bufio.NewScanner(fh)
	}

	lc := 0
	for s.Scan() {
		lc++
		runRe := regexp.MustCompile(`^RUN (CMD|SCRIPT) .*$`)
		putRe := regexp.MustCompile(`^PUT .* \d{3,4}$`)

		l := strings.TrimSpace(s.Text())
		t := &task{
			command: &command{},
			file:    &file{},
		}

		if runRe.MatchString(l) {
			t.task = l
			c := strings.Split(l, " ")

			if c[1] == "CMD" {
				t.command.cmd = strings.Join(c[2:], " ")
				t.command.active = true
			}

			if c[1] == "SCRIPT" {
				dest := "/tmp/" + TmpFn()
				t.command.cmd = dest + "; rm " + dest
				t.command.active = true
				t.file.active = true
				t.file.source = c[2]
				t.file.dest = dest
				t.file.mode = os.FileMode(int(0700))
			}

			tasks = append(tasks, t)
		}

		if putRe.MatchString(l) {
			t.task = l
			t.file.active = true

			p := strings.Split(l, " ")
			if len(p) != 4 {
				return tasks, fmt.Errorf("Invlalid PUT definition on line %d", lc)
			}

			t.file.source = p[1]
			t.file.dest = p[2]

			m, err := strconv.ParseUint(p[3], 8, 32)
			if err != nil {
				return tasks, fmt.Errorf("Could not convert value to integer on line %d - %s", lc, p[3])
			}
			t.file.mode = os.FileMode(m)

			tasks = append(tasks, t)
		}
	}
	if err := s.Err(); err != nil {
		return tasks, fmt.Errorf("Error parsing Efs2file - %s", err)
	}

	return tasks, nil
}

// tmpFn will generate a temporary filename
func TmpFn() string {
	// Snagged from ioutil.TempFile
	r := uint32(time.Now().UnixNano() + int64(os.Getpid()))
	r = r*1664525 + 1013904223
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}
