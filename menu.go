package app

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Menu is the interface that describes a menu.
// Accept only components that contain menu and menuitem tags.
type Menu interface {
	ElementWithComponent

	// Base returns the base menu without any decorators.
	Base() Menu
}

// MenuConfig is a struct that describes a menu.
type MenuConfig struct {
	DefaultURL string

	OnClose func() `json:"-"`
}

// MenuWithLogs returns a decorated version of the given menu that logs
// all the operations.
// Uses the default logger.
func MenuWithLogs(m Menu, name string) Menu {
	return &menuWithLogs{
		name: name,
		base: m,
	}
}

type menuWithLogs struct {
	name string
	base Menu
}

func (m *menuWithLogs) ID() uuid.UUID {
	return m.base.ID()
}

func (m *menuWithLogs) Base() Menu {
	return m.base.Base()
}

func (m *menuWithLogs) Load(url string, v ...interface{}) error {
	fmtURL := fmt.Sprintf(url, v...)
	Logf("%s %s: loading %s", m.name, m.base.ID(), fmtURL)

	err := m.base.Load(url, v...)
	if err != nil {
		Errorf("%s %s: loading %s failed: %s", m.name, m.base.ID(), fmtURL, err)
	}
	return err
}

func (m *menuWithLogs) Component() Component {
	c := m.base.Component()
	Logf("%s %s: mounted component is %T", m.name, m.base.ID(), c)
	return c
}

func (m *menuWithLogs) Contains(c Component) bool {
	ok := m.base.Contains(c)
	Logf("%s %s: contains %T is %v", m.name, m.base.ID(), c, ok)
	return ok
}

func (m *menuWithLogs) Render(c Component) error {
	Logf("%s %s: rendering %T", m.name, m.base.ID(), c)

	err := m.base.Render(c)
	if err != nil {
		Errorf("%s %s: rendering %T failed: %s", m.name, m.base.ID(), c, err)
	}
	return err
}

func (m *menuWithLogs) LastFocus() time.Time {
	return m.base.LastFocus()
}
