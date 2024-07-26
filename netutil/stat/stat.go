package stat

import (
	"fmt"
	"os"
	"time"

	"github.com/prometheus/procfs"
)

type ProcessStats struct {
	fs  procfs.FS
	pid int

	lastCPU     time.Duration
	lastCollect time.Time
}

func NewProcessStats() (*ProcessStats, error) {
	pfs, err := procfs.NewDefaultFS()
	if err != nil {
		return nil, fmt.Errorf("error initialising procfs mount: %w", err)
	}

	ps := &ProcessStats{
		fs:  pfs,
		pid: os.Getpid(),
	}

	// sets initial values
	_, err = ps.CPULoadSinceLastCall()
	if err != nil {
		return nil, fmt.Errorf("error initialising metrics: %w", err)
	}

	return ps, nil
}

func (p *ProcessStats) CPULoadSinceLastCall() (float64, error) {
	proc, err := p.fs.Proc(p.pid)
	if err != nil {
		return 0, fmt.Errorf("error reading procfs for process %d: %w", p.pid, err)
	}

	procStat, err := proc.Stat()
	if err != nil {
		return 0, fmt.Errorf("error reading stat for process %d: %w", p.pid, err)
	}

	timeNow := time.Now()
	totalCPU := time.Duration(procStat.CPUTime()) * time.Second

	elapsedWallTime := timeNow.Sub(p.lastCollect)
	elapsedCPUTime := totalCPU - p.lastCPU

	p.lastCPU = totalCPU
	p.lastCollect = timeNow

	return float64(elapsedCPUTime / elapsedWallTime), nil
}
