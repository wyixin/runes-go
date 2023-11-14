package retry

import "time"

// RunAble
// return bool
// return error
type RunAble func() (bool, error)

func Exec(f RunAble, max int) (b bool, err error) {
	for i := 0; i < max; i++ {
		b, err = f()
		if b {
			return true, nil
		}
		time.Sleep(time.Second)
	}
	return false, err
}
