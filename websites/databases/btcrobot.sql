-- phpMyAdmin SQL Dump
-- version 4.0.9
-- http://www.phpmyadmin.net
--
-- 主机: 127.0.0.1
-- 生成日期: 2014-01-20 10:36:22
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

--
-- 转存表中的数据 `comments`
--

INSERT INTO `comments` (`cid`, `objid`, `objtype`, `content`, `uid`, `floor`, `flag`, `ctime`) VALUES
(1, 7, 0, 'sprite', 1, 1, 0, '2013-12-24 05:51:16'),
(3, 7, 0, '#1楼 @weiyan test', 2, 2, 0, '2013-12-25 05:06:51'),
(4, 7, 0, '支持 Markdown 格式, **粗体**、~~删除线~~、`单行代码`', 2, 3, 0, '2013-12-25 05:11:08'),
(5, 7, 0, 'test', 2, 4, 0, '2013-12-25 05:18:24'),
(6, 7, 0, '#3楼 @root test hi', 2, 5, 0, '2013-12-25 05:18:34'),
(12, 7, 0, '二楼孙子快来打爷我一楼的脸二楼孙子快来打爷我一楼的脸二楼孙子快来打爷我一楼的脸二', 1, 11, 0, '2013-12-25 06:39:43'),
(13, 8, 0, 'goo', 2, 1, 0, '2013-12-25 08:45:13'),
(14, 33, 0, '## Heading ## hi, who are you? guy?', 2, 1, 0, '2013-12-26 06:55:24'),
(15, 40, 0, 'suprise', 2, 1, 0, '2013-12-26 07:07:39'),
(16, 50, 0, '牛xxxx', 9684, 1, 0, '2013-12-26 09:30:58'),
(17, 80, 0, '这是什么？', 9685, 1, 0, '2013-12-30 06:59:06'),
(18, 80, 0, '#1楼 @37184891 分享一句名言：）', 2, 2, 0, '2013-12-30 07:25:49'),
(19, 80, 0, '来学习一下。界面UI不够上档次呀。', 9686, 3, 0, '2013-12-31 08:27:49'),
(20, 82, 0, '咋玩？', 9686, 1, 0, '2013-12-31 08:30:37'),
(21, 82, 0, '#1楼 @nickelchen 做产品啊', 2, 2, 0, '2014-01-02 05:02:23'),
(22, 80, 0, '#3楼 @nickelchen 请问能否参与啊', 2, 4, 0, '2014-01-02 05:03:12'),
(23, 80, 0, '#3楼 @nickelchen 能否指出具体的几点？', 2, 5, 0, '2014-01-02 05:03:30'),
(24, 86, 0, '哦，原来我是全栈程序员，哈哈', 1, 1, 0, '2014-01-19 11:03:19'),
(25, 85, 0, '网站分两步分，够浪前端和监控程序，监控让我关了。。。', 1, 1, 0, '2014-01-19 11:04:07'),
(26, 85, 0, '网站分两部分，够浪前端和监控程序，监控让我关了。。周一我去开启。', 1, 2, 0, '2014-01-19 11:04:30');

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

--
-- 转存表中的数据 `message`
--

INSERT INTO `message` (`id`, `content`, `hasread`, `from`, `fdel`, `to`, `tdel`, `ctime`) VALUES
(1, 'hitest', '已读', 1, '未删', 2, '未删', '2013-12-24 05:42:49'),
(2, 'test', '未读', 1, '未删', 10000, '未删', '2013-12-25 03:44:38'),
(3, 'hi, 你是本站第5000位会员。', '未读', 2, '未删', 9685, '未删', '2013-12-30 07:25:19');

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

--
-- 转存表中的数据 `role`
--

INSERT INTO `role` (`roleid`, `name`, `ctime`) VALUES
(1, '站长', '2013-12-24 03:39:59'),
(2, '副站长', '2013-12-24 03:39:59'),
(3, '超级管理员', '2013-12-24 03:39:59'),
(4, '社区管理员', '2013-12-24 03:39:59'),
(5, '资源管理员', '2013-12-24 03:39:59'),
(6, '酷站管理员', '2013-12-24 03:39:59'),
(7, '高级会员', '2013-12-24 03:39:59'),
(8, '中级会员', '2013-12-24 03:39:59'),
(9, '初级会员', '2013-12-24 03:39:59');

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

