package controller

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sage/src/sage/cmd"
	"sage/src/sage/data"
	"sage/src/sage/server"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
)

func RunCLIController() int {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(`Valid sage commands:
		add <date> <location> <description> <amount>
		log [--start <date>] [--end <date>] [--year <year>] [--month <month>] [-n <limit>] [--page-size <size>] [--page <page>] [--show-id]
		summary [--start <date>] [--end <date>] [--year <year>] [-n <limit>] [--page-size <size>] [--page <page>]
		delete <id>`)
		return 0
	}

	var err error
	command := args[0]
	switch command {
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

		db, err := cmd.ConnectDB("sage.db")
		if err != nil {
			log.Println("error connecting to database: ", err)
			return 1
		}
		addResp := cmd.AddExpense(db, addReq)
		if addResp.Success {
			fmt.Println("Expense added successfully")
		} else {
			fmt.Println("Error adding expense: ", addResp.Error)
			return 1
		}
	case "log":
		logReq := &cmd.LogRequest{ShowId: false}
		if len(args) != 1 {
			logReq, err = parseLogRequest(args[1:])
			if err != nil {
				log.Println("error parsing log request: ", err)
				return 1
			}
		}

		db, err := cmd.ConnectDB("sage.db")
		if err != nil {
			log.Println("error connecting to database: ", err)
			return 1
		}
		logResp := cmd.LogExpenses(db, logReq)
		if logResp.Success {
			defer logResp.Result.Close()

			var date time.Time
			var location string
			var description string
			var amt float64
			for logResp.Result.Next() {
				if logResp.ShowId {
					var id int
					err := logResp.Result.Scan(&id, &date, &location, &description, &amt)
					if err != nil {
						log.Println("error reading retrieved expenses: " + err.Error())
						return 1
					}
					fmt.Printf("%d | %s | %s | %s | $%.2f\n", id, date.Format("2006-01-02"), location, description, amt)
				} else {
					err := logResp.Result.Scan(&date, &location, &description, &amt)
					if err != nil {
						log.Println("error reading retrieved expenses: " + err.Error())
						return 1
					}
					fmt.Printf("%s | %s | %s | $%.2f\n", date.Format("2006-01-02"), location, description, amt)
				}
			}
		} else {
			fmt.Println("Error logging expenses: ", logResp.Error)
			return 1
		}
	case "summary":
		sumReq := &cmd.SummaryRequest{}
		if len(args) != 1 {
			sumReq, err = parseSummaryRequest(args[1:])
			if err != nil {
				log.Println("error parsing summary request: ", err)
				return 1
			}
		}

		db, err := cmd.ConnectDB("sage.db")
		if err != nil {
			log.Println("error connecting to database: ", err)
			return 1
		}
		sumResp := cmd.SummarizeExpenses(db, sumReq)
		if sumResp.Success {
			defer sumResp.Result.Close()

			var month string
			var totalSpent float64
			for sumResp.Result.Next() {
				err = sumResp.Result.Scan(&month, &totalSpent)
				if err != nil {
					log.Println("error reading calculated summary: " + err.Error())
				}
				fmt.Printf("%s: $%.2f\n", month, totalSpent)
			}
		} else {
			fmt.Println("Error summarizing expenses: ", sumResp.Error)
			return 1
		}
	case "delete":
		if len(args) < 2 {
			log.Println("need to provide an ID to delete")
			return 1
		}
		if len(args) > 2 {
			log.Println("can only delete one ID at a time")
			return 1
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			log.Println("invalid ID provided: ", err)
			return 1
		}

		db, err := cmd.ConnectDB("sage.db")
		if err != nil {
			log.Println("error connecting to database: ", err)
			return 1
		}
		deleteResp := cmd.DeleteExpense(db, &cmd.DeleteRequest{
			Id: id,
		})
		if deleteResp.Success {
			fmt.Println("Expense deleted successfully")
		} else {
			fmt.Println("Error deleting expense: ", deleteResp.Error)
			return 1
		}
	case "server":
		err := server.RunServer()
		if err != nil {
			log.Println("error running server: ", err)
			return 1
		}
	default:
		fmt.Println("Invalid command")
		return 1
	}

	return 0
}

// parseAddRequest takes a list of provided fields and constructs the appropriate AddRequest. Assumes 4 fields are provided.
func parseAddRequest(args []string) (*cmd.AddRequest, error) {
	date, err := civil.ParseDate(args[0])
	if err != nil {
		return nil, errors.New("error parsing date: " + err.Error())
	}

	amt, err := cmd.ParseAmount(args[3])
	if err != nil {
		return nil, err
	}

	return &cmd.AddRequest{
		Expense: data.Expense{
			Date:        date,
			Location:    args[1],
			Description: args[2],
			Amount:      amt,
		},
	}, nil
}

// parseLogRequest takes a list of args and constructs the appropriate LogRequest.
func parseLogRequest(args []string) (*cmd.LogRequest, error) {
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

	return cmd.ParseLogArgs(*startStr, *endStr, *year, *month, *limit, *pageSize, *page, *showId)
}

// parseSummaryRequest takes a list of args and constructs the appropriate SummaryRequest.
func parseSummaryRequest(args []string) (*cmd.SummaryRequest, error) {
	summCmd := flag.NewFlagSet("log", flag.ExitOnError)
	startStr := summCmd.String("start", "", "start date")
	endStr := summCmd.String("end", "", "end date")
	year := summCmd.Int("year", 0, "year")
	limit := summCmd.Int("n", 0, "limit")
	pageSize := summCmd.Int("page-size", 0, "page size")
	page := summCmd.Int("page", 0, "page")

	summCmd.Parse(args)

	return cmd.ParseSummaryArgs(*startStr, *endStr, *year, *limit, *pageSize, *page)
}
