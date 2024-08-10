package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type FontReader struct {
	file *os.File
}

func (fr *FontReader) OpenFile(fontPath string) error {
	var err error
	fr.file, err = os.Open(fontPath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	return nil
}

func (fr *FontReader) SkipBytes(n int64) error {
	_, err := fr.file.Seek(n, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("error seeking in file: %w", err)
	}
	return nil
}

func (fr *FontReader) Goto(offset uint32) error {
	_, err := fr.file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return fmt.Errorf("error seeking to offset %d: %w", offset, err)
	}
	return nil
}

func (fr *FontReader) ReadUint16() (uint16, error) {
	var data uint16
	err := binary.Read(fr.file, binary.BigEndian, &data)
	if err != nil {
		return 0, fmt.Errorf("error reading from file: %w", err)
	}
	return data, nil
}

func (fr *FontReader) ReadUint32() (uint32, error) {
	var data uint32
	err := binary.Read(fr.file, binary.BigEndian, &data)
	if err != nil {
		return 0, fmt.Errorf("error reading from file: %w", err)
	}
	return data, nil
}

func (fr *FontReader) ReadTag() (string, error) {
	tagBytes := make([]byte, 4)
	_, err := fr.file.Read(tagBytes)
	if err != nil {
		return "", fmt.Errorf("error reading tag: %w", err)
	}
	return string(tagBytes), nil
}

func (fr *FontReader) CloseFile() {
	if fr.file != nil {
		fr.file.Close()
	}
}
