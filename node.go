package main

import (
	"github.com/google/uuid"
)

// Node is responsible to voting
type Node struct {
	ID        uuid.UUID
	QuorumSet map[string]*Node
}

// CreateNewNode creates a new node
func CreateNewNode() *Node {
	id, _ := uuid.NewUUID()
	node := &Node{
		ID:        id,
		QuorumSet: make(map[string]*Node)}

	return node
}
