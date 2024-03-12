package channels

import (
	"strconv"
	"time"

	"github.com/spurtcms/pkgcontent/categories"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ChannelCreate struct {
	ChannelName        string
	ChannelDescription string
	Sections           []Section
	FieldValues        []Fiedlvalue
	CategoryIds        []string
}

type ChannelUpdate struct {
	ChannelName        string
	ChannelDescription string
	Sections           []Section
	FieldValues        []Fiedlvalue
	Deletesections     []Section
	DeleteFields       []Fiedlvalue
	DeleteOptionsvalue []OptionValues
	CategoryIds        []string
}

type Filter struct {
	Keyword string
}

type EntriesFilter struct {
	Keyword     string
	Title       string
	ChannelName string
	Status      string
	UserName    string
	CategoryId  string
}

type TblFieldType struct {
	Id         int
	TypeName   string
	TypeSlug   string
	IsActive   int
	IsDeleted  int
	CreatedBy  int
	CreatedOn  time.Time
	ModifiedBy int
	ModifiedOn time.Time
}

type TblFieldGroup struct {
	Id         int
	GroupName  string
	CreatedOn  time.Time
	CreatedBy  int
	ModifiedOn time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy int       `gorm:"DEFAULT:NULL"`
	IsDeleted  int       `gorm:"DEFAULT:0"`
	DeletedOn  time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy  int       `gorm:"DEFAULT:NULL"`
}

type TblGroupField struct {
	Id             int
	ChannelId      int
	FieldId        int
	FieldName      string `gorm:"<-:false"`
	FieldTypeId    int    `gorm:"<-:false"`
	MandatoryField int    `gorm:"<-:false"`
	OptionExist    int    `gorm:"<-:false"`
	FoptionId      int    `gorm:"<-:false"`
	OptionName     string `gorm:"<-:false"`
	OptionValue    string `gorm:"<-:false"`
}

type TblField struct {
	Id               int
	FieldName        string
	FieldDesc        string
	FieldTypeId      int
	MandatoryField   int
	OptionExist      int
	InitialValue     string
	Placeholder      string
	CreatedOn        time.Time
	CreatedBy        int
	ModifiedOn       time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy       int       `gorm:"DEFAULT:NULL"`
	IsDeleted        int       `gorm:"DEFAULT:0"`
	DeletedOn        time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy        int       `gorm:"DEFAULT:NULL"`
	OrderIndex       int
	ImagePath        string
	TblFieldOption   []TblFieldOption `gorm:"<-:false; foreignKey:FieldId"`
	DatetimeFormat   string
	TimeFormat       string
	Url              string
	Values           string         `gorm:"-"`
	CheckBoxValue    []FieldValueId `gorm:"-"`
	SectionParentId  int
	FieldTypeName    string `gorm:"column:type_name"`
	CharacterAllowed int
	FieldOptions     []TblFieldOption     `gorm:"-"`
	FieldValue       TblChannelEntryField `gorm:"-"`
}

type TblFieldOption struct {
	Id          int
	OptionName  string
	OptionValue string
	FieldId     int
	CreatedOn   time.Time
	CreatedBy   int
	ModifiedOn  time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy  int       `gorm:"DEFAULT:NULL"`
	IsDeleted   int       `gorm:"DEFAULT:0"`
	DeletedOn   time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy   int       `gorm:"DEFAULT:NULL"`
	Idstring    string    `gorm:"-"`
}

type TblChannel struct {
	Id                 int
	ChannelName        string
	ChannelDescription string
	SlugName           string
	FieldGroupId       int
	IsActive           int
	IsDeleted          int
	CreatedOn          time.Time
	CreatedBy          int
	ModifiedOn         time.Time           `gorm:"DEFAULT:NULL"`
	ModifiedBy         int                 `gorm:"DEFAULT:NULL"`
	DateString         string              `gorm:"-"`
	EntriesCount       int                 `gorm:"-"`
	ChannelEntries     []TblChannelEntries `gorm:"-"`
	ProfileImagePath   string              `gorm:"<-:false"`
	Username           string              `gorm:"<-:false"`
}

type TblChannelCategory struct {
	Id         int
	ChannelId  int
	CategoryId string
	CreatedAt  int
	CreatedOn  time.Time
}

type FieldValueId struct {
	Id     int
	CValue int
}

type Section struct {
	SectionId     int    `json:"SectionId"`
	SectionNewId  int    `json:"SectionNewId"`
	SectionName   string `json:"SectionName"`
	MasterFieldId int    `json:"MasterFieldId"`
	OrderIndex    int    `json:"OrderIndex"`
}

type Fiedlvalue struct {
	MasterFieldId    int            `json:"MasterFieldId"`
	FieldId          int            `json:"FieldId"`
	NewFieldId       int            `json:"NewFieldId"`
	SectionId        int            `json:"SectionId"`
	SectionNewId     int            `json:"SectionNewId"`
	FieldName        string         `json:"FieldName"`
	DateFormat       string         `json:"DateFormat"`
	TimeFormat       string         `json:"TimeFormat"`
	OptionValue      []OptionValues `json:"OptionValue"`
	CharacterAllowed int            `json:"CharacterAllowed"`
	IconPath         string         `json:"IconPath"`
	Url              string         `json:"Url"`
	OrderIndex       int            `json:"OrderIndex"`
	Mandatory        int            `json:"Mandatory"`
}

type OptionValues struct {
	Id         int    `json:"Id"`
	NewId      int    `json:"NewId"`
	FieldId    int    `json:"FieldId"`
	NewFieldId int    `json:"NewFieldId"`
	Value      string `json:"Value"`
}

type TblChannelEntries struct {
	Id                   int
	Title                string `form:"title" binding:"required"`
	Slug                 string `form:"slug" binding:"required"`
	Description          string
	UserId               int
	ChannelId            int
	Status               int //0-draft 1-publish 2-unpublish
	IsActive             int
	IsDeleted            int       `gorm:"DEFAULT:0"`
	DeletedBy            int       `gorm:"DEFAULT:NULL"`
	DeletedOn            time.Time `gorm:"DEFAULT:NULL"`
	CreatedOn            time.Time
	CreatedBy            int
	ModifiedBy           int       `gorm:"DEFAULT:NULL"`
	ModifiedOn           time.Time `gorm:"DEFAULT:NULL"`
	CoverImage           string
	ThumbnailImage       string
	MetaTitle            string `form:"metatitle" binding:"required"`
	MetaDescription      string `form:"metadesc" binding:"required"`
	Keyword              string `form:"keywords" binding:"required"`
	CategoriesId         string
	RelatedArticles      string
	CreatedDate          string                     `gorm:"-"`
	ModifiedDate         string                     `gorm:"-"`
	Username             string                     `gorm:"<-:false"`
	TblChannelEntryField []TblChannelEntryField     `gorm:"<-:false; foreignKey:ChannelEntryId"`
	Category             []categories.TblCategory   `gorm:"<-:false; foreignKey:Id"`
	CategoryGroup        string                     `gorm:"-:migration;<-:false"`
	ChannelName          string                     `gorm:"-:migration;<-:false"`
	Cno                  string                     `gorm:"<-:false"`
	ProfileImagePath     string                     `gorm:"<-:false"`
	EntryStatus          string                     `gorm:"-"`
	Categories           [][]categories.TblCategory `gorm:"-"`
	AdditionalData       string                     `gorm:"-"`
	AuthorDetail         Author                     `gorm:"-"`
	Sections             []TblField                 `gorm:"-"`
	Fields               []TblField                 `gorm:"-"`
	MemberProfiles       []TblMemberProfiles        `gorm:"-"`
}

type TblMemberProfiles struct {
	Id              int               `json:"memberId,omitempty" gorm:"column:id"`
	ProfileName     string            `json:"profileName,omitempty"`
	ProfileSlug     string            `json:"profileSlug,omitempty"`
	ProfilePage     string            `json:"profilePage,omitempty"`
	MemberDetails   datatypes.JSONMap `json:"memberDetails,omitempty" gorm:"column:member_details;type:jsonb"`
	CompanyName     string            `json:"companyName,omitempty"`
	CompanyLocation string            `json:"companyLocation,omitempty"`
	CompanyLogo     string            `json:"companyLogo,omitempty"`
	About           string            `json:"about,omitempty"`
	SeoTitle        string            `json:"seoTitle,omitempty"`
	SeoDescription  string            `json:"seoDescription,omitempty"`
	SeoKeyword      string            `json:"seoKeyword,omitempty"`
}

type Author struct {
	AuthorID         int       `json:"AuthorId" gorm:"column:id"`
	FirstName        string    `json:"FirstName"`
	LastName         string    `json:"LastName"`
	Email            string    `json:"Email"`
	MobileNo         *string   `json:"MobileNo,omitempty"`
	IsActive         *int      `json:"IsActive,omitempty"`
	ProfileImage     *string   `json:"ProfileImage,omitempty"`
	ProfileImagePath *string   `json:"ProfileImagePath,omitempty"`
	CreatedOn        time.Time `json:"CreatedOn"`
	CreatedBy        int       `json:"CreatedBy"`
}

type TblChannelEntryField struct {
	Id             int
	FieldName      string
	FieldValue     string
	ChannelEntryId int
	FieldId        int
	CreatedOn      time.Time
	CreatedBy      int
	ModifiedOn     time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy     int       `gorm:"DEFAULT:NULL"`
	FieldTypeId    int       `gorm:"<-:false"`
	DeletedBy      int       `gorm:"DEFAULT:NULL"`
	DeletedOn      time.Time `gorm:"DEFAULT:NULL"`
}

type SEODetails struct {
	MetaTitle       string
	MetaDescription string
	MetaKeywords    string
	MetaSlug        string
}

type AdditionalFields struct {
	Id            int
	FieldName     string
	FieldValue    string
	FieldId       int
	MultipleValue []string
}

type EntriesRequired struct {
	Title            string
	Content          string
	CoverImage       string
	AdditionalFields []AdditionalFields
	SEODetails       SEODetails
	CategoryIds      string
	ChannelName      string
	Status           int
	ChannelId        int
}

type RecentActivities struct {
	Contenttype string
	Title       string
	User        string
	Imagepath   string
	Createdon   time.Time
	Active      string
	Channelname string
}

/*Get all master fields*/
func (Ch ChannelStruct) GetAllField(channel *[]TblFieldType, DB *gorm.DB) error {

	if err := DB.Table("tbl_field_types").Where("is_deleted=0").Find(&channel).Error; err != nil {

		return err
	}
	return nil
}

/*Create field grup*/
func (Ch ChannelStruct) CreateFieldGroup(fidgroup *TblFieldGroup, DB *gorm.DB) (*TblFieldGroup, error) {

	if err := DB.Table("tbl_field_groups").Create(&fidgroup).Error; err != nil {

		return &TblFieldGroup{}, err

	}

	return fidgroup, nil
}

/*Craete channel */
func (Ch ChannelStruct) CreateChannel(chn *TblChannel, DB *gorm.DB) (*TblChannel, error) {

	if err := DB.Table("tbl_channels").Create(&chn).Error; err != nil {

		return &TblChannel{}, err

	}

	return chn, nil

}

/*create field*/
func (Ch ChannelStruct) CreateFields(flds *TblField, DB *gorm.DB) (*TblField, error) {

	if err := DB.Table("tbl_fields").Create(&flds).Error; err != nil {

		return &TblField{}, err
	}

	return flds, nil
}

/*create option value*/
func (Ch ChannelStruct) CreateFieldOption(optval *TblFieldOption, DB *gorm.DB) error {

	if err := DB.Table("tbl_field_options").Create(&optval).Error; err != nil {

		return err
	}

	return nil
}

func (Ch ChannelStruct) CreateGroupField(grpfield *TblGroupField, DB *gorm.DB) error {

	if err := DB.Table("tbl_group_fields").Create(&grpfield).Error; err != nil {

		return err
	}

	return nil

}

/*Isactive channel*/
func (Ch ChannelStruct) ChannelIsActive(tblch *TblChannel, id, val int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channels").Where("id=?", id).UpdateColumns(map[string]interface{}{"is_active": val, "modified_on": tblch.ModifiedOn, "modified_by": tblch.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

/*channel list*/
func (Ch ChannelStruct) Channellist(chn *[]TblChannel, limit, offset int, filter Filter, activestatus bool, roleid int, flg bool, DB *gorm.DB) (chcount int64, error error) {

	query := DB.Table("tbl_channels").Where("is_deleted = 0").Order("id desc")

	if roleid != 1 && flg {

		query = query.Where("channel_name in (select display_name from tbl_module_permissions inner join tbl_modules on tbl_modules.id = tbl_module_permissions.module_id inner join tbl_role_permissions on tbl_role_permissions.permission_id = tbl_module_permissions.id where role_id =(?) and tbl_modules.module_name='Entries' )", roleid)

	}

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(channel_name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%")
	}

	if activestatus {

		query = query.Where("is_active=1")

	}

	if limit != 0 {

		query.Limit(limit).Offset(offset).Order("id asc").Find(&chn)

	} else {

		query.Find(&chn).Count(&chcount)

		return chcount, nil
	}

	return 0, nil
}

/*Delete Channel*/
func (Ch ChannelStruct) DeleteChannelById(id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channels").Where("id=?", id).UpdateColumns(map[string]interface{}{"is_deleted": 1}).Error; err != nil {

		return err
	}

	return nil
}

/*Delete Channel*/
func (Ch ChannelStruct) DeleteFieldGroupById(tblfieldgrp *TblFieldGroup, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_field_groups").Where("id=?", id).UpdateColumns(map[string]interface{}{"is_deleted": tblfieldgrp.IsDeleted, "deleted_by": tblfieldgrp.DeletedBy, "deleted_on": tblfieldgrp.DeletedOn}).Error; err != nil {

		return err
	}

	return nil
}

