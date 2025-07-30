package logstorage

import (
	"strings"
	"sync"
	"time"
)

// AsyncTaskInfo represents brief information about a background async task.
// It is intended to be exposed via monitoring endpoints / UI pages.
type AsyncTaskInfo struct {
	Seq     uint64             `json:"seq"`
	Type    asyncTaskType      `json:"type"`
	Status  asyncTaskStatus    `json:"status"`
	Tenant  string             `json:"tenant"`
	Payload asyncDeletePayload `json:"payload"`

	CreatedTime int64  `json:"createdTime,omitempty"`
	DoneTime    int64  `json:"doneTime,omitempty"`
	Error       string `json:"error,omitempty"`
}

// AsyncTaskInfoWithSource is AsyncTaskInfo with the source storage name.
type AsyncTaskInfoWithSource struct {
	AsyncTaskInfo `json:",inline"`
	Storage       string `json:"storage"`
}

// asyncTasksCache is a simple cache valid for 5 seconds, used only by the Storage instance.
// It prevents GET /async_tasks from impacting the performance of all parts,
// since the request requires acquiring the parts lock.
// The frontend also refreshes every 5 seconds, so opening multiple tabs
// could cause issues without this cache.
var asyncTasksCache struct {
	mu   sync.Mutex
	ts   time.Time
	data []AsyncTaskInfo
}

// ListAsyncTasks gathers information about all async tasks known to this
// Storage instance. The returned slice isn't sorted.
func (s *Storage) ListAsyncTasks() []AsyncTaskInfo {
	asyncTasksCache.mu.Lock()
	const cacheTTL = 5 * time.Second
	if time.Since(asyncTasksCache.ts) < cacheTTL && asyncTasksCache.data != nil {
		d := append([]AsyncTaskInfo(nil), asyncTasksCache.data...)
		asyncTasksCache.mu.Unlock()
		return d
	}
	asyncTasksCache.mu.Unlock()

	s.partitionsLock.Lock()
	pws := append([]*partitionWrapper(nil), s.partitions...)
	for _, pw := range pws {
		pw.incRef()
	}
	s.partitionsLock.Unlock()
	defer func() {
		for _, p := range pws {
			p.decRef()
		}
	}()

	var out []AsyncTaskInfo
	for _, p := range pws {
		a := p.pt.ats

		a.mu.Lock()
		tasks := append([]asyncTask(nil), a.ts...)
		a.mu.Unlock()

		for _, t := range tasks {
			tn := "*"
			if len(t.TenantIDs) > 0 {
				var b strings.Builder
				for i, id := range t.TenantIDs {
					if i > 0 {
						b.WriteByte(',')
					}
					b.WriteString(id.String())
				}
				tn = b.String()
			}
			out = append(out, AsyncTaskInfo{
				Seq:         t.Seq,
				Type:        t.Type,
				Status:      t.Status,
				Tenant:      tn,
				Payload:     t.Payload,
				CreatedTime: t.CreatedTime,
				DoneTime:    t.DoneTime,
				Error:       t.ErrorMsg,
			})
		}
	}

	asyncTasksCache.mu.Lock()
	asyncTasksCache.ts = time.Now()
	asyncTasksCache.data = out
	asyncTasksCache.mu.Unlock()

	return out
}
