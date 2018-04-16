package handlers

import (
	"github.com/adshao/go-binance"
	"fmt"
	"CIP-exchange-consumer-binance/internal/db"
	"github.com/jinzhu/gorm"
	"time"
	"strconv"
	"log"
)

type PrintHandler struct{
}
	func (p PrintHandler)Handle(event *binance.WsDepthEvent){
		fmt.Println(event)
	}

type ErrHandler struct{
}
	func (e ErrHandler)Handle(err error) {
		fmt.Println(err)
	}

type OrderDbHandler struct{
	Db gorm.DB
	Orderbook db.BinanceOrderBook
}
	func (odb OrderDbHandler) Handle(event *binance.WsDepthEvent){
		for _, ask := range event.Asks {
			price, err := strconv.ParseFloat(ask.Price, 64)
			if err != nil{
				log.Panic(err)
			}

			quantity, err := strconv.ParseFloat(ask.Quantity, 64)
			if err != nil{
				log.Panic(err)
			}

			db.AddOrder(&odb.Db, price, quantity, time.Now(),  odb.Orderbook, false)
		}

		for _, ask := range event.Bids {
			price, err := strconv.ParseFloat(ask.Price, 64)
			if err != nil{
				log.Panic(err)
			}

			quantity, err := strconv.ParseFloat(ask.Quantity, 64)
			if err != nil{
				log.Panic(err)
			}
			db.AddOrder(&odb.Db, price, quantity, time.Now(),  odb.Orderbook, true)
		}
	}

type TickerDbHandler struct{
	Db gorm.DB
}
func (t TickerDbHandler) Handle(price binance.SymbolPrice) {
	priceflt, err := strconv.ParseFloat(price.Price, 64)
	if err != nil{
		log.Panic(err)
	}
	market := db.BinanceMarket{}
	res := t.Db.Where(map[string]interface{}{
		"ticker": price.Symbol[0:3],
		"quote": price.Symbol[len(price.Symbol)-3:]}).Find(&market)
	if res.Error != nil{
		log.Panic(err)
	}

	db.AddTicker(&t.Db, market, priceflt)
}

type TradeDbHandler struct {
	Db gorm.DB
	Market db.BinanceMarket
}

func (t TradeDbHandler) Handle(event *binance.WsAggTradeEvent) {
	price, err := strconv.ParseFloat(event.Price, 64)
	if err != nil{
		log.Panic(err)
	}

	quantity, err := strconv.ParseFloat(event.Quantity, 64)
	if err != nil{
		log.Panic(err)
	}
	db.AddTrade(&t.Db, t.Market, uint64(event.AggTradeID), price, quantity, event.IsBuyerMaker)
}