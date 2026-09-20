package htmltable_test

import (
	"fmt"

	"github.com/nfx/go-htmltable"
)

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
