// Package Channel will help to create a channels in cms
package channels

import (
	"errors"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spurtcms/pkgcontent/categories"
	"github.com/spurtcms/pkgcore/auth"
	authcore "github.com/spurtcms/pkgcore/auth"
	"github.com/spurtcms/pkgcore/member"
	"gorm.io/gorm"
)

func MigrateTable(db *gorm.DB) {

	db.AutoMigrate(
		&TblFieldType{},
		&TblChannel{},
		&TblField{},
		&TblFieldGroup{},
		&TblFieldOption{},
		&TblGroupField{},
		&TblChannelCategory{},
	)
}

/*this struct holds dbconnection ,token*/
type Channel struct {
	Authority *authcore.Authorization
}

type ChannelStruct struct{}

var (
	CH ChannelStruct

	SpecialToken = "%$HEID$#PDGH*&MGEAFCC"
)

/*Get AllChannels*/
func (Ch Channel) GetChannels(limit, offset int, filter Filter, activestatus bool) (channelList []TblChannel, channelcount int, err error) {

	_, roleid, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return []TblChannel{}, 0, checkerr
	}

	check, err := Ch.Authority.IsGranted("Channel", authcore.CRUD)

	if err != nil {

		return []TblChannel{}, 0, err
	}

	if check {

		var channellist []TblChannel

		CH.Channellist(&channellist, limit, offset, filter, activestatus, roleid, false, Ch.Authority.DB)

		var chnallist []TblChannel

		for _, val := range channellist {

			val.SlugName = val.ChannelDescription

			val.ChannelDescription = TruncateDescription(val.ChannelDescription, 130)

			if !val.ModifiedOn.IsZero() {

				val.DateString = val.ModifiedOn.Format("02 Jan 2006 03:04 PM")

			} else {

				val.DateString = val.CreatedOn.Format("02 Jan 2006 03:04 PM")

			}

			entrcount, _ := CH.ChannelEntryList(&[]TblChannelEntries{}, 0, 0, val.Id, EntriesFilter{}, false, roleid, false, Ch.Authority.DB)

			val.EntriesCount = int(entrcount)

			chnallist = append(chnallist, val)

		}

		var chncount []TblChannel

		chcount, _ := CH.Channellist(&chncount, 0, 0, filter, activestatus, roleid, false, Ch.Authority.DB)

		return chnallist, int(chcount), nil

	}

	return []TblChannel{}, 0, errors.New("not authorized")
}

/*Get AllChannels*/
func (Ch Channel) GetPermissionChannels(limit, offset int, filter Filter, activestatus bool) (channelList []TblChannel, channelcount int, err error) {

	_, roleid, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return []TblChannel{}, 0, checkerr
	}

	// check, err := Ch.Authority.IsGranted("Channel", authcore.CRUD)

	// if err != nil {

	// 	return []TblChannel{}, 0, err
	// }

	// if check {

	var channellist []TblChannel

	CH.Channellist(&channellist, limit, offset, filter, activestatus, roleid, true, Ch.Authority.DB)

	var chnallist []TblChannel

	for _, val := range channellist {

		val.SlugName = val.ChannelDescription

		val.ChannelDescription = TruncateDescription(val.ChannelDescription, 130)

		if !val.ModifiedOn.IsZero() {

			val.DateString = val.ModifiedOn.Format("02 Jan 2006 03:04 PM")

		} else {

			val.DateString = val.CreatedOn.Format("02 Jan 2006 03:04 PM")

		}

		entrcount, _ := CH.ChannelEntryList(&[]TblChannelEntries{}, 0, 0, val.Id, EntriesFilter{}, false, roleid, false, Ch.Authority.DB)

		val.EntriesCount = int(entrcount)

		chnallist = append(chnallist, val)

	}

	var chncount []TblChannel

	chcount, _ := CH.Channellist(&chncount, 0, 0, filter, activestatus, roleid, true, Ch.Authority.DB)

	return chnallist, int(chcount), nil

	// }

	// return []TblChannel{}, 0, errors.New("not authorized")
}

/*Get channel by name*/
func (Ch Channel) GetchannelByName(channelname string) (channel TblChannel, err error) {

	_, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return TblChannel{}, checkerr
	}

	var channellist TblChannel

	err1 := CH.GetChannelByChannelName(&channellist, channelname, Ch.Authority.DB)

	if err1 != nil {

		return TblChannel{}, err1
	}

	return channellist, nil

}

/*Get Channels By Id*/
func (Ch Channel) GetChannelsById(channelid int) (channelList TblChannel, section []Section, fields []Fiedlvalue, SelectedCategories []categories.Arrangecategories, err error) {

	_, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return TblChannel{}, []Section{}, []Fiedlvalue{}, []categories.Arrangecategories{}, checkerr
	}

	// check, err := Ch.Authority.IsGranted("Channel", authcore.CRUD)

	// if err != nil {

	// 	return TblChannel{}, []Section{}, []Fiedlvalue{}, []categories.Arrangecategories{}, err
	// }

	// if check {

	var channellist TblChannel

	CH.GetChannelById(&channellist, channelid, Ch.Authority.DB)

	var groupfield []TblGroupField

	CH.GetFieldIdByGroupId(&groupfield, channellist.Id, Ch.Authority.DB)

	var ids []int

	for _, val := range groupfield {

		ids = append(ids, val.FieldId)
	}

	var fieldValue []TblField

	CH.GetFieldAndOptionValue(&fieldValue, ids, Ch.Authority.DB)

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

	err1 := CH.GetSelectedCategoryChannelById(&GetSelectedChannelCateogry, channelid, Ch.Authority.DB)

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

		CH.GetCategoriseById(&GetSelectedCategory, id, Ch.Authority.DB)

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
	// }

	// return TblChannel{}, []Section{}, []Fiedlvalue{}, []categories.Arrangecategories{}, errors.New("not authorized")
}

