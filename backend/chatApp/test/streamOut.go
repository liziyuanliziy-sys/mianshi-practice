package test

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

func ReportStream(reader *schema.StreamReader[*schema.Message]) {
	defer reader.Close()
	for {
		message, err := reader.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("读取流消息失败：", err)
		}
		fmt.Print(message.Content)
	}
}

func Stream(ctx context.Context,
	model model.ToolCallingChatModel,
	messages []*schema.Message) *schema.StreamReader[*schema.Message] {
	result, err := model.Stream(ctx, messages)
	if err != nil {
		log.Fatal("大模型流式输出失败", err)
	}
	return result
}
