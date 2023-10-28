package goerr_test

import (
	"errors"
	"github.com/madhanganesh/goerr"
	"github.com/madhanganesh/goerr/samplesrc"
	"regexp"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	err := samplesrc.Repository()
	if err == nil {
		t.Error("expecting an error")
	}

	want := "error from database"
	got := err.Error()

	if want != got {
		t.Errorf("Want: %s; Got: %s", want, got)
	}
}

func TestNestedErrors(t *testing.T) {
	err := samplesrc.Service()

	want := "service failed"
	got := err.Error()

	if want != got {
		t.Errorf("Want: %s; Got: %s", want, got)
	}
}

func TestStackDetails(t *testing.T) {
	err := samplesrc.Service()

	stacks := goerr.ListStacks(err)
	if len(stacks) != 2 {
		t.Errorf("Nr. of stack entries. Want: %d; Got: %d", 2, len(stacks))
	}

	first := stacks[0]
	if !strings.Contains(first, "service failed") {
		t.Errorf("stack do not contain right error")
	}
	if !strings.Contains(first, "/goerr/samplesrc/samples.go:19") {
		t.Errorf("stack do not contain right file/line number")
	}

	second := stacks[1]
	if !strings.Contains(second, "error from database") {
		t.Errorf("stack do not contain right error")
	}
	if !strings.Contains(second, "/goerr/samplesrc/samples.go:26") {
		t.Errorf("stack do not contain right file/line number")
	}
}

func TestStack(t *testing.T) {
	err := samplesrc.Controller()
	stack := goerr.Stack(err)

	pattern := `controller failed \[.*/goerr/samplesrc/samples.go:11 \(samplesrc.Controller\)\]
\tservice failed \[.*/goerr/samplesrc/samples.go:19 \(samplesrc.Service\)\]
\t\terror from database.* \[.*/goerr/samplesrc/samples.go:26 \(samplesrc.Repository\)\]`
	match, _ := regexp.MatchString(pattern, stack)

	if !match {
		t.Errorf("stack is not matching the expectation")
	}
}

func TestStackNonGoErr(t *testing.T) {
	err := errors.New("some sample error")

	want := "some sample error"
	got := goerr.Stack(err)

	if want != got {
		t.Errorf("Want: %s, Got: %s", want, got)
	}
}
