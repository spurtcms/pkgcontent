package lms

import (
	"time"

	"github.com/spurtcms/pkgcore/auth"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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

type PageLog struct {
	Username string
	Status   string
	Date     time.Time
}

type PageGroups struct {
	GroupId    int
	NewGroupId int
	Name       string
	OrderIndex int `json:"OrderIndex"`
}

type Pages struct {
	PgId        int
	NewPgId     int
	Name        string
	Content     string `json:"Content"`
	Pgroupid    int
	NewGrpId    int
	OrderIndex  int `json:"OrderIndex"`
	ParentId    int
	CreatedDate time.Time
	LastUpdate  time.Time
	Date        string
	Username    string
	Log         []PageLog
}

type SubPages struct {
	SpgId       int
	NewSpId     int
	Name        string
	Content     string
	ParentId    int
	NewParentId int
	PgroupId    int
	NewPgroupId int
	OrderIndex  int `json:"OrderIndex"`
	CreatedDate time.Time
	LastUpdate  time.Time
	Date        string
	Username    string
	Log         []PageLog
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
	Id               int `gorm:"primaryKey;auto_increment"`
	PageId           int
	LanguageId       int
	PageTitle        string
	PageSlug         string
	PageDescription  string
	PublishedOn      time.Time `gorm:"DEFAULT:NULL"`
	Author           string
	Excerpt          string
	FeaturedImages   string
	Access           string
	MetaDetails      datatypes.JSONType[MetaDetails]
	Status           string
	AllowComments    bool
	CreatedOn        time.Time
	CreatedBy        int
	ModifiedOn       time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy       int       `gorm:"DEFAULT:NULL"`
	DeletedOn        time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy        int       `gorm:"DEFAULT:NULL"`
	IsDeleted        int       `gorm:"DEFAULT:0"`
	OrderIndex       int
	PageSuborder     int
	CreatedDate      string `gorm:"-"`
	ModifiedDate     string `gorm:"-"`
	Username         string `gorm:"<-:false"`
	PageGroupId      int    `gorm:"-:migration;<-:false"`
	ParentId         int    `gorm:"-:migration;<-:false"`
	LastRevisionDate time.Time
	LastRevisionNo   int
}

type TblPageAliasesLog struct {
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
	CreatedDate     string    `gorm:"-"`
	ModifiedDate    string    `gorm:"-"`
	Username        string    `gorm:"-:migration;<-:false"`
	PageGroupId     int       `gorm:"-:migration;<-:false"`
	ParentId        int       `gorm:"-:migration;<-:false"`
}

type PageCreate struct {
	SpaceId       int          //spaceid
	NewPages      []Pages      //pages only
	NewGroup      []PageGroups //groups only
	NewSubPage    []SubPages   //subpages only
	UpdatePages   []Pages      //pages only
	UpdateGroup   []PageGroups //groups only
	UpdateSubPage []SubPages   //subpages only
	DeletePages   []Pages      //delete pages only
	DeleteGroup   []PageGroups //delete groups only
	DeleteSubPage []SubPages   //delete subpages only
	Status        string       //publish,draft
}

type TblMemberNotesHighlight struct {
	Id                      int `gorm:"primaryKey;auto_increment"`
	MemberId                int
	PageId                  int
	NotesHighlightsContent  string
	NotesHighlightsType     string
	HighlightsConfiguration datatypes.JSONMap `gorm:"type:jsonb"`
	CreatedBy               int
	CreatedOn               time.Time
	ModifiedOn              time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy              int       `gorm:"DEFAULT:NULL"`
	DeletedOn               time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy               int       `gorm:"DEFAULT:NULL"`
	IsDeleted               int
}

type HighlightsReq struct {
	Pageid       int
	Content      string
	Start        int
	Offset       int
	SelectPara   string
	ContentColor string
}

func (P PageStrut) CreatePageGroup(tblpagegroup *TblPagesGroup, DB *gorm.DB) (*TblPagesGroup, error) {

	if err := DB.Table("tbl_pages_group").Create(&tblpagegroup).Error; err != nil {

		return &TblPagesGroup{}, err
	}

	return tblpagegroup, nil

}

// create page
func (P PageStrut) Createpage(tblpage *TblPage, DB *gorm.DB) error {

	if err := DB.Table("tbl_page").Create(&tblpage).Error; err != nil {

		return err
	}

	return nil

}

// create PageAliases
func (P PageStrut) CreatepageAliases(tblpageAliases *TblPageAliases, DB *gorm.DB) error {

	if err := DB.Table("tbl_page_aliases").Create(&tblpageAliases).Error; err != nil {

		return err
	}

	return nil

}

/*Create PagegroupAliases */
func (P PageStrut) CreatePageGroupAliases(tblpagegroup *TblPagesGroupAliases, DB *gorm.DB) error {

	if err := DB.Table("tbl_pages_group_aliases").Create(&tblpagegroup).Error; err != nil {

		return err
	}

	return nil
}

/*Update pagegroup*/
func (P PageStrut) UpdatePageGroup(tblpagegroup *TblPagesGroup, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_pages_group").Where("id = ?", id).UpdateColumns(map[string]interface{}{"modified_on": tblpagegroup.ModifiedOn, "modified_by": tblpagegroup.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

/*pdate pagegroupAliases */
func (P PageStrut) UpdatePageGroupAliases(tblpagegroup *TblPagesGroupAliases, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_pages_group_aliases").Where("page_group_id = ?", id).UpdateColumns(map[string]interface{}{"group_name": tblpagegroup.GroupName, "group_slug": tblpagegroup.GroupSlug, "group_description": tblpagegroup.GroupDescription, "language_id": tblpagegroup.LanguageId, "modified_on": tblpagegroup.ModifiedOn, "modified_by": tblpagegroup.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

/*update page*/
func (P PageStrut) UpdatePage(tblpage *TblPage, pageid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_page").Where("id=?", pageid).UpdateColumns(map[string]interface{}{"page_group_id": tblpage.PageGroupId, "parent_id": tblpage.ParentId}).Error; err != nil {

		return err
	}

	return nil
}

/*update pagealiases*/
func (P PageStrut) UpdatePageAliase(tblpageali *TblPageAliases, pageid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_page_aliases").Where("page_id=?", pageid).UpdateColumns(map[string]interface{}{
		"page_title": tblpageali.PageTitle, "page_slug": tblpageali.PageSlug, "modified_on": tblpageali.ModifiedOn,
		"modified_by": tblpageali.ModifiedBy, "page_description": tblpageali.PageDescription, "order_index": tblpageali.OrderIndex, "status": tblpageali.Status}).Error; err != nil {
		return err
	}

	return nil
}

