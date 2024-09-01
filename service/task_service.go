package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"novanxyz/models"
	"novanxyz/repository"
	"novanxyz/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/szmcdull/glinq/garray"
)

type TaskServiceInterface interface {
	Create(task models.CreateTaskRequest) models.Task
	Update(task models.UpdateTaskRequest) models.Task
	Delete(taskId uint) models.Task
	FindById(taskId uint) models.Task
	FindAll(filter interface{}) []models.Task
	Mark(taskId uint, status string) int
	AssignTaskFile(taskId uint, uploadedFile *multipart.FileHeader) uint
	GetTaskFile(taskId uint, fileId uint) models.TaskFile
	DeleteTaskFile(taskId uint, fileId uint) int
	GetAllTaskFiles(taskId uint) []uint
}

type TaskService struct {
	TaskRepository     repository.TaskRepositoryInterface
	TaskFileRepository repository.TaskFileRepositoryInterface
	Validator          *validator.Validate
}

func NewTaskService(taskRepository repository.TaskRepositoryInterface, taskFileRepository repository.TaskFileRepositoryInterface, validator *validator.Validate) TaskServiceInterface {
	return &TaskService{
		TaskRepository:     taskRepository,
		TaskFileRepository: taskFileRepository,
		Validator:          validator,
	}
}

// Create implements TaskService
func (t *TaskService) Create(taskRequest models.CreateTaskRequest) models.Task {
	err := t.Validator.Struct(taskRequest)
	utils.ErrorPanic(err)
	taskModel := models.Task{
		Name:   taskRequest.Name,
		Status: "incomplete", //# default value
	}
	task, err := t.TaskRepository.Save(taskModel)
	utils.ErrorPanic(err)
	return task
}

// Delete implements TaskService
func (t *TaskService) Delete(taskId uint) models.Task {
	task, err := t.TaskRepository.FindById(taskId)
	utils.ErrorPanic(err)

	t.TaskRepository.Delete(taskId)
	return task
}

// FindAll implements TaskService
func (t *TaskService) FindAll(filter interface{}) []models.Task {

	filters := filter.(map[string]interface{})
	var page int
	var size int
	var err error
	if p, ok := filters["p"]; ok {
		tmp := p.(string)
		page, err = strconv.Atoi(tmp)
		utils.ErrorPanic(err)
		delete(filters, "p")
	} else {
		page = 1
	}

	if s, ok := filters["s"]; ok {
		tmp := s.(string)
		page, err = strconv.Atoi(tmp)
		utils.ErrorPanic(err)
		delete(filters, "s")
	} else {
		size = 10
	}

	fmt.Println(filter, filters, page, size)
	result := t.TaskRepository.FindAll(filters, page, size)

	return result
}

// FindById implements TaskService
func (t *TaskService) FindById(taskId uint) models.Task {
	taskData, err := t.TaskRepository.FindById(taskId)
	utils.ErrorPanic(err)
	return taskData
}

// Update implements TaskService
func (t *TaskService) Update(task models.UpdateTaskRequest) models.Task {
	taskData, err := t.TaskRepository.FindById(task.Id)
	utils.ErrorPanic(err)
	taskData.Name = task.Name
	taskData.Status = task.Status
	t.TaskRepository.Update(taskData)
	return t.FindById(task.Id)
}

func (t *TaskService) Mark(taskId uint, status string) int {
	taskData, err := t.TaskRepository.FindById(taskId)
	utils.ErrorPanic(err)
	taskData.Status = status
	affected := t.TaskRepository.Update(taskData)
	return affected
}

func (t *TaskService) AssignTaskFile(taskId uint, uploadFile *multipart.FileHeader) uint {

	fmt.Println(uploadFile.Header)
	task, err := t.TaskRepository.FindById(taskId)
	utils.ErrorPanic(err)

	fileHandler, err := uploadFile.Open()
	utils.ErrorPanic(err)

	content, err := io.ReadAll(fileHandler)
	utils.ErrorPanic(err)

	taskFile := models.TaskFile{Filename: uploadFile.Filename, Mime: uploadFile.Header.Get("Content-Type"), ParentTask: &task, Content: content}
	affected := t.TaskFileRepository.Save(taskFile)
	return affected
}

func (t *TaskService) GetTaskFile(taskId uint, fileId uint) models.TaskFile {
	taskFile, err := t.TaskRepository.FindFileById(fileId)
	utils.ErrorPanic(err)
	return taskFile
}

func (t *TaskService) DeleteTaskFile(taskId uint, fileId uint) int {
	taskData, err := t.TaskRepository.FindById(taskId)
	utils.ErrorPanic(err)
	return t.TaskRepository.Update(taskData)
}

func (t *TaskService) GetAllTaskFiles(taskId uint) []uint {
	task, err := t.TaskRepository.FindById(taskId)
	utils.ErrorPanic(err)
	fileIds := garray.MapI(task.Files, func(i int) uint {
		return task.Files[i].Id
	})

	return fileIds
}
