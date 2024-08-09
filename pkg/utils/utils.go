package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// ExtractHostAndPath extracts the host and path from the given URL
func ExtractHostAndPath(serverUrl string) (string, string, error) {
	u, err := url.Parse(serverUrl)
	if err != nil {
		return "", "", err
	}
	return u.Host, u.Path, nil
}

// PrintHttpRequestHeader prints the HTTP request headers
func PrintHttpRequestHeader(request *http.Request) {
	for name, headers := range request.Header {
		//name = strings.ToLower(name)
		for _, h := range headers {
			log.Printf("%v: %v\n", name, h)
		}
	}
}

func PrintJson(jsonS interface{}) {
	jsonData, _ := json.Marshal(jsonS)

	log.Println(string(jsonData))
}

// MapToString 将 map[int]string 转换为单个字符串
func MapToString(m map[int]string) string {
	// 获取 map 的键
	keys := make([]int, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	// 对键进行排序
	sort.Ints(keys)

	// 创建结果切片
	var result []string
	for _, key := range keys {
		result = append(result, fmt.Sprintf("%s", m[key]))
	}

	// 拼接结果
	return strings.Join(result, "")
}
