// © 2014 Steve McCoy under the MIT license. See LICENSE for details.

package pngplus

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"
)

// EncodeIBinary writes to w a chunk with data from s.
// The chunk is given the type "biNy", because I just made it up,
// it's not necessary to decode the image, and it's safe to copy.
// If len(s) cannot fit in a uint32, it is truncated.
func EncodeBinary(w io.Writer, s []byte) error {
	var b bytes.Buffer
	err := binary.Write(&b, binary.BigEndian, uint32(len(s)))
	if err != nil {
		return err
	}

	b.WriteString("biNy")
	b.Write(s)

	buf := b.Bytes()
	err = binary.Write(&b, binary.BigEndian, crc32.ChecksumIEEE(buf[4:]))
	if err != nil {
		return err
	}

	_, err = w.Write(b.Bytes())
	return err
}