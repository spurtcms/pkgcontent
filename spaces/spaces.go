package spaces

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	authcore "github.com/spurtcms/spurtcms-core/auth"
	membercore "github.com/spurtcms/spurtcms-core/member"
	memberaccore "github.com/spurtcms/spurtcms-core/memberaccess"
	"gorm.io/gorm"
)

type Space struct {
	Authority *authcore.Authority
}

var s Space

type MemberSpace struct {
	MemAuth *authcore.Authority
}

func MigrateTable(db *gorm.DB) {

	db.AutoMigrate(
		&TblSpaces{},
		&TblSpacesAliases{},
		&TblPagesCategories{},
		&TblPagesCategoriesAliases{},
		&TblLanguage{},
	)

	db.Exec(`INSERT INTO PUBLIC.TBL_LANGUAGE(ID,LANGUAGE_NAME,LANGUAGE_CODE,JSON_PATH,IS_STATUS,IS_DEFAULT,	CREATED_BY,CREATED_ON,IS_DELETED) VALUES (1,'English', 'en', 'locales/en.json', 1, 1,1, '2023-09-11 11:27:44',0)`)
}

var IST, _ = time.LoadLocation("Asia/Kolkata")

/*spacelist*/
func (s Space) SpaceList(limit, offset int, filter Filter) (tblspace []TblSpacesAliases, totalcount int64, err error) {

	log.Println()

	_, _, checkerr := authcore.VerifyToken(s.Authority.Token, s.Authority.Secret)

	if checkerr != nil {

		return []TblSpacesAliases{}, 0, checkerr
	}

	var default_lang TblLanguage

	GetDefaultLanguage(&default_lang, s.Authority.DB)

	var spaces []TblSpacesAliases

	_, spaceerr := SpaceList(&spaces, default_lang.Id, limit, offset, filter, []int{}, s.Authority.DB)

	if spaceerr != nil {

		return []TblSpacesAliases{}, 0, spaceerr
	}

	var SpaceDetails []TblSpacesAliases

	for _, space := range spaces {

		var parent_page_Category TblPagesCategoriesAliases

		_, parent_page := GetParentPageCategory(&parent_page_Category, space.ParentId,s.Authority.DB)

		space.ParentCategory = parent_page

		if parent_page.Id != 0 {

			var child_page_Categories []TblPagesCategoriesAliases

			_, child_page := GetChildPageCategories(&child_page_Categories, space.PageCategoryId,s.Authority.DB)

			for _, child_category := range child_page {

				space.ChildCategories = append(space.ChildCategories, child_category)
			}

		}

		space.CreatedDate = space.CreatedOn.Format("02 Jan 2006 03:04 PM")

		if !space.ModifiedOn.IsZero() {

			space.ModifiedDate = space.ModifiedOn.Format("02 Jan 2006 03:04 PM")

		} else {

			space.ModifiedDate = space.CreatedOn.Format("02 Jan 2006 03:04 PM")

		}
		SpaceDetails = append(SpaceDetails, space)

	}

	count, _ := SpaceList(&spaces, default_lang.Id, 0, 0, filter, []int{}, s.Authority.DB)

	return SpaceDetails, count, nil

}

/*spacelist*/
func (s MemberSpace) MemberSpaceList(limit, offset int, filter Filter) (tblspace []TblSpacesAliases, totalcount int64, err error) {

	_, _, checkerr := membercore.VerifyToken(s.MemAuth.Token, s.MemAuth.Secret)

	if checkerr != nil {

		return []TblSpacesAliases{}, 0, checkerr
	}

	var mem memberaccore.AccessAuth

	mem.Authority = *s.MemAuth

	spceid, err := mem.GetSpace()

	if err != nil {

		log.Println(err)
	}

	var default_lang TblLanguage

	GetDefaultLanguage(&default_lang, s.MemAuth.DB)

	var spaces []TblSpacesAliases

	_, spaceerr := SpaceList(&spaces, default_lang.Id, limit, offset, filter, spceid, s.MemAuth.DB)

	if spaceerr != nil {

		return []TblSpacesAliases{}, 0, spaceerr
	}

	var SpaceDetails []TblSpacesAliases

	for _, space := range spaces {

		var parent_page_Category TblPagesCategoriesAliases

		_, parent_page := GetParentPageCategory(&parent_page_Category, space.ParentId,s.MemAuth.DB)

		space.ParentCategory = parent_page

		if parent_page.Id != 0 {

			var child_page_Categories []TblPagesCategoriesAliases

			_, child_page := GetChildPageCategories(&child_page_Categories, space.PageCategoryId,s.MemAuth.DB)

			for _, child_category := range child_page {

				space.ChildCategories = append(space.ChildCategories, child_category)
			}

		}

		space.CreatedDate = space.CreatedOn.Format("02 Jan 2006 03:04 PM")

		if !space.ModifiedOn.IsZero() {

			space.ModifiedDate = space.ModifiedOn.Format("02 Jan 2006 03:04 PM")

		} else {

			space.ModifiedDate = space.CreatedOn.Format("02 Jan 2006 03:04 PM")

		}
		SpaceDetails = append(SpaceDetails, space)

	}

	count, _ := SpaceList(&spaces, default_lang.Id, 0, 0, filter, spceid, s.MemAuth.DB)

	return SpaceDetails, count, nil

}

