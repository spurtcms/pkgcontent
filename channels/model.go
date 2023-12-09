package channels

import (
	"time"

	"github.com/spurtcms/spurtcms-content/categories"
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
	CharacterAllowed int
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
	ModifiedOn         time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy         int       `gorm:"DEFAULT:NULL"`
	DateString         string    `gorm:"-"`
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
	CreatedDate          string                   `gorm:"-"`
	ModifiedDate         string                   `gorm:"-"`
	Username             string                   `gorm:"<-:false"`
	TblChannelEntryField []TblChannelEntryField   `gorm:"<-:false; foreignKey:ChannelEntryId"`
	Category             []categories.TblCategory `gorm:"<-:false; foreignKey:Id"`
	CategoryGroup        string                   `gorm:"<-:false"`
	ChannelName          string
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
func (Ch ChannelStruct) Channellist(chn *[]TblChannel, limit, offset int, filter Filter, activestatus bool, DB *gorm.DB) (chcount int64, error error) {

	query := DB.Table("tbl_channels").Where("is_deleted = 0").Order("id desc")

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

/*Get FieldGroupById*/
func (Ch ChannelStruct) GetFieldGroupById(groupfield *[]TblGroupField, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_group_fields").Select("tbl_fields.field_name,tbl_field_options.option_name,tbl_field_options.option_value").Where("field_group_id=?", id).Joins("INNER JOIN TBL_FIELDS ON TBL_GROUP_FIELDS.FIELD_ID = TBL_FIELDS.ID").Joins("LEFT JOIN TBL_FIELD_OPTIONS ON TBL_FIELDS.ID = TBL_FIELD_OPTIONS.FIELD_ID").Find(&groupfield).Error; err != nil {

		return err
	}

	return nil
}

/*Getfieldid using fieldgroupid*/
func (Ch ChannelStruct) GetFieldIdByGroupId(grpfield *[]TblGroupField, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_group_fields").Where("field_group_id=?", id).Find(&grpfield).Error; err != nil {

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

	if err := DB.Table("tbl_fields").Where("id=?", id).UpdateColumns(map[string]interface{}{"field_name": fds.FieldName, "field_desc": fds.FieldDesc, "mandatory_field": fds.MandatoryField, "datetime_format": fds.DatetimeFormat, "time_format": fds.TimeFormat, "initial_value": fds.InitialValue, "placeholder": fds.Placeholder, "modified_on": fds.ModifiedOn, "modified_by": fds.ModifiedBy, "order_index": fds.OrderIndex, "url": fds.Url}).Error; err != nil {

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
func (Ch ChannelStruct) ChannelEntryList(chentry *[]TblChannelEntries, limit, offset, chid int, filter EntriesFilter, DB *gorm.DB) (chentcount int64, err error) {

	query := DB.Table("tbl_channel_entries").Select("tbl_channel_entries.*,tbl_users.username,tbl_channels.channel_name").Joins("inner join tbl_users on tbl_users.id = tbl_channel_entries.user_id").Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Where("tbl_channel_entries.is_deleted=0").Order("id desc")

	if chid != 0 {

		query = query.Where("tbl_channel_entries.channel_id=?", chid)
	}

	if filter.Title != "" {

		query = query.Where("LOWER(TRIM(title)) ILIKE LOWER(TRIM(?)) OR LOWER(TRIM(channel_name)) ILIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")

	}

	if filter.Status != "" {

		query = query.Where("tbl_channel_entries.status=?", filter.Status)

	}
	if filter.Title != "" {

		query = query.Where("LOWER(TRIM(title)) ILIKE LOWER(TRIM(?))", filter.Title)

	}

	if filter.ChannelName != "" {

		query = query.Where("LOWER(TRIM(channel_name)) ILIKE LOWER(TRIM(?))", filter.ChannelName)

	}

	if limit != 0 {

		query.Limit(limit).Offset(offset).Order("id asc").Find(&chentry)

	} else {

		query.Find(&chentry).Count(&chentcount)

		return chentcount, nil
	}

	return 0, nil
}
