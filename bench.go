package main

import (
	"fmt"
	"log"
	"time"
	"math/rand"

	"github.com/dzon2000/benchmark/rest"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func prepend(x [][]string, y string) [][]string {
	x = append(x, []string{""})
	copy(x[1:], x)
	x[0] = []string{y}
	return x
}

func f(url, key string, c chan int64) {
	for {
		api := rest.Api{Url: url, AppKey: key}
		resp := api.Call()

		c <- resp.Time.Milliseconds()
	}
}

func fmtDuration(d time.Duration) string {
    d = d.Round(time.Second)
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
	d -= m * time.Hour
	s := d / time.Second
    return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	header := widgets.NewParagraph()
	header.Text = "Thingworx Benchmark 1.0.0"

	//p.SetRect(0, 0, 25, 5)
	g1 := widgets.NewGauge()
	g1.Title = "Ops / s"
	g1.SetRect(0, 6, 50, 11)
	g1.Percent = 1
	g1.BarColor = ui.ColorGreen
	g1.LabelStyle = ui.NewStyle(ui.ColorYellow)
	g1.TitleStyle.Fg = ui.ColorMagenta
	g1.BorderStyle.Fg = ui.ColorWhite
	g1.Label = fmt.Sprintf("%v ops/s", g1.Percent)

	g2 := widgets.NewGauge()
	g2.Title = "Ops / s"
	g2.SetRect(0, 6, 50, 11)
	g2.Percent = 100
	g2.BarColor = ui.ColorGreen
	g2.LabelStyle = ui.NewStyle(ui.ColorYellow)
	g2.TitleStyle.Fg = ui.ColorMagenta
	g2.BorderStyle.Fg = ui.ColorWhite
	g2.Label = fmt.Sprintf("%v ops/s", g2.Percent)

	table1 := widgets.NewTable()
	table1.Title = "Last request"
	table1.TitleStyle.Fg = ui.ColorMagenta
	table1.TextStyle = ui.NewStyle(ui.ColorWhite)
	table1.Rows = [][]string{
		[]string{""},
	}
	table1.SetRect(0, 0, 60, 10)

	table2 := widgets.NewTable()
	table2.Title = "Last request"
	table2.TitleStyle.Fg = ui.ColorMagenta
	table2.TextStyle = ui.NewStyle(ui.ColorWhite)
	table2.Rows = [][]string{
		[]string{""},
	}
	table2.SetRect(0, 0, 60, 10)

	requests1 := widgets.NewParagraph()
	requests1.Text = "# of requests: "
	requests2 := widgets.NewParagraph()
	requests2.Text = "# of requests: "

	timer := widgets.NewParagraph()
	timer.Text = "00:00:00"
	timer.Border = false

	grid.Set(
		ui.NewRow(
			1.0/10,
			ui.NewCol(1.0, header),
		),
		ui.NewRow(
			1.0/10,
			ui.NewCol(1.0, timer),
		),
		ui.NewRow(
			1.0/10,
			ui.NewCol(1.0/2, requests1),
			ui.NewCol(1.0/2, requests2),
		),
		ui.NewRow(
			1.0/10,
			ui.NewCol(1.0/2, g1),
			ui.NewCol(1.0/2, g2),
		),
		ui.NewRow(
			6.0/10,
			ui.NewCol(1.0/2, table1),
			ui.NewCol(1.0/2, table2),
		),
	)
	ui.Render(grid)

	uiEvents := ui.PollEvents()

	ticker1 := make(chan int64, 1)
	ticker2 := make(chan int64, 1)
	var count1, count2 int64
	var dur1, dur2 int64
	go f("http://thingworx-dev:8480/Thingworx/Things/d64ce9a4-7025-41cd-a7c3-e0626989c40fOutage/Services/GetDataTableEntryCount", "41cec232-325b-4f13-bdd9-0bcaf42e6a5c", ticker1)
	go f("https://edf-energy-nft.cloud.thingworx.com/Thingworx/Things/5b3c1288-d66c-4d8b-858a-57eb1de3a9e6Forecast/Services/GetDataTableEntryCount", "2c773032-5f13-4a8d-b643-326e8b45cabc", ticker2)
	timerBackend := time.NewTicker(time.Second).C
	start := time.Now()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "z":
				data := make([][]float64, 2)
				data[0] = make([]float64, 200)
				data[1] = make([]float64, 200)
				for i := 0; i < 200; i++ {
					data[0][i] = rand.Float64() * 150
					data[1][i] = 6000 + rand.Float64() * (200)
				}

				p1 := widgets.NewPlot()
				p1.Title = "dot-mode line Chart"
				p1.Marker = widgets.MarkerDot
				p1.Data = data
				p1.DotMarkerRune = '▀'
				p1.AxesColor = ui.ColorWhite
				p1.LineColors[0] = ui.ColorYellow
				p1.DrawDirection = widgets.DrawLeft
				grid.Set(
					ui.NewRow(
						1.0,
						ui.NewCol(1.0, p1),
					),
				)
				ui.Render(grid)
			}
		case <- timerBackend:
			timer.Text = fmtDuration(time.Since(start))
			ui.Render(timer)
		case tick := <-ticker1:
			table1.Rows = prepend(table1.Rows, fmt.Sprintf("%v ms", tick))
			ui.Render(table1)
			dur1 += tick
			count1++
			g1.Percent = int(count1 * 1000. / dur1)
			g1.Label = fmt.Sprintf("%v ops/s", g1.Percent)
			ui.Render(g1)
		case tick := <-ticker2:
			table2.Rows = prepend(table2.Rows, fmt.Sprintf("%v ms", tick))
			ui.Render(table2)
			dur2 += tick
			count2++
			g2.Percent = int(count2 * 1000. / dur2)
			g2.Label = fmt.Sprintf("%v ops/s", g2.Percent)
			ui.Render(g2)
		}
	}
}
