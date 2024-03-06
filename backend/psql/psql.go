package psql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Knetic/govaluate"
	_ "github.com/lib/pq"
)

type Expr struct {
	ID           int
	Expression   string `json:"expression"`
	Status       string
	Result       interface{}
	StartingTime time.Time
	EndingTime   time.Time
}

func ConnectToDB(dbInput string) *sql.DB {
	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		"postgres",
		"ri106rom",
		"postgres", // обновленное имя контейнера PostgreSQL
		"5432",     // обновленный порт PostgreSQL
		"go_projects")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	// err = db.Ping()
	// if err != nil {
	// 	fmt.Println("Не удалось подключиться к базе данных")
	// 	panic(err)
	// }
	return db
}

func GetAll(db *sql.DB, table string) []Expr {
	rows, err := db.Query("SELECT * FROM " + table)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	out := make([]Expr, 0)
	for i := 0; rows.Next(); i++ {
		out = append(out, Expr{})
		if err := rows.Scan(&out[i].Expression, &out[i].Status, &out[i].Result, &out[i].StartingTime, &out[i].EndingTime, &out[i].ID); err != nil {
			panic(err)
		}
		out[i].Status, out[i].Result = CheckTiming(db, out[i])
	}
	return out
}

func Set(db *sql.DB, table string, ex *Expr) error {
	query := "INSERT INTO " + table + " (\"Expression\", \"Status\", \"StartingTime\", \"EndingTime\") VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(query, ex.Expression, ex.Status, ex.StartingTime, ex.EndingTime)
	if err != nil {
		return err
	}
	fmt.Println("Expression " + ex.Expression + " inserted")
	return nil
}

// func Update(db *sql.DB, table string, ex *Expr) error {
// 	query := "UPDATE " + table + " SET \"Result\" = $1, \"Status\" = $2 WHERE \"Expression\" = $3"
// 	_, err := db.Exec(query, ex.Result, ex.Status, ex.Expression)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("Expression " + ex.Expression + " updated")
// 	return nil
// }

func CheckTiming(db *sql.DB, ex Expr) (string, interface{}) {
	if !ex.EndingTime.After(time.Now()) {
		updateStmt, err := db.Prepare("UPDATE yandex_final SET \"Result\" = $1, \"Status\" = $2 WHERE \"ID\" = $3")
		if err != nil {
			fmt.Println("Ошибка при подготовке запроса обновления:", err)
			return "failed", nil
		}
		defer updateStmt.Close()
		res, err := Calculate(ex.Expression)
		if err != nil {
			panic(err)
		}
		_, err = updateStmt.Exec(res, "done", ex.ID)
		if err != nil {
			fmt.Println("Ошибка при выполнении запроса обновления:", err)
			return "failed", nil
		}

		fmt.Println("Данные успешно обновлены")
		return "done", res
	} else {
		return "waiting", nil
	}
}

func Calculate(expr string) (float64, error) {
	expression, _ := govaluate.NewEvaluableExpression(expr)
	result, err := expression.Evaluate(nil)
	if err != nil {
		return 0, fmt.Errorf("ошибка вычисления: %v", err)
	}
	return result.(float64), nil
}
