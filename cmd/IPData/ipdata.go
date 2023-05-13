package IPData

import (
	"net"
	"strings"
	"main/cmd"
	"fmt"
	"github.com/bobesa/go-domain-util/domainutil"
	"main/cmd/db"
)

type DNSRecord struct {
	IP          []uint32 `json:"IP"`
	Domains     []string `json:"Domain"`
	DomainNames []string `json:"DomainNames"`
	CNAME       []string `json:"CNAME"`
	MX          []string `json:"MX"`
	NS          []string `json:"NS"`
	TXT         []string `json:"TXT"`
}

func appendUnique(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

func NewEmptyDNSRecord() *DNSRecord {
	return &DNSRecord{}
}

func contains(s []uint32, str uint32) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func addDNSRecordToStruct(a *db.DomainDescription, domain string) {
	a.IPData.CNAME = GetCNAME_str(domain)
	b := GetMX_str(domain)
	if len(b) > 0 {
		a.IPData.MX = b[0]
	}
	b = GetNS_str(domain)
	if len(b) > 0 {
		a.IPData.NS = b[0]
	}
	b = GetTXT_str(domain)
	if len(b) > 0 {
		a.IPData.TXT = b[0]
	}
}

var totalIP_Domain = []uint32{}

func removeTrailingDot(domain string) string {
	if domain[len(domain)-1] == 46 { // 46  = .
		domain = domain[:len(domain)-1]
	}
	return domain
}

func subdomainFirst(description *db.DomainDescription) {

}

func domainFirst(description *db.DomainDescription) {

}

// NewDnsRecord returns a pointer to a new DNSRecord struct
func NewDnsRecord(r *db.DomainDescription) []db.DomainDescription {
	var z []db.DomainDescription
	totalIP_Domain = append(totalIP_Domain, r.IP)
	dnsR := DNSRecord{
		Domains: nil,
		IP:      []uint32{r.IP},
	}
	dnsR.GetDomainNames()

	if len(dnsR.Domains) == 0 {
		GetDomainNames_str(r.IP)
		r.IPData.IP = r.IP
		return z
	}

	var domain string

	for _, v := range dnsR.Domains {
		x := db.DomainDescription{
			SourceIP: r.SourceIP,
		}
		domain = domainutil.Domain(v)
		if len(v) == 0 {
			continue
		}
		x.IP = r.IP
		x.IPData.IP = r.IP
		x.Domain = domain
		x.IPData.Domain = domain
		if domainutil.HasSubdomain(v) {
			x.Subdomain = v
			addDNSRecordToStruct(&x, v)
			z = append(z, x)
		}
		x = db.DomainDescription{
			SourceIP: r.SourceIP,
		}
		x.IP = r.IP
		x.IPData.IP = r.IP
		x.Domain = domain
		x.IPData.Domain = domain
		addDNSRecordToStruct(&x, domain)

		z = append(z, x)
	}

	return z
}

func GetCNAME_str(domain string) string {
	x, _ := net.LookupCNAME(domain)
	return x
}

func GetDomainNames_str(x uint32) []string {
	ip := cmd.IntToIPv4(x)
	a := ip.String()
	addrs, _ := net.LookupAddr(a)
	return addrs

}

// GetMX uses net.LookupMX to get the MX records of the domain
func GetMX_str(domain string) []string {
	b := []string{}
	mx, _ := net.LookupMX(domain)
	for _, v := range mx {
		b = appendUnique(b, v.Host)
	}
	return b
}

// GetNS uses net.LookupNS to get the NS records of the domain
func GetNS_str(domain string) []string {
	ns, _ := net.LookupNS(domain)
	b := []string{}
	for _, v := range ns {
		b = appendUnique(b, v.Host)
	}
	return b
}

// GetTXT uses net.LookupTXT to get the TXT records of the domain
func GetTXT_str(domain string) []string {
	txt, _ := net.LookupTXT(domain)
	return txt
}

// GetIP uses net.LookupIP to get the IP address of the domain
func GetIP_str(domain string) []uint32 {
	ips, _ := net.LookupIP(domain)
	a := []uint32{}
	if len(ips) == 1 {
		x, _ := cmd.IPv4ToInt(ips[0])
		return []uint32{x}
	} else if len(ips) > 1 {

		for _, v := range ips { // Lisandro : modify this part of code to make it more efficient
			x, err := cmd.IPv4ToInt(v)
			if err != nil {
				fmt.Println(err, v)
				continue
			}
			a = append(a, x)
		}
	}
	return a
}

// GetCNAME uses net.LookupCNAME to get the CNAME record of the domain
func (dnsR *DNSRecord) GetCNAME() {
	for _, v := range dnsR.Domains {
		x, _ := net.LookupCNAME(v)
		dnsR.CNAME = appendUnique(dnsR.CNAME, x)
	}
}

// GetMX uses net.LookupMX to get the MX records of the domain
func (dnsR *DNSRecord) GetMX() {
	for _, domain := range dnsR.Domains {
		mx, _ := net.LookupMX(domain)
		for _, v := range mx {
			dnsR.MX = appendUnique(dnsR.MX, v.Host)
		}
	}
}

// GetNS uses net.LookupNS to get the NS records of the domain
func (dnsR *DNSRecord) GetNS() {
	for _, domain := range dnsR.Domains {
		ns, _ := net.LookupNS(domain)
		for _, v := range ns {
			dnsR.NS = appendUnique(dnsR.NS, v.Host)
		}
	}
}

// GetTXT uses net.LookupTXT to get the TXT records of the domain
func (dnsR *DNSRecord) GetTXT() {
	for _, domain := range dnsR.Domains {
		txt, _ := net.LookupTXT(domain)
		for _, v := range txt {
			dnsR.TXT = appendUnique(dnsR.TXT, v)
		}
	}
}

// GetDomainNames uses net.LookupAddr to get the domain names of the IP addresses contained inside the DNSRecord.IP
func (dnsR *DNSRecord) GetDomainNames() {
	for _, ips := range dnsR.IP {
		ip := cmd.IntToIPv4(ips)
		a := ip.String()
		addrs, _ := net.LookupAddr(a)
		for _, v := range addrs {
			if len(v) == 0 {
				continue
			}
			if len(domainutil.Domain(v)) == 0 {
				dnsR.Domains = appendUnique(dnsR.Domains, removeTrailingDot(v))
			} else {
				dnsR.Domains = appendUnique(dnsR.Domains, domainutil.Domain(removeTrailingDot(v)))
			}
		}
	}
}

// GetIP uses net.LookupIP to get the IP address of the domain
func (dnsR *DNSRecord) GetIP() {
	for _, domain := range dnsR.Domains {
		ips, _ := net.LookupIP(domain)
		if len(ips) == 1 {
			x, _ := cmd.IPv4ToInt(ips[0])
			dnsR.IP = append(dnsR.IP, x)
			return
		} else if len(ips) > 1 {
			for _, v := range ips { // Lisandro : modify this part of code to make it more efficient
				x, err := cmd.IPv4ToInt(v)
				if err != nil {
					fmt.Println(err, v)
					continue
				}
				dnsR.IP = append(dnsR.IP, x)
			}
		}
	}
}

func (dnsR *DNSRecord) GetDomainName() {
	for _, domain := range dnsR.Domains {
		domainSplited := strings.Split(domain, ".")
		domainName := domainSplited[len(domainSplited)-3] + "." + domainSplited[len(domainSplited)-2]
		//dnsR.DomainNames = appendUnique(dnsR.DomainNames, domainName)
		dnsR.Domains = appendUnique(dnsR.DomainNames, domainName)
	}
}
