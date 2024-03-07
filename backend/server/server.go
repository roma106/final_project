package server

import (
	"backend/psql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func StartServer(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/getData", GetData)
	mux.HandleFunc("/postData", PostData)

	fmt.Println("server started")
	http.ListenAndServe(port, AddCorsHeaders(mux))
	return nil
}

func AddCorsHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы со всех источников
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Разрешаем отправку куки и заголовков
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Разрешаем методы и заголовки
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			// Предварительный запрос, возвращаем успешный статус
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := psql.ConnectToDB("go_projects")
	data, err := json.Marshal(psql.GetAll(db, "yandex_final"))
	if err != nil {
		panic(err)
	}
	w.Write(data)
	fmt.Println("data sent to frontend")
	// w.WriteHeader(http.StatusOK)
}

func PostData(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("data posted"))
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Body.Close()
	db := psql.ConnectToDB("go_projects")
	// работа со временем операций - не окончена
	expr := &psql.Expr{Status: "waiting", StartingTime: time.Now(), EndingTime: time.Now().Add(2 * time.Minute), Result: nil}
	err = json.Unmarshal(data, expr)
	if err != nil {
		fmt.Println(err)
	}
	err = psql.Set(db, "yandex_final", expr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("data posted", expr)
}
