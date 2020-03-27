package ext

//登录校验
type AuthsLogin interface {
	//账号验证
	Login() bool
	//证书验证
	SslCheck()bool
}