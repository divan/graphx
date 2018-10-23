package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/divan/graphx/formats"
	"github.com/divan/graphx/graph"
	"github.com/divan/graphx/layout"
)

func main() {
	var (
		n               = flag.Int("n", 100, "Number of iterations to run physics simulation")
		input           = flag.String("i", "network.json", "File to read network graph layout from")
		t               = flag.String("type", "preset", "Export type [preset, json, ngraph]")
		verbose         = flag.Bool("v", false, "Be verbose (print forces and positions on each interation)")
		output          = flag.String("o", "positions.json", "Output file")
		repelCoeff      = flag.Float64("repel", -10.0, "Repelling force coefficent")
		springStiffness = flag.Float64("spring", 0.02, "Spring stiffness coefficient")
		springLen       = flag.Float64("springLen", 5.0, "Spring still length")
		dragCoeff       = flag.Float64("drag", 0.02, "Drag force coefficient")
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
	case "preset":
		var res struct {
			Description string             `json:"description"`
			Nodes       []*graph.BasicNode `json:"nodes"`
			Links       []*struct {
				Source string `json:"source"`
				Target string `json:"target"`
			} `json:"links"`
			Positions []*layout.Point `json:"positions"`
		}

		// read input file again
		fd, err := os.Open(*input)
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()

		err = json.NewDecoder(fd).Decode(&res)
		if err != nil {
			log.Fatal(err)
		}

		// add newly calculated positions
		res.Positions = positions

		// write preset to the output file
		fdOut, err := os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		defer fdOut.Close()
		enc := json.NewEncoder(fdOut)
		enc.SetIndent("", "  ")
		err = enc.Encode(res)
	default:
		err = fmt.Errorf("Unknown export format '%s'", *t)
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Written output to %s", *output)
}
