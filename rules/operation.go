package rules

import "math"

type Operation int16

const (
	Exclude Operation = iota
	Include
	Noop
	ExcludeAndTerminate
	Invalid = math.MaxInt16
)
