package main

import (
	"app/domain/entity"
	"app/infra/persistence"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func main() {
	err := persistence.InitDBConnection()
	if err != nil {
		panic("DB接続エラー")
	}

	up := persistence.NewUserPersistence()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		u := entity.NewUser()
		u, _ = up.Create(u)
		return c.String(http.StatusOK, "Hello, World!:"+strconv.Itoa(u.ID))
	})
	e.Logger.Fatal(e.Start(":1323"))
}
