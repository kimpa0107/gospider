package main

import (
	"jasper.com/gospider/config"
	"jasper.com/gospider/engine"
	"jasper.com/gospider/parser"
	"jasper.com/gospider/persist"
	"jasper.com/gospider/scheduler"
)

func main() {
	baseURL := "https:/xxx.com"
	e := &engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorderCount: 100,
		ItemChan: persist.ItemSaver(persist.Option{
			SaveDBType:  config.DB_TYPE_MYSQL,
			WorkerCount: 10,
		}),
	}
	e.Run(engine.Request{
		Url: baseURL + "/xxx/xxx",
		ParserFunc: func(b []byte) engine.ParseResult {
			return parser.ParseList(b, baseURL)
		},
		ParsePagingFunc: func(body []byte) engine.ParseResult {
			return parser.ParseListPaging(body, baseURL)
		},
	})
}
