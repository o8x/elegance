package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/o8x/jk/v2/args"

	"github.com/o8x/elegance/app/loader"
)

func main() {
	a := args.Args{
		Executable: "",
		App: &args.App{
			Name:  "诗词数据导入工具",
			Usage: "go run github.com/o8x/elegance/cmd/elegance -d ../data",
		},
		Flags: []*args.Flag{
			{
				Name:        []string{"-d", "-dir"},
				Description: "扫描目录中的所有文件名以 .txt 结尾的文件并导入数据库，文本格式：第一行名字，第二行年代，第三行作者，其他行为正文",
				Required:    true,
			},
		},
	}

	if err := a.Parse(); err != nil {
		a.PrintHelpExit(err)
	}

	root := a.GetX("dir")
	dir, err := os.ReadDir(root)
	if err != nil {
		a.PrintErrorExit(err)
	}

	log.Printf("扫描到资源总数: %d\n", len(dir))

	for i, d := range dir {
		if !strings.HasSuffix(d.Name(), ".txt") {
			continue
		}

		bs, err := os.ReadFile(fmt.Sprintf("%s/%s", root, d.Name()))
		if err != nil {
			fmt.Println(err)
			continue
		}

		if bytes.Contains(bs, []byte("std")) {
			continue
		}

		data := strings.Split(string(bs), "\n")
		if len(data) < 4 {
			continue
		}

		params := loader.ImportParams{
			Title:    data[0],
			Age:      data[1],
			Name:     data[2],
			Contents: strings.TrimSpace(strings.Join(data[3:], "\n")),
		}

		log.Printf("正在导入资源: %d/%d, 标题: %s\n", i, len(dir), params.Title)
		if err = loader.Import(params); err != nil {
			continue
		}
	}

	log.Println("资源导入完成")
}
