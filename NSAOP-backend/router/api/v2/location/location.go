package location

import (
	"log"
	"net/http"
	"nsaop/model"
	"nsaop/router/resp"
	"nsaop/utils"
	"nsaop/utils/constant"

	"github.com/gin-gonic/gin"
)

type LocationFilter struct {
	Offset int    `json:"offset" form:"offset"`
	Limit  int    `json:"limit" form:"limit"`
	Search string `json:"search" form:"search"`
}

type LocationInfo struct {
	ID      uint   `json:"id"`
	Comment string `json:"comment"`
	Address string `json:"address"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
}

type LocationArray struct {
	Locations []LocationInfo `json:"locations"`
	Count     int            `json:"count"`
}

var inwhich = []string{"comment", "address"}

type LocationNew struct {
	Comment string `json:"comment"`
	Address string `json:"address"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	ReCAPTCHA	string `json:"g_recaptcha_response"`
}

func Location(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodPost:
		NewLocation(c)
	case http.MethodGet:
		GetLocationByFilter(c)
	default:
		resp.ERROR(c, http.StatusMethodNotAllowed, constant.MsgMethodNotAllowed)
	}
}

// @Tags		Location
// @Summary		Create new location
// @Description	with jwt auth, new a location with informations
// @Description	**must fill:**
//@Description  <pre>{
// @Description   comment, // readible comment for customer to remember, len 0~10 chinese
// @Description   address, // location address, len 0~100 chinese
// @Description   contact, // contact name, for privacy reason, could be "赵先生", len 0~20 chinese
// @Description   phone, // contact phone
// @Description }</pre>
// @Description	*can be null:* {
// @Description }
// @Accept		json
// @Product		json
// @Param		g-recaptcha-response query string true "response token got from google reCAPTCHA (expected action: create_location)"
// @Param		Authorization header string true "Bearer [access_token]"
// @Param		body body LocationNew true "location information"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "reCAPTCHA failed<br>format error<br>field empty<br>invalid info<br>permission denied"
// @Router		/location [POST]
func NewLocation(c *gin.Context) {
	var location LocationNew
	if err := c.ShouldBindJSON(&location); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	if reCAPTCHASucceed, err := utils.ReCAPTCHA(location.ReCAPTCHA, "create_location", c.ClientIP()); !reCAPTCHASucceed {
		log.Println(err)
		resp.ERROR(c, http.StatusForbidden, "reCAPTCHA failed")
		return
	}
	if location.Comment == "" || location.Address == "" ||
		location.Contact == "" || location.Phone == "" {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFieldEmpty)
		return
	}
	if !model.ValidCommentLen(location.Comment) || !model.ValidAddressLen(location.Address) ||
		!model.ValidContactLen(location.Contact) || !model.ValidPhoneLen(location.Phone) {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}

	var u model.User
	model.DB.First(&u, c.MustGet("userId").(uint))

	if u.Role != constant.RoleCustomer {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgPermissionDenied)
		return
	}

	l := model.Location{
		Comment: location.Comment,
		Address: location.Address,
		Contact: location.Contact,
		Phone:   location.Phone,
		UserID:  u.ID,
	}
	if err := model.DB.Create(&l).Error; err != nil {
		resp.ERROR(c, http.StatusInternalServerError, err.Error())
	} else {
		resp.OK(c, http.StatusOK, gin.H{})
	}
}

// @Tags		Location
// @Summary		Get location by filter
// @Description	with jwt auth and filter conditions, get user's locations
// @Description **must fill:**
// @Description <pre> {
// @Description 	offset,
// @Description 	limit, // cannot be 0, should be <= 20
// @Descriotion 	// in range [offset, offset+limit)
// @Description } </pre>
// @Description *can be null:*
// @Description <pre> {
// @Description 	search, // key words seperated by space, len<=10, search by OR method // support address & comment
// @Description } </pre>
// @Description	return {
// @Description   "count": count,
// @Description   "locations": [{
// @Description   	id, // location's id
// @Description   	comment, // readible comment for customer to remember
// @Description 	address, // location address
// @Description		contact, // contact name, for privacy reason
// @Description 	phone, // contact phone
// @Description   }], // array of locations order by id
// @Description	}</pre>
// @Accept		json
// @Product		json
// @Param 		Authorization header string true "Bearer [access_token]"
// @Param 		offset query int true "offset"
// @Param 		limit query int true "0<limit<=20"
// @Param 		search query string false "len<=10"
// @Success		200 {object} resp.Response{data=[]LocationInfo} "ok"
// @Failure		400 {object} resp.Response "format error"
// @Router		/location [GET]
func GetLocationByFilter(c *gin.Context) {
	var filter LocationFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	if filter.Limit == 0 || filter.Limit > 20 {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}
	if len(filter.Search) > 30 {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}

	var user model.User
	var locations []model.Location
	var err error
	whereClause := model.DB
	if whereClause, err = utils.GenerateSearchCLause(whereClause, filter.Search, inwhich); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
	}
	if err := model.DB.First(&user, c.MustGet("userId").(uint)).Error; err != nil {
		resp.ERROR(c, http.StatusInternalServerError, constant.MsgTokenCannotFindUser)
		return
	}
	model.DB.Model(&user).Limit(filter.Limit).Offset(filter.Offset).
		Where(whereClause).Association(constant.TableLocation).Find(&locations)
	count := int(model.DB.Model(&user).Where(whereClause).Association(constant.TableLocation).Count())
	var ret = []LocationInfo{}
	for _, loc := range locations {
		ret = append(ret, LocationInfo{
			ID:      loc.ID,
			Comment: loc.Comment,
			Address: loc.Address,
			Contact: loc.Contact,
			Phone:   loc.Phone,
		})
	}
	resp.OK(c, http.StatusOK, LocationArray{
		Locations: ret,
		Count:     count,
	})
}
