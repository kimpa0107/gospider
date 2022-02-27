# 网站爬虫

## Usage

`main.go`

```go
baseURL := "http://xxx.com"
e := &engine.ConcurrentEngine{
  Scheduler:   &scheduler.QueuedScheduler{},
  WorkerCount: 50,
  ItemChan: persist.ItemSaver(persist.Option{
    SaveDBType:  config.DB_TYPE_MYSQL,
    WorkerCount: 10,
  }),
  FetcherOption: fetcher.Option{
    RateLimit: time.NewTicker(50 * time.Millisecond).C,
  },
  // 程序终止逻辑，自从最后一次爬取，过5分钟后仍然没有收到爬取的网页时，程序就会终止
  WaitingForFinish: time.NewTicker(5 * time.Minute).C,
}
e.Run(engine.Request{
  Url: baseURL + "/xxx/xxx",
  ParserFunc: func(b []byte) engine.ParseResult {
    return parser.ParseList(b, baseURL)
  },
  // 没有分页时，此解析器不需要
  ParsePagingFunc: func(body []byte) engine.ParseResult {
    return parser.ParseListPaging(body, baseURL)
  },
})
```

### Parser

在`parser`目录下新建不同网页的解析器，解析器命名`func ParseXxx(contents []byte, baseURL string)`

### Persist

`itemsaver.go`中定义数据持久化