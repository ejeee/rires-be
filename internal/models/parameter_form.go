package models

import "time"

// ParameterForm represents db_parameter_form table
type ParameterForm struct {
	ID            int        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IDKategori    int        `gorm:"column:id_kategori;type:int" json:"id_kategori"`
	NamaParameter string     `gorm:"column:nama_parameter;type:varchar(100)" json:"nama_parameter"`
	Label         string     `gorm:"column:label;type:varchar(200)" json:"label"`
	TipeInput     string     `gorm:"column:tipe_input;type:varchar(50)" json:"tipe_input"` // text, textarea, number, file, radio, select, etc
	Validasi      string     `gorm:"column:validasi;type:text" json:"validasi"`            // JSON: {"required":true,"min":10}
	Placeholder   string     `gorm:"column:placeholder;type:text" json:"placeholder"`
	HelpText      string     `gorm:"column:help_text;type:text" json:"help_text"`
	Opsi          string     `gorm:"column:opsi;type:text" json:"opsi"` // JSON for radio/select options
	Urutan        int        `gorm:"column:urutan;type:int;default:0" json:"urutan"`
	Status        int        `gorm:"column:status;type:int(1);default:1" json:"status"` // 1=active, 2=inactive
	Hapus         int        `gorm:"column:hapus;type:int(1);default:0" json:"-"`       // 0=exists, 1=deleted
	TglInsert     *time.Time `gorm:"column:tgl_insert;type:datetime" json:"tgl_insert"`
	TglUpdate     time.Time  `gorm:"column:tgl_update;type:timestamp;autoUpdateTime" json:"tgl_update"`
	UserUpdate    string     `gorm:"column:user_update;type:text" json:"user_update"`

	// Relations
	Kategori *KategoriPKM `gorm:"foreignKey:IDKategori" json:"kategori,omitempty"`
}

// TableName specifies the table name for ParameterForm model
func (ParameterForm) TableName() string {
	return "db_parameter_form"
}
