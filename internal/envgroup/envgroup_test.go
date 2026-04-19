package envgroup

import (
	"strings"
	"testing"
)

func TestGroup_ByPrefix(t *testing.T) {
	g := New("_")
	secrets := map[string]string{
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"APP_NAME": "vaultline",
		"PLAIN":    "value",
	}
	groups := g.Group(secrets)
	index := map[string]Group{}
	for _, gr := range groups {
		index[gr.Name] = gr
	}
	if len(index["db"].Secrets) != 2 {
		t.Errorf("expected 2 db secrets, got %d", len(index["db"].Secrets))
	}
	if len(index["app"].Secrets) != 1 {
		t.Errorf("expected 1 app secret, got %d", len(index["app"].Secrets))
	}
	if len(index["default"].Secrets) != 1 {
		t.Errorf("expected 1 default secret, got %d", len(index["default"].Secrets))
	}
}

func TestGroup_EmptySecrets(t *testing.T) {
	g := New("_")
	groups := g.Group(map[string]string{})
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}

func TestGroup_NoDelimiter(t *testing.T) {
	g := New("_")
	secrets := map[string]string{"PLAIN": "val"}
	groups := g.Group(secrets)
	if groups[0].Name != "default" {
		t.Errorf("expected default group, got %s", groups[0].Name)
	}
}

func TestGroup_SortedOutput(t *testing.T) {
	g := New("_")
	secrets := map[string]string{
		"Z_KEY": "1",
		"A_KEY": "2",
		"M_KEY": "3",
	}
	groups := g.Group(secrets)
	if groups[0].Name != "a" || groups[1].Name != "m" || groups[2].Name != "z" {
		t.Errorf("groups not sorted: %v", groups)
	}
}

func TestSummary_Format(t *testing.T) {
	groups := []Group{
		{Name: "db", Secrets: map[string]string{"DB_HOST": "x", "DB_PORT": "y"}},
		{Name: "app", Secrets: map[string]string{"APP_NAME": "z"}},
	}
	out := Summary(groups)
	if !strings.Contains(out, "[db] 2 key(s)") {
		t.Errorf("unexpected summary: %s", out)
	}
	if !strings.Contains(out, "[app] 1 key(s)") {
		t.Errorf("unexpected summary: %s", out)
	}
}
