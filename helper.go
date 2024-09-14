package errtrace

func recursive(err *Error, tap func(*Error)) {
	tap(err)

	if err.cause == nil {
		return
	}

	if child, ok := As(err.cause); ok {
		recursive(child, tap)
	}
}
