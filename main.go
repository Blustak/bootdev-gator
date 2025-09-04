package main

import (
	"fmt"

	"github.com/Blustak/bootdev-gator/internal/config"
)

func main() {
    userConfig, err := config.ReadUserConfig()
    if err != nil {
        fmt.Printf("error reading config file: %v\n", err)
        return
    }
}
