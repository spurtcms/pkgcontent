package lms

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/spurtcms/spurtcms-content/categories"
	"github.com/spurtcms/spurtcms-core/auth"
	authcore "github.com/spurtcms/spurtcms-core/auth"
	memberaccore "github.com/spurtcms/spurtcms-core/memberaccess"
	"gorm.io/gorm"
)

type Space struct {
	Authority *authcore.Authorization
}

type SPM struct{}

var SP SPM

type MemberSpace struct {
	MemAuth *authcore.Authorization
}

func MigrateTable(db *gorm.DB) {

	db.AutoMigrate(
		&TblSpaces{},
		&TblSpacesAliases{},
		&TblPagesCategories{},
		&TblPagesCategoriesAliases{},
		&TblLanguage{},
		&TblPage{},
		&TblPageAliases{},
		&TblPagesGroup{},
		&TblPagesGroupAliases{},
		&TblPageAliasesLog{},
		&TblMemberNotesHighlight{},
	)

	db.Exec(`INSERT INTO PUBLIC.TBL_LANGUAGE(ID,LANGUAGE_NAME,LANGUAGE_CODE,JSON_PATH,IS_STATUS,IS_DEFAULT,	CREATED_BY,CREATED_ON,IS_DELETED) VALUES (1,'English', 'en', 'locales/en.json', 1, 1,1, '2023-09-11 11:27:44',0)`)
}

var IST, _ = time.LoadLocation("Asia/Kolkata")

/*SpaceDetail*/
func (s Space) SpaceDetail(spaceid int) (space TblSpaces, err error) {

	_, _, checkerr := authcore.VerifyToken(s.Authority.Token, s.Authority.Secret)

	if checkerr != nil {

		return TblSpaces{}, checkerr
	}

	check, err := s.Authority.IsGranted("Spaces", authcore.CRUD)

	if err != nil {

		return TblSpaces{}, err
	}

	if check {

		var SP SPM

		var spacename TblSpacesAliases

		err1 := SP.GetSpaceName(&spacename, spaceid, s.Authority.DB)

		var tblspace TblSpaces

		SP.GetSpaceDetails(&tblspace, spaceid, s.Authority.DB)

		tblspace.SpaceName = spacename.SpacesName

		tblspace.CreatedDate = tblspace.CreatedOn.Format("02 Jan 2006 3:04 PM")

		if tblspace.ModifiedOn.IsZero() {

			tblspace.ModifiedDate = tblspace.CreatedOn.Format("02 Jan 2006 3:04 PM")

		} else {

			tblspace.ModifiedDate = tblspace.ModifiedOn.Format("02 Jan 2006 3:04 PM")
		}

		return tblspace, err1

	}
	return TblSpaces{}, errors.New("not authorized")
}

/*spacelist*/
func (s Space) SpaceList(limit, offset int, filter Filter) (tblspace []TblSpacesAliases, totalcount int64, err error) {

	_, _, checkerr := authcore.VerifyToken(s.Authority.Token, s.Authority.Secret)

	if checkerr != nil {

		return []TblSpacesAliases{}, 0, checkerr
	}

	check, err := s.Authority.IsGranted("Spaces", authcore.CRUD)

	if err != nil {

		return []TblSpacesAliases{}, 0, err
	}

	if check {

		var default_lang TblLanguage

		SP.GetDefaultLanguage(&default_lang, s.Authority.DB)

		var spaces []TblSpacesAliases

		_, spaceerr := SP.SpaceList(&spaces, default_lang.Id, limit, offset, filter, []int{}, s.Authority.DB)

		if spaceerr != nil {

			return []TblSpacesAliases{}, 0, spaceerr
		}

		var SpaceDetails []TblSpacesAliases

		for _, space := range spaces {

			var child_page_Category categories.TblCategory

			_, child_page := categories.GetChildPageCategoriess(&child_page_Category, space.PageCategoryId, s.Authority.DB)

			var categorynames []categories.TblCategory

			var flg int

			categorynames = append(categorynames, child_page)

			flg = child_page.ParentId

			if flg != 0 {

			CLOOP:

				for {

					var newchildcategory categories.TblCategory

					_, child := categories.GetChildPageCategoriess(&newchildcategory, flg, s.Authority.DB)

					flg = child.ParentId

					if flg != 0 {

						categorynames = append(categorynames, child)

						goto CLOOP

					} else {

						categorynames = append(categorynames, child)

						break
					}

				}

			}

			var reverseCategoryOrder []categories.TblCategory

			for i := len(categorynames) - 1; i >= 0; i-- {

				reverseCategoryOrder = append(reverseCategoryOrder, categorynames[i])

			}

			space.CategoryNames = reverseCategoryOrder

			space.CreatedDate = space.CreatedOn.Format("02 Jan 2006 03:04 PM")

			if !space.ModifiedOn.IsZero() {

				space.ModifiedDate = space.ModifiedOn.Format("02 Jan 2006 03:04 PM")

			} else {

				space.ModifiedDate = space.CreatedOn.Format("02 Jan 2006 03:04 PM")

			}
			SpaceDetails = append(SpaceDetails, space)

		}

		count, _ := SP.SpaceList(&spaces, default_lang.Id, 0, 0, filter, []int{}, s.Authority.DB)

		return SpaceDetails, count, nil

	}
	return []TblSpacesAliases{}, 0, errors.New("not authorized")
}

