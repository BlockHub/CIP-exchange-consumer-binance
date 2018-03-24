package db

type BinanceOrderBook struct {
	ID uint 			`gorm:"primary_key"`
	MarketID uint 		`gorm:"index"`
	Time int64			`gorm:"index"`
}

type BinanceOrder struct {
	ID uint 			`gorm:"primary_key"`
	OrderbookID uint 	`gorm:"index"`
	Rate float64
	//bitfinex supports giving the total number of sell/buyorders.
	//however we should skimp on memory and not add those
	//count int64
	Quantity float64
	Time int64			`gorm:"index"`
}

type BinanceMarket struct {
	ID uint 			`gorm:"primary_key"`
	Ticker string		`gorm:"unique_index:idx_market"`
	Quote string		`gorm:"unique_index:idx_market"`
}

type BinanceTicker struct {
	ID  uint 			`gorm:"primary_key"`
	MarketID uint		`gorm:"index"`
	Price float64
	Volume float64
	Time int64			`gorm:"index"`
}