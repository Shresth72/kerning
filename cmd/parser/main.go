package main

import (
	"fmt"
	"log"
)

func main() {
	ParseFont("/home/shrestha/.fonts/Meditative.ttf")
}

func ParseFont(fontPath string) {
	fontReader, err := NewFontReader(fontPath)
	if err != nil {
		log.Fatalf("Failed to open font file: %v", err)
	}
	defer fontReader.Close()

	fontReader.SkipBytes(4)
	numTables := fontReader.ReadUInt16()
	fontReader.SkipBytes(6)
	fmt.Printf("NumTables: %d\n", numTables)

	tableLocationLookup := make(map[string]uint32)
	for i := 0; i < int(numTables); i++ {
		tag := fontReader.ReadTag()
		_ = fontReader.ReadUInt32() // checksum
		offset := fontReader.ReadUInt32()
		_ = fontReader.ReadUInt32() // length
		tableLocationLookup[tag] = offset
	}

	fontReader.GoTo(tableLocationLookup["glyf"])
	glyph0 := ReadSimpleGlyph(fontReader)
	fmt.Printf("Glyph 0:\n%s", glyph0)

	if err := glyph0.PlotAndSave("glyph0.png"); err != nil {
		log.Fatalf("Failed to save plot: %v", err)
	} else {
		fmt.Printf("Plot saved\n")
	}
}
