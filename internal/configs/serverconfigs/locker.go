package serverconfigs

import "sync"

var sharedLocker = &sync.RWMutex{}
