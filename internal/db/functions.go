package db

import (
	"github.com/jinzhu/gorm"
	"time"
	"strings"
	"log"
)

func CreateOrderBook (db *gorm.DB, market BinanceMarket) BinanceOrderBook{
	// since ID is zero, GORM will override the value and auto increment it.
	orderbook := BinanceOrderBook{0,market.ID, time.Now()}
	err := db.Create(&orderbook).Error
	if err != nil{
		log.Panic(err)
	}
	return orderbook
}


func AddOrder(db *gorm.DB, rate float64, quantity float64, time time.Time, orderbook BinanceOrderBook) BinanceOrder{
	order := BinanceOrder{0,  time, uint64(orderbook.ID), rate, quantity,}
	err := db.Create(&order).Error
	if err != nil{
		log.Panic(err)
	}
	return order
}

func CreateOrGetMarket(db *gorm.DB, ticker string, quote string) BinanceMarket{
	market := BinanceMarket{0, ticker, quote}
	err := db.Create(&market).Error
	if err != nil{
		if ! strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			log.Panic(err)
		}
	}
	return market
}

func AddTicker(db *gorm.DB, market BinanceMarket, price float64){
	ticker := BinanceTicker{0, market.ID, price, time.Now()}
	err := db.Create(&ticker).Error
	if err != nil {
		log.Panic(err)
	}
}

func AddTrade (db *gorm.DB, market BinanceMarket, id uint64, price float64, quantity float64, isbuyermaker bool){
	Trade := BinanceTrade{ID:id, Price:price, Quantity:quantity, MarketID:market.ID, IsBuyerMaker:isbuyermaker, Time:time.Now()}
	err := db.Create(&Trade).Error
	if err != nil{
		log.Panic(err)
	}
}
