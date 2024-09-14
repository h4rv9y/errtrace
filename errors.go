package errtrace

import (
	"fmt"
	"strings"
)

type Error struct {
	cause   error
	msg     string
	context map[string]interface{}
	detail  interface{}

	stacktrace *stacktrace
}

func (e *Error) Error() string {
	if e.cause != nil {
		if e.msg == "" {
			return e.cause.Error()
		}

		return fmt.Sprintf("%s: %s", e.msg, e.cause.Error())
	}

	return e.msg
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Is(err error) bool {
	return e.cause == err
}

func (e *Error) Stacktrace() string {
	var blocks []string
	topFrame := ""

	recursive(e, func(e *Error) {
		if e.stacktrace != nil && len(e.stacktrace.frames) > 0 {
			msg := e.msg
			if msg == "" && e.cause != nil {
				msg = e.cause.Error()
			}
			if msg == "" {
				msg = "Error"
			}
			block := fmt.Sprintf("%s\n%s", msg, e.stacktrace.String(topFrame))

			blocks = append([]string{block}, blocks...)

			topFrame = e.stacktrace.frames[0].String()
		}
	})

	if len(blocks) == 0 {
		return ""
	}

	return "Cause: " + strings.Join(blocks, "\nThrown: ")
}

func (e *Error) Format(s fmt.State, verb rune) {
	if verb == 'v' && s.Flag('+') {
		_, _ = fmt.Fprint(s, e.formatVerbose())
	} else {
		_, _ = fmt.Fprint(s, e.formatSummary())
	}
}

func (e *Error) formatVerbose() string {
	output := new(strings.Builder)

	p := func(format string, args ...interface{}) {
		_, _ = fmt.Fprintf(output, format, args...)
	}

	p("%s\n", e.Error())

	if e.detail != nil {
		p("Detail: %+v\n", e.detail)
	}

	if len(e.context) > 0 {
		p("Context:\n")
		for k, v := range e.context {
			p("  * %s: %v\n", k, v)
		}
	}

	if stacktrace := e.Stacktrace(); stacktrace != "" {
		p("Stacktrace:\n")
		lines := strings.Split(stacktrace, "\n")
		p("  " + strings.Join(lines, "\n  "))
	}

	return output.String()
}

func (e *Error) formatSummary() string {
	return e.Error()
}
