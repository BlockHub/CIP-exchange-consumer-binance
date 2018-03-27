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
)

var (
	apiKey = ""
	secretKey = ""
)

// get orderbook snapshot, start watching and processing orderbookupdates and ticker updates
func Watch(gormdb gorm.DB, client binance.Client, sym binance.SymbolPrice){
	market := db.CreateOrGetMarket(&gormdb, sym.Symbol[0:3], sym.Symbol[len(sym.Symbol)-3:])
	orderbook := db.CreateOrderBook(&gormdb, market)

	fmt.Println(sym.Symbol)
	snapshot, err := client.NewDepthService().Symbol(sym.Symbol).Limit(100).Do(context.Background())
	if err != nil{
		Watch(gormdb, client, sym)
	}

	time := time.Now()
	for _, asks := range snapshot.Asks{
		price, err := strconv.ParseFloat(asks.Price, 64)
		if err != nil{
			panic(err)
		}

		quantity, err := strconv.ParseFloat(asks.Quantity, 64)
		if err != nil{
			panic(err)
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

	gormdb, err := gorm.Open(os.Getenv("DB"), os.Getenv("DB_URL"))
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
	defer gormdb.Close()

	gormdb.AutoMigrate(&db.BinanceMarket{}, &db.BinanceTicker{}, &db.BinanceOrder{}, &db.BinanceOrderBook{})
	err = gormdb.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err = gormdb.Exec("SELECT create_hypertable('binance_orders', 'time', if_not_exists => TRUE)").Error
	if err != nil{
		fmt.Println("binance_orders")
		raven.CaptureErrorAndWait(err, nil)
	}
	err = gormdb.Exec("SELECT create_hypertable('binance_tickers', 'time', if_not_exists => TRUE)").Error
	if err != nil{
		fmt.Println("binance_tickers")
		raven.CaptureErrorAndWait(err, nil)
	}
	err = gormdb.Exec("SELECT create_hypertable('binance_order_books', 'time', if_not_exists => TRUE)").Error
	if err != nil{
		fmt.Println("binance_order_books")
		raven.CaptureErrorAndWait(err, nil)
	}
	gormdb.DB().SetMaxOpenConns(1000)


	// get the different ticker symbols
	client := binance.NewClient(apiKey, secretKey)
	client.Debug = true
	prices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}

	for _, p := range prices {
		go Watch(*gormdb, *client, *p)
	}

	// go Watch needs to create the markets before the ticker handler can start watching them (avoiding a race condition
	// here)
	time.Sleep(10 * time.Second)
	for true{
		prices, err := client.NewListPricesService().Do(context.Background())
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
		}
		for _, price := range prices{
			handler := handlers.TickerDbHandler{*gormdb}
			handler.Handle(*price)
		}
		time.Sleep(60 * time.Second)
	}

}
