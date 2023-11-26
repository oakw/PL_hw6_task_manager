package storage

// Task model representing a record in task_item table
// Id is auto-incremented
type Task struct {
	Id            int
	Title         string
	Description   string
	DueDate       string
	PriorityLevel int
	Completed     bool
}

// For Bubble Tea List filtering
func (i Task) FilterValue() string { return i.Title }
