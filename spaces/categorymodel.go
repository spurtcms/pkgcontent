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

type CategoryCreate struct {
	Id           int
	CategoryName string
	CategorySlug string
	Description  string
	ImagePath    string
	ParentId     int
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

func (c Authstruct) CreateSpaceCategory(category *TblSpaceCategory, DB *gorm.DB) error {

	if err := DB.Create(&category).Error; err != nil {

		return err
	}

	return nil
}

// Update Children list
func (c Authstruct) UpdateCategory(category *TblSpaceCategory, DB *gorm.DB) error {

	if category.ParentId == 0 && category.ImagePath == "" {

		if err := DB.Model(TblSpaceCategory{}).Where("id = ?", category.Id).UpdateColumns(map[string]interface{}{"category_name": category.CategoryName, "category_slug": category.CategorySlug, "description": category.Description, "modified_by": category.ModifiedBy, "modified_on": category.ModifiedOn}).Error; err != nil {

			return err
		}
	} else {
		if err := DB.Model(TblSpaceCategory{}).Where("id = ?", category.Id).UpdateColumns(map[string]interface{}{"category_name": category.CategoryName, "parent_id": category.ParentId, "category_slug": category.CategorySlug, "description": category.Description, "image_path": category.ImagePath, "modified_by": category.ModifiedBy, "modified_on": category.ModifiedOn}).Error; err != nil {

			return err
		}
	}

	return nil
}

func (c Authstruct) DeleteallCategoryById(category *TblSpaceCategory, categoryId []int, spacecatid int, DB *gorm.DB) error {

	if err := DB.Model(TblSpaces{}).Where("page_category_id", spacecatid).Updates(TblSpaceCategory{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	if err := DB.Model(TblSpaceCategory{}).Where("id in(?)", categoryId).Updates(TblSpaceCategory{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

func (c Authstruct) GetCategoryTree(categoryID int, DB *gorm.DB) ([]TblSpaceCategory, error) {
	var categories []TblSpaceCategory
	err := DB.Raw(`
		WITH RECURSIVE cat_tree AS (
			SELECT id, 	CATEGORY_NAME,
			CATEGORY_SLUG,
			PARENT_ID,
			CREATED_ON,
			MODIFIED_ON,
			IS_DELETED
			FROM tbl_space_categories
			WHERE id = ?
			UNION ALL
			SELECT cat.id, CAT.CATEGORY_NAME,
			CAT.CATEGORY_SLUG,
			CAT.PARENT_ID,
			CAT.CREATED_ON,
			CAT.MODIFIED_ON,
			CAT.IS_DELETED
			FROM tbl_space_categories AS cat
			JOIN cat_tree ON cat.parent_id = cat_tree.id
		)
		SELECT *
		FROM cat_tree WHERE IS_DELETED = 0 order by id desc
	`, categoryID).Scan(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}
