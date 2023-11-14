package categories

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spurtcms/spurtcms-core/auth"
	"gorm.io/gorm"
)

var IST, _ = time.LoadLocation("Asia/Kolkata")

type Category struct {
	Authority *auth.Authorization
}

func MigrateTable(DB *gorm.DB) {

	DB.AutoMigrate(
		&TblCategory{},
	)

}

/*List Category Group*/
func (c Category) CategoryGroupList(limit int, offset int, filter Filter) (Categorylist []TblCategory, categorycount int64, err error) {

	_, _, checkerr := auth.VerifyToken(c.Authority.Token, c.Authority.Secret)

	if checkerr != nil {

		return []TblCategory{}, 0, checkerr
	}

	check, err := c.Authority.IsGranted("Categories Group", auth.Read)

	if err != nil {

		return []TblCategory{}, 0, err
	}

	if check {

		var categorylist []TblCategory

		_, Total_categories := GetCategoryList(categorylist, 0, 0, filter, c.Authority.DB)

		categorygrplist, _ := GetCategoryList(categorylist, offset, limit, filter, c.Authority.DB)

		return categorygrplist, Total_categories, nil

	}

	return []TblCategory{}, 0, errors.New("not authorized")
}

/*Add Category Group*/
func (c Category) CreateCategoryGroup(req CategoryCreate) error {

	userid, _, checkerr := auth.VerifyToken(c.Authority.Token, c.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := c.Authority.IsGranted("Categories Group", auth.Create)

	if err != nil {

		return err
	}

	if check {

		if req.CategoryName == "" || req.Description == "" {

			return errors.New("given some values is empty")
		}

		var category TblCategory

		category.CategoryName = req.CategoryName

		category.CategorySlug = strings.ToLower(req.CategoryName)

		category.Description = req.Description

		category.CreatedBy = userid

		category.ParentId = 0

		category.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		err := CreateCategory(&category, c.Authority.DB)

		if err != nil {

			return err
		}

		return nil

	}

	return errors.New("not authorized")
}

/*UpdateCategoryGroup*/
func (c Category) UpdateCategoryGroup(req CategoryCreate) error {

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
		var category TblCategory

		category.Id = req.Id

		category.CategoryName = req.CategoryName

		category.Description = req.Description

		category.CategorySlug = strings.ToLower(req.CategoryName)

		category.ModifiedBy = userid

		category.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		err := UpdateCategory(&category, c.Authority.DB)

		if err != nil {

			return err
		}

		return nil
	}

	return errors.New("not authorized")
}

