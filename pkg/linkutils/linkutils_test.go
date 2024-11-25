package linkutils

import "testing"

func TestToAbsoluteURL(t *testing.T) {
	utils := NewLinkUtils()
	base := "https://example.com/path/page.html"

	tests := []struct {
		href     string
		expected string
	}{
		{"/about", "https://example.com/about"},
		{"contact.html", "https://example.com/path/contact.html"},
		{"https://example.com/services", "https://example.com/services"},
		{"//cdn.example.com/lib.js", "https://cdn.example.com/lib.js"},
	}

	for _, test := range tests {
		result := utils.ToAbsoluteURL(test.href, base)
		if result != test.expected {
			t.Errorf("ToAbsoluteURL(%s, %s) = %s; want %s", test.href, base, result, test.expected)
		}
	}
}

func TestIsInternalLink(t *testing.T) {
	utils := NewLinkUtils()
	base := "https://example.com/path/page.html"

	tests := []struct {
		link     string
		expected bool
	}{
		{"https://example.com/about", true},
		{"https://external.com/page", false},
		{"/contact", true},
		{"//example.com/assets", true},
	}

	for _, test := range tests {
		result := utils.IsInternalLink(test.link, base)
		if result != test.expected {
			t.Errorf("IsInternalLink(%s, %s) = %v; want %v", test.link, base, result, test.expected)
		}
	}
}
