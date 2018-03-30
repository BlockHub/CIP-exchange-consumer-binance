package db

import "time"

type OrderBook struct {
	ID uint 			`gorm:"primary_key"`
	MarketID uint
	Time time.Time		`gorm:"primary_key"`
}

type Order struct {
	ID uint 			`gorm:"primary_key"`
	Time time.Time		`gorm:"primary_key"`
	OrderbookID uint
	Rate float64
	Quantity float64
}

type Market struct {
	ID uint 			`gorm:"primary_key"`
	Ticker string		`gorm:"unique_index:time_idx_market"`
	Quote string		`gorm:"unique_index:time_idx_market"`
}

type Ticker struct {
	ID  uint 			`gorm:"primary_key"`
	MarketID uint
	Price float64
	Time time.Time		`gorm:"primary_key"`
}