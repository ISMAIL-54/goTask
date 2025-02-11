package main

import (
    "fmt"
    "log"
	"encoding/json"
	"net/http"
    "strconv"
    "github.com/gorilla/mux"
)

type APIServer struct {
    listenAddr  string
    storage     Storage
}

func NewAPIServer(addr string, storage Storage) *APIServer {
    return &APIServer{ 
        listenAddr: addr,
        storage: storage,
    }
}

func (server *APIServer) Run() {
    router := mux.NewRouter()

	log.Println("API Server is running on port:", server.listenAddr)

    router.HandleFunc("/api/task", makeHTTPHandlerFunc(server.HandleGetTasks)).Methods("GET")
    router.HandleFunc("/api/task/{id}", makeHTTPHandlerFunc(server.HandleGetTaskByID)).Methods("GET")
    router.HandleFunc("/api/task", makeHTTPHandlerFunc(server.HandleCreateTask)).Methods("POST")
    router.HandleFunc("/api/task/{id}", makeHTTPHandlerFunc(server.HandleUpdateTask)).Methods("PUT")
    router.HandleFunc("/api/task/{id}", makeHTTPHandlerFunc(server.HandleDeleteTask)).Methods("DELETE")

    /*
        GET     /api/task
        GET     /api/task/{id}
        POST    /api/task
        PUT     /api/task/{id}
        DELETE  /api/task/{id}
        
    */
    http.ListenAndServe(server.listenAddr, router)
}

//  func (server *APIServer) HandleTask(w http.ResponseWriter, r *http.Request) error {
//     switch r.Method {
//     case http.MethodGet:
//         return server.HandleGetTasks(w, r)
//     case http.MethodPost:
//         return server.HandleCreateTask(w, r)
//     case http.MethodPut:
//         return server.HandleUpdateTask(w, r)
//     case http.MethodDelete:
//         return server.HandleDeleteTask(w, r)
//     }
//     return fmt.Errorf("Method not allowed: %w", r.Method) 
// }

func (server *APIServer) HandleGetTasks(w http.ResponseWriter, r *http.Request) error {
    tasks, err := server.storage.GetTasks()
    if err != nil {
        return err
    }

    writeJSON(w, http.StatusOK, tasks)
    return nil
}

func (server *APIServer) HandleGetTaskByID(w http.ResponseWriter, r *http.Request) error {
    id, err := getID(r)
    if err != nil {
        return err
    }

    task, err := server.storage.GetTaskByID(id)
    if err != nil {
        return err
    }

    writeJSON(w, http.StatusOK, task)
    return nil
}

func (server *APIServer) HandleCreateTask(w http.ResponseWriter, r *http.Request) error {
    newTask := new(Task)
    if err := json.NewDecoder(r.Body).Decode(newTask); err != nil {
        return err
    }

    if err := server.storage.CreateTask(newTask); err != nil {
        return err
    }

    return nil
}

func (server *APIServer) HandleUpdateTask(w http.ResponseWriter, r *http.Request) error {
    id, err := getID(r)
    if err != nil {
        return err
    }

    modifiedTask := new(Task)
    if err := json.NewDecoder(r.Body).Decode(modifiedTask); err != nil {
        return err
    }
    
    if err := server.storage.UpdateTask(id, modifiedTask); err != nil {
        return err
    }

    return nil
}

func (server *APIServer) HandleDeleteTask(w http.ResponseWriter, r *http.Request) error {
    id, err := getID(r)
    if err != nil {
        return err
    }

    if err := server.storage.DeleteTask(id); err != nil {
        return err
    }

    return writeJSON(w, http.StatusOK, map[string]int{"deleted": id})
}


type apiFunc func(w http.ResponseWriter, r *http.Request) error

type apiError struct {
    Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
    w.WriteHeader(status)
    w.Header().Set("Content-Type", "application/json")
    return json.NewEncoder(w).Encode(data)
}

func makeHTTPHandlerFunc(fun apiFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := fun(w, r); err != nil {
            writeJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
        }
    }
}

func getID(r *http.Request) (int, error) {
    id_str := mux.Vars(r)["id"]
    id, err := strconv.Atoi(id_str)
    if err != nil {
        return id, fmt.Errorf("invalid given id: %s", id_str) 
    }
    return id, nil
}
