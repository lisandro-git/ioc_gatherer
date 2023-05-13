package cowrieParser

import (
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

var total_IPs []string

func getGEOIPData(src_ip string, sipd *db.SourceIPDescription) {
	// Lisandro: return city
	ipdb, err := geoip2.Open("/root/Downloads/GeoLite2-City_20230421/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer ipdb.Close()

	ip := net.ParseIP(src_ip)
	record, err := ipdb.City(ip)
	if err != nil {
		log.Fatal(err)
	}

	sipd.City = record.City.Names["en"]
	sipd.CountryCode = record.Country.IsoCode
	sipd.CountryName = record.Country.Names["en"]
	sipd.Latitude = strconv.FormatFloat(record.Location.Latitude, 'f', 6, 64)
	sipd.Longitude = strconv.FormatFloat(record.Location.Longitude, 'f', 6, 64)
}

func getdata(description *db.DomainDescription) []db.DomainDescription {
	return IPData.NewDnsRecord(description)
}

func appendUnique(ip string) {
	for _, v := range total_IPs {
		if ip == v {
			return
		}
	}
	total_IPs = append(total_IPs, ip)
}

func saveWhitelist(fpDB *db.DB, f db.FileWhiteList) {
	fpDB.Database.Save(&f)
}

func saveBlacklist(fpDB *db.DB, f *db.FileBlacklist) {
	fpDB.Database.Save(&f)
}

var filepathDB *db.DB = db.InitDB(db.FILEDBNAME)

func ReadFile(filePath string) ([]db.DomainDescription, error) {
	var fileScanner *bufio.Scanner = cmd.ImportDataFile(filePath)

	var domArray []db.DomainDescription

	for fileScanner.Scan() {
		eventID, _, _, err := jsonparser.Get([]byte(fileScanner.Text()), "eventid")
		if err != nil {
			fmt.Println("Error while parsing JSON: ", err)
			var fileBlacklist *db.FileBlacklist = db.NewFileBlacklist()
			fileBlacklist.FilePath = filePath
			fileBlacklist.THash = cmd.HashChunk(cmd.GetChunk(filePath))

			saveBlacklist(filepathDB, fileBlacklist)
			return nil, err
		}
		if string(eventID) != "cowrie.session.connect" {
			continue
		}

		{
			var fileWhiteList *db.FileWhiteList = db.NewFileWhiteList()
			fileWhiteList.FilePath = filePath
			fileWhiteList.THash = string(cmd.HashChunk(cmd.GetChunk(filePath)))
			saveWhitelist(filepathDB, *fileWhiteList)
		}

		var dom *db.DomainDescription = db.NewDomainDescription()

		srcIP, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "src_ip")
		proto, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "protocol")
		timestamp, _, _, _ := jsonparser.Get([]byte(fileScanner.Text()), "timestamp")
		getGEOIPData(string(srcIP), &dom.SourceIP)

		intIP, _ := cmd.IPv4ToInt(net.ParseIP(string(srcIP)))
		dom.IP = intIP
		dom.SourceIP.SRCIP = intIP
		dom.SourceIP.Time = string(timestamp)
		dom.SourceIP.Protocol = string(proto)
		dom.SourceIP.EventID = string(eventID)
		appendUnique(string(srcIP))

		// incrementing hit counts for each IP
		for i := 0; i < len(domArray); i++ {
			if domArray[i].IP == intIP {
				domArray[i].SourceIP.HitCount++
				break
			}
		}

		// If no other domain is found, the directly goes to dom, else it loops through the domains,
		// adding them to the final array
		var domainData []db.DomainDescription = getdata(dom)
		if domainData == nil {
			domArray = append(domArray, *dom)
		} else {
			for i := 0; i < len(domainData); i++ {
				domArray = append(domArray, domainData[i])
			}
		}
	}
	fmt.Println("Total IPs: ", len(total_IPs))
	return domArray, nil
}
