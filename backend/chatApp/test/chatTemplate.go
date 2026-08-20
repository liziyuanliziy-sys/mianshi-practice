package test

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func CreateChatTemplate() prompt.ChatTemplate {

	return prompt.FromMessages(schema.FString,
		//系统提示词
		schema.SystemMessage("你是{Role},用{style}的语气回答，帮助用户回答面试上的问题"),
		schema.UserMessage("{question}"),
	)
}

// 给大模型输入
func MessageTemplate() []*schema.Message {
	template := CreateChatTemplate()
	messages, err := template.Format(context.Background(), map[string]any{
		"Role":     "经验丰富的大厂开发面试专家",
		"style":    "温和且专业，条理清晰",
		"question": "什么是Go语言",
	})
	if err != nil {
		log.Fatal("格式化消息模版失败：", err)
	}
	return messages
}
