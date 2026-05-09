package uistate

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kqnade/VRCLaunch/internal/config"
)

type View int

const (
	ViewMain View = iota
	ViewEditProfile
	ViewSettings
)

type ProfileForm struct {
	ID               string // empty = new
	Name             string
	Index            string
	FPS              string
	ScreenWidth      string
	ScreenHeight     string
	ScreenFullscreen bool
	CustomArgs       string
}

type SettingsForm struct {
	LaunchPath string
}

type State struct {
	Config       *config.Config
	View         View
	SelectedID   string
	ProfileForm  ProfileForm
	SettingsForm SettingsForm
	Status       string
}

func NewState(cfg *config.Config) *State {
	if cfg == nil {
		cfg = config.Default()
	}
	return &State{
		Config:     cfg,
		View:       ViewMain,
		SelectedID: cfg.LastSelected,
	}
}

func (s *State) GotoMain() {
	s.View = ViewMain
}

func (s *State) GotoNewProfile() {
	s.ProfileForm = ProfileForm{
		ID:    "",
		Index: strconv.Itoa(s.Config.NextProfileIndex()),
	}
	s.View = ViewEditProfile
}

var ErrProfileNotFound = errors.New("profile not found")

func (s *State) GotoEditProfile(id string) error {
	p := s.Config.FindByID(id)
	if p == nil {
		return ErrProfileNotFound
	}
	s.ProfileForm = ProfileForm{
		ID:               p.ID,
		Name:             p.Name,
		Index:            strconv.Itoa(p.Index),
		FPS:              optionalIntString(p.Options.FPS),
		ScreenWidth:      optionalIntString(p.Options.ScreenWidth),
		ScreenHeight:     optionalIntString(p.Options.ScreenHeight),
		ScreenFullscreen: p.Options.ScreenFullscreen,
		CustomArgs:       p.Options.CustomArgs,
	}
	s.View = ViewEditProfile
	return nil
}

func (s *State) GotoSettings() {
	s.SettingsForm = SettingsForm{LaunchPath: s.Config.LaunchPath}
	s.View = ViewSettings
}

func (s *State) SaveProfileFromForm() error {
	form := s.ProfileForm

	if form.Name == "" {
		return errors.New("name is required")
	}

	idx, err := strconv.Atoi(form.Index)
	if err != nil || idx <= 0 {
		return errors.New("index must be a positive integer")
	}

	for _, p := range s.Config.Profiles {
		if p.ID != form.ID && p.Index == idx {
			return fmt.Errorf("index %d already used by profile %q", idx, p.Name)
		}
	}

	fps, err := optionalAtoi(form.FPS)
	if err != nil {
		return fmt.Errorf("fps: %w", err)
	}
	w, err := optionalAtoi(form.ScreenWidth)
	if err != nil {
		return fmt.Errorf("screen width: %w", err)
	}
	h, err := optionalAtoi(form.ScreenHeight)
	if err != nil {
		return fmt.Errorf("screen height: %w", err)
	}

	options := config.ProfileOptions{
		FPS:              fps,
		ScreenWidth:      w,
		ScreenHeight:     h,
		ScreenFullscreen: form.ScreenFullscreen,
		CustomArgs:       form.CustomArgs,
	}

	if form.ID == "" {
		s.Config.Profiles = append(s.Config.Profiles, config.Profile{
			ID:      config.NewID(),
			Name:    form.Name,
			Index:   idx,
			Options: options,
		})
	} else {
		p := s.Config.FindByID(form.ID)
		if p == nil {
			return ErrProfileNotFound
		}
		p.Name = form.Name
		p.Index = idx
		p.Options = options
	}

	s.View = ViewMain
	return nil
}

func (s *State) DeleteProfile(id string) error {
	for i, p := range s.Config.Profiles {
		if p.ID == id {
			s.Config.Profiles = append(s.Config.Profiles[:i], s.Config.Profiles[i+1:]...)
			if s.SelectedID == id {
				s.SelectedID = ""
				s.Config.LastSelected = ""
			}
			return nil
		}
	}
	return ErrProfileNotFound
}

func (s *State) SelectProfile(id string) {
	s.SelectedID = id
	s.Config.LastSelected = id
}

func (s *State) SaveSettingsFromForm() {
	s.Config.LaunchPath = s.SettingsForm.LaunchPath
	s.View = ViewMain
}

func optionalIntString(v int) string {
	if v == 0 {
		return ""
	}
	return strconv.Itoa(v)
}

func optionalAtoi(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q", s)
	}
	if n < 0 {
		return 0, fmt.Errorf("must be non-negative: %d", n)
	}
	return n, nil
}
