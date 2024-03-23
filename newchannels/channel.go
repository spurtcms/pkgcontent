package channels

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spurtcms/pkgcontent/categories"
	"github.com/spurtcms/pkgcore/auth"
	"gorm.io/gorm"
)

/*Create Instance*/
func DefaultChannel(db *gorm.DB) *Channels {

	err := db.AutoMigrate(
		&TblFieldType{},
		&TblChannel{},
		&TblField{},
		&TblFieldGroup{},
		&TblFieldOption{},
		&TblGroupField{},
		&TblChannelCategory{},
	)

	if err != nil {
		//panic terminate the server
		panic(err)

	}

	channels := new(Channels)

	channels.ChannelRepository.DatabaseConnection = db

	return channels
}

/*Channel Repository defines the behaviour of a Channel*/
type ChannelRepository interface {
	CreateChannel(channelcreate ChannelCreate) error

	CreateChannelFields(channelid int, Sections []Section, FieldValues []Fiedlvalue)

	GetChannels(limit, offset int, filter Filter, activestatus bool) (channelList []TblChannel, channelcount int, err error)

	GetchannelByName(Channelname string) (channel TblChannel, err error)

	DeleteChannel(ChannelId int) error

	GetChannelsById(channelid int) (channelList TblChannel, section []Section, fields []Fiedlvalue, SelectedCategories []categories.Arrangecategories, err error)
}

/*Authentication defines the behaviour of the auth*/
type Authentication interface {
	Authenticate(Authority *auth.Authorization) (flg bool, user, roleid int, err error)
}

type Auth struct {
	Authentication Authentication
}

type BasicChannel struct {
	ChannelRepository  ChannelRepository //ChannelRepository have all methods
	DatabaseConnection *gorm.DB          //var holds db connections string
}

type Channels struct {
	ChannelRepository BasicChannel //ChannelRepository have all method
	Permissions       PermissionAuthentication
	Authentication    Auth //Check jwt tokens only
}

type PermissionAuthentication struct {
}

type Action string

const ( //for permission check
	Create Action = "Create"

	Read Action = "View"

	Update Action = "Update"

	Delete Action = "Delete"

	CRUD Action = "CRUD"
)

type Channelmodel struct{}

var cmod Channelmodel

func (a Auth) Authenticate(Authority *auth.Authorization) (flg bool, user, roleid int, err error) {

	userid, roleid, err := auth.VerifyToken(Authority.Token, Authority.Secret)

	//auth logic
	return true, userid, roleid, err
}

func (per PermissionAuthentication) IsGranted(modulename string, Permission Action) (flag bool, err error) {

	flag, perr := auth.Authorization.IsGranted(auth.Authorization{}, modulename, auth.Action(Permission))

	return flag, perr

}

func (Ch BasicChannel) CreateChannel(channelcreate ChannelCreate) error {

	/*create channel*/
	var channel TblChannel

	channel.ChannelName = channelcreate.ChannelName

	channel.ChannelDescription = channelcreate.ChannelDescription

	channel.SlugName = strings.ToLower(strings.ReplaceAll(channelcreate.ChannelName, " ", " "))

	channel.IsActive = 1

	channel.CreatedBy = channelcreate.CreatedBy

	channel.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	ch, chanerr := cmod.CreateChannel(&channel, Ch.DatabaseConnection)

	if chanerr != nil {

		log.Println(chanerr)

		return nil
	}

	/*This is for module permission creation*/
	var modperms auth.TblModulePermission

	modperms.DisplayName = ch.ChannelName

	modperms.RouteName = "/channel/entrylist/" + strconv.Itoa(ch.Id)

	modperms.SlugName = strings.ReplaceAll(strings.ToLower(ch.ChannelName), " ", "_")

	modperms.CreatedBy = channelcreate.CreatedBy

	modperms.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	modperms.ModuleId = 8

	modperms.FullAccessPermission = 1

	auth.AS.CreateModulePermission(&modperms, Ch.DatabaseConnection)

	for _, categoriesid := range channelcreate.CategoryIds {

		var channelcategory TblChannelCategory

		channelcategory.ChannelId = ch.Id

		channelcategory.CategoryId = categoriesid

		channelcategory.CreatedAt = channelcreate.CreatedBy

		channelcategory.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		err := cmod.CreateChannelCategory(&channelcategory, Ch.DatabaseConnection)

		if err != nil {

			log.Println(err)

		}

	}

	return nil
}

