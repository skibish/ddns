package do

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/skibish/ddns/misc"
)

var name = "digitalocean"

var url = "https://api.digitalocean.com/v2"

// ErrorRequest is returned when some request failed
var ErrorRequest = errors.New("request Failed")

// Record describe record structure
type Record struct {
	ID   uint64 `json:"id"`
	Type string `yaml:"type" json:"type"`
	Name string `yaml:"name" json:"name"`
	Data string `json:"data"`
}

type domainRecords struct {
	Records []Record `json:"domain_records"`
}

type domainRecord struct {
	Record Record `json:"domain_record"`
}

// DigitalOceanInterface should be implemented by DigitalOcean
type DigitalOceanInterface interface {
	GetDomainRecords() ([]Record, error)
	CreateRecord(record Record) (*Record, error)
	UpdateRecord(record Record) (*Record, error)
}

// DigitalOcean is a main structure (implements DigitalOceanIterface)
type DigitalOcean struct {
	c      *http.Client
	token  string
	domain string
}

// New return instance of DigitalOcean
func New(domain, token string, c *http.Client) *DigitalOcean {
	return &DigitalOcean{
		domain: domain,
		token:  token,
		c:      c,
	}
}

// GetDomainRecords return domain records
func (d *DigitalOcean) GetDomainRecords() ([]Record, error) {
	req, errNR := d.prepareRequest("GET", fmt.Sprintf("%s/domains/%s/records", url, d.domain), nil)
	if errNR != nil {
		return nil, fmt.Errorf("%s: %s", name, errNR.Error())
	}

	res, errDo := d.c.Do(req)
	if errDo != nil {
		return nil, fmt.Errorf("%s: %s", name, errDo.Error())
	}

	defer res.Body.Close()

	if !misc.Success(res.StatusCode) {
		return nil, fmt.Errorf("%s: %s", name, ErrorRequest.Error())
	}

	var records domainRecords
	errDecode := json.NewDecoder(res.Body).Decode(&records)
	if errDecode != nil {
		return nil, fmt.Errorf("%s: %s", name, errDecode.Error())
	}

	return records.Records, nil
}

// CreateRecord create record
func (d *DigitalOcean) CreateRecord(record Record) (*Record, error) {
	body, errMarsh := json.Marshal(record)
	if errMarsh != nil {
		return nil, fmt.Errorf("%s: %s", name, errMarsh.Error())
	}

	req, errNR := d.prepareRequest("POST", fmt.Sprintf("%s/domains/%s/records", url, d.domain), bytes.NewBuffer(body))
	if errNR != nil {
		return nil, fmt.Errorf("%s: %s", name, errNR.Error())
	}

	res, errDo := d.c.Do(req)
	if errDo != nil {
		return nil, fmt.Errorf("%s: %s", name, errDo.Error())
	}

	defer res.Body.Close()

	if !misc.Success(res.StatusCode) {
		return nil, fmt.Errorf("%s: %s", name, ErrorRequest.Error())
	}

	var resRecord domainRecord
	errDecode := json.NewDecoder(res.Body).Decode(&resRecord)
	if errDecode != nil {
		return nil, fmt.Errorf("%s: %s", name, errDecode.Error())
	}

	return &resRecord.Record, nil
}

// UpdateRecord updates record
func (d *DigitalOcean) UpdateRecord(record Record) (*Record, error) {
	body, errMarsh := json.Marshal(record)
	if errMarsh != nil {
		return nil, fmt.Errorf("%s: %s", name, errMarsh.Error())
	}

	req, errNR := d.prepareRequest("PUT", fmt.Sprintf("%s/domains/%s/records/%d", url, d.domain, record.ID), bytes.NewBuffer(body))
	if errNR != nil {
		return nil, fmt.Errorf("%s: %s", name, errNR.Error())
	}

	res, errDo := d.c.Do(req)
	if errDo != nil {
		return nil, fmt.Errorf("%s: %s", name, errDo.Error())
	}

	defer res.Body.Close()

	if !misc.Success(res.StatusCode) {
		return nil, fmt.Errorf("%s: %s", name, ErrorRequest.Error())
	}

	var resRecord domainRecord
	errDecode := json.NewDecoder(res.Body).Decode(&resRecord)
	if errDecode != nil {
		return nil, fmt.Errorf("%s: %s", name, errDecode.Error())
	}

	return &resRecord.Record, nil
}

// prepareRequest bootstrap request with needed information
func (d *DigitalOcean) prepareRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, errNR := http.NewRequest(method, url, body)
	if errNR != nil {
		return nil, errNR
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.token))

	return req, nil
}
