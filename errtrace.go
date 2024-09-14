package errtrace

import "errors"

func As(err error) (*Error, bool) {
	var t *Error
	ok := errors.As(err, &t)
	return t, ok
}

func With(kv ...interface{}) *Builder {
	return create().With(kv...)
}

func Detail(d interface{}) *Builder {
	return create().Detail(d)
}

func Errorf(format string, args ...interface{}) error {
	return create().Errorf(format, args...)
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}
	return create().Wrap(err)
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return create().Wrapf(err, format, args...)
}
