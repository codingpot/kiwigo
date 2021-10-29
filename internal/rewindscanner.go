package internal

import (
	"bufio"
	"io"
)

// RewindScanner is a scanner that rewinds the underlying reader.
type RewindScanner struct {
	readSeeker io.ReadSeeker
	scanner    *bufio.Scanner
}

// Rewind rewinds the underlying reader to the beginning.
func (s *RewindScanner) Rewind() {
	s.readSeeker.Seek(0, io.SeekStart)
	s.scanner = bufio.NewScanner(s.readSeeker)
}

func (s *RewindScanner) Scan() bool {
	return s.scanner.Scan()
}

func (s *RewindScanner) Text() string {
	return s.scanner.Text()
}

func NewRewindScanner(seeker io.ReadSeeker) *RewindScanner {
	return &RewindScanner{
		readSeeker: seeker,
		scanner:    bufio.NewScanner(seeker),
	}
}
