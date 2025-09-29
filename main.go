package main

import (
	"fmt"
	"os"

	"github.com/osm/mvdpl/internal/fileutil"
	"github.com/osm/mvdpl/internal/format"
	"github.com/osm/mvdpl/internal/mvdparser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <demo.mvd>\n", os.Args[0])
		os.Exit(1)
	}

	mvdPath := os.Args[1]
	mvdData, err := fileutil.ReadMVD(mvdPath)
	if err != nil {
		fmt.Printf("unable to open %v, %v\n", mvdPath, err)
		os.Exit(1)
	}

	p := mvdparser.New()
	events, err := p.Parse(mvdData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error when parsing %v, %v\n", mvdPath, err)
		os.Exit(1)
	}

	for _, e := range events {
		fmt.Printf("[%s] %s %d%s\n",
			format.Time(e.Timestamp()), e.Name(), e.Value(), e.Suffix())
	}
}
