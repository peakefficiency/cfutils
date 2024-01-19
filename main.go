package main

import (
	"fmt"

	"github.com/peakefficency/cf-zt-utils/cfaccess"
)

func main() {
	const appURL = "example.com"

	response, err := cfaccess.GetWithAccess(appURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Response: %s\n", response)
}
