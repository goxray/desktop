/*
Package netchart implements basic chart to generate small minimalistic chart of network usage.
*/
package netchart

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	canvas2 "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/ajstarks/fc"
	"github.com/ajstarks/fc/chart"
)

type NetChart interface {
	Container() *fyne.Container
	UpdateNamed(stats map[string][]float64, clrs map[string]color.RGBA, order []string) error
	Canvas() fc.Canvas
}

type Chart struct {
	canvas   fc.Canvas
	lineSize float64
}

func New(w, h, lineSize float64) NetChart {
	return &Chart{
		canvas: fc.Canvas{
			Container: container.NewWithoutLayout(canvas2.NewRectangle(color.RGBA{255, 255, 255, 0})),
			Width:     w,
			Height:    h,
		},
		lineSize: lineSize,
	}
}

func (c *Chart) Container() *fyne.Container {
	return c.canvas.Container
}

func (c *Chart) Canvas() fc.Canvas {
	return c.canvas
}

func (c *Chart) UpdateNamed(stats map[string][]float64, clrs map[string]color.RGBA, orders []string) error {
	if len(stats) == 0 {
		return fmt.Errorf("stats is empty")
	}

	if len(clrs) != len(stats) {
		return fmt.Errorf("required %d colors, got %d", len(stats), len(clrs))
	}

	prevLen := 0
	for k, v := range stats {
		if prevLen == 0 {
			prevLen = len(v)
			continue
		}
		if len(v) != prevLen {
			return fmt.Errorf("expected all items to be the same len, got diff on stats[%s]", k)
		}
	}

	charts := make([]chart.ChartBox, 0, len(stats))
	// Preserve order
	for _, name := range orders {
		charts = append(charts, chart.ChartBox{
			Title:     name,
			Data:      make([]chart.NameValue, 0, len(stats[name])),
			Color:     clrs[name],
			Top:       85,
			Bottom:    20,
			Left:      10,
			Right:     95,
			Minvalue:  0,
			Maxvalue:  0, // Will be dynamically updated from stats
			Zerobased: true,
		})
	}

	maxValue := 0.
	for i := range charts {
		// Add data to chart
		data := make([]chart.NameValue, 0, len(stats[charts[i].Title]))
		for _, sv := range stats[charts[i].Title] {
			maxValue = max(maxValue, sv)
			data = append(data, chart.NameValue{
				Value: sv,
				Label: fmt.Sprintf("%.2f", sv),
			})
		}

		charts[i].Data = data
	}

	if maxValue == 0 {
		maxValue = 100
	}
	for i := range charts {
		charts[i].Maxvalue = maxValue
	}

	rect := canvas2.NewRectangle(color.RGBA{200, 255, 255, 0})
	rect.SetMinSize(fyne.NewSize(float32(c.canvas.Width), float32(c.canvas.Height)))
	c.canvas.Container.RemoveAll()
	c.canvas.Container.Add(rect)
	defer c.canvas.Container.Refresh()

	size := 5.
	for i := range charts {
		charts[i].YAxis(c.canvas, size-0.3, 0, maxValue, maxValue/4, "%0.f", true)
		charts[i].Scatter(c.canvas, c.lineSize)
		charts[i].Line(c.canvas, c.lineSize)

	}

	// TODO: dynamically get from Chart, hardcoded for now
	//  (im too lazy to fix this right now, look at the beautifyl graph!!!)
	offset := 17.
	midx := charts[0].Left + ((charts[0].Right - charts[0].Left) / 2) + 14
	c.canvas.CText(midx, charts[0].Bottom-offset, size, charts[0].Title, charts[0].Color)

	midx = charts[1].Left + ((charts[1].Right - charts[1].Left) / 2) - 14
	c.canvas.CText(midx, charts[1].Bottom-offset, size, charts[1].Title, charts[1].Color)

	return nil
}
