package vwap

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/zees-dev/coinbase-vwap/pkg/queue"
)

// datapoint represents a single incoming datapoint of the VWAP calculation.
type datapoint struct {
	Price  string
	Volume string
}

// vwapCounter is a volume-weighted-average-price counter that keeps track of the total price and volume.
type vwapCounter struct {
	totalPriceVolProduct decimal.Decimal
	totalVolume          decimal.Decimal
	datapoints           queue.Queue
	lock                 sync.RWMutex
}

// NewVWAPCounter creates a new vwapCounter with the given sliding window size.
func NewVWAPCounter(size int) *vwapCounter {
	return &vwapCounter{
		datapoints: queue.NewFixedQueue(size),
	}
}

// VWAP returns the current running VWAP value.
func (v *vwapCounter) VWAP() decimal.Decimal {
	v.lock.RLock()
	defer v.lock.RUnlock()

	return v.totalPriceVolProduct.Div(v.totalVolume)
}

// Update updates the VWAP counter with a new datapoint.
// The sliding window moves forward by one datapoint if the datapoints queue is full;
// this effectively drops the datapoint in the head of the queue
func (v *vwapCounter) Update(price, volume string) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	// ensure current datapoint is always added to queue
	defer v.datapoints.Enqueue(datapoint{Price: price, Volume: volume})

	// If queue is full, remove the oldest datapoint and subtract its price and volume from the respective totals.
	// This essentialls moves the sliding window forward
	if v.datapoints.IsFull() {
		item, err := v.datapoints.Dequeue()
		if err != nil {
			return errors.Wrap(err, "dequeue")
		}

		dp := item.(datapoint)

		decPrice, err := decimal.NewFromString(dp.Price)
		if err != nil {
			return errors.Wrap(err, "failed to parse price")
		}

		decVolume, err := decimal.NewFromString(dp.Volume)
		if err != nil {
			return errors.Wrap(err, "failed to parse volume")
		}

		// priceVolume -= price * volume
		v.totalPriceVolProduct = v.totalPriceVolProduct.Sub(decPrice.Mul(decVolume))

		// volume -= volume
		v.totalVolume = v.totalVolume.Sub(decVolume)
	}

	decPrice, err := decimal.NewFromString(price)
	if err != nil {
		return errors.Wrap(err, "failed to parse price")
	}

	decVolume, err := decimal.NewFromString(volume)
	if err != nil {
		return errors.Wrap(err, "failed to parse volume")
	}

	// priceVolume += price * volume
	v.totalPriceVolProduct = v.totalPriceVolProduct.Add(decPrice.Mul(decVolume))

	// volume += volume
	v.totalVolume = v.totalVolume.Add(decVolume)

	return nil
}
