package service

import (
	"ms_gmail/model"
	"ms_gmail/utils"
	"sync"

	"github.com/xuri/excelize/v2"
)

type ExcelWorkerInterface interface {
	Start(concurrency, workLoad int, sheetName, path string) error
}

type UserPool struct {
	concurrency int
	workLoad    int
	tasksChan   chan []model.RegistPayload
	wg          sync.WaitGroup
	mu          sync.Mutex // guards
}

func NewUserPool() ExcelWorkerInterface {
	return &UserPool{
		tasksChan: make(chan []model.RegistPayload),
	}
}

func (p *UserPool) Start(concurrency, workLoad int, sheetName, path string) error {
	p.concurrency = concurrency
	p.workLoad = workLoad

	for i := 1; i <= p.concurrency; i++ {
		// p.wg.Add(1)
		go func(userPool *UserPool) {
			var data []model.RegistPayload
			for j := 1; j <= userPool.workLoad; j++ {
				temp := model.RegistPayload{
					Name:     utils.RandomString(8),
					Email:    utils.RandomString(20),
					Password: "123456",
				}
				data = append(data, temp)
			}
			userPool.tasksChan <- data
			// userPool.wg.Done()
		}(p)
	}
	// p.wg.Wait()
	// Write to the excel
	f := excelize.NewFile()

	index := indexCell{value: 1}
	columName := []interface{}{
		"ID",
		"Name",
		"Email",
		"Password",
	}
	cell, err := excelize.CoordinatesToCellName(1, index.GetValue())
	if err != nil {
		return err
	}
	f.SetSheetRow(sheetName, cell, &columName)
	index.Increase()
	for i := 1; i <= p.concurrency; i++ {
		p.wg.Add(1)
		dataUsers := <-p.tasksChan
		go WriteDataCell(dataUsers, &index, p, f, sheetName)
	}

	p.wg.Wait()
	// Save spreadsheet by the given path.
	if err := f.SaveAs(path); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func WriteDataCell(datas []model.RegistPayload, index *indexCell, pool *UserPool, f *excelize.File, sheetName string) {
	for _, data := range datas {
		pool.mu.Lock()
		idx := index.GetValue()
		index.Increase()
		pool.mu.Unlock()
		cell, err := excelize.CoordinatesToCellName(1, idx)
		if err != nil {
			return
		}

		temp := []interface{}{
			idx - 1,
			data.Name,
			data.Email,
			data.Password,
		}
		f.SetSheetRow(sheetName, cell, &temp)
	}
	pool.wg.Done()
}

type indexCell struct {
	value int
	mu    sync.Mutex // guards
}

func (i *indexCell) GetValue() int {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.value
}

func (i *indexCell) Increase() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.value += 1
}
