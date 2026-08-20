namespace go user

struct RegisterRequest{
    1: required string username (api.form="username")
    2: required string email (api.form="email")
    3: required string password (api.form="password")
}

struct LoginResponse{
    1: required string token
}

//服务的定义
service UserService{
    LoginResponse Register(1: RegisterRequest request)
        (api.post = "/api/user/register", api.category = "user", api.gen_path = "user")
}