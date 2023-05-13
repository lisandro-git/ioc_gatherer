package main

import (
	"main/cmd/cowrieParser"
	"main/cmd/db"
	"fmt"
	"main/cmd"
)

func importAll(domArray []db.DomainDescription, pgDB *db.DB) {
	for i, dom := range domArray {
		fmt.Println(i)
		if dom.Domain == "" {
			dom.Domain = cmd.IntToIPv4(dom.IP).String()
		}
		//intVar, err := strconv.Atoi(dom.Domain)
		//fmt.Println(intVar, err, reflect.TypeOf(intVar))
		x := pgDB.Database.Save(&dom)
		if x.Error != nil {
			fmt.Println(x.Error)
		}
	}
}

func main() () {
	var pgDB *db.DB = db.InitDB()

	domArray := cowrieParser.ReadFile("data_files/cowrie.json")
	fmt.Println("cowrie.json Len of domArray: ", len(domArray))
	importAll(domArray, pgDB)

	domArray = cowrieParser.ReadFile("data_files/cowrie.json.2023-04-30.json")
	fmt.Println("cowrie.json.2023-04-30.json Len of domArray: ", len(domArray))
	importAll(domArray, pgDB)

	domArray = cowrieParser.ReadFile("data_files/cowrie.json.2023-05-05")
	fmt.Println("cowrie.json.2023-05-05 Len of domArray: ", len(domArray))
	importAll(domArray, pgDB)

	domArray = cowrieParser.ReadFile("data_files/cowrie.json.2023-05-09")
	fmt.Println("cowrie.json.2023-05-09 Len of domArray: ", len(domArray))
	importAll(domArray, pgDB)
}
