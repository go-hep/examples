// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build example

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("groot: ")

	var (
		fname  = flag.String("fname", "event.root", "path to ROOT file to create")
		evtmax = flag.Int64("evtmax", 10000, "number of events to generate")
	)

	flag.Parse()

	err := tree0(*fname, *evtmax)
	if err != nil {
		log.Fatalf("could not create ROOT file: %+v", err)
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

func tree0(fname string, evtmax int64) error {
	f, err := groot.Create(fname)
	if err != nil {
		return fmt.Errorf("could not create ROOT file %q: %w", fname, err)
	}
	defer f.Close()

	var e Event
	wvars := []rtree.WriteVar{
		{Name: "evt_i", Value: &e.I},
		{Name: "evt_a_e", Value: &e.A.E},
		{Name: "evt_a_t", Value: &e.A.T},
		{Name: "evt_b_e", Value: &e.B.E},
		{Name: "evt_b_t", Value: &e.B.T},
	}

	tree, err := rtree.NewWriter(f, "tree", wvars)
	if err != nil {
		return fmt.Errorf("could not create tree writer: %w", err)
	}

	log.Printf("-- created tree %q:", tree.Name())
	for i, b := range tree.Branches() {
		log.Printf("branch[%d]: name=%q, title=%q", i, b.Name(), b.Title())
	}

	for i := int64(0); i < evtmax; i++ {
		if i%1000 == 0 {
			log.Printf("processing event %d...", i)
		}
		e.I = i
		e.A.E = rand.NormFloat64()
		e.B.E = rand.NormFloat64()

		e.A.T = rand.Float64()
		e.B.T = e.A.T * rand.NormFloat64()

		if i%1000 == 0 {
			log.Printf("evt.i=   %8d", e.I)
			log.Printf("evt.a.e= %8.3f", e.A.E)
			log.Printf("evt.a.t= %8.3f", e.A.T)
			log.Printf("evt.b.e= %8.3f", e.B.E)
			log.Printf("evt.b.t= %8.3f", e.B.T)
		}
		_, err = tree.Write()
		if err != nil {
			return fmt.Errorf("could not write event %d: %w", i, err)
		}
	}

	err = tree.Close()
	if err != nil {
		return fmt.Errorf("could not close tree writer: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("could not close ROOT file: %w", err)
	}

	return nil
}
