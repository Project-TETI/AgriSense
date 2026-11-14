package catalogseed

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadValidation(t *testing.T) {
	valid := `{"name":"tomato","display_name":"Tomat","diseases":[{"class_name":"spot","display_name":"Bercak","cause":"x","description":"x","recommendation":"x"}]}`
	for _, tc := range []struct {
		name, content string
		duplicate, ok bool
	}{
		{"valid", valid, false, true},
		{"invalid json", "{", false, false},
		{"unknown field", strings.Replace(valid, `"name":`, `"unknown":1,"name":`, 1), false, false},
		{"empty name", strings.Replace(valid, `"name":"tomato"`, `"name":""`, 1), false, false},
		{"invalid class", strings.Replace(valid, `"class_name":"spot"`, `"class_name":"Bad-Class"`, 1), false, false},
		{"empty recommendation", strings.Replace(valid, `"recommendation":"x"`, `"recommendation":"  "`, 1), false, false},
		{"duplicate class", strings.Replace(valid, `}]}`, `},{"class_name":"spot","display_name":"Bercak","cause":"x","description":"x","recommendation":"x"}]}`, 1), false, false},
		{"duplicate plant", valid, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "one.json"), []byte(tc.content), 0600); err != nil {
				t.Fatal(err)
			}
			if tc.duplicate {
				if err := os.WriteFile(filepath.Join(dir, "two.json"), []byte(tc.content), 0600); err != nil {
					t.Fatal(err)
				}
			}
			_, err := Load(dir)
			if (err == nil) != tc.ok {
				t.Fatalf("Load error: %v", err)
			}
		})
	}
	if _, err := Load(t.TempDir()); err == nil {
		t.Fatal("empty directory accepted")
	}
}

func TestTomatoCatalogMatchesModel(t *testing.T) {
	plants, err := Load("../../data/catalog")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("../../../ml/weights/tomato/class_names.json")
	if err != nil {
		t.Fatal(err)
	}
	var labels []string
	if err := json.Unmarshal(data, &labels); err != nil {
		t.Fatal(err)
	}
	var tomato *Plant
	for i := range plants {
		if plants[i].Name == "tomato" {
			tomato = &plants[i]
			break
		}
	}
	if tomato == nil {
		t.Fatal("tomato missing")
	}
	want := make(map[string]bool)
	for _, label := range labels {
		if label != "healthy" {
			want[label] = true
		}
	}
	for _, disease := range tomato.Diseases {
		if !want[disease.ClassName] {
			t.Fatalf("unexpected class %s", disease.ClassName)
		}
		delete(want, disease.ClassName)
	}
	if len(want) != 0 {
		t.Fatalf("missing classes: %v", want)
	}
}
