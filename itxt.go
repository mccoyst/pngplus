// © 2014 Steve McCoy under the MIT license. See LICENSE for details.

package pngplus

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"
)

// EncodeITXT writes to w an iTXt chunk with data from s.
// If len(s) cannot fit in a uint32, it is truncated.
func EncodeITXT(w io.Writer, s string) error {
	var b bytes.Buffer
	err := binary.Write(&b, binary.BigEndian, uint32(len(s)))
	if err != nil {
		return err
	}

	b.WriteString("iTXt")
	b.WriteString(s)

	buf := b.Bytes()
	err = binary.Write(&b, binary.BigEndian, crc32.ChecksumIEEE(buf[4:]))
	if err != nil {
		return err
	}

	_, err = w.Write(b.Bytes())
	return err
}
