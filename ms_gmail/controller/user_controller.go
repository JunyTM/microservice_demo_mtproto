package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"ms_gmail/model"
	"ms_gmail/pb"
	"ms_gmail/utils"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/render"
	"github.com/xuri/excelize/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserController interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	CreateDataUsers(w http.ResponseWriter, r *http.Request)
}

type userController struct{}

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

func (c *userController) CreateDataUsers(w http.ResponseWriter, r *http.Request) {
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

	var wg sync.WaitGroup
	chanelDatas := make(chan []model.RegistPayload, workerCountInt)
	for i := 1; i <= workerCountInt; i++ {
		wg.Add(1)
		go func(chanelData chan<- []model.RegistPayload) {
			var data []model.RegistPayload
			for j := 1; j <= workLoadInt; j++ {
				temp := model.RegistPayload{
					Name:     utils.RandomString(8),
					Email:    utils.RandomString(20),
					Password: "123456",
				}
				data = append(data, temp)
			}
			chanelData <- data
			wg.Done()
		}(chanelDatas)
	}
	wg.Wait()

	// Write to the excel
	f := excelize.NewFile()

	var index int = 1
	columName := []interface{}{
		"ID",
		"Name",
		"Email",
		"Password",
	}
	cell, err := excelize.CoordinatesToCellName(1, index)
	if err != nil {
		InternalServerError(w, r, err)
		return
	}
	f.SetSheetRow("Sheet1", cell, &columName)
	index += 1

	log.Println("Make data successfull")
	for {
		select {
		case chanelData := <-chanelDatas:
			for _, data := range chanelData {
				cell, err := excelize.CoordinatesToCellName(1, index)
				if err != nil {
					InternalServerError(w, r, err)
					return
				}

				temp := []interface{}{
					index,
					data.Name,
					data.Email,
					data.Password,
				}
				f.SetSheetRow("Sheet1", cell, &temp)
				index += 1
			}
		default:
			goto SAVE
		}
	}

SAVE:
	// Save spreadsheet by the given path.
	if err := f.SaveAs("excel/DataBenchmark.xlsx"); err != nil {
		InternalServerError(w, r, err)
		return
	}
	if err := f.Close(); err != nil {
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
	return &userController{}
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
