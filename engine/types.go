package engine

type Request struct {
	Url             string
	ParserFunc      func([]byte) ParseResult
	ParsePagingFunc func([]byte) ParseResult
}

type ParseResult struct {
	Requests []Request
	Items    []Item
}

type Item struct {
	Url     string
	Index   string
	ID      string
	Payload interface{}
}

func NilParser([]byte) ParseResult {
	return ParseResult{}
}
