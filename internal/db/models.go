package db

import "time"

type BinanceOrderBook struct {
	ID uint 			`gorm:"primary_key"`
	MarketID uint64
	Time time.Time		`gorm:"primary_key"`
}

type BinanceOrder struct {
	ID uint64 			`gorm:"primary_key"`
	Time time.Time		`gorm:"primary_key"`
	OrderbookID uint64
	Rate float64
	Quantity float64
}

type BinanceMarket struct {
	ID uint64 			`gorm:"primary_key"`
	Ticker string		`gorm:"unique_index:time_idx_market"`
	Quote string		`gorm:"unique_index:time_idx_market"`
}

type BinanceTicker struct {
	ID  uint64 			`gorm:"primary_key"`
	MarketID uint64
	Price float64
	Time time.Time		`gorm:"primary_key"`
}

type BinanceTrade struct {
	ID  uint64 			`gorm:"primary_key"`
	MarketID uint64
	Price float64
	Quantity float64
	IsBuyerMaker bool
	Time time.Time		`gorm:"primary_key"`
}
