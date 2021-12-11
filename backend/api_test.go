package main

import (
	"fmt"
	"github.com/michaelbironneau/shoppings/backend/api"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	Convey("Given a DB and App", t, func() {
		app, _, err := SetupTestAPI()
		So(err, ShouldBeNil)
		Convey("It should authenticate successfully", func() {
			s := `{"username": "TestUser", "password": "TestPass"}`
			res, err := makeTestRequest(app, "POST", "/token", "", s)
			So(err, ShouldBeNil)
			var token struct {
				Token string `json:"token"`
			}
			So(unmarshalBody(res.Body, &token), ShouldBeNil)
			So(token.Token, ShouldNotBeBlank)
			s = `{"username": "TestUser", "password": "Incorrect"}`
			res, err = makeTestRequest(app, "POST", "/token", "", s)
			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, 401)
			token.Token = ""
			unmarshalBody(res.Body, &token)
			So(token.Token, ShouldBeBlank)
		})
	})
}

func TestItems(t *testing.T) {
	Convey("Given a test app", t, func() {
		app, token, err := SetupTestAPI()
		So(err, ShouldBeNil)
		Convey("Given a test list", func() {
			var l = api.List{
				Name: "Test List",
			}
			// create the list
			res, err := makeTestRequest(app, "POST", "/lists", token, l)
			So(err, ShouldBeNil)
			var listResp struct {
				ID string `json:"id"`
			}
			So(unmarshalBody(res.Body, &listResp), ShouldBeNil)
			So(listResp.ID, ShouldNotBeBlank)
			Convey("It should add an unknown item and update it", func() {
				item := api.ListItem{
					ListID:   listResp.ID,
					Name:     "Test Item",
					Quantity: 1,
				}
				update := api.ListUpdate{Updates: []api.ListItem{item}}
				res, err = makeTestRequest(app, "PATCH", fmt.Sprintf("/lists/%s/updates", listResp.ID), token, update)
				So(err, ShouldBeNil)

				//1. Get all list items
				var items api.ListUpdate
				res, err = makeTestRequest(app, "GET", fmt.Sprintf("/lists/%s/items", listResp.ID), token, nil)
				So(err, ShouldBeNil)
				So(unmarshalBody(res.Body, &items), ShouldBeNil)
				So(items.Updates, ShouldHaveLength, 1)
				So(items.Updates[0].Name, ShouldEqual, item.Name)
				So(items.Updates[0].Quantity, ShouldEqual, 1)
				So(items.UpdateTime, ShouldBeGreaterThan, 0)

				//2. Get updated since in past
				var updates api.ListUpdate
				tt := time.Now().Unix() - int64(time.Hour/time.Second)
				res, err = makeTestRequest(app, "GET", fmt.Sprintf("/lists/%s/updates/%v", listResp.ID, tt), token, nil)
				So(err, ShouldBeNil)
				So(unmarshalBody(res.Body, &updates), ShouldBeNil)
				So(updates.Updates, ShouldHaveLength, 1)
				So(updates.Updates[0].Name, ShouldEqual, item.Name)
				So(updates.UpdateTime, ShouldBeGreaterThan, tt)

				//3. Get list name and check summary
				var lists []api.List
				res, err = makeTestRequest(app, "GET", "/lists/", token, nil)
				So(err, ShouldBeNil)
				So(unmarshalBody(res.Body, &lists), ShouldBeNil)
				So(lists, ShouldHaveLength, 1)
				So(lists[0].Name, ShouldEqual, l.Name)
				So(lists[0].Summary, ShouldNotBeBlank)

				//4. Get updated since with time in future - expect no updates
				tt = updates.UpdateTime + 1
				res, err = makeTestRequest(app, "GET", fmt.Sprintf("/lists/%s/updates/%v", listResp.ID, tt), token, nil)
				So(err, ShouldBeNil)
				So(unmarshalBody(res.Body, &updates), ShouldBeNil)
				So(updates.Updates, ShouldHaveLength, 0)

				//5. Create Item with this name, check that our ListItem has been updated to match
				ii := api.Item{Name: strings.ToLower(item.Name)} // try in a different case, it should still work
				res, err = makeTestRequest(app, "POST", "/items", token, ii)
				So(err, ShouldBeNil)
				res, err = makeTestRequest(app, "GET", fmt.Sprintf("/lists/%s/items", listResp.ID), token, nil)
				So(err, ShouldBeNil)
				So(unmarshalBody(res.Body, &updates), ShouldBeNil)
				So(updates.Updates, ShouldHaveLength, 1)
				So(updates.Updates[0].ItemID, ShouldNotBeBlank)
				So(updates.Updates[0].Name, ShouldEqual, item.Name)

				//6. Add Item with new name
				ii.Name = "New Item"
				res, err = makeTestRequest(app, "POST", "/items", token, ii)
				So(err, ShouldBeNil)
				var newItemID struct {
					ID string `json:"id"`
				}
				So(unmarshalBody(res.Body, &newItemID), ShouldBeNil)
				item.Name = ii.Name
				update = api.ListUpdate{Updates: []api.ListItem{item}}
				res, err = makeTestRequest(app, "PATCH", fmt.Sprintf("/lists/%s/updates", listResp.ID), token, update)
				So(err, ShouldBeNil)

				//7. Receive new update
				res, err = makeTestRequest(app, "GET", fmt.Sprintf("/lists/%s/items", listResp.ID), token, nil)
				So(err, ShouldBeNil)
				So(unmarshalBody(res.Body, &updates), ShouldBeNil)
				So(updates.Updates, ShouldHaveLength, 2)
			})
		})
	})
}

