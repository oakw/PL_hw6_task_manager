package list

import (
	"fmt"
	"io"
	"log"
	"strings"
	"task_manager/src/storage"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// Delegate of the list for Bubble Tea library
// Extends the default one
type listDelegate struct {
	list.DefaultDelegate
	markAsDone key.Binding
	remove     key.Binding
	edit       key.Binding
	add        key.Binding
}

func newlistDelegate() *listDelegate {
	return &listDelegate{
		markAsDone: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "toggle done"),
		),
		remove: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "delete"),
		),
		edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		)}
}

// Add key bindings, so, when rendered, the list view would have the functionality
func addDelegateFunctionality(delegate *listDelegate) *listDelegate {
	help := []key.Binding{delegate.markAsDone, delegate.add, delegate.edit, delegate.remove}

	// Shows functions at the bottom of the screen
	delegate.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{delegate.markAsDone, delegate.add, delegate.edit, delegate.remove}
	}

	// Shows functions when pressing ?
	delegate.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return delegate
}

// Render one list item
//
// Example:
//
//	'[X] Task Title
//	10.01.2023.
//	Your Description'
func (d listDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(storage.Task)
	if !ok {
		return
	}

	// Show X as tick for completed tasks
	checkboxStr := "[ ]"
	if i.Completed {
		checkboxStr = "[X]"
	}

	tt, err := time.Parse("2006-01-02T15:04:05Z0700", i.DueDate)
	if err != nil {
		log.Fatal(err)
	}
	dateFormatted := tt.Format("02.01.2006.")
	str := checkboxStr + fmt.Sprintf(" %s\n    %s\n    %s", i.Title, dateFormatted, i.Description)

	// Change the color of the task based on its priority level
	switch i.PriorityLevel {
	case 0:
		itemStyle.Foreground(lipgloss.Color("#00FF00"))
	case 1:
		itemStyle.Foreground(lipgloss.Color("#FFFF00"))
	case 2:
		itemStyle.Foreground(lipgloss.Color("#FF0000"))
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return itemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
