package main

import (
	"fmt"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

type GlyphData struct {
	XCoordinates      []int
	YCoordinates      []int
	ContourEndIndices []int
}

func NewGlyphData(xCoordinates, yCoordinates, contourEndIndices []int) *GlyphData {
	return &GlyphData{
		XCoordinates:      xCoordinates,
		YCoordinates:      yCoordinates,
		ContourEndIndices: contourEndIndices,
	}
}

func (g *GlyphData) String() string {
	var sb strings.Builder

	for i, idx := range g.ContourEndIndices {
		sb.WriteString(fmt.Sprintf("Contour End Index %d: %d\n", i, idx))
	}

	for i := 0; i < len(g.XCoordinates); i++ {
		sb.WriteString(fmt.Sprintf("Point %d: (%d, %d)\n", i, g.XCoordinates[i], g.YCoordinates[i]))
	}

	return sb.String()
}

func GetAllGlyphLocations(reader *FontReader, lookUp map[string]uint32) []uint32 {
	reader.GoTo(lookUp["maxp"] + 4)
	numGlyphs := reader.ReadUInt16()

	reader.GoTo(lookUp["head"])
	reader.SkipBytes(50) // Skip unused fields
	isTwoByteEntry := reader.ReadUInt16() == 0

	locationTableStart := lookUp["loca"]
	glyphTableStart := lookUp["glyf"]
	allGlyphLocations := make([]uint32, numGlyphs)

	for glyphIndex := 0; glyphIndex < int(numGlyphs); glyphIndex++ {
		var offset uint32
		if isTwoByteEntry {
			reader.GoTo(locationTableStart + uint32(glyphIndex)*2)
			offset = uint32(reader.ReadUInt16()) * 2
		} else {
			reader.GoTo(locationTableStart + uint32(glyphIndex)*4)
			offset = reader.ReadUInt32()
		}
		allGlyphLocations[glyphIndex] = glyphTableStart + offset
	}

	return allGlyphLocations
}

func ReadSimpleGlyph(reader *FontReader) *GlyphData {
	contourEndIndices := make([]int, reader.ReadUInt16())
	reader.SkipBytes(8)

	for i := range contourEndIndices {
		contourEndIndices[i] = int(reader.ReadUInt16())
	}

	numPoints := contourEndIndices[len(contourEndIndices)-1] + 1
	allFlags := make([]byte, numPoints)
	reader.SkipBytes(int(reader.ReadUInt16()))

	for i := 0; i < numPoints; i++ {
		flag, _ := reader.ReadByte()
		allFlags[i] = flag

		if FlagBitIsSet(flag, 3) {
			repeatCount, _ := reader.ReadByte()
			for r := 0; r < int(repeatCount); r++ {
				i++
				allFlags[i] = flag
			}
		}
	}

	coordX := ReadCoordinates(reader, allFlags, true)
	coordY := ReadCoordinates(reader, allFlags, false)

	return NewGlyphData(coordX, coordY, contourEndIndices)
}

func ReadCoordinates(reader *FontReader, allFlags []byte, readingX bool) []int {
	offsetSizeFlagBit := 1
	offsetSignOrSkipBit := 4
	if !readingX {
		offsetSizeFlagBit = 2
		offsetSignOrSkipBit = 5
	}
	coordinates := make([]int, len(allFlags))
	if len(coordinates) > 0 {
		coordinates[0] = 0
	}

	for i := 0; i < len(coordinates); i++ {
		if i > 0 {
			coordinates[i] = coordinates[i-1]
		}
		flag := allFlags[i]
		// onCurve := FlagBitIsSet(flag, 0)

		if FlagBitIsSet(flag, offsetSizeFlagBit) {
			offset, _ := reader.ReadByte()
			sign := -1
			if FlagBitIsSet(flag, offsetSignOrSkipBit) {
				sign = 1
			}
			coordinates[i] += int(offset) * sign
		} else if !FlagBitIsSet(flag, offsetSignOrSkipBit) {
			coordinates[i] += int(reader.ReadUInt16())
		}

	}

	return coordinates
}

func (g *GlyphData) PlotAndSave(filename string) error {
	p := plot.New()
	p.Title.Text = fmt.Sprintf("Glyph Plot (%s)", filename)
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	contourStartIndex := 0
	for _, contourEndIndex := range g.ContourEndIndices {
		numPointsInContour := contourEndIndex - contourStartIndex + 1

		points := make([]Vector2, numPointsInContour)
		for i := 0; i < numPointsInContour; i++ {
			points[i] = Vector2{
				X: float64(g.XCoordinates[contourStartIndex+i]),
				Y: float64(g.YCoordinates[contourStartIndex+i]),
			}
		}

		// Draw Bezier curves between points
		for i := 0; i < numPointsInContour; i+=2 {
			if err := DrawBezier(p, points[i], points[(i+1) % numPointsInContour], points[(i+2) % numPointsInContour], resolution); err != nil {
				return err
			}
		}

		contourStartIndex = contourEndIndex + 1
	}

	if err := p.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
		return err
	}

	return nil
}
