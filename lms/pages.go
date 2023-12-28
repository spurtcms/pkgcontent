package lms

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/spurtcms/spurtcms-core/auth"
	authcore "github.com/spurtcms/spurtcms-core/auth"
	"github.com/spurtcms/spurtcms-core/member"
	memberaccore "github.com/spurtcms/spurtcms-core/memberaccess"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Page struct {
	Authority *authcore.Authorization
}

type MemberPage struct {
	MemAuth *authcore.Authorization
}

type PageStrut struct{}

var PG PageStrut

/*Get Page log*/
func (p Page) PageAliasesLog(spaceid int) (log []PageLog, err error) {

	_, _, checkerr := authcore.VerifyToken(p.Authority.Token, p.Authority.Secret)

	if checkerr != nil {

		return []PageLog{}, checkerr
	}

	var pagelog []TblPageAliasesLog

	err2 := PG.GetPageLogDetails(&pagelog, spaceid, p.Authority.DB)

	var finallog []PageLog

	for _, val := range pagelog {

		var log PageLog

		log.Username = val.Username

		if val.ModifiedOn.IsZero() {

			log.Status = "draft"
		} else {
			log.Status = "Updated"
		}

		if val.Status == "publish" {
			log.Status = val.Status
		}

		log.Date = val.CreatedOn.Format("02-01-2006 15:04 PM")

		finallog = append(finallog, log)

	}

	return finallog, err2

}

/*list page*/
func (p Page) PageList(spaceid int) ([]PageGroups, []Pages, []SubPages, error) {

	_, _, checkerr := authcore.VerifyToken(p.Authority.Token, p.Authority.Secret)

	if checkerr != nil {

		return []PageGroups{}, []Pages{}, []SubPages{}, checkerr
	}

	var group []TblPagesGroup

	var pagegroups []PageGroups

	PG.SelectGroup(&group, spaceid, []int{}, p.Authority.DB)

	for _, group := range group {

		var pagegroup TblPagesGroupAliases

		var page_group PageGroups

		PG.PageGroup(&pagegroup, group.Id, p.Authority.DB)

		page_group.GroupId = pagegroup.PageGroupId

		page_group.Name = pagegroup.GroupName

		page_group.OrderIndex = pagegroup.OrderIndex

		pagegroups = append(pagegroups, page_group)

	}
	var page []TblPage

	var pages []Pages

	var subpages []SubPages

	PG.SelectPage(&page, spaceid, []int{}, p.Authority.DB)

	for _, page := range page {

		var page_aliases TblPageAliases

		if page.ParentId != 0 {

			sid := page.Id

			var subpage SubPages

			PG.PageAliases(&page_aliases, sid, p.Authority.DB)

			subpage.SpgId = page_aliases.PageId

			subpage.Name = page_aliases.PageTitle

			subpage.Content = page_aliases.PageDescription

			subpage.ParentId = page.ParentId

			subpage.OrderIndex = page_aliases.PageSuborder

			subpages = append(subpages, subpage)

		} else {

			pgid := page.Id

			var one_page Pages

			PG.PageAliases(&page_aliases, pgid, p.Authority.DB)

			one_page.PgId = page_aliases.PageId

			one_page.Name = page_aliases.PageTitle

			one_page.Content = page_aliases.PageDescription

			one_page.OrderIndex = page_aliases.OrderIndex

			one_page.Pgroupid = page.PageGroupId

			one_page.ParentId = page.ParentId

			pages = append(pages, one_page)

		}

	}

	return pagegroups, pages, subpages, nil
}

