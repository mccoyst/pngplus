// © 2014 Steve McCoy under the MIT license. See LICENSE for details.

package pngplus

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"testing"
)

func TestEncodeITXT(t *testing.T) {
	var w bytes.Buffer
	s := "hello"

	err := EncodeITXT(&w, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bs := w.Bytes()
	if len(bs) != 4*3 + len(s) {
		t.Fatalf("wrong length for output: %x", bs)
	}

	var n uint32
	err = binary.Read(bytes.NewBuffer(bs), binary.BigEndian, &n)
	if err != nil {
		t.Fatalf("unexpected length read error: %v", err)
	}

	if n != uint32(len(s)) {
		t.Fatalf("wrong chunk length: %d vs. %d", n, len(s))
	}

	g := string(bs[4:8])
	if g != "iTXt" {
		t.Fatalf("wrong chunk type: %q", g)
	}

	g = string(bs[8:8+len(s)])
	if g != s {
		t.Fatalf("wrong chunk content: %q vs. %q", g, s)
	}

	var crc uint32
	err = binary.Read(bytes.NewBuffer(bs[8+len(s):]), binary.BigEndian, &crc)
	if err != nil {
		t.Fatalf("unexpected crc read error: %v", err)
	}

	if crc != crc32.ChecksumIEEE([]byte("iTXthello")) {
		t.Fatalf("wrong crc: %d vs. %d", crc, 666)
	}
}
