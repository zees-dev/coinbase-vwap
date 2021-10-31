package vwap

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_SatisfiesVwapperInterface(t *testing.T) {
	is := assert.New(t)
	is.Implements((*Vwapper)(nil), NewVWAPCounter(0))
}

func Test_vwap(t *testing.T) {
	is := assert.New(t)

	tests := []struct {
		datapoints   []datapoint
		expectedVWAP decimal.Decimal
		windowSize   int
	}{
		{
			datapoints: []datapoint{
				{Price: "90", Volume: "10"},
				{Price: "91", Volume: "10"},
				{Price: "92", Volume: "10"},
				{Price: "100", Volume: "500"},
			},
			expectedVWAP: decimal.NewFromFloat(99.49),
			windowSize:   4,
		},
		{
			datapoints: []datapoint{
				{Price: "61036.26", Volume: "0.00054368"},
				{Price: "61037.18", Volume: "0.00243673"},
				{Price: "61038.24", Volume: "0.01362416"},
				{Price: "61036.46", Volume: "0.00203768"},
			},
			expectedVWAP: decimal.NewFromFloat(61037.84914103179),
			windowSize:   4,
		},
		{
			datapoints: []datapoint{
				{Price: "90", Volume: "10"},
				{Price: "91", Volume: "10"},
				{Price: "92", Volume: "10"},
				{Price: "100", Volume: "500"},
			},
			expectedVWAP: decimal.NewFromFloat(99.84313725490196),
			windowSize:   2,
		},
		{
			datapoints: []datapoint{
				{Price: "61036.26", Volume: "0.00054368"},
				{Price: "61037.18", Volume: "0.00243673"},
				{Price: "61038.24", Volume: "0.01362416"},
				{Price: "61036.46", Volume: "0.00203768"},
			},
			expectedVWAP: decimal.NewFromFloat(61038.00841351974),
			windowSize:   2,
		},
	}
	for _, test := range tests {
		vwap := NewVWAPCounter(test.windowSize)
		for _, dp := range test.datapoints {
			err := vwap.Update(dp.Price, dp.Volume)
			is.NoError(err)
		}
		gotResult := vwap.VWAP()
		is.Equal(test.expectedVWAP.Round(2), gotResult.Round(2))
	}
}
