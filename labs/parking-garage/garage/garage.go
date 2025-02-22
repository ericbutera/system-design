package garage

import (
	"errors"
	"log/slog"
	"math"
	"time"

	"github.com/samber/lo"
)

var (
	ErrNoAvailability = errors.New("no availability")
)

const (
	RATE_AMOUNT   float32 = 1.00
	RATE_DURATION         = 15 * time.Minute
)

type Receipt struct {
	ID       string
	Total    float32 // Ignoring currency & localization
	Checkin  time.Time
	Checkout time.Time
	Garage   *Garage
}

func Charge(checkin time.Time, checkout time.Time) float32 {
	// TODO: support vehicle size and custom rates
	minutes := checkout.Sub(checkin).Minutes()
	intervals := math.Ceil(minutes / 15)
	total := float32(intervals) * RATE_AMOUNT
	return total
}

func NewReceipt(ticket *Ticket) *Receipt {
	checkout := lo.FromPtr(ticket.Checkout)
	return &Receipt{
		Total:    Charge(ticket.Checkin, checkout),
		Checkin:  ticket.Checkin,
		Checkout: checkout,
	}
}

type Garage struct {
	Name   string
	Levels []*Level
}

func New(name string, levels []*Level) *Garage {
	// loc, err := time.LoadLocation("America/Detroit") // TODO: DepInject | Options
	return &Garage{
		Name:   name,
		Levels: levels,
	}
}

type Reservation struct {
	Level *Level
	Spots []*Spot
}

// determine if the provided vehicle can park
func (p *Garage) Reserve(vehicle Vehicle) (*Reservation, error) {
	// TODO: possibly move into the "Level" struct
	requiredSpots := vehicle.Size()
	for _, level := range p.Levels {
		slog.Debug("acquire spots", "level", level)
		// vehicle size requires size contiguous spots
		contiguousSpots := 0
		for spotOffset, spot := range level.Spots {
			slog.Debug("spots", "spotOffset", spotOffset, "contiguousSpots", contiguousSpots)
			if spot.Ticket == nil { // available spot
				contiguousSpots++
			}
			slog.Debug("checking required == contiguous", "required", requiredSpots, "contiguous", contiguousSpots)
			if requiredSpots == contiguousSpots { // found required spots
				start := spotOffset - (requiredSpots - 1)
				end := requiredSpots
				slog.Debug("match!", "start", start, "end", end)
				return &Reservation{
					Level: level,
					Spots: level.Spots[start:end],
				}, nil
			}
			if spot.Ticket != nil { // spot taken, must start over
				slog.Debug("spot reserved, resetting", "spot", spot)
				contiguousSpots = 0
				continue
			}
		}
	}
	return nil, ErrNoAvailability
}

func (p *Garage) Checkin(vehicle Vehicle) (*Ticket, error) {
	reservation, err := p.Reserve(vehicle)
	if err != nil {
		return nil, err
	}

	// TODO: attach reservation to ticket
	// TODO: mark spots as reserved
	// TODO: open gate

	slog.Info("do something with spots ", "reservation", reservation)
	return NewTicket(), nil
}

func (p *Garage) Checkout(ticket *Ticket) (*Receipt, error) {
	if ticket.Checkout == nil {
		ticket.Checkout = lo.ToPtr(time.Now())
	}

	// TODO: payment

	receipt := NewReceipt(ticket)
	receipt.Checkout = *ticket.Checkout

	// TODO: release spots
	// TODO: open gate

	return receipt, nil
}

type Level struct {
	Spots []*Spot
}

func NewLevel(spots []*Spot) *Level {
	return &Level{Spots: spots}
}

type Spot struct {
	Ticket *Ticket
}

func NewSpot() *Spot {
	return &Spot{}
}

type Ticket struct {
	Checkin  time.Time
	Checkout *time.Time
	// TODO: spot
}

func NewTicket() *Ticket {
	return &Ticket{
		Checkin: time.Now(), // TODO use time location
	}
}

type Plate struct {
	country       string
	stateProvince string
	number        string
}

type Vehicle interface {
	Size() int
	Ticket() *Ticket
}

type BaseVehicle struct { // abstract
	ticket *Ticket
}

func (b *BaseVehicle) Ticket() *Ticket {
	return b.ticket
}

type Motorcycle struct {
	BaseVehicle
}

func NewMotorcycle() Vehicle { return &Motorcycle{} }

func (m *Motorcycle) Size() int { return 1 }

type Car struct {
	BaseVehicle
}

func NewCar() Vehicle { return &Car{} }

func (c *Car) Size() int { return 1 }

type Bus struct {
	BaseVehicle
}

func NewBus() Vehicle { return &Bus{} }

func (b *Bus) Size() int { return 5 }