/*Delete Field By Id*/
func (Ch ChannelStruct) DeleteFieldById(field *TblField, id []int, DB *gorm.DB) error {

	if err := DB.Table("tbl_fields").Where("id in(?) ", id).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_by": field.DeletedBy, "deleted_on": field.DeletedOn}).Error; err != nil {

		return err
	}

	return nil
}

/*Delete FieldOption By fieldid*/
func (Ch ChannelStruct) DeleteFieldOptionById(fieldopt *TblFieldOption, id []string, fid int, DB *gorm.DB) error {

	if len(id) > 0 {

		if err := DB.Table("tbl_field_options").Where("option_name not in (?) and field_id=?", id, fid).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_by": fieldopt.DeletedBy, "deleted_on": fieldopt.DeletedOn}).Error; err != nil {

			return err
		}

	} else {

		if err := DB.Table("tbl_field_options").Where("field_id=?", fid).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_by": fieldopt.DeletedBy, "deleted_on": fieldopt.DeletedOn}).Error; err != nil {

			return err
		}
	}

	return nil

}

/*Get Channel*/
func (Ch ChannelStruct) GetChannelById(ch *TblChannel, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channels").Where("id=?", id).First(&ch).Error; err != nil {

		return err
	}

	return nil
}

/*Get Channel*/
func (Ch ChannelStruct) GetChannelByChannelName(ch *TblChannel, name string, DB *gorm.DB) error {

	if err := DB.Table("tbl_channels").Where("channel_name=? and is_deleted=0", name).First(&ch).Error; err != nil {

		return err
	}

	return nil
}