/*spacelist*/
func (s MemberSpace) MemberSpaceList(limit, offset int, filter Filter) (tblspace []TblSpacesAliases, totalcount int64, err error) {

	var mem memberaccore.AccessAuth

	mem.Authority = *s.MemAuth

	var default_lang TblLanguage

	SP.GetDefaultLanguage(&default_lang, s.MemAuth.DB)

	var spaces []TblSpacesAliases

	_, spaceerr := SP.MemberSpaceList(&spaces, default_lang.Id, limit, offset, filter, s.MemAuth.DB)

	if spaceerr != nil {

		return []TblSpacesAliases{}, 0, spaceerr
	}

	var SpaceDetails []TblSpacesAliases

	for _, space := range spaces {

		var child_page_Category categories.TblCategory

		_, child_page := categories.GetChildPageCategoriess(&child_page_Category, space.PageCategoryId, s.MemAuth.DB)

		var categorynames []categories.TblCategory

		var flg int

		categorynames = append(categorynames, child_page)

		flg = child_page.ParentId

		if flg != 0 {

		CLOOP:

			for {

				var newchildcategory categories.TblCategory

				_, child := categories.GetChildPageCategoriess(&newchildcategory, flg, s.MemAuth.DB)

				flg = child.ParentId

				if flg != 0 {

					categorynames = append(categorynames, child)

					goto CLOOP

				} else {

					categorynames = append(categorynames, child)

					break
				}

			}

		}

		var reverseCategoryOrder []categories.TblCategory

		for i := len(categorynames) - 1; i >= 0; i-- {

			reverseCategoryOrder = append(reverseCategoryOrder, categorynames[i])

		}

		space.CategoryNames = reverseCategoryOrder

		space.CreatedDate = space.CreatedOn.Format("02 Jan 2006 03:04 PM")

		if !space.ModifiedOn.IsZero() {

			space.ModifiedDate = space.ModifiedOn.Format("02 Jan 2006 03:04 PM")

		} else {

			space.ModifiedDate = space.CreatedOn.Format("02 Jan 2006 03:04 PM")

		}
		SpaceDetails = append(SpaceDetails, space)

	}

	count, _ := SP.MemberSpaceList(&spaces, default_lang.Id, 0, 0, filter, s.MemAuth.DB)

	return SpaceDetails, count, nil

}

