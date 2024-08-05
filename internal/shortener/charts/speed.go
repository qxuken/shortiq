package shortener_charts

import (
	"fmt"
	"math"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/qxuken/short/internal/shortener"
)

func combineTimings(data []int) map[int]int {
	counts := make(map[int]int)

	for _, value := range data {
		counts[value]++
	}

	return counts
}

func generateScatterData(runs int, count int, xIndex int) []opts.ScatterData {
	timings := make([]int, runs)
	shortener.ShortUrlWithLen(count)
	fmt.Println("Start for len ", count)
	var start time.Time
	var duration time.Duration
	for i := 0; i < runs; i++ {
		start = time.Now()
		shortener.ShortUrlWithLen(count)
		duration = time.Since(start)
		timings[i] = int(duration.Nanoseconds())
	}
	combinedTimings := combineTimings(timings)
	fmt.Println(combinedTimings)
	data := []opts.ScatterData{}
	for value, count := range combinedTimings {
		data = append(data, opts.ScatterData{
			Name:         fmt.Sprintf("%v", count),
			Value:        value,
			Symbol:       "roundRect",
			SymbolSize:   int(math.Min(math.Max(10, float64(count)), 50)),
			SymbolRotate: 10,
			XAxisIndex:   xIndex,
		})
	}
	fmt.Printf("Finished(c=%v) %v\n", count, combinedTimings)
	return data
}

func GenerateSpeedChart(runs int, countGroups [][]int) *charts.Scatter {
	fmt.Println("Generating speed chart")
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Hash speed"}),
	)

	for ig, group := range countGroups {
		for ic, count := range group {
			data := generateScatterData(runs, count, ig+ic+1)
			scatter.AddSeries(fmt.Sprintf("%v", count), data)
		}
	}

	return scatter
}
