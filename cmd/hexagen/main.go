package main

import (
	"encoding/hex"
	"fmt"
	"image/png"
	"os"
	"strings"

	"go.cryptoscope.co/hexagen/v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <hex>\n", os.Args[0])
		return
	}

	data, err := hex.DecodeString(os.Args[1])
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
		return
	}

	g := hexagen.Generate(data, 512)

	// replace slashes, they are not allowed in filesystem context
	f, err := os.Create(strings.Replace(os.Args[1], "/", "|", -1) + ".png")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
		return
	}
	defer f.Close()

	if err := png.Encode(f, g); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
		return
	}

}
