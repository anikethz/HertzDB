package fileio

import "io"

type ByteCounterWriter struct {
	Writer io.Writer
	Count  int64
}

func (bcw *ByteCounterWriter) Write(p []byte) (int, error) {
	n, err := bcw.Writer.Write(p)
	if err == nil {
		bcw.Count += int64(n)
	}
	return n, err
}