// create space
func (s Space) SpaceCreation(c *http.Request) error {

	userid, _, checkerr := authcore.VerifyToken(s.Authority.Token, s.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	_, err := s.Authority.IsGranted("Spaces", authcore.CRUD)

	if err != nil {

		return err
	}

	if c.PostFormValue("spacename") == "" || c.PostFormValue("spacedescription") == "" {

		return errors.New("given some values is empty")
	}

	var spaces TblSpaces

	spaces.PageCategoryId, _ = strconv.Atoi(c.PostFormValue("spacecategory"))

	spaces.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	spaces.CreatedBy = userid

	if err := s.Authority.DB.Table("tbl_spaces").Create(&spaces).Error; err != nil {

		return err
	}

	var space TblSpacesAliases

	space.SpacesName = c.PostFormValue("spacename")

	space.SpacesDescription = c.PostFormValue("spacedescription")

	space.ImagePath = c.PostFormValue("spaceimagepath")

	space.LanguageId, _ = strconv.Atoi(c.PostFormValue("langid"))

	space.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	space.CreatedBy = userid

	space.SpacesSlug = strings.ToLower(space.SpacesName)

	space.SpacesId = spaces.Id

	if err := s.Authority.DB.Table("tbl_spaces_aliases").Create(&space).Error; err != nil {

		return err
	}

	return nil
}

// update space
func (s Space) SpaceUpdate(c *http.Request) error {

	userid, _, checkerr := authcore.VerifyToken(s.Authority.Token, s.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	_, err := s.Authority.IsGranted("Spaces", authcore.CRUD)

	if err != nil {

		return err
	}

	if c.PostFormValue("spacename") == "" || c.PostFormValue("spacedescription") == "" {

		return nil
	}

	var space TblSpaces

	var spaces TblSpacesAliases

	id, _ := strconv.Atoi(c.PostFormValue("id"))

	spaces.Id = id

	space.Id = id

	spaces.SpacesName = c.PostFormValue("spacename")

	spaces.SpacesDescription = c.PostFormValue("spacedescription")

	spaces.SpacesSlug = strings.ToLower(spaces.SpacesName)

	space.PageCategoryId, _ = strconv.Atoi(c.PostFormValue("spacecategory"))

	spaces.ImagePath = c.PostFormValue("spaceimagepath")

	spaces.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	spaces.ModifiedBy = userid

	space.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	space.ModifiedBy = userid

	err1 := EditSpace(&spaces, id,s.Authority.DB)

	if err != nil {

		return err1
	}

	err2 := UpdateSpace(&space, id,s.Authority.DB)

	if err2 != nil {

		return err2
	}

	return nil
}

/*Delete Space*/
func (s Space) DeleteSpace(spaceid int) error {

	userid, _, checkerr := authcore.VerifyToken(s.Authority.Token, s.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	_, err := s.Authority.IsGranted("Spaces", authcore.CRUD)

	if err != nil {

		return err
	}

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

	err1 := DeleteSpaceAliases(&spaces, spaceid,s.Authority.DB)

	if err1 != nil {
		return err
	}

	err2 := DeleteSpace(&space, spaceid,s.Authority.DB)

	if err2 != nil {
		return err2
	}


	var page []TblPage

	GetPageDetailsBySpaceId(&page, spaceid,s.Authority.DB)

	var pid []int

	for _, v := range page {

		pid = append(pid, v.Id)

	}

	var pg TblPage

	pg.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	pg.DeletedBy = userid

	pg.IsDeleted = 1

	DeletePageInSpace(&pg, pid,s.Authority.DB)

	var pgali TblPageAliases

	pgali.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	pgali.DeletedBy = userid

	pgali.IsDeleted = 1

	DeletePageAliInSpace(&pgali, pid,s.Authority.DB)

	return nil
}

// Clone

func (s Space) CloneSpace(c *http.Request) error {

	userid, _, checkerr := authcore.VerifyToken(s.Authority.Token, s.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	if c.PostFormValue("spacename") == "" || c.PostFormValue("spacedescription") == "" {

		return errors.New("given some values is empty")
	}

	var spaces TblSpaces

	var space TblSpacesAliases

	space.SpacesName = c.PostFormValue("spacename")

	space.SpacesDescription = c.PostFormValue("spacedescription")

	space.SpacesSlug = strings.ToLower(space.SpacesName)

	spaces.PageCategoryId, _ = strconv.Atoi(c.PostFormValue("spacecategory"))

	space.ImagePath = c.PostFormValue("spaceimagepath")

	space.LanguageId, _ = strconv.Atoi(c.PostFormValue("langid"))

	spaceid, _ := strconv.Atoi(c.PostFormValue("id"))

	spaces.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	spaces.CreatedBy = userid

	space.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

	space.CreatedBy = userid

	tblspaces, _ := CreateSpace(&spaces,s.Authority.DB)

	space.SpacesId = tblspaces.Id

	err := CreateSpacesAliases(&space,s.Authority.DB)

	if err != nil {

		log.Println(err)
	}

	var pagegroupdata []TblPagesGroupAliases

	GetPageGroupData(&pagegroupdata, spaceid,s.Authority.DB)

	for _, value := range pagegroupdata {

		var group TblPagesGroup

		group.SpacesId = tblspaces.Id

		groups, _ := CloneSpaceInPagesGroup(&group,s.Authority.DB)

		var pagegroup TblPagesGroupAliases

		pagegroup = value

		pagegroup.PageGroupId = groups.Id

		ClonePagesGroup(&pagegroup,s.Authority.DB)
	}

	var pageId []TblPageAliases

	GetPageInPage(&pageId, spaceid,s.Authority.DB) //parentid 0 and groupid 0

	for _, val := range pageId {

		var page TblPage

		page.SpacesId = tblspaces.Id

		page.PageGroupId = 0

		page.ParentId = 0

		pageid, _ := ClonePage(&page,s.Authority.DB)

		var pagesali TblPageAliases

		pagesali = val

		pagesali.PageId = pageid.Id

		ClonePages(&pagesali,s.Authority.DB)

	}

	var pagegroupaldata TblPagesGroupAliases

	GetIdInPage(&pagegroupaldata, spaceid,s.Authority.DB) // parentid = 0 and groupid != 0

	var pagealiase []TblPageAliases

	GetPageAliasesInPage(&pagealiase, spaceid,s.Authority.DB) // parentid = 0 and groupid != 0

	for _, value := range pagealiase {

		var pageal TblPagesGroupAliases

		GetDetailsInPageAli(&pageal, pagegroupaldata.GroupName, tblspaces.Id,s.Authority.DB)

		// var parent TblPage

		// ParentWithChild(&parent, value.PageGroupId, spaceid)

		var page TblPage

		page.SpacesId = tblspaces.Id

		page.PageGroupId = pageal.PageGroupId

		page.ParentId = 0

		pagess, _ := ClonePage(&page,s.Authority.DB)

		var pagesali TblPageAliases

		pagesali = value

		pagesali.PageId = pagess.Id

		ClonePages(&pagesali,s.Authority.DB)

	}

	var pagealiasedata []TblPageAliases

	GetPageAliasesInPageData(&pagealiasedata, spaceid,s.Authority.DB) // parentid != 0 and groupid = 0

	for _, result := range pagealiasedata {

		var newgroupid int

		if result.PageGroupId != 0 {

			var pagesgroupal TblPagesGroupAliases

			GetDetailsInPageAlia(&pagesgroupal, result.PageGroupId, spaceid,s.Authority.DB) // parentid != 0 and groupid = 0

			var pageal TblPagesGroupAliases

			GetDetailsInPageAli(&pageal, pagesgroupal.GroupName, tblspaces.Id,s.Authority.DB)

			newgroupid = pageal.PageGroupId

		}

		var pagealid TblPageAliases

		AliasesInParentId(&pagealid, result.ParentId, spaceid,s.Authority.DB)

		var pageali TblPageAliases

		LastLoopAliasesInPage(&pageali, pagealid.PageTitle, tblspaces.Id,s.Authority.DB)

		var page TblPage

		page.SpacesId = tblspaces.Id

		page.PageGroupId = newgroupid

		page.ParentId = pageali.PageId

		pagealiid, _ := ClonePage(&page,s.Authority.DB)

		var pagesali TblPageAliases

		pagesali = result

		pagesali.PageId = pagealiid.Id

		ClonePages(&pagesali,s.Authority.DB)

	}

	return nil

}

func PageCategoryList() []Arrangecategories {

	var getallparentcat []TblPagesCategoriesAliases

	PageParentCategoryList(&getallparentcat,s.Authority.DB)

	var AllCategorieswithSubCategories []Arrangecategories

	for _, Group := range getallparentcat {

		GetData, _ := GetPageCategoryTree(Group.PageCategoryId,s.Authority.DB)

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
