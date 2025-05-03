package services

import "sync"

type TranslateService struct {
}

var (
	translateServiceOnce     sync.Once
	translateServiceInstance *TranslateService
)

func InitializeTranslateService() {
	translateServiceOnce.Do(func() {
		translateServiceInstance = &TranslateService{}
	})
}

func GetTranslateService() *TranslateService {
	return translateServiceInstance
}
