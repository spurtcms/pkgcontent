package channels

import (
	"log"
	"strings"
	"time"
)

type EntriesRepository interface {
	
	GetAllChannelEntriesList(channelid int, limit, offset int, filter EntriesFilter) (entries []TblChannelEntries, filterentriescount int, totalentriescount int, err error)

	GetPublishedChannelEntriesList(limit, offset int, filter EntriesFilter) (entries []TblChannelEntries, filterentriescount int, totalentriescount int, err error)

	CreateEntry(entriesrequired EntriesRequired) (entry TblChannelEntries, flg bool, err error)

	DeleteEntry(ChannelName string, Entryid int) (bool, error)

	GetAdditionalFieldDataBychannelId(ChannelName string, EntryId int) ([]TblChannelEntryField, error)

	GetEntryDetailsById(ChannelName string, EntryId int) (TblChannelEntries, error)

	UpdateEntryDetailsById(entriesrequired EntriesRequired, ChannelName string, EntryId int) (bool, error)

	EntryStatus(ChannelName string, EntryId int, status int) (bool, error)
}

/*all channel Entries List*/
//if channelid 0 get all channel entries
//if channelid not eq 0 to get particular entries of the channel
func (ch Channels) GetAllChannelEntriesList(channelid int, limit, offset int, filter EntriesFilter) (entries []TblChannelEntries, filterentriescount int, totalentriescount int, err error) {

	if filter.Status == "Draft" {

		filter.Status = "0"

	} else if filter.Status == "Published" {

		filter.Status = "1"

	} else if filter.Status == "Unpublished" {

		filter.Status = "2"
	}
	var chnentry []TblChannelEntries

	cmod.ChannelEntryList(&chnentry, limit, offset, channelid, filter, false, true, ch.DatabaseConnection)

	var chnentry1 []TblChannelEntries

	filtercount, _ := cmod.ChannelEntryList(&chnentry1, 0, 0, channelid, filter, false, true, ch.DatabaseConnection)

	entrcount, _ := cmod.ChannelEntryList(&chnentry1, 0, 0, channelid, EntriesFilter{}, false, true, ch.DatabaseConnection)

	return chnentry, int(filtercount), int(entrcount), nil
}

// Get published entries
func (ch Channels) GetPublishedChannelEntriesList(limit, offset int, filter EntriesFilter) (entries []TblChannelEntries, filterentriescount int, totalentriescount int, err error) {

	var chnentry []TblChannelEntries

	cmod.ChannelEntryList(&chnentry, limit, offset, 0, filter, true, true, ch.DatabaseConnection)

	filtercount, _ := cmod.ChannelEntryList(&chnentry, 0, 0, 0, filter, true, true, ch.DatabaseConnection)

	var chnentry1 []TblChannelEntries

	entrcount, _ := cmod.ChannelEntryList(&chnentry1, 0, 0, 0, EntriesFilter{}, true, true, ch.DatabaseConnection)

	return chnentry, int(filtercount), int(entrcount), nil
}

// create entry
func (ch Channels) CreateEntry(entriesrequired EntriesRequired) (entry TblChannelEntries, flg bool, err error) {

	var Entries TblChannelEntries

	Entries.Title = entriesrequired.Title

	Entries.Description = entriesrequired.Content

	Entries.CoverImage = entriesrequired.CoverImage

	Entries.MetaTitle = entriesrequired.SEODetails.MetaTitle

	Entries.MetaDescription = entriesrequired.SEODetails.MetaDescription

	Entries.Keyword = entriesrequired.SEODetails.MetaKeywords

	if entriesrequired.SEODetails.MetaSlug == "" {

		Entries.Slug = strings.ReplaceAll(strings.ToLower(entriesrequired.Title), " ", "_")

	} else {

		Entries.Slug = entriesrequired.SEODetails.MetaSlug

	}

	Entries.Status = entriesrequired.Status

	Entries.ChannelId = entriesrequired.ChannelId

	Entries.CategoriesId = entriesrequired.CategoryIds

	Entries.CreatedBy = entriesrequired.CreatedBy

	Entries.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	Entriess, err := cmod.CreateChannelEntry(&Entries, ch.DatabaseConnection)

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

			Entrfield.CreatedBy = entriesrequired.CreatedBy

			Entrfield.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			EntriesField = append(EntriesField, Entrfield)

		}

		ferr := cmod.CreateEntrychannelFields(&EntriesField, ch.DatabaseConnection)

		if ferr != nil {

			log.Println(ferr)
		}
	}

	return Entries, true, nil
}

