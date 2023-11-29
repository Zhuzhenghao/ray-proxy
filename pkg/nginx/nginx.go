package nginx

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"syscall"

	"github.com/shirou/gopsutil/v3/process"
)

var reloadLock = sync.Mutex{}

func Reload(ctx context.Context) error {
	reloadLock.Lock()
	defer reloadLock.Unlock()

	processes, err := process.Processes()
	if err != nil {
		return err
	}

	isNginxMaster := func(cmdLine string) bool {
		return strings.Contains(cmdLine, "master") && strings.Contains(cmdLine, "nginx") && !strings.Contains(cmdLine, "grep")
	}

	for _, p := range processes {
		cmdLine, err := p.Cmdline()
		if err != nil {
			continue
		}

		if !isNginxMaster(cmdLine) {
			continue
		}

		return p.SendSignalWithContext(ctx, syscall.SIGHUP)
	}

	return fmt.Errorf("failed to reload nginx")
}
