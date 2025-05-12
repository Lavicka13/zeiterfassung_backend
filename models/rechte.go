package models

// Rechte definiert die Benutzerrechte in der Anwendung
type Rechte struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:50;not null"`
	Beschreibung string `gorm:"size:255"`
}

func (Rechte) TableName() string {
	return "rechte"
}