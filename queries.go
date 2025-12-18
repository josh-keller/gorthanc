package gorthanc

type Level string

const (
	QueryLevelPatient  = "Patient"
	QueryLevelStudy    = "Study"
	QueryLevelSeries   = "Series"
	QueryLevelInstance = "Instance"
)

type Query struct {
	ID     string
	Level  Level
	client *Client
}

func (q *Query) RetrieveAllSynchronous()
