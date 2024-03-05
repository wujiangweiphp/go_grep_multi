package main

import (
	"github.com/urfave/cli/v2"
	scan_project "github.com/wujiangweiphp/go_grep_multi/scan-project"
	"log"
	"os"
)

func main() {
	app := cli.App{
		Name: "scan",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "需要扫描的目录路径",
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "lineCount",
				Aliases: []string{"l"},
				Usage:   "是否需要统计扫描行数",
			},
			&cli.StringFlag{
				Name:    "exclude",
				Aliases: []string{"e"},
				Usage:   "需要跳过的路径标识",
			},
			&cli.StringFlag{
				Name:    "include",
				Aliases: []string{"i"},
				Usage:   "需要包含的路径或文件标识",
			},
			&cli.BoolFlag{
				Name:        "regexp",
				Aliases:     []string{"r"},
				Usage:       "是否使用正则匹配",
				DefaultText: "默认使用字符串",
			},
			&cli.StringSliceFlag{
				Name:     "content",
				Aliases:  []string{"c"},
				Usage:    "需要匹配的内容",
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "onlyFile",
				Aliases: []string{"o"},
				Usage:   "只扫描文件名",
			},
			&cli.BoolFlag{
				Name:    "ignoreCase",
				Aliases: []string{"u"},
				Usage:   "只扫描文件名",
			},
		},
		Action: func(ctx *cli.Context) error {
			sp := scan_project.SmartScan{
				RootPath:       ctx.String("path"),
				IncludePathStr: ctx.String("include"),
				ExcludePathStr: ctx.String("exclude"),
				NeedLineCount:  ctx.Bool("lineCount"),
				NeedUseRegexp:  ctx.Bool("regexp"),
				MatchContent:   ctx.StringSlice("content"),
				OnlyFile:       ctx.Bool("onlyFile"),
				IgnoreCate:     ctx.Bool("ignoreCase"),
			}
			return sp.Scan()
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	return
}