func (Ch BasicChannel) CreateChannelFields(createfield CreateChannelFields) error {

	/*Temp store section id*/
	type tempsection struct {
		Id           int
		SectionId    int
		NewSectionId int
	}

	var TempSections []tempsection

	/*create Section*/
	for _, sectionvalue := range createfield.Sections {

		var cfld TblField

		cfld.FieldName = strings.TrimSpace(sectionvalue.SectionName)

		cfld.FieldTypeId = sectionvalue.MasterFieldId

		cfld.CreatedBy = createfield.CreatedBy

		cfld.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		cfid, fiderr := cmod.CreateFields(&cfld, Ch.DatabaseConnection)

		if fiderr != nil {

			log.Println(fiderr)
		}

		/*create group field*/
		var grpfield TblGroupField

		grpfield.ChannelId = createfield.ChannelId

		grpfield.FieldId = cfid.Id

		grpfielderr := cmod.CreateGroupField(&grpfield, Ch.DatabaseConnection)

		if grpfielderr != nil {

			log.Println(grpfielderr)

		}

		var TempSection tempsection

		TempSection.Id = cfid.Id

		TempSection.SectionId = sectionvalue.SectionId

		TempSection.NewSectionId = sectionvalue.SectionNewId

		TempSections = append(TempSections, TempSection)

	}

	/*create field*/
	for _, val := range createfield.FieldValues {

		var cfld TblField

		cfld.FieldName = strings.TrimSpace(val.FieldName)

		cfld.FieldTypeId = val.MasterFieldId

		cfld.CreatedBy = createfield.CreatedBy

		cfld.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		cfld.OrderIndex = val.OrderIndex

		cfld.ImagePath = val.IconPath

		cfld.CharacterAllowed = val.CharacterAllowed

		cfld.Url = val.Url

		if val.MasterFieldId == 4 {

			cfld.DatetimeFormat = val.DateFormat

			cfld.TimeFormat = val.TimeFormat

		}
		if val.MasterFieldId == 6 {

			cfld.DatetimeFormat = val.DateFormat
		}

		if len(val.OptionValue) > 0 {

			cfld.OptionExist = 1
		}

		for _, sectionid := range TempSections {

			if sectionid.SectionId == val.SectionId && sectionid.NewSectionId == val.SectionNewId {

				cfld.SectionParentId = sectionid.Id

			}

		}

		cfid, fiderr := cmod.CreateFields(&cfld, Ch.DatabaseConnection)

		if fiderr != nil {

			log.Println(fiderr)

		}

		/*option value create*/
		for _, opt := range val.OptionValue {

			var fldopt TblFieldOption

			fldopt.OptionName = opt.Value

			fldopt.OptionValue = opt.Value

			fldopt.FieldId = cfid.Id

			fldopt.CreatedBy = createfield.CreatedBy

			fldopt.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			fopterr := cmod.CreateFieldOption(&fldopt, Ch.DatabaseConnection)

			if fopterr != nil {

				log.Println(fopterr)

			}

		}

		/*create group field*/
		var grpfield TblGroupField

		grpfield.ChannelId = createfield.ChannelId

		grpfield.FieldId = cfid.Id

		grpfielderr := cmod.CreateGroupField(&grpfield, Ch.DatabaseConnection)

		if grpfielderr != nil {

			log.Println(grpfielderr)

		}

	}

	return nil
}

func (ch BasicChannel) DeleteChannel(ChannelId, userid int) error {

	if ChannelId <= 0 {

		return errors.New("invalid channelid cannot delete")
	}

	chid := strconv.Itoa(ChannelId)

	cmod.DeleteEntryByChannelId(ChannelId, ch.DatabaseConnection)

	cmod.DeleteChannelById(ChannelId, ch.DatabaseConnection)

	var chdel TblChannel

	cmod.GetChannelById(&chdel, ChannelId, ch.DatabaseConnection)

	var delfidgrp TblFieldGroup

	delfidgrp.IsDeleted = 1

	delfidgrp.DeletedBy = userid

	delfidgrp.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	cmod.DeleteFieldGroupById(&delfidgrp, chdel.FieldGroupId, ch.DatabaseConnection)

	var checkid auth.TblModulePermission

	auth.AS.GetIdByRouteName(&checkid, chid, ch.DatabaseConnection)

	var DeleteRolepermission auth.TblRolePermission

	auth.AS.Deleterolepermission(&DeleteRolepermission, checkid.Id, ch.DatabaseConnection)

	var modpermission auth.TblModulePermission

	auth.AS.DeleteModulePermissioninEntries(&modpermission, chid, ch.DatabaseConnection)

	//business logic
	return nil
}

