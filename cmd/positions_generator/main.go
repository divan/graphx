package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/divan/graphx/formats"
	"github.com/divan/graphx/layout"
)

func main() {
	var (
		n      = flag.Int("n", 100, "Number of iterations to run physics simulation")
		input  = flag.String("i", "network.json", "File to read network graph layout from")
		t      = flag.String("type", "ngraph", "Export type [json, ngraph]")
		output = flag.String("o", "positions.bin", "File to read network graph layout from")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -i network.json [-type json -o positions.json]", os.Args[0])
	}
	flag.Parse()

	g, err := formats.FromD3JSON(*input)
	if err != nil {
		log.Fatalf("Error reading network layout: %v", err)
	}

	l := layout.NewAuto(g)
	l.CalculateN(*n)

	positions := l.PositionsSlice()
	switch *t {
	case "ngraph":
		err = formats.ToPositionsNGraphFile(positions, *output)
	case "json":
		err = formats.ToPositionsJSONFile(positions, *output)
	default:
		err = fmt.Errorf("Unknown export format '%s'", *t)
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Written output to %s", *output)
}
