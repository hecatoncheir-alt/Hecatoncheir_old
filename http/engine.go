package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// NewEngine is a constructor for HTTPEngine
func NewEngine(apiVersion string) *Engine {
	router := httprouter.New()
	httpEngine := Engine{APIVersion: apiVersion, Router: router}
	return &httpEngine
}

// Engine is a object with http settings
type Engine struct {
	APIVersion string
	Server     *http.Server
	Router     *httprouter.Router
}

func (httpEngine *Engine) PowerUp(host string, port int) {
	httpEngine.Router.GET("/api/version", httpEngine.apiVersionCheckHandler)

	httpEngine.Server = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}
	fmt.Printf("Http server listen on %v, port:%v \n", host, port)

	httpEngine.Server.Handler = httpEngine.Router
	httpEngine.Server.ListenAndServe()
}

func (httpEngine *Engine) apiVersionCheckHandler(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	data := map[string]string{"apiVersion": httpEngine.APIVersion}
	encodedData, _ := json.Marshal(data)

	response.Header().Set("content-type", "application/javascript")
	_, err := response.Write(encodedData)
	if err != nil {
		fmt.Print(err)
	}
}
