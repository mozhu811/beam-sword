package main

import (
	"beam-sword/lark"
	"beam-sword/llm"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/bs", func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		m := make(map[string]string)
		_ = json.Unmarshal(reqBody, &m)
		content, err := llm.Qwen(m["userContent"])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(content)

		err = lark.CreateRecord(content)
		if err != nil {
			fmt.Println(err)
		}
	})

	// 定义服务器的端口
	port := "8080"

	// 启动HTTP服务器
	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
