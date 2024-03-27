package controller

import (
	"context"
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
			for index := yStart * workLoad; index < (yStart+1)*workLoad; index++ {
				// Bypass the column name
				if index == 0 {
					continue
				} else if index > len(rows) {
					break
				}

				// Contact to server auth
				conn, err := grpc.Dial(serverHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func TestLoging(t *testing.T) {
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	require.NoError(t, err)
	require.NotEmpty(t, rows)

	var wg sync.WaitGroup
	for index := 1; index <= 2; index++ {
		wg.Add(1)
		go func(user []string, wg *sync.WaitGroup, index int) {
			// Contact to server auth
			conn, err := grpc.Dial(serverHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			clientAuthRPC := pb.NewAuthenClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second) // Wait for 10sec timeout
			defer cancel()
			// 2 - the column of email, and password default is '123456'
			serverResponse, err := clientAuthRPC.Login(ctx, &pb.LoginMessage{Email: user[2], Password: "123456"})
			if err != nil {
				if err == grpc.ErrClientConnTimeout {
					log.Printf("====> User time out - number: %d / err: %v\n", index, err)
				} else {
					log.Printf("====> RPC error - number: %d / err: %v\n", index, err)
				}
			}
			require.NotNil(t, serverResponse)
			require.NotNil(t, serverResponse.AccessToken)
			log.Println("User - ", serverResponse.SessionId)
			wg.Done()
		}(rows[index], &wg, index)
	}
	wg.Wait()
}
