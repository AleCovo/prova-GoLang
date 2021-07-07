package main

import (
	"database/sql"
	_ "database/sql/driver"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	postgres_address = flag.String("postgres_address", "192.168.2.230", "")
	postgres_port    = flag.Int("postgres_port", 5432, "")
	postgres_dbname  = flag.String("postgres_dbname", "postgres", "")
	postgres_user    = flag.String("postgres_user", "udev01", "")
	postgres_pass    = flag.String("postgres_pass", "QzCtt35QRxZX", "")
)

var gormDB *gorm.DB

func GetDBConnection() *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		*postgres_address,
		*postgres_port,
		*postgres_user,
		*postgres_pass,
		*postgres_dbname,
	)

	dbConnection, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err.Error())
	}

	return dbConnection
}

func GetGormDBConnection() (*gorm.DB, error) {

	sqlDB, _ := gormDB.DB()
	err := sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func InitializeDBConnection() error {
	var err error

	gormDB, err = gorm.Open(postgres.New(
		postgres.Config{
			Conn: GetDBConnection(),
		}),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	)

	return err
}
