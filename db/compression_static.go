package db

import "github.com/server-may-cry/highloadcup_18/structures"

const (
	StatusFreeInt int8 = iota
	StatusHoldInt
	StatusAllComplexInt
)
const (
	SexFInt int8 = iota
	SexMInt
)

const (
	StatusFree       = structures.StatusFree
	StatusHold       = structures.StatusHold
	StatusAllComplex = structures.StatusAllComplex
	SexF             = structures.SexF
	SexM             = structures.SexM
)

var statusToInt = map[string]int8{
	StatusAllComplex: StatusAllComplexInt,
	StatusFree:       StatusFreeInt,
	StatusHold:       StatusHoldInt,
}
var intToStatus = [...]string{
	StatusFree,
	StatusHold,
	StatusAllComplex,
}
var sexToInt = map[string]int8{
	SexF: SexFInt,
	SexM: SexMInt,
}
var intToSex = [...]string{
	SexF,
	SexM,
}
