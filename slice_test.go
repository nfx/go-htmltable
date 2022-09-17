package htmltable

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type nice struct {
	C string `header:"c"`
	D string `header:"d"`
}

func TestNewSliceFromString(t *testing.T) {
	out, err := NewSliceFromString[nice](fixture)
	assertNoError(t, err)
	assertEqual(t, []nice{
		{"2", "5"},
		{"4", "6"},
	}, out)
}

type Ticker struct {
	Symbol   string `header:"Symbol"`
	Security string `header:"Security"`
	CIK      string `header:"CIK"`
}

func TestNewSliceFromUrl(t *testing.T) {
	url := "https://en.wikipedia.org/wiki/List_of_S%26P_500_companies"
	out, err := NewSliceFromURL[Ticker](url)
	assertNoError(t, err)
	assertGreaterOrEqual(t, len(out), 500)
}

func TestNewSliceFromUrl_Fails(t *testing.T) {
	_, err := NewSliceFromURL[Ticker]("https://127.0.0.1")
	assertEqualError(t, err, "Get \"https://127.0.0.1\": dial tcp 127.0.0.1:443: connect: connection refused")
}

func TestNewSliceFromUrl_NoTables(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()
	_, err := NewSliceFromURL[Ticker](server.URL)
	assertEqualError(t, err, "cannot find table with columns: Symbol, Security, CIK")
}

func TestNewSliceInvalidTypes(t *testing.T) {
	type exotic struct {
		A string  `header:""`
		C float32 `header:"c"`
	}
	_, err := NewSliceFromString[exotic](fixture)
	assertEqualError(t, err, "only strings are supported, C is float32")
}
