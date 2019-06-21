// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/vg"
)

// basic4 plots the missing ET of events with at least 2 jets above 40 GeV.
func basic4(fname string) error {
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
		hmet = hbook.NewH1D(100, 0, 2000)
	)

	sc, err := rtree.NewTreeScannerVars(tree,
		rtree.ScanVar{Name: "Jet_pt"},
		rtree.ScanVar{Name: "Jet_eta"},
		rtree.ScanVar{Name: "MET_sumet"},
	)
	if err != nil {
		return errors.Wrap(err, "could not create scanner")
	}
	defer sc.Close()

	for sc.Next() {
		var (
			jetPt  []float32
			jetEta []float32
			met    float32
		)
		err := sc.Scan(&jetPt, &jetEta, &met)
		if err != nil {
			return errors.Wrap(err, "error during scan")
		}
		njets := 0
	loop:
		for i, pt := range jetPt {
			if pt > 40 && math.Abs(float64(jetEta[i])) < 1 {
				njets++
				if njets > 1 {
					break loop
				}
			}
		}
		if njets >= 2 {
			hmet.Fill(float64(met), 1)
		}
	}

	if err := sc.Err(); err != nil {
		return errors.Wrap(err, "could not scan whole file")
	}

	fmt.Printf("hmet: %v\n", hmet.SumW())

	p := hplot.New()
	p.X.Label.Text = "MET [GeV]"
	p.Y.Label.Text = "Nevts"

	p.Add(hplot.NewH1D(hmet))

	err = p.Save(10*vg.Centimeter, -1, "04-basic.png")
	if err != nil {
		return errors.Wrap(err, "could not save plot")
	}

	return nil
}
