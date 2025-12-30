package domain

type Project struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Task  []Task `json:"task"`
}

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}
