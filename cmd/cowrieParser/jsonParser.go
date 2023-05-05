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
)

func getGEOIPData(src_ip string) (string, string, string, string) {
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
	var lat string = strconv.FormatFloat(record.Location.Latitude, 'f', 6, 64)
	var long string = strconv.FormatFloat(record.Location.Longitude, 'f', 6, 64)

	return record.Country.Names["en"], record.Country.IsoCode, lat, long
}

func getdata(ip string) IPData.DNSRecord {
	return IPData.NewDnsRecord(nil, []string{ip})
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
		eventID, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "eventid")
		if string(eventID) != "cowrie.session.connect" {
			continue
		}
		var redirect bool = false
		var sip *db.SourceIPDescription = db.NewSourceIPDescription()
		var ipdata *db.IPDataDescription = db.NewIPDataDescription()

		srcIP, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "src_ip")
		proto, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "protocol")
		timestamp, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "timestamp")
		sip.CountryName, sip.CountryCode, sip.Latitude, sip.Longitude = getGEOIPData(string(srcIP))

		{
			sip.SRCIP = string(srcIP)
			sip.Time = string(timestamp)
			sip.Protocol = string(proto)
			sip.EventID = string(eventID)
		}

		for i := 0; i < len(sipArray); i++ {
			// incrementing hit counts for each IP
			if sipArray[i].SRCIP == string(srcIP) {
				sipArray[i].HitCount++
				// The timestamp is the last time the IP was seen
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

	return sipArray
}
