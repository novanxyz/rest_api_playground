package models

import "github.com/szmcdull/glinq/garray"

type TaskFile struct {
	Id         uint   `gorm:"type:int;primary_key"`
	Filename   string `gorm:"type:varchar(255)"`
	Mime       string `gorm:"type:varchar(32)"`
	Content    []byte `gorm:"type:longblob"`
	ParentTask *Task  `gorm:"ForeignKey:Id"`
}

type Task struct {
	Id     uint        `gorm:"type:int;primary_key"`
	Name   string      `gorm:"type:varchar(255)"`
	Status string      `gorm:"type:varchar(10) default 'incomplete'"`
	Files  []*TaskFile `gorm:"ForeignKey:ParentTask"`
}

func (task *Task) ToResponse() TaskResponse {
	return TaskResponse{
		Id:        task.Id,
		Name:      task.Name,
		Status:    task.Name,
		TaskFiles: garray.MapI(task.Files, func(i int) uint { return task.Files[i].Id }),
	}
}
