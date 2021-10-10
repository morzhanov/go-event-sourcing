package orders

import (
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/morzhanov/go-event-sourcing-example/internal"
	"io/ioutil"
	"net/http"
)

type Order struct {
	Id     string `json:"id,omitempty" db:"id"`
	Name   string `json:"name,omitempty" db:"name"`
	Status string `json:"status,omitempty" db:"status"`
}

type service struct {
	cs internal.CommandStore
	qs internal.QueryStore
}

type Service interface {
	Listen() error
}

func (o *service) handleAddOrder(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	order := Order{}
	if err = json.Unmarshal(jsonData, &order); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if err := o.cs.Add("create_order", &order); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusCreated)
}

func (o *service) handleProcessOrder(c *gin.Context) {
	id := c.Param("id")
	if err := o.cs.Add("process_order", id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (o *service) handleGetOrder(c *gin.Context) {
	id := c.Param("id")
	res, err := o.qs.Get(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func (o *service) Listen() error {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders([]string{"authorization"}...)
	router.Use(cors.New(config))
	r := gin.Default()
	r.POST("/", o.handleAddOrder)
	r.POST("/:id", o.handleProcessOrder)
	r.GET("/:id", o.handleGetOrder)
	return r.Run()
}

func NewService(cs internal.CommandStore, qs internal.QueryStore) Service {
	return &service{cs, qs}
}
