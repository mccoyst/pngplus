// © 2014 Steve McCoy under the MIT license. See LICENSE for details.

package pngplus

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"io/ioutil"
)

// EncodeBinary writes to w a chunk with data from s.
// The chunk is given the type "biNy", because I just made it up,
// it's not necessary to decode the image, and it's safe to copy.
// If len(s) cannot fit in a uint32, it is truncated.
func EncodeBinary(w io.Writer, s []byte) error {
	var b bytes.Buffer
	_ = binary.Write(&b, binary.BigEndian, uint32(len(s)))

	b.WriteString("biNy")
	b.Write(s)

	buf := b.Bytes()
	_ = binary.Write(&b, binary.BigEndian, crc32.ChecksumIEEE(buf[4:]))

	_, err := w.Write(b.Bytes())
	return err
}

var ErrNotBinary = errors.New("not a biNy")

// DecodeBinary tries to read a "biNy" chunk from r.
// You should have decoded a PNG from r before calling this
// to handle parsing of the initial PNG header, etc.
// DecodeBinary will return ErrNotBinary if the current chunk in the reader is not biNy,
// and scan to the next chunk.
func DecodeBinary(r io.Reader) ([]byte, error) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(buf[:4])
	crc := crc32.NewIEEE()
	crc.Write(buf[4:8])

	if string(buf[4:8]) != "biNy" {
		_, err = io.CopyN(ioutil.Discard, r, int64(length) + 4)
		if err != nil {
			return nil, err
		}
		return nil, ErrNotBinary
	}

	b := make([]byte, length)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	crc.Write(b)

	_, err = io.ReadFull(r, buf[:4])
	if err != nil {
		return nil, err
	}
	if binary.BigEndian.Uint32(buf[:4]) != crc.Sum32() {
		return nil, errors.New("Invalid checksum")
	}

	return b, nil
}
