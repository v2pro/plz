//+build !go1.9

package test

func Helper() func() {
	return func() {
	}
}
