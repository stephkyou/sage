package cmd

import (
	"database/sql/driver"
	"sage/src/sage/data"
	"testing"

	"cloud.google.com/go/civil"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Rhymond/go-money"
	"gotest.tools/v3/assert"
)

const SAGE_TEST_DB_NAME string = "sage_test.db"

func TestAddExpense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))

	addResp := AddExpense(db, &AddRequest{
		Expense: data.Expense{
			Date:        civil.Date{Year: 2021, Month: 1, Day: 1},
			Location:    "Test Location",
			Description: "Test Description",
			Amount:      money.NewFromFloat(32.45, money.USD),
		}})
	assert.Assert(t, addResp.Success)
	assert.NilError(t, addResp.Error)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestLogExpenses(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating mock database: %v", err)
	}
	defer db.Close()

	values := [][]driver.Value{
		{
			"2021-01-01", "Test Location", "Test Description", 32.45,
		},
		{
			"2023-06-24", "Test Location 2", "Test Description 2", 9.01,
		},
	}

	rows := sqlmock.NewRows([]string{"date_spent", "location", "description", "amt"}).AddRows(values...)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	logResp := LogExpenses(db, &LogRequest{ShowId: false})
	assert.Assert(t, logResp.Success)
	assert.NilError(t, logResp.Error)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestSummarizeExpenses(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating mock database: %v", err)
	}
	defer db.Close()

	values := [][]driver.Value{
		{
			"2021-01", 44.79,
		},
		{
			"2023-06", 9.01,
		},
	}

	rows := sqlmock.NewRows([]string{"month", "total_spent"}).AddRows(values...)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	sumResp := SummarizeExpenses(db, &SummaryRequest{})
	assert.Assert(t, sumResp.Success)
	assert.NilError(t, sumResp.Error)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestDeleteExpense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE").WillReturnResult(sqlmock.NewResult(1, 1))

	deleteResp := DeleteExpense(db, &DeleteRequest{Id: 1})
	assert.Assert(t, deleteResp.Success)
	assert.NilError(t, deleteResp.Error)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
