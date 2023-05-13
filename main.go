package main

import (
	"main/cmd/cowrieParser"
	"main/cmd/db"
	"fmt"
	"main/cmd"
	"main/cmd/fileDetection"
)

func importAll(domArray []db.DomainDescription, pgDB *db.DB) {
	for _, dom := range domArray {
		if dom.Domain == "" {
			dom.Domain = cmd.IntToIPv4(dom.IP).String()
		}
		x := pgDB.Database.Save(&dom)
		if x.Error != nil {
			fmt.Println(x.Error)
		}
	}
}

func main() () {
	var pgDB *db.DB = db.InitDB(db.IPDBNAME)

	var filePath []string = []string{
		"/root/a.c",
		"data_files/cowrie.json",
		"data_files/cowrie.json.2023-04-30.json",
		"data_files/cowrie.json.2023-05-05",
		"data_files/cowrie.json.2023-05-09",
	}

	//x := fileDetection.GetAllFiles("/data/data/")
	for _, filepath := range filePath {
		fmt.Printf("Reading file: %s\n", filepath)
		domArray, err := cowrieParser.ReadFile(filepath)
		if err != nil {
			fmt.Printf("Error while reading file: %s | %s", filepath, err)
			continue
		}
		fmt.Println(filepath, "Len of domArray: ", len(domArray))
		importAll(domArray, pgDB)
	}

}

func maina() {
	x := fileDetection.GetAllFiles("/data/data/")
	for _, filepath := range x {
		y := cmd.GetChunk(filepath)
		a := cmd.HashChunk(y)
		fmt.Printf("%x\n", a)
	}
}
