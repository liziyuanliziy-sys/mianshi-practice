package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"mianshi/backend/gen/user"
)

const address = "127.0.0.1:9090"

func main() {
	conf := &thrift.TConfiguration{
		ConnectTimeout: 3 * time.Second,
		SocketTimeout:  3 * time.Second,
	}
	socket := thrift.NewTSocketConf(address, conf)
	transport, err := thrift.NewTBufferedTransportFactory(8192).GetTransport(socket)
	if err != nil {
		log.Fatal(err)
	}
	if err := transport.Open(); err != nil {
		log.Fatalf("连接服务端失败（请先启动 service）：%v", err)
	}
	defer transport.Close()

	client := user.NewUserServiceClientFactory(
		transport,
		thrift.NewTBinaryProtocolFactoryConf(conf),
	)
	ctx := context.Background()

	// 1. 发送结构体参数并接收结构体结果。
	bio := "正在学习 Thrift"
	created, err := client.CreateUser(ctx, &user.CreateUserRequest{
		Name:  "小林",
		Email: "xiaolin@example.com",
		Bio:   &bio, // optional 字段使用指针；nil 代表没有传。
		Tags:  []string{"Go", "RPC"},
	})
	must(err)
	fmt.Printf("创建成功：id=%d name=%s status=%s\n", created.ID, created.Name, created.Status)

	// 2. 调用带返回值的查询方法。
	found, err := client.GetUser(ctx, created.ID)
	must(err)
	fmt.Printf("查询成功：%s <%s> tags=%v\n", found.Name, found.Email, found.Tags)

	// 3. optional 筛选条件：传 nil 表示列出全部用户。
	all, err := client.ListUsers(ctx, &user.ListUsersRequest{})
	must(err)
	fmt.Printf("用户总数：%d\n", len(all))

	active := user.UserStatus_ACTIVE
	activeUsers, err := client.ListUsers(ctx, &user.ListUsersRequest{Status: &active})
	must(err)
	fmt.Printf("启用用户数：%d\n", len(activeUsers))

	// 4. 修改状态，并发送一个不等待返回值的 oneway 通知。
	_, err = client.UpdateStatus(ctx, created.ID, user.UserStatus_DISABLED)
	must(err)
	must(client.RecordLogin(ctx, created.ID))
	fmt.Println("状态已更新，登录通知已发送")

	// 5. 声明在 IDL 中的业务异常，可以在客户端按具体类型处理。
	_, err = client.GetUser(ctx, 9999)
	var bizErr *user.BizException
	if errors.As(err, &bizErr) {
		fmt.Printf("收到预期业务异常：code=%d message=%s\n", bizErr.Code, bizErr.Message)
	} else {
		must(err)
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
