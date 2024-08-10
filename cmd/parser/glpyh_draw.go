package main

import (
	"fmt"
	"image/color"
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func DrawLine(p1, p2 Point) {
	p := plot.New()
	p.Title.Text = "Glyph Line Plot"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	line, err := plotter.NewLine(plotter.XYs{
		{X: float64(p1.X), Y: float64(p1.Y)},
		{X: float64(p2.X), Y: float64(p2.Y)},
	})
	if err != nil {
		fmt.Println("error line", p1, p2)
		return
	}

	p.Add(line)
}

func DrawPoint(pt Point) {
	p := plot.New()

	p.Title.Text = "Glyph Point Plot"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	point, err := plotter.NewScatter(plotter.XYs{
		{X: float64(pt.X), Y: float64(pt.Y)},
	})
	if err != nil {
		fmt.Println("error point", pt)
		return
	}

	point.GlyphStyle.Color = color.RGBA{G: 255, A: 255}

	p.Add(point)
}

func (g *GlyphData) DrawTest() {
	p := plot.New()
	p.Title.Text = "Glyph Drawing"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	contourStartIndex := 0
	for _, contourEndIndex := range g.ContourEndIndices {
		numPointsInContour := contourEndIndex - contourStartIndex + 1

		lineData := make(plotter.XYs, numPointsInContour)
		for i := 0; i < numPointsInContour; i++ {
			lineData[i] = plotter.XY{X: float64(g.Points[contourStartIndex+i].X), Y: float64(g.Points[contourStartIndex+i].Y)}
		}

		line, _ := plotter.NewLine(lineData)
		line.LineStyle.Color = color.RGBA{R: 255, A: 255}
		p.Add(line)

		contourStartIndex = contourEndIndex + 1
	}

	pointData := make(plotter.XYs, len(g.Points))
	for i, point := range g.Points {
		pointData[i] = plotter.XY{X: float64(point.X), Y: float64(point.Y)}
	}

	points, _ := plotter.NewScatter(pointData)
	points.GlyphStyle.Color = color.RGBA{G: 255, A: 255} // Green point.
	p.Add(points)

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, "glyph_plot.png"); err != nil {
		log.Fatal(err)
	}
}
