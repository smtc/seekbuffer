package seekbuffer

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"testing"
)

func testSR(t *testing.T, rs *SeekBuffer) {
	var (
		i, n  int
		err   error
		total int
		buf   = make([]byte, 1)
	)

	// 顺序读
	for err == nil {
		_, err = rs.Read(buf)
		if err == nil {
			c := []byte(fmt.Sprint(i + 1))
			if bytes.Compare(buf, c) != 0 {
				t.Logf("Read At %d, Read should be %s, but %s\n", i, string(c), string(buf))
				t.Fail()
			}
		}
		i++
	}

	total = 2 * len(rs.buf)
	for i = 0; i < total; i++ {
		n = rand.Intn(total)
		rs.Seek(int64(n), 0)
		_, err = rs.Read(buf)
		if n < len(rs.buf) {
			c := []byte(fmt.Sprint(n + 1))
			if bytes.Compare(buf, c) != 0 {
				t.Logf("Seek At %d, Read should be [%s], but [%s]\n", n, string(c), string(buf))
			}
		} else {
			if err != io.EOF {
				t.Logf("Seek At %d, Read should be nil, but [%s], seekBuffer: [%s]\n", n, string(buf), string(rs.buf))
				t.Fail()
			}
		}
	}
}

func TestSeekBuffer(t *testing.T) {
	bufs := [][]byte{
		[]byte("1"),
		[]byte("12"),
		[]byte("123"),
		[]byte("1234"),
		[]byte("12345"),
		[]byte("123456"),
	}

	for _, buf := range bufs {
		testSR(t, NewSeekBuffer(buf))
	}
}
