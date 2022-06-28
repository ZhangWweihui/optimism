package derive

import (
	"errors"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-node/eth"
)

var ReorgErr = errors.New("reorg")

// Progress represents the progress of a derivation stage:
// the input L1 block that is being processed, and whether it's fully processed yet.
type Progress struct {
	Origin eth.L1BlockRef
	// Closed means that the Current has no more data that the stage may need.
	Closed bool
}

func (pr *Progress) Update(outer Progress) (changed bool, err error) {
	if pr.Closed {
		if outer.Closed {
			if pr.Origin != outer.Origin {
				return true, fmt.Errorf("outer stage changed origin from %s to %s without opening it", pr.Origin, outer.Origin)
			}
			return false, nil
		} else {
			if pr.Origin.Hash != outer.Origin.ParentHash {
				return true, fmt.Errorf("detected internal pipeline reorg of L1 origin data from %s to %s: %w", pr.Origin, outer.Origin, ReorgErr)
			}
			pr.Origin = outer.Origin
			pr.Closed = false
			return true, nil
		}
	} else {
		if pr.Origin != outer.Origin {
			return true, fmt.Errorf("outer stage changed origin from %s to %s before closing it", pr.Origin, outer.Origin)
		}
		if outer.Closed {
			pr.Closed = true
			return true, nil
		} else {
			return false, nil
		}
	}
}