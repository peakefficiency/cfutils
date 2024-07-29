package main

import (
	"fmt"

	"github.com/peakefficiency/cfutils/cfaccess"
)

func main() {
	const appURL = "https://access-tester.pages.dev"

	response, err := cfaccess.GetWithAccess(appURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Response: %s\n", response)
}
