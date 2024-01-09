package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http/httptest"
	"testing"
)

const (
	testdburi = "mongodb://localhost:27017"
	dbname    = "hotel-reservation-test"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		t.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client, dbname),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "First",
		LastName:  "Last",
		Email:     "test@test.com",
		Password:  "testpass123",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		t.Error(err)
	}
	if len(user.ID) == 0 {
		t.Error("expected user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Error("expecting the EncryptedPassword not to be included in the json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected %s got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected %s got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected %s got %s", params.Email, user.Email)
	}
}