/*update pagealiases*/
func (P PageStrut) UpdatePageAliasePublishStatus(pageid []int, userid int, DB *gorm.DB) error {

	Formatdate, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	if err := DB.Table("tbl_page_aliases").Where("page_id in (?)", pageid).UpdateColumns(map[string]interface{}{"status": "publish", "modified_on": Formatdate, "modified_by": userid}).Error; err != nil {
		return err
	}

	return nil
}

func (P PageStrut) SelectGroup(tblgroup *[]TblPagesGroup, id int, grpid []int, DB *gorm.DB) error {

	query := DB.Table("tbl_pages_group").Where("spaces_id = ? and is_deleted=0", id)

	// if len(grpid) != 0 {

	// 	query = query.Where("id in (?)", grpid)

	// }

	query.Find(&tblgroup)

	if err := query.Error; err != nil {

		return err

	}

	return nil
}

func (P PageStrut) SelectPage(tblpage *[]TblPage, id int, pgid []int, DB *gorm.DB) error {

	query := DB.Table("tbl_page").Where("spaces_id = ? and is_deleted =0 ", id)

	// if len(pgid) != 0 {

	// 	query = query.Where("id in (?)", pgid)

	// }

	query.Find(&tblpage)

	if err := query.Error; err != nil {

		return err

	}

	return nil
}
func (P PageStrut) PageGroup(tblpagegroup *TblPagesGroupAliases, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_pages_group_aliases").Where("is_deleted = 0 and page_group_id = ?", id).First(&tblpagegroup).Error; err != nil {

		return err

	}

	return nil
}

func (P PageStrut) PageAliases(tblpagegroup *TblPageAliases, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_page_aliases").Select("tbl_page_aliases.*,tbl_page.page_group_id,tbl_users.username").Joins("inner join tbl_page on tbl_page.id = tbl_page_aliases.page_id").Joins("inner join tbl_users on tbl_users.id = tbl_page_aliases.created_by").Where("page_id = ? and tbl_page.is_deleted=0 and tbl_page_aliases.is_deleted=0", id).Find(&tblpagegroup).Error; err != nil {

		return err

	}

	return nil
}

