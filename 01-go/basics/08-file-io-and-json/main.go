package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Role  string `json:"-"`
}

func main() {
	dir, err := os.MkdirTemp("", "brm-io-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	u := User{ID: 1, Name: "alice", Email: "alice@example.com", Role: "admin"}
	b, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	fmt.Println("marshal:", string(b))
	fmt.Println("  Role 被 json:\"-\" 忽略；Email 为空时 omitempty 会去掉该字段")

	path := filepath.Join(dir, "users.jsonl")
	if err := os.WriteFile(path, append(b, '\n'), 0o600); err != nil {
		panic(err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read back: %d bytes\n", len(raw))

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var decoded User
		if err := json.Unmarshal(scanner.Bytes(), &decoded); err != nil {
			panic(err)
		}
		fmt.Printf("unmarshal: id=%d name=%s email=%s\n", decoded.ID, decoded.Name, decoded.Email)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}

	// 大文件用流式解码，别一次性 ReadFile 进内存
	stream := `{"id":2,"name":"bob"}` + "\n" + `{"id":3,"name":"carol"}` + "\n"
	dec := json.NewDecoder(strings.NewReader(stream))
	for dec.More() {
		var decoded User
		if err := dec.Decode(&decoded); err != nil {
			panic(err)
		}
		fmt.Println("stream decode:", decoded.Name)
	}

	indented, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("indent:\n" + string(indented))
}