--
-- 转存表中的数据 `system_message`
--

INSERT INTO `system_message` (`id`, `msgtype`, `hasread`, `to`, `ext`, `ctime`) VALUES
(1, 0, '已读', 2, '{"cid":1,"objid":7,"objtype":0,"uid":1}', '2013-12-24 05:51:16'),
(2, 4, '已读', 1, '{"cid":3,"objid":7,"objtype":0,"uid":2}', '2013-12-25 05:06:51'),
(3, 0, '已读', 2, '{"cid":7,"objid":7,"objtype":0,"uid":1}', '2013-12-25 06:35:18'),
(4, 0, '已读', 2, '{"cid":8,"objid":7,"objtype":0,"uid":1}', '2013-12-25 06:36:14'),
(5, 0, '已读', 2, '{"cid":9,"objid":7,"objtype":0,"uid":1}', '2013-12-25 06:38:37'),
(6, 0, '已读', 2, '{"cid":10,"objid":7,"objtype":0,"uid":1}', '2013-12-25 06:39:10'),
(8, 0, '已读', 2, '{"cid":12,"objid":7,"objtype":0,"uid":1}', '2013-12-25 06:39:43'),
(9, 0, '已读', 1, '{"cid":13,"objid":8,"objtype":0,"uid":2}', '2013-12-25 08:45:13'),
(10, 0, '未读', 10000, '{"cid":14,"objid":33,"objtype":0,"uid":2}', '2013-12-26 06:55:25'),
(11, 0, '未读', 10000, '{"cid":15,"objid":40,"objtype":0,"uid":2}', '2013-12-26 07:07:39'),
(12, 0, '未读', 10000, '{"cid":17,"objid":80,"objtype":0,"uid":9685}', '2013-12-30 06:59:06'),
(13, 4, '未读', 9685, '{"cid":18,"objid":80,"objtype":0,"uid":2}', '2013-12-30 07:25:49'),
(14, 0, '未读', 10000, '{"cid":18,"objid":80,"objtype":0,"uid":2}', '2013-12-30 07:25:49'),
(15, 0, '未读', 10000, '{"cid":19,"objid":80,"objtype":0,"uid":9686}', '2013-12-31 08:27:49'),
(16, 0, '已读', 2, '{"cid":20,"objid":82,"objtype":0,"uid":9686}', '2013-12-31 08:30:38'),
(17, 4, '未读', 9686, '{"cid":21,"objid":82,"objtype":0,"uid":2}', '2014-01-02 05:02:24'),
(18, 4, '未读', 9686, '{"cid":22,"objid":80,"objtype":0,"uid":2}', '2014-01-02 05:03:12'),
(19, 0, '未读', 10000, '{"cid":22,"objid":80,"objtype":0,"uid":2}', '2014-01-02 05:03:12'),
(20, 4, '未读', 9686, '{"cid":23,"objid":80,"objtype":0,"uid":2}', '2014-01-02 05:03:30'),
(21, 0, '未读', 10000, '{"cid":23,"objid":80,"objtype":0,"uid":2}', '2014-01-02 05:03:30'),
(22, 0, '未读', 6, '{"cid":24,"objid":86,"objtype":0,"uid":1}', '2014-01-19 11:03:19'),
(23, 0, '未读', 6, '{"cid":25,"objid":85,"objtype":0,"uid":1}', '2014-01-19 11:04:07'),
(24, 0, '未读', 6, '{"cid":26,"objid":85,"objtype":0,"uid":1}', '2014-01-19 11:04:30');

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
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=89 ;

--
-- 转存表中的数据 `topics`
--

