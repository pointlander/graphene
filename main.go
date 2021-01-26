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

// ReadMe is the README
const ReadMe = `# Pyrolytic graphite temperature experiment

## Rational
It has been found in an
[experiment](https://www.nextbigfuture.com/2020/10/tiny-energy-harvested-from-brownian-motion-could-replace-low-power-batteries.html)
that graphene vibrates thus producing usable energy.
Pyrolytic graphite is composed of many layers of graphene, so it should be warmer than the
surrounding environment.

## Experiment Setup
The experiment consists of one [k type thermocouple](https://en.wikipedia.org/wiki/Thermocouple#Type_K) attached to
[pyrolytic graphite](https://en.wikipedia.org/wiki/Pyrolytic_carbon) in a thermally
insulating container. Another thermocouple is outside of the container. Both thermocouples are
attached to a thermometer.

![setup](setup.png?raw=true)

## Experimental Results
The experiment measured the temperature difference between the two thermocouples averaged over a period
of about 80 minutes. On average the difference was found to be ~0.8F. The error in the measurement was
found to be less than 0.1F. The control experiment found the difference to be ~0.5F. This indicates that
the pyrolytic graphite is producing heat. The experiment was conducted at a room temperature of ~70F.

## Materials
* [pyrolytic graphite](https://unitednuclear.com/index.php?main_page=product_info&cPath=16_17_69&products_id=527)
* [meter](https://www.fluke.com/en-us/product/temperature-measurement/ir-thermometers/fluke-54-ii)
* [thermocouple](https://www.fluke.com/en-us/product/accessories/probes/fluke-80pk-1)

## Potential error
Pyrolytic graphite is conductive, so this could impact the thermocouple, but this would probably be noticable
in the meter readings.

## Data

`

func main() {
	output, err := os.Create("README.md")
	if err != nil {
		panic(err)
	}
	defer output.Close()
	_, err = output.WriteString(ReadMe)
	if err != nil {
		panic(err)
	}
	process(false, "Pyrolytic graphite experiment", "log1.csv", output)
	process(true, "Pyrolytic graphite experiment with heat shrink tubing", "log4.csv", output)
	process(true, "Pyrolytic graphite experiment with heat shrink tubing - 8 hours - night", "log5.csv", output)
	process(true, "Pyrolytic graphite experiment with heat shrink tubing - 8 hours - day", "log6.csv", output)
	process(false, "Calibration", "log2.csv", output)
	process(false, "Control", "log3.csv", output)
	process(true, "Control - 8 hours - night", "log7.csv", output)
	process(true, "Thermos Control - 8 hours - day ", "log8.csv", output)
}

func process(fluke bool, title, log string, output *os.File) {
	fmt.Fprintf(output, "### %s - %s\n", title, log)
	input, err := os.Open(log)
	if err != nil {
		panic(err)
	}
	defer input.Close()
	decoder := csv.NewReader(input)
	decoder.FieldsPerRecord = -1
	record, err := decoder.Read()
	if err != nil {
		panic(err)
	}
	t1i, t2i := 1, 2
	if fluke {
		t2i = 3
		for record[0] != "Reading" {
			record, _ = decoder.Read()
		}
	}
	fmt.Println(log, record)
	sum, count := 0.0, 0
	points1, points2 := make(plotter.XYs, 0, 8), make(plotter.XYs, 0, 8)
	record, err = decoder.Read()
	for err == nil {
		t1, err1 := strconv.ParseFloat(strings.TrimSpace(record[t1i]), 64)
		if err1 != nil {
			panic(err1)
		}
		t2, err1 := strconv.ParseFloat(strings.TrimSpace(record[t2i]), 64)
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
	fmt.Fprintf(output, "\n![%s](%s?raw=true)\n\n", log, image)
}
