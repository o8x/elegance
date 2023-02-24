package word

import (
	"github.com/huichen/sego"
)

var (
	DictFile = "dictionary.txt"
)

var x sego.Segmenter

func init() {
	x.LoadDictionary(DictFile)
}

func Parse(content string) []string {
	var result []string
	for _, it := range x.Segment([]byte(content)) {
		result = append(result, it.Token().Text())
	}

	return result
}
