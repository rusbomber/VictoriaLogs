package logstorage

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// asyncTaskType identifies the type of background (asynchronous) task attached to a partition.
// More types can be added in the future (e.g. compaction, ttl, schema-changes).
type asyncTaskType string

const (
	asyncTaskNone   asyncTaskType = ""       // no task
	asyncTaskDelete asyncTaskType = "delete" // delete rows matching a query
)

// Status field tracks the outcome of the task.
type asyncTaskStatus string

const (
	taskPending asyncTaskStatus = "pending"
	taskSuccess asyncTaskStatus = "success"
	taskError   asyncTaskStatus = "error"
)

// asyncDeletePayload contains arguments for async delete tasks.
type asyncDeletePayload struct {
	Query string `json:"query,omitempty"`
}

type asyncTask struct {
	Type        asyncTaskType      `json:"type"`                  // task type (delete, etc.)
	TenantIDs   []TenantID         `json:"tenantIDs,omitempty"`   // affected tenants
	Payload     asyncDeletePayload `json:"payload"`               // task parameters
	Seq         uint64             `json:"seq,omitempty"`         // global sequence number
	Status      asyncTaskStatus    `json:"status,omitempty"`      // pending/success/error
	CreatedTime int64              `json:"createdTime,omitempty"` // creation time (UnixNano)
	DoneTime    int64              `json:"doneTime,omitempty"`    // completion time (UnixNano)
	ErrorMsg    string             `json:"error,omitempty"`       // last error message
}

type asyncTasks struct {
	pt *partition // parent partition that owns these async tasks

	mu  sync.Mutex
	ts  []asyncTask   // list of async tasks for this partition
	seq atomic.Uint64 // sequence number for tracking task progress
}

func newAsyncTasks(pt *partition, tasks []asyncTask) *asyncTasks {
	ast := &asyncTasks{
		pt: pt,
		ts: tasks,
	}
	return ast
}

func (at *asyncTasks) nextPendingTask() asyncTask {
	var result asyncTask

	at.mu.Lock()
	for i := range at.ts {
		task := at.ts[i]
		if task.Status == taskPending {
			result = task
			break
		}
	}
	at.mu.Unlock()

	at.seq.Store(result.Seq)
	return result
}

// resolve updates task status and persists the change to disk.
// It holds the internal mutex only for in‐memory modification; the slow fs write
// is executed after the lock is released to avoid blocking other readers.
func (at *asyncTasks) resolve(seq uint64, err error) {
	status, errMsg := taskSuccess, ""
	if err != nil {
		status, errMsg = taskError, err.Error()
	}

	at.mu.Lock()
	for i := range at.ts {
		t := &at.ts[i]

		if t.Seq < seq {
			continue
		}
		if t.Seq > seq || t.Status != taskPending {
			at.mu.Unlock()
			return // no matching pending task
		}

		t.Status = status
		t.DoneTime = time.Now().UnixNano()
		t.ErrorMsg = errMsg
		at.mu.Unlock()

		at.pt.mustSaveAsyncTasks()
		return
	}
	at.mu.Unlock()
}

// addDeleteTask appends a delete task to the partition's task list
func (at *asyncTasks) addDeleteTaskSync(tenantIDs []TenantID, q *Query, seq uint64) uint64 {
	task := asyncTask{
		Seq:         seq,
		Type:        asyncTaskDelete,
		TenantIDs:   append([]TenantID(nil), tenantIDs...),
		Payload:     asyncDeletePayload{Query: q.String()},
		Status:      taskPending,
		CreatedTime: time.Now().UnixNano(),
	}

	at.mu.Lock()
	at.ts = append(at.ts, task)
	at.mu.Unlock()

	at.pt.mustSaveAsyncTasks()
	return seq
}

// unmarshalAsyncTasks converts JSON data back to async tasks
func unmarshalAsyncTasks(data []byte) ([]asyncTask, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var tasks []asyncTask
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("unmarshal async tasks: %w", err)
	}
	return tasks, nil
}
