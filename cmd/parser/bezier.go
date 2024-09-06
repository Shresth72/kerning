package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type Vector2 struct {
	X float64
	Y float64
}

func LinearInterpolation(start, end Vector2, t float64) Vector2 {
	return Vector2{
		X: start.X + (end.X-start.X)*t,
		Y: start.Y + (end.Y-start.Y)*t,
	}
}

func BezierInterpolation(p0, p1, p2 Vector2, t float64) Vector2 {
	intermediateA := LinearInterpolation(p0, p1, t)
	intermediateB := LinearInterpolation(p1, p2, t)
	return LinearInterpolation(intermediateA, intermediateB, t)
}

func DrawBezier(p *plot.Plot, p0, p1, p2 Vector2, resolution int) error {
	points := make(plotter.XYs, resolution+1)
	points[0] = plotter.XY{X: p0.X, Y: p0.Y}

	for i := 1; i <= resolution; i++ {
		t := float64(i+1) / float64(resolution)
		nextPointOnCurve := BezierInterpolation(p0, p1, p2, t)

		points[i].X = nextPointOnCurve.X
		points[i].Y = nextPointOnCurve.Y
	}

	line, err := plotter.NewLine(points)
	if err != nil {
		return err
	}
	p.Add(line)

	return nil
}
