// htmltable enables structured data extraction from HTML tables and URLs
package htmltable

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// mock for tests
var htmlParse = html.Parse

var maxPossibleHeaderRows = 5

// Page is the container for all tables parseable
type Page struct {
	Tables []*Table

	ctx      context.Context
	rowSpans []int
	colSpans []int
	row      []string
	rows     [][]string
}

// New returns an instance of the page with possibly more than one table
func New(ctx context.Context, r io.Reader) (*Page, error) {
	p := &Page{ctx: ctx}
	return p, p.init(r)
}

// NewFromString is same as New(ctx.Context, io.Reader), but from string
func NewFromString(r string) (*Page, error) {
	return New(context.Background(), strings.NewReader(r))
}

// NewFromResponse is same as New(ctx.Context, io.Reader), but from http.Response.
//
// In case of failure, returns `ResponseError`, that could be further inspected.
func NewFromResponse(resp *http.Response) (*Page, error) {
	p, err := New(resp.Request.Context(), resp.Body)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// NewFromURL is same as New(ctx.Context, io.Reader), but from URL.
//
// In case of failure, returns `ResponseError`, that could be further inspected.
func NewFromURL(url string) (*Page, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return NewFromResponse(resp)
}

// Len returns number of tables found on the page
func (p *Page) Len() int {
	return len(p.Tables)
}

// FindWithColumns performs fuzzy matching of tables by given header column names
func (p *Page) FindWithColumns(columns ...string) (*Table, error) {
	// realistic p won't have this much
	found := 0xfffffff
	for idx, table := range p.Tables {
		matchedColumns := 0
		for _, col := range columns {
			for _, header := range table.Header {
				if col == header {
					// perform fuzzy matching of table headers
					matchedColumns++
				}
			}
		}
		if matchedColumns != len(columns) {
			continue
		}
		if found < len(p.Tables) {
			// and do a best-effort error message, that is cleaner than pandas.read_html
			return nil, fmt.Errorf("more than one table matches columns `%s`: "+
				"[%d] %s and [%d] %s", strings.Join(columns, ", "),
				found, p.Tables[found], idx, p.Tables[idx])
		}
		found = idx
	}
	if found > len(p.Tables) {
		return nil, fmt.Errorf("cannot find table with columns: %s",
			strings.Join(columns, ", "))
	}
	return p.Tables[found], nil
}

// Each row would call func with the value of the table cell from the column
// specified in the first argument.
//
// Returns an error if table has no matching column name.
func (p *Page) Each(a string, f func(a string) error) error {
	table, err := p.FindWithColumns(a)
	if err != nil {
		return err
	}
	offsets := map[string]int{}
	for idx, header := range table.Header {
		offsets[header] = idx
	}
	for idx, row := range table.Rows {
		if len(row) < 1 {
			continue
		}
		err = f(row[offsets[a]])
		if err != nil {
			return fmt.Errorf("row %d: %w", idx, err)
		}
	}
	return nil
}

// Each2 will get two columns specified in the first two arguments
// and call the func with those values for every row in the table.
//
// Returns an error if table has no matching column names.
func (p *Page) Each2(a, b string, f func(a, b string) error) error {
	table, err := p.FindWithColumns(a, b)
	if err != nil {
		return err
	}
	offsets := map[string]int{}
	for idx, header := range table.Header {
		offsets[header] = idx
	}
	_1, _2 := offsets[a], offsets[b]
	for idx, row := range table.Rows {
		if len(row) < 2 {
			continue
		}
		err = f(row[_1], row[_2])
		if err != nil {
			return fmt.Errorf("row %d: %w", idx, err)
		}
	}
	return nil
}

// Each3 will get three columns specified in the first three arguments
// and call the func with those values for every row in the table.
//
// Returns an error if table has no matching column names.
func (p *Page) Each3(a, b, c string, f func(a, b, c string) error) error {
	table, err := p.FindWithColumns(a, b, c)
	if err != nil {
		return err
	}
	offsets := map[string]int{}
	for idx, header := range table.Header {
		offsets[header] = idx
	}
	_1, _2, _3 := offsets[a], offsets[b], offsets[c]
	for idx, row := range table.Rows {
		if len(row) < 3 {
			continue
		}
		err = f(row[_1], row[_2], row[_3])
		if err != nil {
			return fmt.Errorf("row %d: %w", idx, err)
		}
	}
	return nil
}

func (p *Page) init(r io.Reader) error {
	root, err := htmlParse(r)
	if err != nil {
		return err
	}
	p.parse(root)
	p.finishTable()
	return nil
}

func (p *Page) parse(n *html.Node) {
	if n == nil {
		return
	}
	switch n.Data {
	case "td", "th":
		if len(p.rows) <= maxPossibleHeaderRows {
			offset := len(p.row)
			if len(p.colSpans) < offset+1 {
				p.colSpans = append(p.colSpans, 1)
				p.rowSpans = append(p.rowSpans, 1)
			}
			colSpan := p.intAttrOr(n, "colspan", 1)
			if colSpan > p.colSpans[offset] {
				p.colSpans[offset] = colSpan
			}
			rowSpan := p.intAttrOr(n, "rowspan", 1)
			if rowSpan > p.rowSpans[offset] {
				p.rowSpans[offset] = rowSpan
			}
		}
		var sb strings.Builder
		p.innerText(n, &sb)
		p.row = append(p.row, sb.String())
		return
	case "tr":
		p.finishRow()
	case "table":
		p.finishTable()
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.parse(c)
	}
}

func (p *Page) intAttrOr(n *html.Node, attr string, default_ int) int {
	for _, a := range n.Attr {
		if a.Key != attr {
			continue
		}
		val, err := strconv.Atoi(a.Val)
		if err != nil {
			return default_
		}
		return val
	}
	return default_
}

func (p *Page) finishRow() {
	if len(p.row) == 0 {
		return
	}
	p.rows = append(p.rows, p.row[:])
	p.row = []string{}
}

func (p *Page) finishTable() {
	p.finishRow()
	if len(p.rows) == 0 {
		return
	}
	maxRowSpan := 1
	for _, span := range p.rowSpans {
		if span > maxRowSpan {
			maxRowSpan = span
		}
	}
	dataOffset := 1
	header := p.rows[0]
	if maxRowSpan > 1 {
		// only supports 2 for now
		newHeader := []string{}
		si := 0
		for i, text := range p.rows[0] { // initial header
			if p.rowSpans[i] == 2 {
				newHeader = append(newHeader, text)
				continue
			}
			if p.colSpans[i] > 1 {
				ci := 0
				for ci < p.colSpans[i] {
					newHeader = append(newHeader, text+" "+p.rows[1][si+ci])
					ci++
				}
				// store last pos of col
				si = si + ci
				continue
			}
			newHeader = append(newHeader, text) // TODO: add coverage
		}
		header = newHeader
		dataOffset += 1
	}
	Logger(p.ctx, "found table", "columns", header, "count", len(p.rows))
	p.Tables = append(p.Tables, &Table{
		Header: header,
		Rows:   p.rows[dataOffset:],
	})
	p.rows = [][]string{}
	p.colSpans = []int{}
	p.rowSpans = []int{}
}

func (p *Page) innerText(n *html.Node, sb *strings.Builder) {
	if n.Type == html.TextNode {
		sb.WriteString(strings.TrimSpace(n.Data))
		return
	}
	if n.FirstChild == nil {
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.innerText(c, sb)
	}
}

// Table is the low-level representation of raw header and rows.
//
// Every cell string value is truncated of its whitespace.
type Table struct {
	// Header holds names of headers
	Header []string

	// Rows holds slice of string slices
	Rows [][]string
}

func (table *Table) String() string {
	return fmt.Sprintf("Table[%s] (%d rows)", strings.Join(table.Header, ", "), len(table.Rows))
}
