package storage

var currentlyEditedTaskId int

func (dbCon Connection) IsCurrentlyEditingTask() bool {
	return currentlyEditedTaskId >= 0
}

func (dbCon Connection) GetCurrentlyEditedTaskId() int {
	return currentlyEditedTaskId
}

func (dbCon Connection) SetCurrentlyEditedTaskId(taskId interface{}) {
	if taskId != nil && taskId.(int) >= 0 {
		currentlyEditedTaskId = taskId.(int)
	} else {
		currentlyEditedTaskId = -1
	}
}
