package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

// Load an object from a byte slice.
func DecodeObject[T any](content []byte) (*T, error) {
	var o T
	err := json.Unmarshal(content, &o)
	if err != nil {
		fmt.Println("Could not parse JSON.\n" + err.Error())
		return nil, err
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
func WriteObjectToFile[T any](o *T, path string) error {
	// Create/truncate file.
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("Could not open file to save to.\nPath: " + path + "\n" + err.Error())
		return err
	}
	defer f.Close()

	// Write to file.
	return WriteObject(f, o)
}

// Open a file as an object.
func ReadObjectFromFile[T any](path string) (*T, error) {
	// Open file for reading.
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Could not load from " + path + "\n" + err.Error())
		return nil, err
	}
	defer f.Close()

	return ReadObject[T](f)
}
