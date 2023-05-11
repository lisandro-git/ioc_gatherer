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
	"main/cmd"
	"fmt"
	"main/cmd/db"
)

func getGEOIPData(src_ip string) (string, string, string, string) {
	// Lisandro: return city
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

func getdata(description *db.DomainDescription) []db.DomainDescription {
	fmt.Println("IP: ", cmd.IntToIPv4(description.IP).String())
	x := IPData.NewDnsRecord(description)

	return x
}

func ReadFile(filePath string) []db.DomainDescription {
	data, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	var fileScanner *bufio.Scanner = bufio.NewScanner(data)
	fileScanner.Split(bufio.ScanLines)

	var domArray []db.DomainDescription
	for fileScanner.Scan() {
		var redirect bool
		eventID, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "eventid")
		if string(eventID) != "cowrie.session.connect" {
			continue
		}

		var dom *db.DomainDescription = db.NewDomainDescription()
		var sip *db.SourceIPDescription = db.NewSourceIPDescription()

		srcIP, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "src_ip")
		proto, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "protocol")
		timestamp, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "timestamp")
		sip.CountryName, sip.CountryCode, sip.Latitude, sip.Longitude = getGEOIPData(string(srcIP))
		{
			x, _ := cmd.IPv4ToInt(net.ParseIP(string(srcIP)))
			sip.SRCIP = x
			sip.Time = string(timestamp)
			sip.Protocol = string(proto)
			sip.EventID = string(eventID)
		}

		for i := 0; i < len(domArray); i++ {
			// incrementing hit counts for each IP
			x, _ := cmd.IPv4ToInt(net.ParseIP(string(srcIP)))
			if domArray[i].IP == x {
				domArray[i].SourceIP.HitCount++
				// The timestamp is the last time the IP was seen
				//fmt.Printf("IP %s already exists\n", string(srcIP))
				redirect = true
				break
			}
		}
		if redirect {
			continue
		}
		//fmt.Printf("New IP: %s found\n", string(srcIP))
		sip.HitCount++
		dom.IP = sip.SRCIP

		dom.SourceIP = *sip

		a := getdata(dom)
		for i := 0; i < len(a); i++ {
			domArray = append(domArray, a[i])
		}
	}
	return domArray
}
