package db

import (
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/schema"
	"main/cmd"
)

type DB struct {
	Database *gorm.DB
}

func InitDB() *DB {
	db := &DB{}
	db.ConnectToDB()
	db.BuildDatabase()
	return db
}

// ConnectToDB creates a database if it doesn't exist
func (d *DB) ConnectToDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	d.Database = db
}

func (d *DB) BuildDatabase() {
	var err error = d.Database.AutoMigrate(
		&DomainDescription{},
		&SourceIPDescription{},
		&IPDataDescription{},
		//&SubDomainDescription{},
	)
	if err != nil {
		return
	}
}

func getRowID(x *gorm.DB, ip uint32) int {
	y := &SourceIPDescription{}
	b := cmd.IntToIPv4(ip) // Lisandro : Clean this
	a, _ := x.Model(&SourceIPDescription{}).Select("id").Where("src_ip = ?", b).Rows()
	for a.Next() {
		err := x.ScanRows(a, y)
		if err != nil {
			fmt.Println(err)
			return 0
		}
	}
	fmt.Printf("Updated data to database: ID: %d, IP: %s\n", y.ID, ip)
	return y.ID
}

func (d *DB) UpdateDBRow(sourceIPRow SourceIPDescription) {
	rowID := getRowID(d.Database, sourceIPRow.SRCIP)

	d.Database.Model(&SourceIPDescription{}).Where("id = ?", rowID).Updates(
		SourceIPDescription{
			CountryName: sourceIPRow.CountryName,
			CountryCode: sourceIPRow.CountryCode,
			Latitude:    sourceIPRow.Latitude,
			Longitude:   sourceIPRow.Longitude,
			City:        sourceIPRow.City,
			HitCount:    sourceIPRow.HitCount,
			Time:        sourceIPRow.Time,
			Protocol:    sourceIPRow.Protocol,
			EventID:     sourceIPRow.EventID,
		})
}
