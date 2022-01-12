package mediator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mType struct {
	test string
}

var mtr Observable[mType, string]

func init() {
	mtr = New[mType, string]()
}

func TestMediator(t *testing.T) {
	x := func(mp *MediatePayload[mType, string]) {
		assert.Equal(t, mp.Payload.test, "test?")
		mp.Response = "yay"
	}
	mediator := mtr.NewMediator("add", x)
	// need to create locker for obs and subscribers
	//time.Sleep(1 * time.Second)
	resp := mediator.Mediate(mType{test: "test?"})
	assert.Equal(t, resp, "yay")

}

