package markdown

import (
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "見出しレベル1",
			input:    "# Hello World",
			expected: `<h1 id="hello-world">Hello World</h1>`,
		},
		{
			name:     "見出しレベル2",
			input:    "## Subtitle",
			expected: `<h2 id="subtitle">Subtitle</h2>`,
		},
		{
			name:     "段落",
			input:    "This is a paragraph.",
			expected: "<p>This is a paragraph.</p>",
		},
		{
			name:     "太字",
			input:    "This is **bold** text.",
			expected: "<p>This is <strong>bold</strong> text.</p>",
		},
		{
			name:     "イタリック",
			input:    "This is *italic* text.",
			expected: "<p>This is <em>italic</em> text.</p>",
		},
		{
			name:     "リスト",
			input:    "- Item 1\n- Item 2\n- Item 3",
			expected: "<ul>\n<li>Item 1</li>\n<li>Item 2</li>\n<li>Item 3</li>\n</ul>",
		},
		{
			name:     "コードブロック",
			input:    "```go\nfmt.Println(\"Hello\")\n```",
			expected: `<pre><code class="language-go">fmt.Println(&quot;Hello&quot;)
</code></pre>`,
		},
		{
			name:     "インラインコード",
			input:    "Use `fmt.Printf()` function.",
			expected: "<p>Use <code>fmt.Printf()</code> function.</p>",
		},
		{
			name:     "リンク",
			input:    "[Google](https://www.google.com)",
			expected: `<p><a href="https://www.google.com">Google</a></p>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Convert(tt.input)
			// HTMLの改行や空白を正規化
			result = strings.TrimSpace(result)
			expected := strings.TrimSpace(tt.expected)
			
			if result != expected {
				t.Errorf("Convert() = %v, want %v", result, expected)
			}
		})
	}
}

func TestConvertEmpty(t *testing.T) {
	result := Convert("")
	if result != "" {
		t.Errorf("Convert(\"\") = %v, want empty string", result)
	}
}

func TestConvertComplexDocument(t *testing.T) {
	input := `# Title

This is a paragraph with **bold** and *italic* text.

## Subtitle

Here's a list:
- First item
- Second item
- Third item

And some code:

` + "```go" + `
func main() {
    fmt.Println("Hello, World!")
}
` + "```"

	result := Convert(input)
	
	// 必要な要素が含まれているかチェック
	if !strings.Contains(result, `<h1 id="title">Title</h1>`) {
		t.Error("Missing h1 title")
	}
	if !strings.Contains(result, `<h2 id="subtitle">Subtitle</h2>`) {
		t.Error("Missing h2 subtitle")
	}
	if !strings.Contains(result, "<strong>bold</strong>") {
		t.Error("Missing bold text")
	}
	if !strings.Contains(result, "<em>italic</em>") {
		t.Error("Missing italic text")
	}
	if !strings.Contains(result, "<ul>") || !strings.Contains(result, "<li>") {
		t.Error("Missing list elements")
	}
	if !strings.Contains(result, `<code class="language-go">`) {
		t.Error("Missing code block with language")
	}
}