package server

import (
	"net/http"
	"sage/src/sage/cmd"
	"sage/src/sage/data"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
	"github.com/Rhymond/go-money"
	"github.com/gin-gonic/gin"
)

// addHandler handles adding an expense with the given query string parameters
func addHandler(c *gin.Context) {
	dateStr := c.Query("date")
	locationStr := c.Query("location")
	descStr := c.Query("description")
	amtStr := c.Query("amount")

	if dateStr == "" || amtStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "amount and date are required"})
		return
	}

	date, err := civil.ParseDate(dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid date format"})
		return
	}
	fl, err := strconv.ParseFloat(amtStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid amount format"})
		return
	}
	amt := money.NewFromFloat(fl, money.USD)

	addReq := &cmd.AddRequest{
		Expense: data.Expense{
			Date:        date,
			Location:    locationStr,
			Description: descStr,
			Amount:      amt,
		},
	}

	addResp := cmd.AddExpense(db, addReq)
	if addResp.Success {
		c.JSON(http.StatusOK, gin.H{"message": "expense added successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": addResp.Error.Error()})
	}
}

// logHandler handles logging expenses with the given query string parameters
func logHandler(c *gin.Context) {
	startStr := c.Query("start")
	endStr := c.Query("end")
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	limitStr := c.Query("limit")
	pageSizeStr := c.Query("page-size")
	pageStr := c.Query("page")
	showIdStr := c.Query("show-id")
	query := c.Query("query")

	year := 0
	month := 0
	limit := 0
	pageSize := cmd.MAX_PAGE_SIZE
	page := 0
	showId := false
	var err error

	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid year format"})
			return
		}
	}
	if monthStr != "" {
		month, err = strconv.Atoi(monthStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid month format"})
			return
		}
	}
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid limit format"})
			return
		}
	}
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid page size format"})
			return
		}
	}
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid page format"})
			return
		}
	}
	if showIdStr != "" {
		showId, err = strconv.ParseBool(showIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid show ID format"})
			return
		}
	}

	logReq, err := cmd.ParseLogArgs(startStr, endStr, year, month, limit, pageSize, page, showId, query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	logResp := cmd.LogExpenses(db, logReq)
	if logResp.Success {
		defer logResp.Result.Close()

		var results []data.Expense
		var date time.Time
		var location string
		var description string
		var amt money.Amount

		for logResp.Result.Next() {
			if logResp.ShowId {
				var id int
				err := logResp.Result.Scan(&id, &date, &location, &description, &amt)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": err})
				}
				row := data.Expense{
					Id:          id,
					Date:        civil.DateOf(date),
					Location:    location,
					Description: description,
					Amount:      money.New(amt, "USD"),
				}
				results = append(results, row)
			} else {
				err := logResp.Result.Scan(&date, &location, &description, &amt)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": err})
				}
				row := data.Expense{
					Date:        civil.DateOf(date),
					Location:    location,
					Description: description,
					Amount:      money.New(amt, "USD"),
				}
				results = append(results, row)
			}
		}
		c.JSON(http.StatusOK, gin.H{"show_id": logResp.ShowId, "result": results})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": logResp.Error.Error()})
	}
}

// summaryHandler handles summarizing expenses with the given query string parameters
func summaryHandler(c *gin.Context) {
	startStr := c.Query("start")
	endStr := c.Query("end")
	yearStr := c.Query("year")
	limitStr := c.Query("limit")
	pageSizeStr := c.Query("page-size")
	pageStr := c.Query("page")

	year := 0
	limit := 0
	pageSize := cmd.MAX_PAGE_SIZE
	page := 0
	var err error

	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid year format"})
			return
		}
	}
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid limit format"})
			return
		}
	}
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid page size format"})
			return
		}
	}
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid page format"})
			return
		}
	}

	sumReq, err := cmd.ParseSummaryArgs(startStr, endStr, year, limit, pageSize, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	sumResp := cmd.SummarizeExpenses(db, sumReq)
	if sumResp.Success {
		defer sumResp.Result.Close()

		var results []data.Summary
		var month string
		var totalSpent money.Amount

		for sumResp.Result.Next() {
			err := sumResp.Result.Scan(&month, &totalSpent)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			}
			row := data.Summary{
				Month: month,
				Total: money.New(totalSpent, "USD"),
			}
			results = append(results, row)
		}
		c.JSON(http.StatusOK, gin.H{"result": results})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": sumResp.Error.Error()})
	}
}

// deleteHandler handles deleting an expense with the given query string parameters
func deleteHandler(c *gin.Context) {
	idStr := c.Params.ByName("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	deleteResp := cmd.DeleteExpense(db, &cmd.DeleteRequest{
		Id: id,
	})
	if deleteResp.Success {
		c.JSON(http.StatusOK, gin.H{"message": "expense deleted successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": deleteResp.Error.Error()})
	}
}

// countHandler handles counting the number of total expenses
func countHandler(c *gin.Context) {
	typeStr := c.Params.ByName("type")
	if typeStr != "log" && typeStr != "summary" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid type"})
		return
	}

	countResp := cmd.CountExpenses(db, &cmd.CountRequest{Type: typeStr})
	if countResp.Success {
		c.JSON(http.StatusOK, gin.H{"count": countResp.Result})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": countResp.Error.Error()})
	}
}
