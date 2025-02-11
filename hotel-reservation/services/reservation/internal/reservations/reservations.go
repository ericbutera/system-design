package reservations

import (
	"context"
	"time"

	"strconv"

	graphModel "github.com/ericbutera/system-design/hotel-reservation/services/reservation/graph/model"
	dbModel "github.com/ericbutera/system-design/hotel-reservation/services/reservation/internal/db/model"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Reservations struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Reservations {
	return &Reservations{db: db}
}

const UpdateInventorySql = `
UPDATE room_type_inventory
SET total_reserved = total_reserved + ?
WHERE
	room_type_id = ? AND
	hotel_id = ? AND
	date BETWEEN ? AND ?
`

func (r *Reservations) Create(ctx context.Context, reservation *dbModel.Reservation) (*dbModel.Reservation, error) {
	// TODO: validate input
	if err := r.Validate(ctx, reservation); err != nil {
		return nil, err
	}

	// not even querying available rooms before creating the reservation will work during high contention
	// the only way is to actually try to create the reservation and handle the error
	err := r.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(
			UpdateInventorySql,
			reservation.Quantity,
			reservation.RoomTypeID,
			reservation.HotelID,
			reservation.CheckIn,
			reservation.CheckOut,
		)
		if res.Error != nil {
			return res.Error
		}
		if err := tx.Create(&reservation).Error; err != nil {
			return err
		}
		return nil
	})
	return reservation, err
}

func (r *Reservations) Validate(ctx context.Context, reservation *dbModel.Reservation) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(reservation)
}

func TimeFromString(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02")
}

func (r *Reservations) GetByID(ctx context.Context, id string, guestID int) (*dbModel.Reservation, error) {
	var data dbModel.Reservation
	res := r.db.Where("id = ? AND guest_id = ?", id, guestID).First(&data)
	if res.Error != nil {
		return nil, res.Error
	}
	return &data, nil
}

func (r *Reservations) GetByUser(ctx context.Context, guestID int) ([]*dbModel.Reservation, error) {
	var data []*dbModel.Reservation
	res := r.db.Where("guest_id = ?", guestID).Find(&data)
	if res.Error != nil {
		return nil, res.Error
	}
	return data, nil
}

func ReservationToDb(r *graphModel.Reservation) *dbModel.Reservation {
	id, _ := strconv.Atoi(r.ID)
	return &dbModel.Reservation{
		ID:         id,
		RoomTypeID: r.RoomTypeID,
		Quantity:   r.Quantity,
		CheckIn:    TimeFromString(r.CheckIn),
		CheckOut:   TimeFromString(r.CheckOut),
		Status:     r.Status,
		GuestID:    r.GuestID,
		HotelID:    r.HotelID,
		PaymentID:  nil,
	}
}

func ReservationToGraph(r *dbModel.Reservation) *graphModel.Reservation {
	return &graphModel.Reservation{
		ID:         strconv.Itoa(r.ID),
		RoomTypeID: r.RoomTypeID,
		Quantity:   r.Quantity,
		CheckIn:    TimeToString(r.CheckIn),
		CheckOut:   TimeToString(r.CheckOut),
		Status:     r.Status,
		GuestID:    r.GuestID,
		HotelID:    r.HotelID,
	}
}
