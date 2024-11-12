// @title ToDo API
// @version 1.0
// @description This is a simple ToDo API application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	_ "todo-app/docs"
)

var (
	buildTime string
	version   string
)

type config struct {
	port int
	env  string
}
type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	wg       sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	displayVersion := flag.Bool("version", false, "Display version information and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("ToDo-app version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	app := &application{
		config:   cfg,
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	if err := app.startServerWithGracefulShutdown(); err != nil {
		app.errorLog.Fatal(err)
	}
}
