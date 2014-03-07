-- phpMyAdmin SQL Dump
-- version 4.0.9
-- http://www.phpmyadmin.net
--
-- 主机: 127.0.0.1
-- 生成日期: 2014-01-21 08:41:22
-- 服务器版本: 5.5.34
-- PHP 版本: 5.4.22

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;

--
-- 数据库: `btcrobot`
--

-- --------------------------------------------------------

--
-- 表的结构 `comments`
--

CREATE TABLE IF NOT EXISTS `comments` (
  `cid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `objid` int(10) unsigned NOT NULL COMMENT '对象id，属主（评论给谁）',
  `objtype` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '类型,0-帖子;1-博客;2-资源;3-酷站',
  `content` varchar(1024) NOT NULL,
  `uid` int(10) unsigned NOT NULL COMMENT '回复者',
  `floor` int(10) unsigned NOT NULL COMMENT '第几楼',
  `flag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '审核标识,0-未审核;1-已审核;2-审核删除;3-用户自己删除',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`cid`),
  UNIQUE KEY `objid` (`objid`,`objtype`,`floor`),
  KEY `uid` (`uid`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=27 ;

-- --------------------------------------------------------

--
-- 表的结构 `message`
--

CREATE TABLE IF NOT EXISTS `message` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `content` text NOT NULL COMMENT '消息内容',
  `hasread` enum('未读','已读') NOT NULL DEFAULT '未读',
  `from` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '来自谁',
  `fdel` enum('未删','已删') NOT NULL DEFAULT '未删' COMMENT '发送方删除标识',
  `to` int(10) unsigned NOT NULL COMMENT '发给谁',
  `tdel` enum('未删','已删') NOT NULL DEFAULT '未删' COMMENT '接收方删除标识',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `to` (`to`),
  KEY `from` (`from`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='message 短消息（私信）' AUTO_INCREMENT=4 ;

-- --------------------------------------------------------

--
-- 表的结构 `role`
--

CREATE TABLE IF NOT EXISTS `role` (
  `roleid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '角色名',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`roleid`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=10 ;

-- --------------------------------------------------------

--
-- 表的结构 `system_message`
--

CREATE TABLE IF NOT EXISTS `system_message` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `msgtype` tinyint(4) NOT NULL DEFAULT '0' COMMENT '系统消息类型',
  `hasread` enum('未读','已读') NOT NULL DEFAULT '未读',
  `to` int(10) unsigned NOT NULL COMMENT '发给谁',
  `ext` text NOT NULL COMMENT '额外信息',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `to` (`to`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='system_message 系统消息表' AUTO_INCREMENT=25 ;

-- --------------------------------------------------------

--
-- 表的结构 `topics`
--

CREATE TABLE IF NOT EXISTS `topics` (
  `tid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL DEFAULT '',
  `content` varchar(1024) NOT NULL,
  `nid` int(10) unsigned zerofill NOT NULL COMMENT '节点id',
  `uid` int(10) unsigned NOT NULL COMMENT '帖子作者',
  `lastreplyuid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最后回复者',
  `lastreplytime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '最后回复时间',
  `view` int(10) unsigned NOT NULL DEFAULT '0',
  `reply` int(10) unsigned NOT NULL DEFAULT '0',
  `like` int(10) unsigned NOT NULL DEFAULT '0',
  `hate` int(10) unsigned NOT NULL DEFAULT '0',
  `flag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '审核标识,0-未审核;1-已审核;2-审核删除;3-用户自己删除',
  `ctime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`tid`),
  KEY `uid` (`uid`),
  KEY `nid` (`nid`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=90 ;

-- --------------------------------------------------------

--
-- 表的结构 `topics_node`
--

CREATE TABLE IF NOT EXISTS `topics_node` (
  `nid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `parent` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '父节点id，无父节点为0',
  `name` varchar(20) NOT NULL COMMENT '节点名',
  `intro` varchar(50) NOT NULL DEFAULT '' COMMENT '节点简介',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`nid`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=8 ;

-- --------------------------------------------------------

--
-- 表的结构 `user_active`
--

CREATE TABLE IF NOT EXISTS `user_active` (
  `uid` int(10) unsigned NOT NULL,
  `email` varchar(128) NOT NULL,
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `weight` smallint(6) NOT NULL DEFAULT '1' COMMENT '活跃度，越大越活跃',
  `avatar` varchar(128) NOT NULL DEFAULT '' COMMENT '头像(暂时使用http://www.gravatar.com)',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- 表的结构 `user_info`
--

CREATE TABLE IF NOT EXISTS `user_info` (
  `uid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(128) NOT NULL DEFAULT '',
  `open` tinyint(4) NOT NULL DEFAULT '1' COMMENT '邮箱是否公开，默认公开',
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT '姓名',
  `avatar` varchar(128) NOT NULL DEFAULT '' COMMENT '头像(暂时使用http://www.gravatar.com)',
  `city` varchar(10) NOT NULL DEFAULT '',
  `company` varchar(64) NOT NULL DEFAULT '',
  `github` varchar(20) NOT NULL DEFAULT '',
  `weibo` varchar(20) NOT NULL DEFAULT '',
  `website` varchar(50) NOT NULL DEFAULT '' COMMENT '个人主页，博客',
  `status` varchar(140) NOT NULL DEFAULT '' COMMENT '个人状态，签名',
  `introduce` text NOT NULL COMMENT '个人简介',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `lowestPrice` float NOT NULL,
  `highestPrice` float NOT NULL,
  `emailnotice` tinyint(1) NOT NULL,
  `noticeTime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=10001 ;

-- --------------------------------------------------------

--
-- 表的结构 `user_login`
--

CREATE TABLE IF NOT EXISTS `user_login` (
  `uid` int(10) unsigned NOT NULL,
  `email` varchar(128) NOT NULL DEFAULT '',
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `passcode` char(12) NOT NULL DEFAULT '' COMMENT '加密随机数',
  `passwd` char(32) NOT NULL DEFAULT '' COMMENT 'md5密码',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- 表的结构 `user_role`
--

CREATE TABLE IF NOT EXISTS `user_role` (
  `uid` int(10) unsigned NOT NULL,
  `roleid` int(10) unsigned NOT NULL,
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`,`roleid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- 表的结构 `views`
--

CREATE TABLE IF NOT EXISTS `views` (
  `cid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `objid` int(10) unsigned NOT NULL COMMENT '对象id，属主（评论给谁）',
  `objtype` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '类型,0-帖子;1-博客;2-资源;3-酷站',
  `content` text NOT NULL,
  `uid` int(10) unsigned NOT NULL COMMENT '回复者',
  `floor` int(10) unsigned NOT NULL COMMENT '第几楼',
  `flag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '审核标识,0-未审核;1-已审核;2-审核删除;3-用户自己删除',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`cid`),
  UNIQUE KEY `objid` (`objid`,`objtype`,`floor`),
  KEY `uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

-- --------------------------------------------------------

--
-- 表的结构 `vote`
--

CREATE TABLE IF NOT EXISTS `vote` (
  `vid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '绑定的第三方类型',
  `ip` varchar(128) NOT NULL DEFAULT '',
  `tid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '第三方uid',
  PRIMARY KEY (`vid`),
  UNIQUE KEY `uid` (`uid`,`ip`,`tid`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=301 ;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
