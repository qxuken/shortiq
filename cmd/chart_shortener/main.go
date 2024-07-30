package main

import (
	"fmt"
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	"github.com/dustin/go-humanize"

	"github.com/qxuken/short/internal/shortener"
)

const (
	END int = 20_000_000
)

type CollisionData struct {
	count int
	data  []opts.LineData
}

func generateTestMarks() ([]int, []string) {
	runs := []int{}
	for run := 0; run <= END; run += 50_000 {
		runs = append(runs, run)
	}
	xAxis := make([]string, len(runs))
	for i, run := range runs {
		xAxis[i] = humanize.Comma(int64(run))
	}
	fmt.Println("Runs ", xAxis)
	return runs, xAxis
}

func generateLineItems(runs []int, count int, ch chan CollisionData) {
	data := make([]opts.LineData, len(runs))
	handles := map[string]bool{}
	fmt.Println("Runs for len ", count)
	collisions := 0
	for i, run := range runs {
		start := 0
		if i > 0 {
			start = runs[i-1]
		}
		for i := start; i < run; i++ {
			handle := shortener.ShortUrlWithLen(count)
			_, ok := handles[handle]
			if ok {
				collisions++
			} else {
				handles[handle] = true
			}
		}
		fmt.Printf("Finished(c=%v) %v\n", count, humanize.Comma(int64(run)))
		data[i] = opts.LineData{Value: collisions}
	}
	fmt.Println("Done for len ", count)
	ch <- CollisionData{count, data}
}

func newTrue() *bool {
	b := true
	return &b
}

func generateCollisionChart() *charts.Line {
	runs, xAxis := generateTestMarks()

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Collision data",
			Subtitle: "Map handle len to collision",
		}))

	line.SetXAxis(xAxis).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: newTrue()}))

	ch := make(chan CollisionData)
	for i := 3; i < 9; i++ {
		go generateLineItems(runs, i, ch)
	}
	for i := 3; i < 9; i++ {
		data := <-ch
		line.AddSeries(fmt.Sprintf("%v", data.count), data.data)
	}

	return line
}

func main() {
	page := components.NewPage()
	page.AddCharts(generateCollisionChart())

	f, err := os.Create("./tmp/chart.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
