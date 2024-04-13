// Package assert provides functions for specifying assertions.
package assert

// Precondition returns with no value if its first argument is nil.  Otherwise, it
// panics.  A precondition is an expression that must be true at the beginning of a
// function or process.
func Precondition(err error) {
	if err != nil {
		panic(err.Error())
	}
	return
}

// Postcondition return with no value if its first argument is nil.  Otherwise, it
// panics.  A postcondition is an expression that must be true at the end of a function
// or process.  Often the postcondition is an expression involving the calculated
// values from the function.
func Postcondition(err error) {
	if err != nil {
		panic(err.Error())
	}
	return
}

// Invariant returns with no value if its argument is nil.  Otherwise, it panics.
func Invariant(value bool) {
	if !value {
		panic("Invariant is false")
	}
	return
}

// Assert returns with no value if its first argument is true.  Otherwise, it panics.
func Assert(value bool, message string) {
	if !value {
		panic(message)
	}
	return
}
