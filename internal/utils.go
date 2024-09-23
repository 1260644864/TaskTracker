package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type task struct {
	Serial  int    `json:"serial"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

var file_path = "../data/tasks.json"
var quantity = 0

//	func handleError(err error, message string) {
//		if err != nil {
//			log.Fatalf("Error:%v\n", err)
//		} else {
//			log.Fatalf("Error:%v\n", message)
//		}
//	}
func createNewJson() error {
	file, err := os.Create(file_path)
	if err != nil {
		return fmt.Errorf("an error occurred while creating the file:%v", err)
	}

	defer file.Close()

	_, err = file.WriteString("[]")
	if err != nil {
		return fmt.Errorf("an error occurred while writing the file:%v", err)
	}

	return nil
}

func OpenJson() ([]task, error) {
	if _, err := os.Stat(file_path); os.IsNotExist(err) {
		err = createNewJson()

		if err != nil {
			return nil, err
		}
	}

	file, err := os.Open(file_path)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while opening the file:%v", err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read JSON Failed:%v", err)
	}

	var items []task

	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, fmt.Errorf("analyze JSON Failed:%v", err)
	}

	return items, nil
}

func WriteJson(items []task) error {
	file, err := os.OpenFile(file_path, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("an error occurred while opening the file:%v", err)
	}

	defer file.Close()

	newData, err := json.MarshalIndent(items, "", "    ")
	if err != nil {
		return fmt.Errorf("encode JSON Failed: %v", err)
	}

	_, err = file.Write(newData)
	if err != nil {
		return fmt.Errorf("write JSON Failed: %v", err)
	}

	return nil
}

func UpdateQuant() error {
	if _, err := os.Stat(file_path); os.IsNotExist(err) {
		err = createNewJson()

		if err != nil {
			return err
		}
	}

	file, err := os.Open(file_path)
	if err != nil {
		return fmt.Errorf("an error occurs while updating the quantity: %v", err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("an error occurs while updating the quantity: %v", err)
	}

	var items []task

	err = json.Unmarshal(data, &items)
	if err != nil {
		return fmt.Errorf("an error occurs while updating the quantity: %v", err)
	}

	quantity = len(items)

	return nil

}

func AddItem(newContent string) error {
	newItem := task{
		Serial:  quantity + 1,
		Content: newContent,
		Status:  "todo",
	}
	items, err := OpenJson()
	if err != nil {
		return err
	}

	items = append(items, newItem)

	err = WriteJson(items)
	if err != nil {
		return err
	}

	return nil

}

func DeleteItem(serial int) error {
	if !(serial >= 1 && serial <= quantity) {
		return fmt.Errorf("invalid Serial Number Input")
	}

	items, err := OpenJson()
	if err != nil {
		return err
	}

	items = append(items[:serial-1], items[serial:]...)

	for i := range items {
		items[i].Serial = i + 1
	}

	err = WriteJson(items)
	if err != nil {
		return err
	}

	return nil
}

func UpdateItem(serial int, newContent string) error {
	if !(serial >= 1 && serial <= quantity) {
		return fmt.Errorf("invalid Serial Number Input")
	}

	items, err := OpenJson()
	if err != nil {
		return err
	}

	items[serial-1].Content = newContent

	err = WriteJson(items)
	if err != nil {
		return err
	}

	return nil

}

func MarkDone(serial int) error {
	if !(serial >= 1 && serial <= quantity) {
		log.Fatalf("Invalid Serial Number Input")
	}

	items, err := OpenJson()
	if err != nil {
		return err
	}

	items[serial-1].Status = "done"

	err = WriteJson(items)
	if err != nil {
		return err
	}

	return nil

}

func MarkInProgress(serial int) error {
	if !(serial >= 1 && serial <= quantity) {
		log.Fatalf("Invalid Serial Number Input")
	}

	items, err := OpenJson()
	if err != nil {
		return err
	}

	items[serial-1].Status = "in-progress"

	err = WriteJson(items)
	if err != nil {
		return err
	}

	return nil
}

func List() error {
	items, err := OpenJson()
	if err != nil {
		return err
	}

	if len(items) == 0 {
		fmt.Println("No Task")
	} else {
		for _, item := range items {
			fmt.Printf("Task %v\n", item.Serial)
			fmt.Printf("Content: %v\n", item.Content)
			fmt.Printf("Status: %v\n\n", item.Status)
		}
	}

	return nil
}

func ListDone() error {
	items, err := OpenJson()
	if err != nil {
		return err
	}

	if len(items) == 0 {
		fmt.Println("No Task")
	} else {
		for _, item := range items {
			if item.Status == "done" {
				fmt.Printf("Task %v\n", item.Serial)
				fmt.Printf("Content: %v\n", item.Content)
				fmt.Printf("Status: %v\n\n", item.Status)
			}
		}
	}

	return nil
}

func ListTodo() error {
	items, err := OpenJson()
	if err != nil {
		return err
	}

	if len(items) == 0 {
		fmt.Println("No Task")
	} else {
		for _, item := range items {
			if item.Status == "todo" {
				fmt.Printf("Task %v\n", item.Serial)
				fmt.Printf("Content: %v\n", item.Content)
				fmt.Printf("Status: %v\n\n", item.Status)
			}
		}
	}

	return nil

}

func ListInProgress() error {
	items, err := OpenJson()
	if err != nil {
		return err
	}

	if len(items) == 0 {
		fmt.Println("No Task")
	} else {
		for _, item := range items {
			if item.Status == "in-progress" {
				fmt.Printf("Task %v\n", item.Serial)
				fmt.Printf("Content: %v\n", item.Content)
				fmt.Printf("Status: %v\n\n", item.Status)
			}
		}
	}

	return nil
}