/*Create Channel*/
func (Ch Channel) CreateChannel(channelcreate ChannelCreate) (err error) {

	userid, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := Ch.Authority.IsGranted("Channel", authcore.CRUD)

	if err != nil {

		return err
	}

	if check {

		if channelcreate.ChannelName == "" || channelcreate.ChannelDescription == "" {

			return errors.New("empty value")

		}

		/*create channel*/
		var channel TblChannel

		channel.ChannelName = channelcreate.ChannelName

		channel.ChannelDescription = channelcreate.ChannelDescription

		channel.SlugName = strings.ToLower(strings.ReplaceAll(channelcreate.ChannelName, " ", " "))

		channel.IsActive = 1

		channel.CreatedBy = userid

		channel.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		ch, chanerr := CH.CreateChannel(&channel, Ch.Authority.DB)

		if chanerr != nil {

			log.Println(chanerr)

			return
		}

		/*This is for module permission creation*/
		var modperms auth.TblModulePermission

		modperms.DisplayName = ch.ChannelName

		modperms.RouteName = "/channel/entrylist/" + strconv.Itoa(ch.Id)

		modperms.SlugName = strings.ReplaceAll(strings.ToLower(ch.ChannelName), " ", "_")

		modperms.CreatedBy = userid

		modperms.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		modperms.ModuleId = 8

		modperms.FullAccessPermission = 1

		modid, _ := auth.AS.CreateModulePermission(&modperms, Ch.Authority.DB)

		var tblrole auth.TblRolePermission

		tblrole.RoleId = 1

		tblrole.PermissionId = modid.Id

		tblrole.CreatedBy = userid

		tblrole.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		auth.AS.CreateRolePermissionsingle(&tblrole, Ch.Authority.DB)

		/*Temp store section id*/
		type tempsection struct {
			Id           int
			SectionId    int
			NewSectionId int
		}

		var TempSections []tempsection

		/*create Section*/
		for _, sectionvalue := range channelcreate.Sections {

			var cfld TblField

			cfld.FieldName = strings.TrimSpace(sectionvalue.SectionName)

			cfld.FieldTypeId = sectionvalue.MasterFieldId

			cfld.CreatedBy = userid

			cfld.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			cfid, fiderr := CH.CreateFields(&cfld, Ch.Authority.DB)

			if fiderr != nil {

				log.Println(fiderr)
			}

			/*create group field*/
			var grpfield TblGroupField

			grpfield.ChannelId = ch.Id

			grpfield.FieldId = cfid.Id

			grpfielderr := CH.CreateGroupField(&grpfield, Ch.Authority.DB)

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
		for _, val := range channelcreate.FieldValues {

			var cfld TblField

			cfld.FieldName = strings.TrimSpace(val.FieldName)

			cfld.FieldTypeId = val.MasterFieldId

			cfld.CreatedBy = userid

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

			cfid, fiderr := CH.CreateFields(&cfld, Ch.Authority.DB)

			if fiderr != nil {

				log.Println(fiderr)

			}

			/*option value create*/
			for _, opt := range val.OptionValue {

				var fldopt TblFieldOption

				fldopt.OptionName = opt.Value

				fldopt.OptionValue = opt.Value

				fldopt.FieldId = cfid.Id

				fldopt.CreatedBy = userid

				fldopt.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				fopterr := CH.CreateFieldOption(&fldopt, Ch.Authority.DB)

				if fopterr != nil {

					log.Println(fopterr)

				}

			}

			/*create group field*/
			var grpfield TblGroupField

			grpfield.ChannelId = ch.Id

			grpfield.FieldId = cfid.Id

			grpfielderr := CH.CreateGroupField(&grpfield, Ch.Authority.DB)

			if grpfielderr != nil {

				log.Println(grpfielderr)

			}

		}

		for _, categoriesid := range channelcreate.CategoryIds {

			var channelcategory TblChannelCategory

			channelcategory.ChannelId = ch.Id

			channelcategory.CategoryId = categoriesid

			channelcategory.CreatedAt = userid

			channelcategory.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			err := CH.CreateChannelCategory(&channelcategory, Ch.Authority.DB)

			if err != nil {

				log.Println(err)

			}

		}

		return nil

	}

	return errors.New("not authorized")
}

/*Edit channel*/
func (Ch Channel) EditChannel(channelupt ChannelUpdate, channelid int) error {

	userid, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := Ch.Authority.IsGranted("Channel", authcore.CRUD)

	if err != nil {

		return err
	}

	if check {

		var chn TblChannel

		chn.ChannelName = channelupt.ChannelName

		chn.ChannelDescription = channelupt.ChannelDescription

		chn.ModifiedBy = userid

		chn.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		CH.UpdateChannelDetails(&chn, channelid, Ch.Authority.DB)

		var modpermissionupdate auth.TblModulePermission

		modpermissionupdate.SlugName = channelupt.ChannelName

		modpermissionupdate.RouteName = "/channel/entrylist/" + strconv.Itoa(channelid)

		modpermissionupdate.DisplayName = channelupt.ChannelName

		auth.AS.UpdateChannelNameInEntries(&modpermissionupdate, Ch.Authority.DB)

		//delete sections & fields
		var delid []int //temp array for delid
		var optiondelid []int

		for _, val := range channelupt.Deletesections {

			delid = append(delid, val.SectionId)
		}

		for _, val := range channelupt.DeleteFields {

			delid = append(delid, val.FieldId)
		}

		for _, val := range channelupt.DeleteOptionsvalue {

			optiondelid = append(optiondelid, val.Id)

		}

		if len(delid) > 0 || len(optiondelid) > 0 {

			var delsection TblField

			delsection.DeletedBy = userid

			delsection.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			delsection.IsDeleted = 1

			CH.DeleteFieldById(&delsection, delid, Ch.Authority.DB)

			var deloption TblFieldOption

			deloption.DeletedBy = userid

			deloption.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			deloption.IsDeleted = 1

			CH.DeleteOptionById(&deloption, optiondelid, delid, Ch.Authority.DB)

		}

		/*Temp store section id*/
		type tempsection struct {
			Id           int
			SectionId    int
			NewSectionId int
		}

		var TempSections []tempsection

		for _, val := range channelupt.Sections {

			var cfld TblField

			cfld.FieldName = strings.TrimSpace(val.SectionName)

			cfld.FieldTypeId = val.MasterFieldId

			cfld.CreatedBy = userid

			cfld.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			if val.SectionId != 0 {

				CH.UpdateFieldDetails(&cfld, val.SectionId, Ch.Authority.DB)

				var TempSection tempsection

				TempSection.Id = val.SectionId

				TempSection.SectionId = val.SectionId

				TempSection.NewSectionId = val.SectionNewId

				TempSections = append(TempSections, TempSection)

			} else {

				cfid, fiderr := CH.CreateFields(&cfld, Ch.Authority.DB)

				if fiderr != nil {

					log.Println(fiderr)
				}

				/*create group field*/
				var grpfield TblGroupField

				grpfield.ChannelId = channelid

				grpfield.FieldId = cfid.Id

				grpfielderr := CH.CreateGroupField(&grpfield, Ch.Authority.DB)

				if grpfielderr != nil {

					log.Println(grpfielderr)

				}

				var TempSection tempsection

				TempSection.Id = cfid.Id

				TempSection.SectionId = val.SectionId

				TempSection.NewSectionId = val.SectionNewId

				TempSections = append(TempSections, TempSection)

			}

		}

		for _, val := range channelupt.FieldValues {

			var cfld TblField

			cfld.FieldName = strings.TrimSpace(val.FieldName)

			cfld.FieldTypeId = val.MasterFieldId

			cfld.CreatedBy = userid

			cfld.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			cfld.OrderIndex = val.OrderIndex

			cfld.ImagePath = val.IconPath

			cfld.MandatoryField = val.Mandatory

			cfld.Url = val.Url

			cfld.CharacterAllowed = val.CharacterAllowed

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

			var createdchannelid int

			if val.FieldId != 0 {

				cfld.ModifiedBy = userid

				cfld.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				CH.UpdateFieldDetails(&cfld, val.FieldId, Ch.Authority.DB)

				createdchannelid = val.FieldId

			} else {

				cfid, fiderr := CH.CreateFields(&cfld, Ch.Authority.DB)

				if fiderr != nil {

					log.Println(fiderr)

				}

				/*create group field*/
				var grpfield TblGroupField

				grpfield.ChannelId = channelid

				grpfield.FieldId = cfid.Id

				grpfielderr := CH.CreateGroupField(&grpfield, Ch.Authority.DB)

				if grpfielderr != nil {

					log.Println(grpfielderr)

				}

				createdchannelid = cfld.Id

			}
			for _, optv := range val.OptionValue {

				var fldopt TblFieldOption

				fldopt.OptionName = optv.Value

				fldopt.OptionValue = optv.Value

				fldopt.CreatedBy = userid

				fldopt.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				if optv.Id != 0 {

					fldopt.FieldId = optv.FieldId

					fldopt.ModifiedBy = userid

					fldopt.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					CH.UpdateFieldOption(&fldopt, optv.Id, Ch.Authority.DB)

				} else {

					fldopt.FieldId = createdchannelid

					fopterr := CH.CreateFieldOption(&fldopt, Ch.Authority.DB)

					if fopterr != nil {

						log.Println(fopterr)

					}

				}

			}

		}

		/*channel category create if not exist*/
		for _, val := range channelupt.CategoryIds {

			var channcategory TblChannelCategory

			err := CH.CheckChannelCategoryAlreadyExitst(&channcategory, channelid, val, Ch.Authority.DB)

			if errors.Is(err, gorm.ErrRecordNotFound) {

				var createCateogry TblChannelCategory

				createCateogry.ChannelId = channelid

				createCateogry.CategoryId = val

				createCateogry.CreatedAt = userid

				createCateogry.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				CH.CreateChannelCategory(&createCateogry, Ch.Authority.DB)
			}

		}

		/*delete categoryid if not exist in array*/
		var notexistcategory []TblChannelCategory

		CH.GetChannelCategoryNotExist(&notexistcategory, channelid, channelupt.CategoryIds, Ch.Authority.DB)

		for _, val := range notexistcategory {

			var deletechannelcategory TblChannelCategory

			CH.DeleteChannelCategoryByValue(&deletechannelcategory, val.Id, Ch.Authority.DB)

		}

	}

	return errors.New("not authorized")

}

/*Delete Channel*/
func (Ch Channel) DeleteChannel(channelid int) error {

	userid, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := Ch.Authority.IsGranted("Channel", authcore.CRUD)

	if err != nil {

		return err
	}

	if check {

		if channelid <= 0 {

			return errors.New("invalid channelid cannot delete")
		}

		chid := strconv.Itoa(channelid)

		CH.DeleteEntryByChannelId(channelid, Ch.Authority.DB)

		CH.DeleteChannelById(channelid, Ch.Authority.DB)

		var chdel TblChannel

		CH.GetChannelById(&chdel, channelid, Ch.Authority.DB)

		var delfidgrp TblFieldGroup

		delfidgrp.IsDeleted = 1

		delfidgrp.DeletedBy = userid

		delfidgrp.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		CH.DeleteFieldGroupById(&delfidgrp, chdel.FieldGroupId, Ch.Authority.DB)

		var checkid auth.TblModulePermission

		auth.AS.GetIdByRouteName(&checkid, chid, Ch.Authority.DB)

		var DeleteRolepermission auth.TblRolePermission

		auth.AS.Deleterolepermission(&DeleteRolepermission, checkid.Id, Ch.Authority.DB)

		var modpermission auth.TblModulePermission

		auth.AS.DeleteModulePermissioninEntries(&modpermission, chid, Ch.Authority.DB)

		return nil
	}

	return errors.New("not authorized")
}

/*Change Channel status*/
// status 0 = inactive
// status 1 = active
func (Ch Channel) ChangeChannelStatus(channelid int, status int) (bool, error) {

	userid, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	check, err := Ch.Authority.IsGranted("Channel", authcore.CRUD)

	if err != nil {

		return false, err
	}

	if check {

		if channelid <= 0 {

			return false, errors.New("invalid channelid cannot delete")
		}

		var channelstatus TblChannel

		channelstatus.ModifiedBy = userid

		channelstatus.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		CH.ChannelIsActive(&channelstatus, channelid, status, Ch.Authority.DB)

		return true, nil
	}

	return false, errors.New("not authorized")
}

/*Get All Master Field type */
func (Ch Channel) GetAllMasterFieldType() (field []TblFieldType, err error) {

	_, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return []TblFieldType{}, checkerr
	}

	check, err := Ch.Authority.IsGranted("Channel", authcore.CRUD)

	if err != nil {

		return []TblFieldType{}, err
	}

	if check {

		var field []TblFieldType

		CH.GetAllField(&field, Ch.Authority.DB)

		return field, nil
	}

	return []TblFieldType{}, errors.New("not authorized")
}

