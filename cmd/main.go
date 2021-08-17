package main

import (
	"fmt"
	"hw_30_8_1/pkg/storage"
	"hw_30_8_1/pkg/storage/postgres"
	"log"
	"os"
	"time"
)

// db - интерфейс базы данных
var db storage.Interface

func main() {
	var tasks []postgres.Task
	var err error

	dbpass := os.Getenv("taskspass")
	if dbpass == "" {
		os.Exit(1)
	}

	// conn - строка подключения к базе данных
	conn := "postgres://postgres:" + dbpass + "@localhost:5432/tasks"

	// присвоение переменной реализации базы данных
	db, err = postgres.New(conn)
	defer db.ConnClose()

	if err != nil {
		log.Fatal(err)
	}

	// получение задачи по метке
	tasks, err = db.TasksByLabel(3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasks)

	// создание задачи
	t := postgres.Task{
		AuthorID: 1,
		Title:    "Наполнение категорий",
		Content:  "Проставить связи товар - категории для новых товаров",
	}

	id, err := db.NewTask(t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Создана задача №%d\n", id)

	// создание массива задач
	tasks = nil
	t = postgres.Task{
		AuthorID: 5,
		Title:    "Наполнение категорий",
		Content:  "Проставить связи товар - категории для новых товаров",
	}
	tasks = append(tasks, t)

	t = postgres.Task{
		AuthorID: 3,
		Title:    "Добавить производителя",
		Content:  "Добавить в каталог товаров нового производителя 'Jackson'",
	}
	tasks = append(tasks, t)

	ids, err := db.NewTasks(tasks)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Созданы задачи №%d\n", ids)

	// удаление задачи по номеру
	err = db.DeleteTask(1)
	if err != nil {
		log.Fatal(err)
	}

	// получение задачи по номеру и автору
	tasks, err = db.Tasks(3, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasks)

	// обновление задачи
	t = postgres.Task{
		ID:         3,
		AssignedID: 3,
		Title:      "Добавить производителя",
		Content:    "Добавить в каталог товаров нового производителя 'Jackson'",
		Closed:     time.Now().Unix(),
	}

	err = db.UpdateTask(t)
	if err != nil {
		log.Fatal(err)
	}

	// получение всех задач из БД
	tasks, err = db.Tasks(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasks)
}
