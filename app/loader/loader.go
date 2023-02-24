package loader

import (
	"context"
	"fmt"
	"strings"

	pinyin2 "github.com/mozillazg/go-pinyin"
	"github.com/o8x/jk/v2/puresqlite"

	"github.com/o8x/elegance/app/database"
	"github.com/o8x/elegance/app/database/queries"
	"github.com/o8x/elegance/app/pinyin"
	"github.com/o8x/elegance/app/utils"
	"github.com/o8x/elegance/app/word"
)

type ImportParams struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Age      string `json:"age"`
	Contents string `json:"contents"`
}

func Import(data ImportParams) error {
	if strings.Contains(data.Contents, "std") {
		return nil
	}

	ctx := context.Background()
	label := fmt.Sprintf("%s %s %s", data.Title, data.Age, data.Name)

	tx, err := puresqlite.Get().Begin()
	if err != nil {
		return err
	}

	// 不知道为什么开事务后，写性能会提升100倍左右
	q := database.GetQueries().WithTx(tx)
	id, err := q.CreateSource(ctx, queries.CreateSourceParams{
		Label:    label,
		Age:      data.Age,
		Title:    data.Title,
		Author:   data.Name,
		Contents: data.Contents,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	m := make(map[string]any)
	for _, w := range word.Parse(data.Contents) {
		if w == "，" || w == "。" || w == "\n" || w == "" || w == "." || w == "?" || w == "？" || w == "," || w == ";" || w == "；" {
			continue
		}

		if _, ok := m[w]; ok {
			continue
		}

		m[w] = nil

		py, _ := pinyin.ParseLazy(w, pinyin2.Normal)
		tones, _ := pinyin.ParseLazy(w, pinyin2.Tone)
		_, lastPy := pinyin.ParseLazy(w, pinyin2.Finals)
		_, lastTone := pinyin.ParseLazy(w, pinyin2.FinalsTone)
		_, rhymeTone := pinyin.ParseLazy(w, pinyin2.FinalsTone3)

		var toneName int64
		if len(rhymeTone) > 0 {
			toneName = utils.ParseInt64Default(rhymeTone[len(rhymeTone)-1:], 0)
		}

		if err := q.CreateWord(ctx, queries.CreateWordParams{
			SourceID:           id,
			Word:               w,
			Pinyin:             strings.Join(py, " "),
			PinyinTone:         strings.Join(tones, " "),
			LastPinyin:         lastPy,
			LastPinyinTone:     lastTone,
			LastPinyinToneName: toneName,
			WordLength:         int64(len(w) / 3),
		}); err != nil {
			delete(m, w)
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
