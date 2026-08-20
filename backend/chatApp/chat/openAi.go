package chat

import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/joho/godotenv"
)

func CreateOpenAiChatModel(ctx context.Context) (model.ToolCallingChatModel, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("读取 .env 失败：", err)
		return nil, err
	}
	APIKEY := os.Getenv("OPENAI_API_KEY")
	BASEURL := os.Getenv("BASE_URL")
	BASENAME := os.Getenv("BASE_NAME")
	ChatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  APIKEY,
		BaseURL: BASEURL,
		Model:   BASENAME,
	})
	if err != nil {
		log.Fatal("创建大模型失败：", err)
		return nil, err
	}
	return ChatModel, nil

}
