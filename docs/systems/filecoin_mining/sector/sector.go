package sector

import (
	util "github.com/filecoin-project/specs/util"
)

var IMPL_FINISH = util.IMPL_FINISH

type Serialization = util.Serialization

const (
	DeclaredFault StorageFaultType = 1 + iota
	DetectedFault
	TerminatedFault
)
