package consumers

import (
	"github.com/adshao/go-binance"
	"CIP-exchange-consumer-binance/pkg/handlers"
	"CIP-exchange-consumer-binance/internal/db"
	"github.com/jinzhu/gorm"
)

func PrintConsumer(symbol string){
	handler := handlers.PrintHandler{}
	errhandler := handlers.ErrHandler{}
	binance.WsDepthServe(symbol, handler.Handle, errhandler.Handle)
}

func DBConsumer(gorm *gorm.DB, symbol string, book db.BinanceOrderBook){
	handler := handlers.OrderDbHandler{*gorm, book}
	errhandler := handlers.ErrHandler{}

	_, _, err := binance.WsDepthServe(symbol, handler.Handle, errhandler.Handle)
	if err != nil{
		panic(err)
	}
}