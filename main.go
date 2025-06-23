package main

import (
	"fmt"

	// "github.com/Simcha-b/Podcast-Hub/utils"
	"github.com/Simcha-b/Podcast-Hub/services"
)

func main() {
	fmt.Println("welcome to the Podcast-Hub!!")

	//  Logger := utils.NewLogger("info")

	// Example usage of AggregateAllFeeds
	services.RunAggregator()
	
	
}