func (ch BasicChannel) GetChannels(limit, offset int, filter Filter, activestatus bool) (channelList []TblChannel, channelcount int, err error) {

	var channellist []TblChannel

	cmod.Channellist(&channellist, limit, offset, filter, activestatus, ch.DatabaseConnection)

	var chnallist []TblChannel

	for _, val := range channellist {

		val.SlugName = val.ChannelDescription

		val.ChannelDescription = TruncateDescription(val.ChannelDescription, 130)

		if !val.ModifiedOn.IsZero() {

			val.DateString = val.ModifiedOn.Format("02 Jan 2006 03:04 PM")

		} else {

			val.DateString = val.CreatedOn.Format("02 Jan 2006 03:04 PM")

		}

		entrcount, _ := emod.ChannelEntryList(&[]TblChannelEntries{}, 0, 0, val.Id, EntriesFilter{}, false, false, ch.DatabaseConnection)

		val.EntriesCount = int(entrcount)

		chnallist = append(chnallist, val)

	}

	var chncount []TblChannel

	chcount, _ := cmod.Channellist(&chncount, 0, 0, filter, activestatus, ch.DatabaseConnection)

	return chnallist, int(chcount), nil
}

func (ch BasicChannel) GetchannelByName(channelname string) (channel TblChannel, err error) {

	var channellist TblChannel

	err1 := cmod.GetChannelByChannelName(&channellist, channelname, ch.DatabaseConnection)

	if err1 != nil {

		return TblChannel{}, err1
	}

	return channellist, nil
}

func (ch BasicChannel) GetChannelsById(channelid int) (channels TblChannel, section []Section, fields []Fiedlvalue, SelectedCategories []categories.Arrangecategories, err error) {

	var channellist TblChannel

	cmod.GetChannelById(&channellist, channelid, ch.DatabaseConnection)

	var groupfield []TblGroupField

	cmod.GetFieldIdByGroupId(&groupfield, channellist.Id, ch.DatabaseConnection)

	var ids []int

	for _, val := range groupfield {

		ids = append(ids, val.FieldId)
	}

	var fieldValue []TblField

	cmod.GetFieldAndOptionValue(&fieldValue, ids, ch.DatabaseConnection)

	var sections []Section

	var Fieldvalue []Fiedlvalue

	for _, val := range fieldValue {

		var section Section

		var fieldvalue Fiedlvalue

		if val.FieldTypeId == 12 {

			section.SectionId = val.Id

			section.SectionNewId = 0

			section.MasterFieldId = val.FieldTypeId

			section.SectionName = val.FieldName

			sections = append(sections, section)

		} else {

			var optionva []OptionValues

			for _, optionval := range val.TblFieldOption {

				var optiovalue OptionValues

				optiovalue.Id = optionval.Id

				optiovalue.FieldId = optionval.FieldId

				optiovalue.Value = optionval.OptionValue

				optionva = append(optionva, optiovalue)
			}

			fieldvalue.FieldId = val.Id

			fieldvalue.FieldName = val.FieldName

			fieldvalue.CharacterAllowed = val.CharacterAllowed

			fieldvalue.DateFormat = val.DatetimeFormat

			fieldvalue.TimeFormat = val.TimeFormat

			fieldvalue.IconPath = val.ImagePath

			fieldvalue.MasterFieldId = val.FieldTypeId

			fieldvalue.Mandatory = val.MandatoryField

			fieldvalue.OrderIndex = val.OrderIndex

			fieldvalue.SectionId = val.SectionParentId

			fieldvalue.OptionValue = optionva

			Fieldvalue = append(Fieldvalue, fieldvalue)

		}

	}

	var GetSelectedChannelCateogry []TblChannelCategory

	err1 := cmod.GetSelectedCategoryChannelById(&GetSelectedChannelCateogry, channelid, ch.DatabaseConnection)

	if err1 != nil {

		log.Println(err)
	}

	var FinalSelectedCategories []categories.Arrangecategories

	for _, val := range GetSelectedChannelCateogry {

		var id []int

		ids := strings.Split(val.CategoryId, ",")

		for _, cid := range ids {

			convid, _ := strconv.Atoi(cid)

			id = append(id, convid)
		}

		var GetSelectedCategory []categories.TblCategory

		cmod.GetCategoriseById(&GetSelectedCategory, id, ch.DatabaseConnection)

		var addcat categories.Arrangecategories

		var individualid []categories.CatgoriesOrd

		for _, CategoriesArrange := range GetSelectedCategory {

			var individual categories.CatgoriesOrd

			individual.Id = CategoriesArrange.Id

			individual.Category = CategoriesArrange.CategoryName

			individualid = append(individualid, individual)

		}

		addcat.Categories = individualid

		FinalSelectedCategories = append(FinalSelectedCategories, addcat)

	}

	return channellist, sections, Fieldvalue, FinalSelectedCategories, nil
}

// if description is too big show specific lines and after show ...
func TruncateDescription(description string, limit int) string {
	if len(description) <= limit {
		return description
	}

	truncated := description[:limit] + "..."
	return truncated
}
