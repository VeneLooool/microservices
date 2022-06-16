package randService

import "math/rand"

func RAND(err error) error {
	if rand.Intn(10) == 5 {
		return err
	}
	return nil
}
