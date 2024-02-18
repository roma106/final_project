package main

import (
	"backend/server"
)

func main() {
	// db := psql.ConnectToDB("go_projects")
	// defer db.Close()

	// err := psql.Update(db, "yandex_final", &psql.Expr{Expression: "2+2", Status: "failed"})
	// if err != nil {
	// 	fmt.Println(err)
	// }

	server.StartServer(":8080")
}
