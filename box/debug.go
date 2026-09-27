package box

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"
	"unsafe"
)

func WriteMemoryUsage() string {
	// var now = time.Now().Unix()
	// if now-lastStatsRefresh < 1 { // 1 second threshold
	// 	return unsafe.String(unsafe.SliceData(memBuf), len(memBuf))
	// }

	runtime.ReadMemStats(&memStats)
	// lastStatsRefresh = now

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memBuf = memBuf[:0]

	memBuf = appendByteSize(memBuf, int(m.Sys))
	memBuf = append(memBuf, " obtained from OS\n"...)
	memBuf = appendByteSize(memBuf, int(m.HeapSys))
	memBuf = append(memBuf, " heap obtained from OS\n"...)
	memBuf = appendByteSize(memBuf, int(m.TotalAlloc))
	memBuf = append(memBuf, " heap total (allocated since start)\n"...)
	memBuf = appendByteSize(memBuf, int(m.HeapAlloc))
	memBuf = append(memBuf, " heap live objects\n"...)
	memBuf = appendByteSize(memBuf, int(m.HeapInuse))
	memBuf = append(memBuf, " heap spans in use\n"...)
	memBuf = appendByteSize(memBuf, int(m.HeapIdle))
	memBuf = append(memBuf, " heap idle\n"...)
	memBuf = appendByteSize(memBuf, int(m.HeapReleased))
	memBuf = append(memBuf, " heap returned to OS\n"...)
	memBuf = appendThousands(memBuf, m.Mallocs)
	memBuf = append(memBuf, " heap objects total (allocated since start)\n"...)
	memBuf = appendThousands(memBuf, m.HeapObjects)
	memBuf = append(memBuf, " heap objects alive\n"...)
	memBuf = appendThousands(memBuf, m.Frees)
	memBuf = append(memBuf, " heap objects freed\n\n"...)

	memBuf = appendByteSize(memBuf, int(m.StackSys))
	memBuf = append(memBuf, " stack obtained from OS\n"...)
	memBuf = appendByteSize(memBuf, int(m.StackInuse))
	memBuf = append(memBuf, " stack in use\n"...)
	memBuf = appendByteSize(memBuf, int(m.OtherSys))
	memBuf = append(memBuf, " misc runtime overhead\n\n"...)

	memBuf = appendThousands(memBuf, uint64(m.NumGC))
	memBuf = append(memBuf, " GC total triggers\n"...)
	memBuf = strconv.AppendUint(memBuf, uint64(m.NumForcedGC), 10)
	memBuf = append(memBuf, " GC manual triggers\n"...)
	memBuf = appendByteSize(memBuf, int(m.NextGC))
	memBuf = append(memBuf, " GC next heap target\n"...)
	memBuf = strconv.AppendFloat(memBuf, float64(m.PauseTotalNs)/1e9, 'f', 2, 64)
	memBuf = append(memBuf, "s GC total time spent\n"...)
	if m.LastGC == 0 {
		memBuf = append(memBuf, "GC never triggered\n"...)
	} else {
		memBuf = strconv.AppendFloat(memBuf, time.Since(time.Unix(0, int64(m.LastGC))).Seconds(), 'f', 1, 64)
		memBuf = append(memBuf, "s GC since last trigger\n"...)
	}
	return unsafe.String(unsafe.SliceData(memBuf), len(memBuf))
}
func WriteFPS() string {
	statsBuf = statsBuf[:0]
	statsBuf = append(statsBuf, "FPS: "...)
	statsBuf = strconv.AppendInt(statsBuf, int64(CurrentFPS), 10)
	return unsafe.String(unsafe.SliceData(statsBuf), len(statsBuf))
}

func ProfileAllocations(seconds float32) {
	go func() {
		var ts = time.Now().Format("2006-01-02_15-04-05")
		var profileFile = fmt.Sprintf("allocs_%s.prof", ts)

		log.Printf("Allocation profiling: capturing for %.2f seconds...\n", seconds)

		var duration = time.Duration(float64(seconds) * float64(time.Second))
		time.Sleep(duration)

		runtime.GC() // flush pending frees so the snapshot is accurate

		var f, err = os.Create(profileFile)
		if err != nil {
			log.Println("could not create allocs profile:", err)
			return
		}
		defer f.Close()

		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			log.Println("could not write allocs profile:", err)
			return
		}

		log.Println("Allocation profile saved at", profileFile)
		log.Println("Opening browser at http://localhost:8081 ...")

		exec.Command("go", "tool", "pprof", "-http=:8081", profileFile).Start()
	}()
}

// private ========================================================

var memStats runtime.MemStats
var memBuf []byte
var lastStatsRefresh int64
var statsBuf []byte

func formatMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memBuf = memBuf[:0]

	memBuf = appendByteSize(memBuf, int(m.Sys))
	memBuf = append(memBuf, " obtained from OS\n"...)
	memBuf = appendByteSize(memBuf, int(m.HeapSys))
	memBuf = append(memBuf, " heap obtained from OS\n"...)
	memBuf = appendByteSize(memBuf, int(m.TotalAlloc))
	memBuf = append(memBuf, " heap total (allocated since start)\n"...)
	memBuf = appendByteSize(memBuf, int(m.HeapAlloc))
	memBuf = append(memBuf, " heap live objects\n"...)
	memBuf = appendByteSize(memBuf, int(m.HeapInuse))
	memBuf = append(memBuf, " heap spans in use\n"...)
	memBuf = appendByteSize(memBuf, int(m.HeapIdle))
	memBuf = append(memBuf, " heap idle\n"...)
	memBuf = appendByteSize(memBuf, int(m.HeapReleased))
	memBuf = append(memBuf, " heap returned to OS\n"...)
	memBuf = appendThousands(memBuf, m.Mallocs)
	memBuf = append(memBuf, " heap objects total (allocated since start)\n"...)
	memBuf = appendThousands(memBuf, m.HeapObjects)
	memBuf = append(memBuf, " heap objects alive\n"...)
	memBuf = appendThousands(memBuf, m.Frees)
	memBuf = append(memBuf, " heap objects freed\n\n"...)

	memBuf = appendByteSize(memBuf, int(m.StackSys))
	memBuf = append(memBuf, " stack obtained from OS\n"...)
	memBuf = appendByteSize(memBuf, int(m.StackInuse))
	memBuf = append(memBuf, " stack in use\n"...)
	memBuf = appendByteSize(memBuf, int(m.OtherSys))
	memBuf = append(memBuf, " misc runtime overhead\n\n"...)

	memBuf = appendThousands(memBuf, uint64(m.NumGC))
	memBuf = append(memBuf, " GC total triggers\n"...)
	memBuf = strconv.AppendUint(memBuf, uint64(m.NumForcedGC), 10)
	memBuf = append(memBuf, " GC manual triggers\n"...)
	memBuf = appendByteSize(memBuf, int(m.NextGC))
	memBuf = append(memBuf, " GC next heap target\n"...)
	memBuf = strconv.AppendFloat(memBuf, float64(m.PauseTotalNs)/1e9, 'f', 2, 64)
	memBuf = append(memBuf, "s GC total time spent\n"...)
	if m.LastGC == 0 {
		memBuf = append(memBuf, "GC never triggered\n"...)
	} else {
		memBuf = strconv.AppendFloat(memBuf, time.Since(time.Unix(0, int64(m.LastGC))).Seconds(), 'f', 1, 64)
		memBuf = append(memBuf, "s GC since last trigger\n"...)
	}
}

func appendByteSize(buf []byte, n int) []byte {
	const unit = 1024
	if n < unit {
		buf = strconv.AppendInt(buf, int64(n), 10)
		return append(buf, " B"...)
	}

	var exp = 0
	var divisor int64 = 1
	// Find the appropriate unit (K, M, G, etc.)
	for v := n / unit; v >= unit; v /= unit {
		divisor *= unit
		exp++
	}
	// We need one more multiplication because the loop stops early
	divisor *= unit

	// Calculate whole and fractional parts (3 decimal places)
	// Example: 1536 bytes -> 1.500 KB
	// whole = 1536 / 1024 = 1
	// fraction = (1536 % 1024) * 1000 / 1024 = 500
	whole := int64(n) / divisor
	fraction := (int64(n) % divisor) * 1000 / divisor

	buf = strconv.AppendInt(buf, whole, 10)
	buf = append(buf, '.')

	// Ensure leading zeros in fraction (e.g., .005 instead of .5)
	if fraction < 100 {
		buf = append(buf, '0')
	}
	if fraction < 10 {
		buf = append(buf, '0')
	}

	buf = strconv.AppendInt(buf, fraction, 10)
	buf = append(buf, ' ')
	buf = append(buf, "KMGTPE"[exp])
	return append(buf, 'B')
}
func appendThousands(buf []byte, n uint64) []byte {
	var tmp [24]byte
	var digits = strconv.AppendUint(tmp[:0], n, 10)
	var l = len(digits)
	for i, c := range digits {
		if i > 0 && (l-i)%3 == 0 {
			buf = append(buf, ' ')
		}
		buf = append(buf, c)
	}
	return buf
}
func appendFixedPoint(buf []byte, v100 int64) []byte {
	whole := v100 / 100
	frac := v100 % 100

	buf = strconv.AppendInt(buf, whole, 10)
	buf = append(buf, '.')
	if frac < 10 {
		buf = append(buf, '0') // Leading zero for .05 vs .5
	}
	return strconv.AppendInt(buf, frac, 10)
}
