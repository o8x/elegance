package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/o8x/jk/v2/args"

	"github.com/o8x/elegance/app"
)

var StaticRoot = "view/build"

func main() {
	a := args.Args{
		App: &args.App{
			Name:      "Elegance - Name From Poetries",
			Usage:     "./elegance -domain=elegance.stdout.com.cn",
			Copyright: "alex stdout.com.cn",
		},
		Flags: []*args.Flag{
			{
				Name:        []string{"-domain"},
				Description: "listen on tls for this domain.",
				Env:         []string{"ELEGANCE_TLS_DOMAIN"},
			},
		},
	}

	if err := a.Parse(); err != nil {
		a.PrintHelpExit(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	e.Static("/", StaticRoot)

	e.GET("/words/rhyme/:surname", app.SearchRhymeWords)
	e.GET("/words/rhyme_tone/:surname", app.SearchRhymeToneWords)
	e.GET("/words/random/:surname", app.GetRandomWords)

	e.GET("/word/pinyin/:word", app.ConvertPinyin)
	e.GET("/word/sources/:word", app.GetSourcesByWord)

	e.GET("/sources/:id", app.FindSource)
	e.POST("/sources/import", app.ImportSource)

	e.Logger.Infof("use static, root: %s", StaticRoot)

	domain, ok := a.Get("domain")
	if ok {
		e.Logger.Fatal(e.StartAutoTLS(fmt.Sprintf("%s:443", domain)))
	}
	e.Logger.Fatal(e.Start(":1323"))
}
