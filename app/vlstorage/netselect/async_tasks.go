package netselect

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/VictoriaMetrics/VictoriaLogs/lib/logstorage"
	"golang.org/x/sync/errgroup"
)

// ListAsyncTasks gathers all async tasks from every storage node and returns them along with the originating storage address.
func (s *Storage) ListAsyncTasks(ctx context.Context) ([]logstorage.AsyncTaskInfoWithSource, error) {
	if len(s.sns) == 0 {
		return nil, nil
	}

	g, ctx := errgroup.WithContext(ctx)

	// race-free slices
	results := make([][]logstorage.AsyncTaskInfoWithSource, len(s.sns))
	for i, sn := range s.sns {
		i, sn := i, sn
		g.Go(func() error {
			tasks, err := sn.getAsyncTasks(ctx)
			if err != nil {
				return err
			}
			results[i] = tasks
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	var all []logstorage.AsyncTaskInfoWithSource
	for _, ts := range results {
		all = append(all, ts...)
	}

	return all, nil
}

func (sn *storageNode) getAsyncTasks(ctx context.Context) ([]logstorage.AsyncTaskInfoWithSource, error) {
	args := url.Values{}
	args.Set("version", AsyncTasksProtocolVersion)

	reqURL := sn.getRequestURL("/internal/async_tasks", args)
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	if err := sn.ac.SetHeaders(req, true); err != nil {
		return nil, fmt.Errorf("cannot set auth headers for %q: %w", reqURL, err)
	}

	resp, err := sn.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body from %q: %w", reqURL, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code for %q: %d; response: %q", reqURL, resp.StatusCode, body)
	}

	var tasks []logstorage.AsyncTaskInfoWithSource
	if err := json.Unmarshal(body, &tasks); err != nil {
		return nil, fmt.Errorf("cannot decode async tasks response from %q: %w; response body: %q", reqURL, err, body)
	}

	// Attach origin address.
	for i := range tasks {
		if tasks[i].Storage == "" {
			tasks[i].Storage = sn.addr
		}
	}
	return tasks, nil
}
