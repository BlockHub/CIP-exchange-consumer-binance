package main

import (
"os"
"fmt"
"strconv"
"CIP-exchange-consumer-binance/pkg/pushers"
"log"
"github.com/jinzhu/gorm"
"github.com/getsentry/raven-go"
_ "github.com/jinzhu/gorm/dialects/postgres"
"github.com/joho/godotenv"
)

func init(){
	useDotenv := true
	if os.Getenv("PRODUCTION") == "true"{
		useDotenv = false
	}

	// this loads all the constants stored in the .env file (not suitable for production)
	// set variables in supervisor then.
	if useDotenv {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
	}
	raven.SetDSN(os.Getenv("RAVEN_DSN"))
}


func main(){
	localdb, err := gorm.Open(os.Getenv("DB"), os.Getenv("DB_URL"))
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	defer localdb.Close()

	remotedb, err := gorm.Open(os.Getenv("R_DB"), os.Getenv("R_DB_URL"))
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
	defer remotedb.Close()

	fmt.Println("Starting replication worker")
	strlimit := os.Getenv("REPLICATION_LIMIT")
	limit, err := strconv.ParseInt(strlimit, 10, 64)
	fmt.Println(limit)
	if err != nil{
		log.Panic(err)
	}
	replicator := pushers.Replicator{Local:*localdb, Remote:*remotedb, Limit:limit,
		Name:os.Getenv("NAME"), DBlink:os.Getenv("DBlink")}
	replicator.Link()
	defer replicator.Unlink()
	fmt.Println("replicating")
	replicator.PushMarkets()
	go replicator.Start()
	select {}
}