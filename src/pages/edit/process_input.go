package edit

import (
	"log"
	"strconv"
	"task_manager/src/storage"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

func (m *editTaskModel) submitPressed() {
	// Check if all fields are filled
	for i := 0; i < len(m.inputs); i++ {
		if m.inputs[i].Value() == "" {
			m.errorText = "Please fill all fields"
			return
		}
	}

	// Check if priority is valid (0, 1 or 2)
	var priorityLevel int
	if m.inputs[3].Value() != "0" && m.inputs[3].Value() != "1" && m.inputs[3].Value() != "2" {
		m.errorText = "Priority must be 0, 1 or 2"
		return
	} else {
		priorityParsed, err := strconv.Atoi(m.inputs[3].Value())
		if err != nil {
			log.Fatal(err)
		}
		priorityLevel = priorityParsed
	}

	// Check if date is valid
	date, err := time.Parse("02.01.2006", m.inputs[2].Value())
	if err != nil {
		m.errorText = "Date must be in format dd.mm.yyyy"
		return
	}

	if m.store.IsCurrentlyEditingTask() {
		// Update the task
		task, err := m.store.GetTaskById(m.store.GetCurrentlyEditedTaskId())
		if err != nil {
			log.Fatal(err)
		}
		m.store.UpdateTask(storage.Task{
			Title:         m.inputs[0].Value(),
			Description:   m.inputs[1].Value(),
			DueDate:       date.Format("2006-01-02T15:04:05Z0700"),
			PriorityLevel: priorityLevel,
			Completed:     task.Completed,
		}, task.Id)

	} else {
		// Add a new task
		m.store.AddTask(
			storage.Task{
				Title:         m.inputs[0].Value(),
				Description:   m.inputs[1].Value(),
				DueDate:       date.Format("2006-01-02T15:04:05Z0700"),
				PriorityLevel: priorityLevel,
			},
		)
	}

	// Go back to the task list
	reactea.SetCurrentRoute("task_list")
}

// Copied from the example code
func (m *editTaskModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// Process the navigation, when the user presses tab shift+tab, enter, up or down
func (m *editTaskModel) handleNavigation(key string) tea.Cmd {
	// Did the user press enter while the submit button was focused?
	// If so, exit.
	if key == "enter" && m.focusIndex == len(m.inputs) {
		m.submitPressed()
	}

	// Cycle indexes
	if key == "up" || key == "shift+tab" {
		m.focusIndex--
	} else {
		m.focusIndex = m.focusIndex + 1
	}

	// Clamp the index to the size of the inputs slice
	if m.focusIndex > len(m.inputs) {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs)
	}

	var cmd tea.Cmd
	for i := 0; i < len(m.inputs); i++ {
		if i == m.focusIndex {
			// Set focused state
			cmd = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
			continue
		}
		// Remove focused state
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = noStyle
		m.inputs[i].TextStyle = noStyle
	}

	return cmd
}
