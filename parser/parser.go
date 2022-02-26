package parser

import (
	"jasper.com/gospider/engine"
)

// list page parser

func ParseList(contents []byte, baseURL string) engine.ParseResult {
	requests := make([]engine.Request, 0)
	items := make([]engine.Item, 0)

	// TODO: parse list items

	// TODO: append requests and items

	return engine.ParseResult{
		Requests: requests,
		Items:    items,
	}
}

func ParseListPaging(contents []byte, baseURL string) engine.ParseResult {
	requests := []engine.Request{}

	// TODO: parse paging links

	return engine.ParseResult{
		Requests: requests,
	}
}
