package service

import (
	"logger"
	"model"
	"strconv"
)

func FindVote(tid int, uid int, ip string) (err error) {
	condition := "tid=" + strconv.Itoa(tid) + " and uid=" + strconv.Itoa(uid) + " and ip=" + ip
	// 帖子信息
	topic := model.NewTopic()
	err = topic.Where(condition).Find()
	if err != nil {
		logger.Errorln("topic service FindTopicByTid Error:", err)
		return
	}
	if topic.Tid == 0 {
		return
	}

	logger.Traceln(topic)

	return
}

func InsertVote(tid int, uid int, ip string) bool {
	vote := model.NewVote()
	vote.Tid = tid
	vote.Uid = uid
	vote.IP = ip

	if _, err := vote.Insert(); err != nil {
		logger.Errorln("message service InsertVote Error:", err)
		return false
	}

	return true
}
