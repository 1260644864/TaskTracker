package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type task struct {
	Serial  int    `json:"serial"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

var file_path = "tasks.json"
var quantity = 0

var rootCmd = &cobra.Command{
	Use:   "task-cli",
	Short: "A simple CLI to manage tasks",
}

var addCmd = &cobra.Command{
	Use:   "add <taskContent>",
	Short: "To add a task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := addItem(args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("add task successfully!")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <taskSerial> <taskContent>",
	Short: "To update a task",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("the first arg must be Int:%v", err)
		}

		fmt.Printf("%v\n", args[1])

		err = updateItem(serial, args[1])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("update tasks successfully")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <taskSerial>",
	Short: "To delete a task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("the first arg must be Int:%v", err)
		}

		err = deleteItem(serial)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("delete task successfully")
	},
}

var markIPCmd = &cobra.Command{
	Use:   "mark-in-progress <taskSerial>",
	Short: "To mark a task with in-progress status",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("the first arg must be Int:%v", err)
		}

		err = markInProgress(serial)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("mark task successfully")
	},
}

var markDCmd = &cobra.Command{
	Use:   "mark-done <taskSerial>",
	Short: "To mark a task with done status",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("the first arg must be Int:%v", err)
		}

		err = markDone(serial)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("mark task successfully")
	},
}

var listCmd = &cobra.Command{
	Use:   "list [optional]",
	Short: "To list tasks",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := list()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			switch args[0] {
			case "done":
				err := listDone()
				if err != nil {
					log.Fatal(err)
				}
			case "todo":
				err := listTodo()
				if err != nil {
					log.Fatal(err)
				}
			case "in-progress":
				err := listInProgress()
				if err != nil {
					log.Fatal(err)
				}
			default:
				// handleError(nil, "Invalid argument for 'list' command")
				log.Fatal("Invalid argument for 'list' command")
			}
		}
	},
}

func main() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(markIPCmd)
	rootCmd.AddCommand(markDCmd)
	rootCmd.AddCommand(listCmd)

	if err := updateQuant(); err != nil {
		log.Fatalf("Fail to update quantity: %v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		// handleError(err, "Fail to execute command")
		log.Fatal("Fail to execute command")
	}
}

//	func handleError(err error, message string) {
//		if err != nil {
//			log.Fatalf("Error:%v\n", err)
//		} else {
//			log.Fatalf("Error:%v\n", message)
//		}
//	}
func CreateNewJson() error {
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
		err = CreateNewJson()

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

func updateQuant() error {
	if _, err := os.Stat(file_path); os.IsNotExist(err) {
		err = CreateNewJson()

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

func addItem(newContent string) error {
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

func deleteItem(serial int) error {
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

func updateItem(serial int, newContent string) error {
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

func markDone(serial int) error {
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

func markInProgress(serial int) error {
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

func list() error {
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

func listDone() error {
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

func listTodo() error {
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

func listInProgress() error {
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
