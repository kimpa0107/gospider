package fetcher

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"jasper.com/gospider/config"
	"jasper.com/gospider/utils"
)

const storageDir = "storage/"

var (
	rateLimit = time.NewTicker(10 * time.Millisecond).C
	client    *http.Client
)

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

func Fetch(url string) ([]byte, error) {
	// absFilePath := urlToFilepath(url)
	// if utils.IsFileExists(absFilePath) {
	// 	// log.Printf("[INFO] File %s exists, skip fetching.\n", absFilePath)
	// 	return ioutil.ReadFile(absFilePath)
	// }

	<-rateLimit

	client, _ := getProxyClient()
	req, _ := http.NewRequest("GET", url, nil)
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

	// cache(url, body)

	return body, err
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