/*Get FieldGroupById*/
func (Ch ChannelStruct) GetFieldGroupById(groupfield *[]TblGroupField, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_group_fields").Select("tbl_fields.field_name,tbl_field_options.option_name,tbl_field_options.option_value").Where("field_group_id=?", id).Joins("INNER JOIN TBL_FIELDS ON TBL_GROUP_FIELDS.FIELD_ID = TBL_FIELDS.ID").Joins("LEFT JOIN TBL_FIELD_OPTIONS ON TBL_FIELDS.ID = TBL_FIELD_OPTIONS.FIELD_ID").Find(&groupfield).Error; err != nil {

		return err
	}

	return nil
}

/*Getfieldid using fieldgroupid*/
func (Ch ChannelStruct) GetFieldIdByGroupId(grpfield *[]TblGroupField, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_group_fields").Where("channel_id=?", id).Find(&grpfield).Error; err != nil {

		return err
	}

	return nil
}

/*Get optionvalue*/
func (Ch ChannelStruct) GetFieldAndOptionValue(fld *[]TblField, id []int, DB *gorm.DB) error {

	if err := DB.Table("tbl_fields").Where("id in (?) and is_deleted != 1", id).Preload("TblFieldOption", func(db *gorm.DB) *gorm.DB {
		return DB.Where("is_deleted!=1")
	}).Order("order_index asc").Find(&fld).Error; err != nil {

		return err
	}

	return nil
}

/*Get Field Value*/
func (Ch ChannelStruct) GetFieldvalueById(TblField *TblField, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_fields").Where("id=?", id).Preload("TblFieldOption").First(&TblField).Error; err != nil {

		return err
	}

	return nil
}

/*Update Channel Details*/
func (Ch ChannelStruct) UpdateChannelDetails(chn *TblChannel, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channels").Where("id=?", id).UpdateColumns(map[string]interface{}{"channel_name": chn.ChannelName, "channel_description": chn.ChannelDescription, "modified_by": chn.ModifiedBy, "modified_on": chn.ModifiedOn}).Error; err != nil {

		return err
	}

	return nil
}

/*Update Field Details*/
func (Ch ChannelStruct) UpdateFieldDetails(fds *TblField, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_fields").Where("id=?", id).UpdateColumns(map[string]interface{}{"field_name": fds.FieldName, "field_desc": fds.FieldDesc, "mandatory_field": fds.MandatoryField, "datetime_format": fds.DatetimeFormat, "time_format": fds.TimeFormat, "initial_value": fds.InitialValue, "placeholder": fds.Placeholder, "modified_on": fds.ModifiedOn, "modified_by": fds.ModifiedBy, "order_index": fds.OrderIndex, "url": fds.Url, "character_allowed": fds.CharacterAllowed}).Error; err != nil {

		return err
	}

	return nil
}

/*Update Field Option Details*/
func (Ch ChannelStruct) UpdateFieldOption(fdoption *TblFieldOption, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_field_options").Where("id=?", id).UpdateColumns(map[string]interface{}{"option_name": fdoption.OptionName, "option_value": fdoption.OptionValue, "modified_on": fdoption.ModifiedOn, "modified_by": fdoption.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

func (Ch ChannelStruct) UpdateFieldGroup(fldgrp *TblFieldGroup, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_field_groups").Where("id=?", id).UpdateColumns(map[string]interface{}{"modified_by": fldgrp.ModifiedBy, "modified_on": fldgrp.ModifiedOn}).Error; err != nil {

		return err

	}

	return nil
}

/**/
func (Ch ChannelStruct) GetNotInFieldId(group *[]TblGroupField, ids []int, fieldgroupid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_group_fields").Where("field_group_id = ? and field_id not in(?)", fieldgroupid, ids).Find(&group).Error; err != nil {

		return err
	}
	return nil
}

/*Check option already exist by fieldid*/
func (Ch ChannelStruct) CheckOptionAlreadyExist(optval *TblFieldOption, name string, fid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_field_options").Where("option_name=? and is_deleted!=1 and field_id=?", name, fid).First(&optval).Error; err != nil {

		return err
	}

	return nil
}

/*Get All Channel Permission Based*/
func (Ch ChannelStruct) GetAllChannel(chn *[]TblChannel, roleid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channels").Joins("inner join tbl_module_permissions on tbl_module_permissions.display_name = tbl_channels.channel_name").Joins("inner join tbl_role_permissions on tbl_role_permissions.permission_id= tbl_module_permissions.id").Where("tbl_role_permissions.role_id=?  and tbl_channels.is_deleted=0", roleid).Find(&chn).Error; err != nil {

		return err
	}

	return nil
}

