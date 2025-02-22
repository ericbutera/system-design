package garage_test

import (
	"testing"
	"time"

	"parking-garage/garage"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// create a parking lot with the specified level and space counts
func newGarage(t *testing.T, levelCount int, spotCount int) *garage.Garage {
	t.Helper()

	levels := make([]*garage.Level, levelCount)
	for x := 0; x < levelCount; x++ {
		spots := make([]*garage.Spot, spotCount)
		for y := 0; y < spotCount; y++ {
			spots[y] = garage.NewSpot()
		}
		levels[x] = garage.NewLevel(spots)
	}

	return garage.New("super parking lot", levels)
}

func TestReservation(t *testing.T) {
	t.Parallel()

	t.Run("car reservation", func(t *testing.T) {
		t.Parallel()
		lot := newGarage(t, 1, 1)
		level := lot.Levels[0]

		reservation, err := lot.Reserve(garage.NewCar())
		require.NoError(t, err)
		assert.Equal(t, level, reservation.Level)
		assert.Equal(t, level.Spots, reservation.Spots)
	})

	t.Run("bus reservation", func(t *testing.T) {
		t.Parallel()
		lot := newGarage(t, 1, 5)
		level := lot.Levels[0]

		reservation, err := lot.Reserve(garage.NewBus())
		require.NoError(t, err)
		assert.Equal(t, level, reservation.Level)
		assert.Equal(t, level.Spots, reservation.Spots)
	})

	t.Run("bus too big", func(t *testing.T) {
		t.Parallel()
		lot := newGarage(t, 1, 1)

		_, err := lot.Reserve(garage.NewBus())
		require.ErrorIs(t, err, garage.ErrNoAvailability)
	})

	t.Run("ticket resets counter", func(t *testing.T) {
		// 0 open
		// 1 open
		// 2 [CAR] <- car resets spaces
		// 3 [BUS]
		// 4 [BUS]
		// 5 [BUS]
		// 6 [BUS]
		// 7 [BUS]
		lot := newGarage(t, 1, 8)
		lot.Levels[0].Spots[2].Ticket = &garage.Ticket{}
		spots := lot.Levels[0].Spots[3:5]

		bus := garage.NewBus()
		reservation, err := lot.Reserve(bus)
		require.NoError(t, err)
		assert.Equal(t, spots, reservation.Spots)
	})
}

func TestTicketSetsCheckin(t *testing.T) {
	t.Parallel()
	ticket := garage.NewTicket()
	assert.True(t, ticket.Checkin.Before(time.Now()))
}

func TestCheckinMissingSpots(t *testing.T) {
	t.Parallel()

	t.Run("no spots exist", func(t *testing.T) {
		t.Parallel()
		car := garage.NewCar()

		lot := newGarage(t, 1, 0)
		_, err := lot.Checkin(car)
		assert.ErrorIs(t, err, garage.ErrNoAvailability)
	})

	t.Run("spots exist", func(t *testing.T) {
		t.Parallel()
		car := garage.NewCar()

		start := time.Now()
		lot := newGarage(t, 1, 1)
		ticket, err := lot.Checkin(car)
		require.NoError(t, err)
		ticket.Checkin.After(start)
	})
}

func TestCheckout(t *testing.T) {
	t.Parallel()
	ticket := garage.NewTicket()

	lot := newGarage(t, 1, 1)
	receipt, err := lot.Checkout(ticket)
	require.NoError(t, err)

	assert.Equal(t, float32(1), receipt.Total)
}

func TestOneHourRate(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("America/Detroit")

	checkin := time.Date(2025, 01, 01, 00, 0, 0, 0, loc)

	cases := []struct {
		name   string
		ticket *garage.Ticket
		total  float32
	}{
		{
			name: "15 minutes",
			ticket: &garage.Ticket{
				Checkin:  checkin,
				Checkout: lo.ToPtr(checkin.Add(15 * time.Minute)),
			},
			total: float32(1),
		},
		{
			name: "30 minutes",
			ticket: &garage.Ticket{
				Checkin:  checkin,
				Checkout: lo.ToPtr(checkin.Add(30 * time.Minute)),
			},
			total: float32(2),
		},
		{
			name: "45 minutes",
			ticket: &garage.Ticket{
				Checkin:  checkin,
				Checkout: lo.ToPtr(checkin.Add(45 * time.Minute)),
			},
			total: float32(3),
		},
		{
			name: "60 minutes",
			ticket: &garage.Ticket{
				Checkin:  checkin,
				Checkout: lo.ToPtr(checkin.Add(60 * time.Minute)),
			},
			total: float32(4),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			receipt := garage.NewReceipt(tc.ticket)
			assert.Equal(t, tc.total, receipt.Total)
		})
	}
}
