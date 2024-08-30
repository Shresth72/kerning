package main

import (
	"fmt"
	"log"
)

func main() {
	ParseFont("../../assets/Meditative.ttf")
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

	allGlyphLocations := GetAllGlyphLocations(fontReader, tableLocationLookup)

	for i, glyphLocation := range allGlyphLocations {
		fontReader.GoTo(glyphLocation)
		glyphData := ReadSimpleGlyph(fontReader)
		// fmt.Printf("Glyph %d:\n%s", i, glyphData)

		if err := glyphData.PlotAndSave(fmt.Sprintf("glyphs/glyph%d.png", i)); err != nil {
			log.Fatalf("Failed to save plot for glyph %d: %v", i, err)
		} else {
			fmt.Printf("Plot for glyph %d saved\n", i)
		}
	}
}
