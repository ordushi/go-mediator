package mediator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mType struct {
	test string
}

var mtr Observable[mType, string]
var mediator Mediator[mType, string]

func init() {
	x := func(mp *MediatePayload[mType, string]) {
		mp.Response = "yay"
	}
	mtr = New[mType, string]()
	mediator = mtr.NewMediator("add", x)

}

func TestMediator(t *testing.T) {
	x := func(mp *MediatePayload[mType, string]) {
		assert.Equal(t, mp.Payload.test, "test?")
		mp.Response = "yay"
	}
	// need to create locker for obs and subscribers
	//time.Sleep(1 * time.Second)
	mediator := mtr.NewMediator("add", x)

	resp := mediator.Mediate(mType{test: "test?"})
	assert.Equal(t, resp, "yay")

}
func BenchmarkMediator(b *testing.B) {

	for n := 0; n < b.N; n++ {
		_ = mediator.Mediate(mType{test: "test?"})
	}

}
