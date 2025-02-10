package main

import (
    "fmt"
    "database/sql"
    _ "github.com/lib/pq"
)

type Storage interface {
	CreateTask(*Task) error
	DeleteTask(int) error
	UpdateTask(*Task) error
	GetTaskByID(int) (*Task, error)
	GetTasks() ([]*Task, error)
}

type PostgresDB struct {
    db *sql.DB
}

func NewPostgresDB() (*PostgresDB, error) {
    connStr := "user=postgres password=1234 dbname=todo_db sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    return &PostgresDB{
        db: db,
    }, nil
}

func (p *PostgresDB) Init() error {
    return p.CreateTaskTable() 
}

func (p *PostgresDB) CreateTaskTable() error {
    query := `
        CREATE TABLE IF NOT EXISTS task(
            id              SERIAL PRIMARY KEY,
            title           VARCHAR(50) NOT NULL,
            description     TEXT,
            completed       BOOLEAN,
            timestamp       TIMESTAMP
        )
    `
    _, err := p.db.Exec(query)
    return err
}

func (p *PostgresDB) CreateTask(task *Task) error {
    query := `INSERT INTO task(title, description, completed, timestamp) VALUES($1, $2, $3, $4)`
    res, err := p.db.Query(query, task.Title, task.Description, task.Completed, task.Timestamp)
    if err != nil {
        return err
    }
    fmt.Println(res)
    return nil
}

func (p *PostgresDB) GetTasks() ([]*Task, error) {
    rows, err := p.db.Query("SELECT * From task")
    if err != nil {
        return nil, err
    }

    tasks := []*Task{}
    for rows.Next() {
        task := new(Task)
        err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.Timestamp)
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, task)
    }
    return tasks, nil
}

func (p *PostgresDB) GetTaskByID(id int) (*Task, error) {
    rows, err := p.db.Query("SELECT * FROM task WHERE id = $1", id)
    if err != nil {
        return nil, err
    }

    task := new(Task)
    for rows.Next() {
        return scanIntoTask(rows)
    }
    return task, fmt.Errorf("task %d not found", id)
}

func (p *PostgresDB) DeleteTask(id int) error {
    _, err := p.db.Query("DELETE FROM task WHERE id = $1", id)
    return err
}

func (p *PostgresDB) UpdateTask(task *Task) error {
    return nil
}

func scanIntoTask(rows *sql.Rows) (*Task, error){
    task := new(Task)
    err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.Timestamp)
    if err != nil {
        return nil, err
    }
    return task, nil
}
