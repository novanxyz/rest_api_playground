package controller

import (
	"fmt"
	"net/http"
	"novanxyz/models"
	"novanxyz/service"
	"novanxyz/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type TaskController struct {
	taskService service.TaskServiceInterface
}

func NewTaskController(service service.TaskServiceInterface) *TaskController {
	return &TaskController{
		taskService: service,
	}
}

// CreateTask		godoc
// @Summary			Create a new task
// @Description		Save task data in Db.
// @Param			task body models.CreateTaskRequest true "Create task"
// @Param			task body models.CreateTaskRequest true "Create task"
// @Produce			application/json
// @Task			task
// @Success			200 {object} models.Response{}
// @Router			/tasks [post]
func (controller *TaskController) Create(ctx *gin.Context) {
	log.Info().Msg("create task")
	createTaskRequest := models.CreateTaskRequest{}
	err := ctx.ShouldBindJSON(&createTaskRequest)
	utils.ErrorPanic(err)

	newId := controller.taskService.Create(createTaskRequest)
	// utils.ErrorPanic(err)

	webResponse := models.Response{
		Code:   http.StatusCreated,
		Status: "OK",
		Data:   newId,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusCreated, webResponse)
}

// FindAllTask 		godoc
// @Summary			Retrieve a list of task with optional filtering by status (complete|incomplete).
// @Description		Return list of task.
// @Task			task
// @Success			200 {object} models.Response{}
// @Router			/tasks [get]
func (controller *TaskController) FindAll(ctx *gin.Context) {
	log.Info().Msg("findAll task")

	filter := utils.QueryParamMap(ctx.Request.URL.Query())
	tasks := controller.taskService.FindAll(filter)
	// utils.ErrorPanic(err)
	var taskResponses []models.TaskResponse

	for _, task := range tasks {
		taskResponses = append(taskResponses, task.ToResponse())
	}

	webResponse := models.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   taskResponses,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)

}

