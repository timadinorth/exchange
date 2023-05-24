package orderbook

import (
	"container/list"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrInvalidStake      = errors.New("orderbook: invalid order stake")
	ErrInvalidOrderPrice = errors.New("orderbook: invalid order price")
	ErrInvalidLimitPrice = errors.New("orderbook: invalid limit price")
	ErrPriceMismatch     = errors.New("orderbook: limit and order price mismatch")
	ErrOrderExists       = errors.New("orderbook: order id already exists")
)

// Side represents type of the order Back or Lay
type Side int

const (
	Back Side = iota
	Lay
)

func (s Side) String() string {
	return [...]string{"Back", "Lay"}[s]
}

// Order
var minPrice = decimal.NewFromFloat(1.0)

type Order struct {
	Side      Side
	Id        string
	Price     decimal.Decimal
	Stake     decimal.Decimal
	CreatedAt int64
}

func (o Order) String() string {
	return fmt.Sprintf("[Id: %s, Side: %s, Stake: %s, Price: %s]", o.Id, o.Side, o.Stake, o.Price)
}

func NewOrder(side Side, price, stake decimal.Decimal) (*Order, error) {
	if stake.Sign() <= 0 {
		return nil, ErrInvalidStake
	}

	if price.LessThanOrEqual(minPrice) {
		return nil, ErrInvalidOrderPrice
	}

	return &Order{
		Id:        uuid.New().String(), // TODO: replace with actual id
		Side:      side,
		Price:     price,
		Stake:     stake,
		CreatedAt: time.Now().UnixNano(),
	}, nil
}

// Limit - price level in DOM
type Limit struct {
	Price       decimal.Decimal
	TotalVolume decimal.Decimal
	Orders      *list.List
}

func (l Limit) String() string {
	return fmt.Sprintf("[Price: %s, TotalVolume: %s, Len: %d]", l.Price, l.TotalVolume, l.Orders.Len())
}

func NewLimit(price decimal.Decimal) (*Limit, error) {
	if price.LessThanOrEqual(minPrice) {
		return nil, ErrInvalidLimitPrice
	}

	return &Limit{
		Price:       price,
		TotalVolume: decimal.Zero,
		Orders:      list.New(),
	}, nil
}

func (l *Limit) AddOrder(o *Order) (*list.Element, error) {
	if !l.Price.Equal(o.Price) {
		return nil, ErrPriceMismatch
	}

	l.TotalVolume = l.TotalVolume.Add(o.Stake)
	return l.Orders.PushBack(o), nil
}

func (l *Limit) RemoveOrder(e *list.Element) *Order {
	l.TotalVolume = l.TotalVolume.Sub(e.Value.(*Order).Stake)
	return l.Orders.Remove(e).(*Order)
}

// Orderbook
type Orderbook struct {
	orders map[string]*list.Element // from Limit.Orders list

	backLevels map[string]*Limit
	layLevels  map[string]*Limit

	backBest decimal.Decimal
	layBest  decimal.Decimal
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		orders:     map[string]*list.Element{},
		backLevels: make(map[string]*Limit),
		layLevels:  make(map[string]*Limit),
		backBest:   decimal.Zero,
		layBest:    decimal.Zero,
	}
}

// PlaceOrder place order in DOM without filling
func (ob *Orderbook) PlaceOrder(o *Order, limits map[string]*Limit) (e *list.Element, err error) {
	strPrice := o.Price.String()
	limit, ok := limits[strPrice]

	if !ok {
		limit, err = NewLimit(o.Price)
		if err != nil {
			return
		}
		limits[strPrice] = limit
	}

	e, err = limit.AddOrder(o)
	return
}

// FillOrder - filling order by removing liquidity from market
func (ob *Orderbook) FillOrder(o *Order) (*list.Element, error) {
	return nil, nil
}

func (ob *Orderbook) AddOrder(o *Order) (partial *Order, err error) {
	if _, ok := ob.orders[o.Id]; ok {
		return nil, ErrOrderExists
	}

	var e *list.Element

	if o.Side == Back {
		if ob.layBest == decimal.Zero || o.Price.LessThan(ob.layBest) {
			e, err = ob.PlaceOrder(o, ob.backLevels)
			if err != nil {
				return
			}
		} else {
			ob.FillOrder(o)
		}
	} else {
		if ob.backBest == decimal.Zero || o.Price.GreaterThan(ob.backBest) {
			e, err = ob.PlaceOrder(o, ob.layLevels)
			if err != nil {
				return
			}
		} else {
			ob.FillOrder(o)
		}
	}

	ob.orders[o.Id] = e

	return nil, nil
}
