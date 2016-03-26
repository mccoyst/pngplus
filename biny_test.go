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

	// Very short
	_, err := DecodeBinary(&w)
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}

	s := []byte("hello")

	// Normal
	w.Reset()
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
	if w.Len() != 0 {
		t.Fatalf("unexpected remains in w: %d", w.Len())
	}

	// Non-biNy
	w.Reset()
	err = EncodeBinary(&w, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	w.Bytes()[4] = 'q'

	result, err = DecodeBinary(&w)
	if err != ErrNotBinary {
		t.Fatalf("expected ErrNotBinary, got: %v", err)
	}
	if w.Len() != 0 {
		t.Fatalf("expected chunk to be skipped, have %d remaining", w.Len())
	}

	// Non-biNy and short
	w.Reset()
	err = EncodeBinary(&w, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	w.Bytes()[4] = 'q'

	result, err = DecodeBinary(&io.LimitedReader{R: &w, N: int64(w.Len())-1})
	if err != io.EOF {
		t.Fatalf("expected EOF, got: %v", err)
	}

	// biNy and short during content
	w.Reset()
	err = EncodeBinary(&w, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err = DecodeBinary(&io.LimitedReader{R: &w, N: int64(w.Len())-5})
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected UnexpectedEOF, got: %v", err)
	}

	// biNy and short during CRC
	w.Reset()
	err = EncodeBinary(&w, s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err = DecodeBinary(&io.LimitedReader{R: &w, N: int64(w.Len())-1})
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected UnexpectedEOF, got: %v", err)
	}
}
