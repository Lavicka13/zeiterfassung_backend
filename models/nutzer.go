package models

import "time"

type Nutzer struct {
    ID        uint   `gorm:"primaryKey"`
    Vorname   string
    Nachname  string
    Email     string `gorm:"unique"`
    PwHash    string
    RechteID  uint
    CreatedAt time.Time
}

func (Nutzer) TableName() string {
    return "nutzer" // <-- das verhindert das automatische "s"
}
