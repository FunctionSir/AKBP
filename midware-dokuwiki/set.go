package main

type VOID struct{}

type set[Key comparable] map[Key]VOID

func (x set[Key]) Have(key Key) bool {
	_, exists := x[key]
	return exists
}

func (x set[Key]) Ins(key Key) {
	x[key] = VOID{}
}

// func (x set[Key]) Del(key Key) {
// 	delete(x, key)
// }
