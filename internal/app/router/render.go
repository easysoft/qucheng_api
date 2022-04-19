package router

import "github.com/gin-gonic/gin"

func renderError(c *gin.Context, code int, err error) {
	_ = c.Error(err)
	c.JSON(code, gin.H{
		"success": false,
		"message": err.Error(),
	})
}

func renderSuccess(c *gin.Context, code int) {
	c.JSON(code, gin.H{
		"success": true,
		"message": "success",
	})
}

func renderJson(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{
		"success": true,
		"message": "success",
		"data":    data,
	})
}

func renderJsonWithPagination(c *gin.Context, code int, data interface{}, p interface{}) {
	c.JSON(code, gin.H{
		"success":    true,
		"message":    "success",
		"data":       data,
		"pagination": p,
	})
}
