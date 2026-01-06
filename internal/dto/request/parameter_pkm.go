package request

// ParameterRequest represents form parameter answer in pengajuan request
type ParameterRequest struct {
	IDParameter int    `json:"id_parameter" validate:"required"`
	Nilai       string `json:"nilai" validate:"required"`
}