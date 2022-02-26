package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Base64Decode(str string) string {
	b, _ := base64.StdEncoding.DecodeString(str)
	return string(b)
}

func GetRootPath() string {
	pwd, _ := os.Getwd()
	return pwd + "/"
}

func ResolvePath(path string) string {
	rootPath := GetRootPath()
	path = strings.TrimLeft(path, "/")
	return rootPath + path
}

func IsFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

func CalculateCookie(url string) (string, error) {
	key := "3000176000856006061501533003690027800375"
	name := "acw_sc__v2"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
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

	arg1 := getArg1(body)
	if arg1 == "" {
		log.Println("[ERROR] failed to get arg1")
		return "", fmt.Errorf("failed to get arg1")
	}

	calc := unsbox(arg1)
	arg2 := hexXor(key, calc)

	cookie := name + "=" + arg2

	return cookie, nil
}

func getArg1(body []byte) string {
	re := regexp.MustCompile(`var arg1='([^']+)'`)
	match := re.FindSubmatch(body)
	if len(match) != 2 {
		return ""
	}
	return string(match[1])
}

func unsbox(arg string) string {
	/*
		# python code
		arr = [
			0xf, 0x23, 0x1d, 0x18, 0x21, 0x10, 0x1, 0x26,
			0xa, 0x9, 0x13, 0x1f, 0x28, 0x1b, 0x16, 0x17,
			0x19, 0xd, 0x6, 0xb, 0x27, 0x12, 0x14, 0x8,
			0xe, 0x15, 0x20, 0x1a, 0x2, 0x1e, 0x7, 0x4,
			0x11, 0x5, 0x3, 0x1c, 0x22, 0x25, 0xc, 0x24
		]
		tmp = [''] * 40
		res = ''
		for index in range(0, len(arg)):
			char = arg[index]
			for idx in range(0, len(arr)):
				if arr[idx] == index + 0x1:
					tmp[idx] = char
		res = ''.join(tmp)
		return res
	*/
	define := []int{
		0xf, 0x23, 0x1d, 0x18, 0x21, 0x10, 0x1, 0x26,
		0xa, 0x9, 0x13, 0x1f, 0x28, 0x1b, 0x16, 0x17,
		0x19, 0xd, 0x6, 0xb, 0x27, 0x12, 0x14, 0x8,
		0xe, 0x15, 0x20, 0x1a, 0x2, 0x1e, 0x7, 0x4,
		0x11, 0x5, 0x3, 0x1c, 0x22, 0x25, 0xc, 0x24,
	}
	temp := make([]string, len(define))
	for index, char := range arg {
		c := string(char)
		for idx := range define {
			if define[idx] == index+0x1 {
				temp[idx] = c
			}
		}
	}
	return strings.Join(temp, "")
}

func hexXor(key, unsbox string) string {
	/*
		# python code
		res = ''
		index = 0x0
		while index < len(unsbox_str) and index < len(key):
			num1 = int(unsbox_str[index: index + 0x2], 16)
			num2 = int(key[index: index + 0x2], 16)
			tmp = hex(num1 ^ num2)
			if len(tmp) == 0x1:
				tmp = '\x30' + tmp
			res += tmp[2:]

			index += 0x2
		return res
	*/

	res := ""
	index := 0x0
	for index < len(unsbox) && index < len(key) {
		num1 := hexToInt(unsbox[index : index+0x2])
		num2 := hexToInt(key[index : index+0x2])
		tmp := intToHex(num1 ^ num2)
		if len(tmp) == 0x1 {
			tmp = "0" + tmp
		}
		res += tmp[2:]
		index += 0x2
	}
	return res
}

func hexToInt(hex string) int {
	i, _ := strconv.ParseInt(hex, 16, 64)
	return int(i)
}

func intToHex(i int) string {
	return fmt.Sprintf("0x%x", i)
}
