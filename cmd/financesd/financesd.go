package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // because
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/dao/accountdao"
)

func main() {
	db, err := gorm.Open("mysql", "user:password@tcp(host)/schema?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	db.SingularTable(true)

	accounts, err := accountdao.GetAll(db)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	for index := 0; index < len(accounts); index++ {
		log.Printf("%3d '%s' %v, %v\n", accounts[index].ID, accounts[index].Name, accounts[index].Closed, accounts[index].CompanyID)
	}
}
