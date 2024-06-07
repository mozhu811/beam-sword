package util

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetCurrentDate() string {
	return time.Now().Format("今天是2006年01月02日")
}

func Trim(content string) string {
	if strings.Contains(content, "```json") {
		content = content[7 : len(content)-3]
	}

	return strings.Trim(content, "\n")
}
func DoRequest(body *bytes.Buffer) *http.Response {
	req, err := http.NewRequest("POST",
		"https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation",
		bytes.NewBufferString(string(body.Bytes())),
	)

	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("API_KEY")))
	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	return resp
}
