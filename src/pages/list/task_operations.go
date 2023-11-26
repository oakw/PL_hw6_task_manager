package list

import (
	"fmt"
	"log"
	"task_manager/src/storage"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

/**
 * Toggle task done. Change its state in DB and in list shown
 * Return a command to display the status message
 */
func toggleTaskDone(c *Component) tea.Cmd {
	task, ok := c.model.list.SelectedItem().(storage.Task)
	if !ok {
		log.Fatal("Not a task")
	}

	task.Completed = !task.Completed
	err := c.model.store.UpdateTask(task, task.Id)
	if err != nil {
		log.Fatal(err)
	}
	c.model.list.SetItem(c.model.list.Index(), task)

	taskStatus := "not done"
	if task.Completed {
		taskStatus = "done"
	}

	return c.model.list.NewStatusMessage(statusMessageStyle(fmt.Sprintf("%s is %s now", task.Title, taskStatus)))
}

/**
 * Delete task from DB and from list shown
 * Return a command to display the status message
 */
func performTaskDelete(c *Component) tea.Cmd {
	task, ok := c.model.list.SelectedItem().(storage.Task)
	if !ok {
		log.Fatal("Not a task")
	}

	err := c.model.store.DeleteTaskById(task.Id)
	if err != nil {
		log.Fatal(err)
	}
	c.model.list.RemoveItem(c.model.list.Index())

	return c.model.list.NewStatusMessage(statusMessageStyle(fmt.Sprintf("Deleted %s", task.Title)))
}

func getStatisticsString(allTasks []storage.Task) string {
	highPriorityTasksCount := 0
	completedTasksCount := 0
	dueInWeekCount := 0
	lateCount := 0

	for _, task := range allTasks {
		if (task.PriorityLevel == 2) && !task.Completed {
			highPriorityTasksCount++
		}
		if task.Completed {
			completedTasksCount++
		}

		tt, err := time.Parse("2006-01-02T15:04:05Z0700", task.DueDate)
		if err != nil {
			log.Fatal(err)
		}
		if !task.Completed && tt.Compare(time.Now().Add(time.Hour*24*7)) == -1 {
			dueInWeekCount++
		}
		year, month, day := time.Now().Date()
		if !task.Completed && tt.Compare(time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())) == -1 {
			lateCount++
		}
	}

	return fmt.Sprintf(
		"High priority: %d, Completed: %d, Due in week: %d, Late: %d",
		highPriorityTasksCount,
		completedTasksCount,
		dueInWeekCount,
		lateCount,
	)
}
