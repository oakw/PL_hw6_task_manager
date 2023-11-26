package storage

import (
	"database/sql"
	"fmt"
)

// Insert task into task_item table
func (dbCon Connection) AddTask(task Task) error {
	_, err := dbCon.executeStatement(
		"INSERT INTO task_item(Title, Description, DueDate, PriorityLevel, Completed) values(?, ?, ?, ?, ?);",
		task.Title, task.Description, task.DueDate, task.PriorityLevel, task.Completed)

	return err
}

// Retrieve task from task_item table by id
func (dbCon Connection) GetTaskById(taskId int) (Task, error) {
	result, err := dbCon.query("SELECT * FROM task_item WHERE Id = ? LIMIT 1;", taskId)
	if err != nil {
		return Task{}, err
	}

	tasks, err := getTasksFromResult(*result)
	if err != nil {
		return Task{}, err
	}

	if len(tasks) < 1 {
		return Task{}, fmt.Errorf("No task found")

	} else {
		return tasks[0], nil
	}
}

// Retrieve tasks as Task models from sql.Rows object
func getTasksFromResult(result sql.Rows) ([]Task, error) {
	tasks := []Task{}

	err := result.Err()
	if err != nil {
		return nil, err
	}

	for result.Next() {
		var task = Task{}
		err = result.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.PriorityLevel,
			&task.Completed,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Delete task from task_item table by id
func (dbCon Connection) DeleteTaskById(taskId int) error {
	_, err := dbCon.executeStatement(
		"DELETE FROM task_item WHERE Id = ?;", taskId)

	return err
}

// Retrieve client instance to perform operations with
func (dbCon Connection) UpdateTask(task Task, taskId int) error {
	_, err := dbCon.executeStatement(
		"UPDATE task_item SET Title = ?, Description = ?, DueDate = ?, PriorityLevel = ?, Completed = ? WHERE Id = ?;",
		task.Title, task.Description, task.DueDate, task.PriorityLevel, task.Completed, taskId)

	return err
}

// Retrieve all tasks from task_item table
func (dbCon Connection) GetAllTasks() ([]Task, error) {
	result, err := dbCon.query("SELECT * FROM task_item;")
	if err != nil {
		return nil, err
	}

	tasks, err := getTasksFromResult(*result)

	return tasks, err
}