/*all channel Entries List*/
//if channelid 0 get all channel entries
// if channelid not eq 0 to get particular entries of the channel
func (Ch Channel) GetAllChannelEntriesList(channelid int, limit, offset int, filter EntriesFilter) (entries []TblChannelEntries, filterentriescount int, totalentriescount int, err error) {

	_, roleid, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return []TblChannelEntries{}, 0, 0, checkerr
	}

	if filter.Status == "Draft" {

		filter.Status = "0"

	} else if filter.Status == "Published" {

		filter.Status = "1"

	} else if filter.Status == "Unpublished" {

		filter.Status = "2"
	}
	var chnentry []TblChannelEntries

	CH.ChannelEntryList(&chnentry, limit, offset, channelid, filter, false, roleid, true, Ch.Authority.DB)

	var chnentry1 []TblChannelEntries

	filtercount, _ := CH.ChannelEntryList(&chnentry1, 0, 0, channelid, filter, false, roleid, true, Ch.Authority.DB)

	entrcount, _ := CH.ChannelEntryList(&chnentry1, 0, 0, channelid, EntriesFilter{}, false, roleid, true, Ch.Authority.DB)

	return chnentry, int(filtercount), int(entrcount), nil

}

// Get published entries
func (Ch Channel) GetPublishedChannelEntriesList(limit, offset int, filter EntriesFilter) (entries []TblChannelEntries, filterentriescount int, totalentriescount int, err error) {

	_, roleid, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return []TblChannelEntries{}, 0, 0, checkerr
	}

	var chnentry []TblChannelEntries

	CH.ChannelEntryList(&chnentry, limit, offset, 0, filter, true, roleid, true, Ch.Authority.DB)

	filtercount, _ := CH.ChannelEntryList(&chnentry, 0, 0, 0, filter, true, roleid, true, Ch.Authority.DB)

	var chnentry1 []TblChannelEntries

	entrcount, _ := CH.ChannelEntryList(&chnentry1, 0, 0, 0, EntriesFilter{}, true, roleid, true, Ch.Authority.DB)

	return chnentry, int(filtercount), int(entrcount), nil

}

