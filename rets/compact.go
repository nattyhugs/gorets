package rets

import (
	"encoding/xml"
	"strings"
)

// NewCompactData parses a CompactData from a start element
func NewCompactData(start xml.StartElement, decoder *xml.Decoder, delim string) (CompactData, error) {
	cd := CompactData{}
	cd.ID = start.Name.Local
	cd.Delimiter = delim
	cd.Attr = map[string]string{}
	for _, a := range start.Attr {
		cd.Attr[a.Name.Local] = a.Value
	}
	return cd, decoder.DecodeElement(&cd, &start)
}

// CompactData is the common compact decoded structure
type CompactData struct {
	ID        string
	Delimiter string
	Attr      map[string]string
	// parse these values out with decode
	CompactColumns CompactRow   `xml:"COLUMNS"`
	CompactRows    []CompactRow `xml:"DATA"`
}

// Columns parses the compact values for the cols
func (cd CompactData) Columns() Row {
	return cd.CompactColumns.Parse(cd.Delimiter)
}

// Rows provides callback to access each row
func (cd CompactData) Rows(each func(i int, row Row)) {
	for i, row := range cd.CompactRows {
		each(i, row.Parse(cd.Delimiter))
	}
}

// Map turns all rows into maps
func (cd CompactData) Map() map[string]string {
	index := cd.Indexer()
	entries := map[string]string{}
	cols := cd.Columns()
	cd.Rows(func(i int, r Row) {
		for _, c := range cols {
			entries[c] = index(c, r)
		}
	})
	return entries
}

// Indexer provices cached lookup for CompactData
type Indexer func(col string, row Row) string

// Indexer create the cache
func (cd *CompactData) Indexer() Indexer {
	index := make(map[string]int)
	for i, c := range cd.Columns() {
		index[c] = i
	}
	return func(col string, row Row) string {
		return row[index[col]]
	}
}

// Row is a string slice typedef for convenience
type Row []string

// CompactRow ...
type CompactRow string

// Parse ...
func (cr CompactRow) Parse(delim string) Row {
	split := strings.Split(string(cr), delim)
	return split[1 : len(split)-1]
}