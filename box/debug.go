package box

import (
	"runtime"
	"strconv"
	"time"
)

var memStats runtime.MemStats
var statsCacheBuf [4096]byte
var statsCache []byte
var lastStatsRefresh time.Time

// WriteMemoryUsage returns a formatted memory statistics report.
// The report is cached and refreshed at most once per second.
func WriteMemoryUsage() []byte {
	if time.Since(lastStatsRefresh) < time.Second {
		return statsCache
	}
	runtime.ReadMemStats(&memStats)
	lastStatsRefresh = time.Now()
	statsCache = formatMemoryUsage(statsCacheBuf[:0])
	return statsCache
}

func formatMemoryUsage(buf []byte) []byte {
	var m = &memStats

	buf = append(buf, "Memory:\n"...)
	buf = append(buf, "  UsedNow   = "...)
	buf = AppendByteSize(buf, int(m.Alloc))
	buf = append(buf, " (current heap in use)\n"...)
	buf = append(buf, "  UsedTotal = "...)
	buf = AppendByteSize(buf, int(m.TotalAlloc))
	buf = append(buf, " (total allocated since start)\n"...)
	buf = append(buf, "  FromOS    = "...)
	buf = AppendByteSize(buf, int(m.Sys))
	buf = append(buf, " (memory reserved from OS)\n"...)

	buf = append(buf, "\nHeap:\n"...)
	buf = append(buf, "  Used      = "...)
	buf = AppendByteSize(buf, int(m.HeapAlloc))
	buf = append(buf, '\n')
	buf = append(buf, "  Reserved  = "...)
	buf = AppendByteSize(buf, int(m.HeapSys))
	buf = append(buf, '\n')
	buf = append(buf, "  Idle      = "...)
	buf = AppendByteSize(buf, int(m.HeapIdle))
	buf = append(buf, " (not used but still reserved)\n"...)
	buf = append(buf, "  Active    = "...)
	buf = AppendByteSize(buf, int(m.HeapInuse))
	buf = append(buf, " (actively in use)\n"...)
	buf = append(buf, "  Released  = "...)
	buf = AppendByteSize(buf, int(m.HeapReleased))
	buf = append(buf, " (given back to OS)\n"...)

	buf = append(buf, "\nStack:\n"...)
	buf = append(buf, "  Used      = "...)
	buf = AppendByteSize(buf, int(m.StackInuse))
	buf = append(buf, '\n')
	buf = append(buf, "  Reserved  = "...)
	buf = AppendByteSize(buf, int(m.StackSys))
	buf = append(buf, '\n')
	buf = append(buf, "  Other     = "...)
	buf = AppendByteSize(buf, int(m.OtherSys))
	buf = append(buf, " (misc runtime overhead)\n"...)

	buf = append(buf, "\nObjects:\n"...)
	buf = append(buf, "  Allocs    = "...)
	buf = appendSeparateThousands(buf, m.Mallocs)
	buf = append(buf, " (objects allocated)\n"...)
	buf = append(buf, "  Frees     = "...)
	buf = appendSeparateThousands(buf, m.Frees)
	buf = append(buf, " (objects freed)\n"...)
	buf = append(buf, "  Live      = "...)
	buf = appendSeparateThousands(buf, m.HeapObjects)
	buf = append(buf, " (currently alive)\n"...)

	buf = append(buf, "\nGarbage Collection:\n"...)
	buf = append(buf, "  Total     = "...)
	buf = appendSeparateThousands(buf, uint64(m.NumGC))
	buf = append(buf, " (total collections)\n"...)
	buf = append(buf, "  Forced    = "...)
	buf = strconv.AppendUint(buf, uint64(m.NumForcedGC), 10)
	buf = append(buf, " (manual triggers)\n"...)
	buf = append(buf, "  Next      = "...)
	buf = AppendByteSize(buf, int(m.NextGC))
	buf = append(buf, " (target heap size of next GC)\n"...)
	buf = append(buf, "  PauseTotal= "...)
	buf = strconv.AppendFloat(buf, float64(m.PauseTotalNs)/1e9, 'f', 2, 64)
	buf = append(buf, " s (total time in GC)\n"...)
	if m.LastGC == 0 {
		buf = append(buf, "  SinceLast = never\n"...)
	} else {
		buf = append(buf, "  SinceLast = "...)
		buf = strconv.AppendFloat(buf, time.Since(time.Unix(0, int64(m.LastGC))).Seconds(), 'f', 2, 64)
		buf = append(buf, " s\n"...)
	}

	return buf
}

// AppendFPS appends "FPS: <n>" to buf.
func AppendFPS(buf []byte, fps int) []byte {
	buf = append(buf, "FPS: "...)
	return strconv.AppendInt(buf, int64(fps), 10)
}

// AppendIdleTPS spins until the current tick budget is exhausted, then appends
// "Idle: <n>" to buf. This call is what enforces the target TPS rate — if it is
// not called, Running() loops as fast as possible.
func AppendIdleTPS(buf []byte) []byte {
	for idling() {}
	buf = append(buf, "Idle: "...)
	return appendSeparateThousands(buf, uint64(CurrentIdleTPS))
}

// AppendTPS appends "TPS: <n>" to buf.
func AppendTPS(buf []byte, tps int) []byte {
	buf = append(buf, "TPS: "...)
	return strconv.AppendInt(buf, int64(tps), 10)
}

// AppendByteSize appends a human-readable byte size (e.g. "1.500 KB") to buf.
func AppendByteSize(buf []byte, n int) []byte {
	const unit = 1024
	if n < unit {
		buf = strconv.AppendInt(buf, int64(n), 10)
		return append(buf, " B"...)
	}
	var div, exp = int(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	buf = strconv.AppendFloat(buf, float64(n)/float64(div), 'f', 3, 64)
	buf = append(buf, ' ')
	buf = append(buf, "KMGTPE"[exp])
	return append(buf, 'B')
}

func appendSeparateThousands(buf []byte, n uint64) []byte {
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
