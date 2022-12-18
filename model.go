package keg

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rwxrob/choose"
	"github.com/rwxrob/fs"
	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/json"
	"github.com/rwxrob/keg/kegml"
	"github.com/rwxrob/term"
)

const IsoDateFmt = `2006-01-02 15:04:05Z`
const IsoDateExpStr = `\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\dZ`

// Local contains a name to full path mapping for kegs stored locally.
type Local struct {
	Name string
	Path string
}

// DexEntry represents a single line in an index (usually the changes.md
// or nodes.tsv file). All three fields are always required.
type DexEntry struct {
	U    time.Time // updated
	T    string    // title
	N    int       // node id (also see ID)
	HBeg int       // start of highlighted
	HEnd int       // end of highlighted
}

// Update gets the entry for the target keg at kegpath by looking up the
// latest change to any file within it and parsing the title.
func (e *DexEntry) Update(kegpath string) error {
	var err error
	dir := filepath.Join(kegpath, e.ID())
	_, i := fs.LatestChange(dir)
	if i != nil {
		e.U = i.ModTime()
	}
	e.T, err = kegml.ReadTitle(filepath.Join(dir, `README.md`))
	return err
}

// MarshalJSON produces JSON text that contains one DexEntry per line
// that has not been HTML escaped (unlike the default) and that uses
// a consistent DateTime format. Note that the (broken) encoding/json
// encoder is not used at all.
func (e *DexEntry) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 0))
	buf.WriteRune('{')
	buf.WriteString(`"U":"` + e.U.Format(IsoDateFmt) + `",`)
	buf.WriteString(`"N":` + strconv.Itoa(e.N) + `,`)
	buf.WriteString(`"T":"` + json.Escape(e.T) + `"`)
	buf.WriteRune('}')
	return buf.Bytes(), nil
}

func (e *DexEntry) TSV() string {
	return fmt.Sprintf("%v\t%v\t%v", e.N, e.U.Format(IsoDateFmt), e.T)
}

// ID returns the node identifier as a string instead of an integer.
// Returns an empty string if unable to parse the integer.
func (e *DexEntry) ID() string { return strconv.Itoa(e.N) }

// MD returns the entry as a single Markdown list item for inclusion in
// the dex/nodex.md file:
//
//     1. Second last changed in UTC in ISO8601 (RFC3339)
//     2. Current title (always first line of README.md)
//     2. Unique node integer identifier
//
// Note that the second of last change is based on *any* file within the
// node directory changing, not just the README.md or meta files.
func (e *DexEntry) MD() string {
	return fmt.Sprintf(
		"* %v [%v](../%v)",
		e.U.Format(IsoDateFmt),
		e.T, e.N,
	)
}

// String implements fmt.Stringer interface as MD.
func (e DexEntry) String() string { return e.MD() }

// Asinclude returns a KEGML include link list item without the time
// suitable for creating include blocks in node files.
func (e *DexEntry) AsInclude() string {
	return fmt.Sprintf("* [%v](../%v)", e.T, e.N)
}

// Pretty returns a string with pretty colors.
func (e *DexEntry) Pretty() string {
	nwidth := len(e.ID())
	text := e.T
	if e.HBeg > 0 || e.HEnd > 0 {
		before := e.T[0:e.HBeg]
		hilight := e.T[e.HBeg:e.HEnd]
		after := e.T[e.HEnd:]
		text = before + term.Red + hilight + term.White + after
	}
	return fmt.Sprintf(
		"%v%v %v%-"+strconv.Itoa(nwidth)+"v %v%v%v\n",
		term.Black, e.U.Format(`2006-01-02 15:04Z`),
		term.Green, e.N,
		term.White, text,
		term.Reset,
	)
}

// -------------------------------- Dex -------------------------------

// Dex is a collection of DexEntry structs. This allows mapping methods
// for its serialization to different output formats.
type Dex []*DexEntry

// MarshalJSON produces JSON text that contains one DexEntry per line
// that has not been HTML escaped (unlike the default).
func (d *Dex) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 0))
	buf.WriteString("[")
	for _, entry := range *d {
		byt, _ := entry.MarshalJSON()
		buf.Write(byt)
		buf.WriteString(",\n")
	}
	byt := buf.Bytes()
	byt[len(byt)-2] = ']'
	return byt, nil
}

