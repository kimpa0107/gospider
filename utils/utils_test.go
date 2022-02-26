package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestUnsbox(t *testing.T) {
	arg1 := "41A67966D8C2AB516D08C38837711A942E21AE43"

	unsbox := unsbox(arg1)
	fmt.Printf("unsbox: %s\n", unsbox)

	key := "3000176000856006061501533003690027800375"
	arg2 := hexXor(key, unsbox)
	fmt.Printf("arg2: %s\n", arg2)
}

func TestMD5(t *testing.T) {
	u := "https://hdhd114.net/webtoon/3597"
	h := MD5(u)
	fmt.Printf("MD5: %s\n", h)

	filename := "../storage/" + h + ".html"

	if _, err := os.Stat(filename); os.IsExist(err) {
		fmt.Println("File exists")
	}
}
