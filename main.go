package main

import (
	"fmt"

	"github.com/mrjkey/aggregator/internal/config"
)

func main() {
	gator_config := config.Read()
	gator_config.SetUser("jared")
	gator_config_2 := config.Read()
	fmt.Println(gator_config_2)
}
