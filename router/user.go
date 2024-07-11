package router

import (
	"FindTime-Server/model"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type jwtCustonClaims struct {
	UID  int    `json:"uid"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

var signingKey = []byte("secret")

var Config = echojwt.Config{
	SigningKey: signingKey,
	NewClaimsFunc: func(c echo.Context) jwt.Claims {
		return &jwtCustonClaims{}
	},
}

func SignUpHandler(c echo.Context) error {

	user := new(model.User)

	if err := c.Bind(user); err != nil {
		return err
	}

	//バリデーションチェック
	if user.Name == "" || user.Password == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid name or password",
		}
	}

	if u := model.FindUser(&model.User{Name: user.Name}); u.ID != 0 {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "name already exists",
		}
	}

	//データベースに保存
	model.CreateUser(user)
	user.Password = ""
	return c.JSON(http.StatusCreated, user)

}

func LoginHandler(c echo.Context) error {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	user := model.FindUser(&model.User{Name: u.Name})
	if user.ID == 0 || user.Password != u.Password {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "invalid name or password",
		}
	}

	claims := &jwtCustonClaims{
		user.ID,
		user.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(signingKey)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})

}

func userIDFromToken(c echo.Context) int {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustonClaims)
	uid := claims.UID
	return uid
}
