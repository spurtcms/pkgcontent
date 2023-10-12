package spaces

import (
	"time"

	"gorm.io/datatypes"
)

type TblSpaces struct {
	Id             int `gorm:"primaryKey;auto_increment"`
	PageCategoryId int
	CreatedOn      time.Time
	CreatedBy      int
	ModifiedOn     time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy     int       `gorm:"DEFAULT:NULL"`
	DeletedOn      time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy      int       `gorm:"DEFAULT:NULL"`
	IsDeleted      int       `gorm:"DEFAULT:0"`
	Username       string    `gorm:"<-:false"`
	CreatedDate    string    `gorm:"-"`
	ModifiedDate   string    `gorm:"-"`
}

type TblSpacesAliases struct {
	Id                int `gorm:"primaryKey;auto_increment"`
	SpacesId          int
	LanguageId        int
	SpacesName        string
	SpacesSlug        string
	SpacesDescription string
	ImagePath         string
	CreatedOn         time.Time
	CreatedBy         int
	ModifiedOn        time.Time                   `gorm:"DEFAULT:NULL"`
	ModifiedBy        int                         `gorm:"DEFAULT:NULL"`
	DeletedOn         time.Time                   `gorm:"DEFAULT:NULL"`
	DeletedBy         int                         `gorm:"DEFAULT:NULL"`
	IsDeleted         int                         `gorm:"DEFAULT:0"`
	PageCategoryId    int                         `gorm:"column:page_category_id;<-:false"`
	ParentId          int                         `gorm:"column:parent_id;<-:false"`
	ParentCategory    TblPagesCategoriesAliases   `gorm:"-"`
	ChildCategories   []TblPagesCategoriesAliases `gorm:"-"`
	CreatedDate       string                      `gorm:"-"`
	ModifiedDate      string                      `gorm:"-"`
}

type TblPagesCategoriesAliases struct {
	Id                  int `gorm:"primaryKey;auto_increment"`
	PageCategoryId      int
	LanguageId          int
	CategoryName        string
	CategorySlug        string
	CategoryDescription string
	CreatedOn           time.Time
	CreatedBy           int
	ModifiedOn          time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy          int       `gorm:"DEFAULT:NULL"`
	DeletedOn           time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy           int       `gorm:"DEFAULT:NULL"`
	IsDeleted           int       `gorm:"DEFAULT:0"`
	ParentId            int
}

type TblPagesCategories struct {
	Id         int `gorm:"primaryKey;auto_increment"`
	CreatedOn  time.Time
	CreatedBy  int
	ModifiedOn time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy int       `gorm:"DEFAULT:NULL"`
	DeletedOn  time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy  int       `gorm:"DEFAULT:NULL"`
	IsDeleted  int       `gorm:"DEFAULT:0"`
}

type TblLanguage struct {
	Id           int `gorm:"primaryKey;auto_increment"`
	LanguageName string
	LanguageCode string
	CreatedOn    time.Time
	CreatedBy    int
	ModifiedOn   time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy   int       `gorm:"DEFAULT:NULL"`
	DeletedOn    time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy    int       `gorm:"DEFAULT:NULL"`
	IsDeleted    int       `gorm:"DEFAULT:0"`
	ImagePath    string
	IsStatus     int
	IsDefault    int
	JsonPath     string
}

type Filter struct {
	Keyword    string
	CategoryId []int
}

type Arrangecategories struct {
	Categories []CatgoriesOrd
}

type CatgoriesOrd struct {
	Id       int
	Category string
}

type TblPagesGroup struct {
	Id         int `gorm:"primaryKey;auto_increment"`
	SpacesId   int
	CreatedOn  time.Time
	CreatedBy  int
	ModifiedOn time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy int       `gorm:"DEFAULT:NULL"`
	DeletedOn  time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy  int       `gorm:"DEFAULT:NULL"`
	IsDeleted  int       `gorm:"DEFAULT:0"`
}
type TblPagesGroupAliases struct {
	Id               int `gorm:"primaryKey;auto_increment"`
	PageGroupId      int
	LanguageId       int
	GroupName        string
	GroupSlug        string
	GroupDescription string
	CreatedOn        time.Time
	CreatedBy        int
	ModifiedOn       time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy       int       `gorm:"DEFAULT:NULL"`
	DeletedOn        time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy        int       `gorm:"DEFAULT:NULL"`
	IsDeleted        int       `gorm:"DEFAULT:0"`
	OrderIndex       int
}

