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
	dayPlotDateLabels := []string{}
	weekViews := []int{}
	weekPlotDateLabels := []string{}

	fmt.Println("######  How Ya Doing  #############")
	fmt.Println()
	fmt.Println("View Count: ", viewsPerDay.GetCount())
	fmt.Println("Unique Views: ", viewsPerDay.GetUniques())
	fmt.Println()
	fmt.Println("========= Break Down by Week =======")
	for _, tdata := range viewsPerWeek.Views {
		fmt.Println("Week: ", tdata.GetTimestamp())
		fmt.Println("View Count: ", tdata.GetCount())
		weekViews = append(weekViews, tdata.GetUniques())
		weekPlotDateLabels = append(weekPlotDateLabels, string(tdata.GetTimestamp().Format("1/2/06")))
		fmt.Println("Unique Views: ", tdata.GetUniques())
		fmt.Println("====================================")
	}
	fmt.Println("========= Break Down by Day =======")
	for _, tdata := range viewsPerDay.Views {
		fmt.Println("Day: ", tdata.GetTimestamp())
		fmt.Println("View Count: ", tdata.GetCount())
		dayViews = append(dayViews, tdata.GetUniques())
		dayPlotDateLabels = append(dayPlotDateLabels, string(tdata.GetTimestamp().Format("1/2")))
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
	weekWidget.BorderLabel = "Weekly Unique Views"
	weekWidget.Data = weekViews
	weekWidget.Height = 10
	weekWidget.BarColor = termui.ColorBlue
	weekWidget.BarWidth = 7
	weekWidget.BarGap = 2
	weekWidget.DataLabels = weekPlotDateLabels
	weekWidget.TextColor = termui.ColorGreen
	weekWidget.NumColor = termui.ColorBlack
	weekWidget.PaddingTop = 1
	weekWidget.PaddingLeft = 2
	// Daily View Widget
	dayWidget := termui.NewBarChart()
	dayWidget.BorderLabel = "Daily Unique Views"
	dayWidget.Data = dayViews
	dayWidget.Height = 14
	dayWidget.BarColor = termui.ColorMagenta
	dayWidget.BarWidth = 4
	dayWidget.BarGap = 2
	dayWidget.DataLabels = dayPlotDateLabels
	dayWidget.TextColor = termui.ColorGreen
	dayWidget.NumColor = termui.ColorBlack
	// dayWidget.PaddingTop = 1
	dayWidget.PaddingLeft = 2

	// build
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(3, 0, overviewWidget, weekWidget),
			termui.NewCol(9, 0, dayWidget)))

	// calculate layout
	termui.Body.Align()

	termui.Render(termui.Body)

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Loop()
}
