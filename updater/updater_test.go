package updater

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"testing"

	"github.com/skibish/ddns/ipprovider"

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

type TestProvider struct {
	getIP func() (string, error)
}

func (t TestProvider) ForceIPV6() {
}

func (t TestProvider) GetIP() (string, error) {
	return t.getIP()
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

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
			{Type: "TXT", Name: "neo", Data: "{{.IP}} and text"},
		},
	}

	cf.Params = map[string]string{}

	u := &Updater{
		digitalOcean: doT,
		config:       cf,
		storage:      cf,
	}

	allRecords := []do.Record{
		{Type: "A", Name: "test"},
	}

	u.ip = "127.0.0.1"

	var errSync error
	errSync = u.syncRecords(allRecords)
	if errSync != nil {
		t.Error(errSync)
		return
	}

	if u.storage.Records[0].Data != "127.0.0.1" {
		t.Error("IPs should be the same", u.storage.Records[0].Data)
		return
	}
	if u.storage.Records[1].Data != "127.0.0.1 and text" {
		t.Error("IPs should be the same", u.storage.Records[1].Data)
		return
	}
}

func TestSyncRecordsCreateError(t *testing.T) {
	doT := struct{ TestDO }{}
	doT.createRecord = func(record do.Record) (*do.Record, error) {
		return nil, errors.New("Create error")
	}

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
		},
	}

	u := &Updater{
		digitalOcean: doT,
		config:       cf,
		storage:      cf,
	}

	allRecords := []do.Record{
		{Type: "A", Name: "test"},
	}

	u.ip = "127.0.0.1"

	var errSync error
	errSync = u.syncRecords(allRecords)
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

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
		},
	}

	u := &Updater{
		digitalOcean: doT,
		config:       cf,
		storage:      cf,
	}

	allRecords := []do.Record{
		{ID: 123, Type: "A", Name: "test"},
	}

	u.ip = "127.0.0.1"

	var errSync error
	errSync = u.syncRecords(allRecords)
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

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
		},
	}

	u := &Updater{
		digitalOcean: doT,
		config:       cf,
		storage:      cf,
	}

	allRecords := []do.Record{
		{ID: 123, Type: "A", Name: "test"},
	}

	u.ip = "127.0.0.1"

	var errSync error
	errSync = u.syncRecords(allRecords)
	if errSync == nil {
		t.Error("Should be error, but everything is OK.")
		return
	}
}

func TestCheckAndUpdateOnlyCheck(t *testing.T) {
	provT := struct{ TestProvider }{}
	provT.getIP = func() (string, error) {
		return "127.0.0.3", nil
	}

	p := ipprovider.New()
	p.Register(provT)

	doT := struct{ TestDO }{}
	doT.updateRecord = func(record do.Record) (*do.Record, error) {
		return &do.Record{
			ID:   124,
			Type: "TXT",
			Name: "neo",
			Data: "127.0.0.1 and text",
		}, nil
	}

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
		},
	}

	u := &Updater{
		ipprovider:   p,
		digitalOcean: doT,
		config:       cf,
		storage:      cf,
	}

	u.ip = "127.0.0.1"

	var errCheck error
	errCheck = u.checkAndUpdate()
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

	provT := struct{ TestProvider }{}
	provT.getIP = func() (string, error) {
		return "127.0.0.3", nil
	}

	p := ipprovider.New()
	p.Register(provT)

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
			{Type: "TXT", Name: "test", Data: "{{.IP}} is {{.foo}}"},
		},
	}

	cf.Params = map[string]string{}
	cf.Params["foo"] = "bar"

	u := &Updater{
		ipprovider:   p,
		digitalOcean: doT,
		config:       cf,
		storage:      cf,
		ip:           "127.0.0.1",
	}

	var errUpdate error
	errUpdate = u.checkAndUpdate()
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

	provT := struct{ TestProvider }{}
	provT.getIP = func() (string, error) {
		return "127.0.0.3", nil
	}

	p := ipprovider.New()
	p.Register(provT)

	cf := &conf.Configuration{
		Records: []do.Record{
			{Type: "A", Name: "test"},
		},
	}

	u := &Updater{
		ipprovider:   p,
		digitalOcean: doT,
		config:       cf,
		storage:      cf,
		ip:           "127.0.0.1",
	}

	var errUpdate error
	errUpdate = u.checkAndUpdate()
	if errUpdate == nil {
		t.Error("Should be error, but everything is OK")
		return
	}
}