type TblPage struct {
	Id          int `gorm:"primaryKey;auto_increment"`
	SpacesId    int
	PageGroupId int
	ParentId    int
	CreatedOn   time.Time
	CreatedBy   int
	ModifiedOn  time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy  int       `gorm:"DEFAULT:NULL"`
	DeletedOn   time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy   int       `gorm:"DEFAULT:NULL"`
	IsDeleted   int       `gorm:"DEFAULT:0"`
}

type MetaDetails struct {
	MetaTitle       string
	MetaDescription string
	Keywords        string
	Slug            string
}

type TblPageAliases struct {
	Id              int `gorm:"primaryKey;auto_increment"`
	PageId          int
	LanguageId      int
	PageTitle       string
	PageSlug        string
	PageDescription string
	PublishedOn     time.Time `gorm:"DEFAULT:NULL"`
	Author          string
	Excerpt         string
	FeaturedImages  string
	Access          string
	MetaDetails     datatypes.JSONType[MetaDetails]
	Status          string
	AllowComments   bool
	CreatedOn       time.Time
	CreatedBy       int
	ModifiedOn      time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy      int       `gorm:"DEFAULT:NULL"`
	DeletedOn       time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy       int       `gorm:"DEFAULT:NULL"`
	IsDeleted       int       `gorm:"DEFAULT:0"`
	OrderIndex      int
	PageSuborder    int
	CreatedDate     string `gorm:"-"`
	ModifiedDate    string `gorm:"-"`
	Username        string `gorm:"-"`
	PageGroupId     int    `gorm:"<-:false"`
	ParentId        int    `gorm:"<-:false"`
}



/*spaceList*/
func SpaceList(tblspace *[]TblSpacesAliases, langId int, limit int, offset int, filter Filter, spaceid []int) (spacecount int64, err error) {

	query := s.Authority.DB.Table("tbl_spaces_aliases").Select("tbl_spaces_aliases.*,tbl_spaces.page_category_id,tbl_pages_categories.parent_id").
		Joins("inner join tbl_spaces on tbl_spaces_aliases.spaces_id = tbl_spaces.id").
		Joins("inner join tbl_language on tbl_language.id = tbl_spaces_aliases.language_id").
		Joins("inner join tbl_pages_categories on tbl_pages_categories.id = tbl_spaces.page_category_id").
		Where("tbl_spaces.is_deleted = 0 and tbl_spaces_aliases.is_deleted = 0 and tbl_spaces_aliases.language_id = ?", langId)

	if len(spaceid) != 0 {

		query = query.Where("tbl_spaces.id in (?)", spaceid)
	}

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(tbl_spaces_aliases.spaces_name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%")
	}
	if len(filter.CategoryId) > 0 && filter.CategoryId[0] != 0 {
		query = query.Where("tbl_spaces.page_category_id IN (?)", filter.CategoryId)
	}

	if limit != 0 {

		query.Limit(limit).Offset(offset).Order("tbl_spaces.id desc").Find(&tblspace)

	} else {

		query.Find(&tblspace).Count(&spacecount)

		return spacecount, nil
	}

	return 0, nil
}

/*get default lang id*/
func GetDefaultLanguage(default_lang *TblLanguage) error {

	if err := s.Authority.DB.Table("tbl_language").Where("is_deleted=0 and is_default=?", 1).First(&default_lang).Error; err != nil {

		return err
	}

	return nil
}

/*Create Space*/
func CreateSpace(tblspaces *TblSpaces) (*TblSpaces, error) {

	if err := s.Authority.DB.Table("tbl_spaces").Create(&tblspaces).Error; err != nil {

		return &TblSpaces{}, err
	}
	return tblspaces, nil

}

/*Create space*/
func CreateSpacesAliases(tblspace *TblSpacesAliases) error {

	if err := s.Authority.DB.Table("tbl_spaces_aliases").Create(&tblspace).Error; err != nil {

		return err
	}

	return nil
}

// Clone space

func ClonePage(page *TblPage) (*TblPage, error) {

	if err := s.Authority.DB.Table("tbl_page").Create(&page).Error; err != nil {

		return &TblPage{}, err

	}
	return page, nil
}

