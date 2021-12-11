package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"strings"
	"testing"
)

func TestAuth(t *testing.T) {
	Convey("Given a DB and App", t, func(){
		app, err := SetupTestAPI()
		So(err, ShouldBeNil)
		Convey("It should authenticate successfully", func(){
			s := `{"username": "TestUser", "password": "TestPass"}`
			req, _ := http.NewRequest("POST", "/token", strings.NewReader(s))
			res, err := app.Test(req, -1)
			So(err, ShouldBeNil)
			var token struct {
				Token string `json:"token"`
			}
			So(unmarshalBody(res.Body, &token), ShouldBeNil)
			So(token.Token, ShouldNotBeBlank)
		})
	})
}
