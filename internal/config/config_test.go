package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandHome(t *testing.T) {
	tests := []struct {
		name string
		path string
		home string
		want string
	}{
		{
			name: "expands tilde prefix",
			path: "~/Projects/foo",
			home: "/home/user",
			want: "/home/user/Projects/foo",
		},
		{
			name: "no tilde returns unchanged",
			path: "/absolute/path",
			home: "/home/user",
			want: "/absolute/path",
		},
		{
			name: "tilde only in middle is not expanded",
			path: "/some/~/path",
			home: "/home/user",
			want: "/some/~/path",
		},
		{
			name: "empty path returns empty",
			path: "",
			home: "/home/user",
			want: "",
		},
		{
			name: "empty home with tilde",
			path: "~/foo",
			home: "",
			want: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandHome(tt.path, tt.home)
			if got != tt.want {
				t.Errorf("expandHome(%q, %q) = %q, want %q", tt.path, tt.home, got, tt.want)
			}
		})
	}
}

func TestFindProject(t *testing.T) {
	cfg := &Config{
		Projects: []Project{
			{Name: "alpha", Path: "/a"},
			{Name: "beta", Path: "/b"},
			{Name: "gamma", Path: "/c"},
		},
	}

	tests := []struct {
		name      string
		search    string
		wantPath  string
		wantFound bool
	}{
		{name: "finds first project", search: "alpha", wantPath: "/a", wantFound: true},
		{name: "finds last project", search: "gamma", wantPath: "/c", wantFound: true},
		{name: "returns nil for missing", search: "delta", wantFound: false},
		{name: "empty name not found", search: "", wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.FindProject(tt.search)
			if tt.wantFound {
				if got == nil {
					t.Fatalf("FindProject(%q) = nil, want project", tt.search)
				}
				if got.Path != tt.wantPath {
					t.Errorf("FindProject(%q).Path = %q, want %q", tt.search, got.Path, tt.wantPath)
				}
			} else {
				if got != nil {
					t.Errorf("FindProject(%q) = %+v, want nil", tt.search, got)
				}
			}
		})
	}
}

func TestGetLayout(t *testing.T) {
	devLayout := Layout{
		Windows: []Window{{Name: "editor"}, {Name: "server"}},
	}
	cfg := &Config{
		Layouts: map[string]Layout{"dev": devLayout},
		Projects: []Project{
			{Name: "with-layout", Layout: "dev"},
			{Name: "no-layout", Layout: ""},
			{Name: "bad-layout", Layout: "nonexistent"},
		},
	}

	tests := []struct {
		name        string
		projectName string
		wantWindows int
		wantFirst   string
	}{
		{name: "uses named layout", projectName: "with-layout", wantWindows: 2, wantFirst: "editor"},
		{name: "falls back to default", projectName: "no-layout", wantWindows: 1, wantFirst: "main"},
		{name: "missing layout falls back", projectName: "bad-layout", wantWindows: 1, wantFirst: "main"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proj := cfg.FindProject(tt.projectName)
			if proj == nil {
				t.Fatalf("project %q not found", tt.projectName)
			}
			layout := cfg.GetLayout(proj)
			if len(layout.Windows) != tt.wantWindows {
				t.Errorf("GetLayout() windows = %d, want %d", len(layout.Windows), tt.wantWindows)
			}
			if layout.Windows[0].Name != tt.wantFirst {
				t.Errorf("GetLayout() first window = %q, want %q", layout.Windows[0].Name, tt.wantFirst)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		content := `
projects:
  - name: test
    path: /tmp/test
    layout: dev
layouts:
  dev:
    windows:
      - name: editor
        panes:
          - size: "100%"
`
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := Load(path)
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if len(cfg.Projects) != 1 {
			t.Errorf("Load() projects = %d, want 1", len(cfg.Projects))
		}
		if cfg.Projects[0].Name != "test" {
			t.Errorf("Load() project name = %q, want %q", cfg.Projects[0].Name, "test")
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := Load("/nonexistent/path/config.yaml")
		if err == nil {
			t.Error("Load() expected error for missing file, got nil")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		if err := os.WriteFile(path, []byte("{{invalid yaml"), 0644); err != nil {
			t.Fatal(err)
		}

		_, err := Load(path)
		if err == nil {
			t.Error("Load() expected error for invalid YAML, got nil")
		}
	})
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()
	if path == "" {
		t.Error("DefaultConfigPath() returned empty string")
	}
	// Should end with the expected config file name.
	if filepath.Base(path) != ConfigFile {
		t.Errorf("DefaultConfigPath() base = %q, want %q", filepath.Base(path), ConfigFile)
	}
}
