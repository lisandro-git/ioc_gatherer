package main

import (
	"main/cmd/cowrieParser"
)

//func getSRV(domain string) string {
//	srv, _ := net.LookupSRV(domain)
//	return srv
//}
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
	cowrieParser.ReadFile("data_files/cowrie.json")
	//cowrieParser.ReadFile("data_files/cowrie.json.2023-04-30.json")

	//http.HandleFunc("/", rootPage)
	//http.ListenAndServe(":80", nil)

}
