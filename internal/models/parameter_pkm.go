package models

// ParameterPKM represents db_parameter_pkm table
// This stores the actual answers/values from mahasiswa for each parameter
type ParameterPKM struct {
	ID            int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDPengajuan   int    `gorm:"column:id_pengajuan;type:int" json:"id_pengajuan"`       // FK ke db_pengajuan
	IDParameter   int    `gorm:"column:id_parameter;type:int" json:"id_parameter"`       // FK ke db_parameter_form
	Nilai         string `gorm:"column:nilai;type:text" json:"nilai"`                    // Jawaban mahasiswa
	
	// Relations
	Pengajuan     *Pengajuan     `gorm:"foreignKey:IDPengajuan" json:"-"`
	ParameterForm *ParameterForm `gorm:"foreignKey:IDParameter" json:"parameter_form,omitempty"`
}

// TableName specifies the table name for ParameterPKM model
func (ParameterPKM) TableName() string {
	return "db_parameter_pkm"
}