
package mfklib

type Tally struct {
	Marry uint64
	Fuck uint64
	Kill uint64
}

type TripleStats struct {
	Skips uint64

	A Tally
	B Tally
	C Tally
}