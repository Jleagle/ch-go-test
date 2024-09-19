package chive

type Value struct {
	Text      string
	Number    int64
	Bool      bool
	Float     float64
	Timestamp int64 // Unix timestamp
	Array     *RepeatedValues
}

type RepeatedValues struct {
	Strings  []string
	Ints     []int64
	KeyValue map[string][]byte
}
