// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build example

package main

import (
	"flag"
	"log"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/hbook"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("groot: ")

	var (
		fname    = flag.String("o", "hist.root", "path to output ROOT file")
		nentries = flag.Int("n", 1000, "number of entries to generate")
	)

	flag.Parse()
	if *fname == "" {
		flag.Usage()
		log.Fatalf("missing output ROOT file")
	}

	f, err := groot.Create(*fname)
	if err != nil {
		log.Fatalf("could not create ROOT file: %+v", err)
	}
	defer f.Close()

	h1 := hbook.NewH1D(100, -10, 10)
	h1.Annotation()["name"] = "h1"
	h1.Annotation()["title"] = "h1 - math/rand"

	h2 := hbook.NewH1D(100, -10, 10)
	h2.Annotation()["name"] = "h2"
	h2.Annotation()["title"] = "h2 - gonum/distuv"

	const (
		stddev = 5
		mean   = 0
	)

	dist := distuv.Normal{
		Mu:    mean,
		Sigma: stddev,
		Src:   rand.New(rand.NewSource(0)),
	}

	for i := 0; i < *nentries; i++ {
		verbose := i%(*nentries/10) == 0
		if verbose {
			log.Printf("filling entry %d...", i)
		}
		x1 := rand.NormFloat64()*stddev + mean
		x2 := dist.Rand()

		h1.Fill(x1, 1)
		h2.Fill(x2, 1)
	}

	for _, h := range []*hbook.H1D{h1, h2} {
		log.Printf("===== %-15s =====", h.Ann["title"])
		log.Printf("   Entries= %8.3f", h.EffEntries())
		log.Printf("   Mean=    %8.3f", h.XMean())
		log.Printf("   std-dev= %8.3f", h.XStdDev())

		err := f.Put(h.Name(), rhist.NewH1DFrom(h))
		if err != nil {
			log.Fatalf("could not store %q to ROOT file: %+v", h.Name(), err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could not close ROOT file: %w", err)
	}
}
