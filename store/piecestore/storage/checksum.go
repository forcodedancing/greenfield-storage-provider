package storage

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"io"
	"reflect"
	"strconv"

	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
)

var crc32c = crc32.MakeTable(crc32.Castagnoli)

func generateChecksum(rs io.ReadSeeker) string {
	if b, ok := rs.(*bytes.Reader); ok {
		v := reflect.ValueOf(b)
		data := v.Elem().Field(0).Bytes()
		return strconv.Itoa(int(crc32.Update(0, crc32c, data)))
	}

	var hash uint32
	crcBuffer := bufPool.Get().(*[]byte)
	defer bufPool.Put(crcBuffer)
	defer func() { _, _ = rs.Seek(0, io.SeekStart) }()
	for {
		n, err := rs.Read(*crcBuffer)
		hash = crc32.Update(hash, crc32c, (*crcBuffer)[:n])
		if err != nil {
			if err != io.EOF {
				return ""
			}
			break
		}
	}
	return strconv.Itoa(int(hash))
}

type checksumReader struct {
	io.ReadCloser
	expected uint32
	checksum uint32
}

func (c *checksumReader) Read(buf []byte) (n int, err error) {
	n, err = c.ReadCloser.Read(buf)
	c.checksum = crc32.Update(c.checksum, crc32c, buf[:n])
	if err == io.EOF && c.checksum != c.expected {
		return 0, fmt.Errorf("failed to verify checksum: %d != %d", c.checksum, c.expected)
	}
	return
}

func verifyChecksum(rc io.ReadCloser, checksum string) io.ReadCloser {
	if checksum == "" {
		return rc
	}
	expected, err := strconv.Atoi(checksum)
	if err != nil {
		log.Errorf("invalid crc32c: %s", checksum)
		return rc
	}
	return &checksumReader{rc, uint32(expected), 0}
}
