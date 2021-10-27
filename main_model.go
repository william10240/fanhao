package main

type fanhao struct {
	ID int `gorm:"primaryKey"`
	Code string `gorm:"type:string;size:64;not null;index:code_idx"`
	Title string `gorm:"type:string;size:128"`
	Star string `gorm:"type:string;size:64"`
	StarCode string `gorm:"type:string;size:64"`
	Img string `gorm:"type:string;size:128"`
	Fname string `gorm:"type:string;size:128"`
	Ima int `gorm:"type:int;default:0"`
	Iface int `gorm:"type:int;default:0"`
	Starnum int `gorm:"type:int;default:0"`
	Downed int `gorm:"type:int;default:0"`
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
}
func (f *fanhao) TableName() string{
	return "fanhao"
}