package request

// ReviewJudulRequest represents request body for reviewing PKM title
type ReviewJudulRequest struct {
	IDStatusReview int    `json:"id_status_review" validate:"required"` // FK to db_status_review (1=PENDING, 2=ON_REVIEW, 3=ACC, 4=REVISI, 5=TOLAK)
	Catatan        string `json:"catatan" validate:"required,min=10"`
}

// ReviewProposalRequest represents request body for reviewing PKM proposal
type ReviewProposalRequest struct {
	IDStatusReview int    `json:"id_status_review" validate:"required"` // FK to db_status_review
	Catatan        string `json:"catatan" validate:"required,min=10"`
}

// AssignReviewerRequest represents request body for admin to assign reviewer
type AssignReviewerRequest struct {
	IDPegawai int `json:"id_pegawai" validate:"required"`
}

// AnnounceRequest represents request body for admin to announce final result
type AnnounceRequest struct {
	StatusFinal string `json:"status_final" validate:"required,oneof=LOLOS TIDAK_LOLOS"`
}