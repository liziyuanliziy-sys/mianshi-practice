package comprehensive

//一个根据简历提问题的Agent
import (
	"context"
	"fmt"
	"log"
	"mianshi/backend/chatApp/chat"

	"github.com/cloudwego/eino/adk"
)

func NewSchoolAgent() (adk.Agent, error) {
	ctx := context.Background()
	model, err := chat.CreateOpenAiChatModel(ctx)
	if err != nil {
		log.Fatal("创建大模型失败: ", err)
		return nil, err
	}
	BaseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:          "SchoolAgent",
		Description:   "校招综合面试官智能体，全面评估应届毕业生的综合能力",
		Instruction:   SchoolComprehensiveAgentInstruction,
		Model:         model,
		MaxIterations: 3,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create school comprehensive agent: %w", err)
	}
	return BaseAgent, nil
}
