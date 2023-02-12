package contentful

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMembershipsService_List(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/space_memberships")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("membership.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	it, err := cma.Memberships.List(context.Background(), spaceID)
	assertions.Nil(err)
	collection, err := it.Next()
	assertions.Nil(err)
	membership := collection.Items
	assertions.Equal(2, len(membership))
	assertions.Equal("test@contentfulsdk.go", membership[0].Email)
}

func TestMembershipsService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/space_memberships/0xWanD4AZI2AR35wW9q51n")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("membership_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	membership, err := cma.Memberships.Get(context.Background(), spaceID, "0xWanD4AZI2AR35wW9q51n")
	assertions.Nil(err)
	assertions.Equal("0xWanD4AZI2AR35wW9q51n", membership.Sys.ID)
}

func TestMembershipsService_Get_2(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/space_memberships/0xWanD4AZI2AR35wW9q51n")

		checkHeaders(r, assertions)

		w.WriteHeader(400)
		_, _ = fmt.Fprintln(w, readTestData("membership_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Memberships.Get(context.Background(), spaceID, "0xWanD4AZI2AR35wW9q51n")
	assertions.Nil(err)
}

func TestMembershipsService_Upsert_Create(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/space_memberships")

		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)

		email := payload["email"].(string)
		admin := payload["admin"].(bool)
		assertions.Equal("johndoe@nonexistent.com", email)
		assertions.Equal(true, admin)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("membership_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	membership := &Membership{
		Admin: true,
		Roles: []Roles{
			{
				Sys: &Sys{
					ID:       "1ElgCn1mi1UHSBLTP2v4TD",
					Type:     "Link",
					LinkType: "Role",
				},
			},
		},
		Email: "johndoe@nonexistent.com",
	}

	err = cma.Memberships.Upsert(context.Background(), spaceID, membership)
	assertions.Nil(err)
}

func TestMembershipsService_Upsert_Update(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/space_memberships/0xWanD4AZI2AR35wW9q51n")

		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)

		email := payload["email"].(string)
		assertions.Equal("editedmail@examplemail.com", email)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("membership_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	membership, err := membershipFromTestData("membership_1.json")
	assertions.Nil(err)

	membership.Email = "editedmail@examplemail.com"

	err = cma.Memberships.Upsert(context.Background(), spaceID, membership)
	assertions.Nil(err)
}

func TestMembershipsService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/space_memberships/0xWanD4AZI2AR35wW9q51n")
		checkHeaders(r, assertions)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test role
	membership, err := membershipFromTestData("membership_1.json")
	assertions.Nil(err)

	// delete role
	err = cma.Memberships.Delete(context.Background(), spaceID, membership.Sys.ID)
	assertions.Nil(err)
}
