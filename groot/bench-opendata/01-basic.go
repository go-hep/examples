// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/vg"
)

// basic1 plots the MET in an event.
func basic1(fname string) error {
	f, err := groot.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open ROOT file: %w", err)
	}
	defer f.Close()

	o, err := f.Get("Events")
	if err != nil {
		return fmt.Errorf("could not retrieve tree: %w", err)
	}

	tree := o.(rtree.Tree)
	fmt.Printf("tree: %d entries\n", tree.Entries())

	var (
		hmet = hbook.NewH1D(100, 0, 2000)
	)

	sc, err := rtree.NewTreeScannerVars(tree, rtree.ScanVar{Name: "MET_sumet"})
	if err != nil {
		return fmt.Errorf("could not create scanner: %w", err)
	}
	defer sc.Close()

	for sc.Next() {
		var met float32
		err := sc.Scan(&met)
		if err != nil {
			return fmt.Errorf("error during scan: %w", err)
		}

		hmet.Fill(float64(met), 1)
	}

	if err := sc.Err(); err != nil {
		return fmt.Errorf("could not scan whole file: %w", err)
	}

	fmt.Printf("hmet: %v\n", hmet.SumW())

	p := hplot.New()
	p.X.Label.Text = "MET [GeV]"
	p.Y.Label.Text = "Nevts"

	p.Add(hplot.NewH1D(hmet))

	err = p.Save(10*vg.Centimeter, -1, "01-basic.png")
	if err != nil {
		return fmt.Errorf("could not save plot: %w", err)
	}

	return nil
}
