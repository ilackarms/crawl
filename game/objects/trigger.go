package objects

import tl "github.com/ilackarms/termloop"

type Position struct {
	X int        `json:"X"`
	Y int        `json:"Y"`
}

type Trigger interface {
	tl.Physical
	TriggerPositions() []Position //x:y map of spots that trigger
}