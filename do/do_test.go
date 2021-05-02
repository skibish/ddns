package do

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDomainRecordsSuccess(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") != "Bearer amazingtoken" {
			t.Error("Not correct Authorization value")
			return
		}

		if r.Method != "GET" {
			t.Errorf("Method should be GET, got %s instead", r.Method)
			return
		}

		if r.URL.Path != "/domains/example.com/records" {
			t.Errorf("Expected path /domains/example.com/records, got %s", r.URL.Path)
			return
		}

		w.Write([]byte(`{
  "domain_records": [
    {
      "id": 3352895,
      "type": "A",
      "name": "@",
      "data": "1.2.3.4",
      "priority": null,
      "port": null,
      "weight": null
    }
  ]
}`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url = server.URL

	d := New("example.com", "amazingtoken", &http.Client{})

	recs, errGet := d.GetDomainRecords()
	if errGet != nil {
		t.Errorf("Got error : %s", errGet.Error())
		return
	}

	if recs[0].ID != 3352895 {
		t.Error("Got not correct response")
		return
	}
}

func TestGetDomainRecordsIncorrectStatusCode(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte(``))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url = server.URL

	d := New("example.com", "amazingtoken", &http.Client{})

	_, errGet := d.GetDomainRecords()
	if errGet == nil {
		t.Error("Should be error, but everything is OK")
		return
	}
}

func TestGetDomainRecordsParseError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`aaa`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url = server.URL

	d := New("example.com", "amazingtoken", &http.Client{})

	_, errGet := d.GetDomainRecords()
	if errGet.Error() != "digitalocean: invalid character 'a' looking for beginning of value" {
		t.Error("Go not expected value: ", errGet.Error())
		return
	}
}

func TestCreateRecordSuccess(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") != "Bearer amazingtoken" {
			t.Error("Not correct Authorization value")
			return
		}

		if r.Method != "POST" {
			t.Errorf("Method should be POST, got %s instead", r.Method)
			return
		}

		if r.URL.Path != "/domains/example.com/records" {
			t.Errorf("Expected path /domains/example.com/records, got %s", r.URL.Path)
			return
		}

		w.Write([]byte(`{
  "domain_record": {
      "id": 3352895,
      "type": "A",
      "name": "@",
      "data": "1.2.3.4",
      "priority": null,
      "port": null,
      "weight": null
    }
}`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url = server.URL

	d := New("example.com", "amazingtoken", &http.Client{})

	recs, errGet := d.CreateRecord(Record{})
	if errGet != nil {
		t.Errorf("Got error : %s", errGet.Error())
		return
	}

	if recs.ID != 3352895 {
		t.Error("Got not correct response")
		return
	}
}

func TestCreateRecordIncorrectStatusCode(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte(``))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url = server.URL

	d := New("example.com", "amazingtoken", &http.Client{})

	_, errGet := d.CreateRecord(Record{})
	if errGet == nil {
		t.Error("Should be error, but everything is OK")
		return
	}
}

func TestCreateRecordParseError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`aaa`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url = server.URL

	d := New("example.com", "amazingtoken", &http.Client{})

	_, errGet := d.CreateRecord(Record{})
	if errGet.Error() != "digitalocean: invalid character 'a' looking for beginning of value" {
		t.Error("Go not expected value: ", errGet.Error())
		return
	}
}

func TestUpdateRecordSuccess(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") != "Bearer amazingtoken" {
			t.Error("Not correct Authorization value")
			return
		}

		if r.Method != "PUT" {
			t.Errorf("Method should be PUT, got %s instead", r.Method)
			return
		}

		if r.URL.Path != "/domains/example.com/records/0" {
			t.Errorf("Expected path /domains/example.com/records/0, got %s", r.URL.Path)
			return
		}

		w.Write([]byte(`{
  "domain_record": {
      "id": 3352895,
      "type": "A",
      "name": "@",
      "data": "1.2.3.4",
      "priority": null,
      "port": null,
      "weight": null
    }
}`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url = server.URL

	d := New("example.com", "amazingtoken", &http.Client{})

	recs, errGet := d.UpdateRecord(Record{})
	if errGet != nil {
		t.Errorf("Got error : %s", errGet.Error())
		return
	}

	if recs.ID != 3352895 {
		t.Error("Got not correct response")
		return
	}
}

func TestUpdateRecordIncorrectStatusCode(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte(``))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url = server.URL

	d := New("example.com", "amazingtoken", &http.Client{})

	_, errGet := d.UpdateRecord(Record{})
	if errGet == nil {
		t.Error("Should be error, but everything is OK")
		return
	}
}

func TestUpdateRecordParseError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`aaa`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url = server.URL

	d := New("example.com", "amazingtoken", &http.Client{})

	_, errGet := d.UpdateRecord(Record{})
	if errGet.Error() != "digitalocean: invalid character 'a' looking for beginning of value" {
		t.Error("Go not expected value: ", errGet.Error())
		return
	}
}
