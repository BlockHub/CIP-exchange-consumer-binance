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
		fmt.Println(event)
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