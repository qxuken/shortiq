package main

import (
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	"github.com/dustin/go-humanize"

	"github.com/qxuken/short/internal/shortener"
)

type CollisionData struct {
	count int
	data  []opts.LineData
}

func generateTestMarks() ([]int, []string) {
	runs := []int{0, 100_000, 500_000, 1_000_000, 5_000_000, 10_000_000, 50_000_000, 100_000_000}
	xAxis := make([]string, len(runs))
	for i, run := range runs {
		xAxis[i] = humanize.Comma(int64(run))
	}
	fmt.Println("Runs ", xAxis)
	return runs, xAxis
}

func generateLineItems(runs []int, count int, ch chan CollisionData) {
	handles := map[string]bool{}
	data := make([]opts.LineData, len(runs))
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
		fmt.Printf("Finished(c=%v) %v, collisions %v\n", count, humanize.Comma(int64(run)), humanize.Comma(int64(collisions)))
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
	run := func(counts []int) {
		for _, count := range counts {
			go generateLineItems(runs, count, ch)
		}
		for range counts {
			data := <-ch
			line.AddSeries(fmt.Sprintf("%v", data.count), data.data)
		}
	}
	run([]int{5, 6, 7})
	run([]int{8, 9})
	run([]int{9})
	return line
}

func main() {
	page := components.NewPage()
	page.AddCharts(generateCollisionChart())

	f, err := os.Create("./tmp/chart.html")
	if err != nil {
		panic(err)
	}
	page.Render(f)
}
