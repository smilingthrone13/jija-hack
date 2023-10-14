package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

type Service struct {
	ServiceCapability string `json:"serviceCapability"`
	ServiceActivity   string `json:"serviceActivity"`
}

type LoadItem struct {
	Day     int     `json:"day"`
	Loads   [][]int `json:"loads"`
	WorkHrs []int   `json:"workHrs"`
}

type ATM struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	AllDay    bool    `json:"allDay"`
	Services  struct {
		Wheelchair        Service `json:"wheelchair"`
		Blind             Service `json:"blind"`
		NFCForBankCards   Service `json:"nfcForBankCards"`
		QRRead            Service `json:"qrRead"`
		SupportsUSD       Service `json:"supportsUsd"`
		SupportsChargeRUB Service `json:"supportsChargeRub"`
		SupportsEUR       Service `json:"supportsEur"`
		SupportsRUB       Service `json:"supportsRub"`
	} `json:"services"`
	Load []LoadItem `json:"load"`
}

type ATMData struct {
	ATMs []ATM `json:"atms"`
}

func generateRandomLoad(allDay bool) []LoadItem {
	rand.Seed(time.Now().UnixNano())
	var (
		load                []LoadItem
		workHrs             []int
		openHour, closeHour int
	)

	openHour = 0
	closeHour = 23
	if !allDay {
		openHour = 8
		closeHour = 22
	}

	for day := 1; day <= 7; day++ {
		loadDay := make([][]int, 0)
		for hour := openHour; hour <= closeHour; hour++ {
			loadHour := []int{hour, (rand.Intn(10) + 1) * (rand.Intn(3) + 1)}
			loadDay = append(loadDay, loadHour)
			workHrs = append(workHrs, hour)

		}
		load = append(load, LoadItem{
			Day:     day,
			Loads:   loadDay,
			WorkHrs: workHrs,
		})
		workHrs = make([]int, 0)

	}

	return load
}

func main() {
	// Чтение JSON из файла
	inputFileName := "atms.txt"
	jsonData, err := ioutil.ReadFile(inputFileName)
	if err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}

	// Распаковка JSON
	var data ATMData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Printf("Ошибка при распаковке JSON: %v\n", err)
		return
	}

	// Заполнение поля "load" случайными значениями
	for i := range data.ATMs {
		allDay := data.ATMs[i].AllDay
		data.ATMs[i].Load = generateRandomLoad(allDay)
	}

	// Запись обновленных данных обратно в файл
	updatedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Ошибка при маршализации JSON: %v\n", err)
		return
	}

	outputFileName := "atms_rich.txt"
	err = ioutil.WriteFile(outputFileName, updatedJSON, 0644)
	if err != nil {
		fmt.Printf("Ошибка при записи в файл: %v\n", err)
		return
	}

	fmt.Printf("Данные успешно обновлены и сохранены в файл %s\n", outputFileName)
}
