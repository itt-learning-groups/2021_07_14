package main

import (
	"fmt"
	"github.com/us-learn-and-devops/2021-07-01/todo-api/cmd/todo_api_server/handlers"
	"github.com/us-learn-and-devops/2021-07-01/todo-api/configs"
	envcfg "github.com/us-learn-and-devops/2021-07-01/todo-api/pkg/envconfig"
	"log"
	"net/http"
	"os"
)

func main() {
	cfgs := &configs.Settings{}
	err := envcfg.Unmarshal(cfgs)
	if err != nil {
		fmt.Printf("Failed to get configs: %s", err)
		os.Exit(1)
	}

	tl := handlers.NewTodoListHandler()
	r := handlers.NewRouter(tl)

	host := fmt.Sprintf("%s:%s", cfgs.ServerAddr, cfgs.ServerPort)

	fmt.Printf("serving todo-api on %s", host)
	log.Fatal(http.ListenAndServe(host, r))
}
