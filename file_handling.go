package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/cpgillem/garden-planner/models"
)

// Reads a plan from a reader.
func ReadObject[T any](r io.ReadCloser) (*T, error) {
	defer r.Close()

	// Read bytes.
	content, err := io.ReadAll(r)
	if err != nil {
		fmt.Println("Could not read data.\n" + err.Error())
		return nil, err
	}

	// Decode from JSON.
	return DecodeObject[T](content)
}

// Load a plan from a byte slice.
func DecodeObject[T any](content []byte) (*T, error) {
	var o T
	err := json.Unmarshal(content, &o)
	if err != nil {
		fmt.Println("Could not parse JSON.\n" + err.Error())
	}

	return &o, nil
}

// Encodes a plan to JSON.
func EncodeObject[T any](o *T) ([]byte, error) {
	content, err := json.Marshal(o)
	if err != nil {
		fmt.Println("Could not encode to JSON.\n" + err.Error())
		return nil, err
	}

	return content, nil
}

// Writes a plan to a writer.
func WriteObject[T any](w io.WriteCloser, o T) error {
	defer w.Close()

	// Encode JSON.
	content, err := EncodeObject(&o)
	if err != nil {
		return err
	}

	// Write plan.
	_, err = w.Write(content)
	if err != nil {
		fmt.Println("Could not write object.\n" + err.Error())
	}

	return nil
}

// Save a garden plan or create a new one.
func SavePlanToFile(plan *models.Plan, path string) error {
	// Encode to JSON.
	content, _ := EncodeObject(plan)

	// Create/truncate file.
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("Could not open file to save to.\n" + err.Error())
		return err
	}
	defer f.Close()

	// Write to file.
	_, err = f.Write(content)
	if err != nil {
		fmt.Println("Could not save file.\n" + err.Error())
		return err
	}
	f.Sync()

	return nil
}

// Open a file as a garden plan.
func LoadPlanFromFile(path string) (*models.Plan, error) {
	// Open file for reading.
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Could not load plan from " + path + "\n" + err.Error())
		return nil, err
	}
	defer f.Close()

	// Read content into byte slice.
	content, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Could not read file: " + path + "\n" + err.Error())
		return nil, err
	}

	// Decode plan.
	return DecodeObject[models.Plan](content)
}
