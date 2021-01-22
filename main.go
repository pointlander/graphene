// Copyright 2021 The Graphene Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"fmt"
	"image/color/palette"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	input, err := os.Open("log1.csv")
	if err != nil {
		panic(err)
	}
	defer input.Close()
	decoder := csv.NewReader(input)
	record, err := decoder.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println("log1.csv", record)
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
		fmt.Println("log1.csv", t1, t2)
		sum += t1 - t2
		points1 = append(points1, plotter.XY{X: float64(count), Y: float64(t1)})
		points2 = append(points2, plotter.XY{X: float64(count), Y: float64(t2)})
		count++
		record, err = decoder.Read()
	}
	fmt.Println("average=", sum/float64(count))

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

	err = p.Save(8*vg.Inch, 8*vg.Inch, "log1.png")
	if err != nil {
		panic(err)
	}
}
