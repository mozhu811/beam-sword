package util

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetCurrentDate() string {
	return time.Now().Format("2006年01月02日")
}

func Trim(content string) string {
	if strings.Contains(content, "```json") {
		content = content[7 : len(content)-3]
	}

	return strings.Trim(content, "\n")
}

func AskQwen(body *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest("POST",
		"https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation",
		bytes.NewBufferString(string(body.Bytes())),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("API_KEY")))
	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	return resp, nil
}