/*list page*/
func (p MemberPage) MemberPageList(spaceid int) ([]PageGroups, []Pages, []SubPages, error) {

	var mem memberaccore.AccessAuth

	mem.Authority = *p.MemAuth

	var memberrest auth.TblModule

	merr := PG.MemberRestrictActive(&memberrest, p.MemAuth.DB)

	var PageIds []int

	var GroupIds []int

	if gorm.ErrRecordNotFound != merr || memberrest.IsActive == 1 {

		pageid, err := mem.GetPage()

		if err != nil {

			log.Println(err)
		}

		grpid, err1 := mem.GetGroup()

		if err1 != nil {

			log.Println(err1)
		}

		PageIds = pageid

		GroupIds = grpid

	}

	var group []TblPagesGroup

	var pagegroups []PageGroups

	PG.SelectGroup(&group, spaceid, GroupIds, p.MemAuth.DB)

	for _, group := range group {

		var pagegroup TblPagesGroupAliases

		var page_group PageGroups

		PG.PageGroup(&pagegroup, group.Id, p.MemAuth.DB)

		page_group.GroupId = pagegroup.PageGroupId

		page_group.Name = pagegroup.GroupName

		page_group.OrderIndex = pagegroup.OrderIndex

		pagegroups = append(pagegroups, page_group)

	}
	var page []TblPage

	var pages []Pages

	var subpages []SubPages

	PG.SelectPage(&page, spaceid, PageIds, p.MemAuth.DB)

	for _, page := range page {

		var page_aliases TblPageAliases

		if page.ParentId != 0 {

			sid := page.Id

			var subpage SubPages

			PG.PageAliases(&page_aliases, sid, p.MemAuth.DB)

			subpage.SpgId = page_aliases.PageId

			subpage.Name = page_aliases.PageTitle

			// subpage.Content = page_aliases.PageDescription

			subpage.ParentId = page.ParentId

			subpage.OrderIndex = page_aliases.PageSuborder

			subpages = append(subpages, subpage)

		} else {

			pgid := page.Id

			var one_page Pages

			PG.PageAliases(&page_aliases, pgid, p.MemAuth.DB)

			one_page.PgId = page_aliases.PageId

			one_page.Name = page_aliases.PageTitle

			// one_page.Content = page_aliases.PageDescription

			one_page.OrderIndex = page_aliases.OrderIndex

			one_page.Pgroupid = page.PageGroupId

			one_page.ParentId = page.ParentId

			pages = append(pages, one_page)

		}

	}

	return pagegroups, pages, subpages, nil
}

/*Get Page content - PAGE VIEW*/
func (m MemberPage) GetPageContent(pageid int) (TblPageAliases, error) {

	var mem memberaccore.AccessAuth

	mem.Authority = *m.MemAuth

	flg, err := mem.CheckPageLogin(pageid)

	if err != nil {

		log.Println(err)

	}

	if flg {

		var tblpage TblPageAliases

		PG.GetContentByPageId(&tblpage, pageid, m.MemAuth.DB)

		return tblpage, nil
	}

	return TblPageAliases{}, err
}

/*Get Page content - PAGE VIEW*/
func (m MemberPage) GetNotes(pageid int) ([]TblMemberNotesHighlight, error) {

	memberid, _, checkerr := member.VerifyToken(m.MemAuth.Token, m.MemAuth.Secret)

	if checkerr != nil {

		return []TblMemberNotesHighlight{}, checkerr
	}

	var notes []TblMemberNotesHighlight

	PG.GetNotes(&notes, memberid, pageid, m.MemAuth.DB)

	return notes, nil
}

/*Get Page content - PAGE VIEW*/
func (m MemberPage) GetHighlights(pageid int) ([]TblMemberNotesHighlight, error) {

	memberid, _, checkerr := member.VerifyToken(m.MemAuth.Token, m.MemAuth.Secret)

	if checkerr != nil {

		return []TblMemberNotesHighlight{}, checkerr
	}

	var Highlights []TblMemberNotesHighlight

	PG.GetHighlights(&Highlights, memberid, pageid, m.MemAuth.DB)

	return Highlights, nil
}