INSERT INTO `topics` (`tid`, `title`, `content`, `nid`, `uid`, `lastreplyuid`, `lastreplytime`, `view`, `reply`, `like`, `hate`, `flag`, `ctime`, `mtime`) VALUES
(3, '', '怎么能活在已知之中呢？\r\n', 0000000007, 2, 0, '0000-00-00 00:00:00', 14, 0, 2, 0, 0, '2013-12-24 05:08:42', '2014-01-20 07:38:16'),
(5, '', '绝大多数人宁愿相信符合他们自身利益的谬论，也不愿相信不符和他们自身利益的真理。', 0000000002, 2, 0, '0000-00-00 00:00:00', 16, 0, 0, 0, 0, '2013-12-24 05:12:10', '2014-01-20 03:40:11'),
(32, '', '![一张接地气！][1]\r\n\r\n\r\n  [1]: http://p4.zhimg.com/3b/01/3b01135662a652c1604fcfafaaa30d28_m.jpg', 0000000005, 2, 0, '0000-00-00 00:00:00', 12, 0, 0, 0, 0, '2013-12-25 08:46:18', '2014-01-19 23:00:29'),
(41, '', '1.01^365=37.78 -----------\r\n\r\n0.99^365=0.0255', 0000000005, 2, 0, '0000-00-00 00:00:00', 8, 0, 0, 0, 0, '2013-12-26 07:13:16', '2014-01-18 03:58:56'),
(42, '', 'Yesterday you said tomorrow.', 0000000005, 2, 0, '0000-00-00 00:00:00', 11, 0, 2, 0, 0, '2013-12-26 07:19:31', '2014-01-18 03:52:47'),
(82, '', '准备用golang做一个机器人，结合数据挖掘，实现一个比特币信息聚合平台，有意思的一起来。', 0000000005, 2, 2, '2014-01-02 05:02:23', 31, 2, 2, 2, 0, '2013-12-31 07:39:10', '2014-01-20 08:11:00'),
(83, '', '[Here''s The First Job Ad We''ve Seen For A Bitcoin Trader At A Hedge Fund\r\n\r\nRead more: http://www.businessinsider.com/bitcoin-trader-job-2014-1#ixzz2pJsXCYoi][1]\r\n\r\n\r\n  [1]: http://www.businessinsider.com/bitcoin-trader-job-2014-1', 0000000007, 1, 0, '0000-00-00 00:00:00', 5, 0, 2, 0, 0, '2014-01-03 07:50:03', '2014-01-20 08:29:59'),
(85, '', '怎么没看到你的交易记录了，运行得如何？\r\ntest:\r\n#### markdown\r\n\r\n    func main () [\r\n        fmt.Println("test")\r\n    }\r\n', 0000000005, 6, 1, '2014-01-19 11:04:30', 8, 2, 0, 0, 0, '2014-01-19 07:28:14', '2014-01-20 08:11:10'),
(86, '', '兄弟，刚才仔细看了一下你的网站，前端确实很棒，另外我也在学习golang，和做btc自动交易，我很佩服full stack engineer。我会没事来这里看看，学习一下，加油：）', 0000000005, 6, 1, '2014-01-19 11:03:19', 10, 1, 0, 0, 0, '2014-01-19 08:23:31', '2014-01-20 08:10:55'),
(87, '', '发送 0.01 BTC 到 "1NzBorutZGgKuT5VuqmmcY4QD7uTEaRr3g", 并在此留言，即可成为VIP，可以接受买入卖出点通知，若提供API access KEY可按照设定规则自动买入卖出。', 0000000004, 1, 0, '0000-00-00 00:00:00', 1, 0, 0, 0, 0, '2014-01-20 09:14:40', '2014-01-20 09:14:49'),
(88, '', '比特币ABC.,', 0000000002, 10000, 0, '0000-00-00 00:00:00', 0, 0, 0, 0, 0, '2014-01-20 09:32:58', '2014-01-20 09:32:58');

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
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=15 ;

--
-- 转存表中的数据 `topics_node`
--

INSERT INTO `topics_node` (`nid`, `parent`, `name`, `intro`, `ctime`) VALUES
(1, 0, '技术', '技术原理', '2013-12-24 03:39:59'),
(2, 1, '投资分析', '投资分析', '2013-12-24 03:39:59'),
(3, 1, '技术指标', '天注定', '2013-12-24 03:39:59'),
(4, 1, '自动交易', '', '2013-12-24 03:39:59'),
(5, 1, '自动提醒', '你在或者不在', '2013-12-24 03:39:59'),
(6, 0, 'BTC', '', '2013-12-24 03:39:59'),
(7, 0, 'LTC', '', '2013-12-24 03:39:59');

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

