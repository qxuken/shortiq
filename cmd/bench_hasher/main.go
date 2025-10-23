package main

import (
	"fmt"
	"os"
	"path"

	"github.com/go-echarts/go-echarts/v2/components"

	shortener_charts "github.com/qxuken/short/internal/shortener/charts"
)

func main() {
	countGroups := [][]int{
		{5, 6, 7},
		{8, 9},
		{9},
	}

	page := components.NewPage()
	page.AddCharts(
		shortener_charts.GenerateCollisionChart(
			[]int{0, 100_000, 500_000, 1_000_000, 5_000_000}, //, 10_000_000, 50_000_000, 100_000_000},
			countGroups,
		),
		shortener_charts.GenerateSpeedChart(500, countGroups),
	)

	currentDirectory, _ := os.Getwd()
	path := path.Join(currentDirectory, "./tmp/chart.html")
	f, err := os.Create(path)
	fmt.Printf("Rendered file on path %v\n", path)
	if err != nil {
		panic(err)
	}
	page.Render(f)
}
