// Global application component.
// Handles the routing and the global state of the application.
// Similar to https://github.com/londek/reactea/blob/v0.4.2/examples/dynamicRoutes/app/app.go
package app

import (
	"log"
	edit_page "task_manager/src/pages/edit"
	list_page "task_manager/src/pages/list"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"

	"task_manager/src/storage"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	store storage.Connection
}

func New() *Component {
	store, err := storage.CreateConnection()
	if err != nil {
		log.Fatal(err)
	}

	return &Component{
		mainRouter: router.New(),
		store:      store,
	}
}

func (c *Component) Render(width, height int) string {
	return c.mainRouter.Render(width, height)
}

// Return the router map
func (c *Component) Init(reactea.NoProps) tea.Cmd {
	task_list_view := func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := list_page.New(c.store)
		return component, component.Init()
	}

	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default":   task_list_view,
		"task_list": task_list_view,
		"edit_view": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := edit_page.New(c.store)
			return component, component.Init()
		},
	})
}

// Global update function
func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	return c.mainRouter.Update(msg)
}
