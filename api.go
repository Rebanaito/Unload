package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

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
	router.HandleFunc("/start", Play(server))
	router.HandleFunc("/play", Attempt(server))

	log.Println("Server running on:", server.listenAddr)

	go func() { http.ListenAndServe(server.listenAddr, router) }()
	go taskGenerator(server)

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt)
	<-kill

	log.Println("Shutting the server down")
}

// Self-explanatory. Adds tasks every minute as long as there are < 10 active tasks
func taskGenerator(server *APIServer) {
	for {
		count := server.storage.GetActiveTaskCount()
		if count < 10 {
			n := rand.Intn(10)
			for i := 0; i < n; i++ {
				weight := rand.Intn(70) + 10
				server.storage.AddTask(weight)
			}
		}
		time.Sleep(time.Minute)
	}
}

// Navigational pages
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
			if role == "employer" {
				employer, _, _ := server.storage.GetEmployer(username)
				if employer.cash < 0 {
					w.WriteHeader(http.StatusUnauthorized)
					message := fmt.Sprintf(loginError, "Employer is bankrupt")
					tmpl, _ := template.New("token").Parse(message)
					tmpl.Execute(w, nil)
				}
			}
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
					tmpl, _ = template.New("home").Parse(fmt.Sprintf(homeEmployer, username, token, token, token))
				case "worker":
					tmpl, _ = template.New("home").Parse(fmt.Sprintf(homeWorker, username, token, token))
				}
				tmpl.Execute(w, nil)
			}
		}
	})
}

// TODO: implement password hashing
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
			var page string
			switch role {
			case "employer":
				page = fmt.Sprintf(homeEmployer, username, token, token, token)
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

func Play(server *APIServer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, role, token := ValidateJWT(r, "home")
		if username == "" || role != "employer" {
			w.WriteHeader(http.StatusUnauthorized)
			tmpl, _ := template.New("badCredentials").Parse(unauthorizedAccess)
			tmpl.Execute(w, nil)
		} else {
			w.WriteHeader(http.StatusOK)
			tasks := server.storage.GetEmployerTasks()
			workers := server.storage.GetAllWorkers()
			game := getGame(tasks, workers, token)
			page := fmt.Sprintf(playGame, username, token, token, token, game)
			tmpl, _ := template.New("home").Parse(page)
			tmpl.Execute(w, nil)
		}
	})
}

func Attempt(server *APIServer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, role, token := ValidateJWT(r, "home")
		if username == "" || role != "employer" {
			w.WriteHeader(http.StatusUnauthorized)
			tmpl, _ := template.New("badCredentials").Parse(unauthorizedAccess)
			tmpl.Execute(w, nil)
		} else {
			r.ParseForm()
			taskID := r.Form["task"][0]
			workerIDS := r.Form["worker"]
			task := server.storage.GetTask(taskID)
			workers := make([]Worker, len(workerIDS))
			var totalMax int
			var totalWage int
			for i, workerID := range workerIDS {
				workers[i] = server.storage.GetWorkerByID(workerID)
				totalMax += (workers[i].weight * (100 - workers[i].fatigue)) / 100
				totalWage += workers[i].wage
			}
			var message string
			if totalMax >= task.weight {
				message = "Task successful!"
				server.UpdateWin(task, workers, username, totalWage)
			} else {
				message = "Task failed!"
				server.UpdateFail(task, workers, username, totalWage, totalMax)
			}
			page := fmt.Sprintf(attemptResult, username, message, token)
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
			builder.WriteString(`<h3>Available workers</h3><table>
								<tr>
									<th>ID</th>
									<th>Wage</th>
									<th>Fatigue</th>
									<th>Max weight</th>
									<th>Alcoholism</th>
								</tr>`)
		}
		builder.WriteString(fmt.Sprintf(`<tr>
							<td>%d</td>
							<td>%d</td>
							<td>%d</td>
							<td>%d</td>
							<td>%v</td></tr>`, worker.userid, worker.wage, worker.fatigue, worker.weight, worker.drinks))
		if i == len(workers)-1 {
			builder.WriteString(`</table>`)
		}
	}
	if builder.Len() == 0 {
		builder.WriteString("No registered workers")
	}
	return fmt.Sprintf(profileEmployer, username, token, token, token, employer.cash, builder.String())
}

func getWorkerInfo(server *APIServer, username, token string) string {
	worker, err := server.storage.GetWorker(username)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf(profileWorker,
		username,
		token,
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
	return fmt.Sprintf(tasksWorker, username, token, token, builder.String())
}

func getEmployerTasks(server *APIServer, username, token string) string {
	tasks := server.storage.GetEmployerTasks()
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
	return fmt.Sprintf(tasksEmployer, username, token, token, token, builder.String())
}

func getGame(tasks []Task, workers []Worker, token string) string {
	var builder strings.Builder

	if len(tasks) == 0 {
		builder.WriteString("<h2>No available tasks</h2>")
	} else if len(workers) == 0 {
		builder.WriteString("<h2>No available workers</h2>")
	} else {
		for i, task := range tasks {
			if i == 0 {
				builder.WriteString(fmt.Sprintf(`<form action="/play" method="post">
				<input type="hidden" id="token" name="token" value="%s"><table>
								<tr>
									<th>Task ID</th>
									<th>Weight</th>
									<th>Choose task</th>
								</tr>`, token))
			}
			builder.WriteString(fmt.Sprintf(`<tr>
											<td>%d</td>
											<td>%d</td>
											<td><input type="radio" id="%d" name="task" value="%d" required></td>
										</tr>`, task.taskID, task.weight, task.taskID, task.taskID))
			if i == len(tasks)-1 {
				builder.WriteString(`</table>`)
			}
		}
		for i, worker := range workers {
			if i == 0 {
				builder.WriteString(`<h3>Available workers</h3><table>
									<tr>
										<th>ID</th>
										<th>Wage</th>
										<th>Fatigue</th>
										<th>Max weight</th>
										<th>Alcoholism</th>
										<th>Select worker</th>
									</tr>`)
			}
			builder.WriteString(fmt.Sprintf(`<tr>
								<td>%d</td>
								<td>%d</td>
								<td>%d</td>
								<td>%d</td>
								<td>%v</td>
								<td><input type="checkbox" name="worker" value="%d"></td>
								</tr>`, worker.userid, worker.wage, worker.fatigue, worker.weight, worker.drinks, worker.userid))
			if i == len(workers)-1 {
				builder.WriteString(`</table><input type="submit" value="Attempt the task"></form>`)
			}
		}
	}
	return builder.String()
}

func (server *APIServer) UpdateWin(task Task, workers []Worker, username string, wageTotal int) {
	server.storage.AddMoney(username, wageTotal/20)
	server.storage.MarkComplete(task, username, workers)
	server.storage.UpdateWorkers(workers)
}

func (server *APIServer) UpdateFail(task Task, workers []Worker, username string, wageTotal, weight int) {
	server.storage.RemoveMoney(username, wageTotal)
	server.storage.UpdateTask(task, weight)
	server.storage.UpdateWorkers(workers)
}
