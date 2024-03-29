package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/andrii-stp/users-crud/config"
	"github.com/andrii-stp/users-crud/model"
	"github.com/andrii-stp/users-crud/storage"

	sq "github.com/Masterminds/squirrel"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserHandler", Ordered, func() {
	cfg, err := config.Load("../test.env")
	if err != nil {
		panic(fmt.Errorf("failed to load config. %w", err))
	}

	db, err := storage.Connect(cfg.Database)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database. %w", err))
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	repo := storage.NewPostgresRepository(logger, db)

	BeforeAll(func() {
		if err := storage.InitDB(db); err != nil {
			panic(fmt.Errorf("failed to create table users. %w", err))
		}
	})

	AfterAll(func() {
		if _, err := db.Exec(`
		DROP TABLE IF EXISTS users;
		`); err != nil {
			panic(fmt.Errorf("failed to drop table users. %w", err))
		}
	})

	var user *model.User

	BeforeEach(func() {
		if _, err := db.Exec("DELETE FROM users;"); err != nil {
			panic(fmt.Errorf("failed to delete users. %w", err))
		}

		user = &model.User{
			UserName:   "JohnDoe",
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "johndoe@yahoo.com",
			Status:     "A",
			Department: "Accounts",
		}

		err = sq.StatementBuilder.Insert("users").
			Columns("user_name", "first_name", "last_name", "email", "user_status", "department").
			Values(user.UserName, user.FirstName, user.LastName, user.Email, user.Status, user.Department).
			Suffix("RETURNING *").PlaceholderFormat(sq.Dollar).RunWith(db).QueryRow().
			Scan(&user.UserID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.Department)

		if err != nil {
			panic(err)
		}
	})

	AfterEach(func() {
		if _, err := db.Exec("DELETE FROM users;"); err != nil {
			panic(fmt.Errorf("failed to delete users. %w", err))
		}
	})

	url := "/api/v1/users"

	Describe("List", func() {
		var resp *httptest.ResponseRecorder

		JustBeforeEach(func() {
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			resp = ExecuteRequest(logger, req, repo)
		})

		Context("should list all users", func() {

			It("status code should be 200", func() {
				Expect(resp.Code).To(Equal(http.StatusOK))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

			It("body should have equivalent values", func() {
				l, err := DeserializeList(resp.Body.String())
				Expect(err).ToNot(HaveOccurred())
				Expect(l).To(HaveLen(1))
				for _, e := range l {
					Expect(e["id"]).To(Equal(float64(user.UserID)))
					Expect(e["user_name"]).To(Equal(user.UserName))
					Expect(e["first_name"]).To(Equal(user.FirstName))
					Expect(e["last_name"]).To(Equal(user.LastName))
					Expect(e["email"]).To(Equal(user.Email))
					Expect(e["user_status"]).To(Equal(user.Status))
					Expect(e["department"]).To(Equal(user.Department))
				}
			})

		})

	})

	Describe("Create", func() {
		var (
			resp         *httptest.ResponseRecorder
			payload      []byte
			expectedUser *model.User
		)

		BeforeEach(func() {
			payload = []byte(`{
				"user_name": "Pikachu",
				"first_name": "Pika",
				"last_name": "Chu",
				"email": "pikachu@yahoo.com",
				"user_status": "I",
				"department": "Pokemon"
			}`)

			if err := json.Unmarshal(payload, &expectedUser); err != nil {
				panic(err)
			}
		})

		JustBeforeEach(func() {
			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
			req.Header.Add("Content-Type", "application/json")
			resp = ExecuteRequest(logger, req, repo)
		})

		Context("should create a user correctly", func() {

			It("status code should be 201", func() {
				Expect(resp.Code).To(Equal(http.StatusCreated))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

			It("body should have equivalent values", func() {
				e, _ := Deserialize(resp.Body.String())
				Expect(e["id"]).ToNot(Equal(float64(user.UserID)))
				Expect(e["user_name"]).To(Equal(expectedUser.UserName))
				Expect(e["first_name"]).To(Equal(expectedUser.FirstName))
				Expect(e["last_name"]).To(Equal(expectedUser.LastName))
				Expect(e["email"]).To(Equal(expectedUser.Email))
				Expect(e["user_status"]).To(Equal(expectedUser.Status))
				Expect(e["department"]).To(Equal(expectedUser.Department))
			})

		})

		Context("should get an error when create a user without user_name", func() {

			BeforeEach(func() {
				payload = []byte(`{
					"first_name": "Pika",
					"last_name": "Chu",
					"email": "pikachu@yahoo.com",
					"user_status": "I",
					"department": "Pokemon"
				}`)
			})

			It("status code should be 400", func() {
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

		})

		Context("should get an error when create a user with invalid email", func() {

			BeforeEach(func() {
				payload = []byte(`{
					"user_name": "Pikachu",
					"first_name": "Pika",
					"last_name": "Chu",
					"email": "pikachu@invalid-domain",
					"user_status": "I",
					"department": "Pokemon"
				}`)
			})

			It("status code should be 400", func() {
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

		})

		Context("should get an error when create a user with invalid status", func() {

			BeforeEach(func() {
				payload = []byte(`{
					"user_name": "Pikachu",
					"first_name": "Pika",
					"last_name": "Chu",
					"email": "pikachu@yahoo.com",
					"user_status": "_",
					"department": "Pokemon"
				}`)
			})

			It("status code should be 400", func() {
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

		})

	})

	Describe("Update", func() {
		var (
			resp         *httptest.ResponseRecorder
			payload      []byte
			expectedUser *model.User
			id           int64
		)

		BeforeEach(func() {
			payload = []byte(`{
				"user_name": "Pikachu",
				"first_name": "Pika",
				"last_name": "Chu",
				"email": "pikachu@yahoo.com",
				"user_status": "I",
				"department": "Pokemon"
			}`)

			if err := json.Unmarshal(payload, &expectedUser); err != nil {
				panic(err)
			}

			id = user.UserID
		})

		JustBeforeEach(func() {
			path := fmt.Sprintf("%s/%d", url, id)

			req, _ := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(payload))
			req.Header.Add("Content-Type", "application/json")
			resp = ExecuteRequest(logger, req, repo)
		})

		Context("should update a user correctly", func() {

			It("status code should be 200", func() {
				Expect(resp.Code).To(Equal(http.StatusOK))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

			It("body should have equivalent values", func() {
				e, _ := Deserialize(resp.Body.String())
				Expect(e["id"]).To(Equal(float64(user.UserID)))
				Expect(e["user_name"]).To(Equal(expectedUser.UserName))
				Expect(e["first_name"]).To(Equal(expectedUser.FirstName))
				Expect(e["last_name"]).To(Equal(expectedUser.LastName))
				Expect(e["email"]).To(Equal(expectedUser.Email))
				Expect(e["user_status"]).To(Equal(expectedUser.Status))
				Expect(e["department"]).To(Equal(expectedUser.Department))
			})

		})

		Context("should get an error when update a user that does not exists", func() {

			BeforeEach(func() {
				id = -1
				payload = []byte(`{
					"user_name": "Pikachu",
					"first_name": "Pika",
					"last_name": "Chu",
					"email": "pikachu@yahoo.com",
					"user_status": "I",
					"department": "Pokemon"
				}`)
			})

			It("status code should be 404", func() {
				Expect(resp.Code).To(Equal(http.StatusNotFound))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

		})

		Context("should get an error when update a user with existing user_name", func() {

			BeforeEach(func() {
				u := &model.User{
					UserName:   "Kirby",
					FirstName:  "Kir",
					LastName:   "By",
					Email:      "kirby@yahoo.com",
					Status:     "T",
					Department: "Explorer",
				}

				err = sq.StatementBuilder.Insert("users").
					Columns("user_name", "first_name", "last_name", "email", "user_status", "department").
					Values(u.UserName, u.FirstName, u.LastName, u.Email, u.Status, u.Department).
					Suffix("RETURNING *").PlaceholderFormat(sq.Dollar).RunWith(db).QueryRow().
					Scan(&u.UserID, &u.UserName, &u.FirstName, &u.LastName, &u.Email, &u.Status, &u.Department)

				if err != nil {
					panic(err)
				}

				id = user.UserID
				payload = []byte(`{
					"user_name": "Kirby",
					"first_name": "Pika",
					"last_name": "Chu",
					"email": "pikachu@yahoo.com",
					"user_status": "I",
					"department": "Pokemon"
				}`)
			})

			It("status code should be 409", func() {
				Expect(resp.Code).To(Equal(http.StatusConflict))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

		})

		Context("should get an error when update a user without user_name", func() {

			BeforeEach(func() {
				payload = []byte(`{
					"first_name": "Pika",
					"last_name": "Chu",
					"email": "pikachu@yahoo.com",
					"user_status": "I",
					"department": "Pokemon"
				}`)
			})

			It("status code should be 400", func() {
				Expect(resp.Code).To(Equal(400))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

		})

		Context("should get an error when update a user with invalid email", func() {

			BeforeEach(func() {
				payload = []byte(`{
					"user_name": "Pikachu",
					"first_name": "Pika",
					"last_name": "Chu",
					"email": "pikachu@invalid-domain",
					"user_status": "I",
					"department": "Pokemon"
				}`)
			})

			It("status code should be 400", func() {
				Expect(resp.Code).To(Equal(400))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

		})

		Context("should get an error when update a user with invalid status", func() {

			BeforeEach(func() {
				payload = []byte(`{
					"user_name": "Pikachu",
					"first_name": "Pika",
					"last_name": "Chu",
					"email": "pikachu@yahoo.com",
					"user_status": "_",
					"department": "Pokemon"
				}`)
			})

			It("status code should be 400", func() {
				Expect(resp.Code).To(Equal(400))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

		})

		Context("should get a 400 response when sending request with invalid id", func() {

			JustBeforeEach(func() {
				path := url + "/invalid"
				req, _ := http.NewRequest(http.MethodDelete, path, nil)
				resp = ExecuteRequest(logger, req, repo)
			})

			It("status code should be 400", func() {
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

		})

	})

	Describe("Delete", func() {
		var resp *httptest.ResponseRecorder

		BeforeEach(func() {
			u := &model.User{
				UserName:   "Kirby",
				FirstName:  "Kir",
				LastName:   "By",
				Email:      "kirby@yahoo.com",
				Status:     "T",
				Department: "Explorer",
			}

			err = sq.StatementBuilder.Insert("users").
				Columns("user_name", "first_name", "last_name", "email", "user_status", "department").
				Values(u.UserName, u.FirstName, u.LastName, u.Email, u.Status, u.Department).
				Suffix("RETURNING *").PlaceholderFormat(sq.Dollar).RunWith(db).QueryRow().
				Scan(&u.UserID, &u.UserName, &u.FirstName, &u.LastName, &u.Email, &u.Status, &u.Department)

			if err != nil {
				panic(err)
			}
		})

		JustBeforeEach(func() {
			path := fmt.Sprintf("%s/%d", url, user.UserID)

			req, _ := http.NewRequest(http.MethodDelete, path, nil)
			req.Header.Add("Content-Type", "application/json")
			resp = ExecuteRequest(logger, req, repo)
		})

		Context("should delete a user correctly", func() {

			It("status code should be 204", func() {
				Expect(resp.Code).To(Equal(http.StatusNoContent))
			})

			It("body should be nil", func() {
				Expect(resp.Body).To(BeAssignableToTypeOf(&bytes.Buffer{}))
			})

		})

		Context("should get a 400 response when sending request with invalid id", func() {

			JustBeforeEach(func() {
				path := url + "/invalid"
				req, _ := http.NewRequest(http.MethodDelete, path, nil)
				resp = ExecuteRequest(logger, req, repo)
			})

			It("status code should be 400", func() {
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

		})

		Context("should get a 404 response", func() {

			JustBeforeEach(func() {
				path := url + "-1"
				req, _ := http.NewRequest(http.MethodDelete, path, nil)
				resp = ExecuteRequest(logger, req, repo)
			})

			It("status code should be 404", func() {
				Expect(resp.Code).To(Equal(http.StatusNotFound))
			})

			It("body should not be nil", func() {
				Expect(resp.Body).ToNot(BeNil())
			})

		})

	})

})
