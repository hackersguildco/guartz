package endpoint

import (
	"net/http"
	"strconv"

	"gopkg.in/validator.v2"

	"github.com/gin-gonic/gin"
	"github.com/infiniteloopsco/guartz/models"
	"github.com/jinzhu/gorm"
)

//ExecutionCreate serves the route POST /tasks/:task_id/executions
func ExecutionCreate(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		var task models.Task
		if txn.First(&task, c.Param("task_id")); task.ID != "" {
			var execution models.Execution
			if err := c.BindJSON(&execution); err == nil {
				execution.TaskID = task.ID
				if err := validator.Validate(&execution); err == nil {
					if txn.Save(&execution).Error == nil {
						c.JSON(http.StatusOK, execution)
						return true
					} else {
						c.JSON(http.StatusBadRequest, "Execution can't be saved")
					}
				} else {
					c.JSON(http.StatusConflict, err.(validator.ErrorMap))
				}
			}
		} else {
			c.JSON(http.StatusNotFound, "")
		}
		return false
	})
}

//ExecutionList serves the route GET /tasks/:task_id/executions?page=0
func ExecutionList(c *gin.Context) {
	var executions []models.Execution
	page, _ := strconv.Atoi(c.Param("task_id"))
	offset := page * models.ExecutionPage
	models.Gdb.Where("task_id like ?", c.Param("task_id")).Offset(offset).Limit(models.ExecutionPage).Find(&executions)
	c.JSON(http.StatusOK, executions)
}