/*Update Notes*/
func (m MemberPage) UpdateNotes(pageid int, content string) (flg bool, err error) {

	memberid, _, checkerr := member.VerifyToken(m.MemAuth.Token, m.MemAuth.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var notes TblMemberNotesHighlight

	notes.PageId = pageid

	notes.MemberId = memberid

	notes.NotesHighlightsContent = content

	notes.NotesHighlightsType = "notes"

	notes.CreatedBy = memberid

	notes.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err1 := PG.UpdateNotesHighlights(&notes, "notes", m.MemAuth.DB)

	if err1 != nil {

		return false, err1
	}

	return true, nil
}

/*Update Highlights*/
func (m MemberPage) UpdateHighlights(highreq HighlightsReq) (flg bool, err error) {

	memberid, _, checkerr := member.VerifyToken(m.MemAuth.Token, m.MemAuth.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var notes TblMemberNotesHighlight

	notes.PageId = highreq.Pageid

	notes.NotesHighlightsContent = highreq.Content

	notes.NotesHighlightsType = "highlights"

	notes.HighlightsConfiguration = datatypes.JSONMap{"start": highreq.Start, "offset": highreq.Offset, "selectedpara": highreq.SelectPara, "color": highreq.ContentColor}

	notes.MemberId = memberid

	notes.CreatedBy = memberid

	notes.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err1 := PG.UpdateNotesHighlights(&notes, "highlights", m.MemAuth.DB)

	if err1 != nil {

		return false, err1
	}

	return true, nil
}

/*Remove Highlights*/
func (m MemberPage) RemoveHighlightsandNotes(Id int) (flg bool, err error) {

	memberid, _, checkerr := member.VerifyToken(m.MemAuth.Token, m.MemAuth.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	var notes TblMemberNotesHighlight

	notes.Id = Id

	notes.DeletedBy = memberid

	notes.IsDeleted = 1

	notes.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err1 := PG.RemoveHighlights(&notes, m.MemAuth.DB)

	if err1 != nil {

		return false, err1
	}

	return true, nil
}

/*Create page*/
func (p Page) InsertPage(Pagec PageCreate) error {

	userid, _, checkerr := authcore.VerifyToken(p.Authority.Token, p.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	var status string

	if Pagec.Status == "publish" {

		status = "publish"

	} else if Pagec.Status == "save" {

		status = "draft"
	}

	// type pagebind struct {
	// 	NewPages []Pages `json:"newpages"`
	// }

	// type GRP struct {
	// 	NewGroup []PageGroups `json:"newGroup"`
	// }

	// type subpagebind struct {
	// 	SubPage []SubPages `json:"subpage"`
	// }

	// type deletePAGE struct {
	// 	NewPages []Pages `json:"deletePage"`
	// }

	// type deleteGRP struct {
	// 	NewGroup []PageGroups `json:"deletegroup"`
	// }

	// type deleteSUB struct {
	// 	SubPage []SubPages `json:"deletesub"`
	// }

	type TempCheck struct {
		FrontId    int
		NewFrontId int
		DBid       int
	}

	spaceId := Pagec.SpaceId

	// creategroup := Pagec.NewGroup

	// createpages := Pagec.NewPages

	// createsubpage := Pagec.SubPage

	// deletegroup := Pagec.DeleteGroup

	// deletepage := Pagec.DeletePages

	// deletesub := Pagec.DeleteSubPage

	/*CreateFunc*/

	// var createGroup GRP

	// json.Unmarshal([]byte(creategroup), &createGroup)

	// var createPage pagebind

	// json.Unmarshal([]byte(createpages), &createPage)

	// var createSub subpagebind

	// json.Unmarshal([]byte(createsubpage), &createSub)

	var Temparr []TempCheck

	var err error
	for _, val := range Pagec.NewGroup {

		/*check if exists*/
		var ckgroupali TblPagesGroup

		PG.CheckGroupExists(&ckgroupali, val.GroupId, spaceId, p.Authority.DB)

		if ckgroupali.Id == 0 && val.NewGroupId != 0 {

			/*Group create tbl_page_group*/
			var groups TblPagesGroup

			groups.SpacesId = spaceId

			groups.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			groups.CreatedBy = userid

			grpreturn, _ := PG.CreatePageGroup(&groups, p.Authority.DB)

			/*group aliases tbl_page_group_aliases*/
			var groupali TblPagesGroupAliases

			groupali.PageGroupId = grpreturn.Id

			groupali.GroupName = strings.ToUpper(val.Name)

			groupali.GroupSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

			groupali.LanguageId = 1

			groupali.OrderIndex = val.OrderIndex

			groupali.CreatedBy = userid

			groupali.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			err = PG.CreatePageGroupAliases(&groupali, p.Authority.DB)

		} else {

			var uptgroup TblPagesGroupAliases

			uptgroup.GroupName = strings.ToUpper(val.Name)

			uptgroup.GroupSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

			uptgroup.LanguageId = 1

			uptgroup.OrderIndex = val.OrderIndex

			uptgroup.ModifiedBy = userid

			uptgroup.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			err = PG.UpdatePageGroupAliases(&uptgroup, val.GroupId, p.Authority.DB)

		}

	}

	for _, val := range Pagec.NewPages {

		var newgrpid int

		for _, grp := range Pagec.NewGroup {

			if val.Pgroupid == grp.GroupId && val.NewGrpId == grp.NewGroupId {

				var getgid TblPagesGroupAliases

				PG.GetPageGroupByName(&getgid, spaceId, grp.Name, p.Authority.DB)

				newgrpid = getgid.PageGroupId

				break

			}
		}

		var checkpage TblPage

		err := PG.CheckPageExists(&checkpage, val.PgId, spaceId, p.Authority.DB)

		if val.Pgroupid == 0 && val.NewGrpId == 0 && val.ParentId == 0 {

			if errors.Is(err, gorm.ErrRecordNotFound) {

				/*page creation tbl_page*/
				var page TblPage

				page.PageGroupId = 0

				page.SpacesId = spaceId

				page.ParentId = 0

				page.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				page.CreatedBy = userid

				pageret, _ := PG.CreatePage(&page, p.Authority.DB)

				for _, newval := range Pagec.SubPage {

					if newval.ParentId == 0 && newval.NewParentId == val.ParentId || newval.ParentId == 0 && newval.NewParentId == val.NewPgId || newval.ParentId == val.PgId && newval.NewParentId == 0 || newval.ParentId == val.NewPgId && newval.NewParentId == 0 {

						var Temarr TempCheck

						Temarr.FrontId = newval.SpgId

						Temarr.NewFrontId = newval.NewSpId

						Temarr.DBid = pageret.Id

						Temparr = append(Temparr, Temarr)

					}

				}

				/*page creation tbl_page_aliases*/
				var pageali TblPageAliases

				pageali.LanguageId = 1

				pageali.PageId = pageret.Id

				pageali.PageTitle = val.Name

				pageali.PageDescription = val.Content

				pageali.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

				pageali.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				pageali.CreatedBy = userid

				pageali.OrderIndex = val.OrderIndex

				pageali.Status = status

				pageali.LastRevisionDate, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				pageali.LastRevisionNo = 1

				pageali.Access = "public"

				err = PG.CreatepageAliases(&pageali, p.Authority.DB)

				/*This is for log*/
				var pagelog TblPageAliasesLog

				pagelog.PageId = pageret.Id

				pagelog.LanguageId = 1

				pagelog.PageTitle = val.Name

				pagelog.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

				pagelog.PageDescription = val.Content

				pagelog.PageId = pageret.Id

				pagelog.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				pagelog.CreatedBy = userid

				pagelog.Status = status

				pagelog.Access = "public"

				PG.PageAliasesLog(&pagelog, p.Authority.DB)

			} else {

				var uptpage TblPage

				uptpage.PageGroupId = val.Pgroupid

				uptpage.ParentId = val.ParentId

				PG.UpdatePage(&uptpage, val.PgId, p.Authority.DB)

				var uptpageali TblPageAliases

				uptpageali.PageTitle = val.Name

				uptpageali.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

				uptpageali.PageDescription = val.Content

				uptpageali.OrderIndex = val.OrderIndex

				uptpageali.Status = status

				uptpageali.ModifiedBy = userid

				uptpageali.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				uptpageali.LastRevisionDate, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				err = PG.UpdatePageAliase(&uptpageali, val.PgId, p.Authority.DB)

				/*This is for log*/
				var pagelog TblPageAliasesLog

				pagelog.LanguageId = 1

				pagelog.PageId = val.PgId

				pagelog.PageTitle = val.Name

				pagelog.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

				pagelog.PageDescription = val.Content

				pagelog.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				pagelog.CreatedBy = userid

				pagelog.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				pagelog.ModifiedBy = userid

				pagelog.Status = status

				pagelog.Access = "public"

				PG.PageAliasesLog(&pagelog, p.Authority.DB)
			}

		}

		if (val.NewGrpId != 0 && val.Pgroupid == 0) || (val.NewGrpId == 0 && val.Pgroupid != 0) && val.ParentId == 0 {

			if errors.Is(err, gorm.ErrRecordNotFound) {

				/*page creation tbl_page*/
				var page TblPage

				page.PageGroupId = newgrpid

				page.SpacesId = spaceId

				page.ParentId = 0

				page.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				page.CreatedBy = userid

				pageret, _ := PG.CreatePage(&page, p.Authority.DB)

				for _, newval := range Pagec.SubPage {

					if newval.ParentId == 0 && newval.NewParentId == val.ParentId || newval.ParentId == 0 && newval.NewParentId == val.NewPgId || newval.ParentId == val.PgId && newval.NewParentId == 0 || newval.ParentId == val.NewPgId && newval.NewParentId == 0 {

						var Temarr TempCheck

						Temarr.FrontId = newval.SpgId

						Temarr.NewFrontId = newval.NewSpId

						Temarr.DBid = pageret.Id

						Temparr = append(Temparr, Temarr)

					}

				}

				/*page creation tbl_page_aliases*/
				var pageali TblPageAliases

				pageali.LanguageId = 1

				pageali.PageId = pageret.Id

				pageali.PageTitle = val.Name

				pageali.PageDescription = val.Content

				pageali.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

				pageali.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				pageali.CreatedBy = userid

				pageali.OrderIndex = val.OrderIndex

				pageali.Status = status

				pageali.LastRevisionDate, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				pageali.LastRevisionNo = 1

				pageali.Access = "public"

				err = PG.CreatepageAliases(&pageali, p.Authority.DB)

				/*This is for log*/
				var pagelog TblPageAliasesLog

				pagelog.PageId = pageret.Id

				pagelog.LanguageId = 1

				pagelog.PageTitle = val.Name

				pagelog.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

				pagelog.PageDescription = val.Content

				pagelog.PageId = pageret.Id

				pagelog.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				pagelog.CreatedBy = userid

				pagelog.Status = status

				pagelog.Access = "public"

				PG.PageAliasesLog(&pagelog, p.Authority.DB)

			} else {

				var uptpage TblPage

				uptpage.PageGroupId = val.Pgroupid

				uptpage.ParentId = val.ParentId

				PG.UpdatePage(&uptpage, val.PgId, p.Authority.DB)

				var uptpageali TblPageAliases

				uptpageali.PageTitle = val.Name

				uptpageali.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

				uptpageali.PageDescription = val.Content

				uptpageali.OrderIndex = val.OrderIndex

				uptpageali.Status = status

				uptpageali.ModifiedBy = userid

				uptpageali.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				uptpageali.LastRevisionDate, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				err = PG.UpdatePageAliase(&uptpageali, val.PgId, p.Authority.DB)

				/*This is for log*/
				var pagelog TblPageAliasesLog

				pagelog.PageId = val.PgId

				pagelog.LanguageId = 1

				pagelog.PageTitle = val.Name

				pagelog.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

				pagelog.PageDescription = val.Content

				pagelog.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				pagelog.CreatedBy = userid

				pagelog.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				pagelog.ModifiedBy = userid

				pagelog.Status = status

				pagelog.Access = "public"

				PG.PageAliasesLog(&pagelog, p.Authority.DB)
			}

		}

	}

	/*createsub*/
	for _, val := range Pagec.SubPage {

		for _, pg := range Pagec.NewPages {

			var newgrpid int

			var pgid int

			newgrpid = val.PgroupId

			if val.NewPgroupId != 0 {

				for _, grp := range Pagec.NewGroup {

					if val.NewPgroupId == grp.NewGroupId {

						var getgid TblPagesGroupAliases

						PG.GetPageGroupByName(&getgid, spaceId, grp.Name, p.Authority.DB)

						newgrpid = getgid.PageGroupId

						break

					}
				}

			}

			for _, pgd := range Pagec.NewPages {

				if val.NewParentId == pgd.NewPgId && val.NewParentId == pgd.ParentId && val.ParentId == pgd.PgId {

					if val.SpgId == 0 {

						var getpage TblPageAliases

						PG.GetPageDataByName(&getpage, spaceId, pgd.Name, p.Authority.DB)

						pgid = getpage.PageId

						break
					}
				}

			}

			for _, newpgid := range Temparr {

				if newpgid.FrontId == 0 && newpgid.NewFrontId == val.NewSpId || newpgid.FrontId == 0 && newpgid.NewFrontId == val.SpgId || newpgid.FrontId == val.SpgId && newpgid.NewFrontId == 0 || newpgid.FrontId == val.NewSpId && newpgid.NewFrontId == 0 {

					pgid = newpgid.DBid

					break
				}

			}

			if val.NewParentId == pg.NewPgId || val.ParentId == pg.PgId {

				if val.SpgId == 0 {

					/*page creation tbl_page*/
					var page TblPage

					page.PageGroupId = newgrpid

					page.SpacesId = spaceId

					page.ParentId = pgid

					page.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					page.CreatedBy = userid

					pageret, _ := PG.CreatePage(&page, p.Authority.DB)

					/*page creation tbl_page_aliases*/
					var pageali TblPageAliases

					pageali.LanguageId = 1

					pageali.PageId = pageret.Id

					pageali.PageTitle = val.Name

					pageali.PageDescription = val.Content

					pageali.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

					pageali.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					pageali.CreatedBy = userid

					pageali.PageSuborder = val.OrderIndex

					pageali.Status = status

					pageali.Access = "public"

					pageali.LastRevisionDate, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					pageali.LastRevisionNo = 1

					err = PG.CreatepageAliases(&pageali, p.Authority.DB)

					/*This is for log*/
					var pagelog TblPageAliasesLog

					pagelog.LanguageId = 1

					pagelog.PageId = pageret.Id

					pagelog.PageTitle = val.Name

					pagelog.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

					pagelog.PageDescription = val.Content

					pagelog.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					pagelog.CreatedBy = userid

					pagelog.Status = status

					pagelog.Access = "public"

					PG.PageAliasesLog(&pagelog, p.Authority.DB)

				} else {

					var uptpage TblPage

					uptpage.PageGroupId = newgrpid

					uptpage.ParentId = val.ParentId

					PG.UpdatePage(&uptpage, val.SpgId, p.Authority.DB)

					var uptpageali TblPageAliases

					uptpageali.PageTitle = val.Name

					uptpageali.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

					uptpageali.PageDescription = val.Content

					uptpageali.PageSuborder = val.OrderIndex

					uptpageali.Status = status

					uptpageali.ModifiedBy = userid

					uptpageali.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					uptpageali.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					err = PG.UpdatePageAliase(&uptpageali, val.SpgId, p.Authority.DB)

					/*This is for log*/
					var pagelog TblPageAliasesLog

					pagelog.PageId = val.SpgId

					pagelog.LanguageId = 1

					pagelog.PageTitle = val.Name

					pagelog.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

					pagelog.PageDescription = val.Content

					pagelog.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					pagelog.CreatedBy = userid

					pagelog.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					pagelog.ModifiedBy = userid

					pagelog.Status = status

					pagelog.Access = "public"

					PG.PageAliasesLog(&pagelog, p.Authority.DB)
				}

				break
			}

			if val.ParentId == pg.PgId {

				if val.SpgId == 0 {

					/*page creation tbl_page*/
					var page TblPage

					page.PageGroupId = newgrpid

					page.SpacesId = spaceId

					page.ParentId = pgid

					page.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					page.CreatedBy = userid

					pageret, _ := PG.CreatePage(&page, p.Authority.DB)

					/*page creation tbl_page_aliases*/
					var pageali TblPageAliases

					pageali.LanguageId = 1

					pageali.PageId = pageret.Id

					pageali.PageTitle = val.Name

					pageali.PageDescription = val.Content

					pageali.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

					pageali.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					pageali.CreatedBy = userid

					pageali.PageSuborder = val.OrderIndex

					pageali.Status = status

					pageali.Access = "public"

					pageali.LastRevisionNo = 1

					pageali.LastRevisionDate, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					err = PG.CreatepageAliases(&pageali, p.Authority.DB)

					/*This is for log*/
					var pagelog TblPageAliasesLog

					pagelog.PageTitle = val.Name

					pagelog.LanguageId = 1

					pagelog.PageId = val.SpgId

					pagelog.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

					pagelog.PageDescription = val.Content

					pagelog.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					pagelog.CreatedBy = userid

					pagelog.Status = status

					pagelog.Access = "public"

					pageali.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					PG.PageAliasesLog(&pagelog, p.Authority.DB)

				} else {

					var uptpage TblPage

					uptpage.PageGroupId = pg.Pgroupid

					uptpage.ParentId = val.ParentId

					PG.UpdatePage(&uptpage, val.SpgId, p.Authority.DB)

					var uptpageali TblPageAliases

					uptpageali.PageTitle = val.Name

					uptpageali.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

					uptpageali.PageDescription = val.Content

					uptpageali.PageSuborder = val.OrderIndex

					uptpageali.Status = status

					uptpageali.ModifiedBy = userid

					uptpageali.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					uptpageali.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					err = PG.UpdatePageAliase(&uptpageali, val.SpgId, p.Authority.DB)

					/*This is for log*/
					var pagelog TblPageAliasesLog

					pagelog.PageTitle = val.Name

					pagelog.LanguageId = 1

					pagelog.PageId = val.SpgId

					pagelog.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

					pagelog.PageDescription = val.Content

					pagelog.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					pagelog.CreatedBy = userid

					pagelog.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

					pagelog.ModifiedBy = userid

					pagelog.Status = status

					pagelog.Access = "public"

					PG.PageAliasesLog(&pagelog, p.Authority.DB)
				}

				break
			}

		}

	}

	/*DeleteFunc*/

	// var deleteGroup deleteGRP

	// json.Unmarshal([]byte(deletegroup), &deleteGroup)

	for _, val := range Pagec.DeleteGroup {

		var deletegroup TblPagesGroup

		deletegroup.DeletedBy = userid

		deletegroup.IsDeleted = 1

		deletegroup.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		PG.DeletePageGroup(&deletegroup, val.GroupId, p.Authority.DB)

		var deletegroupali TblPagesGroupAliases

		deletegroupali.DeletedBy = userid

		deletegroupali.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		deletegroupali.IsDeleted = 1

		err = PG.DeletePageGroupAliases(&deletegroupali, val.GroupId, p.Authority.DB)

	}

	// var deletePage deletePAGE

	// json.Unmarshal([]byte(deletepage), &deletePage)

	for _, val := range Pagec.DeletePages {

		var deletegroup TblPage

		deletegroup.DeletedBy = userid

		deletegroup.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		PG.DeletePage(&deletegroup, val.PgId, p.Authority.DB)

		var deletegroupali TblPageAliases

		deletegroupali.DeletedBy = userid

		deletegroupali.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		deletegroupali.IsDeleted = 1

		err = PG.DeletePageAliases(&deletegroupali, val.PgId, p.Authority.DB)

		/*This is for log*/
		var pagelog TblPageAliasesLog

		pagelog.PageTitle = val.Name

		pagelog.LanguageId = 1

		pagelog.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

		pagelog.PageDescription = val.Content

		pagelog.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		pagelog.CreatedBy = userid

		pagelog.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		pagelog.DeletedBy = userid

		pagelog.Status = status

		pagelog.Access = "public"

		PG.PageAliasesLog(&pagelog, p.Authority.DB)

	}

	// var deleteSubPage deleteSUB

	// json.Unmarshal([]byte(deletesub), &deleteSubPage)

	for _, val := range Pagec.DeleteSubPage {

		var deletegroup TblPage

		deletegroup.DeletedBy = userid

		deletegroup.IsDeleted = 1

		deletegroup.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		PG.DeletePage(&deletegroup, val.SpgId, p.Authority.DB)

		var deletegroupali TblPageAliases

		deletegroupali.DeletedBy = userid

		deletegroupali.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		deletegroupali.IsDeleted = 1

		err = PG.DeletePageAliase(&deletegroupali, val.SpgId, p.Authority.DB)

		/*This is for log*/
		var pagelog TblPageAliasesLog

		pagelog.PageTitle = val.Name

		pagelog.LanguageId = 1

		pagelog.PageSlug = strings.ToLower(strings.ReplaceAll(val.Name, " ", "_"))

		pagelog.PageDescription = val.Content

		pagelog.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		pagelog.CreatedBy = userid

		pagelog.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		pagelog.DeletedBy = userid

		pagelog.Status = status

		pagelog.Access = "public"

		PG.PageAliasesLog(&pagelog, p.Authority.DB)

	}

	Temparr = nil

	if err != nil {

		return err
	}

	return nil
}
