package channels

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spurtcms/spurtcms-core/auth"
	authcore "github.com/spurtcms/spurtcms-core/auth"
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

type Channel struct {
	Authority *authcore.Authority
}

type ChannelStruct struct{}

var CH ChannelStruct

/*Get AllChannels*/
func (Ch Channel) GetChannels(Channels TblChannel, limit, offset int, filter Filter) (channelList []TblChannel, channelcount int, err error) {

	_, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return []TblChannel{}, 0, checkerr
	}

	check, err := Ch.Authority.IsGranted("Channels", authcore.CRUD)

	if err != nil {

		return []TblChannel{}, 0, err
	}

	if check {

		var channellist []TblChannel

		CH.Channellist(&channellist, limit, offset, filter, Ch.Authority.DB)

		var chnallist []TblChannel

		for _, val := range channellist {

			if !val.ModifiedOn.IsZero() {

				val.DateString = val.ModifiedOn.Format("02 Jan 2006 03:04 PM")

			} else {

				val.DateString = val.CreatedOn.Format("02 Jan 2006 03:04 PM")

			}

			chnallist = append(chnallist, val)

		}

		var chncount []TblChannel

		chcount, _ := CH.Channellist(&chncount, 0, 0, filter, Ch.Authority.DB)

		return chnallist, int(chcount), nil

	}

	return []TblChannel{}, 0, errors.New("not authorized")
}