func ClonePages(pages *TblPageAliases) error {

	if err := s.Authority.DB.Table("tbl_page_aliases").Omit("id").Create(&pages).Error; err != nil {

		return err

	}
	return nil

}

func CloneSpaceInPagesGroup(group *TblPagesGroup) (*TblPagesGroup, error) {

	if err := s.Authority.DB.Table("tbl_pages_group").Create(&group).Error; err != nil {

		return &TblPagesGroup{}, err

	}
	return group, nil

}

func ClonePagesGroup(pagegroup *TblPagesGroupAliases) error {

	if err := s.Authority.DB.Table("tbl_pages_group_aliases").Omit("id").Create(&pagegroup).Error; err != nil {

		return err

	}
	return nil
}

/*Update Space*/
func EditSpace(tblspace *TblSpacesAliases, id int) error {

	if tblspace.ImagePath != "" {
		s.Authority.DB.Table("tbl_spaces_aliases").Where("spaces_id = ?", tblspace.Id).UpdateColumns(map[string]interface{}{"spaces_name": tblspace.SpacesName, "spaces_description": tblspace.SpacesDescription, "spaces_slug": tblspace.SpacesSlug, "image_path": tblspace.ImagePath, "modified_by": tblspace.ModifiedBy, "modified_on": tblspace.ModifiedOn})

	} else {
		s.Authority.DB.Table("tbl_spaces_aliases").Where("spaces_id = ?", tblspace.Id).UpdateColumns(map[string]interface{}{"spaces_name": tblspace.SpacesName, "spaces_description": tblspace.SpacesDescription, "spaces_slug": tblspace.SpacesSlug, "modified_by": tblspace.ModifiedBy, "modified_on": tblspace.ModifiedOn})

	}
	return nil
}

/*Update Space*/
func UpdateSpace(tblspace *TblSpaces, id int) error {

	if tblspace.PageCategoryId != 0 {

		s.Authority.DB.Table("tbl_spaces").Where("id = ?", tblspace.Id).UpdateColumns(map[string]interface{}{"page_category_id": tblspace.PageCategoryId, "modified_by": tblspace.ModifiedBy, "modified_on": tblspace.ModifiedOn})

	} else {

		s.Authority.DB.Table("tbl_spaces").Where("id = ?", tblspace.Id).UpdateColumns(map[string]interface{}{"modified_by": tblspace.ModifiedBy, "modified_on": tblspace.ModifiedOn})

	}
	return nil
}

/*Deleted space*/
func DeleteSpaceAliases(tblspace *TblSpacesAliases, id int) error {

	if err := s.Authority.DB.Table("tbl_spaces_aliases").Where("spaces_id = ?", id).UpdateColumns(map[string]interface{}{"deleted_by": tblspace.DeletedBy, "deleted_on": tblspace.DeletedOn, "is_deleted": tblspace.IsDeleted}).Error; err != nil {

		return err
	}

	return nil
}

/*Deleted space*/
func DeleteSpace(tblspace *TblSpaces, id int) error {

	if err := s.Authority.DB.Table("tbl_spaces").Where("id = ?", id).UpdateColumns(map[string]interface{}{"deleted_by": tblspace.DeletedBy, "deleted_on": tblspace.DeletedOn, "is_deleted": tblspace.IsDeleted}).Error; err != nil {

		return err
	}

	return nil
}

