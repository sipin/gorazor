package gorazor

import "sync"

type layoutManager struct {
	// For gorazor just process on single one gohtml file now
	// we use an singleton map to keep layout relationship
	// Not a good solution but works
	layoutMap map[string][]string
}

var single *layoutManager
var mutexLock sync.RWMutex

// LayoutArgs returns arguments of given layout file
func LayoutArgs(file string) []string {
	mutexLock.RLock()
	defer mutexLock.RUnlock()

	if args, ok := single.layoutMap[file]; ok {
		return args
	}
	return []string{}
}

// SetLayout arguments for layout file
func SetLayout(file string, args []string) {
	mutexLock.Lock()
	defer mutexLock.Unlock()

	single.layoutMap[file] = args
}

func init() {
	single = &layoutManager{}
	single.layoutMap = map[string][]string{}
}
