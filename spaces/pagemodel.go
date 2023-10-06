package spaces

import (
	"time"

	"gorm.io/datatypes"
)

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

type PageGroups struct {
	GroupId    int
	NewGroupId int
	Name       string
	OrderIndex int `json:"OrderIndex"`
}

type Pages struct {
	PgId       int
	NewPgId    int
	Name       string
	Content    string `json:"Content"`
	Pgroupid   int
	OrderIndex int `json:"OrderIndex"`
	ParentId   int
}

type SubPages struct {
	SpgId       int
	NewSpId     int
	Name        string
	Content     string
	ParentId    int
	PgroupId    int
	NewPgroupId int
	OrderIndex  int `json:"OrderIndex"`
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

var s Space

func CreatePageGroup(tblpagegroup *TblPagesGroup) (*TblPagesGroup, error) {

	if err := s.Authority.DB.Table("tbl_pages_group").Create(&tblpagegroup).Error; err != nil {

		return &TblPagesGroup{}, err
	}

	return tblpagegroup, nil

}

// create page
func Createpage(tblpage *TblPage) error {

	if err := s.Authority.DB.Table("tbl_page").Create(&tblpage).Error; err != nil {

		return err
	}

	return nil

}

// create PageAliases
func CreatepageAliases(tblpageAliases *TblPageAliases) error {

	if err := s.Authority.DB.Debug().Table("tbl_page_aliases").Create(&tblpageAliases).Error; err != nil {

		return err
	}

	return nil

}

/*Create PagegroupAliases */
func CreatePageGroupAliases(tblpagegroup *TblPagesGroupAliases) error {

	if err := s.Authority.DB.Debug().Table("tbl_pages_group_aliases").Create(&tblpagegroup).Error; err != nil {

		return err
	}

	return nil
}

/*Update pagegroup*/
func UpdatePageGroup(tblpagegroup *TblPagesGroup, id int) error {

	if err := s.Authority.DB.Table("tbl_pages_group").Where("id = ?", id).UpdateColumns(map[string]interface{}{"modified_on": tblpagegroup.ModifiedOn, "modified_by": tblpagegroup.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

/*pdate pagegroupAliases */
func UpdatePageGroupAliases(tblpagegroup *TblPagesGroupAliases, id int) error {

	if err := s.Authority.DB.Debug().Table("tbl_pages_group_aliases").Where("page_group_id = ?", id).UpdateColumns(map[string]interface{}{"group_name": tblpagegroup.GroupName, "group_slug": tblpagegroup.GroupSlug, "group_description": tblpagegroup.GroupDescription, "language_id": tblpagegroup.LanguageId, "modified_on": tblpagegroup.ModifiedOn, "modified_by": tblpagegroup.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

/*update page*/
func UpdatePage(tblpage *TblPage, pageid int) error {

	if err := s.Authority.DB.Table("tbl_page").Where("id=?", pageid).UpdateColumns(map[string]interface{}{"page_group_id": tblpage.PageGroupId, "parent_id": tblpage.ParentId}).Error; err != nil {

		return err
	}

	return nil
}

/*update pagealiases*/
func UpdatePageAliase(tblpageali *TblPageAliases, pageid int) error {

	if err := s.Authority.DB.Table("tbl_page_aliases").Where("page_id=?", pageid).UpdateColumns(map[string]interface{}{
		"page_title": tblpageali.PageTitle, "page_slug": tblpageali.PageSlug, "modified_on": tblpageali.ModifiedOn,
		"modified_by": tblpageali.ModifiedBy, "page_description": tblpageali.PageDescription}).Error; err != nil {
		return err
	}

	return nil
}

func SelectGroup(tblgroup *[]TblPagesGroup, id int) error {

	if err := s.Authority.DB.Table("tbl_pages_group").Where("spaces_id = ? and is_deleted=0", id).Find(&tblgroup).Error; err != nil {

		return err

	}

	return nil
}

func SelectPage(tblpage *[]TblPage, id int) error {

	if err := s.Authority.DB.Table("tbl_page").Where("spaces_id = ? and is_deleted =0 ", id).Find(&tblpage).Error; err != nil {

		return err

	}

	return nil
}
func PageGroup(tblpagegroup *TblPagesGroupAliases, id int) error {

	if err := s.Authority.DB.Table("tbl_pages_group_aliases").Where("is_deleted = 0 and page_group_id = ?", id).First(&tblpagegroup).Error; err != nil {

		return err

	}

	return nil
}
func PageAliases(tblpagegroup *TblPageAliases, id int) error {

	if err := s.Authority.DB.Table("tbl_page_aliases").Where("page_id = ? and is_deleted=0", id).Find(&tblpagegroup).Error; err != nil {

		return err

	}

	return nil
}

/* Delete group */
func DeletePageGroup(tblpagegroup *TblPagesGroup, id int) error {

	if err := s.Authority.DB.Table("tbl_pages_group").Where("id=?", id).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": tblpagegroup.DeletedOn, "deleted_by": tblpagegroup.DeletedBy}).Error; err != nil {

		return err

	}
	return nil
}

/* Delete Groupaliases */
func DeletePageGroupAliases(tblpagegroup *TblPagesGroupAliases, id int) error {

	if err := s.Authority.DB.Table("tbl_pages_group_aliases").Where("page_group_id=?", id).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": tblpagegroup.DeletedOn, "deleted_by": tblpagegroup.DeletedBy}).Error; err != nil {

		return err

	}
	return nil
}

// Delete Page
func DeletePageAliases(tblpageAliases *TblPageAliases, id int) error {

	if err := s.Authority.DB.Table("tbl_page_aliases").Where("page_id=?", id).UpdateColumns(map[string]interface{}{"deleted_on": tblpageAliases.DeletedOn, "deleted_by": tblpageAliases.DeletedBy}).Error; err != nil {

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

/*Check if groupexist*/
func CheckGroupExists(tblgroup *TblPagesGroup, id int, spaceid int) error {

	if err := s.Authority.DB.Table("tbl_pages_group").Where("id=? and spaces_id=?", id, spaceid).First(&tblgroup).Error; err != nil {

		return err
	}

	return nil
}

/*Check if page exists*/
func CheckPageExists(tblpage *TblPage, pageid int, spaceid int) error {

	if err := s.Authority.DB.Table("tbl_page").Where("id=? and spaces_id=?", pageid, spaceid).First(&tblpage).Error; err != nil {

		return err
	}

	return nil
}

/*Delete PageAliases*/
func DeletePageAliase(tblpage *TblPageAliases, id int) error {

	if err := s.Authority.DB.Table("tbl_page_aliases").Where("page_id=?", id).UpdateColumns(map[string]interface{}{"deleted_on": tblpage.DeletedOn, "deleted_by": tblpage.DeletedBy, "is_deleted": 1}).Error; err != nil {

		return err
	}

	return nil

}

/*Delete PageAliases*/
func DeletePage(tblpage *TblPage, id int) error {

	if err := s.Authority.DB.Table("tbl_page").Where("id=?", id).UpdateColumns(map[string]interface{}{"deleted_on": tblpage.DeletedOn, "deleted_by": tblpage.DeletedBy, "is_deleted": 1}).Error; err != nil {

		return err
	}

	return nil

}

/*PageGroup*/
func GetPageGroupByName(TblPagesGroupAliases *TblPagesGroupAliases, spaceid int, name string) error {

	if err := s.Authority.DB.Table("tbl_pages_group_aliases").Joins("inner join tbl_pages_group on tbl_pages_group.id=tbl_pages_group_aliases.page_group_id").Where("group_name=? and tbl_pages_group.spaces_id=? and tbl_pages_group_aliases.is_deleted=0", name, spaceid).First(&TblPagesGroupAliases).Error; err != nil {

		return err
	}

	return nil
}

/*GetPage*/
func GetPageDataByName(TblPageAliases *TblPageAliases, spaceid int, name string) error {

	if err := s.Authority.DB.Table("tbl_page_aliases").Select("tbl_page_aliases.*").Joins("inner join tbl_page on tbl_page.id=tbl_page_aliases.page_id").Where("page_title=? and tbl_page.spaces_id=?", name, spaceid).First(&TblPageAliases).Error; err != nil {

		return err
	}

	return nil
}

/*CreatePage*/
func CreatePage(tblpage *TblPage) (*TblPage, error) {

	if err := s.Authority.DB.Table("tbl_page").Create(&tblpage).Error; err != nil {

		return &TblPage{}, err
	}
	return tblpage, nil

}
