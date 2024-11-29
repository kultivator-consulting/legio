package main

import (
	"cortex_api/database"
	"github.com/joho/godotenv"
	"go.etcd.io/etcd/client/pkg/v3/fileutil"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

const DbSeedFolder = "seeds"

func LoadEnv() error {
	// load specific environment file, if it exists
	env, valid := os.LookupEnv("APP_ENV")
	if valid && env != "" && fileutil.Exist(".env."+env) {
		return godotenv.Load(".env." + env)
	}
	// load fallback environment file
	if fileutil.Exist(".env") {
		return godotenv.Load(".env")
	}

	return nil
}

func RunDatabaseSeed(seedSource string) error {
	db, _, ctx, err := database.ApiDatabase().Open()
	if err != nil {
		log.Fatalf("DatabaseSeed DB connection error: %v\n", err)
		return err
	}

	defer func(database *database.Model) {
		err := database.Close()
		if err != nil {
			log.Fatalf("DatabaseSeed DB close error: %v\n", err)
			return
		}
	}(database.ApiDatabase())

	if seedSource == "" {
		log.Printf("No seeds data loaded. Manually load seed data if needed. For example './scripts/migrate_up.sh'\n")
		return nil
	}

	log.Printf("Running database seeding for %s\n", seedSource)

	path := filepath.Join(DbSeedFolder, seedSource)

	seedFiles, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Error while reading seeds directory: %v\n", err)
		return err
	}

	sort.Slice(seedFiles, func(i, j int) bool {
		if seedFiles[i].Name()[4:] != seedFiles[j].Name()[4:] {
			return seedFiles[i].Name() < seedFiles[j].Name()
		}
		ii, _ := strconv.Atoi(seedFiles[i].Name()[:4])
		jj, _ := strconv.Atoi(seedFiles[j].Name()[:4])
		return ii < jj
	})

	for _, file := range seedFiles {
		if file.IsDir() {
			continue
		}

		fileName := filepath.Join(path, file.Name())
		seedFileData, ioErr := ioutil.ReadFile(fileName)
		if ioErr != nil {
			log.Printf("Error while reading seed file %s: %v\n", fileName, ioErr)
			return ioErr
		}

		sql := string(seedFileData)
		_, err := db.Exec(ctx, sql)
		if err != nil {
			log.Printf("Error while executing seed file %s: %v\n", fileName, err)
			return err
		}

		log.Printf("Executed seed file %s\n", fileName)
	}

	return nil
}

func main() {
	logFlags := log.LstdFlags | log.LUTC | log.Lshortfile
	log.SetFlags(logFlags)

	err := LoadEnv()
	if err != nil {
		panic(err)
	}

	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("Invalid number of arguments. Second argument must be the seed folder to apply to the DB.\n")
		return
	}
	err = RunDatabaseSeed(args[0])
	if err != nil {
		panic(err)
	}
}
