package pushers

import (
	"github.com/jinzhu/gorm"
	"log"
	"CIP-exchange-consumer-binance/internal/db"
	"fmt"
	"time"
)



type Replicator struct {
	// local db
	Local gorm.DB

	//remote DB (the data warehouse)
	Remote gorm.DB

	//schema related settings

	//replication related settings
	Limit int	// max rows to be fetched from remote and inserted (should be as high as possible)

}
// copy the markets table (should only be done once in a while, as new markets
// are only added once every few months.
func(r *Replicator) Start(){
	for {
		r.Replicate_ticker()
	}
}

// copy the ticker data from a chunk
func (r *Replicator) Replicate_ticker() {
	// an out interface to store lots of Order objects
	backup := r.Remote.Begin()
	local := r.Local.Begin()

	orders := []db.BinanceOrder{}
	tickers := []db.BinanceTicker{}

	if (len(orders) == 0) || (len(tickers) == 0){
		time.Sleep(10* time.Second)
		return
	}

	r.Local.Limit(r.Limit).Find(&orders)
	r.Local.Limit(r.Limit).Find(&tickers)


	for _, order := range orders {
		err := backup.Create(&order).Error
		if err != nil{
			panic(err)
		}
		err = local.Delete(&order).Error
		if err != nil{
			panic(err)
		}
	}

	for _, ticker := range tickers {
		err := backup.Create(&ticker).Error
		if err != nil{
			panic(err)
		}
		err = local.Delete(&ticker).Error
		if err != nil{
			panic(err)
		}
	}

	fmt.Println("committing")
	err := backup.Commit().Error
	if err != nil{
		local.Rollback()
		backup.Rollback()
		log.Panic(err)
	}

	err = local.Commit().Error
	if err != nil{
		local.Rollback()
		log.Panic(err)
	}
}