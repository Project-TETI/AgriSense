package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email         string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	PasswordHash  string         `gorm:"type:text;not null" json:"-"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID" json:"-"`
	Inspections   []Inspection   `gorm:"foreignKey:UserID" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string     `gorm:"type:text;not null" json:"token_hash"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
	RevokedAt *time.Time `gorm:"type:timestamp" json:"revoked_at,omitempty"`
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
}

func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}
	return
}

type Plant struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	DisplayName string       `gorm:"type:varchar(100);not null" json:"display_name"`
	Description *string      `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time    `gorm:"not null" json:"created_at"`
	Diseases    []Disease    `gorm:"foreignKey:PlantID" json:"diseases,omitempty"`
	Inspections []Inspection `gorm:"foreignKey:PlantID" json:"-"`
}

func (p *Plant) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

type Disease struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PlantID        uuid.UUID `gorm:"type:uuid;not null;index" json:"plant_id"`
	ClassName      string    `gorm:"type:varchar(150);uniqueIndex;not null" json:"class_name"`
	DisplayName    string    `gorm:"type:varchar(150);not null" json:"display_name"`
	Cause          string    `gorm:"type:text;not null" json:"cause"`
	Description    string    `gorm:"type:text;not null" json:"description"`
	Recommendation string    `gorm:"type:text;not null" json:"recommendation"`
	Plant          Plant     `gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE;" json:"plant,omitempty"`
}

func (d *Disease) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return
}

type Inspection struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	PlantID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"plant_id"`
	DiseaseID       *uuid.UUID `gorm:"type:uuid;index" json:"disease_id,omitempty"`
	ImageURL        string     `gorm:"type:text;not null" json:"image_url"`
	ConfidenceScore *float64   `gorm:"type:float" json:"confidence_score,omitempty"`
	Status          string     `gorm:"type:varchar(50);not null;index" json:"status"`
	Latitude        *float64   `gorm:"type:double precision" json:"latitude,omitempty"`
	Longitude       *float64   `gorm:"type:double precision" json:"longitude,omitempty"`
	CreatedAt       time.Time  `gorm:"not null;index" json:"created_at"`
	User            User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`
	Plant           Plant      `gorm:"foreignKey:PlantID;" json:"plant,omitempty"`
	Disease         *Disease   `gorm:"foreignKey:DiseaseID;" json:"disease,omitempty"`
}

func (i *Inspection) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return
}
