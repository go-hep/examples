// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/stat/combin"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// basic8 runs the following analysis:
// in events with >=3 leptons and a same-flavour opposite-sign lepton pair,
// find the best same-flavour opposite-sign lepton pair (mass closest to 91.2 GeV),
// and plot the transverse mass of the missing energy and the leading other lepton
func basic8(fname string) error {
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
		hlep = hbook.NewH1D(100, 15, 60)
	)

	sc, err := rtree.NewTreeScannerVars(tree,
		rtree.ScanVar{Name: "Muon_pt"},
		rtree.ScanVar{Name: "Muon_eta"},
		rtree.ScanVar{Name: "Muon_phi"},
		rtree.ScanVar{Name: "Muon_mass"},
		rtree.ScanVar{Name: "Muon_charge"},
		rtree.ScanVar{Name: "Electron_pt"},
		rtree.ScanVar{Name: "Electron_eta"},
		rtree.ScanVar{Name: "Electron_phi"},
		rtree.ScanVar{Name: "Electron_mass"},
		rtree.ScanVar{Name: "Electron_charge"},
		rtree.ScanVar{Name: "MET_sumet"},
	)
	if err != nil {
		return fmt.Errorf("could not create scanner: %w", err)
	}
	defer sc.Close()

	for sc.Next() {
		if sc.Entry() > 1e5 {
			break
		}
		var (
			muPt      []float32
			muEta     []float32
			muPhi     []float32
			muMass    []float32
			muCharge  []int32
			elePt     []float32
			eleEta    []float32
			elePhi    []float32
			eleMass   []float32
			eleCharge []int32
			sumet     float32
		)
		err := sc.Scan(
			&muPt, &muEta, &muPhi, &muMass, &muCharge,
			&elePt, &eleEta, &elePhi, &eleMass, &eleCharge,
			&sumet,
		)
		if err != nil {
			return fmt.Errorf("error during scan: %w", err)
		}

		nleptons := len(muPt) + len(elePt)
		if nleptons < 3 {
			continue
		}

		var (
			imu1, imu2   = findLeptonPair(muPt, muEta, muPhi, muMass, muCharge)
			iele1, iele2 = findLeptonPair(elePt, eleEta, elePhi, eleMass, eleCharge)
		)

		if imu1 < 0 && iele1 < 0 {
			continue
		}

		var lepPt float32
		if imu1 >= 0 {
			for i, pt := range muPt {
				if i == imu1 || i == imu2 {
					continue
				}
				if pt > lepPt {
					lepPt = pt
				}
			}
			for _, pt := range elePt {
				if pt > lepPt {
					lepPt = pt
				}
			}
		}
		if iele1 >= 0 {
			for i, pt := range elePt {
				if i == iele1 || i == iele2 {
					continue
				}
				if pt > lepPt {
					lepPt = pt
				}
			}
			for _, pt := range muPt {
				if pt > lepPt {
					lepPt = pt
				}
			}
		}

		hmet.Fill(float64(sumet), 1)
		hlep.Fill(float64(lepPt), 1)
	}

	if err := sc.Err(); err != nil {
		return fmt.Errorf("could not scan whole file: %w", err)
	}

	fmt.Printf("hmet: %v\n", hmet.SumW())
	fmt.Printf("hlep: %v\n", hlep.SumW())

	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 1, Rows: 2})

	p1 := tp.Plots[0]
	p1.X.Label.Text = "MET [GeV]"
	p1.Y.Label.Text = "Nevts"
	p1.Add(hplot.NewH1D(hmet))

	p2 := tp.Plots[1]
	p2.X.Label.Text = "Lepton Pt [GeV]"
	p2.Y.Label.Text = "Nevts"
	p2.Add(hplot.NewH1D(hlep))

	err = tp.Save(10*vg.Centimeter, -1, "08-basic.png")
	if err != nil {
		return fmt.Errorf("could not save plot: %w", err)
	}

	return nil
}

func findLeptonPair(pt, eta, phi, mass []float32, charge []int32) (int, int) {
	const (
		zMass = 91.2 // or take it from go-hep.org/x/hep/heppdt.PDT[id].Particle.Mass
	)

	cand := struct {
		m      float64
		i1, i2 int
	}{
		m:  math.MaxFloat64,
		i1: -1,
		i2: -1,
	}

	if len(pt) < 2 {
		return cand.i1, cand.i2
	}

	makePtEtaPhiM := func(i int) fmom.PtEtaPhiM {
		return fmom.NewPtEtaPhiM(float64(pt[i]), float64(eta[i]), float64(phi[i]), float64(mass[i]))
	}

	combs := combin.Combinations(len(pt), 2)
	for _, c := range combs {
		i1 := c[0]
		i2 := c[1]
		if charge[i1] == charge[i2] {
			continue
		}
		p1 := makePtEtaPhiM(i1)
		p2 := makePtEtaPhiM(i2)
		mll := fmom.InvMass(&p1, &p2)
		if math.Abs(mll-zMass) < math.Abs(cand.m-zMass) {
			cand.m = mll
			cand.i1 = i1
			cand.i2 = i2
		}
	}

	return cand.i1, cand.i2
}