/*Create Channel*/
func (Ch Channel) CreateChannel(channelcreate ChannelCreate) (err error) {

	userid, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := Ch.Authority.IsGranted("Channels", authcore.CRUD)

	if err != nil {

		return err
	}

	if check {

		type Fiedlvalue struct {
			Id         string   `json:"id"`
			Fid        string   `json:"fid"`
			Name       string   `json:"name"`
			Desc       string   `json:"desc"`
			Mandatory  int      `json:"mandatory"`
			Initial    string   `json:"initial"`
			Placehold  string   `json:"placehold"`
			DateFormat string   `json:"dateformat"`
			TimeFormat string   `json:"timeformat"`
			Optionname []string `json:"optioname"`
			ImgSrc     string   `json:"imgsrc"`
			Url        string   `json:"url"`
		}

		type field struct {
			Fiedlvalue []Fiedlvalue `json:"fiedlvalue"`
		}

		var fieldval field

		json.Unmarshal([]byte(channelcreate.FieldValues), &fieldval)

		if channelcreate.ChannelName == "" || channelcreate.ChannelDescription == "" {

			return errors.New("empty value")

		}

		/*create field group*/
		var cfldgroup TblFieldGroup

		cfldgroup.GroupName = channelcreate.ChannelName

		cfldgroup.CreatedBy = userid

		cfldgroup.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		cfgid, createfieldgrouperr := CH.CreateFieldGroup(&cfldgroup, Ch.Authority.DB)

		if createfieldgrouperr != nil {

			log.Println(createfieldgrouperr)

			// json.NewEncoder(c.Writer).Encode(false)

			return
		}

		/*create channel*/
		var channel TblChannel

		channel.ChannelName = channelcreate.ChannelName

		channel.ChannelDescription = channelcreate.ChannelDescription

		channel.SlugName = strings.ToLower(strings.ReplaceAll(channelcreate.ChannelName, " ", " "))

		channel.FieldGroupId = cfgid.Id

		channel.IsActive = 1

		channel.CreatedBy = userid

		channel.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		ch, chanerr := CH.CreateChannel(&channel, Ch.Authority.DB)

		if chanerr != nil {

			log.Println(chanerr)
		}

		for _, categoriesid := range channelcreate.CategoryIds {

			var channelcategory TblChannelCategory

			channelcategory.ChannelId = ch.Id

			channelcategory.CategoryId = categoriesid

			channelcategory.CreatedAt = userid

			channelcategory.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			err1 := CH.CreateChannelCategory(&channelcategory, Ch.Authority.DB)

			if err1 != nil {

				log.Println(err)

			}

		}

		/*This is for module permission creation*/
		var modperms auth.TblModulePermission

		modperms.DisplayName = ch.ChannelName

		modperms.RouteName = "/channel/entrylist/" + strconv.Itoa(ch.Id)

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

		/*create field*/
		for _, val := range fieldval.Fiedlvalue {

			var cfld TblField

			cfld.FieldName = strings.TrimSpace(val.Name)

			cfld.FieldDesc = val.Desc

			cfld.FieldTypeId, _ = strconv.Atoi(val.Fid)

			cfld.CreatedBy = userid

			cfld.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			cfld.MandatoryField = val.Mandatory

			cfld.OrderIndex, _ = strconv.Atoi(val.Id)

			cfld.InitialValue = val.Initial

			cfld.Placeholder = val.Placehold

			cfld.ImagePath = val.ImgSrc

			cfld.Url = val.Url

			if val.Fid == "4" {

				cfld.DatetimeFormat = val.DateFormat

				cfld.TimeFormat = val.TimeFormat

			}
			if val.Fid == "6" {

				cfld.DatetimeFormat = val.DateFormat
			}

			if len(val.Optionname) > 0 {

				cfld.OptionExist = 1
			}

			cfid, fiderr := CH.CreateFields(&cfld, Ch.Authority.DB)

			if fiderr != nil {

				log.Println(fiderr)
			}

			/*option value create*/
			for _, opt := range val.Optionname {

				var fldopt TblFieldOption

				fldopt.OptionName = opt

				fldopt.OptionValue = opt

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

			grpfield.FieldGroupId = cfgid.Id

			grpfield.FieldId = cfid.Id

			grpfielderr := CH.CreateGroupField(&grpfield, Ch.Authority.DB)

			if grpfielderr != nil {

				log.Println(grpfielderr)

			}

		}

	}

	return errors.New("not authorized")
}

/*Edit channel*/
func (Ch Channel) EditChannel(channelupt ChannelCreate, channelid int) error {

	userid, _, checkerr := authcore.VerifyToken(Ch.Authority.Token, Ch.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := Ch.Authority.IsGranted("Channels", authcore.CRUD)

	if err != nil {

		return err
	}

	if check {

		type Fiedlvalue struct {
			Id         string   `json:"id"`
			Fid        string   `json:"fid"`
			FieldId    int      `json:"FieldId"`
			Name       string   `json:"name"`
			Desc       string   `json:"desc"`
			Mandatory  int      `json:"mandatory"`
			Initial    string   `json:"initial"`
			Placehold  string   `json:"placehold"`
			DateFormat string   `json:"dateformat"`
			TimeFormat string   `json:"timeformat"`
			Optionname []string `json:"optioname"`
			ImgSrc     string   `json:"imgsrc"`
			Url        string   `json:"url"`
		}

		type field struct {
			Fiedlvalue []Fiedlvalue `json:"fiedlvalue"`
		}

		var fieldval field

		json.Unmarshal([]byte(channelupt.FieldValues), &fieldval)

		if channelupt.ChannelName == "" || channelupt.ChannelDescription == "" {

			return errors.New("empty value")

		}

		var chn TblChannel

		chn.ChannelName = channelupt.ChannelName

		chn.ChannelDescription = channelupt.ChannelDescription

		chn.ModifiedBy = userid

		chn.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		CH.UpdateChannelDetails(&chn, channelid, Ch.Authority.DB)

		var modpermissionupdate auth.TblModulePermission

		modpermissionupdate.RouteName = "/channel/entrylist/" + strconv.Itoa(channelid)

		modpermissionupdate.DisplayName = channelupt.ChannelName

		auth.AS.UpdateChannelNameInEntries(&modpermissionupdate, Ch.Authority.DB)

		var getfieldid TblChannel

		CH.GetChannelById(&getfieldid, channelid, Ch.Authority.DB)

		var fldgrpupd TblFieldGroup

		fldgrpupd.ModifiedBy = userid

		fldgrpupd.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		CH.UpdateFieldGroup(&fldgrpupd, getfieldid.FieldGroupId, Ch.Authority.DB)

		var grpfield []TblGroupField

		CH.GetFieldIdByGroupId(&grpfield, getfieldid.FieldGroupId, Ch.Authority.DB)

		var isdeleteids []int

		var optval []string

		for _, gfd := range grpfield {

			for _, val := range fieldval.Fiedlvalue {

				if gfd.FieldId == val.FieldId {

					isdeleteids = append(isdeleteids, val.FieldId)

					optval = append(optval, val.Optionname...)

					var fieldsname TblField

					fieldsname.OrderIndex, _ = strconv.Atoi(val.Id)

					fieldsname.FieldName = strings.TrimSpace(val.Name)

					fieldsname.FieldDesc = val.Desc

					fieldsname.InitialValue = val.Initial

					fieldsname.MandatoryField = val.Mandatory

					fieldsname.Placeholder = val.Placehold

					fieldsname.DatetimeFormat = val.DateFormat

					fieldsname.TimeFormat = val.TimeFormat

					fieldsname.ModifiedBy = userid

					fieldsname.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					fieldsname.Url = val.Url

					CH.UpdateFieldDetails(&fieldsname, val.FieldId, Ch.Authority.DB)

					var fieldoptdel TblFieldOption

					fieldoptdel.DeletedBy = userid

					fieldoptdel.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					CH.DeleteFieldOptionById(&fieldoptdel, val.Optionname, val.FieldId, Ch.Authority.DB)

				}

			}
		}

		var grpfieldnotinfid []TblGroupField

		CH.GetNotInFieldId(&grpfieldnotinfid, isdeleteids, getfieldid.FieldGroupId, Ch.Authority.DB)

		var delids []int

		for _, val := range grpfieldnotinfid {

			delids = append(delids, val.FieldId)

		}

		if len(fieldval.Fiedlvalue) == 0 {

			for _, gfd := range grpfield {

				delids = append(delids, gfd.FieldId)
			}

		}

		var fieldv TblField

		fieldv.DeletedBy = userid

		fieldv.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		DeleteFieldById(&fieldv, delids, Ch.Authority.DB)

		isdeleteids = []int{}

		optval = []string{}

		for _, val := range fieldval.Fiedlvalue {

			var cid int

			if val.FieldId == 0 {

				var cfld TblField

				cfld.FieldName = val.Name

				cfld.FieldDesc = val.Desc

				cfld.FieldTypeId, _ = strconv.Atoi(val.Fid)

				cfld.CreatedBy = userid

				cfld.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				cfld.MandatoryField = val.Mandatory

				cfld.OrderIndex, _ = strconv.Atoi(val.Id)

				cfld.InitialValue = val.Initial

				cfld.Placeholder = val.Placehold

				cfld.ImagePath = val.ImgSrc

				if val.Fid == "4" {

					cfld.DatetimeFormat = val.DateFormat

					cfld.TimeFormat = val.TimeFormat

				}
				if val.Fid == "6" {

					cfld.DatetimeFormat = val.DateFormat
				}

				if len(val.Optionname) > 0 {

					cfld.OptionExist = 1
				}

				cfid, fiderr := CH.CreateFields(&cfld, Ch.Authority.DB)

				if fiderr != nil {

					log.Println(fiderr)

				}

				cid = cfid.Id

				/*create group field*/
				var grpfield TblGroupField

				grpfield.FieldGroupId = getfieldid.FieldGroupId

				grpfield.FieldId = cfid.Id

				grpfielderr := CH.CreateGroupField(&grpfield, Ch.Authority.DB)

				if grpfielderr != nil {

					log.Println(grpfielderr)

				}

			} else {

				cid = val.FieldId
			}

			/*option value create*/
			for _, opt := range val.Optionname {

				var optchk TblFieldOption

				err := CH.CheckOptionAlreadyExist(&optchk, opt, cid, Ch.Authority.DB)

				if errors.Is(err, gorm.ErrRecordNotFound) {

					var fldopt TblFieldOption

					fldopt.OptionName = opt

					fldopt.OptionValue = opt

					fldopt.FieldId = cid

					fldopt.CreatedBy = userid

					fldopt.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

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

	check, err := Ch.Authority.IsGranted("Channels", authcore.CRUD)

	if err != nil {

		return err
	}

	if check {

		if channelid <= 0 {

			return errors.New("invalid channelid cannot delete")
		}

		chid := strconv.Itoa(channelid)

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
	}

	return errors.New("not authorized")
}
