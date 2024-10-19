package client

type Response struct {
	Page          int        `json:"page"`
	PageSize      int        `json:"pageSize"`
	SortKey       string     `json:"sortKey"`
	SortDirection string     `json:"sortDirection"`
	TotalRecords  int        `json:"totalRecords"`
	Records       []Download `json:"records"`
}

type Download struct {
	MovieID                 int      `json:"movieId"`
	Size                    int64    `json:"size"`
	Title                   string   `json:"title"`
	SizeLeft                int64    `json:"sizeleft"`
	TimeLeft                string   `json:"timeleft"`
	EstimatedCompletionTime string   `json:"estimatedCompletionTime"`
	Added                   string   `json:"added"`
	Status                  string   `json:"status"`
	TrackedStatus           string   `json:"trackedDownloadStatus"`
	TrackedState            string   `json:"trackedDownloadState"`
	StatusMessages          []string `json:"statusMessages"`
	DownloadID              string   `json:"downloadId"`
	Protocol                string   `json:"protocol"`
	DownloadClient          string   `json:"downloadClient"`
	OutputPath              string   `json:"outputPath"`
	ID                      int      `json:"id"`
}
