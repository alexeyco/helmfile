package fetcher

import "testing"

func TestParsePlainVersion(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    string
		wantErr bool
	}{
		{"with v prefix", "v1.36.4\n", "1.36.4", false},
		{"with uppercase V prefix", "V1.36.4\n", "1.36.4", false},
		{"without prefix", "1.36.4\n", "1.36.4", false},
		{"surrounding whitespace", "  v1.36.4  \n", "1.36.4", false},
		{"empty body", "", "", true},
		{"whitespace only", " \n\t", "", true},
		{"latest", "latest\n", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePlainVersion([]byte(tt.body))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseTagName(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    string
		wantErr bool
	}{
		{"with v prefix", `{"tag_name": "v3.21.4"}`, "3.21.4", false},
		{"without prefix", `{"tag_name": "3.21.4"}`, "3.21.4", false},
		{"whitespace in tag", `{"tag_name": "  v3.21.4\n"}`, "3.21.4", false},
		{"extra fields", `{"url": "x", "tag_name": "v1.2.3", "name": "release"}`, "1.2.3", false},
		{"missing tag_name", `{"name": "release"}`, "", true},
		{"empty tag_name", `{"tag_name": ""}`, "", true},
		{"latest tag", `{"tag_name": "latest"}`, "", true},
		{"empty body", "", "", true},
		{"invalid JSON", `{"tag_name": `, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTagName([]byte(tt.body))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

const alpineIndexHTML = `<html>
<head><title>Index of /alpine/</title></head>
<body>
<h1>Index of /alpine/</h1><hr><pre><a href="../">../</a>
<a href="edge/">edge/</a>
<a href="latest-stable/">latest-stable/</a>
<a href="v3.19/">v3.19/</a>
<a href="v3.20/">v3.20/</a>
<a href="v3.21/">v3.21/</a>
<a href="v3.22/">v3.22/</a>
<a href="v3.23/">v3.23/</a>
<a href="v3.24/">v3.24/</a>
<a href="v3.9/">v3.9/</a>
</pre><hr></body>
</html>`

func TestParseAlpineIndex(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    string
		wantErr bool
	}{
		{"numeric max wins over lexicographic", alpineIndexHTML, "3.24", false},
		{"no releases", `<html><body>nope</body></html>`, "", true},
		{"empty body", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAlpineIndex([]byte(tt.body))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
