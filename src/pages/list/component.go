// This is a component that displays a list of tasks.
// It is the main page of the application.
// Based on
// https://github.com/charmbracelet/bubbletea/blob/master/examples/list-simple/main.go
package list

import (
	"log"
	"task_manager/src/storage"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
)

// Styles for the list
var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).Render

	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
)

type Component struct {
	reactea.BasicComponent
	model taskListModel
}

func New(store storage.Connection) *Component {
	return &Component{
		model: newTaskListModel(store),
	}
}

func (c *Component) Init() tea.Cmd {
	c.model.store.SetCurrentlyEditedTaskId(nil) // Reset the currently edited task id
	return tea.EnterAltScreen
}

// Return the string displayed in the command-line interface
func (c *Component) Render(int, int) string {
	return appStyle.Render(c.model.list.View())
}

// Strucutre binded to the component and operated as per Bubble Tea library
type taskListModel struct {
	list         list.Model
	delegateKeys *listDelegate
	store        storage.Connection
}

func newTaskListModel(store storage.Connection) taskListModel {
	// Get all tasks from the store
	// Put them in the list
	allTasks, err := store.GetAllTasks()
	if err != nil {
		log.Fatal(err)
	}

	items := make([]list.Item, len(allTasks))
	for i, task := range allTasks {
		items[i] = task
	}

	delegate := newlistDelegate()
	taskList := list.New(items, addDelegateFunctionality(delegate), 0, 15)
	taskList.Title = "To Do List"
	taskList.Styles.Title = titleStyle

	// Show statistics in the status bar using a bit of a hack
	if len(allTasks) != 0 {
		statistics := getStatisticsString(allTasks)
		taskList.SetStatusBarItemName("item\n"+statistics, "items\n"+statistics)
	}

	return taskListModel{
		list:         taskList,
		delegateKeys: delegate,
		store:        store,
	}
}

// Runs on every update
func (c *Component) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		c.model.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if c.model.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, c.model.delegateKeys.markAsDone):
			if c.model.list.SelectedItem() != nil {
				cmds = append(cmds, toggleTaskDone(c))
			}

		case key.Matches(msg, c.model.delegateKeys.remove):
			if c.model.list.SelectedItem() != nil {
				cmds = append(cmds, performTaskDelete(c))
			}

		case key.Matches(msg, c.model.delegateKeys.edit):
			if c.model.list.SelectedItem() != nil {
				task, ok := c.model.list.SelectedItem().(storage.Task)
				if !ok {
					log.Fatal("Not a task")
				}

				c.model.store.SetCurrentlyEditedTaskId(task.Id)
				reactea.SetCurrentRoute("edit_view")
			}
		case key.Matches(msg, c.model.delegateKeys.add):
			reactea.SetCurrentRoute("edit_view")
		}

	}

	// This will also call our delegate's update function.
	newListModel, cmd := c.model.list.Update(msg)
	c.model.list = newListModel
	cmds = append(cmds, cmd)

	// Similarly as in initial setup function
	// Update statistics for the status bar
	allItems := c.model.list.Items()
	allTasks := make([]storage.Task, len(allItems))
	for i, item := range allItems {
		allTasks[i] = item.(storage.Task)
	}
	if len(allTasks) != 0 {
		statistics := getStatisticsString(allTasks)
		c.model.list.SetStatusBarItemName("item\n"+statistics, "items\n"+statistics)
	}

	return tea.Batch(cmds...)
}
