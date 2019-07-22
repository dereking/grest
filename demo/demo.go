package main

import (
	"flag"
	//"log"
	"reflect"

	"github.com/dereking/grest/log"
	"go.uber.org/zap"

	"github.com/dereking/grest"
	"github.com/dereking/grest/demo/controllers"
)

func main() {
	conf := flag.String("conf", "app.conf", "the conf file in conf DIR for this server.")
	flag.Parse()

	log.Logger().Info("Starting server with config file :", zap.Any("conf", *conf))

	s := grest.NewGrestServer(*conf)

	//controller register
	s.AddController("api", reflect.TypeOf(controllers.ApiController{}))
	s.AddController("Home", reflect.TypeOf(controllers.HomeController{}))

	//main loop
	s.Serve()
}
