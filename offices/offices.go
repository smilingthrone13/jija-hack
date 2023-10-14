package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
	"strings"
)

type OpenHours struct {
	Days  string `json:"days"`
	Hours string `json:"hours"`
}

type LoadItem struct {
	Day     int      `json:"day"`
	Loads   [][2]int `json:"loads"`
	WorkHrs []int    `json:"workHrs"`
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
	LoadIndividuals     []LoadItem  `json:"loadIndividuals"`
	Load                []LoadItem  `json:"load"`
}

func getDayName(index int) string {
	days := []string{"пн", "вт", "ср", "чт", "пт", "сб", "вс"}
	if index >= 1 && index <= 7 {
		return days[index-1]
	}
	return ""
}

func getDayIndex(day string) int {
	// Функция для получения индекса дня недели
	daysOfWeek := map[string]int{
		"пн": 1,
		"вт": 2,
		"ср": 3,
		"чт": 4,
		"пт": 5,
		"сб": 6,
		"вс": 7,
	}
	return daysOfWeek[day]
}

func getWorkDays(days string) []string {
	workDays := []string{}
	daysMap := map[string]bool{
		"пн": true, "вт": true, "ср": true, "чт": true, "пт": true, "сб": true, "вс": true,
	}

	daysList := strings.Split(days, ", ")
	for _, day := range daysList {
		if strings.Contains(day, "-") {
			dayRange := strings.Split(day, "-")
			if len(dayRange) == 2 {
				start := dayRange[0]
				end := dayRange[1]

				for day := start; day != end; {
					workDays = append(workDays, day)
					dayIndex := getDayIndex(day)
					nextDayIndex := (dayIndex % 7) + 1
					day = getDayName(nextDayIndex)
				}
				workDays = append(workDays, end)
			}
		} else {
			if daysMap[day] {
				workDays = append(workDays, day)
			}
		}
	}
	return workDays
}

func generateLoadFromOpenHoursIndividual(openHours []OpenHours) []LoadItem {
	var load []LoadItem

	for _, openHour := range openHours {
		workDays := getWorkDays(openHour.Days)
		hours := strings.Split(openHour.Hours, "-")

		// Если "выходной" в поле "hours", пропускаем этот день
		if len(hours) == 2 {
			startHour := 0
			endHour := 0
			fmt.Sscanf(hours[0], "%d:%d", &startHour)
			fmt.Sscanf(hours[1], "%d:%d", &endHour)

			// Добавляем информацию о загруженности для каждого дня
			for _, day := range workDays {
				loadDay := LoadItem{
					Day:     getDayIndex(day),
					Loads:   [][2]int{},
					WorkHrs: []int{},
				}

				// Добавляем рабочие часы и загруженность для каждого часа
				for hour := startHour; hour <= endHour; hour++ {
					loadDay.Loads = append(loadDay.Loads, [2]int{hour, (rand.Intn(10) + 1) * (rand.Intn(5) + 3)})
					loadDay.WorkHrs = append(loadDay.WorkHrs, hour)
				}
				load = append(load, loadDay)
			}
		}
	}

	// Сортировка "loads" по часам
	for i := range load {
		sort.Slice(load[i].Loads, func(j, k int) bool {
			return load[i].Loads[j][0] < load[i].Loads[k][0]
		})
	}

	return load
}

func generateLoadFromOpenHours(openHours []OpenHours) []LoadItem {
	var load []LoadItem

	for _, openHour := range openHours {
		workDays := getWorkDays(openHour.Days)
		hours := strings.Split(openHour.Hours, "-")

		// Если "выходной" в поле "hours", пропускаем этот день
		if len(hours) == 2 {
			startHour := 0
			endHour := 0
			fmt.Sscanf(hours[0], "%d:%d", &startHour)
			fmt.Sscanf(hours[1], "%d:%d", &endHour)

			// Добавляем информацию о загруженности для каждого дня
			for _, day := range workDays {
				loadDay := LoadItem{
					Day:     getDayIndex(day),
					Loads:   [][2]int{},
					WorkHrs: []int{},
				}

				// Добавляем рабочие часы и загруженность для каждого часа
				for hour := startHour; hour <= endHour; hour++ {
					loadDay.Loads = append(loadDay.Loads, [2]int{hour, (rand.Intn(10) + 1) * (rand.Intn(10) + 5)})
					loadDay.WorkHrs = append(loadDay.WorkHrs, hour)
				}
				load = append(load, loadDay)
			}
		}
	}

	// Сортировка "loads" по часам
	for i := range load {
		sort.Slice(load[i].Loads, func(j, k int) bool {
			return load[i].Loads[j][0] < load[i].Loads[k][0]
		})
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
		data[i].LoadIndividuals = generateLoadFromOpenHoursIndividual(data[i].OpenHoursIndividual)
		if data[i].LoadIndividuals == nil {
			data[i].LoadIndividuals = make([]LoadItem, 0)
		}

		data[i].Load = generateLoadFromOpenHours(data[i].OpenHours)
		if data[i].Load == nil {
			data[i].Load = make([]LoadItem, 0)
		}
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
