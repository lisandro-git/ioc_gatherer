package db

const (
	HOST     = "localhost"
	PORT     = 5432
	USER     = "postgres"
	PASSWORD = "123"
	IPDBNAME = "iocdb"

	SrcIPTableName  = "SourceIP"
	IPDataTableName = "IPData"
	DomainTableName = "DomainDescription"

	BlacklistTableName = "BlacklistedFile"
	WhitelistTableName = "WhitelistedFile"
	FILEDBNAME         = "filelistdb"
)
