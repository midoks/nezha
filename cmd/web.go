package cmd

import (
	"context"
	"log"

	"github.com/ory/graceful"
	"github.com/urfave/cli"

	"github.com/midoks/nezha/cmd/dashboard/controller"
	"github.com/midoks/nezha/cmd/dashboard/rpc"
	"github.com/midoks/nezha/model"
	"github.com/midoks/nezha/service/singleton"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "This command starts web services",
	Description: `Start Web services`,
	Action:      runWebService,
}

func init() {
	// 初始化 dao 包
	singleton.InitConfigFromPath("data/config.yaml")
	singleton.InitTimezoneAndCache()
	singleton.InitDBFromPath("data/sqlite.db")
	singleton.InitLocalizer()
	initSystem()
}

func initSystem() {
	// 启动 singleton 包下的所有服务
	singleton.LoadSingleton()

	// 每天的3:30 对 监控记录 和 流量记录 进行清理
	if _, err := singleton.Cron.AddFunc("0 30 3 * * *", singleton.CleanMonitorHistory); err != nil {
		panic(err)
	}

	// 每小时对流量记录进行打点
	if _, err := singleton.Cron.AddFunc("0 0 * * * *", singleton.RecordTransferHourlyUsage); err != nil {
		panic(err)
	}
}

func runWebService(c *cli.Context) error {
	singleton.CleanMonitorHistory()
	go rpc.ServeRPC(singleton.Conf.GRPCPort)
	serviceSentinelDispatchBus := make(chan model.Monitor) // 用于传递服务监控任务信息的channel
	go rpc.DispatchTask(serviceSentinelDispatchBus)
	go rpc.DispatchKeepalive()
	go singleton.AlertSentinelStart()
	singleton.NewServiceSentinel(serviceSentinelDispatchBus)
	srv := controller.ServeWeb(singleton.Conf.HTTPPort)
	graceful.Graceful(func() error {
		return srv.ListenAndServe()
	}, func(c context.Context) error {
		log.Println("NEZHA>> Graceful::START")
		singleton.RecordTransferHourlyUsage()
		log.Println("NEZHA>> Graceful::END")
		srv.Shutdown(c)
		return nil
	})
	return nil
}
