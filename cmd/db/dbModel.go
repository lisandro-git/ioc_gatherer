package db

type DomainDescription struct {
	// Lisandro : make subDomain array ?
	ID        int    `json:"DomainID" gorm:"primaryKey;autoIncrement:true"`
	IP        uint32 `json:"IP"`
	Domain    string `json:"Domain"`
	Subdomain string `json:"Subdomain"`

	FK_SourceIP int                 `json:"-"`
	FK_IPData   int                 `json:"-"`
	SourceIP    SourceIPDescription `json:"SourceIP,omitempty"  gorm:"foreignKey:FK_SourceIP;references:SRCIP"`
	IPData      IPDataDescription   `json:"IPData,omitempty"    gorm:"foreignKey:FK_IPData;references:ID"`
}

type SourceIPDescription struct {
	ID          int    `json:"SourceIPID" gorm:"primaryKey;autoIncrement:true"`
	SRCIP       uint32 `json:"IP" gorm:"UNIQUE"`
	CountryName string `json:"country,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
	City        string `json:"city,omitempty"`
	Latitude    string `json:"latitude,omitempty"`
	Longitude   string `json:"longitude,omitempty"`
	HitCount    int    `json:"hitCount,omitempty"`
	Time        string `json:"time,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	EventID     string `json:"eventID,omitempty"`
}

type IPDataDescription struct {
	ID     int    `json:"IPDataID" gorm:"primaryKey;autoIncrement:true"`
	IP     uint32 `json:"IP,omitempty"`
	Domain string `json:"Domain"`
	CNAME  string `json:"CNAME,omitempty"`
	MX     string `json:"MX,omitempty"`
	NS     string `json:"NS,omitempty"`
	TXT    string `json:"TXT,omitempty"`
}

func NewDomainDescription() *DomainDescription {
	return &DomainDescription{}
}

func NewSourceIPDescription() *SourceIPDescription {
	return &SourceIPDescription{}
}

func (sid *SourceIPDescription) TableName() string {
	return SrcIPTableName
}

func NewIPDataDescription() *IPDataDescription {
	return &IPDataDescription{}
}

func (ipd *IPDataDescription) TableName() string {
	return IPDataTableName
}

/*
func (ipd *IPDataDescription) new_PrepareForDB(dnsR IPData.DNSRecord) {

}

func (ipd *IPDataDescription) old_PrepareForDB(dnsR IPData.DNSRecord) {

	for _, v := range dnsR.Domains {
		ipd.Domains += v + ";"
	}
	for _, v := range dnsR.DomainNames {
		ipd.DomainNames += v + ";"
	}
	for _, v := range dnsR.IP {
		//x, _ := IPv4ToInt(net.ParseIP(string(v)))
		ipd.IP = v
		break // Lisandro : remove the break to add multiple IPs
	}
	for _, v := range dnsR.CNAME {
		ipd.CNAME += v + ";"
	}
	for _, v := range dnsR.MX {
		ipd.MX += v + ";"
	}
	for _, v := range dnsR.NS {
		ipd.NS += v + ";"
	}
	for _, v := range dnsR.TXT {
		ipd.TXT += v + ";"
	}
}
*/
