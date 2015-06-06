package webcounter

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var server *httptest.Server

func init() {
	wc, err := New()
	if err != nil {
		log.Fatalf("cannot create new instance %s\n", err)
	}
	server = httptest.NewServer(wc)
}

func expectCount(url string, val int) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making request to %s - %s\n", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200 for GET to %s received %d\n", url, resp.StatusCode)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body for %s - %s\n", url, err)
	}
	stringVal := strconv.Itoa(val)
	if string(content) != stringVal {
		return fmt.Errorf("expected content %d but received %s\n", val, string(content))
	}
	return nil
}

func TestWebcounterValues(t *testing.T) {
	id := "aa.txt"
	url := server.URL + "/" + id
	for i := 0; i < 3; i++ {
		if err := expectCount(url, i); err != nil {
			t.Error(err)
		} else {
			t.Logf("value %d found as expected\n", i)
		}
	}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("could not make new DELETE request %s\n", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("could not DELETE %s - %s\n", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE expected status 200 but received - %d\n", resp.StatusCode)
	}
	t.Logf("reset counter for %s\n", id)
	for i := 0; i < 11; i++ {
		if err := expectCount(url, i); err != nil {
			t.Error(err)
		} else {
			t.Logf("value %d found as expected\n", i)
		}
	}
}

// TestSuffixes checks that the content-type returned are expected given a request suffix
func TestSuffixes(t *testing.T) {
	tests := []struct {
		suffix      string
		contentType string
	}{
		{".txt", "text/plain"},
		{".png", "image/png"},
		{".jpg", "image/jpeg"},
		{".jpeg", "image/jpeg"},
		{".gif", "image/gif"},
		{"", "image/png"},
	}
	baseURL := server.URL + "/aa"
	for _, test := range tests {
		url := baseURL + test.suffix
		resp, err := http.Get(url)
		if err != nil {
			t.Errorf("error making request to %s - %s\n", url, err)
		}
		if resp.Header.Get("Content-type") != test.contentType {
			t.Errorf("expected Content-type %s for %s but received %s\n", test.contentType, url, resp.Header.Get("Content-type"))
		} else {
			t.Logf("content-type ok for suffix %q\n", test.suffix)
		}
	}
}
