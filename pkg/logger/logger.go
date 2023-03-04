package logger

import (
	"fmt"
	"runtime"
	"strings"
)

// GetStackTrace prints the stack trace with given error by the formatter.
// If the error is not traceable, empty string is returned.
func GetStackTrace() (trace []string) {
	stackBuf := make([]uintptr, 10)
	length := runtime.Callers(1, stackBuf[:])
	if length == 0 {
		// No stackBuf available. Stop now.
		// This can happen if the first argument to runtime.Callers is large.
		return
	}
	// pass only valid pcs to runtime.CallersFrames
	stack := stackBuf[:length]
	frames := runtime.CallersFrames(stack)

	// Loop to get frames.
	// A fixed number of stackBuf can expand to an indefinite number of Frames.
	for {
		frame, more := frames.Next()
		// To keep this output stable
		// even if there are changes in the testing package,
		// stop unwinding when we leave package runtime.
		// if !strings.Contains(frame.File, "runtime/") {
		// continue
		// }
		trace = append(trace, fmt.Sprintf("File: %s, Line: %d. Function: %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	return
}

// WhereAmI return current path and line
func WhereAmI(skipList ...int) string {
	var skip int
	if skipList == nil {
		skip = 1
	} else {
		skip = skipList[0]
	}
	function, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("Function: %s \nFile: %s:%d", chopPath(runtime.FuncForPC(function).Name(), "."), file, line)
}

// return the source filename after the last slash
func chopPath(original string, pathChar string) string {
	i := strings.LastIndex(original, pathChar)
	if i == -1 {
		return original
	}
	return original[i+1:]
}
