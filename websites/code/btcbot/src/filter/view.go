/*
  btcbot is a Bitcoin trading bot for HUOBI.com written
  in golang, it features multiple trading methods using
  technical analysis.

  Disclaimer:

  USE AT YOUR OWN RISK!

  The author of this project is NOT responsible for any damage or loss caused
  by this software. There can be bugs and the bot may not perform as expected
  or specified. Please consider testing it first with paper trading /
  backtesting on historical data. Also look at the code to see what how
  it's working.

  Weibo:http://weibo.com/bocaicfa
*/

package filter

import (
	"config"
	"github.com/gorilla/context"
	"github.com/studygolang/mux"
	"html/template"
	"logger"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"util"
)

// 自定义模板函数
var funcMap = template.FuncMap{
	// 转为前端显示需要的时间格式
	"formatTime": func(i interface{}) string {
		ctime, ok := i.(string)
		if !ok {
			return ""
		}
		t, _ := time.Parse("2006-01-02 15:04:05", ctime)
		return t.Format(time.RFC3339) + "+08:00"
	},
	"substring": func(str string, length int, suffix string) string {
		if length >= len(str) {
			return str
		}
		utf8Str := util.NewString(str)
		if length > utf8Str.RuneCount() {
			return str
		}
		return utf8Str.Slice(0, length) + suffix
	},
	// if 比较
	"eq": func(a, b string) bool {
		if a == b {
			return true
		}
		return false
	},
	"noteq": func(a, b string) bool {
		if a == b {
			return false
		}
		return true
	},
}

// 保存模板路径的key
const CONTENT_TPL_KEY = "__content_tpl"

// 页面展示 过滤器
type ViewFilter struct {
	commonHtmlFiles []string // 通用的html文件
	baseTplName     string   // 第一个基础模板的名称

	// "继承"空实现
	*mux.EmptyFilter
}

func NewViewFilter(files ...string) *ViewFilter {
	viewFilter := new(ViewFilter)
	if len(files) == 0 {
		// 默认使用前端通用模板
		viewFilter.commonHtmlFiles = []string{config.ROOT + "/template/common/base.html"}
		viewFilter.baseTplName = "base.html"
	} else {
		viewFilter.commonHtmlFiles = files
		viewFilter.baseTplName = filepath.Base(files[0])
	}
	return viewFilter
}

func (this *ViewFilter) PreFilter(rw http.ResponseWriter, req *http.Request) bool {
	// ajax请求头设置
	if strings.HasSuffix(req.RequestURI, ".json") {
		logger.Debugln(req.RequestURI)
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	return true
}

// 在逻辑处理完之后，最后展示页面
func (this *ViewFilter) PostFilter(rw http.ResponseWriter, req *http.Request) bool {
	contentHtml := req.FormValue(CONTENT_TPL_KEY)
	if contentHtml == "" {
		return true
	}
	contentHtmls := strings.Split(contentHtml, ",")
	for i, contentHtml := range contentHtmls {
		contentHtmls[i] = config.ROOT + strings.TrimSpace(contentHtml)
	}

	// 为了使用自定义的模板函数，首先New一个以第一个模板文件名为模板名。
	// 这样，在ParseFiles时，新返回的*Template便还是原来的模板实例
	tpl, err := template.New(this.baseTplName).Funcs(funcMap).ParseFiles(append(this.commonHtmlFiles, contentHtmls...)...)
	if err != nil {
		logger.Errorf("解析模板出错（ParseFiles）：[%q] %s\n", req.RequestURI, err)
		return false
	}
	// 如果没有定义css和js模板，则定义之
	if jsTpl := tpl.Lookup("js"); jsTpl == nil {
		tpl.Parse(`{{define "js"}}{{end}}`)
	}
	if jsTpl := tpl.Lookup("css"); jsTpl == nil {
		tpl.Parse(`{{define "css"}}{{end}}`)
	}

	data := GetData(req)

	err = tpl.Execute(rw, data)
	if err != nil {
		logger.Errorf("执行模板出错（Execute）：[%q] %s\n", req.RequestURI, err)
	}
	return true
}

type viewKey int

const datakey viewKey = 0

func GetData(req *http.Request) map[string]interface{} {
	if rv := context.Get(req, datakey); rv != nil {
		// 获取之后立马删除
		context.Delete(req, datakey)
		return rv.(map[string]interface{})
	}
	return make(map[string]interface{})
}

func SetData(req *http.Request, data map[string]interface{}) {
	context.Set(req, datakey, data)
}
