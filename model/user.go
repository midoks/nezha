package model

import (
	"time"
)

type User struct {
	Common
	Login     string `json:"login,omitempty"`      // 登录名
	Password  string `json:"password,omitempty"`   // 密码
	AvatarURL string `json:"avatar_url,omitempty"` // 头像地址
	Name      string `json:"name,omitempty"`       // 昵称
	Blog      string `json:"blog,omitempty"`       // 网站链接
	Email     string `json:"email,omitempty"`      // 邮箱
	Hireable  bool   `json:"hireable,omitempty"`
	Bio       string `json:"bio,omitempty"` // 个人简介

	Token        string    `json:"-"`                       // 认证 Token
	TokenExpired time.Time `json:"token_expired,omitempty"` // Token 过期时间
	SuperAdmin   bool      `json:"super_admin,omitempty"`   // 超级管理员
}

func NewUserFromUP(username string, pwd string) User {
	var u User
	u.Login = username
	u.Password = pwd
	u.AvatarURL = ""
	u.Name = "nezha"
	u.Blog = ""
	u.Email = ""
	u.Bio = ""
	return u
}
