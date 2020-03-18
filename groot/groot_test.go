// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package groot_test

import (
	"os/exec"
	"testing"
)

func TestWriteH1D(t *testing.T) {
	cmd := exec.Command("go", "run", "./ex-hist-00.go")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("could not create histos: %+v", err)
	}
}

func TestRWFlatTree(t *testing.T) {
	cmd := exec.Command("go", "run", "./ex-tree-00.go")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("could not create ROOT tree: %+v", err)
	}

	cmd = exec.Command("go", "run", "./ex-tree-01.go")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("could not read ROOT tree: %+v", err)
	}
}
