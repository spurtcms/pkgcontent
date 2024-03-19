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

// Create Space Category
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

/*getCategory Details*/
func (c Authstruct) GetCategoryById(category *TblSpaceCategory, categoryId int, DB *gorm.DB) (categorylist TblSpaceCategory, err error) {

	if err := DB.Table("tbl_space_categories").Where("is_deleted=0 and id=?", categoryId).First(&category).Error; err != nil {

		return TblSpaceCategory{}, err
	}
	return *category, nil
}

// Children Category List
func (c Authstruct) GetSubSpaceCategoryList(categories *[]TblSpaceCategory, offset int, limit int, filter Filter, parent_id int, flag int, DB *gorm.DB) (categorylist *[]TblSpaceCategory, count int64) {

	var categorycount int64

	res := `WITH RECURSIVE cat_tree AS (
		SELECT id, category_name, category_slug,image_path, parent_id,created_on,modified_on,is_deleted
		FROM tbl_space_categories
		WHERE id = ?
		UNION ALL
		SELECT cat.id, cat.category_name, cat.category_slug, cat.image_path ,cat.parent_id,cat.created_on,cat.modified_on,
		cat.is_deleted
		FROM tbl_space_categories AS cat
		JOIN cat_tree ON cat.parent_id = cat_tree.id )`

	query := DB

	if filter.Keyword != "" {

		if limit == 0 {
			query.Raw(` `+res+` select count(*) from cat_tree where is_deleted = 0 and LOWER(TRIM(category_name)) ILIKE LOWER(TRIM(?)) group by cat_tree.id `, parent_id, "%"+filter.Keyword+"%").Count(&categorycount)

			return categories, categorycount
		}
		query = query.Raw(` `+res+` select * from cat_tree where is_deleted = 0 and LOWER(TRIM(category_name)) ILIKE LOWER(TRIM(?)) limit(?) offset(?) `, parent_id, "%"+filter.Keyword+"%", limit, offset)
	} else if flag == 0 {
		query = query.Raw(``+res+` SELECT * FROM cat_tree where is_deleted = 0 and id not in (?) order by id desc limit(?) offset(?) `, parent_id, parent_id, limit, offset)
	} else if flag == 1 {
		query = query.Raw(``+res+` SELECT * FROM cat_tree where is_deleted = 0 order by id desc `, parent_id)
	}
	if limit != 0 {

		query.Find(&categories)

		return categories, categorycount

	} else {

		DB.Raw(` `+res+` SELECT count(*) FROM cat_tree where is_deleted = 0 and id not in (?)  group by cat_tree.id order by id desc`, parent_id, parent_id).Count(&categorycount)

		return categories, categorycount
	}

}

// Delete Sub Category
func (c Authstruct) DeleteCategoryById(category *TblSpaceCategory, categoryId int, DB *gorm.DB) error {

	if err := DB.Model(&category).Where("id=?", categoryId).Updates(TblSpaceCategory{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

// delete sub category modal
func (c Authstruct) SpaceCategoryDeletePopup(category *TblSpaceCategory, id int, DB *gorm.DB) (categorylist TblSpaceCategory, err error) {

	if err := DB.Table("tbl_space_categories").Where("parent_id=? and is_deleted =0", id).First(category).Error; err != nil {

		return TblSpaceCategory{}, err
	}
	return *category, nil
}

// Get Childern list
func (c Authstruct) GetCategoryDetails(id int, category *TblSpaceCategory, DB *gorm.DB) (categorylist TblSpaceCategory, err error) {

	if err := DB.Table("tbl_space_categories").Where("id=?", id).First(&category).Error; err != nil {

		return TblSpaceCategory{}, err
	}
	return *category, nil

}

// Check category group name already exists
func (c Authstruct) CheckSpaceCategoryGroupName(category TblSpaceCategory, userid int, name string, DB *gorm.DB) error {

	if userid == 0 {

		if err := DB.Model(TblSpaceCategory{}).Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and is_deleted=0", name).First(&category).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Model(TblSpaceCategory{}).Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and id not in (?) and is_deleted=0", name, userid).First(&category).Error; err != nil {

			return err
		}
	}

	return nil
}

/*Get All cateogry with parents and subcategory*/

func (c Authstruct) GetAllParentCategory(categories *[]TblSpaceCategory, DB *gorm.DB) error {

	if err := DB.Table("tbl_space_categories").Where("parent_id=0 and is_deleted=0").Find(&categories).Error; err != nil {

		return err
	}
	return nil
}

// Check sub category name already exists
func (c Authstruct) CheckSubCategoryName(category TblSpaceCategory, categoryid []int, currentid int, name string, DB *gorm.DB) error {

	if len(categoryid) == 0 {

		if err := DB.Model(TblSpaceCategory{}).Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and is_deleted=0", name).First(&category).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Model(TblSpaceCategory{}).Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and id in (?) and id not in (?) and is_deleted=0", name, categoryid, currentid).First(&category).Error; err != nil {

			return err
		}
	}

	return nil
}

// update imagepath
func (c Authstruct) UpdateImagePath(Imagepath string, DB *gorm.DB) error {

	if err := DB.Model(TblSpaceCategory{}).Where("image_path=?", Imagepath).UpdateColumns(map[string]interface{}{
		"image_path": ""}).Error; err != nil {

		return err
	}

	return nil

}