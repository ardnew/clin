// Package clin provides convenience methods that obtain user input for
// command-line Go applications.
//
// It is lightweight (uses packages from the standard library only) and easily
// integrates with complex flag parsing packages like "flag".
//
// There is one type, Input, and three of its methods exported, Args, Fields,
// and Reader. See the godoc comments on each of those methods for details.
//
// A global unexported variable of type Input is also defined, which is the
// target of top-level functions Args, Fields, and Reader.
// The function Default returns an Input initialized with the value of this
// global variable, whose fields can then be modified to fine-tune the behavior
// of each method.
// The default input stream is read from os.Stdin.
// Use field Stream to change this.
package clin

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Input configures the behavior of its exported functions Args and Reader.
//
// The top-level package functions always use the default configuration.
//
// To configure different behavior, either make a different Input or use
// function Default to start with an Input with default configuration.
type Input struct {
	// The default reader to read from when no arguments are given (typically
	// os.Stdin for command-line applications).
	Stream io.Reader
	// If true, always interpret input as a string literal, never a file path.
	Literal bool
	// When Args scans Stream for elements of the returned slice, the input
	// stream is tokenized using ArgsDelim as separator.
	ArgsDelim []byte
	// When Reader returns a strings.NewReader over the given slice args,
	// the elements of args are joined together, with ReadDelim as separator.
	ReadDelim []byte
	// Discard final Scanner token, if empty, when reading Stream in Args.
	skipToken bool
}

// input defines the default configuration and is the target of top-level
// functions Args, Fields, and Reader.
var input = Input{
	Stream:    os.Stdin,
	Literal:   false,
	ArgsDelim: []byte("\n"),
	ReadDelim: []byte(" "),
}

// Default returns an Input with default configuration.
func Default() Input { return input }

// Args returns the given string slice args if non-empty.
// Otherwise, a slice of each token read from Stream is returned, delimited by
// both CR+LF ("\r\n") and LF ("\n").
func Args(args []string) []string { return input.Args(args) }

// Fields is like Args, but with all empty (zeroed) string elements removed.
func Fields(args []string) []string { return input.Fields(args) }

// Reader returns an io.Reader over the string constructed by joining all
// elements in the given non-empty slice args, separated by one space (" ").
// If the given args contains a single element, and that element refers to
// a file path that we can open, then an io.Reader over the content of that
// file is returned.
// Otherwise, args is empty, returns Stream.
func Reader(args []string) io.Reader { return input.Reader(args) }

// Args returns the given string slice args if non-empty.
// Otherwise, a slice of each token read from Stream is returned, delimited by
// ArgsDelim.
func (in *Input) Args(args []string) []string {
	if len(args) == 0 {
		// No arguments: read lines from stdin.
		s := bufio.NewScanner(in.Stream)
		a := []string{}
		s.Split(in.scanArgs)
		in.skipToken = false
		for s.Scan() {
			if !in.skipToken {
				a = append(a, s.Text())
			}
		}
		return a
	}
	return args
}

// Fields is like Args, but with all empty (zeroed) string elements removed.
func (in *Input) Fields(args []string) []string {
	a := make([]string, 0, len(args))
	for _, s := range args {
		if s != "" {
			a = append(a, s)
		}
	}
	return a
}

// Reader returns an io.Reader over the string constructed by joining all
// elements in the given non-empty slice args, separated by ReadDelim.
// If the given args contains a single element, and that element refers to
// a file path that we can open, then an io.Reader over the content of that
// file is returned.
// Otherwise, args is empty, returns Stream.
func (in *Input) Reader(args []string) io.Reader {
	switch len(args) {
	case 0:
		// No arguments: read from Stream.
		return in.Stream
	case 1:
		if !in.Literal {
			// One argument: if it is a file path, read from the file.
			if r, err := os.Open(args[0]); nil == err {
				return r
			}
		}
		// One argument: not a file path, read the string itself.
		return strings.NewReader(args[0])
	default:
		// More than one argument: read from the string constructed by
		// joining all arguments, delimited by ReadDelim.
		return strings.NewReader(strings.Join(args, string(in.ReadDelim)))
	}
}

func (in *Input) scanArgs(data []byte, atEOF bool) (int, []byte, error) {

	n := len(in.ArgsDelim)

	// Split on each UTF-8 rune if ArgsDelim is empty.
	if n == 0 {
		return bufio.ScanRunes(data, atEOF)
	}
	for i := 0; i <= len(data)-n; i++ {
		if string(in.ArgsDelim) == string(data[i:i+n]) {
			// If ArgsDelim is a simple newline, also remove any trailing "\r"
			// that exists, which transparently handles Windows/DOS input.
			// Besides this one possible byte, all other trailing whitespace is
			// preserved in each token.
			j := i
			if i > 0 && data[i-1] == '\r' && n == 1 && in.ArgsDelim[0] == '\n' {
				j--
			}
			return i + n, data[:j], nil
		}
	}
	if !atEOF {
		return 0, nil, nil
	}
	// If the input is terminated with a delimiter, we reach here with a zero-
	// length slice data. Discard this empty, final token.
	// All other empty tokens (consecutive delimiters) are preserved.
	in.skipToken = len(data) == 0
	return 0, data, bufio.ErrFinalToken
}