func (ch Channels) DeleteEntry(Entryid, userid int) (bool, error) {

	var entries TblChannelEntries

	entries.Id = Entryid

	entries.IsDeleted = 1

	entries.DeletedBy = userid

	entries.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := cmod.DeleteChannelEntryId(&entries, Entryid, ch.DatabaseConnection)

	var field TblChannelEntryField

	field.DeletedBy = userid

	field.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err1 := cmod.DeleteChannelEntryFieldId(&field, Entryid, ch.DatabaseConnection)

	if err != nil {

		log.Println(err)
	}

	if err1 != nil {

		log.Println(err)
	}

	return true, nil
}

func (ch Channels) GetAdditionalFieldDataBychannelId(ChannelName string, EntryId int) ([]TblChannelEntryField, error) {

	var EntriesField []TblChannelEntryField

	err := cmod.GetChannelEntryDetailsById(&EntriesField, EntryId, ch.DatabaseConnection)

	if err != nil {

		log.Println(err)
	}

	return EntriesField, nil

}

// get entry details
func (ch Channels) GetEntryDetailsById(ChannelName string, EntryId int) (TblChannelEntries, error) {

	var Entry TblChannelEntries

	err := cmod.GetChannelEntryById(&Entry, EntryId, ch.DatabaseConnection)

	if err != nil {

		log.Println(err)
	}

	return Entry, nil

}

/*update entry details */
func (ch Channels) UpdateEntryDetailsById(entriesrequired EntriesRequired, EntryId int) (bool, error) {

	var Entries TblChannelEntries

	Entries.Title = entriesrequired.Title

	Entries.Description = entriesrequired.Content

	Entries.CoverImage = entriesrequired.CoverImage

	Entries.MetaTitle = entriesrequired.SEODetails.MetaTitle

	Entries.MetaDescription = entriesrequired.SEODetails.MetaDescription

	Entries.Keyword = entriesrequired.SEODetails.MetaKeywords

	if entriesrequired.SEODetails.MetaSlug == "" {

		Entries.Slug = strings.ReplaceAll(strings.ToLower(entriesrequired.Title), " ", "_")

	} else {

		Entries.Slug = entriesrequired.SEODetails.MetaSlug

	}

	Entries.Status = entriesrequired.Status

	Entries.ChannelId = entriesrequired.ChannelId

	Entries.CategoriesId = entriesrequired.CategoryIds

	Entries.ModifiedBy = entriesrequired.CreatedBy

	Entries.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := cmod.UpdateChannelEntryDetails(&Entries, EntryId, ch.DatabaseConnection)

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

			Entrfield.CreatedBy = entriesrequired.CreatedBy

			Entrfield.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			cmod.CreateSingleEntrychannelFields(&Entrfield, ch.DatabaseConnection)

		} else {

			var Entrfield TblChannelEntryField

			Entrfield.Id = val.Id

			Entrfield.ChannelEntryId = EntryId

			Entrfield.FieldName = val.FieldName

			Entrfield.FieldValue = val.FieldValue

			Entrfield.FieldId = val.FieldId

			Entrfield.ModifiedBy = entriesrequired.CreatedBy

			Entrfield.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			cmod.UpdateChannelEntryAdditionalDetails(Entrfield, ch.DatabaseConnection)

		}

	}

	return true, nil

}

// change entries status
func (ch Channels) EntryStatus(ChannelName string, EntryId, userid int, status int) (bool, error) {

	var Entries TblChannelEntries

	Entries.Status = status

	Entries.ModifiedBy = userid

	Entries.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	cmod.PublishQuery(&Entries, EntryId, ch.DatabaseConnection)

	return true, nil

}