// create space
func (s Space) SpaceCreation(SPC SpaceCreation) error {

	userid, _, checkerr := authcore.VerifyToken(s.Authority.Token, s.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := s.Authority.IsGranted("Spaces", authcore.CRUD)

	if err != nil {

		return err
	}

	if check {

		if SPC.Name == "" || SPC.Description == "" || SPC.CategoryId == 0 {

			return errors.New("given some values is empty")
		}

		var spaces TblSpaces

		spaces.PageCategoryId = SPC.CategoryId

		spaces.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		spaces.CreatedBy = userid

		if err := s.Authority.DB.Table("tbl_spaces").Create(&spaces).Error; err != nil {

			return err
		}

		var space TblSpacesAliases

		space.SpacesName = SPC.Name

		space.SpacesDescription = SPC.Description

		space.ImagePath = SPC.ImagePath

		space.LanguageId = SPC.LanguageId

		space.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		space.CreatedBy = userid

		space.SpacesSlug = strings.ToLower(space.SpacesName)

		space.SpacesId = spaces.Id

		if err := s.Authority.DB.Table("tbl_spaces_aliases").Create(&space).Error; err != nil {

			return err
		}

		return nil
	}

	return errors.New("not authorized")
}

// update space
func (s Space) SpaceUpdate(SPC SpaceCreation, spaceid int) error {

	userid, _, checkerr := authcore.VerifyToken(s.Authority.Token, s.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := s.Authority.IsGranted("Spaces", authcore.CRUD)

	if err != nil {

		return err
	}

	if check {

		if SPC.Name == "" || SPC.Description == "" {

			return nil
		}

		var space TblSpaces

		var spaces TblSpacesAliases

		spaces.Id = spaceid

		space.Id = spaceid

		spaces.SpacesName = SPC.Name

		spaces.SpacesDescription = SPC.Description

		spaces.SpacesSlug = strings.ToLower(SPC.Name)

		space.PageCategoryId = SPC.CategoryId

		spaces.ImagePath = SPC.ImagePath

		spaces.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		spaces.ModifiedBy = userid

		space.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		space.ModifiedBy = userid

		err1 := SP.EditSpace(&spaces, spaceid, s.Authority.DB)

		if err != nil {

			return err1
		}

		err2 := SP.UpdateSpace(&space, spaceid, s.Authority.DB)

		if err2 != nil {

			return err2
		}

		return nil

	}
	return errors.New("not authorized")
}

/*Delete Space*/
func (s Space) DeleteSpace(spaceid int) error {

	userid, _, checkerr := authcore.VerifyToken(s.Authority.Token, s.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := s.Authority.IsGranted("Spaces", authcore.CRUD)

	if err != nil {

		return err
	}

	if check {

		var spaces TblSpacesAliases

		var space TblSpaces

		var pageali TblPageAliases

		spaces.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		spaces.DeletedBy = userid

		spaces.IsDeleted = 1

		space.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		space.DeletedBy = userid

		space.IsDeleted = 1

		var deletedon, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		var deletedby = userid

		var isdeleted = 1

		pageali.DeletedOn = deletedon

		pageali.DeletedBy = deletedby

		pageali.IsDeleted = isdeleted

		err1 := SP.DeleteSpaceAliases(&spaces, spaceid, s.Authority.DB)

		if err1 != nil {
			return err
		}

		err2 := SP.DeleteSpace(&space, spaceid, s.Authority.DB)

		if err2 != nil {
			return err2
		}

		var page []TblPage

		SP.GetPageDetailsBySpaceId(&page, spaceid, s.Authority.DB)

		var pid []int

		for _, v := range page {

			pid = append(pid, v.Id)

		}

		var pg TblPage

		pg.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		pg.DeletedBy = userid

		pg.IsDeleted = 1

		SP.DeletePageInSpace(&pg, pid, s.Authority.DB)

		var pgali TblPageAliases

		pgali.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		pgali.DeletedBy = userid

		pgali.IsDeleted = 1

		SP.DeletePageAliInSpace(&pgali, pid, s.Authority.DB)

		var pggroupdel TblPagesGroup

		pggroupdel.DeletedBy = userid

		pggroupdel.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		SP.DeletePageGroup(&pggroupdel, spaceid, s.Authority.DB)

		return nil

	}
	return errors.New("not authorized")
}

// Clone

func (s Space) CloneSpace(SPC SpaceCreation, clonespaceid int) error {

	userid, _, checkerr := authcore.VerifyToken(s.Authority.Token, s.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	if SPC.Name == "" || SPC.Description == "" {

		return errors.New("given some values is empty")
	}

	var spaces TblSpaces

	var space TblSpacesAliases

	space.SpacesName = SPC.Name

	space.SpacesDescription = SPC.Description

	space.SpacesSlug = strings.ToLower(SPC.Name)

	spaces.PageCategoryId = SPC.CategoryId

	space.ImagePath = SPC.ImagePath

	space.LanguageId = SPC.LanguageId

	spaceid := clonespaceid

	spaces.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	spaces.CreatedBy = userid

	space.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	space.CreatedBy = userid

	tblspaces, _ := SP.CreateSpace(&spaces, s.Authority.DB)

	space.SpacesId = tblspaces.Id

	err := SP.CreateSpacesAliases(&space, s.Authority.DB)

	if err != nil {

		log.Println(err)
	}

	var pagegroupdata []TblPagesGroupAliases

	SP.GetPageGroupData(&pagegroupdata, spaceid, s.Authority.DB)

	for _, value := range pagegroupdata {

		var group TblPagesGroup

		group.SpacesId = tblspaces.Id

		groups, _ := SP.CloneSpaceInPagesGroup(&group, s.Authority.DB)

		// var pagegroup TblPagesGroupAliases

		pagegroup := value

		pagegroup.PageGroupId = groups.Id

		SP.ClonePagesGroup(&pagegroup, s.Authority.DB)
	}

	var pageId []TblPageAliases

	SP.GetPageInPage(&pageId, spaceid, s.Authority.DB) //parentid 0 and groupid 0

	for _, val := range pageId {

		var page TblPage

		page.SpacesId = tblspaces.Id

		page.PageGroupId = 0

		page.ParentId = 0

		pageid, _ := SP.ClonePage(&page, s.Authority.DB)

		// var pagesali TblPageAliases

		pagesali := val

		pagesali.PageId = pageid.Id

		SP.ClonePages(&pagesali, s.Authority.DB)

	}

	var pagegroupaldata TblPagesGroupAliases

	SP.GetIdInPage(&pagegroupaldata, spaceid, s.Authority.DB) // parentid = 0 and groupid != 0

	var pagealiase []TblPageAliases

	SP.GetPageAliasesInPage(&pagealiase, spaceid, s.Authority.DB) // parentid = 0 and groupid != 0

	for _, value := range pagealiase {

		var pageal TblPagesGroupAliases

		SP.GetDetailsInPageAli(&pageal, pagegroupaldata.GroupName, tblspaces.Id, s.Authority.DB)

		// var parent TblPage

		// ParentWithChild(&parent, value.PageGroupId, spaceid)

		var page TblPage

		page.SpacesId = tblspaces.Id

		page.PageGroupId = pageal.PageGroupId

		page.ParentId = 0

		pagess, _ := SP.ClonePage(&page, s.Authority.DB)

		// var pagesali TblPageAliases

		pagesali := value

		pagesali.PageId = pagess.Id

		SP.ClonePages(&pagesali, s.Authority.DB)

	}

	var pagealiasedata []TblPageAliases

	SP.GetPageAliasesInPageData(&pagealiasedata, spaceid, s.Authority.DB) // parentid != 0 and groupid = 0

	for _, result := range pagealiasedata {

		var newgroupid int

		if result.PageGroupId != 0 {

			var pagesgroupal TblPagesGroupAliases

			SP.GetDetailsInPageAlia(&pagesgroupal, result.PageGroupId, spaceid, s.Authority.DB) // parentid != 0 and groupid = 0

			var pageal TblPagesGroupAliases

			SP.GetDetailsInPageAli(&pageal, pagesgroupal.GroupName, tblspaces.Id, s.Authority.DB)

			newgroupid = pageal.PageGroupId

		}

		var pagealid TblPageAliases

		SP.AliasesInParentId(&pagealid, result.ParentId, spaceid, s.Authority.DB)

		var pageali TblPageAliases

		SP.LastLoopAliasesInPage(&pageali, pagealid.PageTitle, tblspaces.Id, s.Authority.DB)

		var page TblPage

		page.SpacesId = tblspaces.Id

		page.PageGroupId = newgroupid

		page.ParentId = pageali.PageId

		pagealiid, _ := SP.ClonePage(&page, s.Authority.DB)

		// var pagesali TblPageAliases

		pagesali := result

		pagesali.PageId = pagealiid.Id

		SP.ClonePages(&pagesali, s.Authority.DB)

	}

	return nil

}

func (s Space) PageCategoryList() []Arrangecategories {

	var getallparentcat []TblPagesCategoriesAliases

	SP.PageParentCategoryList(&getallparentcat, s.Authority.DB)

	var AllCategorieswithSubCategories []Arrangecategories

	for _, Group := range getallparentcat {

		GetData, _ := SP.GetPageCategoryTree(Group.PageCategoryId, s.Authority.DB)

		var pid int

		for _, categories := range GetData {

			var addcat Arrangecategories

			var individualid []CatgoriesOrd

			pid = categories.ParentId

		LOOP:
			for _, GetParent := range GetData {

				var indivi CatgoriesOrd

				if pid == GetParent.PageCategoryId {

					pid = GetParent.ParentId

					indivi.Id = GetParent.PageCategoryId

					indivi.Category = GetParent.CategoryName

					individualid = append(individualid, indivi)

					if pid != 0 {

						goto LOOP

					}
				}

			}

			var ReverseOrder Arrangecategories

			addcat.Categories = append(addcat.Categories, individualid...)

			var singlecat []CatgoriesOrd

			for i := len(addcat.Categories) - 1; i >= 0; i-- {

				var Sing CatgoriesOrd

				Sing.Id = addcat.Categories[i].Id

				Sing.Category = addcat.Categories[i].Category

				singlecat = append(singlecat, Sing)

			}

			var newcate CatgoriesOrd

			newcate.Id = categories.PageCategoryId

			newcate.Category = categories.CategoryName

			addcat.Categories = append(addcat.Categories, newcate)

			singlecat = append(singlecat, newcate)

			ReverseOrder.Categories = singlecat

			AllCategorieswithSubCategories = append(AllCategorieswithSubCategories, ReverseOrder)
		}

	}

	/*This for Channel category show also individual group*/
	var FinalCategoryList []Arrangecategories

	for _, val := range AllCategorieswithSubCategories {

		if len(val.Categories) > 1 {

			var infinalarray Arrangecategories

			infinalarray.Categories = append(infinalarray.Categories, val.Categories...)

			FinalCategoryList = append(FinalCategoryList, infinalarray)
		}

	}

	return FinalCategoryList

}

// Check Name is already exits or not
func (s Space) CheckSpaceName(id int, name string) (bool, error) {

	_, _, checkerr := auth.VerifyToken(s.Authority.Token, s.Authority.Secret)

	if checkerr != nil {

		return false, checkerr
	}

	check, err := s.Authority.IsGranted("Spaces", authcore.CRUD)

	if err != nil {

		return false, err
	}

	if check {

		var space TblSpacesAliases

		err := SP.CheckSpaceName(&space, id, name, s.Authority.DB)

		if err != nil {
			return false, err
		}
		if space.Id == 0 {

			return false, err
		}

	}
	return true, nil
}

/*Get Published spaces*/
func (s Space) GetPublishedSpaces(limit int, offset int, filter Filter, languageid int) (spacedetails []TblSpacesAliases, err error) {

	_, _, checkerr := auth.VerifyToken(s.Authority.Token, s.Authority.Secret)

	if checkerr != nil {

		return []TblSpacesAliases{}, checkerr
	}

	check, err := s.Authority.IsGranted("Spaces", authcore.CRUD)

	if err != nil {

		return []TblSpacesAliases{}, err
	}

	if check {

		var spacez []TblSpacesAliases

		SP.PublishPageSpaceList(&spacez, languageid, limit, offset, filter, nil, s.Authority.DB)

		var SpaceDetails []TblSpacesAliases

		for _, space := range spacez {

			var child_page_Category categories.TblCategory

			_, child_page := categories.GetChildPageCategoriess(&child_page_Category, space.PageCategoryId, s.Authority.DB)

			var categorynames []categories.TblCategory

			var flg int

			categorynames = append(categorynames, child_page)

			flg = child_page.ParentId

			if flg != 0 {

			CLOOP:

				var count int //for safe

				for {

					count = count + 1 //for safe

					if count == 200 { //for safe

						break
					}

					var newchildcategory categories.TblCategory

					_, child := categories.GetChildPageCategoriess(&newchildcategory, flg, s.Authority.DB)

					flg = child.ParentId

					if flg != 0 {

						categorynames = append(categorynames, child)

						goto CLOOP

					} else {

						categorynames = append(categorynames, child)

						break
					}

				}

			}

			var reverseCategoryOrder []categories.TblCategory

			for i := len(categorynames) - 1; i >= 0; i-- {

				reverseCategoryOrder = append(reverseCategoryOrder, categorynames[i])

			}

			space.CategoryNames = reverseCategoryOrder

			SpaceDetails = append(SpaceDetails, space)

		}
		return SpaceDetails, nil

	}

	return []TblSpacesAliases{}, errors.New("not authorized")
}
