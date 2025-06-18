package main

import (
	"fmt"
	"os"

	"github.com/osm/mvdpl/internal/mvdparser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <demo.mvd>\n", os.Args[0])
		os.Exit(1)
	}

	mvdPath := os.Args[1]
	mvdData, err := os.ReadFile(mvdPath)
	if err != nil {
		fmt.Printf("unable to open %v, %v\n", mvdPath, err)
		os.Exit(1)
	}

	p := mvdparser.New()
	packetLoss, err := p.Parse(mvdData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error when parsing %v, %v\n", mvdPath, err)
		os.Exit(1)
	}

	for _, pl := range packetLoss {
		fmt.Printf("%s %s %d%% pl\n",
			pl.Name,
			fmt.Sprintf("%02d:%02d", pl.Timestamp/60, pl.Timestamp%60),
			pl.Lossage)
	}
}
