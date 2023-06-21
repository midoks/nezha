package controller

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/midoks/nezha/model"
	"github.com/midoks/nezha/pkg/csrf"
	"github.com/midoks/nezha/pkg/mygin"
	"github.com/midoks/nezha/pkg/utils"
	"github.com/midoks/nezha/service/singleton"
)

type guestPage struct {
	r *gin.Engine
}

type gLoginForm struct {
	Username string
	Password string
}

func (gp *guestPage) serve() {
	gr := gp.r.Group("")
	gr.Use(mygin.Authorize(mygin.AuthorizeOption{
		Guest:    true,
		IsPage:   true,
		Msg:      "您已登录",
		Btn:      "返回首页",
		Redirect: "/",
	}))

	gr.GET("/login", gp.login)
	gr.POST("/login", gp.postLogin)

}

func (gp *guestPage) postLogin(c *gin.Context) {
	csrf.VaildToken(c)

	username := c.PostForm("username")
	password := c.PostForm("password")

	args := mygin.CommonEnvironment(c, gin.H{
		"title": singleton.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "Login"}),
	})

	if strings.EqualFold(username, "") || strings.EqualFold(password, "") {
		args["LoginErrorMessage"] = "用户或密码不能为空!"
	} else {
		var u model.User
		if err := singleton.DB.Where("login = ?", username).First(&u).Error; err == nil {
			if utils.Md5(password) == u.Password {
				sess := sessions.Default(c)
				sess.Set("uid", u.ID)
				sess.Save()
				c.Redirect(http.StatusMovedPermanently, "/server")
				return
			}
			args["LoginErrorMessage"] = "密码错误!"
		} else {
			args["LoginErrorMessage"] = "用户或密码错误!"
		}
	}

	c.HTML(http.StatusOK, "dashboard-"+singleton.Conf.Site.DashboardTheme+"/login", args)
}

func (gp *guestPage) login(c *gin.Context) {
	args := mygin.CommonEnvironment(c, gin.H{
		"Title": singleton.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "Login"}),
	})

	args["Token"] = csrf.GetToken(c)
	c.HTML(http.StatusOK, "dashboard-"+singleton.Conf.Site.DashboardTheme+"/login", args)
}
