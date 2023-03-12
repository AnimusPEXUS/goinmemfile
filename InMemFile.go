package goinmemfile

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const debug = false

type ReadWriteSeekerCloser interface {
	io.Reader
	io.Writer
	io.Closer
	io.Seeker
}

var _ ReadWriteSeekerCloser = &InMemFile{}

type InMemFile struct {
	Buffer              []byte
	GrowOnWriteOverflow bool

	pos    int64
	closed bool

	debugName string
}

func NewInMemFileFromBytes(
	bytes []byte,
	pos int64,
	// this emulates normal file behavior
	GrowOnWriteOverflow bool,
) *InMemFile {
	self := new(InMemFile)
	self.Buffer = bytes
	self.pos = pos
	self.debugName = "[InMemFile]"
	self.GrowOnWriteOverflow = GrowOnWriteOverflow
	return self
}

func (self *InMemFile) SetDebugName(name string) {
	self.debugName = fmt.Sprintf("[%s]", name)
}

func (self *InMemFile) DebugPrintln(data ...any) {
	fmt.Println(append(append([]any{}, self.debugName), data...)...)
}

func (self *InMemFile) DebugPrintfln(format string, data ...any) {
	fmt.Println(append(append([]any{}, self.debugName), fmt.Sprintf(format, data...))...)
}

func (self *InMemFile) Read(p []byte) (n int, err error) {

	if debug {
		self.DebugPrintfln("Read(%d)", len(b))
	}

	defer func() {
		if debug {
			self.DebugPrintfln("Read defer", n, err)

		}
	}()

	if self.closed {
		return 0, os.ErrClosed
	}

	len_b := int64(len(self.Buffer))
	len_p := int64(len(p))

	if self.pos > len_b {
		return 0, errors.New("current position outside file")
	}

	// count := len_p
	len_b_m_pos := len_b - self.pos

	var eof bool = false

	if len_b_m_pos < len_p {
		eof = true
	}

	res := copy(p, self.Buffer)
	self.pos += int64(res)

	if res != len(p) {
		if eof {
			return res, io.EOF
		} else {
			return res, errors.New("copied less than len(p)")
		}
	}

	return res, nil
}

func (self *InMemFile) Write(p []byte) (n int, err error) {

	if self.closed {
		return 0, os.ErrClosed
	}

	len_p := int64(len(p))
	len_b := int64(len(self.Buffer))
	if (len_p + self.pos) > len_b {
		if self.GrowOnWriteOverflow {
			new_buf := make([]byte, len_b-self.pos+len_p)
			copy(new_buf, self.Buffer[:self.pos])
			self.Buffer = new_buf
		}
	}
	res := copy(self.Buffer[self.pos:], p)
	self.pos += int64(res)
	return res, nil
}

func (self *InMemFile) Close() error {
	if !self.closed {
		self.closed = true
		return nil
	} else {
		return os.ErrClosed
	}
}

func (self *InMemFile) Seek(offset int64, whence int) (int64, error) {
	if self.closed {
		return 0, os.ErrClosed
	}

	var new_pos int64

	len_b := int64(len(self.Buffer))

	switch whence {
	default:
		return 0, errors.New("invalid 'whence'")
	case io.SeekStart:
		new_pos = offset
		break
	case io.SeekEnd:
		new_pos = len_b + offset
		break
	case io.SeekCurrent:
		new_pos = self.pos + offset
		break
	}

	if new_pos < 0 || new_pos > len_b {
		return self.pos, os.ErrInvalid
	}

	self.pos = new_pos
	return new_pos, nil
}
