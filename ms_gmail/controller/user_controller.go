package controller

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"ms_gmail/infrastructure"
// 	"ms_gmail/model"
// 	"ms_gmail/pb"
// 	"ms_gmail/service"
// 	"net/http"
// 	"strconv"
// 	"sync"
// 	"time"

// 	"github.com/go-chi/render"
// 	"github.com/xuri/excelize/v2"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// type UserController interface {
// 	Login(w http.ResponseWriter, r *http.Request)
// 	Register(w http.ResponseWriter, r *http.Request)
// 	GetUserProfile(w http.ResponseWriter, r *http.Request)
// 	GenerateUsers(w http.ResponseWriter, r *http.Request)
// 	LoadUsersGenerated(w http.ResponseWriter, r *http.Request)
// }

// // var serverHost string =  "ms_auth:9090"
// var serverHost string = infrastructure.GetServerHost()

// type userController struct {
// 	userWorker service.ExcelWorkerInterface
// }

// func (c *userController) Login(w http.ResponseWriter, r *http.Request) {
// 	var loginPayload model.LoginPayload
// 	err := json.NewDecoder(r.Body).Decode(&loginPayload)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	// Contact to server auth
// 	conn, err := grpc.Dial(serverHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()
// 	clientAuthRPC := pb.NewAuthenClient(conn)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	serverResponse, err := clientAuthRPC.Login(ctx, &pb.LoginMessage{Email: loginPayload.Email, Password: loginPayload.Password})
// 	if err != nil {
// 		log.Fatalf("> gRPC error: %v", err)
// 	}

// 	httpResponse := &model.Response{
// 		Data:    serverResponse,
// 		Success: true,
// 		Message: "Login successful",
// 	}
// 	render.JSON(w, r, httpResponse)
// }

// func (c *userController) Register(w http.ResponseWriter, r *http.Request) {
// 	var registPayload model.RegistPayload
// 	err := json.NewDecoder(r.Body).Decode(&registPayload)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	// Contact to server auth
// 	conn, err := grpc.Dial(serverHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()
// 	clientAuthRPC := pb.NewAuthenClient(conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	serverResponse, err := clientAuthRPC.CreateUser(ctx, &pb.CreateUserMessage{Name: registPayload.Name, Email: registPayload.Email, Password: registPayload.Password})
// 	if err != nil {
// 		log.Fatalf("> gRPC error: %v", err)
// 	}

// 	httpResponse := &model.Response{
// 		Data:    serverResponse,
// 		Success: true,
// 		Message: "Create successful",
// 	}
// 	render.JSON(w, r, httpResponse)
// }

// func (c *userController) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	
// }

// func (c *userController) GenerateUsers(w http.ResponseWriter, r *http.Request) {
// 	/*
// 		workerCount - the number of workers in worker pool.
// 		workLoad - the jobs in single worker.
// 	*/

// 	workerCount := r.URL.Query().Get("workers")
// 	workLoad := r.URL.Query().Get("workload")
// 	if workerCount == "" || workLoad == "" {
// 		BadRequest(w, r, errors.New("param err: workers or workload must be specified"))
// 		return
// 	}

// 	workerCountInt, err := strconv.Atoi(workerCount)
// 	if err != nil {
// 		InternalServerError(w, r, err)
// 		return
// 	}

// 	workLoadInt, err := strconv.Atoi(workLoad)
// 	if err != nil {
// 		InternalServerError(w, r, err)
// 		return
// 	}

// 	if workerCountInt <= 0 || workerCountInt > 10 {
// 		BadRequest(w, r, errors.New("worker count must be between 1 and 10"))
// 		return
// 	}

// 	err = c.userWorker.Start(workerCountInt, workLoadInt, "Sheet1", "excel/DataBenchmark.xlsx")
// 	if err != nil {
// 		InternalServerError(w, r, err)
// 		return
// 	}

// 	res := model.Response{
// 		Data:    nil,
// 		Success: true,
// 		Message: "Make data successful",
// 	}
// 	render.JSON(w, r, res)
// }

