package gocolor

import "fmt"

// Sprint は Attribute 群を先頭にまとめて渡し、残りの引数を色付けして返します。
// 便宜上のパッケージレベル関数で、内部で New(attrs...).Sprint(a...) を呼びます。
//
//	s := gocolor.Sprint(gocolor.FgRed, "エラー発生")
func Sprint(attrs []Attribute, a ...interface{}) string {
	return New(attrs...).Sprint(a...)
}

// Sprintf は Attribute 群と書式文字列を渡し、色付けした文字列を返します。
//
//	s := gocolor.Sprintf([]gocolor.Attribute{gocolor.FgRed, gocolor.Bold}, "エラー: %s", err)
func Sprintf(attrs []Attribute, format string, a ...interface{}) string {
	return New(attrs...).Sprintf(format, a...)
}

// Println は Attribute 群を指定し、色付けした文字列を改行付きで標準出力に書き込みます。
//
//	gocolor.Println([]gocolor.Attribute{gocolor.FgGreen}, "成功しました")
func Println(attrs []Attribute, a ...interface{}) (int, error) {
	return New(attrs...).Println(a...)
}

// 以下は各色専用のショートカット関数です。よく使う色をすぐに使えるようにします。
// 例: gocolor.Red("エラー: %s", err) は赤字で Sprintf した文字列を返します。

func makeColorFunc(attr Attribute) func(format string, a ...interface{}) string {
	c := New(attr)
	return func(format string, a ...interface{}) string {
		if len(a) == 0 {
			return c.Sprint(format)
		}
		return c.Sprintf(format, a...)
	}
}

var (
	// Black は黒色で文字列を装飾します。
	Black = makeColorFunc(FgBlack)
	// Red は赤色で文字列を装飾します。
	Red = makeColorFunc(FgRed)
	// Green は緑色で文字列を装飾します。
	Green = makeColorFunc(FgGreen)
	// Yellow は黄色で文字列を装飾します。
	Yellow = makeColorFunc(FgYellow)
	// Blue は青色で文字列を装飾します。
	Blue = makeColorFunc(FgBlue)
	// Magenta はマゼンタ色で文字列を装飾します。
	Magenta = makeColorFunc(FgMagenta)
	// Cyan はシアン色で文字列を装飾します。
	Cyan = makeColorFunc(FgCyan)
	// White は白色で文字列を装飾します。
	White = makeColorFunc(FgWhite)

	// HiRed は明るい赤色で文字列を装飾します。
	HiRed = makeColorFunc(FgHiRed)
	// HiGreen は明るい緑色で文字列を装飾します。
	HiGreen = makeColorFunc(FgHiGreen)
	// HiYellow は明るい黄色で文字列を装飾します。
	HiYellow = makeColorFunc(FgHiYellow)
	// HiBlue は明るい青色で文字列を装飾します。
	HiBlue = makeColorFunc(FgHiBlue)
)

// Fg256 は 256色パレット（0〜255）を前景色に使う Color を作成します。
//
//	gocolor.Fg256(208).Println("オレンジ色のテキスト")
func Fg256(n int) *Color {
	return New().AddFg256(n)
}

// Bg256 は 256色パレット（0〜255）を背景色に使う Color を作成します。
func Bg256(n int) *Color {
	return New().AddBg256(n)
}

// FgRGB は 24bit トゥルーカラー（R, G, B 各 0〜255）を前景色に使う Color を作成します。
//
//	gocolor.FgRGB(255, 105, 180).Println("ホットピンクのテキスト")
func FgRGB(r, g, b int) *Color {
	return New().AddFgRGB(r, g, b)
}

// BgRGB は 24bit トゥルーカラー（R, G, B 各 0〜255）を背景色に使う Color を作成します。
func BgRGB(r, g, b int) *Color {
	return New().AddBgRGB(r, g, b)
}

// PrintRed など、標準出力にそのまま書き込みたい場合のための補助関数群。
// 必要に応じて自由に増やせます。

// PrintlnColored は色を指定して改行付きで標準出力に書き込む簡易関数です。
//
//	gocolor.PrintlnColored(gocolor.FgRed, "失敗:", err)
func PrintlnColored(attr Attribute, a ...interface{}) {
	fmt.Println(New(attr).Sprint(fmt.Sprint(a...)))
}
