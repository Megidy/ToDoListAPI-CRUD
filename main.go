package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/gorilla/mux"
)

type Task struct {
	Id          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Deadline    time.Duration `json:"deadline"`
}

var tasks []Task

func ErrorChecker(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
func LoadData() {
	filepath := "static/data.json"
	File, err := ioutil.ReadFile(filepath)
	ErrorChecker(err)
	err = json.Unmarshal(File, &tasks)
	ErrorChecker(err)

}
func OpenData() []byte {
	filepath := "static/data.json"
	File, err := os.Open(filepath)
	ErrorChecker(err)
	ByteValue, err := ioutil.ReadAll(File)
	ErrorChecker(err)
	return ByteValue
}
func SaveData(data []byte) {
	filepath := "static/data.json"
	err := ioutil.WriteFile(filepath, data, 0644)
	ErrorChecker(err)
}
func HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(tasks)
	ErrorChecker(err)
}

func HandleCreateTask(w http.ResponseWriter, r *http.Request) {

	err := json.Unmarshal(OpenData(), &tasks)
	ErrorChecker(err)

	w.Header().Set("Content-Type", "application/json")
	var NewTask Task

	err = json.NewDecoder(r.Body).Decode(&NewTask)
	ErrorChecker(err)
	tasks = append(tasks, NewTask)
	err = json.NewEncoder(w).Encode(NewTask)
	ErrorChecker(err)
	UpdatedData, err := json.Marshal(tasks)
	ErrorChecker(err)
	SaveData(UpdatedData)

}
func HandleDeleteTasksById(w http.ResponseWriter, r *http.Request) {

	err := json.Unmarshal(OpenData(), &tasks)
	ErrorChecker(err)

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, task := range tasks {
		if task.Id == params["id"] {
			tasks = append(tasks[:index], tasks[index+1:]...)

			break

		}
		err := json.NewEncoder(w).Encode(tasks)
		ErrorChecker(err)
	}

	UpdatedData, err := json.Marshal(tasks)
	ErrorChecker(err)
	SaveData(UpdatedData)
}
func HandleGetTaskById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for _, task := range tasks {
		if task.Id == params["id"] {
			err := json.NewEncoder(w).Encode(task)
			ErrorChecker(err)
			break
		}
	}
}

func HandleAlterTaskById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.Unmarshal(OpenData(), &tasks)
	ErrorChecker(err)
	var AlteredTask Task
	params := mux.Vars(r)
	err = json.NewDecoder(r.Body).Decode(&AlteredTask)
	ErrorChecker(err)
	for index, task := range tasks {
		if task.Id == params["id"] {
			tasks[index] = AlteredTask
			break
		}
	}
	err = json.NewEncoder(w).Encode(&AlteredTask)
	ErrorChecker(err)
	UpdatedData, err := json.Marshal(tasks)
	ErrorChecker(err)
	SaveData(UpdatedData)

}
func Handlers(r *mux.Router) {
	r.HandleFunc("/Tasks", HandleGetTasks).Methods("GET")
	r.HandleFunc("/Tasks", HandleCreateTask).Methods("POST")
	r.HandleFunc("/Tasks/{id}", HandleDeleteTasksById).Methods("DELETE")
	r.HandleFunc("/Tasks/{id}", HandleGetTaskById).Methods("GET")
	r.HandleFunc("/Tasks/{id}/Edit", HandleAlterTaskById).Methods("PUT")
}

func main() {
	LoadData()

	r := mux.NewRouter()
	Handlers(r)

	log.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", r)
}
