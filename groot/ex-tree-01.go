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

	for sc.Next() && sc.Entry() < evtmax {
		iev := sc.Entry()
		if iev%1000 == 0 {
			log.Printf("processing event %d...", iev)
		}

		err = sc.Scan()
		if err != nil {
			return fmt.Errorf("could not read event #%d: %w", iev, err)
		}

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

	return nil
}
