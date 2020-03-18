// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/stat/combin"
	"gonum.org/v1/plot/vg"
)

// basic5 plots the MEt for events that have an opposite-sign muon pair of mass 60-120 GeV.
func basic5(fname string) error {
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

	sc, err := rtree.NewTreeScannerVars(tree,
		rtree.ScanVar{Name: "Muon_pt"},
		rtree.ScanVar{Name: "Muon_eta"},
		rtree.ScanVar{Name: "Muon_phi"},
		rtree.ScanVar{Name: "Muon_mass"},
		rtree.ScanVar{Name: "Muon_charge"},
		rtree.ScanVar{Name: "MET_sumet"},
	)
	if err != nil {
		return fmt.Errorf("could not create scanner: %w", err)
	}
	defer sc.Close()

	for sc.Next() {
		var (
			muPt     []float32
			muEta    []float32
			muPhi    []float32
			muMass   []float32
			muCharge []int32
			met      float32
		)
		err := sc.Scan(&muPt, &muEta, &muPhi, &muMass, &muCharge, &met)
		if err != nil {
			return fmt.Errorf("error during scan: %w", err)
		}

		nmuons := len(muPt)
		if nmuons < 2 {
			continue
		}

		masses := make([]float64, 0, nmuons)
		combs := combin.Combinations(nmuons, 2)
		for _, c := range combs {
			i1 := c[0]
			i2 := c[1]
			charge1 := muCharge[i1]
			charge2 := muCharge[i2]
			if charge1 == charge2 {
				continue
			}

			p1 := fmom.NewPtEtaPhiM(float64(muPt[i1]), float64(muEta[i1]), float64(muPhi[i1]), float64(muMass[i1]))
			p2 := fmom.NewPtEtaPhiM(float64(muPt[i2]), float64(muEta[i2]), float64(muPhi[i2]), float64(muMass[i2]))
			mass := fmom.InvMass(&p1, &p2)

			if 60 < mass && mass < 100 {
				masses = append(masses, mass)
			}
		}

		if len(masses) > 0 {
			hmet.Fill(float64(met), 1)
		}
	}

	if err := sc.Err(); err != nil {
		return fmt.Errorf("could not scan whole file: %w", err)
	}

	fmt.Printf("hmet: %v\n", hmet.SumW())

	p := hplot.New()
	p.X.Label.Text = "MET [GeV]"
	p.Y.Label.Text = "Nevts"

	p.Add(hplot.NewH1D(hmet))

	err = p.Save(10*vg.Centimeter, -1, "05-basic.png")
	if err != nil {
		return fmt.Errorf("could not save plot: %w", err)
	}

	return nil
}
