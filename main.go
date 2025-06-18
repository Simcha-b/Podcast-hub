package main

import (
	"fmt"
	"github.com/Simcha-b/Podcast-Hub/utils"
)

func main() {
	fmt.Println("welcome to the Podcast-Hub!!")

	logger := utils.NewLogger("info")
	logger.Info("This is an info message")
	logger.Error("This is an error message")
}
