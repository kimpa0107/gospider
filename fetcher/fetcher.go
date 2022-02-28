package fetcher

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"jasper.com/gospider/config"
	"jasper.com/gospider/utils"
)

const storageDir = "storage/"

var (
	client *http.Client

	option     Option
	optionInit = false

	defaultRateLimit = time.NewTicker(50 * time.Millisecond).C
)

type Option struct {
	// wait for fetching each url
	RateLimit <-chan time.Time
	// whether use proxy
	UseProxy bool
	// whether cache the fetched content
	CacheData bool
	// whether read from cache file first
	ReadCache bool
}

// init fetcher option once
func Init(opt Option) {
	if optionInit {
		return
	}

	option = opt

	if option.RateLimit == nil {
		option.RateLimit = defaultRateLimit
	}

	if option.UseProxy {
		client, _ = getProxyClient()
	} else {
		client, _ = getHttpClient()
	}

	optionInit = true
}

func Fetch(url string, method string, data map[string]interface{}) ([]byte, error) {
	if option.ReadCache {
		absFilePath := urlToFilepath(url)
		if utils.IsFileExists(absFilePath) {
			return ioutil.ReadFile(absFilePath)
		}
	}

	<-option.RateLimit

	var payload io.Reader
	if data != nil {
		payload = strings.NewReader(utils.ToJSON(data))
	}

	req, _ := http.NewRequest(method, url, payload)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	r := bufio.NewReader(resp.Body)

	var body []byte

	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			break
		}
		body = append(body, line...)
	}

	if option.CacheData {
		cache(url, body)
	}

	return body, nil
}

func getHttpClient() (*http.Client, error) {
	return &http.Client{}, nil
}

func getProxyClient() (*http.Client, error) {
	if client != nil {
		return client, nil
	}

	parsedURL, err := url.Parse(fmt.Sprintf("%s://%s:%d",
		config.PROXY_SCHEME,
		config.PROXY_HOST,
		config.PROXY_PORT))
	if err != nil {
		return nil, err
	}
	client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(parsedURL),
		},
	}
	return client, nil
}

// write content to file for caching
func cache(url string, content []byte) {
	if content == nil {
		return
	}

	// base64 encode url to make filename
	absFilePath := urlToFilepath(url)

	err := ioutil.WriteFile(absFilePath, content, 0644)
	if err != nil {
		log.Printf("[ERROR] writing to file [%s]: %s", absFilePath, err)
	}
}

func urlToFilepath(url string) string {
	if !utils.IsFileExists(storageDir) {
		err := os.MkdirAll(storageDir, 0755)
		if err != nil {
			log.Printf("[ERROR] creating directory [%s]: %s", storageDir, err)
		}
	}
	filename := utils.MD5(url)
	filepath := storageDir + filename + ".html"
	return utils.ResolvePath(filepath)
}
