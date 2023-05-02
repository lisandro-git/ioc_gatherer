package IPData

import (
	"net"
	"strings"
)

type DNSRecord struct {
	Domains     []string `json:"Domain"`
	DomainNames []string `json:"DomainNames"`
	IP          []string `json:"IP"`
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

// NewDnsRecord returns a pointer to a new DNSRecord struct
func NewDnsRecord(domain, ip []string) DNSRecord {
	dnsR := DNSRecord{
		Domains: domain,
		IP:      ip,
	}

	if domain != nil {
		dnsR.GetIP()
		dnsR.GetDomainNames()
	}
	if ip != nil {
		dnsR.GetDomainNames()
		dnsR.GetIP()
	}

	dnsR.GetCNAME()
	dnsR.GetMX()
	dnsR.GetNS()
	dnsR.GetTXT()
	dnsR.GetIP()

	return dnsR
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
		addrs, _ := net.LookupAddr(ips)
		for _, v := range addrs {
			dnsR.Domains = appendUnique(dnsR.Domains, v)
			domainSplited := strings.Split(v, ".")
			dn := domainSplited[len(domainSplited)-3] + "." + domainSplited[len(domainSplited)-2]
			dnsR.Domains = appendUnique(dnsR.Domains, dn)

			//dnsR.GetDomainName()
		}
	}
}

// GetIP uses net.LookupIP to get the IP address of the domain
func (dnsR *DNSRecord) GetIP() {
	for _, domain := range dnsR.Domains {
		ips, _ := net.LookupIP(domain)
		for _, v := range ips {
			dnsR.IP = appendUnique(dnsR.IP, v.String())
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