// create entry
func (Ch Channel) CreateEntry(entriesrequired EntriesRequired) (entry TblChannelEntries, flg bool, err error) {

	userid, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return TblChannelEntries{}, false, checkerr
	}

	check, err := Ch.Authority.IsGranted(entriesrequired.ChannelName, authcore.CRUD)

	if err != nil {

		return TblChannelEntries{}, false, err
	}

	if check {

		var Entries TblChannelEntries

		Entries.Title = entriesrequired.Title

		Entries.Description = entriesrequired.Content

		Entries.CoverImage = entriesrequired.CoverImage

		Entries.MetaTitle = entriesrequired.SEODetails.MetaTitle

		Entries.MetaDescription = entriesrequired.SEODetails.MetaDescription

		Entries.Keyword = entriesrequired.SEODetails.MetaKeywords

		Entries.Slug = entriesrequired.SEODetails.MetaSlug

		Entries.Status = entriesrequired.Status

		Entries.ChannelId = entriesrequired.ChannelId

		Entries.CategoriesId = entriesrequired.CategoryIds

		Entries.CreatedBy = userid

		Entries.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		Entriess, err := CH.CreateChannelEntry(&Entries, Ch.Authority.DB)

		if err != nil {

			log.Println(err)
		}

		if len(entriesrequired.AdditionalFields) > 0 {

			var EntriesField []TblChannelEntryField

			for _, val := range entriesrequired.AdditionalFields {

				var Entrfield TblChannelEntryField

				Entrfield.ChannelEntryId = Entriess.Id

				Entrfield.FieldName = val.FieldName

				Entrfield.FieldValue = val.FieldValue

				Entrfield.FieldId = val.FieldId

				Entrfield.CreatedBy = userid

				Entrfield.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				EntriesField = append(EntriesField, Entrfield)

			}

			ferr := CH.CreateEntrychannelFields(&EntriesField, Ch.Authority.DB)

			if ferr != nil {

				log.Println(ferr)
			}
		}

		return Entries, true, nil

	}

	return TblChannelEntries{}, false, errors.New("not authorized")
}

