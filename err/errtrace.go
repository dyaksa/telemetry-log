package err

import (
	"path"
	"runtime"
	"strconv"
)

type errStack struct {
	Name string `bson:"name,omitempty" json:"name"`
	File string `bson:"file,omitempty" json:"file"`
	Line string `bson:"line,omitempty" json:"line"`
}

type ErrorTracer struct {
	stackTrace []uintptr
}

func (errTrace *ErrorTracer) Error() (s string) {
	return
}

func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(5, pcs[:])
	return pcs[0:n]
}

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

func (e *ErrorTracer) Err() error {
	return &ErrorTracer{
		stackTrace: callers(),
	}
}
