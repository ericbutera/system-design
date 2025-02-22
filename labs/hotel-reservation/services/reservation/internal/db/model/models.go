package model

import "time"

type Reservation struct {
	ID         int       `gorm:"primaryKey;autoIncrement"`
	Quantity   int       `gorm:"column:quantity;not null" validate:"required,gte=1"`
	CheckIn    time.Time `gorm:"column:checkin;not null" validate:"required,gte=today"`
	CheckOut   time.Time `gorm:"column:checkout;not null" validate:"required,gtfield=CheckIn"`
	Status     string    `gorm:"column:status;not null;default:'PENDING'"`                 // PENDING, CONFIRMED, CANCELLED
	RoomTypeID int       `gorm:"column:room_type_id;not null;index"`                       // Foreign key for room_types
	GuestID    int       `gorm:"column:guest_id;not null;index" validate:"required,gte=1"` // Foreign key for guests
	HotelID    int       `gorm:"column:hotel_id;not null;index" validate:"required,gte=1"` // Foreign key for hotels
	PaymentID  *int      `gorm:"column:payment_id;index"`                                  // Foreign key for payments (nullable)
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`

	// Foreign key associations (optional)
	RoomType RoomType `gorm:"foreignKey:RoomTypeID;constraint:OnDelete:RESTRICT"`
	Guest    Guest    `gorm:"foreignKey:GuestID;constraint:OnDelete:RESTRICT"`
	// Payment  Payment  `gorm:"foreignKey:PaymentID;constraint:OnDelete:SET NULL"`
}

type Guest struct {
	ID    int    `gorm:"primaryKey"`
	Name  string `gorm:"type:varchar(255);not null"`
	Email string `gorm:"type:varchar(255);not null"`
}

type RoomTypeInventory struct {
	HotelID        int       `gorm:"primaryKey;not null"` // Foreign key for hotels
	RoomTypeID     int       `gorm:"primaryKey;not null"` // Foreign key for room_types
	Date           time.Time `gorm:"primaryKey;not null"` // Date field
	TotalInventory int       `gorm:"not null"`
	TotalReserved  int       `gorm:"not null"`

	// Foreign key associations (optional)
	// Hotel    Hotel    `gorm:"foreignKey:HotelID;constraint:OnDelete:CASCADE"`
	// RoomType RoomType `gorm:"foreignKey:RoomTypeID;constraint:OnDelete:CASCADE"`
}

type RoomType struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(255);not null"`
}

// type Payment struct {
// 	ID            int       `gorm:"primaryKey"`
// 	CorrelationID int       `gorm:"not null"`
// 	Amount        float64   `gorm:"type:decimal(10,2);not null"`
// 	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
// }
// type Hotel struct {
// 	ID   int `gorm:"primaryKey"`
// 	Name string
// }
