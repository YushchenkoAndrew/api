package helper

import (
	"io"
)

type StreamWriter struct {
	Result []byte
	io.Writer
}

func (s *StreamWriter) Write(p []byte) (n int, err error) {
	s.Result = append(s.Result, p...)
	return len(p), nil
}