/*Create Channel Categories*/
func (Ch ChannelStruct) CreateChannelCategory(channelcategory *TblChannelCategory, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_category").Create(&channelcategory).Error; err != nil {

		return err
	}

	return nil

}

/**/
func (Ch ChannelStruct) GetSelectedCategoryChannelById(ChannelCategory *[]TblChannelCategory, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_category").Where("channel_id=?", id).Find(&ChannelCategory).Error; err != nil {

		return err
	}

	return nil

}

/*Edit channel category*/
func (Ch ChannelStruct) GetCategoriseById(category *[]categories.TblCategory, id []int, DB *gorm.DB) error {

	if err := DB.Table("tbl_categories").Where("id in (?)", id).Order("id asc").Find(&category).Error; err != nil {

		return err
	}

	return nil

}

/*Delete Channel Category*/
func (Ch ChannelStruct) DeleteChannelCategoryByValue(category *TblChannelCategory, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_category").Where("id=?", id).Delete(&category).Error; err != nil {

		return err
	}

	return nil
}

/*CheckCategoryId Already Exists*/
func (Ch ChannelStruct) CheckChannelCategoryAlreadyExitst(category *TblChannelCategory, channelid int, categoryids string, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_category").Where("channel_id=? and category_id=?", channelid, categoryids).First(&category).Error; err != nil {

		return err
	}

	return nil

}

/**/
func (Ch ChannelStruct) GetChannelCategoryNotExist(category *[]TblChannelCategory, channelid int, categoryids []string, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_category").Where("channel_id=? and category_id not in (?)", channelid, categoryids).Find(&category).Error; err != nil {

		return err
	}

	return nil
}

func (Ch ChannelStruct) GetIdByCategoryValue(category *[]TblChannelCategory, channelid int, categoryids []string, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_category").Where("channel_id=? and category_id in (?)", channelid, categoryids).Find(&category).Error; err != nil {

		return err
	}

	return nil
}

func (Ch ChannelStruct) GetChannelCategoryDetails(category *[]TblChannelCategory, id []int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_category").Where("id in (?)", id).Find(&category).Error; err != nil {

		return err
	}

	return nil

}

func (Ch ChannelStruct) CheckChanName(channel *TblChannel, name string, id int, DB *gorm.DB) error {

	if id == 0 {
		if err := DB.Table("tbl_channels").Where("LOWER(TRIM(channel_name))=LOWER(TRIM(?)) and is_deleted = 0 ", name).First(&channel).Error; err != nil {

			return err
		}
	} else {
		if err := DB.Table("tbl_channels").Where("LOWER(TRIM(channel_name))=LOWER(TRIM(?)) and id not in(?) and is_deleted= 0 ", name, id).First(&channel).Error; err != nil {

			return err
		}
	}
	return nil
}

/*Delete FieldOption By fieldid*/
func (Ch ChannelStruct) DeleteOptionById(fieldopt *TblFieldOption, id []int, fid []int, DB *gorm.DB) error {

	if len(id) > 0 {

		if err := DB.Table("tbl_field_options").Where("id in (?)", id).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_by": fieldopt.DeletedBy, "deleted_on": fieldopt.DeletedOn}).Error; err != nil {

			return err
		}

	} else {

		if err := DB.Table("tbl_field_options").Where("field_id in (?)", fid).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_by": fieldopt.DeletedBy, "deleted_on": fieldopt.DeletedOn}).Error; err != nil {

			return err
		}
	}

	return nil

}

/*List Channel Entry*/
func (Ch ChannelStruct) ChannelEntryList(chentry *[]TblChannelEntries, limit, offset, chid int, filter EntriesFilter, publishedflg bool, RoleId int, activechannel bool, DB *gorm.DB) (chentcount int64, err error) {

	query := DB.Table("tbl_channel_entries").Select("tbl_channel_entries.*,tbl_users.username,tbl_channels.channel_name").Joins("inner join tbl_users on tbl_users.id = tbl_channel_entries.created_by").Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Where("tbl_channel_entries.is_deleted=0").Order("id desc")

	if RoleId != 1 {

		query = query.Where("channel_id in (select id from tbl_channels where channel_name in (select display_name from tbl_module_permissions inner join tbl_modules on tbl_modules.id = tbl_module_permissions.module_id inner join tbl_role_permissions on tbl_role_permissions.permission_id = tbl_module_permissions.id where role_id =(?) and tbl_modules.module_name='Entries' )) ", RoleId)

	}

	if activechannel {

		query = query.Where("tbl_channels.is_active =1")
	}

	if publishedflg {

		query = query.Where("tbl_channel_entries.status=1")

	}

	if chid != 0 {

		query = query.Where("tbl_channel_entries.channel_id=?", chid)
	}

	if filter.UserName != "" {

		query = query.Debug().Where("LOWER(TRIM(tbl_users.username)) ILIKE LOWER(TRIM(?))", "%"+filter.UserName+"%")

	}

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(title)) ILIKE LOWER(TRIM(?)) OR LOWER(TRIM(channel_name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")

	}

	if filter.Status != "" {

		query = query.Where("tbl_channel_entries.status=?", filter.Status)

	}
	if filter.Title != "" {

		query = query.Where("LOWER(TRIM(title)) ILIKE LOWER(TRIM(?))", "%"+filter.Title+"%")

	}

	if filter.ChannelName != "" {

		query = query.Where("LOWER(TRIM(channel_name)) ILIKE LOWER(TRIM(?))", "%"+filter.ChannelName+"%")

	}

	if limit != 0 {

		query.Limit(limit).Offset(offset).Order("id asc").Find(&chentry)

	} else {

		query.Find(&chentry).Count(&chentcount)

		return chentcount, nil
	}

	return 0, nil
}

/*Create channel entry*/
func (Ch ChannelStruct) CreateChannelEntry(entry *TblChannelEntries, DB *gorm.DB) (*TblChannelEntries, error) {

	if err := DB.Table("tbl_channel_entries").Create(&entry).Error; err != nil {

		return &TblChannelEntries{}, err

	}

	return entry, nil

}

/*create channel entry field*/
func (Ch ChannelStruct) CreateEntrychannelFields(entryfield *[]TblChannelEntryField, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entry_fields").Create(&entryfield).Error; err != nil {

		return err
	}

	return nil

}

/*create channel entry field*/
func (Ch ChannelStruct) CreateSingleEntrychannelFields(entryfield *TblChannelEntryField, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entry_fields").Create(&entryfield).Error; err != nil {

		return err
	}

	return nil

}

/*Delete Channel Entry Field*/
func (Ch ChannelStruct) DeleteChannelEntryId(chentry *TblChannelEntries, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entries").Where("id=?", chentry.Id).UpdateColumns(map[string]interface{}{"is_deleted": chentry.IsDeleted, "deleted_by": chentry.DeletedBy, "deleted_on": chentry.DeletedOn}).Error; err != nil {

		return err
	}

	return nil
}

