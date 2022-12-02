package keg_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rwxrob/keg"
)

func ExampleDex_json() {
	date := time.Date(2022, 12, 10, 6, 10, 4, 0, time.UTC)
	d := keg.DexEntry{U: date, N: 2, T: `Some title`}
	byt, err := json.Marshal(d)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(byt))
	// Output:
	// {"U":"2022-12-10T06:10:04Z","T":"Some title","N":2,"HBeg":0,"HEnd":0}

}

func ExampleDex_string() {
	date := time.Date(2022, 12, 10, 6, 10, 4, 0, time.UTC)
	d := keg.DexEntry{U: date, N: 2, T: `Some title`}
	fmt.Println(d)
	fmt.Println(d.MD())
	// Output:
	// * 2022-12-10 06:10:04Z [Some title](../2)
	// * 2022-12-10 06:10:04Z [Some title](../2)
}

func ExampleDex_tsv() {
	date := time.Date(2022, 12, 10, 6, 10, 4, 0, time.UTC)
	d := keg.DexEntry{U: date, N: 2, T: `Some title`}
	fmt.Println(d.TSV())
	// Output:
	// 2	2022-12-10 06:10:04Z	Some title
}

/*
func ExampleDex_Pretty() {
	date := time.Date(2022, 12, 10, 6, 10, 4, 0, time.UTC)
	d1 := keg.DexEntry{U: date, N: 2000, T: `Some title`}
	d2 := keg.DexEntry{U: date, N: 1, T: `Another title`}
	dex := keg.Dex{d1, d2}
	fmt.Println(dex.Pretty())
	// Output:
	// ignored
}
*/

/*
func ExampleDex_Random() {
	date := time.Date(2022, 12, 10, 6, 10, 4, 0, time.UTC)
	dex := keg.Dex{
		{U: date, N: 2, T: `Some title`},
		{U: date, N: 3, T: `Other title`},
		{U: date, N: 4, T: `Another title`},
		{U: date, N: 5, T: `Yet another title`},
	}
	fmt.Println(dex.Random())
	// Output:
	// ignored
}
*/

func ExampleHaveDex() {
	fmt.Println(keg.HaveDex(`testdata/samplekeg`))
	fmt.Println(keg.HaveDex(`testdata/nothing`))
	// Output:
	// true
	// false
}

func ExampleDexEntry_Pretty() {

	e := keg.DexEntry{T: "Some"}
	fmt.Printf("%q\n", e.Pretty())
	e.HBeg = 2
	e.HEnd = 3
	fmt.Printf("%q\n", e.Pretty())

	// Output:
	// "\x1b[30m0001-01-01 00:00Z \x1b[32m0 \x1b[37mSome\x1b[0m\n"
	// "\x1b[30m0001-01-01 00:00Z \x1b[32m0 \x1b[37mSo\x1b[31mm\x1b[37me\x1b[0m\n"

}