/*DeleteCategoryGroup*/
func (c Category) DeleteCategoryGroup(Categoryid int) error {

	userid, _, checkerr := auth.VerifyToken(c.Authority.Token, c.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := c.Authority.IsGranted("Categories Group", auth.Delete)

	if err != nil {

		return err
	}

	if check {

		var category TblCategory

		category.DeletedBy = userid

		category.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		category.IsDeleted = 1

		err := DeleteCategoryById(&category, Categoryid, c.Authority.DB)

		if err != nil {

			return err
		}

		return nil

	}

	return errors.New("not authorized")
}

/*ListCategory*/
func (c Category) ListCategory(offset int, limit int, filter Filter, parent_id int) (tblcat []TblCategory, category []TblCategory, parentcategory TblCategory, categorycount int64, err error) {

	_, _, checkerr := auth.VerifyToken(c.Authority.Token, c.Authority.Secret)

	if checkerr != nil {

		return []TblCategory{}, []TblCategory{}, TblCategory{}, 0, checkerr
	}

	check, err := c.Authority.IsGranted("Categories", auth.Read)

	if err != nil {

		return []TblCategory{}, []TblCategory{}, TblCategory{}, 0, err
	}

	if check {

		var categorylist []TblCategory

		var categorylists []TblCategory

		var categorys []TblCategory

		var category TblCategory

		parentcategory, err1 := GetCategoryById(&category, parent_id, c.Authority.DB)

		if err1 != nil {
			fmt.Println(err)
		}
		_, count := GetSubCategoryList(&categorylist, 0, 0, filter, parent_id, 0, c.Authority.DB)

		fmt.Println("d", count)

		childcategorys, _ := GetSubCategoryList(&categorys, offset, limit, filter, parent_id, 1, c.Authority.DB)

		childcategory, _ := GetSubCategoryList(&categorylist, offset, limit, filter, parent_id, 0, c.Authority.DB)

		for _, val := range *childcategory {

			if !val.ModifiedOn.IsZero() {

				val.DateString = val.ModifiedOn.Format("02 Jan 2006 03:04 PM")

			} else {
				val.DateString = val.CreatedOn.Format("02 Jan 2006 03:04 PM")

			}

			categorylists = append(categorylists, val)

		}
		var AllCategorieswithSubCategories []Arrangecategories

		GetData, _ := GetCategoryTree(parent_id, c.Authority.DB)

		var pid int

		for _, categories := range GetData {

			var addcat Arrangecategories

			var individualid []CatgoriesOrd

			pid = categories.ParentId

		LOOP:
			for _, GetParent := range GetData {

				var indivi CatgoriesOrd

				if pid == GetParent.Id {

					pid = GetParent.ParentId

					indivi.Id = GetParent.Id

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

			newcate.Id = categories.Id

			newcate.Category = categories.CategoryName

			addcat.Categories = append(addcat.Categories, newcate)

			singlecat = append(singlecat, newcate)

			ReverseOrder.Categories = singlecat

			AllCategorieswithSubCategories = append(AllCategorieswithSubCategories, ReverseOrder)

		}

		var FinalCategoryList []Arrangecategories

		for _, val := range AllCategorieswithSubCategories {

			var infinalarray Arrangecategories

			for index, res := range val.Categories {

				if index < len(val.Categories)-1 {

					// var cate CatgoriesOrd

					cate := res

					infinalarray.Categories = append(infinalarray.Categories, cate)

				}

			}
			FinalCategoryList = append(FinalCategoryList, infinalarray)
		}

		var FinalModalCategoryList []Arrangecategories

		for _, val := range AllCategorieswithSubCategories {

			var infinalarray Arrangecategories

			for index, res := range val.Categories {

				if index < len(val.Categories) {

					// var cate CatgoriesOrd

					cate := res

					infinalarray.Categories = append(infinalarray.Categories, cate)
				}
			}
			FinalModalCategoryList = append(FinalModalCategoryList, infinalarray)
		}

		var FinalModalCategoriesList []TblCategory

		for index, val := range *childcategorys {

			// var finalcat TblCategory

			finalcat := val

			for cindex, val2 := range FinalModalCategoryList {

				if index == cindex {

					for _, va3 := range val2.Categories {

						finalcat.Parent = append(finalcat.Parent, va3.Category)
					}
				}
			}
			FinalModalCategoriesList = append(FinalModalCategoriesList, finalcat)
		}
		var FinalCategoriesList []TblCategory

		for index, val := range categorylists {

			// var finalcat TblCategory

			finalcat := val

			for cindex, val2 := range FinalCategoryList {

				if index == cindex {

					for _, va3 := range val2.Categories {

						finalcat.Parent = append(finalcat.Parent, va3.Category)
					}
				}
			}
			FinalCategoriesList = append(FinalCategoriesList, finalcat)
		}

		return FinalCategoriesList, FinalModalCategoriesList, parentcategory, count, nil
	}

	return []TblCategory{}, []TblCategory{}, TblCategory{}, 0, errors.New("not authorized")
}

/*Add Category*/
func (c Category) AddCategory(req CategoryCreate) error {

	userid, _, checkerr := auth.VerifyToken(c.Authority.Token, c.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := c.Authority.IsGranted("Categories", auth.Create)

	if err != nil {

		return err
	}

	if check {
		if req.CategoryName == "" || req.Description == "" {

			return errors.New("given some values is empty")
		}

		var category TblCategory

		category.CategoryName = req.CategoryName

		category.CategorySlug = strings.ToLower(req.CategoryName)

		category.Description = req.Description

		category.ImagePath = req.ImagePath

		category.CreatedBy = userid

		category.ParentId = req.ParentId

		// if req.FormValue("groupid") == "" {
		// 	category.ParentId, _ = strconv.Atoi(req.FormValue("categoryid"))
		// } else {
		// 	category.ParentId, _ = strconv.Atoi(req.FormValue("groupid"))
		// }

		category.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		err := CreateCategory(&category, c.Authority.DB)

		if err != nil {

			return err
		}

		return nil

	}

	return errors.New("not authorized")
}

/*Update Sub category*/
func (c Category) UpdateSubCategory(req CategoryCreate) error {

	userid, _, checkerr := auth.VerifyToken(c.Authority.Token, c.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := c.Authority.IsGranted("Categories", auth.Update)

	if err != nil {

		return err
	}

	if check {

		if req.Id <= 0 || req.CategoryName == "" || req.Description == "" {

			return errors.New("given some values is empty")
		}

		var category TblCategory

		category.CategoryName = req.CategoryName

		category.CategorySlug = strings.ToLower(req.CategoryName)

		category.Description = req.Description

		category.ImagePath = req.ImagePath

		category.ParentId = req.ParentId

		category.CreatedBy = userid

		category.Id = req.Id

		category.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		category.ModifiedBy = userid

		err := UpdateCategory(&category, c.Authority.DB)

		if err != nil {

			return err
		}

		return nil

	}

	return errors.New("not authorized")
}

/*Delete Sub Category*/
func (c Category) DeleteSubCategory(categoryid int) error {

	userid, _, checkerr := auth.VerifyToken(c.Authority.Token, c.Authority.Secret)

	if checkerr != nil {

		return checkerr
	}

	check, err := c.Authority.IsGranted("Categories", auth.Delete)

	if err != nil {

		return err
	}

	if check {

		var category TblCategory

		category.DeletedBy = userid

		category.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().In(IST).Format("2006-01-02 15:04:05"))

		category.IsDeleted = 1

		err := DeleteCategoryById(&category, categoryid, c.Authority.DB)

		if err != nil {

			return err
		}

		return nil
	}

	return errors.New("not authorized")
}

// Get Sub Category List

func (c Category) GetSubCategoryDetails(categoryid int) (categorys TblCategory, err error) {

	_, _, checkerr := auth.VerifyToken(c.Authority.Token, c.Authority.Secret)

	if checkerr != nil {

		return TblCategory{}, checkerr
	}

	check, err := c.Authority.IsGranted("Categories", auth.Delete)

	if err != nil {

		return TblCategory{}, err
	}

	if check {

		var category TblCategory

		category, err := GetCategoryDetails(categoryid, &category, c.Authority.DB)

		if err != nil {

			return TblCategory{}, err
		}

		return category, nil
	}

	return TblCategory{}, errors.New("not authorized")
}

// Filter Category
func (c Category) FilterSubCategory(limit int, filter Filter, parent_id int) (tblcat []TblCategory, categorycount int64, err error) {

	_, _, checkerr := auth.VerifyToken(c.Authority.Token, c.Authority.Secret)

	if checkerr != nil {

		return []TblCategory{}, 0, checkerr
	}

	check, err := c.Authority.IsGranted("Categories", auth.Read)

	if err != nil {

		return []TblCategory{}, 0, err
	}

	if check {

		var categorylist []TblCategory

		var categorylists []TblCategory

		var categorys []TblCategory

		var category TblCategory

		_, err1 := GetCategoryById(&category, parent_id, c.Authority.DB)

		if err1 != nil {
			fmt.Println(err)
		}

		childcategorys, _ := GetSubCategoryList(&categorys, 0, limit, filter, parent_id, 0, c.Authority.DB)

		_, count := GetSubCategoryList(&categorylist, 0, 0, filter, parent_id, 0, c.Authority.DB)

		childcategory, _ := GetSubCategoryList(&categorylist, 0, limit, filter, parent_id, 0, c.Authority.DB)

		for _, val := range *childcategory {

			if !val.ModifiedOn.IsZero() {

				val.DateString = val.ModifiedOn.Format("02 Jan 2006 03:04 PM")

			} else {
				val.DateString = val.CreatedOn.Format("02 Jan 2006 03:04 PM")

			}

			categorylists = append(categorylists, val)

		}
		var AllCategorieswithSubCategories []Arrangecategories

		GetData, _ := GetCategoryTree(parent_id, c.Authority.DB)

		var pid int

		for _, categories := range GetData {

			var addcat Arrangecategories

			var individualid []CatgoriesOrd

			pid = categories.ParentId

		LOOP:
			for _, GetParent := range GetData {

				var indivi CatgoriesOrd

				if pid == GetParent.Id {

					pid = GetParent.ParentId

					indivi.Id = GetParent.Id

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

			newcate.Id = categories.Id

			newcate.Category = categories.CategoryName

			addcat.Categories = append(addcat.Categories, newcate)

			singlecat = append(singlecat, newcate)

			ReverseOrder.Categories = singlecat

			AllCategorieswithSubCategories = append(AllCategorieswithSubCategories, ReverseOrder)

		}

		var FinalCategoryList []Arrangecategories

		for _, val := range AllCategorieswithSubCategories {

			var infinalarray Arrangecategories

			for index, res := range val.Categories {

				if index < len(val.Categories)-1 {

					cate := res

					infinalarray.Categories = append(infinalarray.Categories, cate)

				}

			}
			FinalCategoryList = append(FinalCategoryList, infinalarray)
		}

		var FinalModalCategoryList []Arrangecategories

		for _, val := range AllCategorieswithSubCategories {

			var infinalarray Arrangecategories

			for index, res := range val.Categories {

				if index < len(val.Categories) {

					cate := res

					infinalarray.Categories = append(infinalarray.Categories, cate)
				}
			}
			FinalModalCategoryList = append(FinalModalCategoryList, infinalarray)
		}

		var FinalModalCategoriesList []TblCategory

		for index, val := range *childcategorys {

			// var finalcat TblCategory

			finalcat := val

			for cindex, val2 := range FinalModalCategoryList {

				if index == cindex {

					for _, va3 := range val2.Categories {

						finalcat.Parent = append(finalcat.Parent, va3.Category)
					}
				}
			}
			FinalModalCategoriesList = append(FinalModalCategoriesList, finalcat)
		}
		var FinalCategoriesList []TblCategory

		for index, val := range categorylists {

			// var finalcat TblCategory

			finalcat := val

			for cindex, val2 := range FinalCategoryList {

				if index == cindex {

					for _, va3 := range val2.Categories {

						finalcat.Parent = append(finalcat.Parent, va3.Category)
					}
				}
			}
			FinalCategoriesList = append(FinalCategoriesList, finalcat)
		}

		return FinalCategoriesList, count, nil
	}

	return []TblCategory{}, 0, errors.New("not authorized")

}