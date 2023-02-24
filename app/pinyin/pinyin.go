package pinyin

import (
	"github.com/mozillazg/go-pinyin"
)

func ParseLazy(s string, style int) ([]string, string) {
	a := pinyin.NewArgs()
	a.Style = style
	tones := pinyin.LazyPinyin(s, a)

	lastTone := ""
	if len(tones) != 0 {
		lastTone = tones[len(tones)-1]
	}

	return tones, lastTone
}