/**/
func (Ch Channel) DeleteEntry(ChannelName string, Entryid int) (bool, error) {

	userid, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	check, err := Ch.Authority.IsGranted(ChannelName, authcore.CRUD)

	if err != nil {

		return false, err
	}

	if check {

		var entries TblChannelEntries

		entries.Id = Entryid

		entries.IsDeleted = 1

		entries.DeletedBy = userid

		entries.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		err := CH.DeleteChannelEntryId(&entries, Entryid, Ch.Authority.DB)

		var field TblChannelEntryField

		field.DeletedBy = userid

		field.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		err1 := CH.DeleteChannelEntryFieldId(&field, Entryid, Ch.Authority.DB)

		if err != nil {

			log.Println(err)
		}

		if err1 != nil {

			log.Println(err)
		}

		return true, nil
	}

	return false, errors.New("not authorized")
}

/**/
func (Ch Channel) GetAdditionalFieldDataBychannelId(ChannelName string, EntryId int) ([]TblChannelEntryField, error) {

	_, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return []TblChannelEntryField{}, checkerr
	}

	check, err := Ch.Authority.IsGranted(ChannelName, authcore.CRUD)

	if err != nil {

		return []TblChannelEntryField{}, err
	}

	if check {

		var EntriesField []TblChannelEntryField

		err := CH.GetChannelEntryDetailsById(&EntriesField, EntryId, Ch.Authority.DB)

		if err != nil {

			log.Println(err)
		}

		return EntriesField, nil
	}

	return []TblChannelEntryField{}, errors.New("not authorized")
}

// get entry details
func (Ch Channel) GetEntryDetailsById(ChannelName string, EntryId int) (TblChannelEntries, error) {

	_, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return TblChannelEntries{}, checkerr
	}

	check, err := Ch.Authority.IsGranted(ChannelName, authcore.CRUD)

	if err != nil {

		return TblChannelEntries{}, err
	}

	if check {

		var Entry TblChannelEntries

		err := CH.GetChannelEntryById(&Entry, EntryId, Ch.Authority.DB)

		if err != nil {

			log.Println(err)
		}

		return Entry, nil
	}
	return TblChannelEntries{}, errors.New("not authorized")
}

/*update entry details */
func (Ch Channel) UpdateEntryDetailsById(entriesrequired EntriesRequired, ChannelName string, EntryId int) (bool, error) {

	userid, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	check, err := Ch.Authority.IsGranted(ChannelName, authcore.CRUD)

	if err != nil {

		return false, err
	}

	if check {

		var Entries TblChannelEntries

		Entries.Title = entriesrequired.Title

		Entries.Description = entriesrequired.Content

		Entries.CoverImage = entriesrequired.CoverImage

		Entries.MetaTitle = entriesrequired.SEODetails.MetaTitle

		Entries.MetaDescription = entriesrequired.SEODetails.MetaDescription

		Entries.Keyword = entriesrequired.SEODetails.MetaKeywords

		Entries.Slug = entriesrequired.SEODetails.MetaSlug

		Entries.Status = entriesrequired.Status

		Entries.ChannelId = entriesrequired.ChannelId

		Entries.CategoriesId = entriesrequired.CategoryIds

		Entries.ModifiedBy = userid

		Entries.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		err := CH.UpdateChannelEntryDetails(&Entries, EntryId, Ch.Authority.DB)

		if err != nil {

			log.Println(err)
		}

		for _, val := range entriesrequired.AdditionalFields {

			if val.Id == 0 {

				var Entrfield TblChannelEntryField

				Entrfield.ChannelEntryId = EntryId

				Entrfield.FieldName = val.FieldName

				Entrfield.FieldValue = val.FieldValue

				Entrfield.FieldId = val.FieldId

				Entrfield.CreatedBy = userid

				Entrfield.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				CH.CreateSingleEntrychannelFields(&Entrfield, Ch.Authority.DB)

			} else {

				var Entrfield TblChannelEntryField

				Entrfield.Id = val.Id

				Entrfield.ChannelEntryId = EntryId

				Entrfield.FieldName = val.FieldName

				Entrfield.FieldValue = val.FieldValue

				Entrfield.FieldId = val.FieldId

				Entrfield.ModifiedBy = userid

				Entrfield.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				CH.UpdateChannelEntryAdditionalDetails(Entrfield, *Ch.Authority.DB)

			}

		}

		return true, nil

	}

	return false, errors.New("not authorized")

}