/*Delete Channel Entry Field*/
func (Ch ChannelStruct) DeleteChannelEntryFieldId(chentry *TblChannelEntryField, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entry_fields").Where("channel_entry_id=?", id).UpdateColumns(map[string]interface{}{"deleted_by": chentry.DeletedBy, "deleted_on": chentry.DeletedOn}).Error; err != nil {

		return err
	}

	return nil
}

/*Edit Channel Entry Field*/
func (Ch ChannelStruct) GetChannelEntryDetailsById(tblchanentry *[]TblChannelEntryField, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entry_fields").Select("tbl_channel_entry_fields.*,tbl_fields.field_type_id").Joins("inner join tbl_fields on tbl_fields.Id = tbl_channel_entry_fields.field_id").Find(&tblchanentry).Error; err != nil {

		return err

	}

	return nil
}

/*Edit Channel Entry*/
func (Ch ChannelStruct) GetChannelEntryById(tblchanentry *TblChannelEntries, id int, DB *gorm.DB) error {

	// if err := DB.Table("tbl_channel_entries").Where("id=?", id).First(&tblchanentry).Error; err != nil {

	// 	return err

	// }
	if err := DB.Table("tbl_channel_entries").Where("is_deleted=0 and id=?", id).Preload("TblChannelEntryField", func(db *gorm.DB) *gorm.DB {
		return db.Select("tbl_channel_entry_fields.*,tbl_fields.field_type_id").Joins("inner join tbl_fields on tbl_fields.Id = tbl_channel_entry_fields.field_id")
	}).Find(&tblchanentry).Error; err != nil {

		return err

	}

	return nil
}

/*Update Channel Entry Details*/
func (Ch ChannelStruct) UpdateChannelEntryDetails(entry *TblChannelEntries, entryid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entries").Where("id=?", entryid).UpdateColumns(map[string]interface{}{"title": entry.Title, "description": entry.Description, "slug": entry.Slug, "cover_image": entry.CoverImage, "thumbnail_image": entry.ThumbnailImage, "meta_title": entry.MetaTitle, "meta_description": entry.MetaDescription, "keyword": entry.Keyword, "categories_id": entry.CategoriesId, "related_articles": entry.RelatedArticles, "status": entry.Status, "modified_on": entry.ModifiedOn, "modified_by": entry.ModifiedBy, "user_id": entry.UserId, "channel_id": entry.ChannelId}).Error; err != nil {

		return err
	}

	return nil

}

/*Update Channel Entry Details*/
func (Ch ChannelStruct) UpdateChannelEntryAdditionalDetails(entry TblChannelEntryField, DB gorm.DB) error {

	if err := DB.Table("tbl_channel_entry_fields").Where("id=?", entry.Id).UpdateColumns(map[string]interface{}{"field_name": entry.FieldName, "field_value": entry.FieldValue, "modified_by": entry.ModifiedBy, "modified_on": entry.ModifiedOn}).Error; err != nil {

		return err
	}

	return nil
}

func (Ch ChannelStruct) PublishQuery(chl *TblChannelEntries, id int, DB gorm.DB) error {

	if err := DB.Table("tbl_channel_entries").Where("id =?", id).UpdateColumns(map[string]interface{}{"status": chl.Status, "modified_on": chl.ModifiedOn, "modified_by": chl.ModifiedBy}).Error; err != nil {

		return err

	}

	return nil
}

func (Ch ChannelStruct) AllentryCount(DB *gorm.DB) (count int64, err error) {

	if err := DB.Table("tbl_channel_entries").Where("is_deleted = 0 ").Count(&count).Error; err != nil {

		return 0, err
	}

	return count, nil
}

func (Ch ChannelStruct) NewentryCount(DB *gorm.DB) (count int64, err error) {

	if err := DB.Table("tbl_channel_entries").Where("is_deleted = 0 AND created_on >=?", time.Now().AddDate(0, 0, -10)).Count(&count).Error; err != nil {

		return 0, err
	}

	return count, nil
}

func (Ch ChannelStruct) DeleteEntryByChannelId(id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entries").Where("channel_id=?", id).UpdateColumns(map[string]interface{}{"is_deleted": 1}).Error; err != nil {

		return err
	}

	return nil

}

func (Ch ChannelStruct) Newchannels(DB *gorm.DB) (chn []TblChannel, err error) {

	if err := DB.Table("tbl_channels").Select("tbl_channels.*,tbl_users.username,tbl_users.profile_image_path").
		Joins("inner join tbl_users on tbl_users.id = tbl_channels.created_by").
		Where("tbl_channels.is_deleted=0 and tbl_channels.is_active=1 and tbl_channels.created_on >= ?", time.Now().Add(-24*time.Hour).Format("2006-01-02 15:04:05")).
		Order("created_on desc").Limit(6).Find(&chn).Error; err != nil {

		return []TblChannel{}, err
	}

	return chn, nil

}

func (Ch ChannelStruct) Newentries(DB *gorm.DB) (entries []TblChannelEntries, err error) {

	if err := DB.Table("tbl_channel_entries").Select("tbl_channel_entries.*,tbl_users.username,tbl_users.profile_image_path").
		Joins("inner join tbl_users on tbl_users.id = tbl_channel_entries.created_by").Where("tbl_channel_entries.is_deleted=0 and tbl_channel_entries.created_on >=?", time.Now().Add(-24*time.Hour).Format("2006-01-02 15:04:05")).
		Order("created_on desc").Limit(6).Find(&entries).Error; err != nil {

		return []TblChannelEntries{}, err
	}

	return entries, nil

}

// update imagepath
func (Ch ChannelStruct) UpdateImagePath(Imagepath string, DB *gorm.DB) error {

	if err := DB.Model(TblChannelEntries{}).Where("cover_image=?", Imagepath).UpdateColumns(map[string]interface{}{
		"cover_image": ""}).Error; err != nil {

		return err
	}

	return nil

}

