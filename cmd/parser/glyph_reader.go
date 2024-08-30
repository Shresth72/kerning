package main

import (
	"fmt"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
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

	for i := 0; i < len(coordinates); i++ {
		coordinates[i] = coordinates[max(0, i-1)]
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
	p.Title.Text = "Glyph Plot"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Draw Points
	pts := make(plotter.XYs, len(g.XCoordinates))
	for i := range pts {
		pts[i].X = float64(g.XCoordinates[i])
		pts[i].Y = float64(g.YCoordinates[i])
	}

	contourStartIndex := 0
	for _, contourEndIndex := range g.ContourEndIndices {
		numPointsInContour := contourEndIndex - contourStartIndex + 1
		points := make(plotter.XYs, numPointsInContour+1) // +1 to close the contour

		for i := 0; i < numPointsInContour; i++ {
			points[i].X = float64(g.XCoordinates[contourStartIndex+i])
			points[i].Y = float64(g.YCoordinates[contourStartIndex+i])
		}
		// Close the contour
		points[numPointsInContour] = points[0]

		line, err := plotter.NewLine(points)
		if err != nil {
			return err
		}
		p.Add(line)

		contourStartIndex = contourEndIndex + 1
	}
	if err := p.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
		return err
	}

	return nil
}
