package reservations

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"strconv"

	graphModel "github.com/ericbutera/system-design/hotel-reservation/services/reservation/graph/model"
	dbModel "github.com/ericbutera/system-design/hotel-reservation/services/reservation/internal/db/model"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var (
	ErrNotEnoughInventory = errors.New("Not enough inventory")
)

type Reservations struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Reservations {
	return &Reservations{db: db}
}

func (r *Reservations) Create(ctx context.Context, reservation *dbModel.Reservation) (*dbModel.Reservation, error) {
	// TODO: validate input
	if err := r.Validate(ctx, reservation); err != nil {
		return nil, err
	}

	// not even querying available rooms before creating the reservation will work during high contention
	// the only way is to actually try to create the reservation and handle the error
	// note: this is a range update that will increment the total_reserved column for all inventory records in the date range
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.updateInventory(tx, reservation); err != nil {
			return err
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

const UpdateInventorySql = `
UPDATE room_type_inventory
SET total_reserved = total_reserved + ?
WHERE
	hotel_id = ? AND
	room_type_id = ? AND
	date BETWEEN ? AND ?
`

func (r *Reservations) updateInventory(tx *gorm.DB, reservation *dbModel.Reservation) error {
	days := reservation.CheckOut.Sub(reservation.CheckIn).Hours() / 24
	end := reservation.CheckOut.AddDate(0, 0, -1) // reservation inventory is inclusive of check-in date but exclusive of check-out date
	slog.Info("creating reservation", "days", days, "quantity", reservation.Quantity)

	res := tx.Exec(
		UpdateInventorySql,
		reservation.Quantity,
		reservation.RoomTypeID,
		reservation.HotelID,
		TimeToString(reservation.CheckIn),
		TimeToString(end),
	)
	if res.Error != nil {
		if strings.Contains(res.Error.Error(), "check_room_count") {
			slog.Info("err check_room_constraint", "rows", res.RowsAffected, "days", days)
			return ErrNotEnoughInventory
		}
		return res.Error
	}
	if res.RowsAffected != int64(days) {
		slog.Info("err partial reservation", "rows", res.RowsAffected, "days", days)
		tx.Rollback()
		return ErrNotEnoughInventory
	}
	return nil
}

func (r *Reservations) CreateInventory() error {
	sql := `
	WITH RECURSIVE DateSeries AS (
		SELECT CURRENT_DATE::timestamp AS date
		UNION ALL
		SELECT date + INTERVAL '1 day' FROM DateSeries WHERE date < CURRENT_DATE + INTERVAL '10 days'
	)
	INSERT INTO room_type_inventory (hotel_id, room_type_id, date, total_inventory, total_reserved)
	SELECT r.hotel_id,
		r.room_type_id,
		ds.date::date,
		COUNT(*) AS total_inventory,
		0 AS total_reserved
	FROM rooms r
	JOIN DateSeries ds
	ON ds.date >= CURRENT_DATE
	GROUP BY r.hotel_id, r.room_type_id, ds.date
	ON CONFLICT (hotel_id, room_type_id, date) DO NOTHING;
	`
	return r.db.Exec(sql).Error
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