// change entries status
func (Ch Channel) EntryStatus(ChannelName string, EntryId int, status int) (bool, error) {

	userid, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	check, err := Ch.Authority.IsGranted(ChannelName, authcore.CRUD)

	if err != nil {

		return false, err
	}

	if check {

		var Entries TblChannelEntries

		Entries.Status = status

		Entries.ModifiedBy = userid

		Entries.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		CH.PublishQuery(&Entries, EntryId, *Ch.Authority.DB)

		return true, nil

	}

	return false, errors.New("not authorized")

}

// if description is too big show specific lines and after show ...
func TruncateDescription(description string, limit int) string {
	if len(description) <= limit {
		return description
	}

	truncated := description[:limit] + "..."
	return truncated
}

// DashboardEntry count function
func (Ch Channel) DashboardEntriesCount() (totalcount int, lasttendayscount int, err error) {

	_, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return 0, 0, checkerr
	}

	allentrycount, err := CH.AllentryCount(Ch.Authority.DB)

	if err != nil {

		return 0, 0, err
	}

	entrycount, err := CH.NewentryCount(Ch.Authority.DB)

	if err != nil {

		return 0, 0, err
	}

	return int(allentrycount), int(entrycount), nil
}

func (Ch Channel) DashboardChannellist() (channelList []TblChannel, err error) {

	_, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return []TblChannel{}, checkerr
	}

	Newchannels, _ := CH.Newchannels(Ch.Authority.DB)

	if err != nil {

		return []TblChannel{}, checkerr

	}

	return Newchannels, nil

}

/*DashboardEntries */
func (Ch Channel) DashboardEntrieslist() (entries []TblChannelEntries, err error) {

	_, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return []TblChannelEntries{}, checkerr
	}

	Newentries, _ := CH.Newentries(Ch.Authority.DB)

	if err != nil {

		return []TblChannelEntries{}, checkerr

	}

	return Newentries, nil

}

/*Recent activites for dashboard*/
func (Ch Channel) DashboardRecentActivites() (entries []RecentActivities, err error) {

	_, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return []RecentActivities{}, checkerr
	}

	Newentries, _ := CH.Newentries(Ch.Authority.DB)

	var Newrecords []RecentActivities

	for _, val := range Newentries {

		newrecord := RecentActivities{Contenttype: "entry", Title: val.Title, User: val.Username, Imagepath: val.ProfileImagePath, Createdon: val.CreatedOn, Channelname: val.ChannelName}

		Newrecords = append(Newrecords, newrecord)
	}

	Newchannel, _ := CH.Newchannels(Ch.Authority.DB)

	for _, val := range Newchannel {

		newrecord := RecentActivities{Contenttype: "channel", Title: val.ChannelName, User: val.Username, Imagepath: val.ProfileImagePath, Createdon: val.CreatedOn, Channelname: val.ChannelName}

		Newrecords = append(Newrecords, newrecord)
	}
	sort.Slice(Newrecords, func(i, j int) bool {

		return Newrecords[i].Createdon.After(Newrecords[j].Createdon)

	})
	maxRec := 5

	if len(Newrecords) < maxRec {

		maxRec = len(Newrecords)

	}
	recentActive := Newrecords[:maxRec]

	var newactive RecentActivities

	var NewActive []RecentActivities

	for _, val := range recentActive {

		difference := time.Now().Sub(val.Createdon)

		hour := int(difference.Hours())

		min := int(difference.Minutes())

		if hour >= 1 {

			newactive.Contenttype = val.Contenttype

			newactive.Title = val.Title

			newactive.User = val.User

			newactive.Imagepath = val.Imagepath

			newactive.Createdon = val.Createdon

			newactive.Channelname = val.Channelname

			newactive.Active = strconv.Itoa(hour) + " " + "hrs"
		} else {
			newactive.Contenttype = val.Contenttype

			newactive.Title = val.Title

			newactive.User = val.User

			newactive.Imagepath = val.Imagepath

			newactive.Createdon = val.Createdon

			newactive.Channelname = val.Channelname

			newactive.Active = strconv.Itoa(min) + " " + "mins"

		}

		NewActive = append(NewActive, newactive)

	}

	return NewActive, nil
}

/*Remove entries cover image if media image delete*/
func (Ch Channel) RemoveEntriesCoverImage(ImagePath string) error {

	_, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	err := CH.UpdateImagePath(ImagePath, Ch.Authority.DB)

	if err != nil {

		return err
	}

	return nil

}

