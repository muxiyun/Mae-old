package main

import (

	"github.com/muxiyun/Mae/config"
	"github.com/muxiyun/Mae/model"
	"github.com/muxiyun/Mae/pkg/casbin"
	"github.com/muxiyun/Mae/pkg/mail"
	"github.com/muxiyun/Mae/router"

	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"github.com/lexkong/log"
)


//for test
func newApp() *iris.Application {
	//init the config and log
	if err := config.Init("./conf/config.yaml"); err != nil {
		panic(err)
	}

	//init db
	model.DB.Init()

	// init casbin
	casbin.Init()
	//Mae app
	app := iris.Default()

	//register routers to Mae app
	app = router.Load(app)

	//add init casbin policy to db
	casbin.InitPolicy()

	// setup mailservice
	mail.Setup()

	return app
}

func main() {

	printLogo()

	app := newApp()

	// start mail daemon
	go mail.StartMailDaemon()

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()

	// Start to listening the incoming requests.
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")
	if cert != "" && key != "" {
		go func() {
			log.Infof("Start to listeni5ng the incoming requests on https address: %s", viper.GetString("tls.addr"))
			log.Info(app.Run(iris.TLS(viper.GetString("tls.addr"), cert, key), iris.WithoutVersionChecker).Error())
		}()
	}

	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Info(app.Run(iris.Addr(viper.GetString("addr")), iris.WithoutVersionChecker).Error())


	model.DB.RWdb.Close()
}
