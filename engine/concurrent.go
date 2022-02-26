package engine

import (
	"log"
	"runtime"
	"time"

	"jasper.com/gospider/fetcher"
)

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorderCount int
	ItemChan    chan Item
}

type Scheduler interface {
	ReadyNotifier

	Submit(Request)
	WorkerChan() chan Request
	Run()
}

type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	start := time.Now()
	defer func() {
		log.Printf("Engine Run Time: %s\n", time.Since(start))
	}()

	r := seeds[0]

	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		return
	}

	requests := make([]Request, 0)
	var parseResult ParseResult

	// parse paging
	if r.ParsePagingFunc != nil {
		parseResult = r.ParsePagingFunc(body)
		requests = append(requests, parseResult.Requests...)
	}

	// parse first seed
	parseResult = r.ParserFunc(body)
	requests = append(requests, parseResult.Requests...)

	for _, item := range parseResult.Items {
		go func(item Item) {
			e.ItemChan <- item
		}(item)
	}

	// start engine
	e.run(requests...)
}

func (e *ConcurrentEngine) run(seeds ...Request) {
	out := make(chan ParseResult)
	e.Scheduler.Run()

	for i := 0; i < e.WorderCount; i++ {
		createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}

	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	for {
		select {
		case result := <-out:
			for _, item := range result.Items {
				go func(item Item) {
					e.ItemChan <- item
				}(item)
			}

			for _, r := range result.Requests {
				e.Scheduler.Submit(r)
			}

		// after all request done, break
		case <-time.After(5 * time.Minute):
			runtime.GC()
			return
		}
	}
}
