package main

import (
	"github.com/adshao/go-binance"
	"CIP-exchange-consumer-binance/pkg/consumers"
	"context"
	"github.com/joho/godotenv"
	"github.com/jinzhu/gorm"
	"CIP-exchange-consumer-binance/internal/db"
	"time"
	"strconv"
	"os"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"CIP-exchange-consumer-binance/pkg/handlers"
	"fmt"
	"github.com/getsentry/raven-go"
	"strings"
	"log"
)

var (
	apiKey = ""
	secretKey = ""
)

// get orderbook snapshot, start watching and processing orderbookupdates and ticker updates
func Watch(gormdb gorm.DB, client binance.Client, sym binance.SymbolPrice){
	market := db.CreateOrGetMarket(&gormdb, sym.Symbol[0:3], sym.Symbol[len(sym.Symbol)-3:])
	orderbook := db.CreateOrderBook(&gormdb, market)

	snapshot, err := client.NewDepthService().Symbol(sym.Symbol).Limit(100).Do(context.Background())
	if err != nil{
			if ! strings.Contains(sym.Symbol, "WPR") {
				Watch(gormdb, client, sym)
			} else {
				fmt.Println("weird binance error")
				return
			}
		}

	time := time.Now()
	for _, asks := range snapshot.Asks{
		price, err := strconv.ParseFloat(asks.Price, 64)
		if err != nil{
			log.Panic(err)
		}

		quantity, err := strconv.ParseFloat(asks.Quantity, 64)
		if err != nil{
			log.Panic(err)
		}
		db.AddOrder(&gormdb, price, quantity, time, orderbook)
	}
	consumers.DBConsumer(&gormdb, sym.Symbol, orderbook)
}

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
			raven.CaptureErrorAndWait(err, nil)
		}
	}
	raven.SetDSN(os.Getenv("RAVEN_DSN"))
}


func main() {

	// our local connection
	localdb, err := gorm.Open(os.Getenv("DB"), os.Getenv("DB_URL"))
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
	defer localdb.Close()

	remotedb, err := gorm.Open(os.Getenv("R_DB"), os.Getenv("R_DB_URL"))
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
	defer remotedb.Close()

	db.Migrate(*localdb, *remotedb)

	localdb.DB().SetMaxOpenConns(1000)

	// get the different ticker symbols
	client := binance.NewClient(apiKey, secretKey)
	prices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
	workersleep, err := strconv.ParseInt(os.Getenv("WORKER_SLEEP"), 10, 64)
	tickersleep, err := strconv.ParseInt(os.Getenv("TICKER_SLEEP"), 10, 64)


	// this function can kind of blast the DB.
	for _, p := range prices {
		go Watch(*localdb, *client, *p)
		time.Sleep(time.Duration(workersleep) * time.Millisecond)
	}

	// go Watch needs to create the markets before the ticker handler can start watching them (avoiding a race condition
	// here)
	time.Sleep(10 * time.Second)
	handler := handlers.TickerDbHandler{*localdb}
	// now that markets are created we can send them to the remote
	for true{
		prices, err := client.NewListPricesService().Do(context.Background())
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
		}
		for _, price := range prices{
			handler.Handle(*price)
		}
		time.Sleep(time.Duration(tickersleep) * time.Second)
	}

}
