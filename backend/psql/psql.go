package psql

import (
	"database/sql"
	"fmt"
	"log"
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
	db, err := sql.Open("postgres", "host=postgres port=5432 user=postgres password=ri106rom dbname=go_projects sslmode=disable")
	if err != nil {
		log.Printf("[ERROR]: Spawn connection to database was failed: %v", err)
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Не удалось подключиться к базе данных")
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS yandex_final (
			expression character varying COLLATE pg_catalog."default" NOT NULL,
			status character varying COLLATE pg_catalog."default" NOT NULL,
			result integer,
			startingTime timestamp with time zone NOT NULL,
			endingTime timestamp with time zone NOT NULL,
			id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 1000 CACHE 1 ),
			CONSTRAINT yandex_final_pkey PRIMARY KEY (id)
	)`)
	if err != nil {
		fmt.Println("Не удалось создать таблицу")
		panic(err)
	}
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
	query := "INSERT INTO " + table + " (expression, status, startingTime, endingTime) VALUES ($1, $2, $3, $4)"
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
		updateStmt, err := db.Prepare("UPDATE yandex_final SET result = $1, status = $2 WHERE id = $3")
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
