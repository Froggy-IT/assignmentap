package worker

import (
	"fmt"
	"time"
)

type StatsProvider interface {
	RequestCount() int64
	KeyCount() int
}

func StartWorker(stop <-chan struct{}, provider StatsProvider) {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Printf("[WORKER] requests=%d keys=%d\n",
					provider.RequestCount(),
					provider.KeyCount(),
				)
			case <-stop:
				fmt.Println("[WORKER] stopped")
				return
			}
		}
	}()
}
