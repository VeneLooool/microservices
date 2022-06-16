package yellowBank

import (
	"errors"
	consts "gitlab.ozon.dev/VeneLooool/homework-3/const"
	"math/rand"
)

type YellowBank struct {
}

func (YellowBank) ReserveFound(order *consts.Order) error {
	return nil
}
func (YellowBank) Charge(order *consts.Order) error {
	if rand.Intn(10) == 5 {
		return errors.New("denied payment")
	}
	return nil
}
