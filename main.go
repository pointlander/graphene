// Copyright 2021 The Graphene Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"fmt"
	"image/color/palette"
	"math"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	output, err := os.Create("README.md")
	if err != nil {
		panic(err)
	}
	defer output.Close()
	process("Pyrolytic graphite experiment", "log1.csv", output)
	process("Calibration", "log2.csv", output)
	process("Control", "log3.csv", output)
}

func process(title, log string, output *os.File) {
	fmt.Fprintf(output, "## %s - %s\n", title, log)
	input, err := os.Open(log)
	if err != nil {
		panic(err)
	}
	defer input.Close()
	decoder := csv.NewReader(input)
	record, err := decoder.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(log, record)
	sum, count := 0.0, 0
	points1, points2 := make(plotter.XYs, 0, 8), make(plotter.XYs, 0, 8)
	record, err = decoder.Read()
	for err == nil {
		t1, err1 := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		if err1 != nil {
			panic(err1)
		}
		t2, err1 := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)
		if err1 != nil {
			panic(err1)
		}
		fmt.Println(log, t1, t2)
		sum += math.Abs(t1 - t2)
		points1 = append(points1, plotter.XY{X: float64(count), Y: float64(t1)})
		points2 = append(points2, plotter.XY{X: float64(count), Y: float64(t2)})
		count++
		record, err = decoder.Read()
	}
	fmt.Println("average=", sum/float64(count))
	fmt.Fprintf(output, "* average=%f\n", sum/float64(count))

	deviation := func(values plotter.XYs) float64 {
		a, b, count := 0.0, 0.0, 0
		for _, value := range values {
			a += value.Y * value.Y
			b += value.Y
			count++
		}
		return math.Sqrt((a - b*b/float64(count)) / float64(count))
	}
	sigma1 := deviation(points1)
	sigma2 := deviation(points2)
	fmt.Println("sigma1=", sigma1)
	fmt.Println("sigma2=", sigma2)
	average := func(values plotter.XYs) float64 {
		sum, count := 0.0, 0
		for _, value := range values {
			sum += value.Y
			count++
		}
		return sum / float64(count)
	}
	average1 := average(points1)
	average2 := average(points2)
	fmt.Println("average1=", average1)
	fmt.Println("average2=", average2)
	corr, count := 0.0, 0
	for i := range points1 {
		corr += (points1[i].Y - average1) * (points2[i].Y - average2)
		count++
	}
	corr /= float64(count) * sigma1 * sigma2
	fmt.Println("corr=", corr)
	fmt.Fprintf(output, "* corr=%f\n", corr)
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "temperature vs time"
	p.X.Label.Text = "time"
	p.Y.Label.Text = "temperature"

	scatter, err := plotter.NewScatter(points1)
	if err != nil {
		panic(err)
	}
	scatter.GlyphStyle.Radius = vg.Length(1)
	scatter.GlyphStyle.Shape = draw.CircleGlyph{}
	scatter.GlyphStyle.Color = palette.WebSafe[0x00F]
	p.Add(scatter)

	scatter, err = plotter.NewScatter(points2)
	if err != nil {
		panic(err)
	}
	scatter.GlyphStyle.Radius = vg.Length(1)
	scatter.GlyphStyle.Shape = draw.CircleGlyph{}
	scatter.GlyphStyle.Color = palette.WebSafe[0x030]

	p.Add(scatter)

	image := strings.Replace(log, ".csv", ".png", 1)
	err = p.Save(8*vg.Inch, 8*vg.Inch, image)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(output, "\n![%s](%s?raw=true)\n", log, image)
}
