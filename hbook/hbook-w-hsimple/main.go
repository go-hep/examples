package main

import (
	"compress/flate"
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/go-hep/hbook"
	"github.com/go-hep/rio"
)

var (
	fname = flag.String("o", "hsimple.rio", "output file name")
)

func main() {
	flag.Parse()

	fmt.Printf(":: creating output file [%s]...\n", *fname)
	f, err := os.Create(*fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w, err := rio.NewWriter(f)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	err = w.SetCompressor(rio.CompressZlib, flate.DefaultCompression)
	if err != nil {
		panic(err)
	}

	h1 := hbook.NewH1D(100, -10., 10.)
	if h1 == nil {
		panic(fmt.Errorf("nil pointer to H1D"))
	}

	h1.Annotation()["name"] = "h1"

	h2 := hbook.NewH1D(100, -10., 10.)
	if h2 == nil {
		panic(fmt.Errorf("nil pointer to H1D"))
	}

	h2.Annotation()["name"] = "h2"

	// draw random values from a gaussian distribution
	const n = 1000
	fmt.Printf(":: filling histos with %d events...\n", n)
	rand.Seed(0)
	for i := 0; i < n; i++ {
		v1 := rand.NormFloat64()
		h1.Fill(v1, 1.0)

		v2 := rand.NormFloat64()
		h2.Fill(v2+2, 1.0)
	}
	fmt.Printf("%s: entries=%v mean=%+8.3f RMS=%+8.3f\n",
		h1.Name(), h1.Entries(), h1.Mean(), h1.RMS(),
	)
	fmt.Printf("%s: entries=%v mean=%+8.3f RMS=%+8.3f\n",
		h2.Name(), h2.Entries(), h2.Mean(), h2.RMS(),
	)

	fmt.Printf(":: saving histos to [%s]...\n", *fname)
	err = w.WriteValue("h1", h1)
	if err != nil {
		panic(err)
	}

	err = w.WriteValue("h2", h2)
	if err != nil {
		panic(err)
	}

	err = w.Close()
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

}
