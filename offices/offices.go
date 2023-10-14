package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
)

type OpenHours struct {
	Days  string `json:"days"`
	Hours string `json:"hours"`
}

type LoadItem struct {
	Days  string `json:"days"`
	Loads []int  `json:"loads"`
}

type DataItem struct {
	SalePointName       string      `json:"salePointName"`
	Address             string      `json:"address"`
	Status              string      `json:"status"`
	OpenHours           []OpenHours `json:"openHours"`
	RKO                 string      `json:"rko"`
	OpenHoursIndividual []OpenHours `json:"openHoursIndividual"`
	OfficeType          string      `json:"officeType"`
	SalePointFormat     string      `json:"salePointFormat"`
	SuoAvailability     string      `json:"suoAvailability"`
	HasRamp             string      `json:"hasRamp"`
	Latitude            float64     `json:"latitude"`
	Longitude           float64     `json:"longitude"`
	MetroStation        string      `json:"metroStation"`
	Distance            int         `json:"distance"`
	Kep                 bool        `json:"kep"`
	MyBranch            bool        `json:"myBranch"`
	Load                []LoadItem  `json:"load"`
}

func generateLoad(openHours []OpenHours) []LoadItem {
	var load []LoadItem

	for _, openHour := range openHours {
		days := strings.Split(openHour.Days, ", ")
		hours := strings.Split(openHour.Hours, "-")

		// Если "выходной" в поле "hours", пропускаем этот день
		if len(hours) == 2 {
			startHour := 0
			endHour := 0
			fmt.Sscanf(hours[0], "%d:%d", &startHour)
			fmt.Sscanf(hours[1], "%d:%d", &endHour)
			dayLoads := []int{}
			for hour := startHour; hour <= endHour; hour++ {
				dayLoads = append(dayLoads, rand.Intn(10)+1)
			}
			for _, day := range days {
				load = append(load, LoadItem{
					Days:  day,
					Loads: dayLoads,
				})
			}
		}
	}

	return load
}

func main() {
	// Чтение JSON из файла
	inputFileName := "offices.txt"
	jsonData, err := ioutil.ReadFile(inputFileName)
	if err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}

	// Распаковка JSON в массив объектов
	var data []DataItem
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Printf("Ошибка при распаковке JSON: %v\n", err)
		return
	}

	// Для каждого объекта в массиве data создаем поле "Load" на основе "openHoursIndividual"
	for i := range data {
		data[i].Load = generateLoad(data[i].OpenHoursIndividual)
	}

	// Затем можно маршализовать обновленные данные и сохранить их в файл, если это необходимо
	updatedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Ошибка при маршализации JSON: %v\n", err)
		return
	}

	outputFileName := "offices_rich.txt"
	err = ioutil.WriteFile(outputFileName, updatedJSON, 0644)
	if err != nil {
		fmt.Printf("Ошибка при записи в файл: %v\n", err)
		return
	}

	fmt.Printf("Данные успешно обновлены и сохранены в файл %s\n", outputFileName)
}
