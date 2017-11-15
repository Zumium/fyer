package fragmngr

import (
	"sync"

	"github.com/Zumium/fyer/cfg"
)

var holder fragManagerHolder

type fragManagerHolder struct {
	lock        sync.Mutex
	fragManager FragManager
}

func (fmh *fragManagerHolder) setFragManager(fm FragManager) bool {
	fmh.lock.Lock()
	defer fmh.lock.Unlock()

	if fmh.fragManager != nil {
		return false
	}
	fmh.fragManager = fm
	return true
}

func (fmh *fragManagerHolder) instance() FragManager {
	return fmh.fragManager
}

//====================================================================

//FMInstance returns the current FragManager instance
func FMInstance() FragManager {
	return holder.instance()
}

//Init initializes the fragment manager module
func Init() error {
	return InitSimpleFSFragManager(cfg.FragBasePath())
}
