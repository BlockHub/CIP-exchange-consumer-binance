package consumers

import (
	"github.com/adshao/go-binance"
	"CIP-exchange-consumer-binance/pkg/handlers"
	"CIP-exchange-consumer-binance/internal/db"
	"github.com/jinzhu/gorm"
	"log"
)

func PrintConsumer(symbol string){
	handler := handlers.PrintHandler{}
	errhandler := handlers.ErrHandler{}
	binance.WsDepthServe(symbol, handler.Handle, errhandler.Handle)
}

func DBConsumer(gorm *gorm.DB, symbol string, book db.BinanceOrderBook, market db.BinanceMarket){
	orderhandler := handlers.OrderDbHandler{*gorm, book}
	tradehandler := handlers.TradeDbHandler{*gorm, market}
	errhandler := handlers.ErrHandler{}

	_, _, err := binance.WsDepthServe(symbol, orderhandler.Handle, errhandler.Handle)
	if err != nil{
		log.Panic(err)
	}
	_, _, err = binance.WsAggTradeServe(symbol, tradehandler.Handle, errhandler.Handle)
	if err != nil {
		log.Panic(err)
	}

}