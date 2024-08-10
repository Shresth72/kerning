package main

import (
	"fmt"
)

type Point struct {
	X int
	Y int
}

type GlyphData struct {
	Points            []Point
	ContourEndIndices []int
}

func (g *GlyphData) Display() {
	for i, index := range g.ContourEndIndices {
		fmt.Printf("Contour End Index %d: %d\n", i, index)
	}

	for i, point := range g.Points {
		fmt.Printf("Point %d: (%d, %d)\n", i, point.X, point.Y)
	}
}

// Read
func ReadSimpleGlyph(reader *FontReader) (GlyphData, error) {
	// Read contour and indices
	contourEndIndices := make([]int, 0)
	numContours, err := reader.ReadUint16()
	if err != nil {
		return GlyphData{}, err
	}
	contourEndIndices = make([]int, numContours)

	// Skip bounds size (8 bytes)
	err = reader.SkipBytes(8)
	if err != nil {
		return GlyphData{}, err
	}

	for i := 0; i < len(contourEndIndices); i++ {
		index, err := reader.ReadUint16()
		if err != nil {
			return GlyphData{}, err
		}
		contourEndIndices[i] = int(index)
	}

	numPoints := contourEndIndices[len(contourEndIndices)-1] + 1
	allFlags := make([]byte, numPoints)

	// Skip instructions
	numFlags, err := reader.ReadUint16()
	if err != nil {
		return GlyphData{}, err
	}
	err = reader.SkipBytes(int64(numFlags))
	if err != nil {
		return GlyphData{}, err
	}

	for i := 0; i < numPoints; i++ {
		flag, err := reader.ReadUint8()
		if err != nil {
			return GlyphData{}, err
		}
		allFlags[i] = flag

		if FlagBitIsSet(flag, 3) {
			repeatCount, err := reader.ReadUint8()
			if err != nil {
				return GlyphData{}, err
			}

			for r := 0; r < int(repeatCount); r++ {
				if i+1 < len(allFlags) {
					allFlags[i+1] = flag
					i++
				}
			}
		}
	}

	coordsX := ReadCoordinates(reader, allFlags, true)
	coordsY := ReadCoordinates(reader, allFlags, false)

	points := make([]Point, numPoints)
	for i := 0; i < numPoints; i++ {
		points[i] = Point{X: coordsX[i], Y: coordsY[i]}
	}

	return GlyphData{
		Points:            points,
		ContourEndIndices: contourEndIndices,
	}, nil
}

func FlagBitIsSet(flag uint8, bitIndex int) bool {
	return (flag >> uint8(bitIndex) & 1) == 1
}

func ReadCoordinates(reader *FontReader, allFlags []byte, readingX bool) []int {
	var offsetSizeFlagBit, offsetSignOrSkipBit int

	if readingX {
		offsetSizeFlagBit = 1
		offsetSignOrSkipBit = 4
	} else {
		offsetSizeFlagBit = 2
		offsetSignOrSkipBit = 5
	}

	coordinates := make([]int, len(allFlags))
	for i := 0; i < len(coordinates); i++ {
		// Coordinate starts at previous value (0 if first coordinate)
		if i > 0 {
			coordinates[i] = coordinates[i-1]
		}
		flag := allFlags[i]

		// onCurve (TODO)
		_ = FlagBitIsSet(flag, 0)

		if FlagBitIsSet(flag, offsetSizeFlagBit) {
			offset, err := reader.ReadUint8()
			if err != nil {
				return nil
			}

			sign := 1
			if !FlagBitIsSet(flag, offsetSignOrSkipBit) {
				sign = -1
			}

			coordinates[i] += int(offset) * sign
		} else if !FlagBitIsSet(flag, offsetSignOrSkipBit) {
			offset, err := reader.ReadUint16()
			if err != nil {
				return nil
			}
			coordinates[i] += int(offset)
		}
	}

	return coordinates
}
