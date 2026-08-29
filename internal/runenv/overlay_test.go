package runenv

import (
	"reflect"
	"runtime"
	"testing"
)

func TestOverlayApplyRemovesAmbientSelectorsAndSetsProfile(t *testing.T) {
	base := []string{
		"PATH=/usr/bin",
		"GH_TOKEN=ambient-token",
		"GH_HOST=wrong.example.com",
		"GH_CONFIG_DIR=/ambient/gh",
		"KEEP=value",
	}
	overlay := Overlay{
		Set: map[string]string{"GH_CONFIG_DIR": "/profiles/personal"},
		Unset: []string{
			"GH_TOKEN",
			"GITHUB_TOKEN",
			"GH_HOST",
			"GH_REPO",
		},
	}

	got := overlay.Apply(base)
	want := []string{
		"PATH=/usr/bin",
		"KEEP=value",
		"GH_CONFIG_DIR=/profiles/personal",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Apply() = %#v, want %#v", got, want)
	}
}

func TestOverlayApplyEdgeCases(t *testing.T) {
	t.Run("empty overlay preserves base", func(t *testing.T) {
		base := []string{"B=2", "A=1"}
		if got := (Overlay{}).Apply(base); !reflect.DeepEqual(got, base) {
			t.Fatalf("Apply() = %#v, want %#v", got, base)
		}
	})

	t.Run("nil base uses process environment", func(t *testing.T) {
		t.Setenv("NO_MISTAKES_OVERLAY_TEST", "ambient")
		got := (Overlay{Set: map[string]string{"NO_MISTAKES_OVERLAY_TEST": "scoped"}}).Apply(nil)
		if !containsEnvironmentEntry(got, "NO_MISTAKES_OVERLAY_TEST=scoped") {
			t.Fatalf("Apply(nil) did not replace the ambient value: %#v", got)
		}
	})

	t.Run("set wins over unset and keys are sorted", func(t *testing.T) {
		overlay := Overlay{
			Set:   map[string]string{"ZED": "last", "BOTH": "set", "ALPHA": "first"},
			Unset: []string{"BOTH", "REMOVE"},
		}
		got := overlay.Apply([]string{"REMOVE=old", "BOTH=old", "KEEP=yes"})
		want := []string{"KEEP=yes", "ALPHA=first", "BOTH=set", "ZED=last"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Apply() = %#v, want %#v", got, want)
		}
	})

	t.Run("key matching follows platform semantics", func(t *testing.T) {
		got := (Overlay{Unset: []string{"mixed_key"}}).Apply([]string{"MIXED_KEY=value"})
		wantRemoved := runtime.GOOS == "windows"
		if removed := len(got) == 0; removed != wantRemoved {
			t.Fatalf("removed = %v, want %v on %s; result %#v", removed, wantRemoved, runtime.GOOS, got)
		}
	})
}

func TestOverlayEmptyAndClone(t *testing.T) {
	if !(Overlay{}).Empty() {
		t.Fatal("zero overlay is not empty")
	}
	if (Overlay{Set: map[string]string{"A": "1"}}).Empty() {
		t.Fatal("overlay with Set entry reported empty")
	}
	if (Overlay{Unset: []string{"A"}}).Empty() {
		t.Fatal("overlay with Unset entry reported empty")
	}

	original := Overlay{Set: map[string]string{"A": "1"}, Unset: []string{"B"}}
	cloned := original.Clone()
	cloned.Set["A"] = "changed"
	cloned.Unset[0] = "changed"
	if original.Set["A"] != "1" || original.Unset[0] != "B" {
		t.Fatalf("Clone shares mutable state with original: %#v", original)
	}
}

func containsEnvironmentEntry(environment []string, want string) bool {
	for _, entry := range environment {
		if entry == want {
			return true
		}
	}
	return false
}