// func (c *userController) LoadUsersGenerated(w http.ResponseWriter, r *http.Request) {
// 	workLoad := r.URL.Query().Get("workload")
// 	if workLoad == "" {
// 		BadRequest(w, r, errors.New("param err: workers or workload must be specified"))
// 		return
// 	}

// 	workLoadInt, err := strconv.Atoi(workLoad)
// 	if err != nil {
// 		InternalServerError(w, r, err)
// 		return
// 	}

// 	f, err := excelize.OpenFile("excel/DataBenchmark.xlsx")
// 	if err != nil {
// 		InternalServerError(w, r, err)
// 		return
// 	}

// 	// Get all the rows in the Sheet1.
// 	rows, err := f.GetRows("Sheet1")
// 	if err != nil {
// 		InternalServerError(w, r, err)
// 		return
// 	}

// 	var wg sync.WaitGroup
// 	errChan := make(chan string)
// 	for index := 0; index < 10; index++ {
// 		wg.Add(1)
// 		go func(rows [][]string, yStart, workLoad int, wg *sync.WaitGroup, errChan chan string) {
// 			for index := yStart * workLoad; index < (yStart+1)*workLoad; index++ {
// 				// Bypass the column name
// 				if index == 0 {
// 					continue
// 				} else if index > len(rows) {
// 					break
// 				}

// 				// Contact to server auth
// 				conn, err := grpc.Dial(serverHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 				if err != nil {
// 					log.Fatalf("did not connect: %v", err)
// 				}
// 				defer conn.Close()
// 				clientAuthRPC := pb.NewAuthenClient(conn)

// 				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 				defer cancel()

// 				// 2 - the column of email, and password default is '123456'
// 				serverResponse, err := clientAuthRPC.CreateUser(ctx, &pb.CreateUserMessage{Name: rows[index][1], Email: rows[index][2], Password: "123456"})
// 				if err != nil {
// 					if err == grpc.ErrClientConnTimeout {
// 						// log.Printf("===> error timeout while login Id = %d", index)
// 						errChan <- fmt.Sprintf("=> error timeout while login Id = %d", index)
// 						break
// 					}
// 				}
// 				if serverResponse == nil {
// 					errChan <- fmt.Sprintf("=> error timeout while login Id = %d", index)
// 				}
// 			}
// 			wg.Done()
// 		}(rows, index, workLoadInt, &wg, errChan)
// 	}

// 	wg.Wait()
// 	close(errChan)

// 	errors := []string{}
// 	for err := range errChan {
// 		errors = append(errors, err)
// 	}

// 	res := model.Response{
// 		Data:    errors,
// 		Success: true,
// 		Message: "Request finished",
// 	}
// 	render.JSON(w, r, res)
// }

// func NewUserController() UserController {
// 	return &userController{
// 		userWorker: service.NewUserPool(),
// 	}
// }

// func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
// 	http.Error(w, "Bad Request", http.StatusBadRequest)
// 	w.WriteHeader(http.StatusBadRequest)
// 	res := &model.Response{
// 		Data:    nil,
// 		Success: false,
// 		Message: err.Error(),
// 	}
// 	render.JSON(w, r, res)
// }

// func NotFound(w http.ResponseWriter, r *http.Request, err error) {
// 	http.Error(w, "Not Found", http.StatusNotFound)
// 	w.WriteHeader(http.StatusNotFound)
// 	res := &model.Response{
// 		Data:    nil,
// 		Success: false,
// 		Message: err.Error(),
// 	}
// 	render.JSON(w, r, res)
// }

// func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
// 	// http.Error(w, "Internal Server Error\n", http.StatusInternalServerError)
// 	w.WriteHeader(http.StatusInternalServerError)
// 	res := &model.Response{
// 		Data:    nil,
// 		Success: false,
// 		Message: err.Error(),
// 	}
// 	render.JSON(w, r, res)
// }
