package utils

import (
	"runtime"
	"time"
)

const Memory64KB uint32 = 64 * 1024
const Memory128KB uint32 = 128 * 1024
const Memory256KB uint32 = 256 * 1024
const Memory512KB uint32 = 512 * 1024
const Memory1MB uint32 = 1 * 1024 * 1024
const Memory2MB uint32 = 2 * 1024 * 1024
const Memory64MB uint32 = 64 * 1024 * 1024

const DefaultSessionTTL = TimeTtl30Minutes
const TimeTtl30Minutes time.Duration = 30 * time.Minute

func NumThreads(requested_n_threads int) uint8 {
	if requested_n_threads <= 0 {
		requested_n_threads = runtime.GOMAXPROCS(0)
	} //                                                  /*NumCPU()*/
	return uint8(min(max(requested_n_threads, 1), runtime.GOMAXPROCS(0)))
}
