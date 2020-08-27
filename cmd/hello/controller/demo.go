package controller

import (
	"github.com/gin-gonic/gin"
)

type Demo struct {

}

func (a *Demo) Query(c *gin.Context) {
	c.Writer.WriteString("Hello World!")
}