package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/civil"
	"github.com/Rhymond/go-money"
)

type Expense struct {
	Id          int
	Date        civil.Date
	Location    string
	Description string
	Amount      *money.Money
}

var createTableQuery string = `CREATE TABLE IF NOT EXISTS expenses (
	id INTEGER PRIMARY KEY,
	date_spent DATE NOT NULL,
	location VARCHAR(255),
	description VARCHAR(255),
	amt DECIMAL(19,4) NOT NULL
	)`

func main() {
	os.Exit(runMain())
}

func runMain() int {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(`Valid sage commands:
		add <date> <location> <description> <amount>
		log [--start <date>] [--end <date>] [--year <year>] [--month <month>] [-n <limit>] [--page-size <size>] [--page <page>] [--show-id]
		summary [--start <date>] [--end <date>] [--year <year>] [-n <limit>] [--page-size <size>] [--page <page>]
		delete <id>`)
		return 0
	}

	return execRequest(args)
}

// execRequest accepts a list of args, the first element of which is assumed to be the command name.
// Validation is performed on the given args and the appropriate command is executed.
func execRequest(args []string) int {
	cmd := args[0]
	switch cmd {
	case "add":
		if len(args) < 5 {
			log.Println("not enough fields provided")
			return 1
		}

		addReq, err := parseAddRequest(args[1:])
		if err != nil {
			log.Println("error parsing add request", err)
			return 1
		}
		err = verifyDatabase()
		if err != nil {
			log.Println("error verifying database", err)
			return 1
		}

		return AddExpense(addReq)
	case "log":
		if len(args) == 1 {
			return LogExpenses(&LogRequest{
				ShowId: false,
			})
		}

		logReq, err := parseLogRequest(args[1:])
		if err != nil {
			log.Println("error parsing log request: ", err)
			return 1
		}
		err = verifyDatabase()
		if err != nil {
			log.Println("error verifying database", err)
			return 1
		}

		return LogExpenses(logReq)
	case "summary":
		if len(args) == 1 {
			return SummarizeExpenses(&SummaryRequest{})
		}

		sumReq, err := parseSummaryRequest(args[1:])
		if err != nil {
			log.Println("error parsing summary request: ", err)
			return 1
		}
		err = verifyDatabase()
		if err != nil {
			log.Println("error verifying database", err)
			return 1
		}

		return SummarizeExpenses(sumReq)
	case "delete":
		if len(args) < 2 {
			log.Println("need to provide an ID to delete")
			return 1
		}
		if len(args) > 2 {
			log.Println("can only delete one ID at a time")
			return 1
		}

		err := verifyDatabase()
		if err != nil {
			log.Println("error verifying database", err)
			return 1
		}
		return DeleteExpense(&DeleteRequest{
			Id: args[1],
		})
	default:
		fmt.Println("Invalid command")
		return 1
	}
}

// verifyDatabase checks if the sage folder and sage.db database exists. Creates the necessary folder and SQLite file
// if it doesn't.
func verifyDatabase() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return errors.New("error getting user home directory: " + err.Error())
	}

	if _, err := os.Stat(dirname + "/sage"); errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(dirname+"/sage", 0755)
		if err != nil {
			return errors.New("error creating sage directory: " + err.Error())
		}
	}

	if _, err := os.Stat(dirname + "/sage/sage.db"); errors.Is(err, fs.ErrNotExist) {
		file, err := os.Create(dirname + "/sage/sage.db")
		if err != nil {
			return errors.New("error creating database file: " + err.Error())
		}
		file.Close()
	}

	return nil
}

// parseAddRequest takes a list of provided fields and constructs the appropriate AddRequest. Assumes 4 fields are provided.
func parseAddRequest(args []string) (*AddRequest, error) {
	date, err := civil.ParseDate(args[0])
	if err != nil {
		return nil, errors.New("error parsing date: " + err.Error())
	}

	inputAmt := args[3]
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

	return &AddRequest{
		Expense: Expense{
			Date:        date,
			Location:    args[1],
			Description: args[2],
			Amount:      amt,
		},
	}, nil
}

// parseLogRequest takes a list of args and constructs the appropriate LogRequest.
func parseLogRequest(args []string) (*LogRequest, error) {
	var err error

	logCmd := flag.NewFlagSet("log", flag.ExitOnError)
	startStr := logCmd.String("start", "", "start date")
	endStr := logCmd.String("end", "", "end date")
	year := logCmd.Int("year", 0, "year")
	month := logCmd.Int("month", 0, "month")
	limit := logCmd.Int("n", 0, "limit")
	pageSize := logCmd.Int("page-size", 0, "page size")
	page := logCmd.Int("page", 0, "page")
	showId := logCmd.Bool("show-id", false, "show the expense ID")

	logCmd.Parse(args)

	start := civil.Date{}
	if *startStr != "" {
		start, err = civil.ParseDate(*startStr)
		if err != nil {
			return nil, errors.New("error parsing start date: " + err.Error())
		}
	}
	end := civil.Date{}
	if *endStr != "" {
		end, err = civil.ParseDate(*endStr)
		if err != nil {
			return nil, errors.New("error parsing end date: " + err.Error())
		}
	}

	if !start.IsZero() && !end.IsZero() {
		if start.After(end) {
			return nil, errors.New("start date is after end date")
		}
	}
	if *limit != 0 && (*pageSize != 0 || *page != 0) {
		return nil, errors.New("cannot provide limit with page size or page")
	}
	if *page != 0 && *pageSize == 0 {
		return nil, errors.New("must provide page size with page")
	}

	return &LogRequest{
		Start:    start,
		End:      end,
		Year:     *year,
		Month:    *month,
		Limit:    *limit,
		PageSize: *pageSize,
		Page:     *page,
		ShowId:   *showId,
	}, nil
}

// parseSummaryRequest takes a list of args and constructs the appropriate SummaryRequest.
func parseSummaryRequest(args []string) (*SummaryRequest, error) {
	var err error

	summCmd := flag.NewFlagSet("log", flag.ExitOnError)
	startStr := summCmd.String("start", "", "start date")
	endStr := summCmd.String("end", "", "end date")
	year := summCmd.Int("year", 0, "year")
	limit := summCmd.Int("n", 0, "limit")
	pageSize := summCmd.Int("page-size", 0, "page size")
	page := summCmd.Int("page", 0, "page")

	summCmd.Parse(args)

	start := civil.Date{}
	if *startStr != "" {
		start, err = civil.ParseDate(*startStr)
		if err != nil {
			return nil, errors.New("error parsing start date: " + err.Error())
		}
	}
	end := civil.Date{}
	if *endStr != "" {
		end, err = civil.ParseDate(*endStr)
		if err != nil {
			return nil, errors.New("error parsing end date: " + err.Error())
		}
	}

	if !start.IsZero() && !end.IsZero() {
		if start.After(end) {
			return nil, errors.New("start date is after end date")
		}
	}
	if *limit != 0 && (*pageSize != 0 || *page != 0) {
		return nil, errors.New("cannot provide limit with page size or page")
	}
	if *page != 0 && *pageSize == 0 {
		return nil, errors.New("must provide page size with page")
	}

	return &SummaryRequest{
		Start:    start,
		End:      end,
		Year:     *year,
		Limit:    *limit,
		PageSize: *pageSize,
		Page:     *page,
	}, nil
}
