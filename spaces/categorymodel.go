package lms

import (
	"time"

	"gorm.io/gorm"
)

type TblSpaceCategory struct {
	Id                 int
	CategoryName       string
	CategorySlug       string
	Description        string
	ImagePath          string
	CreatedOn          time.Time
	CreatedBy          int
	ModifiedOn         time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy         int       `gorm:"DEFAULT:NULL"`
	IsDeleted          int
	DeletedOn          time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy          int       `gorm:"DEFAULT:NULL"`
	ParentId           int
	CreatedDate        string   `gorm:"-"`
	ModifiedDate       string   `gorm:"-"`
	DateString         string   `gorm:"-"`
	ParentCategoryName string   `gorm:"-"`
	Parent             []string `gorm:"-"`
	ParentWithChild    []Result `gorm:"-"`
}

type Result struct {
	CategoryName string
}

// Parent Category List
func (c Authstruct) GetSpaceCategoryList(categories []TblSpaceCategory, offset int, limit int, filter Filter, DB *gorm.DB) (category []TblSpaceCategory, count int64) {

	var categorycount int64

	query := DB.Model(TblSpaceCategory{}).Where("is_deleted = 0 and parent_id=0").Order("id desc")

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(category_name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%")
	}

	if limit != 0 {

		query.Limit(limit).Offset(offset).Find(&categories)

		return categories, categorycount

	}

	query.Find(&categories).Count(&categorycount)

	return categories, categorycount

}