/*List Channel Entry*/
func (Ch ChannelStruct) ChannelEntryListForTemplates(chentry *[]TblChannelEntries, limit, offset int, filter EntriesFilter, DB *gorm.DB) (chentcount int64, err error) {

	query := DB.Table("tbl_channel_entries").Select("tbl_channel_entries.*,tbl_users.username,tbl_channels.channel_name").Joins("inner join tbl_users on tbl_users.id = tbl_channel_entries.created_by").Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Where("tbl_channel_entries.is_deleted=0 and tbl_channel_entries.status=1 and LOWER(TRIM(tbl_channels.channel_name))=LOWER(TRIM(?))", filter.ChannelName).Order("id desc")

	if filter.CategoryId != "" {

		query = query.Where("tbl_channel_entries.id in (SELECT id FROM (SELECT id, unnest(string_to_array(categories_id, ',')) AS categoryid	FROM tbl_channel_entries) AS subquery WHERE categoryid = (?))", filter.CategoryId)
	}

	if limit != 0 {

		query.Limit(limit).Offset(offset).Order("id asc").Find(&chentry)

	} else {

		query.Find(&chentry).Count(&chentcount)

		return chentcount, nil
	}

	return 0, nil
}
func (ch ChannelStruct) GetGraphqlChannelList(DB *gorm.DB, memberid, limit, offset int) (channellist []TblChannel, channelCount int64, err error) {

	if memberid > 0 {

		query := DB.Table("tbl_channels").Select("distinct on (tbl_channels.id) tbl_channels.*").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
			Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ?", memberid)

		if err := query.Limit(limit).Offset(offset).Find(&channellist).Error; err != nil {

			return []TblChannel{}, 0, err
		}

		countquery := DB.Table("tbl_channels").Distinct("tbl_channels.id").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
			Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ?", memberid)

		if err := countquery.Count(&channelCount).Error; err != nil {

			return []TblChannel{}, 0, err

		}

	} else {

		query := DB.Table("tbl_channels").Select("distinct on (tbl_channels.id) tbl_channels.*").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
			Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1")

		if err := query.Limit(limit).Offset(offset).Find(&channellist).Error; err != nil {

			return []TblChannel{}, 0, err
		}

		countquery := DB.Table("tbl_channels").Distinct("tbl_channels.id").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
			Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1")

		if err := countquery.Count(&channelCount).Error; err != nil {

			return []TblChannel{}, 0, err

		}
	}

	return channellist, channelCount, nil
}

func (ch ChannelStruct) GetGraphqlChannelDetailsById(DB *gorm.DB, memberid, channelid int) (channel TblChannel, err error) {

	if memberid > 0 {

		if err = DB.Table("tbl_channels").Select("distinct on (tbl_channels.id) tbl_channels.*").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
			Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ? and  tbl_channels.id = ?", memberid, channelid).First(&channel).Error; err != nil {

			return TblChannel{}, err
		}

	} else {

		if err := DB.Table("tbl_channels").Select("distinct on (tbl_channels.id) tbl_channels.*").Joins("inner join tbl_channel_entries on tbl_channel_entries.channel_id = tbl_channels.id").
			Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_channels.id = ?", channelid).First(&channel).Error; err != nil {

			return TblChannel{}, err
		}

	}

	return channel, nil
}

func (ch ChannelStruct) GetGraphqlChannelEntryDetailsById(DB *gorm.DB, memberid int, channelEntryId int, channelId, categoryId *int) (channelEntry TblChannelEntries, err error) {

	var query *gorm.DB

	if memberid > 0 {

		if channelId != nil {

			query = DB.Table("tbl_channel_entries").Select("tbl_channel_entries.id,tbl_channel_entries.title,tbl_channel_entries.slug,tbl_channel_entries.description,tbl_channel_entries.user_id,tbl_channel_entries.channel_id,tbl_channel_entries.status,tbl_channel_entries.is_active,tbl_channel_entries.deleted_by,tbl_channel_entries.deleted_on,tbl_channel_entries.created_on,tbl_channel_entries.created_by,tbl_channel_entries.modified_by,tbl_channel_entries.modified_on,tbl_channel_entries.cover_image,tbl_channel_entries.thumbnail_image,tbl_channel_entries.meta_title,tbl_channel_entries.meta_description,tbl_channel_entries.keyword,tbl_channel_entries.categories_id,tbl_channel_entries.related_articles").
				Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
				Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
				Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ? and tbl_channel_entries.channel_id = ? and tbl_channel_entries.id = ?", memberid, channelId, channelEntryId)

			if categoryId != nil {

				query = query.Where("'" + strconv.Itoa(*categoryId) + "' = ANY(STRING_TO_ARRAY(tbl_channel_entries.categories_id," + "','" + "))")

			}

			if err = query.First(&channelEntry).Error; err != nil {

				return TblChannelEntries{}, err

			}

		} else {

			query = DB.Table("tbl_channel_entries").Select("tbl_channel_entries.id,tbl_channel_entries.title,tbl_channel_entries.slug,tbl_channel_entries.description,tbl_channel_entries.user_id,tbl_channel_entries.channel_id,tbl_channel_entries.status,tbl_channel_entries.is_active,tbl_channel_entries.deleted_by,tbl_channel_entries.deleted_on,tbl_channel_entries.created_on,tbl_channel_entries.created_by,tbl_channel_entries.modified_by,tbl_channel_entries.modified_on,tbl_channel_entries.cover_image,tbl_channel_entries.thumbnail_image,tbl_channel_entries.meta_title,tbl_channel_entries.meta_description,tbl_channel_entries.keyword,tbl_channel_entries.categories_id,tbl_channel_entries.related_articles").
				Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
				Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
				Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ? and tbl_channel_entries.id = ?", memberid, channelEntryId)

			if categoryId != nil {

				query = query.Where("'" + strconv.Itoa(*categoryId) + "' = ANY(STRING_TO_ARRAY(tbl_channel_entries.categories_id," + "','" + "))")

			}

			if err = query.First(&channelEntry).Error; err != nil {

				return TblChannelEntries{}, err

			}

		}

	} else {

		if channelId != nil {

			query = DB.Table("tbl_channel_entries").Select("tbl_channel_entries.id,tbl_channel_entries.title,tbl_channel_entries.slug,tbl_channel_entries.description,tbl_channel_entries.user_id,tbl_channel_entries.channel_id,tbl_channel_entries.status,tbl_channel_entries.is_active,tbl_channel_entries.deleted_by,tbl_channel_entries.deleted_on,tbl_channel_entries.created_on,tbl_channel_entries.created_by,tbl_channel_entries.modified_by,tbl_channel_entries.modified_on,tbl_channel_entries.cover_image,tbl_channel_entries.thumbnail_image,tbl_channel_entries.meta_title,tbl_channel_entries.meta_description,tbl_channel_entries.keyword,tbl_channel_entries.categories_id,tbl_channel_entries.related_articles").
				Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_channel_entries.channel_id = ? and tbl_channel_entries.id = ?", channelId, channelEntryId)

			if categoryId != nil {

				query = query.Where("'" + strconv.Itoa(*categoryId) + "' = ANY(STRING_TO_ARRAY(tbl_channel_entries.categories_id," + "','" + "))")

			}

			if err = query.First(&channelEntry).Error; err != nil {

				return TblChannelEntries{}, err

			}

		} else {

			query = DB.Table("tbl_channel_entries").Select("tbl_channel_entries.id,tbl_channel_entries.title,tbl_channel_entries.slug,tbl_channel_entries.description,tbl_channel_entries.user_id,tbl_channel_entries.channel_id,tbl_channel_entries.status,tbl_channel_entries.is_active,tbl_channel_entries.deleted_by,tbl_channel_entries.deleted_on,tbl_channel_entries.created_on,tbl_channel_entries.created_by,tbl_channel_entries.modified_by,tbl_channel_entries.modified_on,tbl_channel_entries.cover_image,tbl_channel_entries.thumbnail_image,tbl_channel_entries.meta_title,tbl_channel_entries.meta_description,tbl_channel_entries.keyword,tbl_channel_entries.categories_id,tbl_channel_entries.related_articles").
				Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_channel_entries.id = ?", channelEntryId)

			if categoryId != nil {

				query = query.Where("'" + strconv.Itoa(*categoryId) + "' = ANY(STRING_TO_ARRAY(tbl_channel_entries.categories_id," + "','" + "))")

			}

			if err = query.First(&channelEntry).Error; err != nil {

				return TblChannelEntries{}, err

			}

		}
	}

	return channelEntry, nil
}

