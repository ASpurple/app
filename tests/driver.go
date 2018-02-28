package tests

import (
	"testing"

	"github.com/murlokswarm/app"
)

// TestDriver is a test suite that ensure that all driver implementations behave
// the same.
func TestDriver(t *testing.T, setup func(onRun func()) app.Driver, shutdown func()) {
	var driver app.Driver

	app.Import(&Hello{})
	app.Import(&World{})
	app.Import(&Menu{})
	app.Import(&Menubar{})

	onRun := func() {
		defer shutdown()

		t.Log("testing driver", driver.Name())
		t.Run("window", func(t *testing.T) { testWindow(t, driver) })
		t.Run("page", func(t *testing.T) { testPage(t, driver) })
		t.Run("context menu", func(t *testing.T) { testContextMenu(t, driver) })
		t.Run("dock", func(t *testing.T) { testDockTile(t, driver) })
	}

	driver = setup(onRun)

	if err := app.Run(driver); err != nil {
		t.Error(err)
	}
}
