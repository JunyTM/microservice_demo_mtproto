package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"ms_gmail/model"
	"ms_gmail/pb"
	"ms_gmail/utils"
	"net/http"
	"strconv"
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

	chanelDatas := make(map[int]chan model.RegistPayload)
	for i := 1; i <= workerCountInt; i++ {
		chanelDatas[i] = make(chan model.RegistPayload, workLoadInt)
		go func(chanelData chan<- model.RegistPayload) {
			for j := 1; j <= workLoadInt; j++ {
				temp := model.RegistPayload{
					Name:     utils.RandomString(8),
					Email:    utils.RandomString(20),
					Password: "123456",
				}
				chanelData <- temp
			}
		}(chanelDatas[i])
	}

	errExport := WriteToExcel(&chanelDatas, workerCountInt, workLoadInt)
	if errExport != nil {
		InternalServerError(w, r, errExport)
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

func WriteToExcel(datas *map[int]chan model.RegistPayload, workerCount, workLoad int) error {
	// Write to the excel
	f := excelize.NewFile()

	// Declear the columns names
	columName := []interface{}{
		"ID",
		"Name",
		"Email",
		"Password",
	}

	var indexRow int = 1
	cell, err := excelize.CoordinatesToCellName(1, indexRow)
	if err != nil {
		return errors.New("error get CoordinatesToCellName for column name")
	}
	indexRow += 1
	f.SetSheetRow("Sheet1", cell, &columName)

	for goroutine, chanelData := range *datas {
		log.Printf("==> Execl routine %d running", goroutine)
		go func(chanel <-chan model.RegistPayload) {
			// Write user data to excel
			for temp := range chanel {
				if indexRow > (workerCount * workLoad) {
					return
				}
				slices := []interface{}{
					indexRow,
					temp.Name,
					temp.Email,
					temp.Password,
				}
				cell, err := excelize.CoordinatesToCellName(1, indexRow)
				if err != nil {
					log.Println("Error get CoordinatesToCellName - user:", temp.Name)
				}
				indexRow += 1
				f.SetSheetRow("Sheet1", cell, &slices)
			}
		}(chanelData)
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs("excel/DataBenchmark.xlsx"); err != nil {
		return fmt.Errorf("failed to save excel file: %v", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close excel file: %v", err)
	}
	return nil
}
