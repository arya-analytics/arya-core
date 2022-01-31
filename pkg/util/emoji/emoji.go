package emoji

import (
	log "github.com/sirupsen/logrus"
	"strconv"
)

// New creates and converts a print-able emoji from a unicode string.
func New(s string) string {
	r, err := strconv.ParseInt(s, 16, 32)
	if err != nil {
		log.Fatal("Unable to parse emoji")
	}
	return string(r)
}

// || EMOJI VARS ||
// Pulled from http://www.unicode.org/emoji/charts/full-emoji-list.html.
var (
	Frog    = New("1F438")
	Check   = New("2705")
	Flame   = New("1F525")
	Drop    = New("1F4A7")
	Sparks  = New("2728")
	Rainbow = New("1F308")
	Bolt    = New("26A1")
	Tools   = New("1F6E0")
	Bison   = New("1F9AC")
)
