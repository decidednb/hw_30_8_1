package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

type Task struct {
	ID         int    // идентификатор
	Opened     int64  // время Unix постановки
	Closed     int64  // время Unix закрытия
	AuthorID   int    // идентификатор автора
	AssignedID int    // идентификатор исполнителя
	Title      string // заголовок
	Content    string // содержание
}

// New - конструктор, conn - строка подключения к БД
func New(conn string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), conn)

	if err != nil {
		return nil, err
	}

	s := Storage{
		db: db,
	}

	return &s, nil
}

// ConnClose - закрывает соединение к базе данных
func (s *Storage) ConnClose() {
	s.db.Close()
}

// Tasks - список задач из базы данных
func (s *Storage) Tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT
			id,
			opened,
			closed,
			author_id,
			COALESCE(assigned_id, 0) AS assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2);
	`, taskID, authorID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, rows.Err()
}

// TasksByLabel - список задач по метке
func (s *Storage) TasksByLabel(labelID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT
			id,
			opened,
			closed,
			author_id,
			COALESCE(assigned_id, 0) AS assigned_id,
			title,
			content
		FROM tasks
		JOIN tasks_labels ON tasks_labels.task_id = tasks.id
		WHERE ($1 = 0 OR tasks_labels.label_id = $1);
	`, labelID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, rows.Err()
}

// NewTask - метод для создания задачи
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (author_id, title, content) VALUES
			($1, $2, $3) RETURNING id;
	`, t.AuthorID, t.Title, t.Content).Scan(&id)

	return id, err
}

// NewTasks - создание массива задач
func (s *Storage) NewTasks(tasks []Task) ([]int, error) {
	// массив идентификаторов, созданных задач
	var ids []int

	// начало транзакции
	ctx := context.Background()
	tx, err := s.db.Begin(ctx)

	if err != nil {
		return nil, err
	}

	// в случае ошибки - отмена транзакции
	defer tx.Rollback(ctx)

	// подготовка пакетного запроса
	batch := new(pgx.Batch)

	// добавление запросов в пакет
	for _, t := range tasks {
		batch.Queue(`
		INSERT INTO tasks(author_id, title, content)
		VALUES ($1, $2, $3) RETURNING id;
		`, t.AuthorID, t.Title, t.Content)
	}

	// отправка пакета запросов
	res := tx.SendBatch(ctx, batch)

	for i := 0; i < len(tasks); i++ {
		var id int
		err := res.QueryRow().Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	// закрытие соединения
	err = res.Close()
	if err != nil {
		return nil, err
	}

	// подтверждение транзакции
	return ids, tx.Commit(ctx)
}

// UpdateTask - обновление задачи по идентификатору
func (s *Storage) UpdateTask(t Task) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE tasks SET assigned_id = $1, title = $2, content = $3, closed = $4 WHERE id=$5;
	`, t.AssignedID, t.Title, t.Content, t.Closed, t.ID)

	return err
}

// DeleteTask - удаление задачи по идентификатору и удаление
// записей в таблице звязей task_labels
func (s *Storage) DeleteTask(id int) error {
	// начало транзакции
	ctx := context.Background()
	tx, err := s.db.Begin(ctx)

	if err != nil {
		return err
	}

	// в случае ошибки - отмена транзакции
	defer tx.Rollback(ctx)

	// подготовка пакетного запроса
	batch := new(pgx.Batch)

	// добавление в пакет запроса на удаление связей
	batch.Queue(`DELETE FROM tasks_labels WHERE task_id = $1;`, id)

	// добавление в пакет запроса на удаление задачи
	batch.Queue(`DELETE FROM tasks WHERE id = $1;`, id)

	// отправка пакета запросов
	res := tx.SendBatch(ctx, batch)

	// закрытие соединения
	err = res.Close()
	if err != nil {
		return err
	}

	// подтверждение транзакции
	return tx.Commit(ctx)
}
