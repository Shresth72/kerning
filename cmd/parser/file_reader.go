package main

import (
	"encoding/binary"
	"io"
	"os"
)

type FontReader struct {
	file *os.File
}

func NewFontReader(pathToFont string) (*FontReader, error) {
	file, err := os.Open(pathToFont)
	if err != nil {
		return nil, err
	}
	return &FontReader{file: file}, nil
}

func (r *FontReader) ReadTag() string {
	tag := make([]byte, 4)
	r.file.Read(tag)
	return string(tag)
}

func (r *FontReader) ReadByte() (byte, error) {
	var b [1]byte
	r.file.Read(b[:])
	return b[0], nil
}

func (r *FontReader) ReadUInt16() uint16 {
	var value uint16
	binary.Read(r.file, binary.BigEndian, &value)
	return value
}

func (r *FontReader) ReadUInt32() uint32 {
	var value uint32
	binary.Read(r.file, binary.BigEndian, &value)
	return value
}

func (r *FontReader) SkipBytes(num int) {
	r.file.Seek(int64(num), io.SeekCurrent)
}

func (r *FontReader) GoTo(position uint32) {
	r.file.Seek(int64(position), io.SeekStart)
}

func (r *FontReader) Close() {
	r.file.Close()
}
