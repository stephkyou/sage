package cmd

import (
	"errors"
	"fmt"

	"cloud.google.com/go/civil"
)

const MAX_PAGE_SIZE = 100

// ParseLogArgs takes a list of args and constructs the appropriate LogRequest. year, month, limit, pageSize, and page
// default to 0. showId defaults to false.
func ParseLogArgs(startStr, endStr string, year, month, limit, pageSize, page int, showId bool, query string) (*LogRequest, error) {
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
		Query:    query,
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
