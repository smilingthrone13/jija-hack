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
	Days  string `json:"days"`
	Loads []int  `json:"loads"`
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

func generateRandomLoad() []LoadItem {
	rand.Seed(time.Now().UnixNano())
	days := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	var load []LoadItem

	for _, day := range days {
		loadDay := make([]int, 24)
		for hour := 9; hour <= 20; hour++ {
			loadDay[hour] = rand.Intn(10) + 1
		}
		load = append(load, LoadItem{
			Days:  day,
			Loads: loadDay[9:21], // Ограничиваем часы от 9 до 20 включительно
		})
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
		data.ATMs[i].Load = generateRandomLoad()
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
