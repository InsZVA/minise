package main

import "testing"

func TestExample1(t *testing.T) {
	var c [20]byte
	hashCode("abc", &c)
	var d [20]byte
	d = c
	t.Log(c, d)
	c[0] = 8
	t.Log(c, d)
}