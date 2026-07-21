// Package gocolor は、ターミナル出力に ANSI エスケープシーケンスを使って
// 色や装飾（太字・下線など）を付けるための軽量なライブラリです。
//
// 基本的な使い方:
//
//	gocolor.Println(gocolor.FgGreen, "成功しました")
//	msg := gocolor.Sprintf(gocolor.FgRed, gocolor.Bold, "エラー: %s", err)
//
// より細かく制御したい場合は Color 型を直接使います。
//
//	c := gocolor.New(gocolor.FgCyan, gocolor.Underline)
//	c.Println("見出しテキスト")
//
// NO_COLOR 環境変数が設定されている場合、または出力先が端末でない場合
// （パイプやファイルへのリダイレクトなど）は、自動的に色付けが無効になります。
package gocolor

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Attribute は ANSI エスケープシーケンスのパラメータ（色や装飾）を表します。
type Attribute int

const (
	// スタイル
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

const (
	// 前景色（通常）
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

const (
	// 背景色（通常）
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

const (
	// 前景色（明るい版）
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

const (
	// 背景色（明るい版）
	BgHiBlack Attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

const escape = "\x1b"

// NoColor は色出力を全体で無効化するかどうかを制御するグローバルフラグです。
// デフォルトでは、NO_COLOR 環境変数が設定されているか、標準出力が端末でない
// 場合に自動的に true が設定されます。
var NoColor = noColorDefault()

func noColorDefault() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return true
	}
	return !isTerminal(os.Stdout)
}

// isTerminal は与えられたファイルが端末に接続されているかどうかを判定します。
// 外部依存を避けるため、標準ライブラリのみで簡易判定します。
func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// Color は一連の Attribute（色・装飾）をまとめて保持し、
// 文字列の色付けや出力を行うための型です。
type Color struct {
	params  []Attribute
	extra   []string // 256色・RGB(トゥルーカラー)など、単純な Attribute で表せないコード
	noColor *bool    // nil の場合はパッケージ変数 NoColor に従う
}

// New は指定した Attribute を持つ新しい Color を作成します。
//
//	c := gocolor.New(gocolor.FgRed, gocolor.Bold)
func New(value ...Attribute) *Color {
	return &Color{params: append([]Attribute{}, value...)}
}

// Add は Attribute を追加し、メソッドチェーンできるように自身を返します。
//
//	c := gocolor.New(gocolor.FgRed).Add(gocolor.Bold).Add(gocolor.Underline)
func (c *Color) Add(value ...Attribute) *Color {
	c.params = append(c.params, value...)
	return c
}

// AddFg256 は 256色パレット（0〜255）から前景色を追加します。
// 0-7: 通常色, 8-15: 明るい色, 16-231: 6x6x6 のカラーキューブ, 232-255: グレースケール。
//
//	c := gocolor.New().AddFg256(208) // オレンジ色
func (c *Color) AddFg256(n int) *Color {
	c.extra = append(c.extra, fmt.Sprintf("38;5;%d", clamp256(n)))
	return c
}

// AddBg256 は 256色パレット（0〜255）から背景色を追加します。
func (c *Color) AddBg256(n int) *Color {
	c.extra = append(c.extra, fmt.Sprintf("48;5;%d", clamp256(n)))
	return c
}

// AddFgRGB は 24bit トゥルーカラー（R, G, B 各 0〜255）で前景色を追加します。
// 対応する端末（多くのモダンターミナル）でのみ正しく表示されます。
//
//	c := gocolor.New().AddFgRGB(255, 105, 180) // ホットピンク
func (c *Color) AddFgRGB(r, g, b int) *Color {
	c.extra = append(c.extra, fmt.Sprintf("38;2;%d;%d;%d", clamp255(r), clamp255(g), clamp255(b)))
	return c
}

// AddBgRGB は 24bit トゥルーカラー（R, G, B 各 0〜255）で背景色を追加します。
func (c *Color) AddBgRGB(r, g, b int) *Color {
	c.extra = append(c.extra, fmt.Sprintf("48;2;%d;%d;%d", clamp255(r), clamp255(g), clamp255(b)))
	return c
}

func clamp256(n int) int {
	if n < 0 {
		return 0
	}
	if n > 255 {
		return 255
	}
	return n
}

func clamp255(n int) int { return clamp256(n) }

// EnableColor はこの Color インスタンスに限り、色出力を強制的に有効にします。
func (c *Color) EnableColor() *Color {
	v := false
	c.noColor = &v
	return c
}

// DisableColor はこの Color インスタンスに限り、色出力を強制的に無効にします。
func (c *Color) DisableColor() *Color {
	v := true
	c.noColor = &v
	return c
}

func (c *Color) isNoColor() bool {
	if c.noColor != nil {
		return *c.noColor
	}
	return NoColor
}

// sequence はこの Color が持つ Attribute と拡張コード（256色・RGB）から
// ANSI エスケープシーケンスのパラメータ部分（例: "1;31" や "38;5;208"）を組み立てます。
func (c *Color) sequence() string {
	parts := make([]string, 0, len(c.params)+len(c.extra))
	for _, v := range c.params {
		parts = append(parts, strconv.Itoa(int(v)))
	}
	parts = append(parts, c.extra...)
	return strings.Join(parts, ";")
}

func (c *Color) wrap(s string) string {
	if c.isNoColor() || (len(c.params) == 0 && len(c.extra) == 0) {
		return s
	}
	return fmt.Sprintf("%s[%sm%s%s[%dm", escape, c.sequence(), s, escape, Reset)
}

// Sprint は引数を fmt.Sprint と同様に連結し、色付けした文字列を返します。
func (c *Color) Sprint(a ...interface{}) string {
	return c.wrap(fmt.Sprint(a...))
}

// Sprintln は引数を fmt.Sprintln と同様に連結し、色付けした文字列を返します。
func (c *Color) Sprintln(a ...interface{}) string {
	return c.wrap(fmt.Sprintln(a...))
}

// Sprintf は書式付きで色付けした文字列を返します。
func (c *Color) Sprintf(format string, a ...interface{}) string {
	return c.wrap(fmt.Sprintf(format, a...))
}

// Fprint は色付けした文字列を指定した io.Writer に書き込みます。
func (c *Color) Fprint(w io.Writer, a ...interface{}) (int, error) {
	return fmt.Fprint(w, c.Sprint(a...))
}

// Fprintln は色付けした文字列を改行付きで指定した io.Writer に書き込みます。
func (c *Color) Fprintln(w io.Writer, a ...interface{}) (int, error) {
	return fmt.Fprintln(w, c.wrap(fmt.Sprint(a...)))
}

// Fprintf は書式付きで色付けした文字列を指定した io.Writer に書き込みます。
func (c *Color) Fprintf(w io.Writer, format string, a ...interface{}) (int, error) {
	return fmt.Fprint(w, c.Sprintf(format, a...))
}

// Print は色付けした文字列を標準出力に書き込みます。
func (c *Color) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(os.Stdout, c.Sprint(a...))
}

// Println は色付けした文字列を改行付きで標準出力に書き込みます。
func (c *Color) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(os.Stdout, c.wrap(fmt.Sprint(a...)))
}

// Printf は書式付きで色付けした文字列を標準出力に書き込みます。
func (c *Color) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprint(os.Stdout, c.Sprintf(format, a...))
}