/* Delete group */
func (P PageStrut) DeletePageGroup(tblpagegroup *TblPagesGroup, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_pages_group").Where("id=?", id).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": tblpagegroup.DeletedOn, "deleted_by": tblpagegroup.DeletedBy}).Error; err != nil {

		return err

	}
	return nil
}

/* Delete Groupaliases */
func (P PageStrut) DeletePageGroupAliases(tblpagegroup *TblPagesGroupAliases, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_pages_group_aliases").Where("page_group_id=?", id).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": tblpagegroup.DeletedOn, "deleted_by": tblpagegroup.DeletedBy}).Error; err != nil {

		return err

	}
	return nil
}

// Delete Page
func (P PageStrut) DeletePageAliases(tblpageAliases *TblPageAliases, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_page_aliases").Where("page_id=?", id).UpdateColumns(map[string]interface{}{"deleted_on": tblpageAliases.DeletedOn, "deleted_by": tblpageAliases.DeletedBy}).Error; err != nil {

		return err

	}
	return nil
}

/*Check if groupexist*/
func (P PageStrut) CheckGroupExists(tblgroup *TblPagesGroup, id int, spaceid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_pages_group").Where("id=? and spaces_id=?", id, spaceid).First(&tblgroup).Error; err != nil {

		return err
	}

	return nil
}

/*Check if page exists*/
func (P PageStrut) CheckPageExists(tblpage *TblPage, pageid int, spaceid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_page").Where("id=? and spaces_id=?", pageid, spaceid).First(&tblpage).Error; err != nil {

		return err
	}

	return nil
}

/*Delete PageAliases*/
func (P PageStrut) DeletePageAliase(tblpage *TblPageAliases, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_page_aliases").Where("page_id=?", id).UpdateColumns(map[string]interface{}{"deleted_on": tblpage.DeletedOn, "deleted_by": tblpage.DeletedBy, "is_deleted": 1}).Error; err != nil {

		return err
	}

	return nil

}

/*Delete PageAliases*/
func (P PageStrut) DeletePage(tblpage *TblPage, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_page").Where("id=?", id).UpdateColumns(map[string]interface{}{"deleted_on": tblpage.DeletedOn, "deleted_by": tblpage.DeletedBy, "is_deleted": 1}).Error; err != nil {

		return err
	}

	return nil

}

/*PageGroup*/
func (P PageStrut) GetPageGroupByName(TblPagesGroupAliases *TblPagesGroupAliases, spaceid int, name string, DB *gorm.DB) error {

	if err := DB.Table("tbl_pages_group_aliases").Joins("inner join tbl_pages_group on tbl_pages_group.id=tbl_pages_group_aliases.page_group_id").Where("LOWER(TRIM(group_name))=LOWER(TRIM(?)) and tbl_pages_group.spaces_id=? and tbl_pages_group_aliases.is_deleted=0", name, spaceid).Last(&TblPagesGroupAliases).Error; err != nil {

		return err
	}

	return nil
}

/*GetPage*/
func (P PageStrut) GetPageDataByName(TblPageAliases *TblPageAliases, spaceid int, name string, DB *gorm.DB) error {

	if err := DB.Table("tbl_page_aliases").Select("tbl_page_aliases.*").Joins("inner join tbl_page on tbl_page.id=tbl_page_aliases.page_id").Where("page_title=? and tbl_page.spaces_id=? and tbl_page_aliases.is_deleted=0", name, spaceid).Last(&TblPageAliases).Error; err != nil {

		return err
	}

	return nil
}

/*CreatePage*/
func (P PageStrut) CreatePage(tblpage *TblPage, DB *gorm.DB) (*TblPage, error) {

	if err := DB.Table("tbl_page").Create(&tblpage).Error; err != nil {

		return &TblPage{}, err
	}
	return tblpage, nil

}

/*Create page log*/
func (P PageStrut) PageAliasesLog(tblpagelog *TblPageAliasesLog, DB *gorm.DB) error {

	if err := DB.Table("tbl_page_aliases_log").Create(&tblpagelog).Error; err != nil {

		return err
	}

	return nil
}

