package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var data = SafeMap{
	data: make(map[string]TimeAndString),
}

type TimeAndString struct {
	time  time.Time
	value string
}

type SafeMap struct {
	data map[string]TimeAndString
	mu   sync.RWMutex
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]TimeAndString),
	}
}

func (m *SafeMap) Add(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[key]; !ok {
		m.data[key] = TimeAndString{
			time:  time.Now().Add(1 * time.Minute),
			value: value,
		}
		return nil
	}

	return fmt.Errorf("Already Used!")
}

func (m *SafeMap) refresh_cache() {
	now := time.Now()
	for key, val := range m.data {
		if now.Sub(val.time) > 0 {
			delete(m.data, key)
		}
	}
	time.Sleep(5 * time.Second)
	m.refresh_cache()
}

func AddNew(c echo.Context) error {
	key := c.QueryParam("key")
	value := c.QueryParam("value")

	err := data.Add(key, value)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"result": "1",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": "0",
	})
}

func main() {

	go data.refresh_cache()
	e := echo.New()
	e.POST("/add", AddNew)
	e.Logger.Fatal(e.Start(":6969"))
}