// Lookup does a linear search through the Dex for one with the passed
// id and if found returns along with the current position in the Dex
// slice, otherwise returns nil and negative 1.
func (d Dex) Lookup(id int) (int, *DexEntry) {
	for i, node := range d {
		if node.N == id {
			return i, d[i]
		}
	}
	return -1, nil
}

// String fulfills the fmt.Stringer interface as JSON. Any error returns
// a "null" string.
func (e Dex) String() string { return e.TSV() }

// MD renders the entire Dex as a Markdown list suitable for the
// standard dex/changes.md file.
func (e Dex) MD() string {
	var str string
	for _, entry := range e {
		str += entry.MD() + "\n"
	}
	return str
}

// AsIncludes renders the entire Dex as a KEGML include list (markdown
// bulleted list) and cab be useful from within editing sessions to
// include from the current keg without leaving the terminal editor.
func (e Dex) AsIncludes() string {
	var str string
	for _, entry := range e {
		str += entry.AsInclude() + "\n"
	}
	return str
}

// TSV renders the entire Dex as a loadable tab-separated values file.
func (e Dex) TSV() string {
	var str string
	for _, entry := range e {
		str += entry.TSV() + "\n"
	}
	return str
}

// Last returns the DexEntry with the highest integer value identifier.
func (d Dex) Last() *DexEntry {
	entry := new(DexEntry)
	for _, e := range d {
		if e.N > entry.N {
			entry = e
		}
	}
	return entry
}

// First returns the first content node created. Returns nil if it
// cannot determine. The first node is always the one with the lowest
// integer identifier greater than 0.
func (d Dex) First() *DexEntry {
	entry := new(DexEntry)
	for _, e := range d {
		if e.N > 0 && (entry.N == 0 || entry.N > e.N) {
			entry = e
		}
	}
	return entry
}

// LastChanged returns the entry with the most recent modification time.
func (d Dex) LastChanged() *DexEntry {
	entry := new(DexEntry)
	for _, e := range d {
		if e.U.After(entry.U) {
			entry = e
		}
	}
	return entry
}

// FirstChanged returns the entry with the oldest modification time.
func (d Dex) FirstChanged() *DexEntry {
	entry := new(DexEntry)
	for _, e := range d {
		if e.U.Before(entry.U) {
			entry = e
		}
	}
	return entry
}

// LastIdString returns Last as string.
func (d Dex) LastIdString() string { return strconv.Itoa(d.Last().N) }

// LastIdWidth returns width of Last integer identifier.
func (d Dex) LastIdWidth() int { return len(d.LastIdString()) }

// LastChangedIdString returns Last as string.
func (d Dex) LastChangedIdString() string { return strconv.Itoa(d.LastChanged().N) }

// LastChangedIdWidth returns width of Last integer identifier.
func (d Dex) LastChangedIdWidth() int { return len(d.LastChangedIdString()) }

// Pretty returns a string with pretty color string with no time stamps
// rendered in more readable way.
func (d Dex) Pretty() string {
	var str string
	nwidth := d.LastIdWidth()
	for _, e := range d {
		text := e.T
		if e.HBeg > 0 || e.HEnd > 0 {
			before := e.T[0:e.HBeg]
			hilight := e.T[e.HBeg:e.HEnd]
			after := e.T[e.HEnd:]
			text = before + term.Red + hilight + term.White + after
		}
		str += fmt.Sprintf(
			"%v%"+strconv.Itoa(nwidth)+"v %v%v%v\n",
			term.Green, e.N,
			term.White, text,
			term.Reset,
		)
	}
	return str
}

// PrettyLines returns Pretty but each line separate and without line
// return.
func (d Dex) PrettyLines() []string {
	lines := make([]string, 0, len(d))
	nwidth := d.LastIdWidth()
	for _, e := range d {
		text := e.T
		if e.HBeg > 0 || e.HEnd > 0 {
			before := e.T[0:e.HBeg]
			hilight := e.T[e.HBeg:e.HEnd]
			after := e.T[e.HEnd:]
			text = before + term.Red + hilight + term.White + after
		}
		lines = append(lines, fmt.Sprintf(
			"%v%-"+strconv.Itoa(nwidth)+"v %v%v%v",
			term.Green, e.N,
			term.White, text,
			term.Reset,
		))
	}
	return lines
}

// ByID sorts the Dex from lowest to highest node ID integer. A pointer
// to self is returned for convenience.
func (d Dex) ByID() Dex {
	sort.Slice(d, func(i, j int) bool {
		return d[i].N < d[j].N
	})
	return d
}