func GetPageCategoryTree(id int) ([]TblPagesCategoriesAliases, error) {
	var categories []TblPagesCategoriesAliases
	err := s.Authority.DB.Raw(`
	WITH RECURSIVE cat_tree AS (
		SELECT id,PAGE_CATEGORY_ID,
		LANGUAGE_ID,
		CATEGORY_NAME,
		CATEGORY_SLUG,
		DESCRIPTION,
		PARENT_ID,
		CREATED_ON,
		MODIFIED_ON,
		IS_DELETED
		FROM tbl_pages_categories_aliases
		WHERE id = ?
		UNION ALL
		SELECT cat.id,CAT.PAGE_CATEGORY_ID,
		CAT.LANGUAGE_ID,
		CAT.CATEGORY_NAME,
		CAT.CATEGORY_SLUG,
		CAT.DESCRIPTION,
		CAT.PARENT_ID,
		CAT.CREATED_ON,
		CAT.MODIFIED_ON,
		CAT.IS_DELETED
		FROM tbl_pages_categories_aliases AS cat
		JOIN cat_tree ON cat.parent_id = cat_tree.id
	)
	SELECT *
	FROM cat_tree WHERE IS_DELETED = 0 order by id desc
	`, id).Scan(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func GetParentPageCategory(pagecategory *TblPagesCategoriesAliases, page_category_id int) (error, TblPagesCategoriesAliases) {

	if err := s.Authority.DB.Table("tbl_pages_categories_aliases").Where("is_deleted = 0 and page_category_id=?", page_category_id).Find(&pagecategory).Error; err != nil {

		return err, TblPagesCategoriesAliases{}
	}

	return nil, *pagecategory
}

func GetChildPageCategories(pagecategory *[]TblPagesCategoriesAliases, parent_id int) (error, []TblPagesCategoriesAliases) {

	if err := s.Authority.DB.Table("tbl_pages_categories_aliases").Where("is_deleted=0 and page_category_id=?", parent_id).Find(&pagecategory).Error; err != nil {

		return err, []TblPagesCategoriesAliases{}
	}

	return nil, *pagecategory
}

// Category based space list
func GetSpacesData(spaces *[]TblSpaces, id int) error {

	if err := s.Authority.DB.Table("tbl_spaces").Where("is_deleted=? and page_category_id = ?", id, 0).Find(&spaces).Error; err != nil {

		return err

	}

	return nil
}

// space data

func GetPageData(page *[]TblPageAliases, id int) error {

	if err := s.Authority.DB.Table("tbl_page_aliases").Select("tbl_page_aliases.language_id,tbl_page_aliases.page_title,tbl_page_aliases.page_slug,tbl_page_aliases.page_description,tbl_page_aliases.published_on,tbl_page_aliases.author,tbl_page_aliases.excerpt,tbl_page_aliases.featured_images,tbl_page_aliases.access,tbl_page_aliases.meta_details,tbl_page_aliases.status,tbl_page_aliases.allow_comments,tbl_page_aliases.created_on,tbl_page_aliases.created_by,tbl_page_aliases.order_index,tbl_page_aliases.page_suborder,tbl_page.id").
		Joins("inner join tbl_page on tbl_page_aliases.page_id = tbl_page.id").Where("tbl_page.spaces_id = ?", id).Find(&page).Error; err != nil {

		return err

	}

	return nil

}

func GetPageGroupData(group *[]TblPagesGroupAliases, id int) error {

	if err := s.Authority.DB.Table("tbl_pages_group_aliases").Joins("inner join tbl_pages_group on tbl_pages_group_aliases.page_group_id = tbl_pages_group.id").Where("tbl_pages_group.is_deleted = ? and tbl_pages_group_aliases.is_deleted = ? and tbl_pages_group.spaces_id = ?", 0, 0, id).Find(&group).Error; err != nil {

		return err

	}

	return nil

}

func GetIdInPage(pageid *TblPagesGroupAliases, spaceid int) error {

	if err := s.Authority.DB.Table("tbl_pages_group_aliases").Joins("inner join tbl_page on tbl_pages_group_aliases.page_group_id = tbl_page.page_group_id").Where("tbl_page.is_deleted = ? and tbl_pages_group_aliases.is_deleted = ? and tbl_page.page_group_id != ? and  tbl_page.parent_id = ? and  tbl_page.spaces_id = ?", 0, 0, 0, 0, spaceid).First(&pageid).Error; err != nil {
		return err

	}

	return nil
}

func GetPageInPage(pageid *[]TblPageAliases, spaceid int) error {

	if err := s.Authority.DB.Table("tbl_page_aliases").Select("tbl_page_aliases.*").Joins("inner join tbl_page on tbl_page_aliases.page_id = tbl_page.id").Where("tbl_page.is_deleted = ? and tbl_page_aliases.is_deleted = ? and tbl_page.page_group_id = ? and  parent_id = ? and  tbl_page.spaces_id = ?", 0, 0, 0, 0, spaceid).Find(&pageid).Error; err != nil {

		return err

	}

	return nil
}

func GetDetailsInPageAli(pagedetails *TblPagesGroupAliases, groupname string, spaceid int) error {

	if err := s.Authority.DB.Table("tbl_pages_group_aliases").Joins("inner join tbl_pages_group on tbl_pages_group_aliases.page_group_id = tbl_pages_group.id").Where("tbl_pages_group.is_deleted = ? and  tbl_pages_group_aliases.is_deleted = ? and  tbl_pages_group_aliases.group_name = ? and tbl_pages_group.spaces_id = ? ", 0, 0, groupname, spaceid).Find(&pagedetails).Error; err != nil {

		return err

	}

	return nil
}

func GetDetailsInPageAlia(pageid *TblPagesGroupAliases, pagegroupid int, spaceid int) error {

	if err := s.Authority.DB.Table("tbl_pages_group_aliases").Joins("inner join tbl_pages_group on tbl_pages_group_aliases.page_group_id = tbl_pages_group.id").Where("tbl_pages_group_aliases.page_group_id = ? and  tbl_pages_group.spaces_id = ?", pagegroupid, spaceid).First(&pageid).Error; err != nil {

		return err

	}

	return nil
}

func GetPageAliasesInPage(data *[]TblPageAliases, spacid int) error {

	if err := s.Authority.DB.Table("tbl_page_aliases").Select("tbl_page_aliases.*,tbl_page.parent_id,tbl_page.page_group_id,tbl_page.spaces_id").Joins("inner join tbl_page on tbl_page_aliases.page_id = tbl_page.id").Where("tbl_page_aliases.is_deleted = ? and tbl_page.is_deleted = ? and tbl_page.page_group_id != ? and  tbl_page.parent_id = ? and  tbl_page.spaces_id = ?", 0, 0, 0, 0, spacid).Find(&data).Error; err != nil {

		return err

	}

	return nil
}

func GetPageAliasesInPageData(result *[]TblPageAliases, spacid int) error {

	if err := s.Authority.DB.Table("tbl_page_aliases").Select("tbl_page_aliases.*,tbl_page.page_group_id,tbl_page.parent_id").Joins("inner join tbl_page on tbl_page_aliases.page_id = tbl_page.id").Where("tbl_page.parent_id != ? and  tbl_page.spaces_id = ?", 0, spacid).Find(&result).Error; err != nil {

		return err

	}

	return nil
}

func LastLoopAliasesInPage(data *TblPageAliases, pagetitle string, spacid int) error {

	if err := s.Authority.DB.Table("tbl_page_aliases").Select("tbl_page_aliases.*").Joins("inner join tbl_page on tbl_page_aliases.page_id = tbl_page.id").Where("tbl_page_aliases.page_title = ? and  tbl_page.spaces_id = ?", pagetitle, spacid).First(&data).Error; err != nil {

		return err

	}

	return nil
}

func AliasesInParentId(data *TblPageAliases, parentid int, spacid int) error {

	if err := s.Authority.DB.Table("tbl_page_aliases").Select("tbl_page_aliases.*,tbl_page.id,tbl_page.parent_id").Joins("inner join tbl_page on tbl_page_aliases.page_id = tbl_page.id").Where("tbl_page_aliases.page_id = ? and  tbl_page.spaces_id = ?", parentid, spacid).First(&data).Error; err != nil {

		return err

	}

	return nil
}

func ParentWithChild(parent *TblPage, id int, spaceid int) error {

	if err := s.Authority.DB.Table("tbl_page").Select("tbl_page.*,tbl_pages_group.id").Joins("inner join tbl_pages_group on tbl_page.page_group_id = tbl_pages_group.id").Where("tbl_pages_group.is_deleted = ? and tbl_page.is_deleted = ? and tbl_pages_group.id = ? and tbl_page.spaces_id = ? ", 0, 0, id, spaceid).First(&parent).Error; err != nil {

		return err

	}

	return nil
}

func PageParentCategoryList(pagecategory *[]TblPagesCategoriesAliases) error {

	if err := s.Authority.DB.Table("tbl_pages_categories_aliases").Where("is_deleted = 0").Find(&pagecategory).Error; err != nil {

		return err
	}
	return nil

}

/*spacename*/
func GetSpaceName(TblSpacesAliases *TblSpacesAliases, spaceid int) error {

	if err := s.Authority.DB.Table("tbl_spaces_aliases").Where("spaces_id=?", spaceid).First(&TblSpacesAliases).Error; err != nil {

		return err
	}

	return nil
}
