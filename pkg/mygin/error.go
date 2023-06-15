package mygin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/midoks/nezha/model"
	"github.com/midoks/nezha/service/singleton"
)

type ErrInfo struct {
	Code  int
	Title string
	Msg   string
	Link  string
	Btn   string
}

func ShowErrorPage(c *gin.Context, i ErrInfo, isPage bool) {
	if isPage {
		fmt.Println("dashboard-" + singleton.Conf.Site.DashboardTheme + "/error")
		c.HTML(i.Code, "dashboard-"+singleton.Conf.Site.DashboardTheme+"/error", CommonEnvironment(c, gin.H{
			"Code":  i.Code,
			"Title": i.Title,
			"Msg":   i.Msg,
			"Link":  i.Link,
			"Btn":   i.Btn,
		}))
	} else {
		c.JSON(http.StatusOK, model.Response{
			Code:    i.Code,
			Message: i.Msg,
		})
	}
	c.Abort()
}