// FindByIdTask 		godoc
// @Summary				Retrive a specific  task by id.
// @Param				taskId path int true "find task by taskId"
// @Description			Return the task which id value matches taskId.
// @Produce				application/json
// @Task				task
// @Success				200 {object} models.Response{}
// @Router				/tasks/{taskId} [get]
func (controller *TaskController) FindById(ctx *gin.Context) {
	log.Info().Msg("findbyid task")
	taskId := ctx.Param("taskId")
	id, err := strconv.Atoi(taskId)
	utils.ErrorPanic(err)

	taskResponse := controller.taskService.FindById(uint(id))

	webResponse := models.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   taskResponse,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// UpdateTask		godoc
// @Summary			Update specific task by Id
// @Description		Update task data.
// @Param			taskId path int true "update task by id"
// @Param			task body models.UpdateTaskRequest true  "Update task"
// @Task			task
// @Produce			application/json
// @Success			200 {object} models.Response{}
// @Router			/tasks/{taskId} [put]
func (controller *TaskController) Update(ctx *gin.Context) {
	log.Info().Msg("update task")
	updateTaskRequest := models.UpdateTaskRequest{}
	err := ctx.ShouldBindJSON(&updateTaskRequest)
	utils.ErrorPanic(err)

	taskId := ctx.Param("taskId")
	id, err := strconv.Atoi(taskId)
	utils.ErrorPanic(err)

	updateTaskRequest.Id = uint(id)

	affected := controller.taskService.Update(updateTaskRequest)

	webResponse := models.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   affected,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// DeleteTask		godoc
// @Summary			Delete a specific task by Id
// @Description		Remove task data by id.
// @Produce			application/json
// @Task			task
// @Success			200 {object} models.Response{}
// @Router			/tasks/{taskId} [delete]
func (controller *TaskController) Delete(ctx *gin.Context) {
	log.Info().Msg("delete task")
	taskId := ctx.Param("taskId")
	id, err := strconv.Atoi(taskId)
	utils.ErrorPanic(err)

	deleted := controller.taskService.Delete(uint(id))

	// if affected == nil {
	// 	utils.ResponseError(errors.New("delete failed"), ctx)
	// }

	webResponse := models.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   deleted,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Mark 				godoc
// @Summary				Mark a specific task as complete or incomplete.
// @Param				taskId path int true "update task status by id"
// @Param				status path string true "mark task as  complete|incomplete status" enum[complete, incomplete]
// @Description			Return the task which Id value matches taskId.
// @Produce				application/json
// @Task				task
// @Success				200 {object} models.Response{}
// @Router				/tasks/{taskId}/{status} [patch]
func (controller *TaskController) Mark(ctx *gin.Context) {
	log.Info().Msg("mark task")
	taskId := ctx.Param("taskId")
	status := ctx.Param("status")
	id, err := strconv.Atoi(taskId)
	utils.ErrorPanic(err)

	affected := controller.taskService.Mark(uint(id), status)

	webResponse := models.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   affected,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Mark 				godoc
// @Summary				Upload a specific file associated with a task.
// @Param				taskId path int true "update task by id"
// @Param				file formData file true "file to be upload"
// @Description			Return file Id of uploaded file, and task it's associated.
// @Consumes			multipart/form-data
// @Produce				application/json
// @Task				task
// @Success				200 {object} models.Response{}
// @Router				/tasks/{taskId}/files [post]
func (controller *TaskController) UploadTaskFile(ctx *gin.Context) {
	log.Info().Msg("upload task file")
	taskId := ctx.Param("taskId")
	uploadedFile, err := ctx.FormFile("file")

	fmt.Println(uploadedFile.Header)
	utils.ErrorPanic(err)

	id, err := strconv.Atoi(taskId)
	utils.ErrorPanic(err)

	taskResponse := controller.taskService.AssignTaskFile(uint(id), uploadedFile)

	webResponse := models.Response{
		Code:   http.StatusCreated,
		Status: "OK",
		Data:   taskResponse,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Mark 				godoc
// @Summary				Download a specific file associated with a task .
// @Param				taskId path int true "task id which file belong to"
// @Param				fileId path int true "file id to delete"
// @Description			Download file stream with content disposition.
// @Produce				application/octet-stream
// @Task				task
// @Success				200 {object} models.Response{}
// @Router				/tasks/{taskId}/files/{fileId} [get]
func (controller *TaskController) DownloadTaskFile(ctx *gin.Context) {
	log.Info().Msg("download task file")
	taskId := ctx.Param("taskId")
	fileId := ctx.Param("fileId")
	id, err := strconv.Atoi(taskId)
	utils.ErrorPanic(err)

	fid, err := strconv.Atoi(fileId)
	utils.ErrorPanic(err)

	taskFile := controller.taskService.GetTaskFile(uint(id), uint(fid))

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", taskFile.Filename))
	ctx.Data(http.StatusOK, taskFile.Mime, taskFile.Content)

}

// Mark 				godoc
// @Summary				Delete a specific file associated with a task by Id.
// @Param				taskId path int true "task id which file belong to"
// @Param				fileId path int true "file id to delete"
// @Description			return success if file deleted.
// @Produce				application/json
// @Task				task
// @Success				200 {object} models.Response{}
// @Router				/tasks/{taskId}/files/{fileId} [delete]
func (controller *TaskController) DeleteTaskFile(ctx *gin.Context) {
	log.Info().Msg("delete task file")
	taskId := ctx.Param("taskId")
	fileId := ctx.Param("fileId")
	id, err := strconv.Atoi(taskId)
	utils.ErrorPanic(err)

	fid, err := strconv.Atoi(fileId)
	utils.ErrorPanic(err)

	deleted := controller.taskService.DeleteTaskFile(uint(id), uint(fid))

	webResponse := models.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   deleted,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// Mark 				godoc
// @Summary				Get a File assigned to task by Id .
// @Param				taskId path int true "task id which the files belong to "
// @Description			Return list of filename and url of files in models.data.
// @Produce				application/json
// @Task				task
// @Success				200 {object} models.Response{}
// @Router				/tasks/{taskId}/files [get]
func (controller *TaskController) GetTaskFiles(ctx *gin.Context) {
	log.Info().Msg("mark task")
	taskId := ctx.Param("taskId")
	id, err := strconv.Atoi(taskId)
	utils.ErrorPanic(err)

	fileids := controller.taskService.GetAllTaskFiles(uint(id))

	webResponse := models.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   fileids,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
