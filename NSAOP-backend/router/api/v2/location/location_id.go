package location

import (
	"net/http"
	"nsaop/model"
	"nsaop/router/owned"
	"nsaop/router/resp"
	"nsaop/utils/constant"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LocationChange struct {
	Comment string `json:"comment"`
	Address string `json:address`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
}

func LocationId(c *gin.Context) {
	if id, err := strconv.ParseUint(c.Param("id"), 10, 32); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
	} else {
		var locations []model.Location
		if err := owned.GetUserOwnedDataByID(c, id, constant.TableLocation, &locations); err != nil {
			return
		}
		location := locations[0]
		switch c.Request.Method {
		case http.MethodGet:
			GetLocationId(c, location)
		case http.MethodPut:
			ChangeLocationId(c, location)
		case http.MethodDelete:
			DeleteLocationId(c, location)
		default:
			resp.ERROR(c, http.StatusMethodNotAllowed, constant.MsgMethodNotAllowed)
		}
	}
}

// @Tags		Location
// @Summary		Get location info by id
// @Description	with jwt auth, get user's locations by id
// @Description	return:
// @Description <pre> {
// @Description   	comment, // readible comment for customer to remember
// @Description 	address, // location address
// @Description		contact, // contact name, for privacy reason
// @Description 	phone, // contact phone
// @Description	}</pre>
// @Accept		json
// @Product		json
// @Param		id path int true "location id"
// @Param 		Authorization header string true "Bearer [access_token]"
// @Success		200 {object} resp.Response{data=LocationInfo} "ok"
// @Failure		400 {object} resp.Response "format error<br>id not available"
// @Router		/location/{id} [GET]
func GetLocationId(c *gin.Context, location model.Location) {
	resp.OK(c, http.StatusOK, LocationInfo{
		ID:      location.ID,
		Comment: location.Comment,
		Address: location.Address,
		Contact: location.Contact,
		Phone:   location.Phone,
	})
}

// @Tags		Location
// @Summary		Change location info by id
// @Description	with jwt auth, change user's location (by id)
// @Description *can be null:*
// @Description <pre>{
// @Description   comment, // readible comment for customer to remember
// @Description   address, // location address
// @Description   contact, // contact name, for privacy reason
// @Description   phone, // contact phone
// @Description }</pre>
// @Accept		json
// @Product		json
// @Param 		Authorization header string true "Bearer [access_token]"
// @Param 		id path int true "location id"
// @Param 		body body LocationChange true "fill what you need to modify, and don't fill what you don't want to modify"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "format error<br>invalid info<br>id not available"
// @Router		/location/{id} [PUT]
func ChangeLocationId(c *gin.Context, location model.Location) {
	var change LocationChange
	if err := c.ShouldBindJSON(&change); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	if (change.Comment != "" && !model.ValidCommentLen(change.Comment)) ||
		(change.Address != "" && !model.ValidAddressLen(change.Address)) ||
		(change.Contact != "" && !model.ValidContactLen(change.Contact)) ||
		(change.Phone != "" && !model.ValidPhoneLen(change.Phone)) {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}
	if change.Comment != "" {
		location.Comment = change.Comment
	}
	if change.Address != "" {
		location.Address = change.Address
	}
	if change.Contact != "" {
		location.Contact = change.Contact
	}
	if change.Phone != "" {
		location.Phone = change.Phone
	}
	model.DB.Save(&location)
	resp.OK(c, http.StatusOK, gin.H{})
}

// @Tags		Location
// @Summary		Delete location by id
// @Description	with jwt auth, delete user's unused location (by id)
// @Accept		json
// @Product		json
// @Param 		id path int true "location id"
// @Param 		Authorization header string true "Bearer [access_token]"
// @Param 		body body LocationChange true "location id"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "format error<br>id not available"
// @Router		/location/{id} [DELETE]
func DeleteLocationId(c *gin.Context, location model.Location) {
	model.DB.Delete(&location)
	resp.OK(c, http.StatusOK, gin.H{})
}
