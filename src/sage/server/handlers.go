package server

import (
	"fmt"
	"net/http"
	"sage/src/sage/cmd"
	"sage/src/sage/data"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
	"github.com/gin-gonic/gin"
)

// addHandler handles adding an expense with the given query string parameters
func addHandler(c *gin.Context) {
	dateStr := c.Query("date")
	locationStr := c.Query("location")
	descStr := c.Query("description")
	amtStr := c.Query("amount")

	if dateStr == "" || amtStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount and date are required"})
		return
	}

	date, err := civil.ParseDate(dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}
	amt, err := cmd.ParseAmount(amtStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount format"})
		return
	}

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
		c.JSON(http.StatusOK, gin.H{"success": "expense added successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": addResp.Error.Error()})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year format"})
			return
		}
	}
	if monthStr != "" {
		month, err = strconv.Atoi(monthStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month format"})
			return
		}
	}
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit format"})
			return
		}
	}
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page size format"})
			return
		}
	}
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page format"})
			return
		}
	}
	if showIdStr != "" {
		showId, err = strconv.ParseBool(showIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid show ID format"})
			return
		}
	}

	logReq, err := cmd.ParseLogArgs(startStr, endStr, year, month, limit, pageSize, page, showId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logResp := cmd.LogExpenses(db, logReq)
	if logResp.Success {
		defer logResp.Result.Close()

		var results [][]string

		var date time.Time
		var location string
		var description string
		var amt float64
		for logResp.Result.Next() {
			if logResp.ShowId {
				var id int
				err := logResp.Result.Scan(&id, &date, &location, &description, &amt)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				}
				row := []string{strconv.Itoa(id), date.Format("2006-01-02"), location, description, fmt.Sprintf("%.2f", amt)}
				results = append(results, row)
			} else {
				err := logResp.Result.Scan(&date, &location, &description, &amt)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				}
				row := []string{date.Format("2006-01-02"), location, description, fmt.Sprintf("%.2f", amt)}
				results = append(results, row)
			}
		}
		c.JSON(http.StatusOK, gin.H{"show_id": logResp.ShowId, "result": results})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": logResp.Error.Error()})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year format"})
			return
		}
	}
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit format"})
			return
		}
	}
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page size format"})
			return
		}
	}
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page format"})
			return
		}
	}

	sumReq, err := cmd.ParseSummaryArgs(startStr, endStr, year, limit, pageSize, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sumResp := cmd.SummarizeExpenses(db, sumReq)
	if sumResp.Success {
		defer sumResp.Result.Close()

		var results [][]string

		var month string
		var totalSpent float64
		for sumResp.Result.Next() {
			err := sumResp.Result.Scan(&month, &totalSpent)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			}
			row := []string{month, fmt.Sprintf("%.2f", totalSpent)}
			results = append(results, row)
		}
		c.JSON(http.StatusOK, gin.H{"result": results})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": sumResp.Error.Error()})
	}
}

// deleteHandler handles deleting an expense with the given query string parameters
func deleteHandler(c *gin.Context) {
	idStr := c.Params.ByName("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	deleteResp := cmd.DeleteExpense(db, &cmd.DeleteRequest{
		Id: id,
	})
	if deleteResp.Success {
		c.JSON(http.StatusOK, gin.H{"success": "expense deleted successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": deleteResp.Error.Error()})
	}
}
