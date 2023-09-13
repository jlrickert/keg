package markdown_test

import (
	"fmt"
	"github.com/rwxrob/keg/markdown"

	"github.com/rwxrob/pegn/scanner"
)

func ExampleScanTitle_short() {

	s := scanner.New(`# A short title`)

	fmt.Println(markdown.ScanTitle(s, nil))
	s.Print()

	// Output:
	// true
	// 'e' 14-15 ""

}

func ExampleTitle_parsed_Short() {

	s := scanner.New(`# A short title`)

	title := make([]rune, 0, 100)
	fmt.Println(markdown.ScanTitle(s, &title))

	s.Print()
	fmt.Printf("%q", string(title))

	// Output:
	// true
	// 'e' 14-15 ""
	// "A short title"

}

func ExampleTitle() {
	title, _ := markdown.ReadTitle(`testdata/sample-node`)
	fmt.Println(title)
	// Output:
	// This is a title
}

func ExampleTitle_no_README() {
	title, _ := markdown.ReadTitle(`testdata/sample-node`)
	fmt.Println(title)
	// Output:
	// This is a title
}

func ExampleTitleWithFrontMater() {
	title, _ := markdown.ReadTitle(`testdata/front-mater-node`)
	fmt.Println(title)
	// Output:
	// This is the title
}

func ExampleTitleWithTitleInFrontMater() {
	title, _ := markdown.ReadTitle(`testdata/front-mater-node-with-title`)
	fmt.Println(title)
	// Output:
	// Example
}

/*
func ExampleTitle_long() {

	s := scanner.New(`# A really, really long title that is more than 72 runes long but doesn't get truncated`)
	fmt.Println(markdown.Title.Scan(s, nil))
	s.Print()

	// Output:
	// false
	// 't' 72-73 " get truncated"

}

func ExampleTitle_parsed_Long() {

	s := scanner.New(`# A really, really long title that is more than 72 runes long but doesn't get truncated`)
	title := make([]rune, 0, 70)
	fmt.Println(markdown.Title.Scan(s, &title))
	s.Print()
	fmt.Printf("%q", string(title))

	// Output:
	// false
	// 't' 72-73 " get truncated"
	// "A really, really long title that is more than 72 runes long but doesn'"

}
*/