/*Get page log*/
func (P PageStrut) GetPageLogDetails(tblpagelog *[]TblPageAliasesLog, spaceid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_page_aliases_log").Select("tbl_page_aliases_log.created_by,tbl_page_aliases_log.created_on,tbl_page_aliases_log.status,tbl_users.username,max(TBL_PAGE_ALIASES_LOG.modified_by) as modified_by,max(TBL_PAGE_ALIASES_LOG.modified_on) as modified_on").Joins("inner join tbl_page on tbl_page.id = tbl_page_aliases_log.page_id").Joins("inner join tbl_users on tbl_users.id = tbl_page_aliases_log.created_by").Where("tbl_page.spaces_id=?", spaceid).Group("tbl_page_aliases_log.created_by,tbl_page_aliases_log.created_on,tbl_page_aliases_log.status,tbl_users.username").Order("tbl_page_aliases_log.created_on desc").Find(&tblpagelog).Error; err != nil {

		return err
	}

	return nil
}

/*Get page log*/
func (P PageStrut) GetPageLogDetailsByPageId(tblpagelog *[]TblPageAliasesLog, spaceid int, pageid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_page_aliases_log").Select("tbl_page_aliases_log.created_by,tbl_page_aliases_log.created_on,tbl_page_aliases_log.status,tbl_users.username,max(TBL_PAGE_ALIASES_LOG.modified_by) as modified_by,max(TBL_PAGE_ALIASES_LOG.modified_on) as modified_on").Joins("inner join tbl_page on tbl_page.id = tbl_page_aliases_log.page_id").Joins("inner join tbl_users on tbl_users.id = tbl_page_aliases_log.created_by").Where("tbl_page.spaces_id=? and page_id = ? ", spaceid, pageid).Group("tbl_page_aliases_log.created_by,tbl_page_aliases_log.created_on,tbl_page_aliases_log.status,tbl_users.username").Order("tbl_page_aliases_log.created_on desc").Find(&tblpagelog).Error; err != nil {

		return err
	}

	return nil
}

/*Get Content*/
func (P PageStrut) GetContentByPageId(tblpage *TblPageAliases, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_page_aliases").Where("page_id=?", id).First(&tblpage).Error; err != nil {

		return err
	}

	return nil
}

/*GET NOTES*/
func (p PageStrut) GetNotes(notes *[]TblMemberNotesHighlight, memberid int, pageid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_member_notes_highlights").Where("member_id=? and page_id=? and notes_highlights_type='notes' and is_deleted=0", memberid, pageid).Find(&notes).Error; err != nil {
		return err
	}

	return nil
}

/*GET NOTES*/
func (p PageStrut) GetHighlights(notes *[]TblMemberNotesHighlight, memberid int, pageid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_member_notes_highlights").Where("member_id=? and page_id=? and notes_highlights_type='highlights' and is_deleted=0", memberid, pageid).Find(&notes).Error; err != nil {
		return err
	}

	return nil
}

/*Insert Notes*/
func (p PageStrut) UpdateNotesHighlights(notes *TblMemberNotesHighlight, contype string, DB *gorm.DB) error {

	if contype == "notes" {

		if err := DB.Model(TblMemberNotesHighlight{}).Create(notes).Error; err != nil {

			return err

		}

	} else if contype == "highlights" {

		if err := DB.Model(TblMemberNotesHighlight{}).Create(notes).Error; err != nil {

			return err

		}

	}

	return nil

}

/*Remove Highligts*/
func (P PageStrut) RemoveHighlights(high *TblMemberNotesHighlight, DB *gorm.DB) error {

	if err := DB.Model(TblMemberNotesHighlight{}).Where("id=?", high.Id).UpdateColumns(map[string]interface{}{"is_deleted": high.IsDeleted, "deleted_by": high.DeletedBy, "deleted_on": high.DeletedOn}).Error; err != nil {

		return err
	}

	return nil
}

func (P PageStrut) MemberRestrictActive(Mod *auth.TblModule, DB *gorm.DB) error {

	if err := DB.Model(auth.TblModule{}).Where("(LOWER(TRIM(module_name)) ILIKE LOWER(TRIM('Member Restrict')))").First(&Mod).Error; err != nil {

		return err

	}

	return nil

}
