package shortener_charts

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	"github.com/dustin/go-humanize"

	"github.com/qxuken/short/internal/shortener"
)

type collisionData struct {
	count int
	data  []opts.LineData
}

func generateCollisionTestMarks(runs []int) []string {
	xAxis := make([]string, len(runs))
	for i, run := range runs {
		xAxis[i] = humanize.Comma(int64(run))
	}
	fmt.Println("Runs ", xAxis)
	return xAxis
}

func generateCollisionLineItems(runs []int, count int, ch chan collisionData) {
	handles := map[string]bool{}
	data := make([]opts.LineData, len(runs))
	fmt.Println("Start for len ", count)
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
	ch <- collisionData{count, data}
}

func newTrue() *bool {
	b := true
	return &b
}

func GenerateCollisionChart(runs []int, countGroups [][]int) *charts.Line {
	fmt.Println("Generating collision chart")
	xAxis := generateCollisionTestMarks(runs)

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Collision data",
			Subtitle: "Map handle len to collision",
		}))

	line.SetXAxis(xAxis).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: newTrue()}))

	ch := make(chan collisionData)
	run := func(group []int) {
		for _, count := range group {
			go generateCollisionLineItems(runs, count, ch)
		}
		for range group {
			data := <-ch
			line.AddSeries(fmt.Sprintf("%v", data.count), data.data)
		}
	}
	for _, group := range countGroups {
		run(group)
	}

	return line
}
