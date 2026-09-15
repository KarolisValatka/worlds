package main

import (
	"fmt"
	"os"

	"github.com/KarolisValatka/worlds/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "worlds: %v\n", err)
		os.Exit(1)
	}
}
