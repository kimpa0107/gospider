package engine

import (
	"jasper.com/gospider/fetcher"
)

type SimpleEngine struct {
	ItemChan chan interface{}
}

func (e SimpleEngine) Run(seeds ...Request) {
	r := seeds[0]

	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		return
	}

	requests := make([]Request, 0)
	var parseResult ParseResult

	// parse paging
	if r.ParserFunc != nil {
		parseResult = r.ParsePagingFunc(body)
		requests = append(requests, parseResult.Requests...)
	}

	// parse first seed
	parseResult = r.ParserFunc(body)
	requests = append(requests, parseResult.Requests...)

	for _, item := range parseResult.Items {
		go func(item interface{}) {
			e.ItemChan <- item
		}(item)
	}

	// start engine
	e.run(requests...)
}

func (e SimpleEngine) run(seeds ...Request) {
	var requests []Request
	requests = append(requests, seeds...)

	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]

		parseResult, err := worker(r)
		if err != nil {
			continue
		}

		requests = append(requests, parseResult.Requests...)

		for _, item := range parseResult.Items {
			go func(item interface{}) {
				e.ItemChan <- item
			}(item)
		}
	}
}
