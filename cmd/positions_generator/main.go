package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/divan/graphx/formats"
	"github.com/divan/graphx/layout"
)

func main() {
	var (
		n               = flag.Int("n", 100, "Number of iterations to run physics simulation")
		input           = flag.String("i", "network.json", "File to read network graph layout from")
		t               = flag.String("type", "ngraph", "Export type [json, ngraph]")
		verbose         = flag.Bool("v", false, "Be verbose (print forces and positions on each interation)")
		output          = flag.String("o", "positions.bin", "File to read network graph layout from")
		repelCoeff      = flag.Float64("repel", -10.0, "Repelling force coefficent")
		springStiffness = flag.Float64("spring", 0.02, "Spring stiffness coefficient")
		springLen       = flag.Float64("springLen", 5.0, "Spring still length")
		dragCoeff       = flag.Float64("drag", 0.8, "Drag force coefficient")
	)

	flag.Parse()

	g, err := formats.FromD3JSON(*input)
	if err != nil {
		log.Fatalf("Error reading network layout: %v", err)
	}

	config := layout.Config{
		Repelling:       *repelCoeff,
		SpringStiffness: *springStiffness,
		SpringLen:       *springLen,
		DragCoeff:       *dragCoeff,
	}
	l := layout.New(g, config)

	if *verbose {
		for i := 0; i < *n; i++ {
			for _, point := range l.PositionsSlice() {
				fmt.Printf("%s ", point)
			}
			fmt.Println()
			l.UpdatePositions()
		}
	} else {
		l.CalculateN(*n)
	}

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
