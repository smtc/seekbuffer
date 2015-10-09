package seekbuffer

import (
	"fmt"
	"io"
)

// SeekBuffer is read Only
type SeekBuffer struct {
	buf []byte
	off int
}

func NewSeekBuffer(buf []byte) *SeekBuffer {
	rs := &SeekBuffer{
		buf: make([]byte, len(buf)),
	}
	copy(rs.buf, buf[0:len(buf)])
	return rs
}

func (b *SeekBuffer) Len() int {
	return len(b.buf) - b.off
}

func (b *SeekBuffer) Bytes() []byte { return b.buf[b.off:] }
func (b *SeekBuffer) String() string {
	if b == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	return string(b.buf[b.off:])
}

func (b *SeekBuffer) Read(p []byte) (n int, err error) {
	if b.off >= len(b.buf) {
		if len(p) == 0 {
			return
		}
		return 0, io.EOF
	}
	n = copy(p, b.buf[b.off:])
	b.off += n

	return
}

//
// Seek sets the offset for the next Read or Write to offset, interpreted
// according to whence:
//   0 means relative to the origin of the file,
//   1 means relative to the current offset,
//   and 2 means relative to the end.
// Seek returns the new offset and an error, if any.
//
func (b *SeekBuffer) Seek(offset int64, whence int) (int64, error) {
	noffset := offset
	total := int64(len(b.buf))

	switch whence {
	case 0:
		if offset >= total {
			noffset = total
		}
		b.off = int(noffset)
	case 1:
		noffset += int64(b.off)
		if noffset >= total {
			noffset = total
		}
		b.off = int(noffset)
	case 2:
		if offset >= total {
			noffset = 0
			b.off = 0
		} else {
			noffset = total - offset
			b.off = int(noffset)
		}

	default:
		return 0, fmt.Errorf("invalid param whence: %d", whence)
	}

	return noffset, nil
}
