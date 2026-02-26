package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	proxyv1 "github.com/soltiHQ/control-plane/api/proxy/v1"
)

const (
	httpV1Timeout = 10 * time.Second

	v1PathTasks       = "/api/v1/tasks"
	v1PathTasksSubmit = "/api/v1/tasks"
	v1PathTasksExport = "/api/v1/tasks/export"
)

// httpProxyV1 implements AgentProxy over HTTP for API v1.
type httpProxyV1 struct {
	endpoint string
	client   *http.Client
}

func (p *httpProxyV1) ListTasks(ctx context.Context, f TaskFilter) (*proxyv1.TaskListResponse, error) {
	u, err := url.Parse(p.endpoint + v1PathTasks)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadEndpointURL, err)
	}

	q := u.Query()
	if f.Slot != "" {
		q.Set("slot", f.Slot)
	}
	if f.Status != "" {
		q.Set("status", f.Status)
	}
	if f.Limit > 0 {
		q.Set("limit", strconv.Itoa(f.Limit))
	}
	if f.Offset > 0 {
		q.Set("offset", strconv.Itoa(f.Offset))
	}
	u.RawQuery = q.Encode()

	ctx, cancel := context.WithTimeout(ctx, httpV1Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreateRequest, err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode)
	}

	var body proxyv1.TaskListResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecode, err)
	}

	return &body, nil
}

func (p *httpProxyV1) SubmitTask(ctx context.Context, sub TaskSubmission) error {
	u, err := url.Parse(p.endpoint + v1PathTasksSubmit)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBadEndpointURL, err)
	}

	ctx, cancel := context.WithTimeout(ctx, httpV1Timeout)
	defer cancel()

	body := map[string]any{"spec": sub.Spec}
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCreateRequest, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCreateRequest, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSubmitTask, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode)
	}
	return nil
}

func (p *httpProxyV1) ExportSpecs(ctx context.Context) ([]SpecExport, error) {
	u, err := url.Parse(p.endpoint + v1PathTasksExport)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadEndpointURL, err)
	}

	ctx, cancel := context.WithTimeout(ctx, httpV1Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreateRequest, err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrExportSpecs, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode)
	}

	var specs []SpecExport
	if err = json.NewDecoder(resp.Body).Decode(&specs); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecode, err)
	}
	return specs, nil
}
