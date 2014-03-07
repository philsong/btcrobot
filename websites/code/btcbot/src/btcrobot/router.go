// Copyright 2014 The btcbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/
// Author：Phil	78623269@qq.com

package main

import (
	"config"
	. "controller"
	"controller/admin"
	"filter"
	"github.com/studygolang/mux"
)

func initRouter() *mux.Router {
	// 登录校验过滤器
	loginFilter := new(filter.LoginFilter)
	loginFilterChain := mux.NewFilterChain(loginFilter)

	router := mux.NewRouter()
	// 所有的页面都需要先检查用户cookie是否存在，以便在没登录时自动登录
	cookieFilter := new(filter.CookieFilter)
	// 大部分handler都需要页面展示
	frontViewFilter := filter.NewViewFilter()
	// 表单校验过滤器（配置了验证规则就会执行）
	formValidateFilter := new(filter.FormValidateFilter)
	router.FilterChain(mux.NewFilterChain([]mux.Filter{cookieFilter, formValidateFilter, frontViewFilter}...))

	router.HandleFunc("/", WelcomeHandler)
	router.HandleFunc("/licence{json:(|.json)}", LicenceHandler)
	router.HandleFunc("/secret{json:(|.json)}", SecretHandler)
	router.HandleFunc("/engine{msgtype:(|ajax|get|post)}{json:(|.json)}", EngineHandler)
	router.HandleFunc("/trade{msgtype:(|ajax|dobuy|dosell)}{json:(|.json)}", TradeHandler)
	router.HandleFunc("/remind{msgtype:(|ajax|get|post)}{json:(|.json)}", RemindHandler)

	router.HandleFunc("/topics{view:(|/popular|/no_reply|/last)}", TopicsBriefHandler)
	router.HandleFunc("/topicslist{view:(|/popular|/no_reply|/last)}", TopicsListHandler)
	router.HandleFunc("/topics/{tid:[0-9]+}", TopicDetailHandler)
	router.HandleFunc("/topics/new{json:(|.json)}", NewTopicHandler)
	router.HandleFunc("/topics/inspect{json:(|.json)}", InspectTopicHandler)
	router.HandleFunc("/topics/inspect/{tid:[0-9]+}/{val:[0-1]}", InspectTopicHandler)
	// 某个节点下的话题
	router.HandleFunc("/topics/node{nid:[0-9]+}", NodesHandler)

	// 注册
	router.HandleFunc("/account/register{json:(|.json)}", RegisterHandler)
	// 登录
	router.HandleFunc("/account/login", LoginHandler)
	router.HandleFunc("/account/logout", LogoutHandler)

	router.HandleFunc("/account/edit{json:(|.json)}", AccountEditHandler).AppendFilterChain(loginFilterChain)
	router.HandleFunc("/account/changepwd.json", ChangePwdHandler).AppendFilterChain(loginFilterChain)

	router.HandleFunc("/account/forgetpwd", ForgetPasswdHandler)
	router.HandleFunc("/account/resetpwd", ResetPasswdHandler)
	router.HandleFunc("/account/reminder{json:(|.json)}", ReminderHandler).AppendFilterChain(loginFilterChain)

	// 用户相关
	router.HandleFunc("/users", UsersHandler)
	router.HandleFunc("/user/{username:\\w+}", UserHomeHandler)

	// 评论
	router.HandleFunc("/comment/{objid:[0-9]+}.json", CommentHandler).AppendFilterChain(loginFilterChain)

	// 消息相关
	router.HandleFunc("/message/send{json:(|.json)}", SendMessageHandler).AppendFilterChain(loginFilterChain)
	router.HandleFunc("/message/{msgtype:(system|inbox|outbox)}", MessageHandler).AppendFilterChain(loginFilterChain)
	router.HandleFunc("/message/delete.json", DeleteMessageHandler).AppendFilterChain(loginFilterChain)

	/////////////////// 异步请求 开始///////////////////////
	// 某节点下其他帖子
	router.HandleFunc("/topics/others/{nid:[0-9]+}_{tid:[0-9]+}.json", OtherTopicsHandler)
	// 热门节点
	router.HandleFunc("/nodes/hot.json", HotNodesHandler)
	router.HandleFunc("/topics/vote.{tid:[0-9]+}_{val:[0-1]}.json", TopicVoteHandler)
	/////////////////// 异步请求 结束 ///////////////////////

	// 管理后台权限检查过滤器
	adminFilter := new(filter.AdminFilter)
	backViewFilter := filter.NewViewFilter(config.ROOT + "/template/admin/common.html")
	adminFilterChain := mux.NewFilterChain([]mux.Filter{loginFilter, adminFilter, formValidateFilter, backViewFilter}...)
	// admin 子系统
	// router.HandleFunc("/admin", admin.IndexHandler).AppendFilterChain(loginFilterChain) // 支持"/admin访问"
	subrouter := router.PathPrefix("/admin").Subrouter()
	// 所有后台需要的过滤器链
	subrouter.FilterChain(adminFilterChain)
	subrouter.HandleFunc("/", admin.IndexHandler)

	// 帖子管理
	subrouter.HandleFunc("/topics", admin.TopicsHandler)
	subrouter.HandleFunc("/nodes", admin.NodesHandler)

	// 用户管理
	subrouter.HandleFunc("/users", admin.UsersHandler)
	subrouter.HandleFunc("/newuser", admin.NewUserHandler)
	subrouter.HandleFunc("/adduser", admin.AddUserHandler)

	// 错误处理handler
	router.HandleFunc("/noauthorize", NoAuthorizeHandler) // 无权限handler
	// 404页面
	router.HandleFunc("/{*}", NotFoundHandler)

	return router
}
