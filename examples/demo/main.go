package main

import (
	"fmt"

	"github.com/gvhrerf/gocolor"
)

func main() {
	// 常に色を出す(パイプ出力でも見えるように)ためのデモ用Color
	demo := func(c *gocolor.Color) *gocolor.Color { return c.EnableColor() }

	fmt.Println("=== 基本の色 ===")
	demo(gocolor.New(gocolor.FgRed)).Println("赤色のテキスト")
	demo(gocolor.New(gocolor.FgGreen, gocolor.Bold)).Println("太字の緑色")
	fmt.Println(demo(gocolor.New(gocolor.FgCyan, gocolor.Underline)).Sprint("下線付きシアン"))

	fmt.Println("=== ショートカット関数 ===")
	fmt.Println(demo(gocolor.New(gocolor.FgRed)).Sprintf("Red: %s", "hello"))

	fmt.Println("=== 256色パレット ===")
	for _, n := range []int{196, 208, 226, 46, 51, 21, 129} {
		demo(gocolor.Fg256(n)).Printf("[%3d] ", n)
	}
	fmt.Println()

	fmt.Println("=== トゥルーカラー(RGB) ===")
	demo(gocolor.FgRGB(255, 105, 180)).Println("ホットピンク")
	demo(gocolor.New().AddFgRGB(0, 191, 255).AddBgRGB(20, 20, 20)).Println("背景付きスカイブルー")
}
