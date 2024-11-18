package client

type Response struct {
	TotalRecords int        `json:"totalRecords"`
	Records      []Download `json:"records"`
	//Page          int        `json:"page"`
	//PageSize      int        `json:"pageSize"`
	//SortKey       string     `json:"sortKey"`
	//SortDirection string     `json:"sortDirection"`
}

type Download struct {
	Size     int64  `json:"size"`
	Title    string `json:"title"`
	SizeLeft int64  `json:"sizeleft"`
	Added    string `json:"added"`
	Status   string `json:"status"`
	ID       int    `json:"id"`
	//MovieID                 int      `json:"movieId"`
	//TimeLeft                string   `json:"timeleft"`
	//EstimatedCompletionTime string   `json:"estimatedCompletionTime"`
	//TrackedStatus           string   `json:"trackedDownloadStatus"`
	//TrackedState            string   `json:"trackedDownloadState"`
	//StatusMessages          []string `json:"statusMessages"`
	//DownloadID              string   `json:"downloadId"`
	//Protocol                string   `json:"protocol"`
	//DownloadClient          string   `json:"downloadClient"`
	//OutputPath              string   `json:"outputPath"`
}
