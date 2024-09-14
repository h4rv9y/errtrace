package errtrace

import (
	"fmt"
)

type Builder Error

func create() *Builder {
	return &Builder{
		cause:   nil,
		msg:     "",
		context: make(map[string]interface{}),
		detail:  nil,
	}
}

func (b *Builder) copy() *Builder {
	nb := &Builder{
		cause:   b.cause,
		msg:     b.msg,
		context: make(map[string]interface{}),
		detail:  b.detail,
	}

	for k, v := range b.context {
		nb.context[k] = v
	}

	return nb
}

func (b *Builder) Detail(d interface{}) *Builder {
	nb := b.copy()
	nb.detail = d
	return nb
}

func (b *Builder) With(kv ...interface{}) *Builder {
	nb := b.copy()

	for i := 0; i < len(kv)-1; i += 2 {
		k := kv[i]
		v := kv[i+1]

		if key, ok := k.(string); ok {
			nb.context[key] = v
		}
	}

	return nb
}

func (b *Builder) Errorf(format string, args ...interface{}) error {
	nb := b.copy()
	nb.msg = fmt.Sprintf(format, args...)
	nb.stacktrace = newStacktrace()
	return (*Error)(nb)
}

func (b *Builder) Wrap(err error) error {
	if err == nil {
		return nil
	}

	nb := b.copy()
	nb.cause = err
	nb.stacktrace = newStacktrace()
	return (*Error)(nb)
}

func (b *Builder) Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	nb := b.copy()
	nb.cause = err
	nb.msg = fmt.Sprintf(format, args...)
	nb.stacktrace = newStacktrace()
	return (*Error)(nb)
}
