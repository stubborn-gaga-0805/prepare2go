package runtime

import (
	"context"
	"fmt"
	"github.com/apcera/termtables"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var systemStatus *SystemStatus

type SystemStatus struct {
	ctx        context.Context
	memStats   runtime.MemStats
	stopSignal chan struct{}
}

func NewSystemStatus(ctx context.Context) *SystemStatus {
	if systemStatus == nil {
		systemStatus = &SystemStatus{
			ctx:        ctx,
			stopSignal: make(chan struct{}, 1),
		}
	}
	return systemStatus
}

func (sys *SystemStatus) StartMonitor() {
	fmt.Println("\nSystem Monitor Started...")
	time.Sleep(time.Millisecond * 500)
	sys.printSystemStatus()

	signalCtx, signalStop := signal.NotifyContext(sys.ctx, syscall.SIGINT, syscall.SIGTERM)
	defer signalStop()

	ticker := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-ticker.C:
			sys.printSystemStatus()
		case <-signalCtx.Done():
			fmt.Println("\nSystem Monitor Stopped...")
			ticker.Stop()
			return
		}
	}
}

func (sys *SystemStatus) printSystemStatus() {
	runtime.ReadMemStats(&sys.memStats)
	table := termtables.CreateTable()
	table.AddTitle("current system running status")
	table.AddHeaders("cpu cores", "running goroutines", "allocated memory (bytes)", "total allocated memory(bytes)", "system allocated memory(bytes)", "heap object allocations(bytes)")
	table.AddRow(runtime.NumCPU(), runtime.NumGoroutine(), sys.memStats.Alloc, sys.memStats.TotalAlloc, sys.memStats.Sys, sys.memStats.HeapObjects)
	fmt.Println("\n" + table.Render())
}
