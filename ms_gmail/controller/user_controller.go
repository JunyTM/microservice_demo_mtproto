package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"ms_gmail/model"
	"ms_gmail/pb"
	"ms_gmail/service"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserController interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	GenerateUsers(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	userWorker service.ExcelWorkerInterface
}

func (c *userController) Login(w http.ResponseWriter, r *http.Request) {
	var loginPayload model.LoginPayload
	err := json.NewDecoder(r.Body).Decode(&loginPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Contact to server auth
	conn, err := grpc.Dial("0.0.0.0:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	clientAuthRPC := pb.NewAuthenClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	serverResponse, err := clientAuthRPC.Login(ctx, &pb.LoginMessage{Email: loginPayload.Email, Password: loginPayload.Password})
	if err != nil {
		log.Fatalf("> gRPC error: %v", err)
	}

	httpResponse := &model.Response{
		Data:    serverResponse,
		Success: true,
		Message: "Login successful",
	}
	render.JSON(w, r, httpResponse)
}

func (c *userController) Register(w http.ResponseWriter, r *http.Request) {
	var registPayload model.RegistPayload
	err := json.NewDecoder(r.Body).Decode(&registPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Contact to server auth
	conn, err := grpc.Dial("0.0.0.0:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	clientAuthRPC := pb.NewAuthenClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	serverResponse, err := clientAuthRPC.CreateUser(ctx, &pb.CreateUserMessage{Name: registPayload.Name, Email: registPayload.Email, Password: registPayload.Password})
	if err != nil {
		log.Fatalf("> gRPC error: %v", err)
	}

	httpResponse := &model.Response{
		Data:    serverResponse,
		Success: true,
		Message: "Create successful",
	}
	render.JSON(w, r, httpResponse)
}
func (c *userController) GenerateUsers(w http.ResponseWriter, r *http.Request) {
	/*
		workerCount - the number of workers in worker pool.
		workLoad - the jobs in single worker.
	*/

	workerCount := r.URL.Query().Get("workers")
	workLoad := r.URL.Query().Get("workload")
	if workerCount == "" || workLoad == "" {
		BadRequest(w, r, errors.New("param err: workers or workload must be specified"))
		return
	}

	workerCountInt, err := strconv.Atoi(workerCount)
	if err != nil {
		InternalServerError(w, r, err)
		return
	}

	workLoadInt, err := strconv.Atoi(workLoad)
	if err != nil {
		InternalServerError(w, r, err)
		return
	}

	if workerCountInt <= 0 || workerCountInt > 10 {
		BadRequest(w, r, errors.New("worker count must be between 1 and 10"))
		return
	}

	err = c.userWorker.Start(workerCountInt, workLoadInt, "Sheet1", "excel/DataBenchmark.xlsx")
	if err != nil {
		InternalServerError(w, r, err)
		return
	}

	res := model.Response{
		Data:    nil,
		Success: true,
		Message: "Make data successful",
	}
	render.JSON(w, r, res)
}

func NewUserController() UserController {
	return &userController{
		userWorker: service.NewUserPool(),
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
	w.WriteHeader(http.StatusBadRequest)
	res := &model.Response{
		Data:    nil,
		Success: false,
		Message: err.Error(),
	}
	render.JSON(w, r, res)
}

func NotFound(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Not Found", http.StatusNotFound)
	w.WriteHeader(http.StatusNotFound)
	res := &model.Response{
		Data:    nil,
		Success: false,
		Message: err.Error(),
	}
	render.JSON(w, r, res)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	w.WriteHeader(http.StatusInternalServerError)
	res := &model.Response{
		Data:    nil,
		Success: false,
		Message: err.Error(),
	}
	render.JSON(w, r, res)
}
