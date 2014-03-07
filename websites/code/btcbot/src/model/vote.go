// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package model

import (
	"logger"
	"util"
)

type Vote struct {
	Vid int    `json:"vid"`
	Uid int    `json:"uid"`
	IP  string `json:"ip"`
	Tid int    `json:"tid"`

	// 数据库访问对象
	*Dao
}

func NewVote() *Vote {
	return &Vote{
		Dao: &Dao{tablename: "vote"},
	}
}

func (this *Vote) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *Vote) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Vote) FindAll(selectCol ...string) ([]*Vote, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	voteList := make([]*Vote, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		vote := NewVote()
		err = this.Scan(rows, colNum, vote.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Vote FindAll Scan Error:", err)
			continue
		}
		voteList = append(voteList, vote)
	}
	return voteList, nil
}

func (this *Vote) prepareInsertData() {
	this.columns = []string{"uid", "ip", "tid"}
	this.colValues = []interface{}{this.Uid, this.IP, this.Tid}
}

func (this *Vote) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"vid": &this.Vid,
		"uid": &this.Uid,
		"ip":  &this.IP,
		"tid": &this.Tid,
	}
}
