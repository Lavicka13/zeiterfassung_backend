package models

import "time"



type Arbeitszeit struct {
    ID          uint       `gorm:"primaryKey" json:"id"`
    NutzerID    uint       `json:"nutzer_id"`
    Anfangszeit time.Time  `json:"anfangszeit"`
    Endzeit     *time.Time `json:"endzeit,omitempty"`
    Pause       int        `json:"pause"`
    Datum       time.Time  `json:"datum"`
    ErstelltAm  time.Time  `json:"erstellt_am"`        // KEIN Pointer!
    Bearbeitet  *time.Time `json:"bearbeitet,omitempty"` // Pointer, weil optional
}


func (Arbeitszeit) TableName() string {
    return "arbeitszeiten"
}