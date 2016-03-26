// © 2014 Steve McCoy under the MIT license. See LICENSE for details.

package pngplus

import (
	"bytes"
	"encoding/binary"
	"io"
	"hash/crc32"
	"testing"
)

func TestEncodeBinary(t *testing.T) {
	var w bytes.Buffer
	s := []byte("hello")

	err := EncodeBinary(&w, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bs := w.Bytes()

	var n uint32
	err = binary.Read(bytes.NewBuffer(bs), binary.BigEndian, &n)
	if err != nil {
		t.Fatalf("unexpected length read error: %v", err)
	}

	if n != uint32(len(s)) {
		t.Fatalf("wrong chunk length: %d vs. %d", n, len(s))
	}

	content := string(bs[4:len(bs)-4])
	expected := "biNyhello"
	if content != expected {
		t.Fatalf("wrong content: %q vs. %q", content, expected)
	}

	var crc uint32
	err = binary.Read(bytes.NewBuffer(bs[8+len(s):]), binary.BigEndian, &crc)
	if err != nil {
		t.Fatalf("unexpected crc read error: %v", err)
	}

	excrc := crc32.ChecksumIEEE([]byte(expected))
	if crc != excrc {
		t.Fatalf("wrong crc: %d vs. %d", crc, excrc)
	}
}

func TestDecodeBinary(t *testing.T) {
	var w bytes.Buffer

	_, err := DecodeBinary(&w)
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}

	s := []byte("hello")

	err = EncodeBinary(&w, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := DecodeBinary(&w)
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if string(result) != "hello" {
		t.Fatalf("unexpected result: %q", string(result))
	}
}
