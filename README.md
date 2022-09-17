# HTML table data extractor for Go

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg)](https://pkg.go.dev/mod/github.com/nfx/go-htmltable)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/nfx/go-htmltable/blob/main/LICENSE)
[![codecov](https://codecov.io/gh/nfx/go-htmltable/branch/main/graph/badge.svg)](https://codecov.io/gh/nfx/go-htmltable)
[![build](https://github.com/nfx/go-htmltable/workflows/build/badge.svg?branch=main)](https://github.com/nfx/go-htmltable/actions?query=workflow%3Abuild+branch%3Amain)


`htmltable` enables structured data extraction from HTML tables and URLs and requires almost no external dependencies.

## Usage

You can retrieve a slice of `header`-annotated types using the `NewSlice*` contructors:

```go
import "github.com/nfx/go-htmltable"

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
```

An error would be thrown if there's no matching page with the specified columns:

```go
page, _ := htmltable.NewFromURL("https://en.wikipedia.org/wiki/List_of_S%26P_500_companies")
_, err := page.FindWithColumns("invalid", "column", "names")
fmt.Println(err)

// Output: 
// cannot find table with columns: invalid, column, names
```

And you can use more low-level API to work with extracted data:

```go
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
```

And the last note: you're encouraged to plug your own structured logger:

```go
htmltable.Logger = func(_ context.Context, msg string, fields ...any) {
    fmt.Printf("[INFO] %s %v\n", msg, fields)
}
htmltable.NewFromURL("https://en.wikipedia.org/wiki/List_of_S%26P_500_companies")

// Output:
// [INFO] found table [columns [Symbol Security SEC filings GICSSector GICS Sub-Industry Headquarters Location Date first added CIK Founded] count 504]
// [INFO] found table [columns [Date Added Ticker Added Security Removed Ticker Removed Security Reason] count 308]
```

## Inspiration

This library aims to be something like [pandas.read_html](https://pandas.pydata.org/docs/reference/api/pandas.read_html.html) or [table_extract](https://docs.rs/table-extract/latest/table_extract/) Rust crate, but more idiomatic for Go.