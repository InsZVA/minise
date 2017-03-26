package main

import "testing"

func TestHashCodeSlice_Less(t *testing.T) {
	a := hashCodeSlice{[20]byte{1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7,8,9,0}, [20]byte{1,2,3,4,5,6,7,8,9,1,1,2,3,4,5,6,7,8,9,1}}
	t.Log(a.Less(1, 0))
}