// Get published entries
func (Ch Channel) GetPublishedChannelEntriesListForTemplate(limit, offset int, filter EntriesFilter) (entries []TblChannelEntries, count int, err error) {

	var chnentry []TblChannelEntries

	CH.ChannelEntryListForTemplates(&chnentry, limit, offset, filter, Ch.Authority.DB)

	filtercount, _ := CH.ChannelEntryListForTemplates(&chnentry, 0, 0, filter, Ch.Authority.DB)

	return chnentry, int(filtercount), nil

}

// this function provides channel list for accessible members if the channel contains published entries
func (ch Channel) GetGraphqlChannelList(limit, offset int) (channelList []TblChannel, count int64, err error) {

	var memberid int

	if ch.Authority.Token == SpecialToken {

		memberid = 0

	} else {

		memberid, _, err = member.VerifyToken(ch.Authority.Token, ch.Authority.Secret)

		if err != nil {

			return []TblChannel{}, 0, err

		}

	}

	channelList, count, err = CH.GetGraphqlChannelList(ch.Authority.DB, memberid, limit, offset)

	if err != nil {

		return []TblChannel{}, 0, err

	}

	return channelList, count, nil
}

// given entry id returns related entry
func (ch Channel) GetGraphqlChannelEntriesDetails(channelEntryId, channelId, categoryId *int, pathUrl string) (channelEntry TblChannelEntries, err error) {

	var memberid int

	if ch.Authority.Token == SpecialToken {

		memberid = 0

	} else {

		memberid, _, err = member.VerifyToken(ch.Authority.Token, ch.Authority.Secret)

		if err != nil {

			return TblChannelEntries{}, err

		}

	}

	channelEntry, err = CH.GetGraphqlChannelEntryDetailsById(ch.Authority.DB, memberid, channelEntryId, channelId, categoryId)

	if err != nil {

		return TblChannelEntries{}, err
	}

	modified_path := strings.TrimPrefix(channelEntry.CoverImage, "/")

	channelEntry.CoverImage = pathUrl + modified_path

	splittedArr := strings.Split(channelEntry.CategoriesId, ",")

	var parentCatId int

	var indivCategories [][]categories.TblCategory

	for _, catId := range splittedArr {

		var indivCategory []categories.TblCategory

		conv_id, _ := strconv.Atoi(catId)

		category, _ := CH.GetGraphqlEntriesCategoryByParentId(ch.Authority.DB, conv_id)

		if category.Id != 0 {

			indivCategory = append(indivCategory, category)
		}

		parentCatId = category.ParentId

		if parentCatId != 0 {

			var count int

		LOOP:

			for {

				count = count + 1 //count increment used to check how many times the loop gets executed

				parentCategory, _ := CH.GetGraphqlEntriesCategoryByParentId(ch.Authority.DB, parentCatId)

				if parentCategory.Id != 0 {

					indivCategory = append(indivCategory, parentCategory)
				}

				parentCatId = parentCategory.ParentId

				if parentCatId != 0 { //mannuall condition to break the loop in overlooping situations

					goto LOOP

				} else if count > 49 {

					break //use to break the loop if infinite loop doesn't break ,So forcing the loop to break at overlooping conditions

				} else {

					break
				}

			}

		}

		if len(indivCategory) > 0 {

			sort.SliceStable(indivCategory, func(i, j int) bool {

				return indivCategory[i].Id < indivCategory[j].Id

			})

			indivCategories = append(indivCategories, indivCategory)
		}

	}

	channelEntry.Categories = indivCategories

	return channelEntry, nil
}

// given channel id returns all published entries
func (ch Channel) GetGraphqlChannelEntriesByChannelId(channelId, categoryId, limit, offset *int, pathUrl string) (channelEntries []TblChannelEntries, count int64, err error) {

	var memberid int

	if ch.Authority.Token == SpecialToken {

		memberid = 0

	} else {

		memberid, _, err = member.VerifyToken(ch.Authority.Token, ch.Authority.Secret)

		if err != nil {

			return []TblChannelEntries{}, 0, err

		}

	}

	channelEntries, count, err = CH.GetGraphqlChannelEntrieslistByChannelId(ch.Authority.DB, memberid, channelId, categoryId, limit, offset)

	if err != nil {

		return []TblChannelEntries{}, 0, err
	}

	var final_entries_list []TblChannelEntries

	for _, entry := range channelEntries {

		modified_path := strings.TrimPrefix(entry.CoverImage, "/")

		entry.CoverImage = pathUrl + modified_path

		splittedArr := strings.Split(entry.CategoriesId, ",")

		var parentCatId int

		var indivCategories [][]categories.TblCategory

		for _, catId := range splittedArr {

			var indivCategory []categories.TblCategory

			conv_id, _ := strconv.Atoi(catId)

			category, _ := CH.GetGraphqlEntriesCategoryByParentId(ch.Authority.DB, conv_id)

			if category.Id != 0 {

				indivCategory = append(indivCategory, category)
			}

			parentCatId = category.ParentId

			if parentCatId != 0 {

				var count int

			LOOP:

				for {

					count = count + 1 //count increment used to check how many times the loop gets executed

					parentCategory, _ := CH.GetGraphqlEntriesCategoryByParentId(ch.Authority.DB, parentCatId)

					if parentCategory.Id != 0 {

						indivCategory = append(indivCategory, parentCategory)
					}

					parentCatId = parentCategory.ParentId

					if parentCatId != 0 { //mannuall condition to break the loop in overlooping situations

						goto LOOP

					} else if count > 49 {

						break //use to break the loop if infinite loop doesn't break ,So forcing the loop to break at overlooping conditions

					} else {

						break
					}

				}

			}

			if len(indivCategory) > 0 {

				sort.SliceStable(indivCategory, func(i, j int) bool {

					return indivCategory[i].Id < indivCategory[j].Id

				})

				indivCategories = append(indivCategories, indivCategory)
			}

		}

		entry.Categories = indivCategories

		final_entries_list = append(final_entries_list, entry)
	}

	return final_entries_list, count, nil
}

