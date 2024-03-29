package categories

import (
	"time"

	"gorm.io/gorm"
)

type TblCategory struct {
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

type Filter struct {
	Keyword  string
	Category string
	Status   string
	FromDate string
	ToDate   string
}

type Result struct {
	CategoryName string
}

type Arrangecategories struct {
	Categories []CatgoriesOrd
}

type CatgoriesOrd struct {
	Id       int
	Category string
}

type CategoryCreate struct {
	Id           int
	CategoryName string
	CategorySlug string
	Description  string
	ImagePath    string
	ParentId     int
}

type TblChannelCategory struct {
	Id           int
	ChannelId    int
	CategoriesId string
	CreatedAt    int
	CreatedOn    time.Time
}

// Parent Category List
func (c Authstruct) GetCategoryList(categories []TblCategory, offset int, limit int, filter Filter, DB *gorm.DB) (category []TblCategory, count int64) {

	var categorycount int64

	query := DB.Table("tbl_categories").Where("is_deleted = 0 and parent_id=0").Order("id desc")

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

// Children Category List
func (c Authstruct) GetSubCategoryList(categories *[]TblCategory, offset int, limit int, filter Filter, parent_id int, flag int, DB *gorm.DB) (categorylist *[]TblCategory, count int64) {

	var categorycount int64

	res := `WITH RECURSIVE cat_tree AS (
		SELECT id, category_name, category_slug,image_path, parent_id,created_on,modified_on,is_deleted
		FROM tbl_categories
		WHERE id = ?
		UNION ALL
		SELECT cat.id, cat.category_name, cat.category_slug, cat.image_path ,cat.parent_id,cat.created_on,cat.modified_on,
		cat.is_deleted
		FROM tbl_categories AS cat
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

func (c Authstruct) CreateCategory(category *TblCategory, DB *gorm.DB) error {

	if err := DB.Create(&category).Error; err != nil {

		return err
	}

	return nil
}

// Update Children list
func (c Authstruct) UpdateCategory(category *TblCategory, DB *gorm.DB) error {

	if category.ParentId == 0 && category.ImagePath == "" {

		if err := DB.Table("tbl_categories").Where("id = ?", category.Id).UpdateColumns(map[string]interface{}{"category_name": category.CategoryName, "category_slug": category.CategorySlug, "description": category.Description, "modified_by": category.ModifiedBy, "modified_on": category.ModifiedOn}).Error; err != nil {

			return err
		}
	} else {
		if err := DB.Table("tbl_categories").Where("id = ?", category.Id).UpdateColumns(map[string]interface{}{"category_name": category.CategoryName, "parent_id": category.ParentId, "category_slug": category.CategorySlug, "description": category.Description, "image_path": category.ImagePath, "modified_by": category.ModifiedBy, "modified_on": category.ModifiedOn}).Error; err != nil {

			return err
		}
	}

	return nil
}

// delete sub category
func (c Authstruct) DeletePopup(category *TblCategory, id int, DB *gorm.DB) (categorylist TblCategory, err error) {

	if err := DB.Table("tbl_categories").Where("parent_id=? and is_deleted =0", id).First(category).Error; err != nil {
		return TblCategory{}, err
	}
	return *category, nil
}

func (c Authstruct) DeleteCategoryById(category *TblCategory, categoryId int, DB *gorm.DB) error {

	if err := DB.Model(&category).Where("id=?", categoryId).Updates(TblCategory{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

/*getCategory Details*/
func (c Authstruct) GetCategoryById(category *TblCategory, categoryId int, DB *gorm.DB) (categorylist TblCategory, err error) {

	if err := DB.Table("tbl_categories").Where("is_deleted=0 and id=?", categoryId).First(&category).Error; err != nil {

		return TblCategory{}, err
	}
	return *category, nil
}

// Get Childern list
func (c Authstruct) GetCategoryDetails(id int, category *TblCategory, DB *gorm.DB) (categorylist TblCategory, err error) {

	if err := DB.Table("tbl_categories").Where("id=?", id).First(&category).Error; err != nil {

		return TblCategory{}, err
	}
	return *category, nil

}

func (c Authstruct) GetChildCategoriesById(childCategories *[]TblCategory, parentId int, DB *gorm.DB) error {

	if err := DB.Model(&childCategories).Where("is_deleted=0 and parent_id=?", parentId).Find(&childCategories).Error; err != nil {

		return err
	}
	return nil
}

func (c Authstruct) GetAllParentCategory(categories *[]TblCategory, DB *gorm.DB) error {

	if err := DB.Table("tbl_categories").Where("parent_id=0 and is_deleted=0").Find(&categories).Error; err != nil {

		return err
	}
	return nil
}

func (c Authstruct) GetCategoryTree(categoryID int, DB *gorm.DB) ([]TblCategory, error) {
	var categories []TblCategory
	err := DB.Raw(`
		WITH RECURSIVE cat_tree AS (
			SELECT id, 	CATEGORY_NAME,
			CATEGORY_SLUG,
			PARENT_ID,
			CREATED_ON,
			MODIFIED_ON,
			IS_DELETED
			FROM tbl_categories
			WHERE id = ?
			UNION ALL
			SELECT cat.id, CAT.CATEGORY_NAME,
			CAT.CATEGORY_SLUG,
			CAT.PARENT_ID,
			CAT.CREATED_ON,
			CAT.MODIFIED_ON,
			CAT.IS_DELETED
			FROM tbl_categories AS cat
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

func GetChildPageCategoriess(category *TblCategory, category_id int, DB *gorm.DB) (error, TblCategory) {

	if err := DB.Table("tbl_categories").Where("is_deleted=0 and id=?", category_id).Find(&category).Error; err != nil {

		return err, TblCategory{}
	}

	return nil, *category
}

func (c Authstruct) GetParentPageCategorys(category *[]TblCategory, parent_id int, DB *gorm.DB) (error, []TblCategory) {

	if err := DB.Table("tbl_categories").Where("is_deleted = 0 and id=?", parent_id).Find(&category).Error; err != nil {

		return err, []TblCategory{}
	}

	return nil, *category
}

// Check category group name already exists
func (c Authstruct) CheckCategoryGroupName(category TblCategory, userid int, name string, DB *gorm.DB) error {

	if userid == 0 {

		if err := DB.Model(TblCategory{}).Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and is_deleted=0", name).First(&category).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Model(TblCategory{}).Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and id not in (?) and is_deleted=0", name, userid).First(&category).Error; err != nil {

			return err
		}
	}

	return nil
}

// Check sub category name already exists
func (c Authstruct) CheckSubCategoryName(category TblCategory, categoryid []int, currentid int, name string, DB *gorm.DB) error {

	if len(categoryid) == 0 {

		if err := DB.Model(TblCategory{}).Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and is_deleted=0", name).First(&category).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Model(TblCategory{}).Where("LOWER(TRIM(category_name))=LOWER(TRIM(?)) and id in (?) and id not in (?) and is_deleted=0", name, categoryid, currentid).First(&category).Error; err != nil {

			return err
		}
	}

	return nil
}

func (c Authstruct) DeleteallCategoryById(category *TblCategory, categoryId []int, spacecatid int, DB *gorm.DB) error {

	if err := DB.Debug().Table("tbl_spaces").Where("page_category_id", spacecatid).Updates(TblCategory{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	if err := DB.Debug().Table("tbl_categories").Where("id in(?)", categoryId).Updates(TblCategory{IsDeleted: category.IsDeleted, DeletedOn: category.DeletedOn, DeletedBy: category.DeletedBy}).Error; err != nil {

		return err

	}

	return nil
}

func (c Authstruct) DeleteEntriesCategoryids(cid string, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entries").Where("categories_id LIKE ?", "%"+cid+"%").Update("categories_id", gorm.Expr("REPLACE(categories_id, ?, '')", cid)).Error; err != nil {
		return err
	}

	return nil
}

func (c Authstruct) DeleteChannelCategoryids(channelcategory *TblChannelCategory, cid string, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_category").Where("category_id LIKE ?", "%"+cid+"%").Delete(&channelcategory).Error; err != nil {
		return err
	}

	return nil
}

// update imagepath
func (c Authstruct) UpdateImagePath(Imagepath string, DB *gorm.DB) error {

	if err := DB.Model(TblCategory{}).Where("image_path=?", Imagepath).UpdateColumns(map[string]interface{}{
		"image_path": ""}).Error; err != nil {

		return err
	}

	return nil

}

func (c Authstruct) GetParentTreeByChild(childID int, DB *gorm.DB) ([]TblCategory, error) {
	var categories []TblCategory
	err := DB.Raw(`
		WITH RECURSIVE cat_tree AS (
			SELECT id, 	CATEGORY_NAME,
			CATEGORY_SLUG,
			PARENT_ID,
			CREATED_ON,
			MODIFIED_ON,
			IS_DELETED
			FROM tbl_categories
			WHERE id = ?
			UNION ALL
			SELECT cat.id, CAT.CATEGORY_NAME,
			CAT.CATEGORY_SLUG,
			CAT.PARENT_ID,
			CAT.CREATED_ON,
			CAT.MODIFIED_ON,
			CAT.IS_DELETED
			FROM tbl_categories AS cat
			JOIN cat_tree ON cat.id = cat_tree.parent_id
		)
		SELECT *
		FROM cat_tree WHERE IS_DELETED = 0 order by id desc
	`, childID).Scan(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}
