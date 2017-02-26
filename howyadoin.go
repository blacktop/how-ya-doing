package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	tbOptions := github.TrafficBreakdownOptions{Per: "week"}

	views, _, err := client.Repositories.ListTrafficViews(ctx, "maliceio", "malice", &tbOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("######  How Ya Doing  #############")
	fmt.Println()
	fmt.Println("View Count: ", views.GetCount())
	fmt.Println("Unique Views: ", views.GetUniques())
	fmt.Println()
	fmt.Println("========= Break Down by Week =======")
	for _, tdata := range views.Views {
		fmt.Println("Week: ", tdata.GetTimestamp())
		fmt.Println("View Count: ", tdata.GetCount())
		fmt.Println("Unique Views: ", tdata.GetUniques())
		fmt.Println("====================================")
	}
}
