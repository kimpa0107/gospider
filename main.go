package main

import (
	"time"

	"jasper.com/gospider/config"
	"jasper.com/gospider/engine"
	"jasper.com/gospider/fetcher"
	"jasper.com/gospider/parser"
	"jasper.com/gospider/persist"
	"jasper.com/gospider/scheduler"
)

func main() {
	baseURL := "https:/xxx.com"
	e := &engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 100,
		ItemChan: persist.ItemSaver(persist.Option{
			SaveDBType:  config.DB_TYPE_MYSQL,
			WorkerCount: 10,
		}),
		FetcherOption: fetcher.Option{
			RateLimit: time.NewTicker(50 * time.Millisecond).C,
		},
		WaitingForFinish: time.NewTicker(5 * time.Minute).C,
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
