package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/gin-gonic/gin"
	"github.com/leekchan/accounting"
)

type TableStatus struct {
	Name string
	ID   int
}

type ParentTableStatus struct {
	ID       int
	StatusID int
}

type StatusClient struct {
	Err    error
	Status int
	Jwt    string
	IP     [][]string
}

type LastID struct {
	ID uint64
}

type buffer struct {
	r         []byte
	runeBytes [utf8.UTFMax]byte
}

type FindAllConditionParam struct {
	FindAll      bool
	CountPage    int
	CountSize    int
	StartData    int
	ErrorStatus  int
	ErrorMessage string
	Error        error
}

const (
	serverKeyPushNotification = "AAAAZvT7Vs0:APA91bFs6wlz6vyM5GksKZ9Jdd00qrw4QrLVApsI9vdvaUoAFKwHR6Xszc_z1XQIabeZFPK5Ic0MUnttd2Ht3i0VPDRgK3IJmhl38762Cg7oFDbd1F659XYAukLqHE6BFOW4fF1nofSK"
	passwordSalt              = "a99VVoWzmd1C9ujcitK0fIVNE0I5I61AC47C852RoLTsHDyLCltvP+ZHEkIl/2hkzTOW90c3ZEjtYRkdfTWJ1Q=="
)

// find all
func FilterFindAll(c *gin.Context) (string, string) {
	page := c.Query("Page")
	size := c.Query("Size")
	if c.Query("Page") == "" {
		page = "-1"
	}
	if c.Query("Size") == "" {
		size = "10"
	}

	return page, size
}

// find all multifunction
func FilterFindAllParam(c *gin.Context) types.FindAllParams {
	var statusID string
	var businessID string
	var companyID string
	var outletID string
	var Outlets []string
	var sort string

	// userType := appcontext.Type(c)
	// if userType != nil {
	// 	if *userType != "Web" {
	// 		businessID = fmt.Sprintf("%d", appcontext.CompanyID(c))
	// 	}
	// }

	if c.Query("CompanyID") != "" {
		companyID = fmt.Sprintf("%v", c.Query("CompanyID"))
	}

	if c.Query("BusinessID") != "" {
		businessID = fmt.Sprintf("%v", c.Query("BusinessID"))
	}

	findallparams := types.FindAllParams{-1, 10, "", "code", "desc", "", "", "", Outlets}
	sortName := Underscore(c.Query("SortName"))
	sortBy := strings.ToLower(c.Query("SortBy"))

	if c.Query("SortName") == "" {
		sortName = "id"
	}

	if c.Query("SortBy") == "" {
		sortBy = "desc"
	}

	if c.Query("StatusID") == "" {
		statusID = c.Query("Status ID")
	} else {
		statusID = c.Query("StatusID")
	}

	explodeOutlets := strings.Split(outletID, ",")
	for _, vOutlet := range explodeOutlets {
		if vOutlet != "-1" && vOutlet != "" {
			Outlets = append(Outlets, vOutlet)
		}
	}

	explodeStatus := strings.Split(statusID, ",")
	for _, vStatus := range explodeStatus {
		if vStatus != "-1" && vStatus != "" {
			JoinStringStatus := strings.Join(explodeStatus, ",")
			statusID = "status_id IN (" + JoinStringStatus + ")"
			break
		} else {
			statusID = ""
			break
		}
	}

	cID, _ := MultiValueUUIDCheck(companyID) // make sure its all UUID
	if cID != "" {
		explodeCompany := strings.Split(cID, ",")
		for idx, b := range explodeCompany {
			if b != "-1" && b != "" && b != "0" {
				explodeCompany[idx] = fmt.Sprintf(`"%s"`, b)
			}
		}
		JoinStringCompany := strings.Join(explodeCompany, ",")
		companyID = "company_id IN (" + JoinStringCompany + ")"
	}

	bID, _ := MultiValueUUIDCheck(businessID) // make sure its all UUID
	if bID != "" {
		explodeBusiness := strings.Split(bID, ",")
		fmt.Println(">>", explodeBusiness)
		for idx, b := range explodeBusiness {
			if b != "-1" && b != "" && b != "0" {
				explodeBusiness[idx] = fmt.Sprintf(`"%s"`, b)
			}
		}
		JoinStringBusiness := strings.Join(explodeBusiness, ",")
		businessID = "business_id IN (" + JoinStringBusiness + ")"
	}

	if sortName != "" {
		sort = GetSortBy(sortName, sortBy)
	}

	dataFinder := DataFinder(c.Query("KeywordName"), c.Query("Keyword"))
	page, _ := strconv.Atoi(c.Query("Page"))
	size, _ := strconv.Atoi(c.Query("Size"))
	findallparams = types.FindAllParams{Page: page, Size: size, StatusID: statusID, DataFinder: dataFinder, SortName: sortName, SortBy: sort, BusinessID: businessID, CompanyID: companyID, Outlets: Outlets}
	return findallparams
}

