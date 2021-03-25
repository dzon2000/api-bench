package main

import (
	"fmt"
	"github.com/dzon2000/benchmark/rest"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
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
	table1.Rows = [][]string{
		[]string{"120 ms"},
		[]string{"80 ms"},
		[]string{"60 ms"},
		[]string{"63 ms"},
		[]string{"81 ms"},
		[]string{"77 ms"},
		[]string{"54 ms"},
	}
	table1.TextStyle = ui.NewStyle(ui.ColorWhite)
	table1.SetRect(0, 0, 60, 10)

	table2 := widgets.NewTable()
	table2.Title = "Last request"
	table2.TitleStyle.Fg = ui.ColorMagenta
	table2.Rows = [][]string{
		[]string{"6411 ms"},
	}
	table2.TextStyle = ui.NewStyle(ui.ColorWhite)
	table2.SetRect(0, 0, 60, 10)

	grid.Set(
		ui.NewRow(
			1.0/10,
			ui.NewCol(1.0, header),
		),
		ui.NewRow(
			2.0/10,
			ui.NewCol(1.0/2, g1),
			ui.NewCol(1.0/2, g2),
		),
		ui.NewRow(
			7.0/10,
			ui.NewCol(1.0/2, table2),
			ui.NewCol(1.0/2, table1),
		),
	)
	ui.Render(grid)

	uiEvents := ui.PollEvents()

	ticker := make(chan int64, 1)
	go f("https://edf-energy-nft.cloud.thingworx.com/Thingworx/Things/5b3c1288-d66c-4d8b-858a-57eb1de3a9e6Forecast/Services/GetDataTableEntryCount", "2c773032-5f13-4a8d-b643-326e8b45cabc", ticker)
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case tick := <-ticker:
			table1.Rows = prepend(table1.Rows, fmt.Sprint(tick))
			ui.Render(table1)
		}
	}
}
