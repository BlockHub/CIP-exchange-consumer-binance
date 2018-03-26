package db

import "time"

type BinanceOrderBook struct {
	ID uint 			`gorm:"primary_key"`
	MarketID uint 		`gorm:"index"`
	Time time.Time		`gorm:"primary_key"`
}

type BinanceOrder struct {
	ID uint 			`gorm:"primary_key"`
	Time time.Time		`gorm:"primary_key"`
	OrderbookID uint 	`gorm:"index:book"`
	Rate float64
	Quantity float64
}

type BinanceMarket struct {
	ID uint 			`gorm:"primary_key"`
	Ticker string		`gorm:"unique_index:time_idx_market"`
	Quote string		`gorm:"unique_index:time_idx_market"`
}

type BinanceTicker struct {
	ID  uint 			`gorm:"primary_key"`
	MarketID uint		`gorm:"index"`
	Price float64
	Time time.Time		`gorm:"primary_key"`
}