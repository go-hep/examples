// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"os/exec"
	"testing"
)

func TestReadWriteHBook(t *testing.T) {
	cmd := exec.Command("go", "run", "./hbook-w-hsimple.go")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("could not create hbook data: %+v", err)
	}

	cmd = exec.Command("go", "run", "./hbook-r-hsimple.go")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("could not read hbook data: %+v", err)
	}
}
