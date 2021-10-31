package queue

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_SatisfiesQueueInterface(t *testing.T) {
	is := assert.New(t)
	is.Implements((*Queue)(nil), NewFixedQueue(0))
}

func Test_queue(t *testing.T) {
	is := assert.New(t)

	tests := []struct {
		name   string
		size   int
		points []decimal.Decimal
		head   decimal.Decimal
		isFull bool
	}{
		{
			name: "single item enqueue capacity reached",
			size: 1,
			points: []decimal.Decimal{
				decimal.NewFromFloat(1),
			},
			head:   decimal.NewFromFloat(1),
			isFull: true,
		},
		{
			name: "single item enqueue capacity available",
			size: 2,
			points: []decimal.Decimal{
				decimal.NewFromFloat(1),
			},
			head:   decimal.NewFromFloat(1),
			isFull: false,
		},
		{
			name: "slide window with 1 item capacity",
			size: 1,
			points: []decimal.Decimal{
				decimal.NewFromFloat(1),
				decimal.NewFromFloat(2),
			},
			head:   decimal.NewFromFloat(2),
			isFull: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q := NewFixedQueue(test.size)
			for _, point := range test.points {
				q.Enqueue(point)
			}
			is.Equal(test.isFull, q.IsFull())

			head, err := q.Head()
			is.NoError(err)

			is.Equal(head.(decimal.Decimal), test.head)
		})
	}
}