func (ch ChannelStruct) GetGraphqlEntriesCategoryByParentId(DB *gorm.DB, categoryId int) (category categories.TblCategory, err error) {

	if err = DB.Table("tbl_categories").Where("is_deleted = 0 and id = ?", categoryId).First(&category).Error; err != nil {

		return categories.TblCategory{}, err
	}

	return category, nil
}

func (ch ChannelStruct) GetGraphqlChannelEntrieslistByChannelId(DB *gorm.DB, memberid int, channelid, categoryid *int, limit, offset int) (channelEntries []TblChannelEntries, count int64, err error) {

	if memberid > 0 {

		query := DB.Table("tbl_channel_entries").Select("tbl_channel_entries.id,tbl_channel_entries.title,tbl_channel_entries.slug,tbl_channel_entries.description,tbl_channel_entries.user_id,tbl_channel_entries.channel_id,tbl_channel_entries.status,tbl_channel_entries.is_active,tbl_channel_entries.deleted_by,tbl_channel_entries.deleted_on,tbl_channel_entries.created_on,tbl_channel_entries.created_by,tbl_channel_entries.modified_by,tbl_channel_entries.modified_on,tbl_channel_entries.cover_image,tbl_channel_entries.thumbnail_image,tbl_channel_entries.meta_title,tbl_channel_entries.meta_description,tbl_channel_entries.keyword,tbl_channel_entries.categories_id,tbl_channel_entries.related_articles").
			Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ? and tbl_channel_entries.channel_id = ?", memberid, channelid).Limit(limit).Offset(offset).Order("tbl_channel_entries.id desc")

		countquery := DB.Table("tbl_channel_entries").Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ? and tbl_channel_entries.channel_id = ?", memberid, channelid)

		if categoryid != nil {

			query = query.Where("'" + strconv.Itoa(*categoryid) + "' = ANY(STRING_TO_ARRAY(tbl_channel_entries.categories_id," + "','" + "))")

			countquery = countquery.Where("'" + strconv.Itoa(*categoryid) + "' = ANY(STRING_TO_ARRAY(tbl_channel_entries.categories_id," + "','" + "))")
		}

		if err = query.Find(&channelEntries).Error; err != nil {

			return []TblChannelEntries{}, 0, err
		}

		if err = countquery.Count(&count).Error; err != nil {

			return []TblChannelEntries{}, 0, err
		}

	} else {

		query := DB.Table("tbl_channel_entries").Select("tbl_channel_entries.id,tbl_channel_entries.title,tbl_channel_entries.slug,tbl_channel_entries.description,tbl_channel_entries.user_id,tbl_channel_entries.channel_id,tbl_channel_entries.status,tbl_channel_entries.is_active,tbl_channel_entries.deleted_by,tbl_channel_entries.deleted_on,tbl_channel_entries.created_on,tbl_channel_entries.created_by,tbl_channel_entries.modified_by,tbl_channel_entries.modified_on,tbl_channel_entries.cover_image,tbl_channel_entries.thumbnail_image,tbl_channel_entries.meta_title,tbl_channel_entries.meta_description,tbl_channel_entries.keyword,tbl_channel_entries.categories_id,tbl_channel_entries.related_articles").
			Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").
			Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_channel_entries.channel_id = ?", channelid).Limit(limit).Offset(offset).Order("tbl_channel_entries.id desc")

		countquery := DB.Table("tbl_channel_entries").Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_channel_entries.channel_id = ?", channelid)

		if categoryid != nil {

			query = query.Where("'" + strconv.Itoa(*categoryid) + "' = ANY(STRING_TO_ARRAY(tbl_channel_entries.categories_id," + "','" + "))")

			countquery = countquery.Where("'" + strconv.Itoa(*categoryid) + "' = ANY(STRING_TO_ARRAY(tbl_channel_entries.categories_id," + "','" + "))")
		}

		if err = query.Find(&channelEntries).Error; err != nil {

			return []TblChannelEntries{}, 0, err
		}

		if err = countquery.Count(&count).Error; err != nil {

			return []TblChannelEntries{}, 0, err
		}

	}

	return channelEntries, count, nil
}

