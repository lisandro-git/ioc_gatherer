package cowrieParser

import (
	"os"
	"bufio"
	"github.com/buger/jsonparser"
	"github.com/oschwald/geoip2-golang"
	"net"
	"log"
	"strconv"
	"main/cmd/IPData"
	"main/cmd/db"
	"gorm.io/gorm"
	"fmt"
)

func getCity(src_ip string) (string, string, string) {
	db, err := geoip2.Open("/root/Downloads/GeoLite2-City_20230428/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	ip := net.ParseIP(src_ip)
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Coordinates: ", record.Location.Latitude, record.Location.Longitude)

	var lat string = strconv.FormatFloat(record.Location.Latitude, 'f', 6, 64)
	var long string = strconv.FormatFloat(record.Location.Longitude, 'f', 6, 64)
	return lat, long, record.City.Names["en"]
}

func getCountry(src_ip string) (string, string) {
	db, err := geoip2.Open("/root/Downloads/GeoLite2-City_20230421/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP(src_ip)
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	return record.Country.Names["en"], record.Country.IsoCode
}

func getdata(ip string) IPData.DNSRecord {
	return IPData.NewDnsRecord(nil, []string{ip})
}

func getID(x *gorm.DB, ip string) int {
	y := &db.SourceIPDescription{}
	a, _ := x.Model(&db.SourceIPDescription{}).Select("id").Where("src_ip = ?", ip).Rows()
	for a.Next() {
		x.ScanRows(a, y)
	}
	fmt.Printf("Updated data to database: ID: %d, IP: %s", y.ID, ip)
	return y.ID
}

func ReadFile(filePath string) []db.SourceIPDescription {
	data, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	var fileScanner *bufio.Scanner = bufio.NewScanner(data)
	fileScanner.Split(bufio.ScanLines)

	var sipArray []db.SourceIPDescription

	for fileScanner.Scan() {
		var redirect bool = false
		var sip *db.SourceIPDescription = db.NewSourceIPDescription()
		var ipdata *db.IPDataDescription = db.NewIPDataDescription()
		//var country *db.CountryDescription = db.NewCountryDescription()
		//var hittime *db.HitTimeDescription = db.NewHitTimeDescription()

		eventID, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "eventid")
		if string(eventID) != "cowrie.session.connect" {
			continue
		}
		srcIP, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "src_ip")
		proto, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "protocol")
		timestamp, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "timestamp")
		sip.CountryName, sip.CountryCode = getCountry(string(srcIP))

		{
			sip.Time = string(timestamp)
			sip.Protocol = string(proto)
			sip.EventID = string(eventID)
			sip.SRCIP = string(srcIP)
			sip.Latitude, sip.Longitude, sip.City = getCity(string(srcIP))
		}

		for i := 0; i < len(sipArray); i++ {
			// incrementing hit counts for each IP
			if sipArray[i].SRCIP == string(srcIP) {
				sipArray[i].HitCount++
				// The timestamp is the last time the IP was seen
				//countryArray[i].SourceIP[j].Timestamp = string(timestamp)
				redirect = true
				break
			}

		}
		if redirect {
			continue
		}
		sip.HitCount++
		ipdata.PrepareForDB(getdata(sip.SRCIP))

		sip.IPData = *ipdata
		sipArray = append(sipArray, *sip)
	}

	x := db.ConnectToDB()
	db.BuildDatabase(x)
	for i := 1; i < len(sipArray); i++ {
		fmt.Println("IP: ", sipArray[i].SRCIP)
		//fmt.Println(errors.Is(err, gorm.ErrRecordNotFound))r(
		a := x.Save(&sipArray[i])
		//sipArray[i].Latitude = "123"

		// filter if error = duplicate
		if a.Error != nil {
			z := getID(x, sipArray[i].SRCIP)
			// update data
			x.Model(&db.SourceIPDescription{}).Where("id = ?", z).Update("country_name", sipArray[i].CountryName)
			x.Model(&db.SourceIPDescription{}).Where("id = ?", z).Update("Country_Code", sipArray[i].CountryCode)
			x.Model(&db.SourceIPDescription{}).Where("id = ?", z).Update("Latitude", sipArray[i].Latitude)
			x.Model(&db.SourceIPDescription{}).Where("id = ?", z).Update("Longitude", sipArray[i].Longitude)
			x.Model(&db.SourceIPDescription{}).Where("id = ?", z).Update("City", sipArray[i].City)
			x.Model(&db.SourceIPDescription{}).Where("id = ?", z).Update("Hit_Count", sipArray[i].HitCount)
			x.Model(&db.SourceIPDescription{}).Where("id = ?", z).Update("Time", sipArray[i].Time)
			x.Model(&db.SourceIPDescription{}).Where("id = ?", z).Update("Protocol", sipArray[i].Protocol)
			x.Model(&db.SourceIPDescription{}).Where("id = ?", z).Update("Event_ID", sipArray[i].EventID)
		}
	}
	return sipArray
}
