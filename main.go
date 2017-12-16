package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"encoding/json"
	"os"
	"strings"
)

// Common section begin
func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

// Common section end

// DAO section begin
type DataStore interface {
	IntentExists(intent string) bool
	Put(intent string, key string, value string)
	Get(intent string, key string) string
}

type JSONStore struct {
	DataDIR       string `default:"/var/e2c/data"`
	DefaultIntent string `default:"default"`
}

func (js *JSONStore) LoadFile(intent string) *os.File {
	file, err := os.OpenFile(js.DataDIR+"/"+intent+".json", os.O_RDWR, os.ModePerm)
	if os.IsNotExist(err) {
		// fmt.Println(err.Error())
		file, err = os.Create(js.DataDIR + "/" + intent + ".json")
	}
	CheckError(err)
	return file
}

func (js *JSONStore) IntentExists(intent string) bool {
	if _, err := os.Stat(js.DataDIR + "/" + intent + ".json"); os.IsNotExist(err) {
		return false
	}
	return true
}

func (js *JSONStore) LoadJSON(file *os.File) *map[string]string {
	data := make(map[string]string)
	err := json.NewDecoder(file).Decode(&data)
	CheckError(err)
	// fmt.Println(data)
	return &data
}

func (js *JSONStore) SaveJSON(file *os.File, data *map[string]string) {
	// file := LoadFile(intent)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(data)
	CheckError(err)
	defer file.Close()
}

func (js *JSONStore) Put(intent string, key string, value string) {
	file := js.LoadFile(intent)
	data := js.LoadJSON(file)
	file.Seek(0, 0)
	(*data)[key] = value
	js.SaveJSON(file, data)
	defer file.Close()
}

func (js *JSONStore) Get(intent string, key string) string {
	file := js.LoadFile(intent)
	data := js.LoadJSON(file)
	defer file.Close()
	value := (*data)[key]
	return value
}

// DAO section end

// Dict section begin
func NewDataStore() DataStore {
	return &JSONStore{DataDIR: "json/data"}
}

func GetValue(intent string, key string) string {
	store := NewDataStore()
	key = strings.ToLower(key)
	if store.IntentExists(intent) {
		return store.Get(intent, key)
	} else {
		return store.Get("default", key)
	}
}

func SetValue(intent string, key string, value string) {
	store := NewDataStore()
	key = strings.ToLower(key)
	if store.IntentExists(intent) {
		store.Put(intent, key, value)
	} else {
		store.Put("default", key, value)
	}
}

func GetValueFromDefault(key string) string {
	return GetValue("default", key)
}

func SetValueToDefault(key string, value string) {
	SetValue("default", key, value)
}

// Dict section end

// Command section end
func main() {
	var cmdGet = &cobra.Command{
		Use:   "get <key>",
		Short: "读取翻译值",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(GetValueFromDefault(args[0]))
		},
	}

	var cmdPut = &cobra.Command{
		Use:   "put <key> <value>",
		Short: "添加翻译值",
		Run: func(cmd *cobra.Command, args []string) {
			SetValueToDefault(args[0], args[1])
		},
	}

	var rootCmd = &cobra.Command{Use: "e2c"}
	rootCmd.AddCommand(cmdGet)
	rootCmd.AddCommand(cmdPut)

	rootCmd.Execute()
}
