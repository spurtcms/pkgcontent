package lms

import (
	"github.com/spurtcms/pkgcore/auth"
)

/*this struct holds dbconnection ,token*/
type SpaceCategory struct {
	Authority *auth.Authorization
}

type Authstruct struct{}

var C Authstruct

/*List Category Group*/
func (c SpaceCategory) CategoryGroupList(limit int, offset int, filter Filter) (Categorylist []TblSpaceCategory, categorycount int64, err error) {

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
