package keg_test

import (
	"fmt"
	"path/filepath"

	"github.com/rwxrob/keg"
)

func ExampleNodePaths() {

	dirs, low, high := keg.NodePaths("testdata/samplekeg")

	fmt.Println(low)
	fmt.Println(high)

	for _, d := range dirs {
		fmt.Printf("%v ", filepath.Base(d.Path))
	}

	// Output:
	// 0
	// 12
	// 0 1 10 11 12 2 3 4 5 6 7 8 9

}

func ExampleUpdatedString() {
	fmt.Println(keg.UpdatedString(`testdata/samplekeg`))
	// Output:
	// 2022-11-26 19:33:24Z
}

/*
func ExampleDex_WithTitleText() {
	dex, _ := keg.ReadDex(`testdata/samplekeg`)
	fmt.Println(dex.WithTitleText(`5`).TSV())
	// Output:
	// ignored
}
*/

/*
func ExampleReadDex() {
	dex, _ := keg.ReadDex(`testdata/samplekeg`)
	fmt.Println(dex.TSV())
	// Output:
	// ignored
}
*/

/*
func ExampleUpdateUpdated() {
	err := keg.UpdateUpdated(`testdata/samplekeg`)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// ignored
}
*/

/*
func ExampleScanDex() {
	fmt.Println(keg.ScanDex("testdata/samplekeg"))

	// Output:
	// ignored
}
*/

/*
func ExampleMakeDex() {
	fmt.Println(keg.MakeDex("testdata/samplekeg"))

	// Output:
	// ignored
}
*/

/*
func ExampleMakeNode() {
	fmt.Println(keg.MakeNode("testdata/samplekeg"))

	// Output:
	// ignored
}
*/

/*
func ExampleWriteDex() {
	dex, _ := keg.ReadDex(`testdata/samplekeg`)
	fmt.Println(keg.WriteDex("testdata/newkeg", dex))

	// Output:
	// ignored
}
*/

func ExampleNextNew() {
	fmt.Println(keg.NextNew(`testdata/samplekeg`).N)
	fmt.Println(keg.NextNew(`testdata/noexist`))
	// Output:
	// 13
	// <nil>
}

func ExampleLast() {
	fmt.Println(keg.Last(`testdata/samplekeg`).N)
	fmt.Println(keg.Last(`testdata/noexist`))
	// Output:
	// 12
	// <nil>
}

func ExampleFirst() {
	fmt.Println(keg.First(`testdata/samplekeg`).N)
	fmt.Println(keg.First(`testdata/noexist`))
	// Output:
	// 1
	// <nil>
}

func ExampleReadTags() {
	tags, err := keg.ReadTags(`testdata/samplekeg`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tags[`foo`][1])
	// Output:
	// 6
}

func ExampleGrepTags() {
	str, err := keg.GrepTags(`testdata/samplekeg`, `foo`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
	// Output:
	// foo 2 6 3
}

func ExampleGrepTags_multiple() {
	str, err := keg.GrepTags(`testdata/samplekeg`, `foo,bar`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
	// Output:
	// foo 2 6 3
	// bar 8
}
