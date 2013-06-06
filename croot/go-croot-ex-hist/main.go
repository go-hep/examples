package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/go-hep/croot"
)

var (
	fname    = flag.String("o", "hist.root", "path to the file holding the saved histos")
	nentries = flag.Int("nentries", 1000, "number of entries to generate")
)

func main() {
	flag.Parse()
	if *fname == "" {
		fmt.Fprintf(os.Stderr, "**error** you need to give a path to a ROOT file to save histograms into\n")
		os.Exit(1)
	}

	f, err := croot.OpenFile(*fname, "recreate", "go-croot-ex-hist", 1, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "**error** creating ROOT file [%s]: %v\n", *fname, err)
		os.Exit(1)
	}
	defer f.Close("")

	h1 := croot.NewH1F("h1", "h1", 100, -10, 10)
	if h1 == nil {
		fmt.Fprintf(os.Stderr, "**error** creating histogram [h1]\n")
		os.Exit(1)
	}

	h2 := croot.NewH1F("h2", "h2", 100, -10, 10)
	if h2 == nil {
		fmt.Fprintf(os.Stderr, "**error** creating histogram [h2]\n")
		os.Exit(1)
	}

	stddev_1 := 5.0
	mean_1 := 0.0

	stddev_2 := 3.0
	mean_2 := 2.5

	for i := 0; i < *nentries; i++ {
		do_display := i%(*nentries/10) == 0
		if do_display {
			fmt.Printf(":: filling entry %d...\n", i)
		}
		x_1 := rand.NormFloat64()*stddev_1 + mean_1
		x_2 := rand.NormFloat64()*stddev_2 + mean_2

		h1.Fill(x_1, 1)
		h2.Fill(x_2, 1)
	}

	for _, h := range []croot.H1F{h1, h2} {
		fmt.Printf(":: %s ===============\n", h.GetTitle())
		fmt.Printf("   Entries= %8.3f\n", h.GetEntries())
		fmt.Printf("   Mean=    %8.3f\n", h.GetMean())
		fmt.Printf("   RMS=     %8.3f\n", h.GetRMS())
	}

	if o := f.Write("", 0, 0); o < 0 {
		fmt.Fprintf(os.Stderr, "**error** problem committing to file: %d\n", o)
		os.Exit(1)
	}

	return
}

// EOF
