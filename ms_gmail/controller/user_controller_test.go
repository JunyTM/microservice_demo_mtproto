package controller

import (
	"context"
	"fmt"
	"log"
	"ms_gmail/pb"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var f *excelize.File

func TestMain(t *testing.M) {
	var err error
	f, err = excelize.OpenFile("../excel/DataBenchmark.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	t.Run()
}

func TestRegisterGRPC(t *testing.T) {
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	require.NoError(t, err)
	require.NotEmpty(t, rows)

	var workload int = 500
	var wg sync.WaitGroup
	for index := 0; index < 10; index++ {
		wg.Add(1)
		go func(rows [][]string, yStart, workLoad int, t *testing.T, wg *sync.WaitGroup) {
			for index := yStart * workLoad; index < (yStart+1)*workLoad || index >= len(rows); index++ {
				// Bypass the column name
				if index == 0 {
					continue
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

				// 2 - the column of email, and password default is '123456'
				serverResponse, err := clientAuthRPC.CreateUser(ctx, &pb.CreateUserMessage{Name: rows[index][1], Email: rows[index][2], Password: "123456"})
				if err != nil {
					if err == grpc.ErrClientConnTimeout {
						log.Printf("===> error timeout while login Id = %d", index)
						break
					}
				}
				require.NotNil(t, serverResponse)
			}
			wg.Done()
		}(rows, index, workload, t, &wg)
	}
	wg.Wait()
}

func TestLoginGRPC(t *testing.T) {
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	require.NoError(t, err)
	require.NotEmpty(t, rows)

	var errChan chan error
	var workload int = 200
	var wg sync.WaitGroup
	for index := 0; index < 10; index++ {
		go func(rows [][]string, yStart, workLoad int, t *testing.T, wg *sync.WaitGroup, errChan chan error) {
			for index := yStart * workLoad; index < (yStart+1)*workLoad || index >= len(rows); index++ {
				// Bypass the column name
				if index == 0 {
					continue
				}
				go func(user []string, errChan chan<- error, index int) {
					// Contact to server auth
					conn, err := grpc.Dial("0.0.0.0:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						log.Fatalf("did not connect: %v", err)
					}
					defer conn.Close()
					clientAuthRPC := pb.NewAuthenClient(conn)

					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					// 2 - the column of email, and password default is '123456'
					serverResponse, err := clientAuthRPC.Login(ctx, &pb.LoginMessage{Email: user[2], Password: "123456"})
					if err != nil {
						errChan <- fmt.Errorf("====> User time out - number: %d / err: %v", index, err)
					}
					require.NotNil(t, serverResponse)
					require.NotNil(t, serverResponse.AccessToken)
				}(rows[index], errChan, index)
			}
		}(rows, index, workload, t, &wg, errChan)
	}

	for {
		select {
		case err := <-errChan:
			log.Printf("===> error timeout while login Id = %v", err)
			break
		default:
		}
	}

}
