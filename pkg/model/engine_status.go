package model

import (
	"time"
)

const (
	EngineStatusFetching = "fetching"
	EngineStatusSuccess  = "success"
	EngineStatusFailed   = "failed"
)

type EngineStatus struct {
	ErrorMessage string    `json:"error_message"`
	Status       string    `json:"status"`
	FetchedAt    time.Time `json:"fetched_at"`
}

func MakeEngineStatusFetching() EngineStatus {
	return EngineStatus{
		ErrorMessage: "",
		Status:       EngineStatusFetching,
		FetchedAt:    time.Now(),
	}
}

func MakeEngineStatusSuccess() EngineStatus {
	return EngineStatus{
		ErrorMessage: "",
		Status:       EngineStatusSuccess,
		FetchedAt:    time.Now(),
	}
}

func MakeEngineStatusFailed(err error) EngineStatus {
	return EngineStatus{
		ErrorMessage: err.Error(),
		Status:       EngineStatusFailed,
		FetchedAt:    time.Now(),
	}
}
