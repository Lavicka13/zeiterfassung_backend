package models

import "time"

// ArbeitszeitsAudit speichert Änderungen an Arbeitszeiten
type ArbeitszeitsAudit struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ArbeitszeitID uint      `json:"arbeitszeit_id"`
	NutzerID     uint      `json:"nutzer_id"`
	Feld         string    `json:"feld"`
	AlterWert    string    `json:"alter_wert"`
	NeuerWert    string    `json:"neuer_wert"`
	BearbeitetAm time.Time `json:"bearbeitet_am"`
}

func (ArbeitszeitsAudit) TableName() string {
	return "arbeitszeiten_audit"
}