func sanitize(text string) string {
	return strings.NewReplacer("'", "", `"`, "").Replace(text)
}

// keyword like full text search
func DataFinder(keywordname string, keyword string) string {
	str := "1=1"
	if keywordname != "" && keyword != "" {
		ExplodeParam := strings.Split(keywordname, ",")
		// ExplodeKeyword := strings.Split(keyword, " ")
		// for _, vKeyword := range ExplodeKeyword {
		str += " AND ( "
		strTmp := ""
		for _, vParam := range ExplodeParam {
			date := strings.Contains(vParam, "date")
			if date {
				t, errDate := time.Parse("2006-01-02", keyword)
				if errDate == nil {
					keyword = t.Format("2006-01-02")
				}

				t, errDate = time.Parse("02-01-2006", keyword)
				if errDate == nil {
					keyword = t.Format("2006-01-02")
				}
			}

			if strTmp != "" {
				strTmp += " OR "
			}

			strTmp += " " + sanitize(Underscore(vParam)) + " LIKE '%" + keyword + "%' "
		}
		str += strTmp
		str += " )"
	}
	// }

	return str
}

func GetSortBy(sortName string, sortBy string) string {
	var sort string
	var sortNameArr []string
	var sortByArr []string

	checkMultipleSortName := strings.Contains(sortName, ",")
	checkMultipleSortBy := strings.Contains(sortBy, ",")
	if checkMultipleSortName {
		explodeSortName := strings.Split(sortName, ",")
		for _, vSortName := range explodeSortName {
			sortNameArr = append(sortNameArr, vSortName)
		}
	} else {
		sortNameArr = append(sortNameArr, sortName)
	}

	if checkMultipleSortBy {
		explodeSortBy := strings.Split(sortBy, ",")
		for _, vSortBy := range explodeSortBy {
			sortByArr = append(sortByArr, vSortBy)
		}
	} else {
		sortByArr = append(sortByArr, sortBy)
	}

	for k, v := range sortNameArr {
		var str string
		lenSortBy := len(sortByArr)
		lenSortName := len(sortNameArr)
		if lenSortBy-1 >= k {
			str = v + " " + sortByArr[k]
		} else {
			str = v + " " + sortByArr[lenSortBy-1]
		}

		if lenSortName-1 != k {
			str = str + ","
		} else {
			str = str
		}
		sort = sort + str
	}

	return sort
}

func (b *buffer) write(r rune) {
	if r < utf8.RuneSelf {
		b.r = append(b.r, byte(r))
		return
	}
	n := utf8.EncodeRune(b.runeBytes[0:], r)
	b.r = append(b.r, b.runeBytes[0:n]...)
}

func (b *buffer) indent() {
	if len(b.r) > 0 {
		b.r = append(b.r, '_')
	}
}

func (b *buffer) indentSpace() {
	if len(b.r) > 0 {
		b.r = append(b.r, ' ')
	}
}

// set camelcase model name to table name with underscore
func Underscore(s string) string {
	b := buffer{
		r: make([]byte, 0, len(s)),
	}
	var m rune
	var w bool
	for _, ch := range s {
		if unicode.IsUpper(ch) {
			if m != 0 {
				if !w {
					b.indent()
					w = true
				}
				b.write(m)
			}
			m = unicode.ToLower(ch)
		} else if unicode.IsSpace(ch) {
			if m != 0 {
				b.indentSpace()
				m = 0
				w = false
			}
		} else {
			if m != 0 {
				b.indent()
				b.write(m)
				m = 0
				w = false
			}
			b.write(ch)
		}
	}
	if m != 0 {
		if !w {
			b.indent()
		}
		b.write(m)
	}

	// handle ID camel case
	strReplace := []byte(string(b.r))
	countID := strings.Count(string(strReplace), "i_d")
	if countID >= 1 {
		len := len(strReplace)
		for i := 0; i < len; i++ {
			if strReplace[i] == 'i' {
				if strReplace[i+1] == '_' {
					if strReplace[i+2] == 'd' {
						strReplace[i+1] = ' '
					}
				}
			}
		}
	}
	return strings.Replace(string(strReplace), " ", "", -1)
}

// format rupiah
func ConvertRupiah(value int, symbol bool) string {
	var strSymbol string
	if symbol {
		strSymbol = "Rp. "
	}
	ac := accounting.Accounting{Symbol: strSymbol, Precision: 2}

	Strings := ac.FormatMoney(value)

	return Strings
}
