package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNeighbors(t *testing.T) {
	c := [][]bool{
		{false, true, false},
		{false, false, false},
		{false, false, false},
	}

	assert.Equal(t, 1, getNeighbors(c, 1, 1))
	assert.Equal(t, 0, getNeighbors(c, 0, 1))
	assert.Equal(t, 0, getNeighbors(c, 2, 2))
	assert.Equal(t, 1, getNeighbors(c, 1, 0))
	assert.Equal(t, 1, getNeighbors(c, 0, 2))
	assert.Equal(t, 1, getNeighbors(c, 1, 2))
	assert.Equal(t, 1, getNeighbors(c, 0, 0))
}

func TestCanvasNext(t *testing.T) {
	c := [][]bool{
		{false, false, false},
		{true, true, true},
		{false, false, false},
	}

	next := canvasNext(c)

	assert.Equal(t, [][]bool{
		{false, true, false},
		{false, true, false},
		{false, true, false},
	}, next)
}
