package service

import (
	"fmt"
)

type Collector interface {
	GetSystemData() (string, error)
}

type collector struct {
}

func (c *collector) GetSystemData() (string, error) {
	return "", fmt.Errorf("test")
}
