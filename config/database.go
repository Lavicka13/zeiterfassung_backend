package config

import (
    "log"
    

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
    dsn := "root:Passw0rd@tcp(127.0.0.1:3306)/zeiterfassung?parseTime=true"
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("DB Fehler:", err)
    }
}
