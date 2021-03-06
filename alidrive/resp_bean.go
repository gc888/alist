package alidrive

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type RespError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type UserInfo struct {
	RespError
	DomainId       string                 `json:"domain_id"`
	UserId         string                 `json:"user_id"`
	Avatar         string                 `json:"avatar"`
	CreatedAt      int                    `json:"created_at"`
	UpdatedAt      int                    `json:"updated_at"`
	Email          string                 `json:"email"`
	NickName       string                 `json:"nick_name"`
	Phone          string                 `json:"phone"`
	Role           string                 `json:"role"`
	Status         string                 `json:"status"`
	UserName       string                 `json:"user_name"`
	Description    string                 `json:"description"`
	DefaultDriveId string                 `json:"default_drive_id"`
	UserData       map[string]interface{} `json:"user_data"`
}

type Files struct {
	RespError
	Items      []File `json:"items"`
	NextMarker string `json:"next_marker"`
	Readme     string `json:"readme"`
	Paths       []Path `json:"paths"`
}

type Path struct {
	Name   string `json:"name"`
	FileId string `json:"file_id"`
}

type File struct {
	RespError
	DriveId       string     `json:"drive_id"`
	CreatedAt     *time.Time `json:"created_at"`
	DomainId      string     `json:"domain_id"`
	EncryptMode   string     `json:"encrypt_mode"`
	FileExtension string     `json:"file_extension"`
	FileId        string     `json:"file_id"`
	Hidden        bool       `json:"hidden"`
	Name          string     `json:"name"`
	ParentFileId  string     `json:"parent_file_id"`
	Starred       bool       `json:"starred"`
	Status        string     `json:"status"`
	Type          string     `json:"type"`
	UpdatedAt     *time.Time `json:"updated_at"`
	// 文件多出部分
	Category           string                 `json:"category"`
	ContentHash        string                 `json:"content_hash"`
	ContentHashName    string                 `json:"content_hash_name"`
	ContentType        string                 `json:"content_type"`
	Crc64Hash          string                 `json:"crc_64_hash"`
	DownloadUrl        string                 `json:"download_url"`
	PunishFlag         int                    `json:"punish_flag"`
	Size               int                    `json:"size"`
	Thumbnail          string                 `json:"thumbnail"`
	Url                string                 `json:"url"`
	ImageMediaMetadata map[string]interface{} `json:"image_media_metadata"`

	Paths       []Path `json:"paths"`
}

func (resp *RespError) IsAvailable() bool {
	return resp.Code == ""
}

type TokenLoginResp struct {
	RespError
	Goto string `json:"goto"`
}

type TokenResp struct {
	RespError
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`

	UserInfo

	DefaultSboxDriveId string        `json:"default_sbox_drive_id"`
	ExpireTime         *time.Time    `json:"expire_time"`
	State              string        `json:"state"`
	ExistLink          []interface{} `json:"exist_link"`
	NeedLink           bool          `json:"need_link"`
	PinSetup           bool          `json:"pin_setup"`
	IsFirstLogin       bool          `json:"is_first_login"`
	NeedRpVerify       bool          `json:"need_rp_verify"`
	DeviceId           string        `json:"device_id"`
}

func HasPassword(files *Files) string {
	fileList := files.Items
	for _, file := range fileList {
		if strings.HasPrefix(file.Name, ".password-") {
			return file.Name[10:]
		}
	}
	return ""
}

func HasReadme(files *Files) string {
	fileList := files.Items
	for _, file := range fileList {
		if file.Name == "Readme.md" {
			resp, err := http.Get(file.Url)
			if err != nil {
				log.Errorf("Get Readme出错:%s", err.Error())
				return ""
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("读取 Readme出错:%s", err.Error())
				return ""
			}
			return string(data)
		}
	}
	return ""
}
