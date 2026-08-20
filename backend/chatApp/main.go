package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"mianshi/backend/chatApp/agent/comprehensive"
	"mianshi/backend/chatApp/agent/resume"
	"os"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

func main() {

	ctx := context.Background()

	schoolagent, err := comprehensive.NewSchoolAgent()
	if err != nil {
		log.Fatalf("校招智能体创建失败: %v", err)
	}
	resumeAgent, err := resume.NewResumeAgent()
	if err != nil {
		log.Fatalf("简历智能体创建失败: %v", err)
	}

	//按顺序执行校招智能体和简历智能体
	sequentialAgent, err := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        "SequentialAgent",
		Description: "按顺序执行校招智能体和简历智能体",
		SubAgents:   []adk.Agent{resumeAgent, schoolagent},
	})
	if err != nil {
		log.Fatalf("顺序智能体创建失败: %v", err)
	}

	AgentTest(sequentialAgent)

}

func AgentTest(agent adk.Agent) {

	ctx := context.Background()

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
	})

	fmt.Println("====== Agent 交互测试 ======")
	fmt.Println("输入问题进行测试，输入 'exit' 或 'quit' 退出")
	fmt.Println("=====================================")

	processInput := func(input string) {
		message := []*schema.Message{
			schema.UserMessage(input),
		}

		//启动智能体
		iter := runner.Run(ctx, message)

		fmt.Print("Agent:")

		for {
			event, ok := iter.Next()

			if !ok {
				break
			}

			if event.Err != nil {
				log.Fatalf("智能体运行失败: %v", event.Err)
				break
			}

			// 打印智能体输出
			if event.Output != nil && event.Output.MessageOutput != nil {
				content := event.Output.MessageOutput.Message.Content
				if content != "" {
					fmt.Printf("%s", content)
				}
			}
		}
		fmt.Println()
	}

	autoInput := "D:\\NewProject\\mianshi\\backend\\chatApp\\GoTest.pdf"

	if autoInput != "" {
		fmt.Printf("自动输入: %s\n", autoInput)
		processInput(autoInput)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n 用户输入：")

		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("用户输入失败: %v", err)
			continue
		}
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("退出测试")
			break
		}
		processInput(input)
	}

}
