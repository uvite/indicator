package indicator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uvite/indicator/types"
)

func Test_VIDYA(t *testing.T) {
	vidya := &VIDYA{IntervalWindow: types.IntervalWindow{Window: 16}}
	vidya.Update(1)
	assert.Equal(t, vidya.Last(), 1.)
	vidya.Update(2)
	newV := 2./17.*2. + 1.*(1.-2./17.)
	assert.Equal(t, vidya.Last(), newV)
	vidya.Update(1)
	assert.Equal(t, vidya.Last(), vidya.Index(1))
}
