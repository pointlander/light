// Copyright 2025 The Light Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Prompt is a llm prompt
type Prompt struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// Query submits a query to the llm
func Query(query string) string {
	prompt := Prompt{
		Model:  "llama3.2",
		Prompt: query,
	}
	data, err := json.Marshal(prompt)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(data)
	response, err := http.Post("http://10.0.0.54:11434/api/generate", "application/json", buffer)
	if err != nil {
		panic(err)
	}
	reader, answer := bufio.NewReader(response.Body), ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		data := map[string]interface{}{}
		err = json.Unmarshal([]byte(line), &data)
		text := data["response"].(string)
		answer += text
	}
	return answer
}

//go:embed prompts/*
var Prompts embed.FS

func main() {
	file, err := Prompts.Open("prompts/1.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	raw, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	prompt := fmt.Sprintf(string(raw), "the number 1337")
	goja := NewGOJA()
	const (
		begin = "javascript"
		end   = "```"
	)
	result, i := Query(prompt), 0
	for {
		index := strings.Index(result, begin)
		if index == -1 {
			fmt.Print(result)
			break
		}
		fmt.Print(result[:index+len(begin)])
		result = result[index+len(begin):]
		index = strings.Index(result, end)
		fmt.Println(result[:index+len(end)])
		fmt.Println("```goja")
		err := goja.Run(i, result[:index])
		if err != nil {
			fmt.Print("<<<")
			fmt.Println(err)
		}
		i++
		fmt.Println("```")
		result = result[index+len(end):]
	}
}
