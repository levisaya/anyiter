package anyiter

// Iter is the main iterator type
type Iter interface {
	Value() SafeValue
	Type() SafeType
	Next() Iter
}
