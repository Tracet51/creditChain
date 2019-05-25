package main

import (
	"github.com/google/uuid"
)

// Node is responsible to voting
type Node struct {
	ID        uuid.UUID
	QuorumSet map[string]*Node
}
