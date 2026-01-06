package request

// UploadProposalRequest represents request for uploading proposal
// File will be handled separately via multipart/form-data
type UploadProposalRequest struct {
	IDPengajuan int `json:"id_pengajuan" validate:"required"`
	// File will be uploaded via c.FormFile("file")
}

// Note: No struct needed for file upload, just validation in controller