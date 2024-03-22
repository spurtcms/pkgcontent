package channels

import (
	"time"

	"github.com/spurtcms/pkgcontent/categories"
	"gorm.io/gorm"
)

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
	Feature              int                        `gorm:"DEFAULT:0"`
	ViewCount            int                        `gorm:"DEFAULT:0"`
}

/*List Channel Entry*/
func (Ch Channelmodel) ChannelEntryList(chentry *[]TblChannelEntries, limit, offset, chid int, filter EntriesFilter, publishedflg bool, activechannel bool, DB *gorm.DB) (chentcount int64, err error) {

	query := DB.Table("tbl_channel_entries").Select("tbl_channel_entries.*,tbl_users.username,tbl_channels.channel_name").Joins("inner join tbl_users on tbl_users.id = tbl_channel_entries.created_by").Joins("inner join tbl_channels on tbl_channels.id = tbl_channel_entries.channel_id").Where("tbl_channel_entries.is_deleted=0").Order("id desc")

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
func (Ch Channelmodel) CreateChannelEntry(entry *TblChannelEntries, DB *gorm.DB) (*TblChannelEntries, error) {

	if err := DB.Table("tbl_channel_entries").Create(&entry).Error; err != nil {

		return &TblChannelEntries{}, err

	}

	return entry, nil

}

/*create channel entry field*/
func (Ch Channelmodel) CreateEntrychannelFields(entryfield *[]TblChannelEntryField, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entry_fields").Create(&entryfield).Error; err != nil {

		return err
	}

	return nil

}

/*create channel entry field*/
func (Ch Channelmodel) CreateSingleEntrychannelFields(entryfield *TblChannelEntryField, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entry_fields").Create(&entryfield).Error; err != nil {

		return err
	}

	return nil

}

/*Delete Channel Entry Field*/
func (Ch Channelmodel) DeleteChannelEntryId(chentry *TblChannelEntries, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entries").Where("id=?", chentry.Id).UpdateColumns(map[string]interface{}{"is_deleted": chentry.IsDeleted, "deleted_by": chentry.DeletedBy, "deleted_on": chentry.DeletedOn}).Error; err != nil {

		return err
	}

	return nil
}

/*Delete Channel Entry Field*/
func (Ch Channelmodel) DeleteChannelEntryFieldId(chentry *TblChannelEntryField, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entry_fields").Where("channel_entry_id=?", id).UpdateColumns(map[string]interface{}{"deleted_by": chentry.DeletedBy, "deleted_on": chentry.DeletedOn}).Error; err != nil {

		return err
	}

	return nil
}

/*Edit Channel Entry Field*/
func (Ch Channelmodel) GetChannelEntryDetailsById(tblchanentry *[]TblChannelEntryField, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entry_fields").Select("tbl_channel_entry_fields.*,tbl_fields.field_type_id").Joins("inner join tbl_fields on tbl_fields.Id = tbl_channel_entry_fields.field_id").Find(&tblchanentry).Error; err != nil {

		return err

	}

	return nil
}

/*Edit Channel Entry*/
func (Ch Channelmodel) GetChannelEntryById(tblchanentry *TblChannelEntries, id int, DB *gorm.DB) error {

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
func (Ch Channelmodel) UpdateChannelEntryDetails(entry *TblChannelEntries, entryid int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entries").Where("id=?", entryid).UpdateColumns(map[string]interface{}{"title": entry.Title, "description": entry.Description, "slug": entry.Slug, "cover_image": entry.CoverImage, "thumbnail_image": entry.ThumbnailImage, "meta_title": entry.MetaTitle, "meta_description": entry.MetaDescription, "keyword": entry.Keyword, "categories_id": entry.CategoriesId, "related_articles": entry.RelatedArticles, "status": entry.Status, "modified_on": entry.ModifiedOn, "modified_by": entry.ModifiedBy, "user_id": entry.UserId, "channel_id": entry.ChannelId}).Error; err != nil {

		return err
	}

	return nil

}

/*Update Channel Entry Details*/
func (Ch Channelmodel) UpdateChannelEntryAdditionalDetails(entry TblChannelEntryField, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entry_fields").Where("id=?", entry.Id).UpdateColumns(map[string]interface{}{"field_name": entry.FieldName, "field_value": entry.FieldValue, "modified_by": entry.ModifiedBy, "modified_on": entry.ModifiedOn}).Error; err != nil {

		return err
	}

	return nil
}

func (Ch Channelmodel) PublishQuery(chl *TblChannelEntries, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_channel_entries").Where("id =?", id).UpdateColumns(map[string]interface{}{"status": chl.Status, "modified_on": chl.ModifiedOn, "modified_by": chl.ModifiedBy}).Error; err != nil {

		return err

	}

	return nil
}
