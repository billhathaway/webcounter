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

// init sets up our handler on a test instance
func init() {
	wc, err := New()
	if err != nil {
		log.Fatalf("cannot create new instance %s\n", err)
	}
	server = httptest.NewServer(wc)
}

// TestInvalidRequests will verify proper responses to requests we don't support
// including invalid URLs and request methods
func TestInvalidRequests(t *testing.T) {
	urls := []string{"/", "/favicon.ico"}
	for _, url := range urls {
		resp, err := http.Get(server.URL + url)
		if err != nil {
			t.Errorf("error hitting %s - %s\n", url, err)
			continue
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected %d for %s but received %d\n", http.StatusNotFound, url, resp.StatusCode)
			continue
		}
		t.Logf("received %d as expected for %s\n", http.StatusNotFound, url)
		resp.Body.Close()
	}
	url := server.URL + "/a.txt"
	resp, err := http.Post(url, "plain/text", nil)
	if err != nil {
		t.Errorf("error hitting %s - %s\n", url, err)
		return
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected %d for POST but received %d\n", http.StatusBadRequest, resp.StatusCode)
	}
	resp.Body.Close()

}

// TestWebcounterValues validates the basic ounter works including resets
func TestWebcounterValues(t *testing.T) {
	id := "aa.txt"
	url := server.URL + "/" + id
	for i := 0; i < 3; i++ {
		if err := expectCount(url, "", i+1); err != nil {
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
	for i := 0; i < 3; i++ {
		if err := expectCount(url, "", i+1); err != nil {
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

// TestReferers validates that counts are kept separate per HTTP referer
func TestReferers(t *testing.T) {
	tests := []struct {
		referer string
		count   int
	}{
		{"", 1},
		{"", 2},
		{"google.com", 1},
		{"google.com", 2},
		{"", 3},
	}
	url := server.URL + "/referer.txt"
	for _, test := range tests {
		err := expectCount(url, test.referer, test.count)
		if err != nil {
			t.Error(err)
		}
	}
}

// expectCount is a helper function to see if the number returned is correct
func expectCount(url string, referer string, val int) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", referer)
	resp, err := http.DefaultClient.Do(req)
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
