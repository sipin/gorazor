package gorazor

import "sync"

// LayoutManager is the layout manager
type LayoutManager struct {
	// For gorazor just process on single one gohtml file now
	// we use an singleton map to keep layout relationship
	// Not a good solution but works
	layoutMap map[string][]string
}

var single *LayoutManager
var mutexLock sync.RWMutex

// LayoutArgs returns arguments of given layout file
func LayoutArgs(file string) []string {
	mutexLock.RLock()
	defer mutexLock.RUnlock()
	manager := newManager()
	if args, ok := manager.layoutMap[file]; ok {
		return args
	}
	return []string{}
}

// SetLayout arguments for layout file
func SetLayout(file string, args []string) {
	mutexLock.Lock()
	manager := newManager()
	manager.layoutMap[file] = args
	mutexLock.Unlock()
}

func newManager() *LayoutManager {
	if single != nil {
		return single
	}
	lay := &LayoutManager{}
	lay.layoutMap = map[string][]string{}
	single = lay
	return lay
}
