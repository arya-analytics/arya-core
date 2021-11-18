package telem

type Value struct {
	TimeStamp float64
	Value     float64
}

type Slice map[int32]Value
