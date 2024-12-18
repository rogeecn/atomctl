package utils

import (
	"bufio"
	"io"
)

// NewLogBuffer creates a buffer that can be used to capture output stream
// and write to a logger in real time
func NewLogBuffer(output func(string)) io.Writer {
	reader, writer := io.Pipe()

	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			output(scanner.Text())
		}
	}()

	return writer
}

// NewCombinedBuffer combines multiple io.Writers
func NewCombinedBuffer(writers ...io.Writer) io.Writer {
	return io.MultiWriter(writers...)
}
