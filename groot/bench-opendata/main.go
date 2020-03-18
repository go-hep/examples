// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/pkg/profile"

	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
)

var (
	// IDs of OpenData benchmarks
	benchIDs = map[string]func(name string) error{
		"01-basic": basic1,
		"02-basic": basic2,
		"03-basic": basic3,
		"04-basic": basic4,
		"05-basic": basic5,
		"06-basic": basic6,
		"07-basic": basic7,
		"08-basic": basic8,
		"01-rsql":  rsql1,
	}
	benchNames []string
)

func main() {
	log.SetPrefix("bench-opendata: ")
	log.SetFlags(0)

	var (
		benchFlag = flag.String(
			"bench", "",
			"comma-separated list of opendata benchmark examples to run ("+strings.Join(benchNames, ",")+")",
		)
		fnameFlag = flag.String(
			"f", "root://eospublic.cern.ch//eos/root-eos/benchmark/Run2012B_SingleMu.root",
			"input file to analyze",
		)
		listFlag = flag.Bool("list", false, "list all available benchmarks and exits")
		profFlag = flag.Bool("profile", false, "enable/disable CPU profiling")
	)

	flag.Parse()

	if *listFlag {
		log.Printf("available OpenData benchmark examples: %q", benchNames)
		os.Exit(0)
	}

	if *profFlag {
		defer profile.Start(profile.CPUProfile).Stop()
	}

	var benchs []string

	switch *benchFlag {
	case "":
		benchs = benchNames
	default:
		names := strings.Split(*benchFlag, ",")
		for _, name := range names {
			benchs = append(benchs, strings.TrimSpace(name))
		}
	}

	log.Printf("running benchs: %q", benchs)

	allGood := true
	for _, name := range benchs {
		err := runBench(name, *fnameFlag)
		if err != nil {
			log.Printf("could not run bench %q: %v", name, err)
			allGood = false
		}
	}

	if !allGood {
		log.Fatalf("at least one benchmark failed")
	}
}

func runBench(id, fname string) error {
	fct, ok := benchIDs[id]
	if !ok {
		return fmt.Errorf("no such OpenData example: %q", id)
	}

	log.Printf("running %q...", id)
	beg := time.Now()
	err := fct(fname)
	end := time.Now()
	log.Printf("running %q... [err=%v] delta=%v", id, err, end.Sub(beg))

	if err != nil {
		return fmt.Errorf("could not run bench %q: %w", id, err)
	}

	return err
}

func init() {
	benchNames = make([]string, 0, len(benchIDs))
	for k := range benchIDs {
		benchNames = append(benchNames, k)
	}
	sort.Strings(benchNames)
}
