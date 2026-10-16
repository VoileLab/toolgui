package tccontent

import (
	"testing"

	"github.com/voilelab/toolgui/toolgui/tgcomp/internal/tcutil"
	"github.com/voilelab/toolgui/toolgui/tgframe"
)

func TestTextProps(t *testing.T) {
	cases := []struct {
		name string
		text string
	}{
		{"plain", "hello"},
		{"empty", ""},
		{"unicode", "你好世界"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			props := addMetric(t, func(c *tgframe.Container) {
				Text(c, tc.text)
			})
			if props["text"] != tc.text {
				t.Errorf("text = %v, want %v", props["text"], tc.text)
			}
			if props["name"] != "text_component" {
				t.Errorf("name = %v, want text_component", props["name"])
			}
		})
	}
}

func TestTitleProps(t *testing.T) {
	cases := []struct {
		name string
		text string
	}{
		{"plain", "My Title"},
		{"empty", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			props := addMetric(t, func(c *tgframe.Container) {
				Title(c, tc.text)
			})
			if props["text"] != tc.text {
				t.Errorf("text = %v, want %v", props["text"], tc.text)
			}
			if props["name"] != "title_component" {
				t.Errorf("name = %v, want title_component", props["name"])
			}
		})
	}
}

func TestSubtitleProps(t *testing.T) {
	props := addMetric(t, func(c *tgframe.Container) {
		Subtitle(c, "a subtitle")
	})
	if props["text"] != "a subtitle" {
		t.Errorf("text = %v, want a subtitle", props["text"])
	}
	if props["name"] != "subtitle_component" {
		t.Errorf("name = %v, want subtitle_component", props["name"])
	}
}

func TestCaptionProps(t *testing.T) {
	props := addMetric(t, func(c *tgframe.Container) {
		Caption(c, "a caption")
	})
	if props["text"] != "a caption" {
		t.Errorf("text = %v, want a caption", props["text"])
	}
	if props["name"] != "caption_component" {
		t.Errorf("name = %v, want caption_component", props["name"])
	}
}

func TestLinkProps(t *testing.T) {
	cases := []struct {
		name string
		text string
		url  string
	}{
		{"basic", "docs", "https://example.com"},
		{"empty text", "", "https://example.com"},
		{"empty url", "link", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			props := addMetric(t, func(c *tgframe.Container) {
				Link(c, tc.text, tc.url)
			})
			if props["text"] != tc.text {
				t.Errorf("text = %v, want %v", props["text"], tc.text)
			}
			if props["url"] != tc.url {
				t.Errorf("url = %v, want %v", props["url"], tc.url)
			}
		})
	}
}

func TestCodeProps(t *testing.T) {
	cases := []struct {
		name string
		code string
		lang string
		want string
	}{
		{"default lang", "fmt.Println()", "", "go"},
		{"explicit lang", "console.log()", "javascript", "javascript"},
		{"empty code", "", "", "go"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			props := addMetric(t, func(c *tgframe.Container) {
				Code(c, tc.code, &CodeConf{Language: tc.lang})
			})
			if props["code"] != tc.code {
				t.Errorf("code = %v, want %v", props["code"], tc.code)
			}
			if props["lang"] != tc.want {
				t.Errorf("lang = %v, want %v", props["lang"], tc.want)
			}
		})
	}
}

func TestBadgeProps(t *testing.T) {
	props := addMetric(t, func(c *tgframe.Container) {
		Badge(c, "active", &BadgeConf{Color: tcutil.ColorSuccess})
	})
	if props["text"] != "active" {
		t.Errorf("text = %v, want active", props["text"])
	}
	if props["color"] != "success" {
		t.Errorf("color = %v, want success", props["color"])
	}
}

func TestDividerProps(t *testing.T) {
	props := addMetric(t, func(c *tgframe.Container) {
		Divider(c)
	})
	if props["name"] != "divider_component" {
		t.Errorf("name = %v, want divider_component", props["name"])
	}
}
