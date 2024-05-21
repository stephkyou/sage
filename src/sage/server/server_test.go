package server

import (
	"encoding/json"
	"net/http/httptest"
	"sage/src/sage/cmd"
	"testing"

	"github.com/gin-gonic/gin"
	"gotest.tools/v3/assert"
)

func teardown() {
	db.Exec("DELETE FROM expenses")
	db.Close()
}

func TestAddHandler(t *testing.T) {
	defer teardown()

	db, _ = cmd.ConnectDB(cmd.TEST_DB_NAME)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/add", nil)
	q := c.Request.URL.Query()
	q.Add("date", "2021-01-01")
	q.Add("location", "Test Location")
	q.Add("description", "Test Description")
	q.Add("amount", "20.12")
	c.Request.URL.RawQuery = q.Encode()

	addHandler(c)
	assert.Equal(t, 200, w.Code)
}

func TestLogHandler(t *testing.T) {
	defer teardown()

	db, _ = cmd.ConnectDB(cmd.TEST_DB_NAME)
	db.Exec("INSERT INTO expenses (date_spent, location, description, amt) VALUES ('2021-01-01', 'Test Location', 'Test Description', 20.12)")
	db.Exec("INSERT INTO expenses (date_spent, location, description, amt) VALUES ('2022-04-16', 'Test Location 2', 'Test Description 2', 69.24)")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/log", nil)

	logHandler(c)
	assert.Equal(t, 200, w.Code)

	body := gin.H{
		"result": [][]string{
			{"2021-01-01", "Test Location", "Test Description", "20.12"},
			{"2022-04-16", "Test Location 2", "Test Description 2", "69.24"},
		},
		"show_id": false,
	}
	response, err := json.Marshal(body)

	assert.NilError(t, err)
	assert.Equal(t, w.Body.String(), string(response))
}

func TestSummaryHandler(t *testing.T) {
	defer teardown()

	db, _ = cmd.ConnectDB(cmd.TEST_DB_NAME)
	db.Exec("INSERT INTO expenses (date_spent, location, description, amt) VALUES ('2021-01-01', 'Test Location', 'Test Description', 20.12)")
	db.Exec("INSERT INTO expenses (date_spent, location, description, amt) VALUES ('2022-04-16', 'Test Location 2', 'Test Description 2', 2.00)")
	db.Exec("INSERT INTO expenses (date_spent, location, description, amt) VALUES ('2022-04-25', 'Test Location 3', 'Test Description 3', 69.24)")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/summary", nil)

	summaryHandler(c)
	assert.Equal(t, 200, w.Code)

	body := gin.H{
		"result": [][]string{
			{"2021-01", "20.12"},
			{"2022-04", "71.24"},
		},
	}
	response, err := json.Marshal(body)

	assert.NilError(t, err)
	assert.Equal(t, w.Body.String(), string(response))
}

func TestDeleteHandler(t *testing.T) {
	defer teardown()

	db, _ = cmd.ConnectDB(cmd.TEST_DB_NAME)
	db.Exec("INSERT INTO expenses (date_spent, location, description, amt) VALUES ('2021-01-01', 'Test Location', 'Test Description', 20.12)")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/delete", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	deleteHandler(c)
	assert.Equal(t, 200, w.Code)
}
