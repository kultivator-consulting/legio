package main

import (
	"context"
	"cortex_api/database"
	"cortex_api/services/file_service/config"
	"cortex_api/services/file_service/utils"
	"flag"
	"github.com/joho/godotenv"
	"go.etcd.io/etcd/client/pkg/v3/fileutil"
	"go.uber.org/fx"
	"log"
	"os"
	"time"
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

func FileListContains(fileList []string, file string) bool {
	for _, f := range fileList {
		if utils.GetFileStorageName(f) == file {
			return true
		}
	}
	return false
}

func PurgeStaleFiles() {
	_, queries, ctx, err := database.ApiDatabase().Open()
	if err != nil {
		log.Fatalf("HouseKeeping DB connection error: %v\n", err)
	}

	defer func(database *database.Model) {
		err := database.Close()
		if err != nil {
			log.Fatalf("HouseKeeping DB close error: %v\n", err)
		}
	}(database.ApiDatabase())

	purgeInterval, err := time.ParseDuration(os.Getenv("FILE_STORE_PURGE_INTERVAL"))
	if err != nil {
		log.Printf("Error while parsing file store purge interval: %v\n", err)
	}

	for {
		log.Printf("Running file store purge\n")

		fileListRecords, err := queries.ListStoredFiles(ctx)
		if err != nil {
			log.Fatalf("Error while getting list of active stored files: %v\n", err)
		}

		storedFiles, err := os.ReadDir(os.Getenv("FILE_STORE_PATH"))
		if err != nil {
			log.Fatalf("Error while reading file storage directory: %v\n", err)
		}

		for _, file := range storedFiles {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			if !FileListContains(fileListRecords, fileName) {
				log.Printf("File %s exists in file storage but not in database, removing from file storage\n", fileName)
				err := os.Remove(utils.GetFileStoragePath(fileName))
				if err != nil {
					log.Printf("Error while removing file %s: %v\n", fileName, err)
				}
			}
		}

		time.Sleep(purgeInterval)
	}
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

	err = os.Mkdir(os.Getenv("FILE_STORE_PATH"), 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	go func() {
		PurgeStaleFiles()
	}()

	app := fx.New(
		Module,
		fx.Invoke(RegisterApiServer),
	)
	app.Run()
}
