package services

type Service string

const (
	Blob  Service = "b"
	Queue Service = "q"
	Table Service = "t"
	File  Service = "f"
)
