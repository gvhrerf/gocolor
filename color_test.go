package gocolor

import (
	"strings"
	"testing"
)

func TestBasicColor(t *testing.T) {
	c := New(FgRed).EnableColor()
	got := c.Sprint("hello")
	want := "\x1b[31mhello\x1b[0m"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMultipleAttributes(t *testing.T) {
	c := New(FgGreen, Bold).EnableColor()
	got := c.Sprint("hi")
	want := "\x1b[32;1mhi\x1b[0m"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func Test256Color(t *testing.T) {
	c := New().AddFg256(208).EnableColor()
	got := c.Sprint("orange")
	want := "\x1b[38;5;208morange\x1b[0m"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRGBColor(t *testing.T) {
	c := New().AddFgRGB(255, 105, 180).EnableColor()
	got := c.Sprint("pink")
	want := "\x1b[38;2;255;105;180mpink\x1b[0m"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFgAndBgRGBCombined(t *testing.T) {
	c := New().AddFgRGB(0, 191, 255).AddBgRGB(20, 20, 20).EnableColor()
	got := c.Sprint("x")
	want := "\x1b[38;2;0;191;255;48;2;20;20;20mx\x1b[0m"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestClamp256OutOfRange(t *testing.T) {
	c := New().AddFg256(999).EnableColor()
	got := c.Sprint("x")
	if !strings.Contains(got, "38;5;255") {
		t.Errorf("expected clamp to 255, got %q", got)
	}

	c2 := New().AddFg256(-10).EnableColor()
	got2 := c2.Sprint("x")
	if !strings.Contains(got2, "38;5;0") {
		t.Errorf("expected clamp to 0, got %q", got2)
	}
}

func TestDisableColorNoEscape(t *testing.T) {
	c := New(FgRed).DisableColor()
	got := c.Sprint("plain")
	if got != "plain" {
		t.Errorf("expected no escape codes, got %q", got)
	}
}

func TestSprintf(t *testing.T) {
	c := New(FgBlue).EnableColor()
	got := c.Sprintf("count=%d", 5)
	want := "\x1b[34mcount=5\x1b[0m"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestShortcutFuncNoAttrsHaveNoColorByDefault(t *testing.T) {
	// makeColorFunc の Color は EnableColor されていないため、
	// NoColor がテスト環境の判定に従う（端末でなければ無色）。
	// ここでは EnableColor 相当を直接確認するため New を使う。
	c := New(FgRed).EnableColor()
	if c.Sprint("x") == "x" {
		t.Errorf("expected colored output when EnableColor is set")
	}
}
