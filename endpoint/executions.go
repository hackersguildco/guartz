package endpoint

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/infiniteloopsco/guartz/models"
	"github.com/jinzhu/gorm"
)

//ExecutionCreate serves the route POST /tasks/:task_id/executions
func ExecutionCreate(c *gin.Context) {
	models.InTx(func(txn *gorm.DB) bool {
		var task models.Task
		if txn.Where("id like ? ", c.Param("task_id")).First(&task); task.ID != "" {
			var execution models.Execution
			if err := c.BindJSON(&execution); err == nil {
				execution.TaskID = task.ID
				if valid, errMap := models.ValidStruct(&execution); valid {
					if txn.Create(&execution).Error == nil {
						c.JSON(http.StatusOK, execution)
						return true
					} else {
						c.JSON(http.StatusBadRequest, "Execution can't be saved")
					}
				} else {
					c.JSON(http.StatusConflict, errMap)
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
	page, _ := strconv.Atoi(c.Param("page"))
	offset := page * models.ExecutionPage
	models.Gdb.Where("task_id like ?", c.Param("task_id")).Offset(offset).Limit(models.ExecutionPage).Find(&executions)
	c.JSON(http.StatusOK, executions)
}
