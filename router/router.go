package router

import (
	"os"

	"github.com/labstack/echo/v4/middleware"

	_ "net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// Routingを設定する関数　引数はecho.echo型であり、戻り値はerror型
func SetRouter(e *echo.Echo) error {

	// 諸々の設定(*1)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339_nano} ${host} ${method} ${uri} ${status} ${header}\n",
		Output: os.Stdout,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	//新規登録用
	e.POST("/signup", SignUpHandler)
	//ログイン用
	e.POST("/login", LoginHandler)

	//グループ化することでjwt認証を必須にする
	api := e.Group("/api")
	//自身の空き時間に関して、
	api.Use(echojwt.WithConfig(Config))
	api.GET("/events", GetEventsHandler)
	api.POST("/events", AddEventHandler)
	api.PUT("/events/:eventID", ChangeEventHandler)
	api.DELETE("/events/:eventID", DeleteEventHandler)

	// バンド作成POST /api/bands
	api.POST("/bands", AddBandHandler)

	//ユーザが入っているバンド一覧
	api.GET("/bands", GetBandHandler)

	//自分がバンドを抜ける
	api.POST("/bands/:bandID/leave", LeaveBandMemberHandler)

	// バンドメンバー追加POST /api/bands/{bandId}/members
	api.POST("/bands/:bandID/members", AddBandMemberHandler)

	// バンドのメンバ一覧GET /api/bands/{bandId}/members
	api.GET("/bands/:bandID/members", GetBandMembersHandler)

	// バンドの空き時間取得POST /api/bands/{bandId}/availabilities
	api.POST("/bands/:bandID/freetimes", GetBandMembersFreeTimeHandler)

	//バンドをお気に入り登録
	api.POST("/bands/:bandID/favorite", FavoriteBandHandler)

	//バンドのお気に入り削除
	api.DELETE("/bands/:bandID/favorite", RemoveFavoriteBandHandler)

	// 8000番のポートを開く(*2)
	err := e.Start(":8000")
	return err
}