// ByChanges sorts the Dex from most recently changed to oldest. A pointer
// to self is returned for convenience.
func (d Dex) ByChanges() Dex {
	sort.Slice(d, func(i, j int) bool {
		return d[i].U.After(d[j].U)
	})
	return d
}

// Add appends the entry to the Dex.
func (d *Dex) Add(entry *DexEntry) {
	(*d) = append((*d), entry)
}

// WithTitleText returns a new Dex from self with all nodes that do not
// contain the text substring in the title filtered out.
func (e *Dex) WithTitleText(keyword string) Dex {
	dex := Dex{}
	for _, d := range *e {
		i := strings.Index(strings.ToLower(d.T), strings.ToLower(keyword))
		if i >= 0 {
			d.HBeg = i
			d.HEnd = i + len(keyword)
			dex = append(dex, d)
		}
	}
	return dex
}

// WithTitleTextExp returns a new Dex from self with all nodes that do not
// contain the regular expression matches in the title filtered out.
func (e *Dex) WithTitleTextExp(re *regexp.Regexp) Dex {
	dex := Dex{}
	for _, d := range *e {
		i := re.FindStringIndex(d.T)
		if i != nil {
			d.HBeg = i[0]
			d.HEnd = i[1]
			dex = append(dex, d)
		}
	}
	return dex
}

// ChooseWithTitleText returns a single *DexEntry for the keyword
// passed. If there are more than one then user is prompted to choose
// from list sent to the terminal.
func (d *Dex) ChooseWithTitleText(key string) *DexEntry {
	hits := d.WithTitleText(key)
	switch len(hits) {
	case 1:
		return hits[0]
	case 0:
		return nil
	default:
		i, _, err := choose.From(hits.PrettyLines())
		if err != nil {
			return nil
		}
		if i < 0 {
			return nil
		}
		return hits[i]
	}
}

// ChooseWithTitleTextExp returns a single *DexEntry for the regular
// expression matches passed. If there are more than one then user is
// prompted to choose from list sent to the terminal.
func (d *Dex) ChooseWithTitleTextExp(re *regexp.Regexp) *DexEntry {
	hits := d.WithTitleTextExp(re)
	switch len(hits) {
	case 1:
		return hits[0]
	case 0:
		return nil
	default:
		i, _, err := choose.From(hits.PrettyLines())
		if err != nil {
			return nil
		}
		if i < 0 {
			return nil
		}
		return hits[i]
	}
}

// Random returns a random entry.
func (d Dex) Random() *DexEntry {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(d))
	return d[i]
}

// Delete removes the specified entry shrinking the slice without
// changing the underlying size of the supporting array while
// maintaining references to each DexEntry. Note that the persisted
// content node directory may still exist. This method only affects the
// Dex itself. Stops at first match.
func (d *Dex) Delete(entry *DexEntry) {
	var index int
	for i, it := range *d {
		if it == entry || it.N == entry.N {
			index = i
			break
		}
	}
	for i := index; i < len(*d)-1; i++ {
		(*d)[i] = (*d)[i+1]
	}
	(*d) = (*d)[:len(*d)-1]
}

// ----------------------------- TagsList -----------------------------

type TagsMap map[string][]string

func (tl TagsMap) String() string {
	var str string
	for k, v := range tl {
		str += k + " " + strings.Join(v, " ") + "\n"
	}
	return str
}

func (tl TagsMap) MarshalText() ([]byte, error) {
	var str string
	for k, v := range tl {
		str += k + " " + strings.Join(v, " ") + "\n"
	}
	return []byte(str), nil
}

//Write writes the marshaled text of a TagsMap to the file at path.
func (tl TagsMap) Write(path string) error {
	return file.Overwrite(path, tl.String())
}

// UnmarshalText parses the tag lines items from the bytes buffer and
// sets the key pair for that tag to the values overwriting any that
// were already set.
func (tl TagsMap) UnmarshalText(buf []byte) error {
	s := bufio.NewScanner(strings.NewReader(string(buf)))
	for s.Scan() {
		line := s.Text()
		f := strings.Split(line, " ")
		switch len(f) {
		case 0:
			return fmt.Errorf(_InvalidTagLine, line)
		case 1:
			return nil
		default:
			tl[f[0]] = f[1:]
		}
	}
	return nil
}
