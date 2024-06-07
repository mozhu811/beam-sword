package llm

import "encoding/json"

type Function struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments"`
}

type Message struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Name      string `json:"name,omitempty"`
	ToolCalls []struct {
		Func Function `json:"function,omitempty"`
		Id   string   `json:"id,omitempty"`
		Type string   `json:"type,omitempty"`
	} `json:"tool_calls,omitempty"`
}

type Params struct {
	Prompt   string
	Messages []Message
}

func (m *Message) String() string {
	s, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(s)
}
