package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/muxiyun/Mae/config"
	"github.com/muxiyun/Mae/model"
	"github.com/muxiyun/Mae/router"

	"github.com/kataras/iris"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("url") + "/api/v1.0/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}


//for test
func newApp() *iris.Application{
	//init the config and log
	if err := config.Init("./conf/config.yaml"); err != nil {
		panic(err)
	}

	//init db
	model.DB.Init()

	//Mae app
	app := iris.Default()

	//register routers to Mae app
	app = router.Load(app)
	return app
}

func main() {

	app:=newApp()

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
			log.Info(app.Run(iris.TLS(viper.GetString("tls.addr"), cert, key),iris.WithoutVersionChecker).Error())
		}()
	}

	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Info(app.Run(iris.Addr(viper.GetString("addr")), iris.WithoutVersionChecker).Error())

	model.DB.RWdb.Close()
}

