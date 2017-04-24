package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gizak/termui"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func getStarDateSparklineData(firstStarDate time.Time, starMap map[string]int) []int {

	data := []int{}

	for d := firstStarDate; d.Unix() < time.Now().Unix(); d = d.AddDate(0, 0, 1) {
		count, exist := starMap[string(d.Format("1/2/06"))]
		if exist {
			data = append(data, count)
		} else {
			data = append(data, 0)
		}
	}
	return data
}

func histogramStarDates(list []*github.Stargazer) map[string]int {

	dupFrequency := make(map[string]int)

	for _, item := range list {
		// check if the item/element exist in the dupFrequency map
		dateStr := string(item.GetStarredAt().Format("1/2/06"))
		_, exist := dupFrequency[dateStr]

		if exist {
			dupFrequency[dateStr]++
		} else {
			dupFrequency[dateStr] = 1
		}
	}
	return dupFrequency
}

func main() {

	var err error
	var owner, repo string

	if len(os.Args) < 2 {
		log.Fatal(fmt.Errorf("please supply a repo in the format: owner/repo"))
	}

	repoParts := strings.Split(os.Args[1], "/")

	if len(repoParts) < 2 {
		log.Fatal(fmt.Errorf("please supply a repo in the format: owner/repo"))
	}

	owner = repoParts[0]
	repo = repoParts[1]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	viewsPerDay, _, err := client.Repositories.ListTrafficViews(ctx, owner, repo, &github.TrafficBreakdownOptions{Per: "day"})
	if err != nil {
		log.Fatal(err)
	}

	viewsPerWeek, _, err := client.Repositories.ListTrafficViews(ctx, owner, repo, &github.TrafficBreakdownOptions{Per: "week"})
	if err != nil {
		log.Fatal(err)
	}
	// get all star gazers
	var allStargazer []*github.Stargazer
	opt := &github.ListOptions{PerPage: 50}
	for {
		starGazers, resp, err := client.Activity.ListStargazers(ctx, owner, repo, opt)
		if err != nil {
			log.Fatal(err)
		}
		allStargazer = append(allStargazer, starGazers...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// fmt.Println("histogramStarDates: ", histogramStarDates(allStargazer))
	// fmt.Println("getStarDateSparklineData: ", getStarDateSparklineData(allStargazer[0].GetStarredAt().Time, histogramStarDates(allStargazer)))
	// fmt.Println("len(getStarDateSparklineData): ", len(getStarDateSparklineData(allStargazer[0].GetStarredAt().Time, histogramStarDates(allStargazer))))

	dayViews := []int{}
	dayPlotDateLabels := []string{}
	weekViews := []int{}
	weekPlotDateLabels := []string{}

	// fmt.Println("######  How Ya Doing  #############")
	// fmt.Println()
	// fmt.Println("View Count: ", viewsPerDay.GetCount())
	// fmt.Println("Unique Views: ", viewsPerDay.GetUniques())
	// fmt.Println()
	// fmt.Println("========= Break Down by Week =======")
	for _, tdata := range viewsPerWeek.Views {
		// fmt.Println("Week: ", tdata.GetTimestamp())
		// fmt.Println("View Count: ", tdata.GetCount())
		weekViews = append(weekViews, tdata.GetUniques())
		weekPlotDateLabels = append(weekPlotDateLabels, string(tdata.GetTimestamp().Format("1/2/06")))
		// fmt.Println("Unique Views: ", tdata.GetUniques())
		// fmt.Println("====================================")
	}
	// fmt.Println("========= Break Down by Day =======")
	for _, tdata := range viewsPerDay.Views {
		// 	fmt.Println("Day: ", tdata.GetTimestamp())
		// 	fmt.Println("View Count: ", tdata.GetCount())
		dayViews = append(dayViews, tdata.GetUniques())
		dayPlotDateLabels = append(dayPlotDateLabels, string(tdata.GetTimestamp().Format("1/2")))
		// fmt.Println("Unique Views: ", tdata.GetUniques())
		// fmt.Println("====================================")
	}

	// fmt.Println("dayViews: ", dayViews)
	// fmt.Println("weekViews: ", weekViews)

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
	weekWidget.PaddingLeft = 5
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
	// Star Widget
	starSparkline := termui.NewSparkline()
	starData := getStarDateSparklineData(
		allStargazer[0].GetStarredAt().Time,
		histogramStarDates(allStargazer),
	)
	// starSparkline.Data = starData
	starSparkline.Data = starData[len(starData)-110:]
	starSparkline.Title = "Last 100 days"
	starSparkline.Height = 7
	starSparkline.LineColor = termui.ColorYellow

	starWidget := termui.NewSparklines(starSparkline)
	starWidget.Height = 10
	starWidget.BorderLabel = "Stars"
	starWidget.PaddingLeft = 2

	// build
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(3, 0, overviewWidget, weekWidget),
			termui.NewCol(9, 0, dayWidget)),
		termui.NewRow(termui.NewCol(12, 0, starWidget)))

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
