package v2

import (
	"fmt"
	"net/http"

	"github.com/dappsteros-io/DappsterOS-Common/utils"
	"github.com/dappsteros-io/DappsterOS/codegen"
	"github.com/dappsteros-io/DappsterOS/common"
	"github.com/dappsteros-io/DappsterOS/pkg/utils/httper"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

func (s *DappsterOS) SetZerotierNetworkStatus(ctx echo.Context, networkId string) error {

	return ctx.JSON(http.StatusOK, nil)
}
func (s *DappsterOS) GetZerotierInfo(ctx echo.Context) error {
	info := codegen.GetZTInfoOK{}
	respBody, err := httper.ZTGet("/controller/network")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, codegen.BaseResponse{Message: utils.Ptr(err.Error())})
	}

	networkNames := gjson.ParseBytes(respBody).Array()
	for _, v := range networkNames {
		res, err := httper.ZTGet("/controller/network/" + v.Str)
		if err != nil {
			fmt.Println(err)
			return ctx.JSON(http.StatusInternalServerError, codegen.BaseResponse{Message: utils.Ptr(err.Error())})
		}
		name := gjson.GetBytes(res, "name").Str
		if name == common.RANW_NAME {
			via := gjson.GetBytes(res, "routes.0.via").Str
			info.Id = utils.Ptr(gjson.GetBytes(res, "id").Str)
			info.Name = &name
			if len(via) == 0 {
				info.Status = utils.Ptr("online")
			} else {
				info.Status = utils.Ptr("offline")
			}
			break
		}
	}
	return ctx.JSON(http.StatusOK, info)
}
