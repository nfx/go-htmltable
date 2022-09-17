package htmltable

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

const fixture = `<body>
<h1>foo</h2>
<table>
	<tr><td>a</td><td>b</td></tr>
	<tr><td> 1 </td><td>2</td></tr>
	<tr><td>3  </td><td>4   </td></tr>
</table>
<h1>bar</h2>
<table>
	<tr><th>b</th><th>c</th><th>d</th></tr>
	<tr><td>1</td><td>2</td><td>5</td></tr>
	<tr><td>3</td><td>4</td><td>6</td></tr>
</table>
</body>`

func TestFindsAllTables(t *testing.T) {
	p, err := NewFromString(fixture)
	assertNoError(t, err)
	assertEqual(t, p.Len(), 2)
}

// added public domain data from https://en.wikipedia.org/wiki/List_of_S&P_500_companies
const fixtureColspans = `<table>
	<thead>
		<tr>
			<th rowspan="2">Date</th>
			<th colspan="2">Added</th>
			<th colspan="2">Removed</th>
			<th rowspan="2">Reason</th>
		</tr>
		<tr>
			<th rowspan="@#$%^&">Ticker</th>
			<th>Security</th>
			<th>Ticker</th>
			<th>Security</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td>June 21, 2022</td>
			<td>KDP</td>
			<td><a href="/wiki/Keurig_Dr_Pepper" title="Keurig Dr Pepper">Keurig Dr Pepper</a></td>
			<td>UA/UAA</td>
			<td><a href="/wiki/Under_Armour" title="Under Armour">Under Armour</a></td>
			<td>Market capitalization change.<sup id="cite_ref-sp20220603_4-0" class="reference"><a href="#cite_note-sp20220603-4">[4]</a></sup></td>
		</tr>
		<tr>
			<td>June 21, 2022</td>
			<td>ON</td>
			<td><a href="/wiki/ON_Semiconductor" class="mw-redirect" title="ON Semiconductor">ON Semiconductor</a></td>
			<td>IPGP</td>
			<td><a href="/wiki/IPG_Photonics" title="IPG Photonics">IPG Photonics</a></td>
			<td>Market capitalization change.<sup id="cite_ref-sp20220603_4-1" class="reference"><a href="#cite_note-sp20220603-4">[4]</a></sup></td>
		</tr>
	</tbody>
</table>`

func TestFindsWithColspans(t *testing.T) {
	p, err := NewFromString(fixtureColspans)
	assertNoError(t, err)
	assertEqual(t, p.Len(), 1)
	assertEqual(t, "Added Ticker", p.tables[0].header[1])
	assertEqual(t, "Market capitalization change.[4]", p.tables[0].rows[0][5])
}

func TestInitFails(t *testing.T) {
	prev := htmlParse
	t.Cleanup(func() {
		htmlParse = prev
	})
	htmlParse = func(r io.Reader) (*html.Node, error) {
		return nil, fmt.Errorf("nope")
	}
	_, err := New(context.Background(), strings.NewReader(".."))

	assertEqualError(t, err, "nope")
}

func TestNewFromHttpResponseError(t *testing.T) {
	prev := htmlParse
	t.Cleanup(func() {
		htmlParse = prev
	})
	htmlParse = func(r io.Reader) (*html.Node, error) {
		return nil, fmt.Errorf("nope")
	}
	_, err := NewFromResponse(&http.Response{
		Request: &http.Request{},
	})
	assertEqualError(t, err, "nope")
}

func TestRealPageFound(t *testing.T) {
	wiki, err := http.Get("https://en.wikipedia.org/wiki/List_of_S%26P_500_companies")
	assertNoError(t, err)
	p, err := NewFromResponse(wiki)
	assertNoError(t, err)
	snp, err := p.FindWithColumns("Symbol", "Security", "CIK")
	assertNoError(t, err)
	assertGreaterOrEqual(t, len(snp.rows), 500)
}

func TestRealPageFound_BasicRowColSpans(t *testing.T) {
	wiki, err := http.Get("https://en.wikipedia.org/wiki/List_of_S%26P_500_companies")
	assertNoError(t, err)
	p, err := NewFromResponse(wiki)
	assertNoError(t, err)
	snp, err := p.FindWithColumns("Date", "Added Ticker", "Removed Ticker")
	assertNoError(t, err)
	assertGreaterOrEqual(t, len(snp.rows), 250)
}

func TestFindsTableByColumnNames(t *testing.T) {
	p, err := NewFromString(fixture)
	assertNoError(t, err)

	cd, err := p.FindWithColumns("c", "d")
	assertNoError(t, err)
	assertEqual(t, 2, len(cd.rows))
}

func TestEach(t *testing.T) {
	p, err := NewFromString(fixture)
	assertNoError(t, err)
	err = p.Each("a", func(a string) error {
		t.Logf("%s", a)
		return nil
	})
	assertNoError(t, err)
}

func TestEachFails(t *testing.T) {
	p, err := NewFromString(fixture)
	assertNoError(t, err)
	err = p.Each("a", func(a string) error {
		return fmt.Errorf("nope")
	})
	assertEqualError(t, err, "row 0: nope")
}

func TestEachFailsNoCols(t *testing.T) {
	p, err := NewFromString(fixture)
	assertNoError(t, err)
	err = p.Each("x", func(a string) error {
		return nil
	})
	assertEqualError(t, err, "cannot find table with columns: x")
}

func TestEach2(t *testing.T) {
	p, err := NewFromString(fixture)
	assertNoError(t, err)
	err = p.Each2("b", "c", func(b, c string) error {
		t.Logf("%s %s", b, c)
		return nil
	})
	assertNoError(t, err)
}

func TestEach2Fails(t *testing.T) {
	p, err := NewFromString(fixture)
	assertNoError(t, err)
	err = p.Each2("b", "c", func(b, c string) error {
		return fmt.Errorf("nope")
	})
	assertEqualError(t, err, "row 0: nope")
}

func TestEach2FailsNoCols(t *testing.T) {
	p, err := NewFromString(fixture)
	assertNoError(t, err)
	err = p.Each2("x", "y", func(b, c string) error {
		return nil
	})
	assertEqualError(t, err, "cannot find table with columns: x, y")
}

func TestEach3(t *testing.T) {
	p, err := NewFromString(fixture)
	assertNoError(t, err)
	err = p.Each3("b", "c", "d", func(b, c, d string) error {
		t.Logf("%s %s %s", b, c, d)
		return nil
	})
	assertNoError(t, err)
}

func TestEach3Fails(t *testing.T) {
	p, err := NewFromString(fixture)
	assertNoError(t, err)
	err = p.Each3("b", "c", "d", func(b, c, d string) error {
		return fmt.Errorf("nope")
	})
	assertEqualError(t, err, "row 0: nope")
}

func TestEach3FailsNoCols(t *testing.T) {
	p, err := NewFromString(fixture)
	assertNoError(t, err)
	err = p.Each3("x", "y", "z", func(b, c, d string) error {
		return nil
	})
	assertEqualError(t, err, "cannot find table with columns: x, y, z")
}

func TestMoreThanOneTableFoundErrors(t *testing.T) {
	p, err := NewFromString(fixture)
	assertNoError(t, err)

	_, err = p.FindWithColumns("b")
	assertError(t, err)
}

func TestNoTablesFoundErrors(t *testing.T) {
	p, err := NewFromString(fixture)
	assertNoError(t, err)

	_, err = p.FindWithColumns("z")
	assertError(t, err)
}

func TestNilNodeReturns(t *testing.T) {
	p := &page{}
	p.parse(nil)
}
