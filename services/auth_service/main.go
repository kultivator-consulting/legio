package main

import (
	"context"
	"cortex_api/services/auth_service/config"
	"flag"
	"github.com/joho/godotenv"
	"go.etcd.io/etcd/client/pkg/v3/fileutil"
	"go.uber.org/fx"
	"log"
	"os"
)

var spyGodotenvLoad = godotenv.Load

func LoadEnv() error {
	// load specific environment file, if it exists
	env, valid := os.LookupEnv("APP_ENV")
	if valid && env != "" && fileutil.Exist(".env."+env) {
		return spyGodotenvLoad(".env." + env)
	}
	// load fallback environment file
	if fileutil.Exist(".env") {
		return spyGodotenvLoad(".env")
	}

	return nil
}

func RegisterApiServer(
	lifeCycle fx.Lifecycle,
	webServer *Model,
) {
	lifeCycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := webServer.StartServer(); err != nil {
					log.Fatalf("Unable to start API server. Error : %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Println("Stopping API server...")
			return webServer.StopServer()
		},
	})
}

func main() {
	configPathPtr := flag.String("config", "./config.yaml", "The path to the configuration (config.yaml) file for the service")
	flag.Parse()

	logFlags := log.LstdFlags | log.LUTC | log.Lshortfile
	log.SetFlags(logFlags)

	err := LoadEnv()
	if err != nil {
		panic(err)
	}

	err = config.AppConfig().InitializeConfig(*configPathPtr)
	if err != nil {
		panic(err)
	}

	app := fx.New(
		Module,
		fx.Invoke(RegisterApiServer),
	)
	app.Run()
}