--
-- 转存表中的数据 `user_active`
--

INSERT INTO `user_active` (`uid`, `email`, `username`, `weight`, `avatar`, `mtime`) VALUES
(2, 'songbohr@gmail.com', 'root', 288, '', '2014-01-07 08:07:37'),
(6, 'donge@donge.org', 'donge', 22, '', '2014-01-19 08:23:31'),
(10000, 'test@163.com', 'test', 12, '', '2014-01-20 09:32:59');

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
  PRIMARY KEY (`uid`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=10001 ;

--
-- 转存表中的数据 `user_info`
--

INSERT INTO `user_info` (`uid`, `email`, `open`, `username`, `name`, `avatar`, `city`, `company`, `github`, `weibo`, `website`, `status`, `introduce`, `ctime`) VALUES
(1, '78623269@qq.com', 1, 'btcrobot', '', '', '', '', '', '', '', '', '', '2013-12-24 03:40:47'),
(2, 'songbohr@gmail.com', 1, 'root', 'god', '', 'Qingdao', '', '', 'bocaicfa', '', '', '', '2013-12-24 04:35:50'),
(6, 'donge@donge.org', 1, 'donge', '', '', '', '', '', '', '', '', '', '2014-01-19 07:23:41'),
(10000, 'test@163.com', 1, 'test', '', '', '', '', '', '', '', '', '', '2013-12-25 03:16:26');

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

--
-- 转存表中的数据 `user_login`
--

INSERT INTO `user_login` (`uid`, `email`, `username`, `passcode`, `passwd`) VALUES
(1, '78623269@qq.com', 'btcrobot', '5124654', 'a17e729fb11cfc74011b33d5a7012b55'),
(2, 'songbohr@gmail.com', 'root', '3f791a55', '7ffc385413b29223dd0614a31930e333'),
(6, 'donge@donge.org', 'donge', '7b3e718c', '62c0cea1e2fc3b0ab01cd0eb7843fe33'),
(10000, 'test@163.com', 'test', '45fb879f', '1b2529e8a85ab0583eb98e249f20a933');

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

--
-- 转存表中的数据 `user_role`
--

INSERT INTO `user_role` (`uid`, `roleid`, `ctime`) VALUES
(2, 1, '2013-12-24 04:35:50'),
(5, 9, '2013-12-25 03:16:26'),
(6, 9, '2014-01-19 07:23:42'),
(9684, 9, '2013-12-26 09:29:22'),
(9685, 9, '2013-12-30 06:58:40'),
(9686, 9, '2013-12-31 08:26:53');

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
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=299 ;

--
-- 转存表中的数据 `vote`
--

INSERT INTO `vote` (`vid`, `uid`, `ip`, `tid`) VALUES
(170, 1, '219.147.23.114', 35),
(137, 2, '127.0.0.1', 3),
(143, 2, '127.0.0.1', 11),
(135, 2, '127.0.0.1', 24),
(136, 2, '127.0.0.1', 26),
(146, 2, '127.0.0.1', 27),
(148, 2, '127.0.0.1', 28),
(149, 2, '127.0.0.1', 29),
(150, 2, '127.0.0.1', 30),
(220, 2, '127.0.0.1', 42),
(142, 2, '127.0.0.1', 50),
(141, 2, '127.0.0.1', 60),
(140, 2, '127.0.0.1', 61),
(234, 2, '127.0.0.1', 82),
(290, 10000, '101.226.102.97', 26),
(291, 10000, '101.226.167.198', 46),
(261, 10000, '101.226.33.217', 83),
(196, 10000, '101.226.33.226', 43),
(240, 10000, '112.15.173.240', 73),
(226, 10000, '112.254.98.199', 24),
(289, 10000, '112.64.235.91', 24),
(218, 10000, '114.255.192.96', 79),
(262, 10000, '117.25.127.235', 31),
(255, 10000, '118.192.168.153', 62),
(257, 10000, '119.147.146.189', 82),
(256, 10000, '123.15.49.66', 82),
(134, 10000, '127.0.0.1', 1),
(154, 10000, '127.0.0.1', 2),
(155, 10000, '127.0.0.1', 24),
(130, 10000, '127.0.0.1', 33),
(131, 10000, '127.0.0.1', 34),
(132, 10000, '127.0.0.1', 35),
(133, 10000, '127.0.0.1', 36),
(168, 10000, '159.226.43.61', 28),
(197, 10000, '180.153.214.152', 75),
(295, 10000, '182.118.20.162', 26),
(296, 10000, '182.118.20.167', 46),
(294, 10000, '182.118.20.175', 26),
(292, 10000, '182.118.20.183', 24),
(293, 10000, '182.118.20.187', 24),
(259, 10000, '182.118.22.247', 24),
(258, 10000, '182.118.22.249', 24),
(298, 10000, '182.118.25.233', 24),
(185, 10000, '183.61.117.17', 11),
(186, 10000, '183.61.117.17', 24),
(187, 10000, '183.61.117.17', 26),
(188, 10000, '183.61.117.17', 27),
(189, 10000, '183.61.117.17', 28),
(190, 10000, '183.61.117.17', 29),
(191, 10000, '183.61.117.17', 30),
(233, 10000, '210.76.108.133', 80),
(273, 10000, '219.140.166.34', 38),
(280, 10000, '219.140.166.34', 39),
(279, 10000, '219.140.166.34', 79),
(283, 10000, '219.146.252.174', 24),
(284, 10000, '219.146.252.174', 26),
(160, 10000, '219.147.23.114', 1),
(163, 10000, '219.147.23.114', 2),
(165, 10000, '219.147.23.114', 3),
(166, 10000, '219.147.23.114', 4),
(157, 10000, '219.147.23.114', 8),
(173, 10000, '219.147.23.114', 24),
(174, 10000, '219.147.23.114', 26),
(175, 10000, '219.147.23.114', 27),
(176, 10000, '219.147.23.114', 28),
(177, 10000, '219.147.23.114', 29),
(178, 10000, '219.147.23.114', 30),
(179, 10000, '219.147.23.114', 31),
(158, 10000, '219.147.23.114', 36),
(211, 10000, '219.147.23.114', 44),
(212, 10000, '219.147.23.114', 45),
(213, 10000, '219.147.23.114', 46),
(214, 10000, '219.147.23.114', 47),
(215, 10000, '219.147.23.114', 48),
(216, 10000, '219.147.23.114', 49),
(217, 10000, '219.147.23.114', 50),
(171, 10000, '219.147.23.114', 63),
(172, 10000, '219.147.23.114', 72),
(254, 10000, '220.178.118.197', 80),
(198, 10000, '221.11.183.25', 43),
(199, 10000, '221.11.183.25', 44),
(200, 10000, '221.11.183.25', 45),
(201, 10000, '221.11.183.25', 46),
(202, 10000, '221.11.183.25', 47),
(203, 10000, '221.11.183.25', 48),
(204, 10000, '221.11.183.25', 49),
(205, 10000, '221.11.183.25', 50),
(225, 10000, '222.69.242.53', 2),
(180, 10000, '27.210.72.230', 35),
(192, 10000, '27.210.72.230', 43),
(181, 10000, '27.210.72.230', 45),
(182, 10000, '27.210.72.230', 46),
(183, 10000, '27.210.72.230', 49),
(195, 10000, '27.210.72.230', 75),
(194, 10000, '27.210.72.230', 76),
(184, 10000, '27.210.72.230', 77),
(260, 10000, '27.210.72.230', 83),
(221, 10000, '58.240.65.50', 81),
(206, 10000, '58.57.73.75', 48),
(251, 10000, '59.66.123.198', 82),
(231, 10000, '61.135.24.210', 8),
(232, 10000, '61.135.24.210', 11),
(228, 10000, '61.135.24.210', 45),
(235, 10000, '61.142.131.119', 42),
(239, 10000, '61.142.131.119', 61);

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
