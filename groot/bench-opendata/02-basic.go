// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/vg"
)

// basic2 plots the Jet pT of all jets in an event.
func basic2(fname string) error {
	f, err := groot.Open(fname)
	if err != nil {
		return errors.Wrap(err, "could not open ROOT file")
	}
	defer f.Close()

	o, err := f.Get("Events")
	if err != nil {
		return errors.Wrap(err, "could not retrieve tree")
	}

	tree := o.(rtree.Tree)
	fmt.Printf("tree: %d entries\n", tree.Entries())

	var (
		hJetPt = hbook.NewH1D(100, 15, 60)
	)

	sc, err := rtree.NewTreeScannerVars(tree, rtree.ScanVar{Name: "Jet_pt"})
	if err != nil {
		return errors.Wrap(err, "could not create scanner")
	}
	defer sc.Close()

	for sc.Next() {
		var jetPt []float32
		err := sc.Scan(&jetPt)
		if err != nil {
			return errors.Wrap(err, "error during scan")
		}
		for _, pt := range jetPt {
			hJetPt.Fill(float64(pt), 1)
		}
	}

	if err := sc.Err(); err != nil {
		return errors.Wrap(err, "could not scan whole file")
	}

	fmt.Printf("hJetPt: %v\n", hJetPt.SumW())

	p := hplot.New()
	p.X.Label.Text = "Jet Pt [GeV]"
	p.Y.Label.Text = "Nevts"

	p.Add(hplot.NewH1D(hJetPt))

	err = p.Save(10*vg.Centimeter, -1, "02-basic.png")
	if err != nil {
		return errors.Wrap(err, "could not save plot")
	}

	return nil
}
