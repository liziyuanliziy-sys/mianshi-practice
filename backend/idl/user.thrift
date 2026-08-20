// 在项目根目录重新生成：
// thrift --gen go:package_prefix=mianshi/backend/gen/ -out backend/gen backend/idl/user.thrift
// namespace go user 决定生成到 backend/gen/user 包中。
namespace go user

// typedef 能给基础类型补充业务含义，生成到 Go 后仍然是一个独立类型。
typedef i64 UserID

enum UserStatus {
  ACTIVE = 1,
  DISABLED = 2
}

// Thrift 字段都有编号。协议传输靠编号识别字段，而不是靠字段名。
// 已发布的字段编号不要复用，这样服务升级时才能保持兼容。
struct User {
  1: required UserID id,
  2: required string name,
  3: required string email,
  4: required UserStatus status,
  5: optional string bio,
  6: list<string> tags,
  7: map<string, string> metadata
}

struct CreateUserRequest {
  1: required string name,
  2: required string email,
  3: optional string bio,
  4: list<string> tags
}

struct ListUsersRequest {
  // optional 字段在 Go 中会生成为指针；nil 表示客户端没有传该条件。
  1: optional UserStatus status
}

// 业务错误应定义为 exception，而不是让客户端猜测普通字符串。
exception BizException {
  1: required i32 code,
  2: required string message
}

service UserService {
  User CreateUser(1: CreateUserRequest request)
    throws (1: BizException bizErr),

  User GetUser(1: UserID id)
    throws (1: BizException bizErr),

  // 筛选条件放进请求结构体，方便以后兼容地增加分页等字段。
  list<User> ListUsers(1: ListUsersRequest request),

  bool UpdateStatus(1: UserID id, 2: UserStatus status)
    throws (1: BizException bizErr),

  // oneway 只负责发送，不等待返回值，适合不要求立即确认的通知。
  oneway void RecordLogin(1: UserID id)
}
