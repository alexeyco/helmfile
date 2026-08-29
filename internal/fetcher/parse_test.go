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

func TestParseNodeTags(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    string
		wantErr bool
	}{
		{"happy path", `{"results":[{"name":"26.8.1-alpine","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]}]}`, "26.8.1", false},
		{"semver max beats list order", `{"results":[{"name":"24.20.0-alpine","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]},{"name":"26.8.1-alpine","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]}]}`, "26.8.1", false},
		{"numeric minor compare", `{"results":[{"name":"26.9.3-alpine","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]},{"name":"26.10.0-alpine","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]}]}`, "26.10.0", false},
		{"skips non plain alpine tags", `{"results":[{"name":"26.8.1-alpine3.24","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]},{"name":"26.8.1-bookworm","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]},{"name":"26.8.1-slim","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]},{"name":"26.8.1","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]},{"name":"latest-alpine","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]},{"name":"27.0.0-rc.1-alpine","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]},{"name":"26.8.1-alpine","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]}]}`, "26.8.1", false},
		{"ignores non linux and unknown images", `{"results":[{"name":"26.8.1-alpine","images":[{"os":"windows","architecture":"amd64"},{"architecture":"arm64"},{"os":"linux","architecture":"riscv64"},{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]}]}`, "26.8.1", false},
		{"skips amd64 only tag", `{"results":[{"name":"26.8.1-alpine","images":[{"os":"linux","architecture":"amd64"}]},{"name":"26.8.0-alpine","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]}]}`, "26.8.0", false},
		{"amd64 only candidate", `{"results":[{"name":"26.8.1-alpine","images":[{"os":"linux","architecture":"amd64"}]}]}`, "", true},
		{"null images", `{"results":[{"name":"26.8.1-alpine","images":null}]}`, "", true},
		{"no qualifying tags", `{"results":[{"name":"26.8.1-alpine3.24","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]},{"name":"26.8.1-bookworm","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]},{"name":"26.8.1-slim","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]},{"name":"26.8.1","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]},{"name":"latest-alpine","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]},{"name":"27.0.0-rc.1-alpine","images":[{"os":"linux","architecture":"amd64"},{"os":"linux","architecture":"arm64"}]}]}`, "", true},
		{"empty results", `{"results":[]}`, "", true},
		{"invalid JSON", `{"results": `, "", true},
		{"empty body", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNodeTags([]byte(tt.body))
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
