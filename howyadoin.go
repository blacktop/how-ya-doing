package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gizak/termui"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {

	var err error

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	viewsPerDay, _, err := client.Repositories.ListTrafficViews(ctx, "maliceio", "malice", &github.TrafficBreakdownOptions{Per: "day"})
	if err != nil {
		log.Fatal(err)
	}

	viewsPerWeek, _, err := client.Repositories.ListTrafficViews(ctx, "maliceio", "malice", &github.TrafficBreakdownOptions{Per: "week"})
	if err != nil {
		log.Fatal(err)
	}

	dayViews := []int{}
	weekViews := []int{}

	fmt.Println("######  How Ya Doing  #############")
	fmt.Println()
	fmt.Println("View Count: ", viewsPerDay.GetCount())
	fmt.Println("Unique Views: ", viewsPerDay.GetUniques())
	fmt.Println()
	fmt.Println("========= Break Down by Week =======")
	for _, tdata := range viewsPerWeek.Views {
		fmt.Println("Week: ", tdata.GetTimestamp())
		fmt.Println("View Count: ", tdata.GetCount())
		weekViews = append(weekViews, tdata.GetCount())
		fmt.Println("Unique Views: ", tdata.GetUniques())
		fmt.Println("====================================")
	}
	fmt.Println("========= Break Down by Day =======")
	for _, tdata := range viewsPerDay.Views {
		fmt.Println("Day: ", tdata.GetTimestamp())
		fmt.Println("View Count: ", tdata.GetCount())
		dayViews = append(dayViews, tdata.GetCount())
		fmt.Println("Unique Views: ", tdata.GetUniques())
		fmt.Println("====================================")
	}

	fmt.Println("dayViews: ", dayViews)
	fmt.Println("weekViews: ", weekViews)

	err = termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	// View Overview
	overviewWidget := termui.NewList()
	overviewWidget.Items = []string{
		fmt.Sprintf("View Count: %d", viewsPerDay.GetCount()),
		fmt.Sprintf("Unique Views: %d", viewsPerDay.GetUniques()),
	}
	overviewWidget.ItemFgColor = termui.ColorYellow
	overviewWidget.BorderLabel = "Overview"
	overviewWidget.Align()
	overviewWidget.PaddingLeft = 2
	overviewWidget.Height = 4
	// Weekly Views Widget
	weekWidget := termui.NewBarChart()
	bclabels := []string{"W0", "W1", "W2", "W3"}
	weekWidget.BorderLabel = "Weekly Views"
	weekWidget.Data = weekViews
	weekWidget.Height = 10
	weekWidget.BarColor = termui.ColorBlue
	weekWidget.BarWidth = 5
	weekWidget.DataLabels = bclabels
	weekWidget.TextColor = termui.ColorGreen
	weekWidget.NumColor = termui.ColorBlack
	// Daily View Widget
	dayWidget := termui.NewBarChart()
	dayWidget.BorderLabel = "Daily Views"
	dayWidget.Data = dayViews
	dayWidget.Height = 10
	dayWidget.BarColor = termui.ColorMagenta
	dayWidget.BarWidth = 5
	dayWidget.DataLabels = []string{"D0", "D1", "D2", "D3", "D4", "D5", "D6", "D7", "D8", "D9", "D10", "D11", "D12", "D13", "D14", "D15", "D16", "D17", "D18"}
	dayWidget.TextColor = termui.ColorGreen
	dayWidget.NumColor = termui.ColorBlack

	// spl3 := termui.NewSparkline()
	// spl3.Data = dayViews
	// spl3.Title = "Enlarged Sparkline"
	// spl3.Height = 7
	// spl3.
	// spl3.LineColor = termui.ColorYellow
	//
	// spls2 := termui.NewSparklines(spl3)
	// spls2.Height = 10
	// // spls2.Width = 30
	// spls2.BorderFg = termui.ColorCyan
	// // spls2.X = 21
	// spls2.BorderLabel = "Daily View"

	// build
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(6, 0, overviewWidget)),
		termui.NewRow(
			termui.NewCol(2, 0, weekWidget),
			termui.NewCol(10, 0, dayWidget)))

	// calculate layout
	termui.Body.Align()

	termui.Render(termui.Body)

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Loop()
}
