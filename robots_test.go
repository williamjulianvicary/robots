package robots

import (
	"os"
	"testing"
)

func TestAgentPrecedence(t *testing.T) {
	fname := "testdata/agent_precedence.txt"
	data, err := os.Open(fname)
	if err != nil {
		t.Errorf("couldn't open test data %s", fname)
	}

	var tests = []struct {
		input string
		want  string
	}{
		{"Googlebot-News", "googlebot-news"},
		{"Googlebot", "googlebot"},
		{"Googlebot-Image", "googlebot"},
		{"Bingbot", "*"},
	}

	r, err := From(200, data)
	if err != nil {
		t.Errorf("couldn't read from test data %s", fname)
	}

	for _, test := range tests {
		if got, _ := r.bestAgent(test.input); got.name != test.want {
			t.Errorf("r.bestAgent(%q).name = %v", test.input, got.name)
		}
	}
}

func TestGrouping(t *testing.T) {
	fname := "testdata/grouping.txt"
	data, err := os.Open(fname)
	if err != nil {
		t.Errorf("couldn't open test data %s", fname)
	}

	var tests = []struct {
		name string
		path string
		want bool
	}{
		{"a", "/c", false},
		{"a", "/d", true},
		{"a", "/g", true},
		{"b", "/c", true},
		{"b", "/d", false},
		{"b", "/q", true},
		{"e", "/c", true},
		{"e", "/d", true},
		{"e", "/g", false},
		{"f", "/c", true},
		{"f", "/d", true},
		{"f", "/g", false},
	}

	r, err := From(200, data)
	if err != nil {
		t.Errorf("couldn't read from test data %s", fname)
	}

	for _, test := range tests {
		if got := r.Test(test.name, test.path); got != test.want {
			t.Errorf("r.Test(%q, %q) = %t", test.name, test.path, got)
		}
	}
}

func TestPathMatching(t *testing.T) {
	var tests = []struct {
		source string
		test   string
		want   bool
	}{
		{"/", "/", true},
		{"/", "/lower/level", true},
		{"/*", "/", true},
		{"/*", "/lower/level", true},
		{"/fish", "/fish", true},
		{"/fish", "/fish.html", true},
		{"/fish", "/fish/salmon.html", true},
		{"/fish", "/fishheads", true},
		{"/fish", "/fishheads/yummy.html", true},
		{"/fish", "/fish.php?id=anything", true},
		{"/fish", "/Fish.asp", false},
		{"/fish", "/catfish", false},
		{"/fish", "/?id=fish", false},
		{"/fish*", "/fish", true},
		{"/fish*", "/fish.html", true},
		{"/fish*", "/fish/salmon.html", true},
		{"/fish*", "/fishheads", true},
		{"/fish*", "/fishheads/yummy.html", true},
		{"/fish*", "/fish.php?id=anything", true},
		{"/fish*", "/Fish.asp", false},
		{"/fish*", "/catfish", false},
		{"/fish*", "/?id=fish", false},
		{"/fish/", "/fish/", true},
		{"/fish/", "/fish/?id=anything", true},
		{"/fish/", "/fish/salmon.htm", true},
		{"/fish/", "/fish", false},
		{"/fish/", "/fish.html", false},
		{"/fish/", "/Fish/Salmon.asp", false},
		{"/*.php", "/filename.php", true},
		{"/*.php", "/folder/filename.php", true},
		{"/*.php", "/folder/filename.php?parameters", true},
		{"/*.php", "/folder/any.php.file.html", true},
		{"/*.php", "/filename.php/", true},
		{"/*.php", "/", false},
		{"/*.php", "/windows.PHP", false},
		{"/*.php$", "/filename.php", true},
		{"/*.php$", "/folder/filename.php", true},
		{"/*.php$", "/folder/filename.php?parameters", false},
		{"/*.php$", "/filename.php/", false},
		{"/*.php$", "/filename.php5", false},
		{"/*.php$", "/windows.PHP", false},
		{"/fish*.php", "/fish.php", true},
		{"/fish*.php", "/fishheads/catfish.php?parameters", true},
		{"/fish*.php", "/Fish.PHP", false},
	}

	for _, test := range tests {
		m := &member{
			path: test.source,
		}
		m.compile()
		if got := m.match(test.test); got != test.want {
			t.Errorf(
				"for pattern %q, m.match(%q) = %t",
				test.source,
				test.test,
				got,
			)
		}
	}
}

