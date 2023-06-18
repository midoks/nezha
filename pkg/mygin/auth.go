package mygin

import (
	"fmt"
	"net/http"
	// "strings"
	// "time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/midoks/nezha/model"
	"github.com/midoks/nezha/service/singleton"
)

type AuthorizeOption struct {
	Guest    bool
	Member   bool
	IsPage   bool
	AllowAPI bool
	Msg      string
	Redirect string
	Btn      string
}

func Authorize(opt AuthorizeOption) func(*gin.Context) {
	return func(c *gin.Context) {

		var code = http.StatusForbidden
		if opt.Guest {
			code = http.StatusBadRequest
		}

		commonErr := ErrInfo{
			Title: "访问受限",
			Code:  code,
			Msg:   opt.Msg,
			Link:  opt.Redirect,
			Btn:   opt.Btn,
		}
		var isLogin bool

		// 用户鉴权
		sess := sessions.Default(c)
		uid := sess.Get("uid")
		fmt.Println(uid)
		if uid != nil {
			var u model.User
			if err := singleton.DB.Where("id = ?", uid).First(&u).Error; err == nil {
				fmt.Println(u)
				isLogin = true
			}
			if isLogin {
				c.Set(model.CtxKeyAuthorizedUser, &u)
			}
		}

		// API鉴权
		if opt.AllowAPI {
			apiToken := c.GetHeader("Authorization")
			if apiToken != "" {
				var u model.User
				singleton.ApiLock.RLock()
				if _, ok := singleton.ApiTokenList[apiToken]; ok {
					err := singleton.DB.First(&u).Where("id = ?", singleton.ApiTokenList[apiToken].UserID).Error
					isLogin = err == nil
				}
				singleton.ApiLock.RUnlock()
				if isLogin {
					c.Set(model.CtxKeyAuthorizedUser, &u)
					c.Set("isAPI", true)
				}
			}
		}

		// 已登录且只能游客访问
		if isLogin && opt.Guest {
			ShowErrorPage(c, commonErr, opt.IsPage)
			return
		}
		// 未登录且需要登录
		if !isLogin && opt.Member {
			ShowErrorPage(c, commonErr, opt.IsPage)
			return
		}
	}
}
