package cmd

import (
	// "context"
	// "log"
	"fmt"
	"strings"

	"github.com/urfave/cli"

	"github.com/midoks/nezha/model"
	"github.com/midoks/nezha/pkg/utils"
	"github.com/midoks/nezha/service/singleton"
)

var Conf = cli.Command{
	Name:        "conf",
	Usage:       "This command set/show info",
	Description: `set/show info`,
	Action:      runConf,
	Flags: []cli.Flag{
		stringFlag("username, u", "", "show username"),
		stringFlag("set_username, su", "", "set username"),
		stringFlag("password, p", "", "set password"),
	},
}

func initDb() {
	// 初始化 dao 包
	singleton.InitConfigFromPath("data/config.yaml")
	singleton.InitTimezoneAndCache()
	singleton.InitDBFromPath("data/sqlite.db")

}

func runConf(c *cli.Context) error {
	initDb()

	var user model.User
	if err := singleton.DB.Where("id = ?", 1).Take(&user).Error; err == nil {
		username := c.String("username")
		if strings.EqualFold(username, "") {
			fmt.Println(user.Login)
		}

		set_username := c.String("set_username")
		if !strings.EqualFold(set_username, "") {
			user.Login = set_username
			singleton.DB.Save(&user)
			fmt.Println("Login modification completed")
		}

		password := c.String("password")
		if !strings.EqualFold(password, "") {
			password = utils.Md5(password)
			user.Password = password
			singleton.DB.Save(&user)
			fmt.Println("Password modification completed")
		}
	}

	return nil
}
