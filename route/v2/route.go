package v2

import (
	"github.com/dappsteros-io/DappsterOS/codegen"
	"github.com/dappsteros-io/DappsterOS/service"
)

type DappsterOS struct {
	fileUploadService *service.FileUploadService
}

func NewDappsterOS() codegen.ServerInterface {
	return &DappsterOS{
		fileUploadService: service.NewFileUploadService(),
	}
}
