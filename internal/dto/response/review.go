package response

import "time"

// ReviewResponse represents review history data
type ReviewResponse struct {
	ID             int              `json:"id"`
	TipeReview     string           `json:"tipe_review"` // JUDUL or PROPOSAL
	StatusReview   string           `json:"status_review"` // PENDING, ON_REVIEW, ACC, REVISI, TOLAK
	Catatan        string           `json:"catatan"`
	TglReview      *time.Time       `json:"tgl_review"`
	Reviewer       *PegawaiResponse `json:"reviewer,omitempty"`
}

// PlottingResponse represents reviewer assignment data
type PlottingResponse struct {
	ID          int              `json:"id"`
	Tipe        string           `json:"tipe"` // JUDUL or PROPOSAL
	Status      string           `json:"status"` // ASSIGNED, REVIEWED
	TglAssign   *time.Time       `json:"tgl_assign"`
	Reviewer    *PegawaiResponse `json:"reviewer,omitempty"`
}