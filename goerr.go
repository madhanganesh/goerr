package goerr

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
)

var MaxStackDepth = 50

type errorEx struct {
	err     error
	message string
	stack   []uintptr
	frames  []StackFrame
}

func New(nested error, message ...any) error {
	msg := "error"

	if nested != nil {
		msg = nested.Error()
	}
	if len(message) == 1 {
		msg = message[0].(string)
	}
	if len(message) > 1 {
		msg = fmt.Sprintf(message[0].(string), message[1:]...)
	}

	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2, stack[:])

	frames := make([]StackFrame, len(stack))

	for i, pc := range stack {
		frames[i] = NewStackFrame(pc)
	}

	return &errorEx{
		err:     nested,
		message: msg,
		stack:   stack[:length],
		frames:  frames,
	}

}

func (e *errorEx) Error() string {
	//return fmt.Sprintf("%s (%s:%d)\n", e.message, e.frames[0].File, e.frames[0].LineNumber)
	return e.message
}

// A StackFrame contains all necessary information about to generate a line
// in a callstack.
type StackFrame struct {
	// The path to the file containing this ProgramCounter
	File string
	// The LineNumber in that file
	LineNumber int
	// The Name of the function that contains this ProgramCounter
	Name string
	// The Package that contains this function
	Package string
	// The underlying ProgramCounter
	ProgramCounter uintptr
}

// NewStackFrame popoulates a stack frame object from the program counter.
func NewStackFrame(pc uintptr) (frame StackFrame) {

	frame = StackFrame{ProgramCounter: pc}
	if frame.Func() == nil {
		return
	}
	frame.Package, frame.Name = packageAndName(frame.Func())

	// pc -1 because the program counters we use are usually return addresses,
	// and we want to show the line that corresponds to the function call
	frame.File, frame.LineNumber = frame.Func().FileLine(pc - 1)
	return

}

// Func returns the function that contained this frame.
func (frame *StackFrame) Func() *runtime.Func {
	if frame.ProgramCounter == 0 {
		return nil
	}
	return runtime.FuncForPC(frame.ProgramCounter)
}

// String returns the stackframe formatted in the same way as go does
// in runtime/debug.Stack()
func (frame *StackFrame) String() string {
	str := fmt.Sprintf("%s:%d (0x%x)\n", frame.File, frame.LineNumber, frame.ProgramCounter)

	source, err := frame.sourceLine()
	if err != nil {
		return str
	}

	return str + fmt.Sprintf("\t%s: %s\n", frame.Name, source)
}

// SourceLine gets the line of code (from File and Line) of the original source if possible.
func (frame *StackFrame) SourceLine() (string, error) {
	source, err := frame.sourceLine()
	if err != nil {
		return source, err
	}
	return source, err
}

func (frame *StackFrame) sourceLine() (string, error) {
	if frame.LineNumber <= 0 {
		return "???", nil
	}

	file, err := os.Open(frame.File)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 1
	for scanner.Scan() {
		if currentLine == frame.LineNumber {
			return string(bytes.Trim(scanner.Bytes(), " \t")), nil
		}
		currentLine++
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "???", nil
}

func packageAndName(fn *runtime.Func) (string, string) {
	name := fn.Name()
	pkg := ""

	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//  runtime/debug.*T·ptrmethod
	// and want
	//  *T.ptrmethod
	// Since the package path might contains dots (e.g. code.google.com/...),
	// we first remove the path prefix if there is one.
	if lastslash := strings.LastIndex(name, "/"); lastslash >= 0 {
		pkg += name[:lastslash] + "/"
		name = name[lastslash+1:]
	}
	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}

	name = strings.Replace(name, "·", ".", -1)
	return pkg, name
}

func ListStacks(err error) []string {
	var result []string
	e, ok := err.(*errorEx)
	if !ok {
		result = append(result, err.Error())
		return result
	}
	packageParts := strings.Split(e.frames[0].Func().Name(), "/")
	funcName := packageParts[len(packageParts)-1]
	result = append(result, fmt.Sprintf("%s [%s:%d (%s)]", e.message, e.frames[0].File, e.frames[0].LineNumber, funcName))
	if e.err == nil {
		return result
	}

	return append(result, ListStacks(e.err)...)
}

func Stack(err error) string {
	stacks := ListStacks(err)
	if len(stacks) == 0 {
		return ""
	}
	if len(stacks) == 1 {
		return stacks[0]
	}

	var stack string
	for i, line := range stacks {
		stack += "\n"
		for j := 0; j < i; j++ {
			stack += "\t"
		}
		stack += line
	}
	return stack
}
