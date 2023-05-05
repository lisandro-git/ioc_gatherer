package db

import (
	"main/cmd/IPData"
	"gorm.io/gorm"
)

type SourceIPDescription struct {
	gorm.Model
	ID          int               `json:"locationID"`
	SRCIP       string            `json:"SRCIP" gorm:"UNIQUE"`
	CountryName string            `json:"country"`
	CountryCode string            `json:"countryCode"`
	City        string            `json:"city"`
	Latitude    string            `json:"latitude,omitempty"`
	Longitude   string            `json:"longitude,omitempty"`
	HitCount    int               `json:"hitCount"`
	Time        string            `json:"time"`
	Protocol    string            `json:"protocol"`
	EventID     string            `json:"eventID"`
	IPDataID    int               `json:"-"`
	IPData      IPDataDescription `json:"IPData,omitempty" gorm:"foreignKey:IPDataID;references:ID"`
}

func NewSourceIPDescription() *SourceIPDescription {
	return &SourceIPDescription{}
}

func (sid *SourceIPDescription) TableName() string {
	return SrcIPTableName
}

type IPDataDescription struct {
	ID          int64  `json:"IPDataID" gorm:"primaryKey;autoIncrement:true"`
	Domains     string `json:"Domain"`
	DomainNames string `json:"DomainNames"`
	IP          string `json:"IP"`
	CNAME       string `json:"CNAME"`
	MX          string `json:"MX"`
	NS          string `json:"NS"`
	TXT         string `json:"TXT"`
}

func NewIPDataDescription() *IPDataDescription {
	return &IPDataDescription{}
}

func (ipd *IPDataDescription) TableName() string {
	return IPDataTableName
}

func (ipd *IPDataDescription) PrepareForDB(dnsR IPData.DNSRecord) {
	for _, v := range dnsR.Domains {
		ipd.Domains += v + ";"
	}
	for _, v := range dnsR.DomainNames {
		ipd.DomainNames += v + ";"
	}
	for _, v := range dnsR.IP {
		ipd.IP += v + ";"
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
