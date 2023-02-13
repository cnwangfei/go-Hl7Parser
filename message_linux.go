//go:build !windows
// +build !windows

package golevel7

import (
	"strings"
)

func (m *Message) parse() error {
	m.Value = []rune(strings.Trim(string(m.Value), "\n\r\x1c\x0b"))
	if err := m.parseSep(); err != nil {
		return err
	}
	r := strings.NewReader(string(m.Value))
	i := 0
	ii := 0
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			ch = eof
		}
		ii++
		switch {
		case ch == eof || (ch == endMsg && m.Delimeters.LFTermMsg):
			//just for safety: cannot reproduce this on windows
			safeii := map[bool]int{true: len(m.Value), false: ii}[ii > len(m.Value)]
			v := m.Value[i:safeii]
			if len(v) > 4 { // seg name + field sep
				seg := Segment{Value: v}
				seg.parse(&m.Delimeters)
				m.Segments = append(m.Segments, seg)
			}
			return nil
		case ch == segTerm:
			seg := Segment{Value: m.Value[i : ii-1]}
			seg.parse(&m.Delimeters)
			m.Segments = append(m.Segments, seg)
			i = ii
		case ch == m.Delimeters.Escape:
			ii++
			r.ReadRune()
		}
	}
}