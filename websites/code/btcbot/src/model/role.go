// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package model

import (
	"logger"
	"time"
	"util"
)

// 角色分界点：roleid小于该值，则没有管理权限
const AdminMinRoleId int = 6

// 角色信息
type Role struct {
	Roleid int    `json:"roleid"`
	Name   string `json:"name"`
	ctime  time.Time

	// 数据库访问对象
	*Dao
}

func NewRole() *Role {
	return &Role{
		Dao: &Dao{tablename: "role"},
	}
}

func (this *Role) Insert() (int64, error) {
	this.prepareInsertData()
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *Role) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *Role) FindAll(selectCol ...string) ([]*Role, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	roleList := make([]*Role, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		role := NewRole()
		err = this.Scan(rows, colNum, role.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("Role FindAll Scan Error:", err)
			continue
		}
		roleList = append(roleList, role)
	}
	return roleList, nil
}

func (this *Role) prepareInsertData() {
	this.columns = []string{"name"}
	this.colValues = []interface{}{this.Name}
}

func (this *Role) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"roleid": &this.Roleid,
		"name":   &this.Name,
		"ctime":  &this.ctime,
	}
}
