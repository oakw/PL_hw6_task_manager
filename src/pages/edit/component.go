// This is a component that contains a form to either create a new task or edit an existing one.
// Based on
// https://github.com/charmbracelet/bubbletea/blob/master/examples/textinputs/main.go

package edit

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"task_manager/src/storage"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
)

// Styles of the various UI elements
var (
	appStyle       = lipgloss.NewStyle().Padding(1, 4)
	focusedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	errorTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	cursorStyle    = focusedStyle.Copy()
	noStyle        = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type Component struct {
	reactea.BasicComponent
	model editTaskModel
}

func New(store storage.Connection) *Component {
	return &Component{model: newEditTaskModel(store)}
}

func (c *Component) Init() tea.Cmd {
	return textinput.Blink
}

func (c *Component) Render(int, int) string {
	return appStyle.Render(c.model.View())
}

type editTaskModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	errorText  string
	store      storage.Connection
}

func newEditTaskModel(store storage.Connection) editTaskModel {
	model := editTaskModel{
		inputs:     make([]textinput.Model, 4),
		focusIndex: 0,
		store:      store,
	}

	// The page is opened to edit an existing task
	var existingTask storage.Task
	if store.IsCurrentlyEditingTask() {
		task, err := store.GetTaskById(store.GetCurrentlyEditedTaskId())
		if err != nil {
			log.Fatal(err)
		}
		existingTask = task
	}

	// Create the text inputs
	for i := range model.inputs {
		t := textinput.New()
		t.Cursor.Style = cursorStyle

		switch i {
		case 0:
			t.Prompt = "Task Title: "
			t.Placeholder = "My Task"
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.Focus()
			if store.IsCurrentlyEditingTask() {
				t.SetValue(existingTask.Title)
			}
		case 1:
			t.Prompt = "Description: "
			t.Placeholder = "Description"
			if store.IsCurrentlyEditingTask() {
				t.SetValue(existingTask.Description)
			}

		case 2:
			t.Prompt = "Due Date: "
			t.Placeholder = time.Now().Format("02.01.2006")
			if store.IsCurrentlyEditingTask() {
				tt, err := time.Parse("2006-01-02T15:04:05Z0700", existingTask.DueDate)
				if err != nil {
					log.Fatal(err)
				}
				t.SetValue(tt.Format("02.01.2006"))
			}

		case 3:
			t.Placeholder = "0"
			t.Prompt = "Priority (0 - 2): "
			if store.IsCurrentlyEditingTask() {
				t.SetValue(strconv.Itoa(existingTask.PriorityLevel))
			}
		}

		model.inputs[i] = t
	}

	return model
}

// Runs on every update
func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			reactea.SetCurrentRoute("task_list") // Move back to task list

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			c.model.handleNavigation(msg.String())
		}
	}

	// Handle character input and blinking
	cmd := c.model.updateInputs(msg)

	return cmd
}

// Custom form that is shown to the user
func (m editTaskModel) View() string {
	var b strings.Builder

	// Title
	title := "Create Task"
	if m.store.IsCurrentlyEditingTask() {
		title = "Edit Task"
	}
	fmt.Fprintf(&b, "%s\n\n",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1).
			Render(title))

	// Show the inputs
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)+1 {
			b.WriteRune('\n')
		}
	}

	// Show error message if there is one
	if m.errorText != "" {
		fmt.Fprintf(&b, "\n\n%s\n", errorTextStyle.Render(m.errorText))
	}

	// Show the submit button
	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
