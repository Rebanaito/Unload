package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIServer struct {
	listenAddr string
	storage    Storage
}

func NewAPIServer(listenAddr string, conn *pgxpool.Pool) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		storage:    PostgreSQL{conn: conn},
	}
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func (server *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/", DefaultPage)
	router.HandleFunc("/login_page", LoginPage)
	router.HandleFunc("/login", Login(server))
	router.HandleFunc("/register_page", RegisterPage)
	router.HandleFunc("/register", Register(server))
	// router.HandleFunc("/register", makeHTTPHandleFunc(server.Register))
	// router.HandleFunc("/home", Home(server))
	// router.HandleFunc("/me", ProfileInfo(server))
	// router.HandleFunc("/tasks", TaskInfo(server))

	log.Println("Server running on:", server.listenAddr)

	go func() { http.ListenAndServe(server.listenAddr, router) }()

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt)
	<-kill

	log.Println("Shutting the server down")
}

func DefaultPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.New("default").Parse(defaultPage)
	tmpl.Execute(w, nil)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.New("login").Parse(loginPage)
	tmpl.Execute(w, nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.New("register").Parse(registerPage)
	tmpl.Execute(w, nil)
}

func Login(server *APIServer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username := r.FormValue("username")
		password := r.FormValue("password")

		id, role := server.storage.AuthUser(username, password)

		if id == -1 {
			w.WriteHeader(http.StatusUnauthorized)
			message := fmt.Sprintf(loginError, "Wrong username/password")
			tmpl, _ := template.New("badCredentials").Parse(message)
			tmpl.Execute(w, nil)
		} else {

			token, err := CreateJWT(id, role)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				message := fmt.Sprintf(loginError, "Something went wrong, try again")
				tmpl, _ := template.New("badCredentials").Parse(message)
				tmpl.Execute(w, nil)
			} else {
				http.SetCookie(w, &http.Cookie{
					Name:  "token",
					Value: token},
				)
				w.WriteHeader(http.StatusOK)

				var tmpl *template.Template
				switch role {
				case "employer":
					tmpl, _ = template.New("home").Parse(homeEmployer)
				case "worker":
					tmpl, _ = template.New("home").Parse(homeWorker)
				}

				tmpl.Execute(w, nil)
			}
		}
	})
}

func Register(server *APIServer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username := r.FormValue("username")
		password := r.FormValue("password")
		role := r.FormValue("role")

		err := server.storage.CreateUser(username, password, role)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			message := fmt.Sprintf(loginError, err.Error())
			tmpl, _ := template.New("badCredentials").Parse(message)
			tmpl.Execute(w, nil)
		} else {
			w.WriteHeader(http.StatusOK)
			tmpl, _ := template.New("login").Parse(loginPage)
			tmpl.Execute(w, nil)
		}
	})
}
