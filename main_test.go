package main

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"testing"

	"github.com/skibish/ddns/conf"
	"github.com/skibish/ddns/do"
)

type TestDO struct {
	getDomainRecords func() ([]do.Record, error)
	createRecord     func(record do.Record) (*do.Record, error)
	updateRecord     func(record do.Record) (*do.Record, error)
}

func (t TestDO) GetDomainRecords() ([]do.Record, error) {
	return t.getDomainRecords()
}
func (t TestDO) CreateRecord(record do.Record) (*do.Record, error) {
	return t.createRecord(record)
}
func (t TestDO) UpdateRecord(record do.Record) (*do.Record, error) {
	return t.updateRecord(record)
}

func TestSyncRecordsCreateNew(t *testing.T) {
	doT := struct{ TestDO }{}
	doT.createRecord = func(record do.Record) (*do.Record, error) {
		if record.Type == "A" {
			return &do.Record{
				ID:   123,
				Type: "A",
				Name: "test",
				Data: "127.0.0.1",
			}, nil
		}

		return &do.Record{
			ID:   124,
			Type: "TXT",
			Name: "neo",
			Data: "127.0.0.1 and text",
		}, nil
	}

	doT.updateRecord = func(record do.Record) (*do.Record, error) {
		return &do.Record{
			ID:   124,
			Type: "TXT",
			Name: "neo",
			Data: "127.0.0.1 and text",
		}, nil
	}

	digio = doT

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
			{Type: "TXT", Name: "neo", Data: "{{.IP}} and text"},
		},
	}

	cf.Params = map[string]string{}

	storage := &conf.Configuration{}
	*storage = *cf

	allRecords := []do.Record{
		{Type: "A", Name: "test"},
	}

	currentIP = "127.0.0.1"

	var errSync error
	errSync = syncRecords(storage, cf, allRecords)
	if errSync != nil {
		t.Error(errSync)
		return
	}

	if storage.Records[0].Data != "127.0.0.1" {
		t.Error("IPs should be the same", storage.Records[0].Data)
		return
	}
	if storage.Records[1].Data != "127.0.0.1 and text" {
		t.Error("IPs should be the same", storage.Records[1].Data)
		return
	}
}

func TestSyncRecordsCreateError(t *testing.T) {
	doT := struct{ TestDO }{}
	doT.createRecord = func(record do.Record) (*do.Record, error) {
		return nil, errors.New("Create error")
	}

	digio = doT

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
		},
	}

	allRecords := []do.Record{
		{Type: "A", Name: "test"},
	}

	currentIP = "127.0.0.1"

	var errSync error
	errSync = syncRecords(cf, cf, allRecords)
	if errSync == nil {
		t.Error("Should be error, but everything is OK.")
		return
	}
}

func TestSyncRecordsUpdateRecord(t *testing.T) {
	doT := struct{ TestDO }{}
	doT.updateRecord = func(record do.Record) (*do.Record, error) {
		return &do.Record{
			ID:   123,
			Type: "A",
			Name: "test",
			Data: "127.0.0.1",
		}, nil
	}

	digio = doT

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
		},
	}

	allRecords := []do.Record{
		{ID: 123, Type: "A", Name: "test"},
	}

	currentIP = "127.0.0.1"

	var errSync error
	errSync = syncRecords(cf, cf, allRecords)
	if errSync != nil {
		t.Error(errSync)
		return
	}

	if cf.Records[0].Data != "127.0.0.1" {
		t.Error("IPs should be the same", cf.Records[0].Data)
		return
	}
}

func TestSyncRecordsUpdateError(t *testing.T) {
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	log.SetOutput(bw)
	defer bw.Flush()

	doT := struct{ TestDO }{}
	doT.updateRecord = func(record do.Record) (*do.Record, error) {
		return nil, errors.New("Update error")
	}

	digio = doT

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
		},
	}

	allRecords := []do.Record{
		{ID: 123, Type: "A", Name: "test"},
	}

	currentIP = "127.0.0.1"

	var errSync error
	errSync = syncRecords(cf, cf, allRecords)
	if errSync == nil {
		t.Error("Should be error, but everything is OK.")
		return
	}
}

func TestCheckAndUpdateOnlyCheck(t *testing.T) {
	currentIP = "127.0.0.1"

	tf := func() string {
		return "127.0.0.1"
	}

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
		},
	}

	var errCheck error
	errCheck = checkAndUpdate(cf, cf, tf)
	if errCheck != nil {
		t.Error(errCheck)
		return
	}
}

func TestCheckAndUpdateOnlyUpdate(t *testing.T) {
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	log.SetOutput(bw)
	defer bw.Flush()

	doT := struct{ TestDO }{}
	doT.updateRecord = func(record do.Record) (*do.Record, error) {
		return &do.Record{
			ID:   123,
			Type: "A",
			Name: "test",
			Data: "127.0.0.1",
		}, nil
	}

	digio = doT
	currentIP = "127.0.0.1"

	tf := func() string {
		return "127.0.0.3"
	}

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
			{Type: "TXT", Name: "test", Data: "{{.IP}} is {{.foo}}"},
		},
	}

	cf.Params = map[string]string{}
	cf.Params["foo"] = "bar"

	storage = &conf.Configuration{}
	*storage = *cf

	var errUpdate error
	errUpdate = checkAndUpdate(storage, cf, tf)
	if errUpdate != nil {
		t.Error(errUpdate)
		return
	}
}

func TestCheckAndUpdateError(t *testing.T) {
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	log.SetOutput(bw)
	defer bw.Flush()

	doT := struct{ TestDO }{}
	doT.updateRecord = func(record do.Record) (*do.Record, error) {
		return nil, errors.New("Update Error")
	}

	digio = doT
	currentIP = "127.0.0.1"

	tf := func() string {
		return "127.0.0.3"
	}

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
		},
	}

	var errUpdate error
	errUpdate = checkAndUpdate(cf, cf, tf)
	if errUpdate == nil {
		t.Error("Should be error, but everything is OK")
		return
	}
}
