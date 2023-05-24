package orderbook

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestSideString(t *testing.T) {
	assert.Equal(t, Back.String(), "Back")
	assert.Equal(t, Lay.String(), "Lay")
}

func TestNewOrderBack(t *testing.T) {
	price := decimal.NewFromFloat(1.95)
	stake := decimal.NewFromFloat(100.0)
	o, err := NewOrder(Back, price, stake)

	assert.Nil(t, err)
	assert.Equal(t, o.Side, Back)
	assert.NotEmpty(t, o.Id)
	assert.LessOrEqual(t, o.CreatedAt, time.Now().UnixNano())
	assert.Equal(t, o.Stake, stake)
	assert.Equal(t, o.Price, price)
}

func TestNewOrderLay(t *testing.T) {
	price := decimal.NewFromFloat(1.95)
	stake := decimal.NewFromFloat(100.0)
	o, err := NewOrder(Lay, price, stake)

	assert.Nil(t, err)
	assert.Equal(t, o.Side, Lay)
}

func TestNewOrderZeroStake(t *testing.T) {
	price := decimal.NewFromFloat(1.96)
	stake := decimal.NewFromFloat(0.0)

	o, err := NewOrder(Back, price, stake)

	assert.Nil(t, o)
	assert.Equal(t, err, ErrInvalidStake)
}

func TestNewOrderInvalidPrice(t *testing.T) {
	price := decimal.NewFromFloat(1.0)
	stake := decimal.NewFromFloat(100.0)

	o, err := NewOrder(Back, price, stake)

	assert.Nil(t, o)
	assert.Equal(t, err, ErrInvalidOrderPrice)
}

func TestNewLimit(t *testing.T) {
	price := decimal.NewFromFloat(1.95)
	l, err := NewLimit(price)

	assert.Nil(t, err)
	assert.NotNil(t, l)
	assert.Equal(t, l.TotalVolume, decimal.Zero)
	assert.Equal(t, l.Orders.Len(), 0)
}

func TestNewLimitInvalidPrice(t *testing.T) {
	price := decimal.NewFromFloat(0.95)
	l, err := NewLimit(price)

	assert.Nil(t, l)
	assert.Equal(t, err, ErrInvalidLimitPrice)
}

func TestAddOrder(t *testing.T) {
	price := decimal.NewFromFloat(1.95)
	l, _ := NewLimit(price)
	volumeBefore := l.TotalVolume
	o := createTestOrder(t)

	e, err := l.AddOrder(o)

	assert.Nil(t, err)
	assert.Equal(t, e.Value.(*Order), o)
	assert.True(t, l.TotalVolume.Equal(o.Stake))
	assert.True(t, l.TotalVolume.GreaterThan(volumeBefore))
}

func TestAddOrderPriceMismatch(t *testing.T) {
	limitPrice := decimal.NewFromFloat(1.96)

	l, _ := NewLimit(limitPrice)
	o := createTestOrder(t)

	e, err := l.AddOrder(o)

	assert.Nil(t, e)
	assert.Equal(t, err, ErrPriceMismatch)
}

func TestRemoveOrder(t *testing.T) {
	price := decimal.NewFromFloat(1.95)
	l, _ := NewLimit(price)
	o := createTestOrder(t)

	e, _ := l.AddOrder(o)
	volumeBefore := l.TotalVolume

	removedOrder := l.RemoveOrder(e)
	assert.Equal(t, removedOrder, o)
	assert.True(t, l.TotalVolume.Equal(decimal.Zero))
	assert.True(t, l.TotalVolume.LessThan(volumeBefore))
}

func TestNewOrderbook(t *testing.T) {
	ob := NewOrderbook()

	assert.NotNil(t, ob)
}

func TestOrderbookAddOrder(t *testing.T) {
	ob := NewOrderbook()
	o := createTestOrder(t)

	p, err := ob.AddOrder(o)

	assert.Nil(t, p)
	assert.Nil(t, err)
}

func TestOrderbookAddOrderDuplicate(t *testing.T) {
	ob := NewOrderbook()
	o1 := createTestOrder(t)

	_, err := ob.AddOrder(o1)
	assert.Nil(t, err)

	_, err = ob.AddOrder(o1)
	assert.Equal(t, err, ErrOrderExists)
}

func TestPlaceOrderWithoutLimit(t *testing.T) {
	ob := NewOrderbook()
	o := createTestOrder(t)
	limits := make(map[string]*Limit)

	e, err := ob.PlaceOrder(o, limits)
	assert.Nil(t, err)
	assert.NotNil(t, e)
}

func TestPlaceOrderWithLimit(t *testing.T) {
	ob := NewOrderbook()
	o := createTestOrder(t)
	limits := make(map[string]*Limit)

	limit, err := NewLimit(o.Price)
	assert.Nil(t, err)

	strPrice := o.Price.String()
	limits[strPrice] = limit

	e, err := ob.PlaceOrder(o, limits)
	assert.Nil(t, err)
	assert.NotNil(t, e)
}

// createTestOrder - helper to create proper test order
func createTestOrder(t *testing.T) (o *Order) {
	t.Helper()

	price := decimal.NewFromFloat(1.95)
	stake := decimal.NewFromFloat(100.0)
	o, err := NewOrder(Lay, price, stake)

	if err != nil {
		t.Fatalf("Error creating test order: %s", err)
	}
	return
}
