package handlers

import (
	"github.com/adshao/go-binance"
	"fmt"
	"CIP-exchange-consumer-binance/internal/db"
	"github.com/jinzhu/gorm"

	"time"
	"strconv"
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
				panic(err)
			}

			quantity, err := strconv.ParseFloat(ask.Quantity, 64)
			if err != nil{
				panic(err)
			}

			db.AddOrder(&odb.Db, price, quantity, int64(time.Now().Unix()),  odb.Orderbook)
		}
	}

type TickerDbHandler struct{
	Db gorm.DB
}
func (t TickerDbHandler) Handle(price binance.SymbolPrice) {
	priceflt, err := strconv.ParseFloat(price.Price, 64)
	if err != nil{
		panic(err)
	}
	market := db.BinanceMarket{}
	res := t.Db.Where(map[string]interface{}{
		"ticker": price.Symbol[0:3],
		"quote": price.Symbol[len(price.Symbol)-3:]}).Find(&market)
	if res.Error != nil{
		panic(res.Error)
	}

	db.AddTicker(&t.Db, market, priceflt)

}