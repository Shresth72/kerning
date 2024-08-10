package main

import (
	"fmt"
)

func ParseFont(fontPath string) error {
	fr := &FontReader{}

	err := fr.OpenFile(fontPath)
	if err != nil {
		return err
	}
	defer fr.CloseFile()

	err = fr.SkipBytes(4)
	if err != nil {
		return err
	}

	numTables, err := fr.ReadUint16()
	if err != nil {
		return err
	}

	err = fr.SkipBytes(6)
	if err != nil {
		return err
	}

	tagOffsetMap := make(map[string]uint32)

	for i := 0; i < int(numTables); i++ {
		tag, err := fr.ReadTag()
		if err != nil {
			return err
		}

		_, err = fr.ReadUint32()
		if err != nil {
			return err
		}

		offset, err := fr.ReadUint32()
		if err != nil {
			return err
		}

		_, err = fr.ReadUint32()
		if err != nil {
			return err
		}

		tagOffsetMap[tag] = offset
	}

	if offset, ok := tagOffsetMap["glyf"]; ok {
		err = fr.Goto(offset)
		if err != nil {
			return err
		}
	}

	glpyh0, err := ReadSimpleGlyph(fr)
	if err != nil {
		return err
	}
	glpyh0.Display()
	glpyh0.DrawTest()

	return nil
}

func main() {
	fontPath := "/home/shrestha/.fonts/Meditative.ttf"
	err := ParseFont(fontPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

}
