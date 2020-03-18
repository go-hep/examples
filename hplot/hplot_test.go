// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"os/exec"
	"testing"
)

func TestHistogram(t *testing.T) {
	cmd := exec.Command("go", "run", "./ex-00-hist.go")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("could not run example: %+v", err)
	}
}
