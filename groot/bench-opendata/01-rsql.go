// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rsql"
	_ "go-hep.org/x/hep/groot/rsql/rsqldrv"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/vg"
)

// rsql1 plots the MET in an event.
func rsql1(fname string) error {
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

	hmet, err := rsql.ScanH1D(tree, "SELECT MET_sumet FROM Events", hbook.NewH1D(100, 0, 2000))
	if err != nil {
		return fmt.Errorf("could not scan tree: %w", err)
	}

	fmt.Printf("hmet: %v\n", hmet.SumW())

	p := hplot.New()
	p.X.Label.Text = "MET [GeV]"
	p.Y.Label.Text = "Nevts"

	p.Add(hplot.NewH1D(hmet))

	err = p.Save(10*vg.Centimeter, -1, "01-rsql.png")
	if err != nil {
		return fmt.Errorf("could not save plot: %w", err)
	}

	return nil
}
