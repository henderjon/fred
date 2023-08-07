package main

type search struct {
	reverse bool // this is NOT forward because forward is the default, but the zero value of a bool false; "forward" should be used in surrounding code
	pattern string
}

// func (s search) isEmpty() bool {
// 	return len(s.pattern) == 0
// }
