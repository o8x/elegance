package app

import (
	"context"
	"math"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	py "github.com/mozillazg/go-pinyin"
	"github.com/o8x/jk/v2/response"

	"github.com/o8x/elegance/app/database"
	"github.com/o8x/elegance/app/database/queries"
	"github.com/o8x/elegance/app/loader"
	"github.com/o8x/elegance/app/pinyin"
	"github.com/o8x/elegance/app/utils"
)

func FindSource(c echo.Context) error {
	id := utils.ParseInt64Default(c.Param("id"), 0)
	if id == 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	source, err := database.GetQueries().FindSource(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, response.OK(source))
}

func ConvertPinyin(c echo.Context) error {
	word := c.Param("word")
	lazy, _ := pinyin.ParseLazy(word, py.Normal)
	lazyTone, _ := pinyin.ParseLazy(word, py.Tone)
	_, lastPy := pinyin.ParseLazy(word, py.Finals)
	_, lastTone := pinyin.ParseLazy(word, py.FinalsTone)
	_, initials := pinyin.ParseLazy(word, py.Initials)
	_, rhymeTone := pinyin.ParseLazy(word, py.FinalsTone3)

	return c.JSON(http.StatusOK, response.OK(map[string]any{
		"word":            word,
		"pinyin":          strings.Join(lazy, " "),
		"pinyin_tone":     strings.Join(lazyTone, " "),
		"rhyme":           lastPy,
		"rhyme_tone":      lastTone,
		"rhyme_initials":  initials,
		"rhyme_tone_name": rhymeTone[len(rhymeTone)-1:],
	}))
}

func SearchRhymeWords(c echo.Context) error {
	surname := c.Param("surname")
	page := utils.ParseInt64Default(c.QueryParam("page"), 1)
	pagesize := utils.ParseInt64Default(c.QueryParam("pagesize"), 20)
	length := utils.ParseInt64Default(c.QueryParam("length"), 2)

	_, lastPy := pinyin.ParseLazy(surname, py.Finals)

	ctx := context.Background()
	q := database.GetQueries()
	count, err := q.GetSearchRhymeWordsCount(ctx, queries.GetSearchRhymeWordsCountParams{
		WordLength: length,
		LastPinyin: lastPy,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err))
	}

	words, err := q.SearchRhymeWords(ctx, queries.SearchRhymeWordsParams{
		WordLength: length,
		LastPinyin: lastPy,
		Limit:      pagesize,
		Offset:     page * pagesize,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err))
	}

	return c.JSON(http.StatusOK, response.OK(map[string]any{
		"surname":   surname,
		"vowel":     lastPy,
		"sum_page":  math.Ceil(float64(count) / float64(pagesize)),
		"sum_count": count,
		"words":     words,
	}))
}

func SearchRhymeToneWords(c echo.Context) error {
	surname := c.Param("surname")
	page := utils.ParseInt64Default(c.QueryParam("page"), 1)
	pagesize := utils.ParseInt64Default(c.QueryParam("pagesize"), 20)
	length := utils.ParseInt64Default(c.QueryParam("length"), 2)

	lazy, _ := pinyin.ParseLazy(surname, py.Tone)
	_, lastPy := pinyin.ParseLazy(surname, py.Finals)
	_, lastTone := pinyin.ParseLazy(surname, py.FinalsTone)

	ctx := context.Background()
	q := database.GetQueries()
	count, err := q.GetSearchRhymeToneWordsCount(ctx, queries.GetSearchRhymeToneWordsCountParams{
		WordLength:     length,
		LastPinyinTone: lastPy,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err))
	}

	words, err := q.SearchRhymeToneWords(ctx, queries.SearchRhymeToneWordsParams{
		WordLength:     length,
		LastPinyinTone: lastPy,
		Limit:          pagesize,
		Offset:         page * pagesize,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err))
	}

	return c.JSON(http.StatusOK, response.OK(map[string]any{
		"surname":             surname,
		"surname_pinyin_tone": strings.Join(lazy, " "),
		"vowel":               lastPy,
		"vowel_tone":          lastTone,
		"sum_page":            math.Ceil(float64(count) / float64(pagesize)),
		"sum_count":           count,
		"words":               words,
	}))
}

func GetRandomWords(c echo.Context) error {
	surname := c.Param("surname")
	page := utils.ParseInt64Default(c.QueryParam("page"), 1)
	pagesize := utils.ParseInt64Default(c.QueryParam("pagesize"), 20)
	length := utils.ParseInt64Default(c.QueryParam("length"), 2)

	lazy, _ := pinyin.ParseLazy(surname, py.Tone)
	_, lastPy := pinyin.ParseLazy(surname, py.Finals)
	_, lastTone := pinyin.ParseLazy(surname, py.FinalsTone)

	ctx := context.Background()
	q := database.GetQueries()
	count, err := q.GetWordsCount(ctx, length)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err))
	}

	words, err := q.GetRandomWords(ctx, queries.GetRandomWordsParams{
		WordLength: length,
		Limit:      pagesize,
		Offset:     page * pagesize,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err))
	}

	return c.JSON(http.StatusOK, response.OK(map[string]any{
		"surname":             surname,
		"surname_pinyin_tone": strings.Join(lazy, " "),
		"vowel":               lastPy,
		"vowel_tone":          lastTone,
		"sum_page":            math.Ceil(float64(count) / float64(pagesize)),
		"sum_count":           count,
		"words":               words,
	}))
}

func GetSourcesByWord(c echo.Context) error {
	word := c.Param("word")
	page := utils.ParseInt64Default(c.QueryParam("page"), 1)
	pagesize := utils.ParseInt64Default(c.QueryParam("pagesize"), 20)

	ctx := context.Background()
	q := database.GetQueries()
	count, err := q.GetSourcesByWordCount(ctx, word)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err))
	}

	list, err := q.GetSourcesByWord(ctx, queries.GetSourcesByWordParams{
		Word:   word,
		Limit:  pagesize,
		Offset: page * pagesize,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err))
	}

	return c.JSON(http.StatusOK, response.OK(map[string]any{
		"sum_page":  math.Ceil(float64(count) / float64(pagesize)),
		"sum_count": count,
		"list":      list,
	}))
}

func ImportSource(c echo.Context) error {
	var params loader.ImportParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(
			err))
	}

	if err := loader.Import(params); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.Error(err),
		)
	}

	return c.JSON(http.StatusOK, response.NoContent)
}
