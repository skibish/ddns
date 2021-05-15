package do

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/skibish/ddns/misc"
)

// Record describe record structure
type Record struct {
	ID       uint64 `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Data     string `json:"data"`
	TTL      uint64 `json:"ttl,omitempty"`
	Priority uint64 `json:"priority,omitempty"`
	Port     uint64 `json:"port,omitempty"`
	Weight   uint64 `json:"weight,omitempty"`
	Flags    uint64 `json:"flags,omitempty"`
	Tag      string `json:"tag,omitempty"`
}

type domainRecords struct {
	Records []Record `json:"domain_records"`
}

// DomainsService is an interface to interact with DNS records.
type DomainsService interface {
	List(context.Context, string) ([]Record, error)
	Create(context.Context, string, Record) error
	Update(context.Context, string, Record) error
}

// DigitalOcean hold
type DigitalOcean struct {
	c       *http.Client
	token   string
	url     string
	timeout time.Duration
}

// New return instance of DigitalOcean.
func New(token string, timeout time.Duration) *DigitalOcean {
	return &DigitalOcean{
		token:   token,
		c:       &http.Client{},
		url:     "https://api.digitalocean.com/v2",
		timeout: timeout,
	}
}

// List return domain DNS records.
func (d *DigitalOcean) List(ctx context.Context, domain string) ([]Record, error) {
	req, err := d.prepareRequest(http.MethodGet, fmt.Sprintf("/domains/%s/records", domain), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare a request: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	req = req.WithContext(ctx)

	res, err := d.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do a request: %w", err)
	}

	defer res.Body.Close()

	if !misc.Success(res.StatusCode) {
		return nil, fmt.Errorf("unexpected response with status code %d", res.StatusCode)
	}

	var records domainRecords
	if err := json.NewDecoder(res.Body).Decode(&records); err != nil {
		return nil, fmt.Errorf("failed to decode the response: %w", err)
	}

	return records.Records, nil
}

// Create creates DNS record.
func (d *DigitalOcean) Create(ctx context.Context, domain string, record Record) error {
	body, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal the record: %w", err)
	}

	req, err := d.prepareRequest(http.MethodPost, fmt.Sprintf("/domains/%s/records", domain), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to prepare a request: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	req = req.WithContext(ctx)

	res, err := d.c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do a request: %w", err)
	}

	defer res.Body.Close()

	if !misc.Success(res.StatusCode) {
		return fmt.Errorf("unexpected response with status code %d", res.StatusCode)
	}

	return nil
}

// Update updates DNS record.
func (d *DigitalOcean) Update(ctx context.Context, domain string, record Record) error {
	body, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal the record %w", err)
	}

	req, err := d.prepareRequest(http.MethodPut, fmt.Sprintf("/domains/%s/records/%d", domain, record.ID), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to prepare a request: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	req = req.WithContext(ctx)

	res, err := d.c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do a request: %w", err)
	}

	defer res.Body.Close()

	if !misc.Success(res.StatusCode) {
		return fmt.Errorf("unexpected response with status code %d", res.StatusCode)
	}

	return nil
}

func (d *DigitalOcean) prepareRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, d.url+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.token))

	return req, nil
}
