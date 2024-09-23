package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/1260644864/TaskTracker/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task-cli",
	Short: "A simple CLI to manage tasks",
}

var addCmd = &cobra.Command{
	Use:   "add <taskContent>",
	Short: "To add a task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.AddItem(args[0])
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

		err = internal.UpdateItem(serial, args[1])
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

		err = internal.DeleteItem(serial)
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

		err = internal.MarkInProgress(serial)
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

		err = internal.MarkDone(serial)
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
			err := internal.List()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			switch args[0] {
			case "done":
				err := internal.ListDone()
				if err != nil {
					log.Fatal(err)
				}
			case "todo":
				err := internal.ListTodo()
				if err != nil {
					log.Fatal(err)
				}
			case "in-progress":
				err := internal.ListInProgress()
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

	if err := internal.UpdateQuant(); err != nil {
		log.Fatalf("Fail to update quantity: %v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		// handleError(err, "Fail to execute command")
		log.Fatal("Fail to execute command")
	}
}
