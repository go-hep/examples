// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/stat/combin"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// basic6 plots the pt of the tri-jet system with mass closest to 172.5 GeV,
// and the leading b-tag discriminator among the 3 jets in the triplet.
func basic6(fname string) error {
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
		h1 = hbook.NewH1D(100, 15, 40)
		h2 = hbook.NewH1D(100, 0, 1)
	)

	sc, err := rtree.NewTreeScannerVars(tree,
		rtree.ScanVar{Name: "Jet_pt"},
		rtree.ScanVar{Name: "Jet_eta"},
		rtree.ScanVar{Name: "Jet_phi"},
		rtree.ScanVar{Name: "Jet_mass"},
		rtree.ScanVar{Name: "Jet_btag"},
	)
	if err != nil {
		return errors.Wrap(err, "could not create scanner")
	}
	defer sc.Close()

	for sc.Next() {
		if sc.Entry() > 1e5 {
			break
		}
		var (
			jetPt   []float32
			jetEta  []float32
			jetPhi  []float32
			jetMass []float32
			jetBtag []float32
		)
		err := sc.Scan(&jetPt, &jetEta, &jetPhi, &jetMass, &jetBtag)
		if err != nil {
			return errors.Wrap(err, "error during scan")
		}

		njets := len(jetPt)
		if njets < 3 {
			continue
		}

		idx := findTriJet(jetPt, jetEta, jetPhi, jetMass)
		btag := 0.0
		for _, i := range idx {
			h1.Fill(float64(jetPt[i]), 1)
			if v := float64(jetBtag[i]); v > btag {
				btag = v
			}
		}
		h2.Fill(btag, 1)
	}

	if err := sc.Err(); err != nil {
		return errors.Wrap(err, "could not scan whole file")
	}

	fmt.Printf("h1: %v\n", h1.SumW())
	fmt.Printf("h2: %v\n", h2.SumW())

	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 1, Rows: 2})

	p1 := tp.Plots[0]
	p1.X.Label.Text = "Trijet Pt [GeV]"
	p1.Y.Label.Text = "Nevts"
	p1.Add(hplot.NewH1D(h1))

	p2 := tp.Plots[1]
	p2.X.Label.Text = "Trijet leading b-tag"
	p2.Y.Label.Text = "Nevts"
	p2.Add(hplot.NewH1D(h2))

	err = tp.Save(10*vg.Centimeter, -1, "06-basic.png")
	if err != nil {
		return errors.Wrap(err, "could not save plot")
	}

	return nil
}

func findTriJet(pt, eta, phi, mass []float32) [3]int {
	makePtEtaPhiM := func(i int) *fmom.PtEtaPhiM {
		p := fmom.NewPtEtaPhiM(float64(pt[i]), float64(eta[i]), float64(phi[i]), float64(mass[i]))
		return &p
	}

	var (
		combs = combin.Combinations(len(pt), 3)
		delta = math.MaxFloat64
	)

	const topMass = 172.5 // could use go-hep.org/x/hep/heppdt.Particle.Mass, though.

	idx := 0
	for i, c := range combs {
		p1 := makePtEtaPhiM(c[0])
		p2 := makePtEtaPhiM(c[1])
		p3 := makePtEtaPhiM(c[2])

		m := fmom.Add(fmom.Add(p1, p2), p3).M()
		if d := math.Abs(m - topMass); d < delta {
			delta = d
			idx = i
		}
	}

	tri := combs[idx]
	return [3]int{tri[0], tri[1], tri[2]}
}
