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
	key := "blah"
	lang := "en-US"
	s := "hello"

	err := EncodeITXT(&w, false, key, lang, s)
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
	expected := "iTXtblah\x00\x00\x00en-US\x00\x00hello"
	if content != expected {
		t.Fatalf("wrong content: %q vs. %q", content, expected)
	}

	var crc uint32
	err = binary.Read(bytes.NewBuffer(bs[len(bs)-4:]), binary.BigEndian, &crc)
	if err != nil {
		t.Fatalf("unexpected crc read error: %v", err)
	}

	excrc := crc32.ChecksumIEEE([]byte(expected))
	if crc != excrc {
		t.Fatalf("wrong crc: %d vs. %d", crc, excrc)
	}
}

func TestEncodeITXTCompressed(t *testing.T) {
	var w bytes.Buffer
	key := "blah"
	lang := "en-US"
	s := "hello"

	err := EncodeITXT(&w, true, key, lang, s)
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
	expected := "iTXtblah\x00\x01\x00en-US\x00\x00\x78\x9c\xca\x48\xcd\xc9\xc9\x07\x04\x00\x00\xff\xff\x06\x2c\x02\x15"
	if content != expected {
		t.Fatalf("wrong content:\n%x\nvs.\n%x", content, expected)
	}

	var crc uint32
	err = binary.Read(bytes.NewBuffer(bs[len(bs)-4:]), binary.BigEndian, &crc)
	if err != nil {
		t.Fatalf("unexpected crc read error: %v", err)
	}

	excrc := crc32.ChecksumIEEE([]byte(expected))
	if crc != excrc {
		t.Fatalf("wrong crc: %d vs. %d", crc, excrc)
	}
}
