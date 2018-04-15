package db

import (
	"github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"

)

func Migrate(Local gorm.DB, Remote gorm.DB){
	// migrations are only performed by GORM if a table/column/index does not exist.
	err := Local.AutoMigrate(	&BinanceMarket{},
		&BinanceOrder{},
		//&BinanceTicker{},
		&BinanceTrade{},
		&BinanceOrderBook{}).Error

	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err = Local.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}

	err = Local.Exec("CREATE EXTENSION IF NOT EXISTS dblink;").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err = Local.Exec("SELECT create_hypertable('bitfinex_orders', 'time',  'orderbook_id', if_not_exists => TRUE)").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err = Local.Exec("SELECT create_hypertable('bitfinex_tickers', 'time', 'market_id', if_not_exists => TRUE)").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err =Local.Exec("SELECT create_hypertable('bitfinex_order_books', 'time', 'market_id', if_not_exists => TRUE)").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}

	err = Remote.AutoMigrate(	&BinanceMarket{},
								&BinanceOrder{},
								//&BinanceTicker{},
								&BinanceTrade{},
								&BinanceOrderBook{}).Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err = Remote.Exec("CREATE EXTENSION IF NOT EXISTS dblink;").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}

	err = Remote.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err = Remote.Exec("SELECT create_hypertable('bitfinex_orders', 'time',  'orderbook_id', if_not_exists => TRUE)").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err = Remote.Exec("SELECT create_hypertable('bitfinex_tickers', 'time', 'market_id', if_not_exists => TRUE)").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err =Remote.Exec("SELECT create_hypertable('bitfinex_order_books', 'time', 'market_id', if_not_exists => TRUE)").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
}