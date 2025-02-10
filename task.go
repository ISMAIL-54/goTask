package main

import "time"

type Task struct {
    ID              uint        `json:"id"`
    Title           string      `json:"title"`
    Description     string      `json:"description"`
    Completed       bool        `json:"completed"`
    Timestamp       time.Time   `json:"create_at"`
}

func NewTask(title string, description string, completed bool) *Task {
    return &Task{
        Title: title,
        Description: description,
        Completed: completed,
        Timestamp: time.Now().UTC(),
    }    
}
