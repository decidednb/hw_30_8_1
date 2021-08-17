-- Схема базы данных приложения "Задачник"

DROP TABLE IF EXISTS tasks_labels, labels, users, tasks;

-- Пользователи
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- Метки задач
CREATE TABLE labels (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- Задачи
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    opened BIGINT DEFAULT EXTRACT(epoch FROM clock_timestamp()),
    closed BIGINT DEFAULT 0,
    author_id INTEGER REFERENCES users(id) NOT NULL,
    assigned_id INTEGER REFERENCES users(id),
    title TEXT NOT NULL,
    content TEXT NOT NULL
);

-- Связи задачи - метки
CREATE TABLE tasks_labels (
    task_id INTEGER REFERENCES tasks(id),
    label_id INTEGER REFERENCES labels(id)
);

-- Очистка таблиц перед начальным заполнение БД
TRUNCATE TABLE tasks_labels, labels, users, tasks;

-- Начальное заполнение таблиц БД
INSERT INTO labels(name) VALUES
    ('Контент'),
    ('Каталог'),
    ('Сервис');

INSERT INTO users(name) VALUES
    ('Иванов Иван'),
    ('Петров Петр'),
    ('Сидоров Александр'),
    ('Жуков Степан'),
    ('Андреев Юрий');

INSERT INTO tasks(author_id, assigned_id, title, content) VALUES
    (1, 2, 'Добавить фотографии товаров', 'Подготовлены фотографии, необходимо добавить в каталог'),
    (1, 3, 'Добавить описание для товаров', 'Для новых товаров добавить описание'),
    (1, 4, 'Новый метод доставки', 'Добавить новый метод доставки курьерской службой СДЭК');

INSERT INTO tasks_labels(task_id, label_id) VALUES
    (1, 1), (1, 2), (2, 1), (2, 2), (3, 3);