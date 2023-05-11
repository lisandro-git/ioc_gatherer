package main

import (
	"main/cmd/cowrieParser"
	"main/cmd/db"
	"fmt"
	"main/cmd"
	"strconv"
	"reflect"
)

func main() () {
	//http.HandleFunc("/", rootPage)
	//http.ListenAndServe(":80", nil)
	domArray := cowrieParser.ReadFile("data_files/cowrie.json.2023-05-05")
	//cowrieParser.ReadFile("data_files/cowrie.json.2023-04-30.json")
	fmt.Println("Len of domArray: ", len(domArray))

	var pgDB *db.DB = db.InitDB()
	for i, dom := range domArray {
		fmt.Println(i)
		if dom.Domain == "" {
			dom.Domain = cmd.IntToIPv4(dom.IP).String()
		}
		intVar, err := strconv.Atoi(dom.Domain)
		fmt.Println(intVar, err, reflect.TypeOf(intVar))
		x := pgDB.Database.Save(&dom)
		if x.Error != nil {
			fmt.Println(x)
		}
		fmt.Println(x)

	}

}
