// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong
// Author：Btcrobot

package model_test

import (
	. "model"
	"testing"
)

func TestNewRole(t *testing.T) {
	roleList, err := NewRole().FindAll()
	for _, tmpUser := range roleList {
		t.logger(tmpUser.Roleid)
		t.logger("===")
	}
	if err != nil {
		t.Fatal(err)
	}
}
