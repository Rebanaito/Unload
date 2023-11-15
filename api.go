package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

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

func (server *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/", DefaultPage)
	router.HandleFunc("/login_page", LoginPage)
	router.HandleFunc("/login", Login(server))
	router.HandleFunc("/register_page", RegisterPage)
	router.HandleFunc("/register_employer", RegisterEmployer)
	router.HandleFunc("/register_worker", RegisterWorker)
	router.HandleFunc("/register", Register(server))
	router.HandleFunc("/home", Home(server))
	router.HandleFunc("/me", ProfileInfo(server))
	router.HandleFunc("/tasks", TaskInfo(server))

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

func RegisterEmployer(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.New("register").Parse(registerEmployer)
	tmpl.Execute(w, nil)
}

func RegisterWorker(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.New("register").Parse(registerWorker)
	tmpl.Execute(w, nil)
}

func Login(server *APIServer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username := r.FormValue("username")
		password := r.FormValue("password")

		role := server.storage.AuthUser(username, password)

		if role == "" {
			w.WriteHeader(http.StatusUnauthorized)
			message := fmt.Sprintf(loginError, "Wrong username/password")
			tmpl, _ := template.New("badCredentials").Parse(message)
			tmpl.Execute(w, nil)
		} else {

			token, err := CreateJWT(username, role)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				message := fmt.Sprintf(loginError, "Something went wrong, try again")
				tmpl, _ := template.New("token").Parse(message)
				tmpl.Execute(w, nil)
			} else {
				w.WriteHeader(http.StatusOK)

				var tmpl *template.Template
				switch role {
				case "employer":
					tmpl, _ = template.New("home").Parse(fmt.Sprintf(homeEmployer, username, token))
				case "worker":
					tmpl, _ = template.New("home").Parse(fmt.Sprintf(homeWorker, username, token, token))
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

		var err error
		switch role {
		case "employer":
			cash, _ := strconv.Atoi(r.FormValue("cash"))
			err = server.storage.CreateEmployer(username, password, cash)
		case "worker":
			weight, _ := strconv.Atoi(r.FormValue("weight"))
			wage, _ := strconv.Atoi(r.FormValue("wage"))
			drinks := r.FormValue("drinks")
			err = server.storage.CreateWorker(username, password, weight, wage, drinks == "true")
		}

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

func Home(server *APIServer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, role, token := ValidateJWT(r, "home")
		if username == "" {
			w.WriteHeader(http.StatusUnauthorized)
			tmpl, _ := template.New("badCredentials").Parse(unauthorizedAccess)
			tmpl.Execute(w, nil)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("token", token)
			var page string
			switch role {
			case "employer":
				page = fmt.Sprintf(homeEmployer, username, token)
			case "worker":
				page = fmt.Sprintf(homeWorker, username, token, token)
			}
			tmpl, _ := template.New("home").Parse(page)
			tmpl.Execute(w, nil)
		}
	})
}

func ProfileInfo(server *APIServer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, role, token := ValidateJWT(r, "home")
		if username == "" {
			w.WriteHeader(http.StatusUnauthorized)
			tmpl, _ := template.New("badCredentials").Parse(unauthorizedAccess)
			tmpl.Execute(w, nil)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("token", token)
			var page string
			switch role {
			case "employer":
				page = getEmployerInfo(server, username, token)
			case "worker":
				page = getWorkerInfo(server, username, token)
			}
			tmpl, _ := template.New("home").Parse(page)
			tmpl.Execute(w, nil)
		}
	})
}

func TaskInfo(server *APIServer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, role, token := ValidateJWT(r, "home")
		if username == "" {
			w.WriteHeader(http.StatusUnauthorized)
			tmpl, _ := template.New("badCredentials").Parse(unauthorizedAccess)
			tmpl.Execute(w, nil)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("token", token)
			var page string
			switch role {
			case "employer":
				page = getEmployerTasks(server, username, token)
			case "worker":
				page = getWorkerTasks(server, username, token)
			}
			tmpl, _ := template.New("home").Parse(page)
			tmpl.Execute(w, nil)
		}
	})
}

func getEmployerInfo(server *APIServer, username, token string) string {
	employer, workers, err := server.storage.GetEmployer(username)
	if err != nil {
		return err.Error()
	}
	var builder strings.Builder
	for i, worker := range workers {
		if i == 0 {
			builder.WriteString(`<table>
								<tr>
									<th>Wage</th>
									<th>Fatigue</th>
									<th>Max weight</th>
									<th>Alcoholism</th>
								</tr>`)
		}
		builder.WriteString(fmt.Sprintf(`<tr><td>%d</td>
							<td>%d</td>
							<td>%d</td>
							<td>%v</td></tr>`, worker.wage, worker.fatigue, worker.weight, worker.drinks))
		if i == len(workers)-1 {
			builder.WriteString(`</table>`)
		}
	}
	if builder.Len() == 0 {
		builder.WriteString("No registered workers")
	}
	return fmt.Sprintf(profileEmployer, username, token, employer.cash, builder.String())
}

func getWorkerInfo(server *APIServer, username, token string) string {
	worker, err := server.storage.GetWorker(username)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf(profileWorker,
		username,
		token,
		worker.userid,
		worker.wage,
		worker.fatigue,
		worker.weight,
		worker.drinks)
}

func getWorkerTasks(server *APIServer, username, token string) string {
	tasks := server.storage.GetWorkerTasks(username)
	var builder strings.Builder
	for i, task := range tasks {
		if i == 0 {
			builder.WriteString(`<table>
								<tr>
									<th>Task ID</th>
									<th>Employer</th>
									<th>Weight</th>
								</tr>`)
		}
		builder.WriteString(fmt.Sprintf(`<tr>
											<td>%d</td>
											<td>%d</td>
											<td>%d</td>
										</tr>`, task.taskID, task.employer, task.weight))
		if i == len(tasks)-1 {
			builder.WriteString(`</table>`)
		}
	}
	if builder.Len() == 0 {
		builder.WriteString("No completed tasks")
	}
	return fmt.Sprintf(tasksWorker, username, token, builder.String())
}

func getEmployerTasks(server *APIServer, username, token string) string {
	tasks := server.storage.GetEmployerTasks(username)
	var builder strings.Builder
	for i, task := range tasks {
		if i == 0 {
			builder.WriteString(`<table>
								<tr>
									<th>Task ID</th>
									<th>Weight</th>
								</tr>`)
		}
		builder.WriteString(fmt.Sprintf(`<tr>
											<td>%d</td>
											<td>%d</td>
										</tr>`, task.taskID, task.weight))
		if i == len(tasks)-1 {
			builder.WriteString(`</table>`)
		}
	}
	if builder.Len() == 0 {
		builder.WriteString("No available tasks")
	}
	return fmt.Sprintf(tasksEmployer, username, token, builder.String())
}
