package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"cloud.google.com/go/civil"
	"github.com/Rhymond/go-money"
)

const MAX_PAGE_SIZE = 100

// ParseAmount takes a string and constructs a money.Money object.
func ParseAmount(inputAmt string) (*money.Money, error) {
	parts := strings.Split(inputAmt, ".")
	if len(parts) > 2 {
		return nil, errors.New("amount must be in the format X.YY")
	}
	if len(parts) == 2 {
		if len(parts[1]) > 2 {
			return nil, errors.New("amount must have at most 2 decimal places")
		}
		if len(parts[1]) == 1 {
			inputAmt += "0"
		} else if len(parts[1]) == 0 {
			return nil, errors.New("amount must be in the format X.YY")
		}
	} else if len(parts) == 1 {
		inputAmt += ".00"
	}

	i, err := strconv.ParseInt(strings.ReplaceAll(inputAmt, ".", ""), 10, 64)
	if err != nil {
		return nil, errors.New("error parsing amount: " + err.Error())
	}
	if i < 0 {
		return nil, errors.New("amount cannot be negative")
	}
	if i == 0 {
		return nil, errors.New("amount cannot be zero")
	}
	amt := money.New(i, money.USD)

	return amt, nil
}

// ParseLogArgs takes a list of args and constructs the appropriate LogRequest. year, month, limit, pageSize, and page
// default to 0. showId defaults to false.
func ParseLogArgs(startStr, endStr string, year, month, limit, pageSize, page int, showId bool) (*LogRequest, error) {
	var err error

	start := civil.Date{}
	if startStr != "" {
		start, err = civil.ParseDate(startStr)
		if err != nil {
			return nil, errors.New("error parsing start date: " + err.Error())
		}
	}
	end := civil.Date{}
	if endStr != "" {
		end, err = civil.ParseDate(endStr)
		if err != nil {
			return nil, errors.New("error parsing end date: " + err.Error())
		}
	}

	if !start.IsZero() && !end.IsZero() {
		if start.After(end) {
			return nil, errors.New("start date is after end date")
		}
	}
	if month < 0 || month > 12 {
		return nil, errors.New("month must be between 1 and 12")
	}
	if year < 0 {
		return nil, errors.New("year must be positive")
	}
	if limit != 0 && (pageSize != 0 || page != 0) {
		return nil, errors.New("cannot provide limit with page size or page")
	}
	if page != 0 && pageSize == 0 {
		return nil, errors.New("must provide page size with page")
	}
	if pageSize > MAX_PAGE_SIZE {
		return nil, fmt.Errorf("page size must be at most %d", MAX_PAGE_SIZE)
	}
	if pageSize < 0 {
		return nil, errors.New("page size must be positive")
	}
	if limit < 0 {
		return nil, errors.New("limit must be positive")
	}
	if page < 0 {
		return nil, errors.New("page must be positive")
	}

	return &LogRequest{
		Start:    start,
		End:      end,
		Year:     year,
		Month:    month,
		Limit:    limit,
		PageSize: pageSize,
		Page:     page,
		ShowId:   showId,
	}, nil
}

// ParseSummaryArgs takes a list of args and constructs the appropriate SummaryRequest. year, limit, and page default to
// 0. pageSize defaults to 100.
func ParseSummaryArgs(startStr, endStr string, year, limit, pageSize, page int) (*SummaryRequest, error) {
	var err error

	start := civil.Date{}
	if startStr != "" {
		start, err = civil.ParseDate(startStr)
		if err != nil {
			return nil, errors.New("error parsing start date: " + err.Error())
		}
	}
	end := civil.Date{}
	if endStr != "" {
		end, err = civil.ParseDate(endStr)
		if err != nil {
			return nil, errors.New("error parsing end date: " + err.Error())
		}
	}

	if !start.IsZero() && !end.IsZero() {
		if start.After(end) {
			return nil, errors.New("start date is after end date")
		}
	}
	if limit != 0 && (pageSize != 0 || page != 0) {
		return nil, errors.New("cannot provide limit with page size or page")
	}
	if page != 0 && pageSize == 0 {
		return nil, errors.New("must provide page size with page")
	}
	if pageSize > MAX_PAGE_SIZE {
		return nil, fmt.Errorf("page size must be at most %d", MAX_PAGE_SIZE)
	}
	if pageSize < 0 {
		return nil, errors.New("page size must be positive")
	}
	if limit < 0 {
		return nil, errors.New("limit must be positive")
	}
	if page < 0 {
		return nil, errors.New("page must be positive")
	}

	return &SummaryRequest{
		Start:    start,
		End:      end,
		Year:     year,
		Limit:    limit,
		PageSize: pageSize,
		Page:     page,
	}, nil
}
