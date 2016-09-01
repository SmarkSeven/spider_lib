package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common"    //选用
	//信息输出
	// net包
	"net/http" //设置http.Header
)

func init() {
	HaoJiaoLianNews.Register()
}

var HaoJiaoLianNews = &Spider{
	Name:        "好教练新闻搜索",
	Description: "好教练新闻搜索",
	// Pausetime: 300,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			// 调用指定Rule下辅助函数AidFunc()。
			ctx.Aid(map[string]interface{}{"loop": [2]int{0, 1}, "Rule": "新闻列表"}, "新闻列表")
		},

		Trunk: map[string]*Rule{

			"新闻列表": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					// for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
					ctx.AddQueue(&request.Request{
						Url:    "http://news.haojiaolian.com/xueche/bj/list_2.html",
						Rule:   aid["Rule"].(string),
						Header: http.Header{"Content-Type": []string{"text/html; charset=gbk"}},
					})
					// }
					return nil
				},

				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					query.Find("#info_list ul").Each(func(i int, s *goquery.Selection) {
						// 新闻标题
						title := s.Find("li > h2").Text()
						// 获取URL
						url, _ := s.Find("li > h2 a").First().Attr("href")
						// 摘要
						summary := s.Find("li > span").Text()

						ctx.AddQueue(&request.Request{
							Url:    url,
							Header: http.Header{"Content-Type": []string{"text/html; charset=gbk"}},
							Rule:   "新闻详情",
							Temp: map[string]interface{}{
								"title":   title,
								"url":     url,
								"summary": summary,
							},
						})
					})
				},
			},

			"新闻详情": {
				// 	//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"标题",
					"链接",
					"摘要",
					"时间",
					"内容",
				},

				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					newsTime := query.Find("#news_time").Text()
					content := query.Find("#content_news").Text()

					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0: ctx.GetTemp("title", ""),
						1: ctx.GetTemp("url", ""),
						2: ctx.GetTemp("summary", ""),
						3: newsTime,
						4: content,
					})
				},
			},
		},
	},
}
