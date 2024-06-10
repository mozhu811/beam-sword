package llm

import (
	"beam-sword/util"
	"bytes"
	"encoding/json"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"strings"
	"text/template"
)

func Qwen(userContent string) (string, error) {
	payload := `
	{
		"model": "qwen-max",
		"input":{
			"messages": [{{ generateMessages .Messages }}]
		},
		"parameters": {
			"result_format": "message",
			"tools": 
				[
					{
						"type": "function",
						"function": {
							"name": "GetCurrentDate",
							"description": "当你想知道今天的日期时非常有用，如果你需要我告诉你今天的日期，请你调用这个函数。",
							"parameters": {}
						}
					}
				]
		}
	}
	`
	funcMap := template.FuncMap{}
	funcMap["generateMessages"] = generateMessages

	p := &Params{
		Messages: []Message{
			{Role: "system", Content: `
			你是一位记账的管家，现在请你为我记录我的日常账单。要求：
			1. 如果你需要我提供今天的日期，请你调用函数GetCurrentDate。
			2. 你需要对账单类型进行分类，相似度高的分类整合为一类。如购买猫砂的账单和猫粮的账单，同意归类为宠物消费。点外卖和聚餐都归类为餐饮
			3. 如果是支出，金额使用负数表示。
			4. 所有账单类型必须使用JSON数组返回，并且只需要提供JSON数据，并且消息类型为纯文本，不要使用markdown。
			5. 你不需要输出其他文字，我只需要你提供我的账单JSON数据。请你严格按照以下两个示例进行回复，只输出结果。
			6. 以下示例中的日期date属性不是真实的数据，如果你需要真实的数据，请你调用函数GetCurrentDate
			示例1:
			输入：
			"""
			我今天点外卖花了20快钱
			"""

			输出
			"""
			[{"event": "点外卖", type": "支出", "amount": -20, "tag": "餐饮", "date":"2024/06/07"}]
			"""
			
			示例2:
			输入：
			"""
			我昨天点外卖花了20快钱，然后花了80块钱买奶茶
			"""

			输出
			"""
			[{"event": "点外卖", "type": "支出", "amount": -20, "tag": "餐饮", "date":"2024/06/06"},{"event": "买奶茶", "type": "支出", "amount": -80, "tag": "餐饮", "date":"2024/06/06"}]
			"""

			示例3:
			输入：
			"""
			今天发了两万的工资
			"""

			输出
			"""
			[{"event": "发工资", "type": "收入", "amount": 20000, "tag": "工资", "date":"2024/06/06"}]
			"""
`},
			{Role: "user", Content: "今天是" + util.GetCurrentDate() + "，" + userContent}},
	}

	tpl, err := template.New("template").Funcs(funcMap).Parse(payload)
	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer
	err = tpl.Execute(&buff, p)
	if err != nil {
		panic(err)
	}
	resp := util.DoRequest(&buff)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(string(body))
	}
	data := string(body)

	message := gjson.Get(data, "output.choices.0.message")
	aiContent := message.Get("content").String()
	if !strings.Contains(message.Raw, "tool_calls") {
		return util.Trim(aiContent), nil
	}
	p.Messages = append(p.Messages, Message{
		Role:    message.Get("role").String(),
		Content: aiContent,
		ToolCalls: []struct {
			Func Function `json:"function,omitempty"`
			Id   string   `json:"id,omitempty"`
			Type string   `json:"type,omitempty"`
		}{
			{
				Func: Function{
					Name:      message.Get("tool_calls.0.function.name").String(),
					Arguments: message.Get("tool_calls.0.function.arguments").String(),
				},
				Id:   message.Get("id").String(),
				Type: message.Get("type").String(),
			},
		},
	})
	// 需要调用函数
	// 获取函数名
	funcs := message.Get("tool_calls.#.function")

	funcs.ForEach(func(key, value gjson.Result) bool {
		funcName := value.Get("name")
		_ = value.Get("arguments")
		toolMsg := Message{
			Role: "tool",
			Name: funcName.String(),
		}
		if funcName.String() == "GetCurrentDate" {
			toolMsg.Content = util.GetCurrentDate()
		}
		p.Messages = append(p.Messages, toolMsg)
		return true
	})

	buff.Reset()
	err = tpl.Execute(&buff, p)
	if err != nil {
		panic(err)
	}
	resp = util.DoRequest(&buff)
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(string(body))
	}
	data = string(body)
	aiContent = gjson.Parse(data).Get("output.choices.0.message.content").String()

	return util.Trim(aiContent), nil
}

func generateMessages(messages []Message) (string, error) {
	var sb strings.Builder
	for i, m := range messages {
		jsonByte, err := json.Marshal(m)
		if err != nil {
			return "", err
		}
		sb.WriteString(string(jsonByte))
		if i != len(messages)-1 {
			sb.WriteString(",")
		}
	}
	return sb.String(), nil
}
