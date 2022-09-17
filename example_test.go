package htmltable_test

import (
	"context"
	"fmt"

	"github.com/nfx/go-htmltable"
)

func ExampleNewSliceFromUrl() {
	type Ticker struct {
		Symbol   string `header:"Symbol"`
		Security string `header:"Security"`
		CIK      string `header:"CIK"`
	}
	url := "https://en.wikipedia.org/wiki/List_of_S%26P_500_companies"
	out, _ := htmltable.NewSliceFromURL[Ticker](url)
	fmt.Println(out[0].Symbol)
	fmt.Println(out[0].Security)

	// Output: 
	// MMM
	// 3M
}

func ExampleNewFromString() {
	page, _ := htmltable.NewFromString(`<body>
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
	</body>`)

	fmt.Printf("found %d tables\n", page.Len())
	_ = page.Each2("c", "d", func(c, d string) error {
		fmt.Printf("c:%s d:%s\n", c, d)
		return nil
	})

	// Output: 
	// found 2 tables
	// c:2 d:5
	// c:4 d:6
}

func ExampleNewFromURL() {
	page, _ := htmltable.NewFromURL("https://en.wikipedia.org/wiki/List_of_S%26P_500_companies")
	_, err := page.FindWithColumns("invalid", "column", "names")
	fmt.Println(err)

	// Output: 
	// cannot find table with columns: invalid, column, names
}

func ExampleLogger() {
	htmltable.Logger = func(_ context.Context, msg string, fields ...any) {
		fmt.Printf("[INFO] %s %v\n", msg, fields)
	}
	_, _ = htmltable.NewFromURL("https://en.wikipedia.org/wiki/List_of_S%26P_500_companies")

	// Output:
	// [INFO] found table [columns [Symbol Security SEC filings GICSSector GICS Sub-Industry Headquarters Location Date first added CIK Founded] count 504]
	// [INFO] found table [columns [Date Added Ticker Added Security Removed Ticker Removed Security Reason] count 308]
}