func (ch ChannelStruct) GetGraphqlChannelEntriesList(DB *gorm.DB, memberid int, categoryid *int, limit, offset int) (channelEntries []TblChannelEntries, count int64, err error) {

	if memberid > 0 {

		query := DB.Table("tbl_channel_entries").Select("tbl_channel_entries.id,tbl_channel_entries.title,tbl_channel_entries.slug,tbl_channel_entries.description,tbl_channel_entries.user_id,tbl_channel_entries.channel_id,tbl_channel_entries.status,tbl_channel_entries.is_active,tbl_channel_entries.deleted_by,tbl_channel_entries.deleted_on,tbl_channel_entries.created_on,tbl_channel_entries.created_by,tbl_channel_entries.modified_by,tbl_channel_entries.modified_on,tbl_channel_entries.cover_image,tbl_channel_entries.thumbnail_image,tbl_channel_entries.meta_title,tbl_channel_entries.meta_description,tbl_channel_entries.keyword,tbl_channel_entries.categories_id,tbl_channel_entries.related_articles").
			Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ?", memberid).Limit(limit).Offset(offset).Order("tbl_channel_entries.id desc")

		countquery := DB.Table("tbl_channel_entries").Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Joins("inner join tbl_access_control_pages on tbl_access_control_pages.entry_id = tbl_channel_entries.id").Joins("inner join tbl_access_control_user_group on tbl_access_control_user_group.id = tbl_access_control_pages.access_control_user_group_id").
			Joins("inner join tbl_member_groups on tbl_member_groups.id = tbl_access_control_user_group.member_group_id").Joins("inner join tbl_members on tbl_members.member_group_id = tbl_member_groups.id").
			Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1 and tbl_members.is_deleted = 0 and tbl_member_groups.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0  and tbl_access_control_user_group.is_deleted = 0 and tbl_members.id = ?", memberid)

		if categoryid != nil {

			query = query.Where("'" + strconv.Itoa(*categoryid) + "' = ANY(STRING_TO_ARRAY(tbl_channel_entries.categories_id," + "','" + "))")

			countquery = countquery.Where("'" + strconv.Itoa(*categoryid) + "' = ANY(STRING_TO_ARRAY(tbl_channel_entries.categories_id," + "','" + "))")
		}

		if err = query.Find(&channelEntries).Error; err != nil {

			return []TblChannelEntries{}, 0, err
		}

		if err = countquery.Count(&count).Error; err != nil {

			return []TblChannelEntries{}, 0, err
		}

	} else {

		query := DB.Table("tbl_channel_entries").Select("tbl_channel_entries.id,tbl_channel_entries.title,tbl_channel_entries.slug,tbl_channel_entries.description,tbl_channel_entries.user_id,tbl_channel_entries.channel_id,tbl_channel_entries.status,tbl_channel_entries.is_active,tbl_channel_entries.deleted_by,tbl_channel_entries.deleted_on,tbl_channel_entries.created_on,tbl_channel_entries.created_by,tbl_channel_entries.modified_by,tbl_channel_entries.modified_on,tbl_channel_entries.cover_image,tbl_channel_entries.thumbnail_image,tbl_channel_entries.meta_title,tbl_channel_entries.meta_description,tbl_channel_entries.keyword,tbl_channel_entries.categories_id,tbl_channel_entries.related_articles").
			Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1").Limit(limit).Offset(offset).Order("tbl_channel_entries.id desc")

		countquery := DB.Table("tbl_channel_entries").Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Where("tbl_channels.is_deleted = 0 and tbl_channels.is_active = 1 and tbl_channel_entries.is_deleted = 0 and tbl_channel_entries.status = 1")

		if categoryid != nil {

			query = query.Where("'" + strconv.Itoa(*categoryid) + "' = ANY(STRING_TO_ARRAY(tbl_channel_entries.categories_id," + "','" + "))")

			countquery = countquery.Where("'" + strconv.Itoa(*categoryid) + "' = ANY(STRING_TO_ARRAY(tbl_channel_entries.categories_id," + "','" + "))")
		}

		if err = query.Find(&channelEntries).Error; err != nil {

			return []TblChannelEntries{}, 0, err
		}

		if err = countquery.Count(&count).Error; err != nil {

			return []TblChannelEntries{}, 0, err
		}

	}

	return channelEntries, count, nil
}

func (Ch ChannelStruct) GetChannelCategoryDetailsByChannelId(category *[]TblChannelCategory, id []int, DB *gorm.DB) error {

	if err := DB.Debug().Table("tbl_channel_category").Where("channel_id in (?)", id).Find(&category).Error; err != nil {

		return err
	}

	return nil

}

func (Ch ChannelStruct) GetAuthorDetails(DB *gorm.DB, authorId int) (authorDetail Author, err error) {

	if err := DB.Table("tbl_users").Where("tbl_users.is_deleted = 0 and tbl_users.id = ?", authorId).First(&authorDetail).Error; err != nil {

		return Author{}, err
	}

	return authorDetail, nil
}

func (Ch ChannelStruct) GetSectionsUnderEntries(DB *gorm.DB, channelId, sectionTypeId int) (sections []TblField, err error) {

	if err = DB.Table("tbl_group_fields").Select("tbl_fields.*,tbl_field_types.type_name").Joins("inner join tbl_fields on tbl_fields.id = tbl_group_fields.field_id").Joins("inner join tbl_field_types on tbl_field_types.id = tbl_fields.field_type_id").
		Where("tbl_fields.is_deleted = 0 and tbl_field_types.is_deleted = 0 and tbl_fields.field_type_id = ? and tbl_group_fields.channel_id = ?", sectionTypeId, channelId).Find(&sections).Error; err != nil {

		return []TblField{}, err
	}

	return sections, nil
}

func (Ch ChannelStruct) GetFieldsInEntries(DB *gorm.DB, channelId, sectionTypeId int) (fields []TblField, err error) {

	if err = DB.Table("tbl_group_fields").Select("tbl_fields.*,tbl_field_types.type_name").Joins("inner join tbl_fields on tbl_fields.id = tbl_group_fields.field_id").Joins("inner join tbl_field_types on tbl_field_types.id = tbl_fields.field_type_id").
		Where("tbl_fields.is_deleted = 0 and tbl_field_types.is_deleted = 0 and tbl_fields.field_type_id != ? and tbl_group_fields.channel_id = ?", sectionTypeId, channelId).Find(&fields).Error; err != nil {

		return []TblField{}, err
	}

	return fields, nil
}

func (Ch ChannelStruct) GetFieldValue(DB *gorm.DB, fieldId, entryId int) (fieldvalue TblChannelEntryField, err error) {

	if err = DB.Table("tbl_channel_entry_fields").Where("tbl_channel_entry_fields.field_id = ? and tbl_channel_entry_fields.channel_entry_id = ?", fieldId, entryId).First(&fieldvalue).Error; err != nil {

		return TblChannelEntryField{}, err
	}

	return fieldvalue, nil
}

func (ch ChannelStruct) GetFieldOptions(DB *gorm.DB, fieldId int) (fieldOptions []TblFieldOption, err error) {

	if err = DB.Table("tbl_field_options").Where("tbl_field_options.is_deleted = 0 and tbl_field_options.field_id = ?", fieldId).Find(&fieldOptions).Error; err != nil {

		return []TblFieldOption{}, err
	}

	return fieldOptions, nil
}

func (ch ChannelStruct) GetMemberProfile(DB *gorm.DB, memberid int) (memberProfile TblMemberProfiles, err error) {

	if err = DB.Table("tbl_member_profiles").Select("tbl_member_profiles.*").Joins("inner join tbl_members on tbl_members.id = tbl_member_profiles.member_id").Where("tbl_members.is_deleted = 0 and tbl_members.id = ?", memberid).First(&memberProfile).Error; err != nil {

		return TblMemberProfiles{}, err
	}

	return memberProfile, nil

}

func (ch ChannelStruct) MakeFeature(channelid, entryid int, DB *gorm.DB) (err error) {

	DB.Model(TblChannelEntries{}).UpdateColumns(map[string]interface{}{"feature": 0}).Where("channel_id=?", channelid)

	if err := DB.Model(TblChannelEntries{}).UpdateColumns(map[string]interface{}{"feature": 1}).Where("id=? and channel_id=?", entryid, channelid).Error; err != nil {

		return err
	}

	return nil
}
