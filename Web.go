package main

import (
	"Web/PhFactory"
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmApiResolver"
	"github.com/alfredyang1986/BmServiceDef/BmConfig"
	"github.com/alfredyang1986/BmServiceDef/BmPodsDefine"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
	"github.com/rs/cors"
	"net/http"
	"os"
)

func main() {
	version := "v0"
	prodEnv := "PHARBERS_WEB_HOME"
	fmt.Println("WEB pods archi begins, version =", version)

	fac := PhFactory.PhTable{}
	var pod = BmPodsDefine.Pod{Name: "new PHARBERS_WEB", Factory: fac}
	envHome := os.Getenv(prodEnv)
	pod.RegisterSerFromYAML(envHome + "/resource/service-def.yaml")

	var bmRouter BmConfig.BmRouterConfig
	bmRouter.GenerateConfig(prodEnv)

	addr := bmRouter.Host + ":" + bmRouter.Port
	fmt.Println("Listening on ", addr)
	api := api2go.NewAPIWithResolver(version, &BmApiResolver.RequestURL{Addr: addr})
	pod.RegisterAllResource(api)

	pod.RegisterAllFunctions(version, api)
	pod.RegisterAllMiddleware(api)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})

	handler := api.Handler().(*httprouter.Router)

	pod.RegisterPanicHandler(handler)
	http.ListenAndServe(":"+bmRouter.Port, c.Handler(handler))

	fmt.Println("WEB pods archi ends, version =", version)
}
