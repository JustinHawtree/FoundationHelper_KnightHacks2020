package controller

import (
	"bytes"
	"context"
	"FoundationHelper_KnightHacks2020/model"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

// SyncBuf data structure
type SyncBuf struct {
	mu       sync.Mutex
	buf      bytes.Buffer
	output   string
	overflow bool
}

// Reset function resets the buffer
func (s *SyncBuf) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buf.Reset()
}

func (s *SyncBuf) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Println("Buff Length:", s.buf.Len())
	if s.buf.Len() > 10000 {
		s.overflow = true
		s.buf.Truncate(10000)
	}
	return s.buf.String()
}

func (s *SyncBuf) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

// RunCode function to get a bike
func RunCode(c *fiber.Ctx) error {
	start := time.Now()
	// New incoming code
	userCode := new(model.UserCode)

	// Parse body into struct
	if bodyErr := c.BodyParser(userCode); bodyErr != nil {
		return c.Status(400).SendString(bodyErr.Error())
	}
	//log.Println(userCode.Input)

	f, fileErr := os.Create("hello.c")
	if fileErr != nil {
		return c.Status(500).SendString(fileErr.Error())
	}
	_, writeErr := f.WriteString(userCode.Input)
	if writeErr != nil {
		return c.Status(500).SendString(writeErr.Error())
	}

	cmd := exec.Command("gcc", "hello.c")
	out, err := cmd.CombinedOutput()

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": string(out)})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd = exec.CommandContext(ctx, "./a.out")
	var buf SyncBuf
	cmd.Stdout = &buf

	err = cmd.Run()

	if err != nil {
		log.Println("Timed out?")
		//log.Println(buf.String())
		buf.Reset()
		buf = SyncBuf{}
		return c.Status(400).JSON(fiber.Map{"error": "Process Timed Out (Code exceeded 3 Seconds)"})
	}

	done := time.Now()
	elapsed := done.Sub(start)

	result := fiber.Map{"output": buf.String(), "time_taken": elapsed, "overflow": buf.overflow}
	buf.Reset()
	return c.JSON(result)
}


func GetCode(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{"problem":
`#include <stdlib.h>
#include <stdio.h>

typedef struct node {
	int data;
	struct node* next, *prev;
} node;
		 
typedef struct queue {
	int size;
	struct node *head, *tail;
} queue;
		
int dequeue(queue*thisQ) {

	// Insert Code Here










}`, "problem_topic":"Linked Lists", "problem_name":"Fall 2020 Part A", "problem_summary":"Suppose we have a queue implemented as a doubly linked list using the structures shown below.  Use head for the front of the queue and tail for the end of the queue.\n\nstruct node {\n\tint data;\n\tstruct node* next, *prev;\n}\n\nstruct queue {\n\tint size;\n\tstruct node *head, *tail;\n}\n\nWrite a dequeue function for this queue. If the queue is NULL or is already empty, return 0 and take no other action. If the queue isn't empty, dequeue the appropriate value, make the necessary adjustments, and return the dequeued value. (Note: You must free the node that previously stored the dequeued value.)"})
}