// function give all channel entries list
func (ch Channel) GetGraphqlAllChannelEntriesList(categoryid, limit, offset *int, pathUrl string) (channelEntries []TblChannelEntries, count int64, err error) {

	var memberid int

	if ch.Authority.Token == SpecialToken {

		memberid = 0

	} else {

		memberid, _, err = member.VerifyToken(ch.Authority.Token, ch.Authority.Secret)

		if err != nil {

			return []TblChannelEntries{}, 0, err

		}

	}

	channelEntries, count, err = CH.GetGraphqlChannelEntriesList(ch.Authority.DB, memberid, categoryid, limit, offset)

	if err != nil {

		return []TblChannelEntries{}, 0, err
	}

	var final_entries_list []TblChannelEntries

	for _, entry := range channelEntries {

		modified_path := strings.TrimPrefix(entry.CoverImage, "/")

		entry.CoverImage = pathUrl + modified_path

		splittedArr := strings.Split(entry.CategoriesId, ",")

		var parentCatId int

		var indivCategories [][]categories.TblCategory

		for _, catId := range splittedArr {

			conv_id, _ := strconv.Atoi(catId)

			var indivCategory []categories.TblCategory

			category, _ := CH.GetGraphqlEntriesCategoryByParentId(ch.Authority.DB, conv_id)

			if category.Id != 0 {

				indivCategory = append(indivCategory, category)
			}

			parentCatId = category.ParentId

			if parentCatId != 0 {

				var count int

			LOOP:

				for {

					count = count + 1 //count increment used to check how many times the loop gets executed

					parentCategory, _ := CH.GetGraphqlEntriesCategoryByParentId(ch.Authority.DB, parentCatId)

					if parentCategory.Id != 0 {

						indivCategory = append(indivCategory, parentCategory)
					}

					parentCatId = parentCategory.ParentId

					if parentCatId != 0 { //mannuall condition to break the loop in overlooping situations

						goto LOOP

					} else if count > 49 {

						break //use to break the loop if infinite loop doesn't break ,So forcing the loop to break at overlooping conditions

					} else {

						break
					}

				}

			}

			if len(indivCategory) > 0 {

				sort.SliceStable(indivCategory, func(i, j int) bool {

					return indivCategory[i].Id < indivCategory[j].Id

				})

				indivCategories = append(indivCategories, indivCategory)
			}

		}

		entry.Categories = indivCategories

		final_entries_list = append(final_entries_list, entry)
	}

	return final_entries_list, count, nil
}

// this function provides channel detail for accessible members if the channel contains published entries
func (ch Channel) GetGraphqlChannelDetails(channelid int) (channel TblChannel, err error) {

	var memberid int

	if ch.Authority.Token == SpecialToken {

		memberid = 0

	} else {

		memberid, _, err = member.VerifyToken(ch.Authority.Token, ch.Authority.Secret)

		if err != nil {

			return TblChannel{}, err

		}

	}

	channel, err = CH.GetGraphqlChannelDetailsById(ch.Authority.DB, memberid, channelid)

	if err != nil {

		return TblChannel{}, err

	}

	return channel, nil
}

// get entry details
func (Ch Channel) GetEntryDetailsByIdTemplates(EntryId int) (TblChannelEntries, error) {

	var Entry TblChannelEntries

	err := CH.GetChannelEntryById(&Entry, EntryId, Ch.Authority.DB)

	if err != nil {

		log.Println(err)
	}

	return Entry, nil
}

func (Ch Channel) GetChannelCategoryByIdTemplates(channelid int) ([][]categories.CatgoriesOrd, error) {

	var category []TblChannelCategory

	err := CH.GetChannelCategoryDetailsByChannelId(&category, []int{channelid}, Ch.Authority.DB)

	if err != nil {

		log.Println(err)
	}

	var chancategory [][]categories.CatgoriesOrd

	for _, val := range category {

		ids := strings.Split(val.CategoryId, ",")

		var intid []int

		for _, val2 := range ids {

			i, _ := strconv.Atoi(val2)

			intid = append(intid, i)
		}

		var singcategory []categories.CatgoriesOrd

		var categoriss []categories.TblCategory

		categories.C.GetCategoryByIds(&categoriss, intid, Ch.Authority.DB)

		for _, val3 := range categoriss {

			var indsingcategory categories.CatgoriesOrd

			indsingcategory.Id = val3.Id

			indsingcategory.Category = val3.CategoryName

			singcategory = append(singcategory, indsingcategory)

		}

		chancategory = append(chancategory, singcategory)

	}

	return chancategory, nil
}
