package main

import (
	"encoding/json"
	"fmt"
)

// TodoStatus represents the state of a todo item.
type TodoStatus string

const (
	TodoStatusPending    TodoStatus = "pending"
	TodoStatusInProgress TodoStatus = "in_progress"
	TodoStatusCompleted  TodoStatus = "completed"
)

// Todo represents a single task.
type Todo struct {
	ActiveForm string     `json:"activeForm"`
	Content    string     `json:"content"`
	Status     TodoStatus `json:"status"`
}

// TodoList represents the root structure containing todos.
type TodoList struct {
	Todos []Todo `json:"todos"`
}

// ParseTodoFromMap parses a todo list from a map[string]any structure.
// Expected input format:
//
//	map[string]any{
//		"todos": []any{
//			map[string]any{"activeForm": "...", "content": "...", "status": "..."},
//			...
//		},
//	}
func ParseTodoFromMap(m map[string]any) (*TodoList, error) {
	todosAny, ok := m["todos"]
	if !ok {
		return nil, fmt.Errorf("missing 'todos' key in map")
	}

	todosSlice, ok := todosAny.([]any)
	if !ok {
		return nil, fmt.Errorf("'todos' is not a slice: %T", todosAny)
	}

	var todos []Todo
	for i, item := range todosSlice {
		todoMap, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("todo item %d is not a map: %T", i, item)
		}

		todo := Todo{}
		if v, ok := todoMap["activeForm"].(string); ok {
			todo.ActiveForm = v
		}
		if v, ok := todoMap["content"].(string); ok {
			todo.Content = v
		}
		if v, ok := todoMap["status"].(string); ok {
			todo.Status = TodoStatus(v)
		}

		todos = append(todos, todo)
	}

	return &TodoList{Todos: todos}, nil
}

// MarshalJSON implements json.Marshaler.
func (t TodoStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(t))
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *TodoStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = TodoStatus(s)
	return nil
}
