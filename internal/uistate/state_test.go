package uistate

import (
	"errors"
	"testing"

	"github.com/kqnade/VRCLaunch/internal/config"
)

func newConfigWith(profiles ...config.Profile) *config.Config {
	c := config.Default()
	c.Profiles = append(c.Profiles, profiles...)
	return c
}

func TestNewState_StartsAtMainView(t *testing.T) {
	s := NewState(config.Default())
	if s.View != ViewMain {
		t.Errorf("View = %v, want ViewMain", s.View)
	}
}

func TestNewState_NilConfigSafe(t *testing.T) {
	s := NewState(nil)
	if s.Config == nil {
		t.Fatal("Config should be non-nil even when initialized with nil")
	}
	if s.View != ViewMain {
		t.Errorf("View = %v, want ViewMain", s.View)
	}
}

func TestNewState_RestoresLastSelected(t *testing.T) {
	c := newConfigWith(config.Profile{ID: "x", Name: "X", Index: 1})
	c.LastSelected = "x"
	s := NewState(c)
	if s.SelectedID != "x" {
		t.Errorf("SelectedID = %q, want %q", s.SelectedID, "x")
	}
}

func TestGotoMain(t *testing.T) {
	s := NewState(config.Default())
	s.View = ViewSettings
	s.GotoMain()
	if s.View != ViewMain {
		t.Errorf("View = %v, want ViewMain", s.View)
	}
}

func TestGotoNewProfile_SetsFormAndSuggestsNextIndex(t *testing.T) {
	c := newConfigWith(
		config.Profile{ID: "a", Name: "A", Index: 1},
		config.Profile{ID: "b", Name: "B", Index: 2},
	)
	s := NewState(c)
	s.GotoNewProfile()

	if s.View != ViewEditProfile {
		t.Errorf("View = %v, want ViewEditProfile", s.View)
	}
	if s.ProfileForm.ID != "" {
		t.Errorf("expected empty ID for new profile, got %q", s.ProfileForm.ID)
	}
	if s.ProfileForm.Index != "3" {
		t.Errorf("Index = %q, want %q", s.ProfileForm.Index, "3")
	}
}

func TestGotoEditProfile_PopulatesForm(t *testing.T) {
	p := config.Profile{
		ID:    "x",
		Name:  "Main",
		Index: 1,
		Options: config.ProfileOptions{
			FPS:              90,
			ScreenWidth:      1280,
			ScreenHeight:     720,
			ScreenFullscreen: true,
			CustomArgs:       "--foo",
		},
	}
	s := NewState(newConfigWith(p))

	if err := s.GotoEditProfile("x"); err != nil {
		t.Fatalf("GotoEditProfile: %v", err)
	}
	if s.View != ViewEditProfile {
		t.Errorf("View = %v, want ViewEditProfile", s.View)
	}
	want := ProfileForm{
		ID:               "x",
		Name:             "Main",
		Index:            "1",
		FPS:              "90",
		ScreenWidth:      "1280",
		ScreenHeight:     "720",
		ScreenFullscreen: true,
		CustomArgs:       "--foo",
	}
	if s.ProfileForm != want {
		t.Errorf("ProfileForm:\n got %+v\nwant %+v", s.ProfileForm, want)
	}
}

func TestGotoEditProfile_UnknownIDReturnsError(t *testing.T) {
	s := NewState(config.Default())
	err := s.GotoEditProfile("missing")
	if !errors.Is(err, ErrProfileNotFound) {
		t.Errorf("got %v, want ErrProfileNotFound", err)
	}
}

func TestGotoEditProfile_OmitsZeroIntsAsEmptyStrings(t *testing.T) {
	p := config.Profile{ID: "x", Name: "X", Index: 1}
	s := NewState(newConfigWith(p))
	if err := s.GotoEditProfile("x"); err != nil {
		t.Fatal(err)
	}
	if s.ProfileForm.FPS != "" || s.ProfileForm.ScreenWidth != "" || s.ProfileForm.ScreenHeight != "" {
		t.Errorf("zero options should map to empty strings: %+v", s.ProfileForm)
	}
}

func TestGotoSettings_PopulatesForm(t *testing.T) {
	c := config.Default()
	c.LaunchPath = "/x/launch.exe"
	s := NewState(c)
	s.GotoSettings()

	if s.View != ViewSettings {
		t.Errorf("View = %v, want ViewSettings", s.View)
	}
	if s.SettingsForm.LaunchPath != "/x/launch.exe" {
		t.Errorf("LaunchPath = %q, want %q", s.SettingsForm.LaunchPath, "/x/launch.exe")
	}
}

