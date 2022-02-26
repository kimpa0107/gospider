package engine

import (
	"log"

	"jasper.com/gospider/fetcher"
)

func worker(r Request) (ParseResult, error) {
	log.Printf("Fetching %s\n", r.Url)
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("Fetcher: error fetching url %s: %v\n", r.Url, err)
		return ParseResult{}, err
	}

	return r.ParserFunc(body), nil
}

func createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			// tell scheduler that worker is ready
			ready.WorkerReady(in)
			request := <-in
			result, err := worker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
