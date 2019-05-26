package main

import (
	"fmt"

	"go_module/hello"

	"github.com/leekchan/accounting"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	fmt.Println(ac.FormatMoney(123.123))
	hello.Sayhello()
}
