// Package err provides functionality for error tracing.
package err

import (
	"path"
	"runtime"
	"strconv"
)

// errStack is a struct that holds the necessary information for an error stack.
type errStack struct {
	Name string `bson:"name,omitempty" json:"name"` // Name is the name of the function where the error occurred.
	File string `bson:"file,omitempty" json:"file"` // File is the name of the file where the error occurred.
	Line string `bson:"line,omitempty" json:"line"` // Line is the line number where the error occurred.
}

// ErrorTracer is a struct that holds the necessary information for error tracing.
type ErrorTracer struct {
	stackTrace []uintptr // stackTrace is a slice of program counters.
}

// Error is a method that returns an empty string.
// It is required to satisfy the error interface.
func (errTrace *ErrorTracer) Error() (s string) {
	return
}

// callers is a function that returns a slice of program counters.
// It skips the first 5 callers.
func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(5, pcs[:])
	return pcs[0:n]
}

// Print is a method that returns a slice of errStack.
// It creates an errStack for each program counter in the stack trace.
func (e *ErrorTracer) Print() []errStack {
	var traces []errStack

	for k := range e.stackTrace {
		v := e.stackTrace[k] - 1
		f := runtime.FuncForPC(v)
		file, line := f.FileLine(v)

		traces = append(traces, errStack{Name: f.Name(), File: path.Base(file), Line: strconv.Itoa(line)})

	}

	return traces
}

// Err is a method that returns a new ErrorTracer with the current stack trace.
func (e *ErrorTracer) Err() error {
	return &ErrorTracer{
		stackTrace: callers(),
	}
}