func TestMemberPrecedence(t *testing.T) {
	fname := "testdata/agent_precedence.txt"
	data, err := os.Open(fname)
	if err != nil {
		t.Errorf("couldn't open test data %s", fname)
	}

	var tests = []struct {
		input string
		want  bool
	}{
		{"/page", true},
		// Under the specification I think this should
		// also be undefined.
		{"/folder/page", true},
		{"/", true},
	}

	r, err := From(200, data)
	if err != nil {
		t.Errorf("couldn't read from test data %s", fname)
	}

	for _, test := range tests {
		if got := r.Test("crawler", test.input); got != test.want {
			t.Errorf("r.Test(\"crawler\", %q) = %t", test.input, got)
		}
	}
}

func TestLocate(t *testing.T) {
	var tests = []struct {
		robots string
		path   string
		want   bool
	}{
		{"http://example.com/robots.txt", "http://example.com/", true},
		{"http://example.com/robots.txt", "http://example.com/folder/file", true},
		{"http://example.com/robots.txt", "http://other.example.com/", false},
		{"http://example.com/robots.txt", "https://example.com/", false},
		{"http://example.com/robots.txt", "http://example.com:8181/", false},
		{"http://www.example.com/robots.txt", "http://www.example.com/", true},
		{"http://www.example.com/robots.txt", "http://example.com/", false},
		{"http://www.example.com/robots.txt", "http://shop.www.example.com/", false},
		{"http://www.example.com/robots.txt", "http://www.shop.example.com/", false},
		{"http://www.müller.eu/robots.txt", "http://www.müller.eu/", true},
		{"http://www.müller.eu/robots.txt", "http://www.xn--mller-kva.eu/", true},
		{"http://www.müller.eu/robots.txt", "http://www.muller.eu/", false},
		{"ftp://example.com/robots.txt", "ftp://example.com/", true},
		{"ftp://example.com/robots.txt", "http://example.com/", false},
		{"http://example.com:8181/robots.txt", "http://example.com:8181/", true},
		{"http://example.com:8181/robots.txt", "http://example.com/", false},
	}

	for _, test := range tests {
		rpath, _ := Locate(test.path)
		if test.want {
			if got := rpath == test.robots; !got {
				t.Errorf("Locate(%q) should be %v", test.path, test.robots)
			}
		} else {
			if got := rpath == test.robots; got {
				t.Errorf("Locate(%q) should not be %v", test.path, test.robots)
			}
		}
	}
}

func TestLocateCase(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"hTtP://eXamPle.cOm:1234/somefile.html", "http://example.com:1234/robots.txt"},
	}

	for _, test := range tests {
		if got, _ := Locate(test.input); got != test.want {
			t.Errorf("Locate(%q) should be %v, got %v", test.input, test.want, got)
		}
	}
}

func TestPathological(t *testing.T) {
	f, err := os.Open("testdata/pathological.txt")
	if err != nil {
		t.Errorf("%v", err)
	}
	r, _ := From(200, f)
	if r.Test("Googlebot", "/exact-match") {
		t.Errorf("googlebot shouldn't crawl /exact-match")
	}
	seen := false
	for _, s := range r.Sitemaps() {
		if s == "http://www.example.com/sitemap.xml" {
			seen = true
		}
	}
	if !seen {
		t.Errorf("didn't discover http://www.example.com/sitemap.xml")
	}
}

func TestRobotsAllow(t *testing.T) {
	var tests = []struct {
		status int
		want   bool
	}{
		{200, true},
		{301, true},
		{404, true},
		{503, false},
	}

	for _, test := range tests {
		r, err := From(test.status, nil)
		if err != nil {
			t.Errorf("couldn't make Robots with status %d", test.status)
		}
		if got := r.Test("Crawlerbot", "/somepath"); got != test.want {
			t.Errorf("%d status should default to %t, got %t",
				test.status, test.want, got)
		}
	}
}
