package storage

import "hw_30_8_1/pkg/storage/postgres"

type Interface interface {
	// ConnClose - закрывает соединение к базе данных
	ConnClose()

	// Tasks - список задач из базы данных
	Tasks(taskID, authorID int) ([]postgres.Task, error)

	// TasksByLabel - список задач по идентификатору метки
	TasksByLabel(labelID int) ([]postgres.Task, error)

	// NewTask - создание задачи
	NewTask(postgres.Task) (int, error)

	// NewTasks - создание массива задач
	NewTasks([]postgres.Task) ([]int, error)

	// UpdateTask - обновление задачи по идентификатору
	UpdateTask(postgres.Task) error

	//DeleteTask - удаление задачи по идентификатору и удаление
	//записей в таблице звязей task_labels
	DeleteTask(id int) error
}