func TestStores(t *testing.T) {
	Convey("Given a test app", t, func() {
		app, token, err := SetupTestAPI()
		So(err, ShouldBeNil)
		Convey("It should list stores", func() {
			res, err := makeTestRequest(app, "GET", "/stores", token, nil)
			So(err, ShouldBeNil)
			var stores []api.Store
			So(unmarshalBody(res.Body, &stores), ShouldBeNil)
			So(stores, ShouldHaveLength, 2)
			So(stores[0].Name, ShouldNotBeBlank)
			So(stores[1].Name, ShouldNotBeBlank)
			So(stores[0].ID, ShouldNotBeBlank)
			So(stores[1].ID, ShouldNotBeBlank)
		})
	})
}

func TestLists(t *testing.T) {
	Convey("Given a test app", t, func() {
		app, token, err := SetupTestAPI()
		So(err, ShouldBeNil)
		Convey("It should create and retrieve a new list", func() {
			var l = api.List{
				Name: "Test List",
			}
			// create the list
			res, err := makeTestRequest(app, "POST", "/lists", token, l)
			So(err, ShouldBeNil)

			// retrieve the list again
			res, err = makeTestRequest(app, "GET", "/lists", token, nil)
			var lists []api.List
			So(unmarshalBody(res.Body, &lists), ShouldBeNil)
			So(lists, ShouldHaveLength, 1)
			So(lists[0].Name, ShouldEqual, "Test List")
			So(lists[0].Archived, ShouldBeFalse)
			So(lists[0].ID, ShouldNotBeBlank)

			// archive the list
			res, err = makeTestRequest(app, "POST", fmt.Sprintf("/lists/%s/archive", lists[0].ID), token, nil)
			So(err, ShouldBeNil)

			// check that the list no longer comes up
			res, err = makeTestRequest(app, "GET", "/lists", token, nil)
			So(err, ShouldBeNil)
			testListID := lists[0].ID
			lists = []api.List{}
			So(unmarshalBody(res.Body, &lists), ShouldBeNil)
			So(lists, ShouldHaveLength, 0)

			//unarchive the list
			res, err = makeTestRequest(app, "POST", fmt.Sprintf("/lists/%s/unarchive", testListID), token, nil)
			So(err, ShouldBeNil)

			// check that it comes up again
			res, err = makeTestRequest(app, "GET", "/lists", token, nil)
			So(err, ShouldBeNil)
			lists = []api.List{}
			So(unmarshalBody(res.Body, &lists), ShouldBeNil)
			So(lists, ShouldHaveLength, 1)
		})
	})
}
