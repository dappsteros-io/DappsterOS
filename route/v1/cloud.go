package v1

import (
	"strings"

	"github.com/dappsteros-io/DappsterOS-Common/utils/logger"
	"github.com/dappsteros-io/DappsterOS/drivers/dropbox"
	"github.com/dappsteros-io/DappsterOS/drivers/google_drive"
	"github.com/dappsteros-io/DappsterOS/drivers/onedrive"
	"github.com/dappsteros-io/DappsterOS/model"
	"github.com/dappsteros-io/DappsterOS/pkg/utils/common_err"
	"github.com/dappsteros-io/DappsterOS/pkg/utils/httper"
	"github.com/dappsteros-io/DappsterOS/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ListStorages(ctx echo.Context) error {
	// var req model.PageReq
	// if err := ctx.Bind(&req); err != nil {
	// 	return ctx.JSON(common_err.SUCCESS, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: err.Error()})
	// 	return
	// }
	// req.Validate()

	// logger.Info("ListStorages", zap.Any("req", req))
	// storages, total, err := service.MyService.Storage().GetStorages(req.Page, req.PerPage)
	// if err != nil {
	// 	return ctx.JSON(common_err.SUCCESS, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
	// 	return
	// }
	// return ctx.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: model.PageResp{
	// 	Content: storages,
	// 	Total:   total,
	// }})
	r, err := service.MyService.Storage().GetStorages()
	if err != nil {
		return ctx.JSON(common_err.SUCCESS, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
	}

	for i := 0; i < len(r.MountPoints); i++ {
		dataMap, err := service.MyService.Storage().GetConfigByName(r.MountPoints[i].Fs)
		if err != nil {
			logger.Error("GetConfigByName", zap.Any("err", err))
			continue
		}
		if dataMap["type"] == "drive" {
			r.MountPoints[i].Icon = google_drive.ICONURL
		}
		if dataMap["type"] == "dropbox" {
			r.MountPoints[i].Icon = dropbox.ICONURL
		}
		if dataMap["type"] == "onedrive" {
			r.MountPoints[i].Icon = onedrive.ICONURL
		}
		r.MountPoints[i].Name = dataMap["username"]
	}
	list := []httper.MountPoint{}

	for _, v := range r.MountPoints {
		list = append(list, httper.MountPoint{
			Fs:         v.Fs,
			Icon:       v.Icon,
			MountPoint: v.MountPoint,
			Name:       v.Name,
		})
	}

	return ctx.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: list})
}

func UmountStorage(ctx echo.Context) error {
	json := make(map[string]string)
	ctx.Bind(&json)
	mountPoint := json["mount_point"]
	if mountPoint == "" {
		return ctx.JSON(common_err.CLIENT_ERROR, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: "mount_point is empty"})
	}
	err := service.MyService.Storage().UnmountStorage(mountPoint)
	if err != nil {
		return ctx.JSON(common_err.SERVICE_ERROR, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
	}
	service.MyService.Storage().DeleteConfigByName(strings.ReplaceAll(mountPoint, "/mnt/", ""))
	return ctx.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: "success"})
}

func GetStorage(ctx echo.Context) error {
	// idStr := ctx.QueryParam("id")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	return ctx.JSON(common_err.SUCCESS, model.Result{Success: common_err.CLIENT_ERROR, Message: common_err.GetMsg(common_err.CLIENT_ERROR), Data: err.Error()})
	// 	return
	// }
	// storage, err := service.MyService.Storage().GetStorageById(uint(id))
	// if err != nil {
	// 	return ctx.JSON(common_err.SUCCESS, model.Result{Success: common_err.SERVICE_ERROR, Message: common_err.GetMsg(common_err.SERVICE_ERROR), Data: err.Error()})
	// 	return
	// }
	// return ctx.JSON(common_err.SUCCESS, model.Result{Success: common_err.SUCCESS, Message: common_err.GetMsg(common_err.SUCCESS), Data: storage})
	return nil
}
