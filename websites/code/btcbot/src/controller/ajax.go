// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package controller

import (
	"encoding/json"
	"filter"
	"fmt"
	"github.com/studygolang/mux"
	"logger"
	"net/http"
	"service"
	"strconv"
	"strings"
)

// 侧边栏的内容通过异步请求获取

// 某节点下其他帖子
func OtherTopicsHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	topics := service.FindTopicsByNid(vars["nid"], vars["tid"])
	topics = service.JSEscape(topics)
	data, err := json.Marshal(topics)
	if err != nil {
		logger.Errorln("[OtherTopicsHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"errno": 1, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "topics":`+string(data)+`}`)
}

// 社区热门节点
// uri: /nodes/hot.json
func HotNodesHandler(rw http.ResponseWriter, req *http.Request) {
	nodes := service.FindHotNodes()
	hotNodes, err := json.Marshal(nodes)
	if err != nil {
		logger.Errorln("[HotNodesHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"errno": 1, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "nodes":`+string(hotNodes)+`}`)
}

func TopicVoteHandler(rw http.ResponseWriter, req *http.Request) {
	logger.Traceln(req.RemoteAddr)

	vars := mux.Vars(req)
	uid, _ := strconv.Atoi(Config["auid"])
	user, ok := filter.CurrentUser(req)
	if ok {
		uid = user["uid"].(int)
	}

	tid, _ := strconv.Atoi(vars["tid"])

	client := strings.Split(req.RemoteAddr, ":")

	logger.Traceln(tid, uid, client[0])
	if service.InsertVote(tid, uid, client[0]) == false {
		fmt.Fprint(rw, `{"errno": 1, "error":"已经投过票"}`)
		return
	}

	valInt, err := strconv.Atoi(vars["val"])
	if err != nil {
		fmt.Fprint(rw, `{"errno": 1, "error": "invalid request data"}`)
		return
	}

	if valInt == 0 {
		service.IncrTopicHate(vars["tid"], uid)
	} else {
		// 增加like量
		service.IncrTopicLike(vars["tid"], uid)
	}

	like, hate, err := service.FindTopicPopular(vars["tid"])
	logger.Traceln(like)
	logger.Traceln(hate)
	if err != nil {
		// TODO:
	}

	fmt.Fprint(rw, `{"errno": 0, "like":`+strconv.Itoa(like)+`, "hate":`+strconv.Itoa(hate)+`}`)
}
