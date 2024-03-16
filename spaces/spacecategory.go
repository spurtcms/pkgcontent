package lms

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spurtcms/pkgcore/auth"
)

/*this struct holds dbconnection ,token*/
type SpaceCategory struct {
	Authority *auth.Authorization
}

type Authstruct struct{}

var C Authstruct

/*List Category Group*/
func (c SpaceCategory) CategorySpaceGroupList(limit int, offset int, filter Filter) (Categorylist []TblSpaceCategory, categorycount int64, err error) {

	_, _, checkerr := auth.VerifyToken(c.Authority.Token, c.Authority.Secret)

	if checkerr != nil {

		return []TblSpaceCategory{}, 0, checkerr
	}

	// check, err := c.Authority.IsGranted("Categories Group", auth.Read)

	// if err != nil {

	// 	return []TblSpaceCategory{}, 0, err
	// }

	// if check {

	var categorylist []TblSpaceCategory

	_, Total_categories := C.GetSpaceCategoryList(categorylist, 0, 0, filter, c.Authority.DB)

	categorygrplist, _ := C.GetSpaceCategoryList(categorylist, offset, limit, filter, c.Authority.DB)

	var categorylists []TblSpaceCategory

	for _, val := range categorygrplist {

		if !val.ModifiedOn.IsZero() {

			val.DateString = val.ModifiedOn.Format("02 Jan 2006 03:04 PM")

		} else {
			val.DateString = val.CreatedOn.Format("02 Jan 2006 03:04 PM")

		}

		categorylists = append(categorylists, val)

	}

	return categorylists, Total_categories, nil

	// }

	// return []TblCategory{}, 0, errors.New("not authorized")
}

/*Add Category Group*/
func (c SpaceCategory) CreateCategoryGroup(req CategoryCreate) error {

	userid, _, checkerr := auth.VerifyToken(c.Authority.Token, c.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	// check, err := c.Authority.IsGranted("Categories Group", auth.Create)

	// if err != nil {

	// 	return err
	// }

	// if check {

		if req.CategoryName == "" || req.Description == "" {

			return errors.New("given some values is empty")
		}

		var category TblSpaceCategory

		category.CategoryName = req.CategoryName

		category.CategorySlug = strings.ToLower(req.CategoryName)

		category.Description = req.Description

		category.CreatedBy = userid

		category.ParentId = 0

		category.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		err := C.CreateSpaceCategory(&category, c.Authority.DB)

		if err != nil {

			return err
		}

		return nil

	// }

	// return errors.New("not authorized")
}

/*UpdateCategoryGroup*/
func (c SpaceCategory) UpdateCategoryGroup(req CategoryCreate) error {

	userid, _, checkerr := auth.VerifyToken(c.Authority.Token, c.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := c.Authority.IsGranted("Categories Group", auth.Update)

	if err != nil {

		return err
	}

	if check {

		if req.Id <= 0 || req.CategoryName == "" || req.Description == "" {

			return errors.New("given some values is empty")
		}
		var category TblSpaceCategory

		category.Id = req.Id

		category.CategoryName = req.CategoryName

		category.Description = req.Description

		category.CategorySlug = strings.ToLower(req.CategoryName)

		category.ModifiedBy = userid

		category.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		err := C.UpdateCategory(&category, c.Authority.DB)

		if err != nil {

			return err
		}

		return nil
	}

	return errors.New("not authorized")
}

/*DeleteCategoryGroup*/
func (c SpaceCategory) DeleteCategoryGroup(Categoryid int) error {

	userid, _, checkerr := auth.VerifyToken(c.Authority.Token, c.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	// check, err := c.Authority.IsGranted("Categories Group", auth.Delete)

	// if err != nil {

	// 	return err
	// }

	// if check {

		GetData, _ := C.GetCategoryTree(Categoryid, c.Authority.DB)

		var individualid []int

		for _, GetParent := range GetData {

			indivi := GetParent.Id

			individualid = append(individualid, indivi)

			fmt.Println(individualid, "categoryids")
		}

		spacecategory := individualid[0]

		// spacecategoryStr := fmt.Sprintf("%d", spacecategory)

		var category TblSpaceCategory

		category.DeletedBy = userid

		category.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		category.IsDeleted = 1

		err := C.DeleteallCategoryById(&category, individualid, spacecategory, c.Authority.DB)

		if err != nil {

			return err
		}

		return nil

	// }

	// return errors.New("not authorized")
}