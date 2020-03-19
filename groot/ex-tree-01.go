// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build example

package main

import (
	"flag"
	"fmt"
	"log"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("groot: ")

	var (
		fname  = flag.String("fname", "event.root", "path to ROOT file to create")
		evtmax = flag.Int64("evtmax", -1, "number of events to generate")
	)

	flag.Parse()

	err := tree1(*fname, *evtmax)
	if err != nil {
		log.Fatalf("could not read ROOT file: %+v", err)
	}
}

type Event struct {
	I int64
	A Det
	B Det
}

type Det struct {
	E float64
	T float64
}

func tree1(fname string, evtmax int64) error {
	f, err := groot.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open ROOT file %q: %w", fname, err)
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		return fmt.Errorf("could not retrieve object: %w", err)
	}

	tree := obj.(rtree.Tree)

	var e Event
	rvars := []rtree.ScanVar{
		{Name: "evt_i", Value: &e.I},
		{Name: "evt_a_e", Value: &e.A.E},
		{Name: "evt_a_t", Value: &e.A.T},
		{Name: "evt_b_e", Value: &e.B.E},
		{Name: "evt_b_t", Value: &e.B.T},
	}

	sc, err := rtree.NewScannerVars(tree, rvars...)
	if err != nil {
		return fmt.Errorf("could not create tree scanner: %w", err)
	}
	defer sc.Close()

	if evtmax < 0 || evtmax > tree.Entries() {
		evtmax = tree.Entries()
	}

	var (
		h00 = hbook.NewH1D(100, -5, 5)
		h01 = hbook.NewH1D(100, -5, 5)
		h10 = hbook.NewH2D(100, -5, 5, 100, -5, 5)
		h20 = hbook.NewH1D(100, 0, 1)
		h21 = hbook.NewH1D(100, -0.5, 1.5)
		h30 = hbook.NewH2D(100, 0, 1, 100, -0.5, 1.5)
	)

	for sc.Next() && sc.Entry() < evtmax {
		iev := sc.Entry()
		if iev%1000 == 0 {
			log.Printf("processing event %d...", iev)
		}

		err = sc.Scan()
		if err != nil {
			return fmt.Errorf("could not read event #%d: %w", iev, err)
		}

		h00.Fill(e.A.E, 1)
		if -0.2 < e.B.E && e.B.E < 0.2 {
			h01.Fill(e.A.E, 3)
		}
		h10.Fill(e.A.E, e.B.E, 1)
		h20.Fill(e.A.T, 1)
		h21.Fill(e.B.T, 1)
		h30.Fill(e.A.T, e.B.T, 1)

		if iev%1000 == 0 {
			log.Printf("ievt: %d", iev)
			log.Printf("evt.a.e= %8.3f", e.A.E)
			log.Printf("evt.a.t= %8.3f", e.A.T)
			log.Printf("evt.b.e= %8.3f", e.B.E)
			log.Printf("evt.b.t= %8.3f", e.B.T)
		}
	}

	err = sc.Err()
	if err != nil {
		return fmt.Errorf("could not read tree: %w", err)
	}

	{
		tp := hplot.NewTiledPlot(draw.Tiles{
			Rows: 2, Cols: 2,
		})
		tp.Align = true

		p0 := tp.Plot(0, 0)
		p0.Title.Text = "A.E"
		p0.X.Label.Text = "A.E"
		hh00 := hplot.NewH1D(h00)
		hh00.LineStyle.Color = plotutil.SoftColors[0]
		hh01 := hplot.NewH1D(h01)
		hh01.LineStyle.Color = plotutil.SoftColors[2]
		p0.Add(hh00, hh01)
		p0.Add(hplot.NewGrid())
		p0.Legend.Add("A.E", hh00)
		p0.Legend.Add("A.E (w/ cut)", hh01)
		p0.Legend.Top = true

		p1 := tp.Plot(0, 1)
		p1.Title.Text = "B.E:A.E"
		p1.X.Label.Text = "A.E"
		p1.Y.Label.Text = "B.E"
		p1.Add(hplot.NewH2D(h10, nil))

		p2 := tp.Plot(1, 0)
		p2.Title.Text = "Time"
		p2.X.Label.Text = "time"
		hh20 := hplot.NewH1D(h20)
		hh20.LineStyle.Color = plotutil.SoftColors[0]
		hh21 := hplot.NewH1D(h21)
		hh21.LineStyle.Color = plotutil.SoftColors[2]
		p2.Add(hh20, hh21)
		p2.Add(hplot.NewGrid())
		p2.Legend.Add("A.T", hh20)
		p2.Legend.Add("B.T", hh21)
		p2.Legend.Top = true

		p3 := tp.Plot(1, 1)
		p3.Title.Text = "B.T:A.T"
		p3.X.Label.Text = "A.T"
		p3.Y.Label.Text = "B.T"
		p3.Add(hplot.NewH2D(h30, nil))

		err = tp.Save(20*vg.Centimeter, -1, "testdata/event.png")
		if err != nil {
			return fmt.Errorf("could not save plot: %w", err)
		}
	}
	return nil
}

func plot(fname string) error {
	f, err := groot.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open ROOT file: %w", err)
	}
	defer f.Close()

	return nil
}
