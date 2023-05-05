package main

import (
	"main/cmd/cowrieParser"
	"fmt"
	"main/cmd/db"
)

/*
func getIOCData() *cowrieParser.IOCs {
	return cowrieParser.ReadFile("data_files/cowrie.json")
}

func rootPage(w http.ResponseWriter, req *http.Request) {
	var tmpl *template.Template = template.Must(template.ParseFiles("templates/index.html"))
	var fs http.Handler = http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Gathering data...")
	var data *cowrieParser.IOCs = getIOCData()
	fmt.Println("Gathering data... Done!")

	tmpl.Execute(w, data)
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	f, err := os.Create("data_files/x.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.Write(b)
}

*/

func main() () {
	//http.HandleFunc("/", rootPage)
	//http.ListenAndServe(":80", nil)

	sipArray := cowrieParser.ReadFile("data_files/cowrie.json")
	//cowrieParser.ReadFile("data_files/cowrie.json.2023-04-30.json")

	var pgDB *db.DB = db.InitDB()
	for i := 0; i < len(sipArray); i++ {
		fmt.Println("IP: ", sipArray[i].SRCIP)
		a := pgDB.Database.Save(&sipArray[i])

		// filter if error = duplicate
		if a.Error != nil {
			pgDB.UpdateDBRow(sipArray[i])
		}
	}

}