func TestSaveProfile_EmptyNameRejected(t *testing.T) {
	s := NewState(config.Default())
	s.ProfileForm = ProfileForm{Name: "", Index: "1"}
	if err := s.SaveProfileFromForm(); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestSaveProfile_InvalidIndexRejected(t *testing.T) {
	cases := []string{"", "0", "-1", "abc"}
	for _, idx := range cases {
		t.Run(idx, func(t *testing.T) {
			s := NewState(config.Default())
			s.ProfileForm = ProfileForm{Name: "X", Index: idx}
			if err := s.SaveProfileFromForm(); err == nil {
				t.Errorf("expected error for index %q", idx)
			}
		})
	}
}

func TestSaveProfile_DuplicateIndexRejected(t *testing.T) {
	c := newConfigWith(config.Profile{ID: "a", Name: "A", Index: 1})
	s := NewState(c)
	s.ProfileForm = ProfileForm{Name: "B", Index: "1"}
	if err := s.SaveProfileFromForm(); err == nil {
		t.Error("expected error for duplicate index")
	}
}

func TestSaveProfile_DuplicateIndexAllowedForSameProfile(t *testing.T) {
	c := newConfigWith(config.Profile{ID: "a", Name: "A", Index: 1})
	s := NewState(c)
	s.ProfileForm = ProfileForm{ID: "a", Name: "A renamed", Index: "1"}
	if err := s.SaveProfileFromForm(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.Profiles[0].Name != "A renamed" {
		t.Errorf("Name = %q, want %q", c.Profiles[0].Name, "A renamed")
	}
}

func TestSaveProfile_NewProfileGetsIDAndAppended(t *testing.T) {
	c := config.Default()
	s := NewState(c)
	s.ProfileForm = ProfileForm{Name: "Main", Index: "1"}
	if err := s.SaveProfileFromForm(); err != nil {
		t.Fatal(err)
	}
	if len(c.Profiles) != 1 {
		t.Fatalf("Profiles len = %d, want 1", len(c.Profiles))
	}
	if c.Profiles[0].ID == "" {
		t.Error("new profile should be assigned an ID")
	}
	if s.View != ViewMain {
		t.Errorf("View = %v, want ViewMain after save", s.View)
	}
}

func TestSaveProfile_InvalidNumericFieldsRejected(t *testing.T) {
	cases := []ProfileForm{
		{Name: "X", Index: "1", FPS: "abc"},
		{Name: "X", Index: "1", ScreenWidth: "-1"},
		{Name: "X", Index: "1", ScreenHeight: "abc"},
	}
	for i, f := range cases {
		t.Run("", func(t *testing.T) {
			s := NewState(config.Default())
			s.ProfileForm = f
			if err := s.SaveProfileFromForm(); err == nil {
				t.Errorf("case %d: expected error, got nil for %+v", i, f)
			}
		})
	}
}

func TestSaveProfile_UpdateMissingIDReturnsError(t *testing.T) {
	s := NewState(config.Default())
	s.ProfileForm = ProfileForm{ID: "ghost", Name: "X", Index: "1"}
	if err := s.SaveProfileFromForm(); !errors.Is(err, ErrProfileNotFound) {
		t.Errorf("got %v, want ErrProfileNotFound", err)
	}
}

func TestSaveProfile_PersistsAllOptions(t *testing.T) {
	c := config.Default()
	s := NewState(c)
	s.ProfileForm = ProfileForm{
		Name:             "X",
		Index:            "5",
		FPS:              "120",
		ScreenWidth:      "2560",
		ScreenHeight:     "1440",
		ScreenFullscreen: true,
		CustomArgs:       "--foo",
	}
	if err := s.SaveProfileFromForm(); err != nil {
		t.Fatal(err)
	}
	got := c.Profiles[0].Options
	want := config.ProfileOptions{
		FPS: 120, ScreenWidth: 2560, ScreenHeight: 1440,
		ScreenFullscreen: true, CustomArgs: "--foo",
	}
	if got != want {
		t.Errorf("Options:\n got %+v\nwant %+v", got, want)
	}
}

func TestDeleteProfile_RemovesAndClearsSelection(t *testing.T) {
	c := newConfigWith(
		config.Profile{ID: "a", Name: "A", Index: 1},
		config.Profile{ID: "b", Name: "B", Index: 2},
	)
	c.LastSelected = "a"
	s := NewState(c)
	s.SelectedID = "a"

	if err := s.DeleteProfile("a"); err != nil {
		t.Fatal(err)
	}
	if len(c.Profiles) != 1 || c.Profiles[0].ID != "b" {
		t.Errorf("Profiles after delete: %+v", c.Profiles)
	}
	if s.SelectedID != "" {
		t.Errorf("SelectedID = %q, want empty", s.SelectedID)
	}
	if c.LastSelected != "" {
		t.Errorf("LastSelected = %q, want empty", c.LastSelected)
	}
}

func TestDeleteProfile_KeepsSelectionWhenDifferent(t *testing.T) {
	c := newConfigWith(
		config.Profile{ID: "a", Name: "A", Index: 1},
		config.Profile{ID: "b", Name: "B", Index: 2},
	)
	s := NewState(c)
	s.SelectedID = "b"

	if err := s.DeleteProfile("a"); err != nil {
		t.Fatal(err)
	}
	if s.SelectedID != "b" {
		t.Errorf("SelectedID should remain %q, got %q", "b", s.SelectedID)
	}
}

func TestDeleteProfile_UnknownIDReturnsError(t *testing.T) {
	s := NewState(config.Default())
	if err := s.DeleteProfile("ghost"); !errors.Is(err, ErrProfileNotFound) {
		t.Errorf("got %v, want ErrProfileNotFound", err)
	}
}

func TestSelectProfile_UpdatesLastSelected(t *testing.T) {
	c := newConfigWith(config.Profile{ID: "x", Name: "X", Index: 1})
	s := NewState(c)
	s.SelectProfile("x")
	if s.SelectedID != "x" || c.LastSelected != "x" {
		t.Errorf("SelectedID=%q LastSelected=%q", s.SelectedID, c.LastSelected)
	}
}

func TestSaveSettings_AppliesLaunchPath(t *testing.T) {
	c := config.Default()
	s := NewState(c)
	s.SettingsForm.LaunchPath = "/games/launch.exe"
	s.SaveSettingsFromForm()
	if c.LaunchPath != "/games/launch.exe" {
		t.Errorf("LaunchPath = %q, want %q", c.LaunchPath, "/games/launch.exe")
	}
	if s.View != ViewMain {
		t.Errorf("View = %v, want ViewMain after save", s.View)
	}
}
