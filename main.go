package main

import (
	"github.com/adshao/go-binance"
	"CIP-exchange-consumer-binance/pkg/consumers"
	"context"
	"github.com/joho/godotenv"
	"log"
	"github.com/jinzhu/gorm"
	"CIP-exchange-consumer-binance/internal/db"
	"time"
	"strconv"
	"os"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"CIP-exchange-consumer-binance/pkg/handlers"
)

var (
	apiKey = ""
	secretKey = ""
)

// get orderbook snapshot, start watching and processing orderbookupdates and ticker updates
func Watch(gormdb gorm.DB, client binance.Client, sym binance.SymbolPrice){
	market := db.CreateOrGetMarket(&gormdb, sym.Symbol[0:3], sym.Symbol[len(sym.Symbol)-3:])
	orderbook := db.CreateOrderBook(&gormdb, market)

	snapshot, err := client.NewDepthService().Symbol(sym.Symbol).Do(context.Background())
	if err != nil{
		panic(err)
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


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gormdb, err := gorm.Open(os.Getenv("DB"), os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	defer gormdb.Close()

	gormdb.AutoMigrate(&db.BinanceMarket{}, &db.BinanceTicker{}, &db.BinanceOrder{}, &db.BinanceOrderBook{})
	err = gormdb.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;").Error
	if err != nil{
		panic(err)
	}
	err = gormdb.Exec("SELECT create_hypertable('binance_orders', 'time', if_not_exists => TRUE)").Error
	if err != nil{
		panic(err)
	}
	err = gormdb.Exec("SELECT create_hypertable('binance_tickers', 'time', if_not_exists => TRUE)").Error
	if err != nil{
		panic(err)
	}
	err = gormdb.Exec("SELECT create_hypertable('binance_order_books', 'time', if_not_exists => TRUE)").Error
	if err != nil{
		panic(err)
	}
	gormdb.DB().SetMaxOpenConns(1000)


	// get the different ticker symbols
	client := binance.NewClient(apiKey, secretKey)
	prices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		panic(err)
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
			panic(err)
		}
		for _, price := range prices{
			handler := handlers.TickerDbHandler{*gormdb}
			handler.Handle(*price)
		}
		time.Sleep(60 * time.Second)
	}

}
