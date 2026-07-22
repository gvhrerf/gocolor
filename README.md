
# gocolor

Go のターミナル出力に色や装飾（太字・下線など）を付けるための軽量ライブラリです。
外部依存なし（標準ライブラリのみ）。16色 / 256色パレット / 24bit トゥルーカラー(RGB) に対応しています。

## インストール

```bash
go get github.com/gvhrerf/gocolor
```

## 使い方

### 基本の色

```go
package main

import "github.com/gvhrerf/gocolor"

func main() {
	gocolor.New(gocolor.FgRed).Println("エラーが発生しました")
	gocolor.New(gocolor.FgGreen, gocolor.Bold).Println("成功しました")

	s := gocolor.New(gocolor.FgCyan, gocolor.Underline).Sprintf("処理数: %d", 42)
	println(s)
}
```

### ショートカット関数

```go
fmt.Println(gocolor.Red("エラー: %s", err))
fmt.Println(gocolor.Green("OK"))
```

### 256色パレット

```go
gocolor.Fg256(208).Println("オレンジ色のテキスト") // 0-255 の256色パレット
gocolor.New().AddFg256(196).AddBg256(235).Println("背景付き")
```

### 24bit トゥルーカラー(RGB)

対応しているターミナル（多くのモダンターミナル）でのみ正しく表示されます。

```go
gocolor.FgRGB(255, 105, 180).Println("ホットピンク")
gocolor.New().AddFgRGB(0, 191, 255).AddBgRGB(20, 20, 20).Println("背景付きスカイブルー")
```

## 色を無効化する

- `NO_COLOR` 環境変数がセットされている場合
- 標準出力が端末でない場合（パイプ・ファイルへのリダイレクトなど）

上記のいずれかに該当すると、自動的に色付けが無効になります。

個別の `Color` インスタンスで強制的に有効/無効にすることもできます。

```go
c := gocolor.New(gocolor.FgRed)
c.DisableColor() // 常に無色
c.EnableColor()  // 常に色付き
```

## 対応している装飾（Attribute）

| 種類 | 定数 |
| --- | --- |
| スタイル | `Bold`, `Faint`, `Italic`, `Underline`, `BlinkSlow`, `BlinkRapid`, `ReverseVideo`, `Concealed`, `CrossedOut` |
| 前景色 | `FgBlack`〜`FgWhite`, `FgHiBlack`〜`FgHiWhite` |
| 背景色 | `BgBlack`〜`BgWhite`, `BgHiBlack`〜`BgHiWhite` |
| 256色 | `AddFg256(n)`, `AddBg256(n)` (0〜255) |
| トゥルーカラー | `AddFgRGB(r,g,b)`, `AddBgRGB(r,g,b)` |

## サンプル

`examples/demo/main.go` に一通りの使用例があります。

```bash
go run ./examples/demo
```

## ライセンス

MIT License. 詳細は [LICENSE](./LICENSE) を参照してください。
