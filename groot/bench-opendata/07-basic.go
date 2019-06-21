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

// basic7 plots the sum of the pt of all jets of pt > 30 GeV
// that are not within DR 0.4 from a lepton of pt > 10 GeV.
func basic7(fname string) error {
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
		h1 = hbook.NewH1D(100, 15, 200)
	)

	sc, err := rtree.NewTreeScannerVars(tree,
		rtree.ScanVar{Name: "Jet_pt"},
		rtree.ScanVar{Name: "Jet_eta"},
		rtree.ScanVar{Name: "Jet_phi"},
		rtree.ScanVar{Name: "Muon_pt"},
		rtree.ScanVar{Name: "Muon_eta"},
		rtree.ScanVar{Name: "Muon_phi"},
		rtree.ScanVar{Name: "Electron_pt"},
		rtree.ScanVar{Name: "Electron_eta"},
		rtree.ScanVar{Name: "Electron_phi"},
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
			jetPt  []float32
			jetEta []float32
			jetPhi []float32
			muPt   []float32
			muEta  []float32
			muPhi  []float32
			elePt  []float32
			eleEta []float32
			elePhi []float32
		)
		err := sc.Scan(&jetPt, &jetEta, &jetPhi, &muPt, &muEta, &muPhi, &elePt, &eleEta, &elePhi)
		if err != nil {
			return errors.Wrap(err, "error during scan")
		}

		njets := len(jetPt)
		if njets < 1 {
			continue
		}

		var (
			jets    []int
			muJets  = goodJets(jetPt, jetEta, jetPhi, muPt, muEta, muPhi)
			eleJets = goodJets(jetPt, jetEta, jetPhi, elePt, eleEta, elePhi)
		)
		switch {
		case len(muJets) > 0:
			jets = muJets
		case len(eleJets) > 0:
			jets = eleJets
		default:
			continue
		}

		pt := 0.0
		for _, i := range jets {
			pt += float64(jetPt[i])
		}
		h1.Fill(pt, 1)
	}

	if err := sc.Err(); err != nil {
		return errors.Wrap(err, "could not scan whole file")
	}

	fmt.Printf("h1: %v\n", h1.SumW())

	p := hplot.New()
	p.X.Label.Text = "Jet Pt sum [GeV]"
	p.Y.Label.Text = "Nevts"
	p.Add(hplot.NewH1D(h1))

	err = p.Save(10*vg.Centimeter, -1, "07-basic.png")
	if err != nil {
		return errors.Wrap(err, "could not save plot")
	}

	return nil
}

func goodJets(pt1, eta1, phi1, pt2, eta2, phi2 []float32) []int {
	const (
		twopi  = 2 * math.Pi
		dr2Min = 0.4 * 0.4
	)

	var jets = make([]int, 0, len(pt1))

	for ijet, jetPt := range pt1 {
		if jetPt <= 30 {
			continue
		}
		for ilep, lepPt := range pt2 {
			if lepPt <= 10 {
				continue
			}
			dphi := -math.Remainder(float64(phi1[ijet]-phi2[ilep]), twopi)
			deta := float64(eta1[ijet] - eta2[ilep])
			dr2 := dphi*dphi + deta*deta
			if dr2 > dr2Min {
				jets = append(jets, ijet)
			}
		}
	}

	return jets
}
