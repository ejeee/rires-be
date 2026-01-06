package models

// import "time"

// ParameterPKM represents db_parameter_form table
type ParameterPKM struct {
	ID            int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDPengajuan   int        `gorm:"column:id_pengajuan;type:int" json:"id_pengajuan"`
	IDParameter   int        `gorm:"column:id_parameter;type:int" json:"id_parameter"` // FK ke db_kategori_pkm
	Nilai		 string     `gorm:"column:nilai;type:varchar(100)" json:"nilai"`

	// Relations
	Kategori *KategoriPKM `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
}

// TableName specifies the table name for ParameterPKM model
func (ParameterPKM) TableName() string {
	return "db_parameter_pkm"
}