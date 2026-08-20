package main

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/apache/thrift/lib/go/thrift"
	"mianshi/backend/gen/user"
)

const address = "127.0.0.1:9090"

// userHandler 是 IDL 中 UserService 的具体业务实现。
// 数据只保存在内存中，重启服务后会清空，便于把注意力放在 RPC 调用流程上。
type userHandler struct {
	mu     sync.RWMutex
	nextID user.UserID
	users  map[user.UserID]*user.User
}

func newUserHandler() *userHandler {
	return &userHandler{
		nextID: 1,
		users:  make(map[user.UserID]*user.User),
	}
}

func (h *userHandler) CreateUser(_ context.Context, req *user.CreateUserRequest) (*user.User, error) {
	if req == nil || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
		return nil, bizError(400, "name 和 email 不能为空")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for _, existing := range h.users {
		if existing.Email == req.Email {
			return nil, bizError(409, "email 已存在")
		}
	}

	u := &user.User{
		ID:       h.nextID,
		Name:     req.Name,
		Email:    req.Email,
		Status:   user.UserStatus_ACTIVE,
		Bio:      req.Bio,
		Tags:     append([]string(nil), req.Tags...),
		Metadata: map[string]string{"source": "thrift-client"},
	}
	h.users[u.ID] = u
	h.nextID++
	return cloneUser(u), nil
}

func (h *userHandler) GetUser(_ context.Context, id user.UserID) (*user.User, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	u, ok := h.users[id]
	if !ok {
		return nil, bizError(404, "用户不存在")
	}
	return cloneUser(u), nil
}

func (h *userHandler) ListUsers(_ context.Context, req *user.ListUsersRequest) ([]*user.User, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]*user.User, 0, len(h.users))
	for _, u := range h.users {
		if req != nil && req.Status != nil && u.Status != *req.Status {
			continue
		}
		result = append(result, cloneUser(u))
	}
	return result, nil
}

func (h *userHandler) UpdateStatus(_ context.Context, id user.UserID, status user.UserStatus) (bool, error) {
	if status != user.UserStatus_ACTIVE && status != user.UserStatus_DISABLED {
		return false, bizError(400, "无效的用户状态")
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	u, ok := h.users[id]
	if !ok {
		return false, bizError(404, "用户不存在")
	}
	u.Status = status
	return true, nil
}

func (h *userHandler) RecordLogin(_ context.Context, id user.UserID) error {
	// oneway 方法不应承载必须让客户端知道成功与否的业务。
	log.Printf("收到用户 %d 的登录通知", id)
	return nil
}

func bizError(code int32, message string) *user.BizException {
	return &user.BizException{Code: code, Message: message}
}

// 返回副本，避免服务端 map 内的数据被其他协程意外修改。
func cloneUser(src *user.User) *user.User {
	dst := *src
	dst.Tags = append([]string(nil), src.Tags...)
	dst.Metadata = make(map[string]string, len(src.Metadata))
	for k, v := range src.Metadata {
		dst.Metadata[k] = v
	}
	return &dst
}

func main() {
	serverTransport, err := thrift.NewTServerSocket(address)
	if err != nil {
		log.Fatal(err)
	}

	processor := user.NewUserServiceProcessor(newUserHandler())
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	protocolFactory := thrift.NewTBinaryProtocolFactoryConf(nil)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	log.Printf("UserService 已启动：%s", address)
	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
