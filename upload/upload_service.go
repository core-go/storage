package upload

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

const contentTypeHeader = "Content-Type"

type UploadService interface {
	// UploadFile(fileName string, contentType string, data []byte, r *http.Request) (rs string, errorRespone error)
	// DeleteFile(url string, r *http.Request) (bool, error)
	UploadGallery(data Upload, r *http.Request) ([]UploadInfo, error)
	DeleteGalleryFile(id string, url string, r *http.Request) (int64, error)
	UploadCover(id string, data []UploadData, contentType string, r *http.Request) (string, error)
	UploadImage(id string, data []UploadData, contentType string, r *http.Request) (string, error)
	UpdateGallery(data map[string][]UploadInfo, id string, r *http.Request) (int64, error)
	GetGallery(id string, r *http.Request) ([]UploadInfo, error)
	// addExternalResource(id string, data UploadInfo, r *http.Request) (bool, error)
	// deleteExternalResource(id string, url string, r *http.Request) (bool, error)
}

type StorageService interface {
	Upload(ctx context.Context, directory string, filename string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type uploadService struct {
	repository       StorageRepository
	Service          StorageService
	Provider         string
	GeneralDirectory string
	Directory        string
	KeyFile          string
	SizesImage       []int
	SizesCover       []int
	Columns          UploadFieldColumn
}
type UploadInfo struct {
	Source string `json:"source,omitempty" gorm:"column:source" bson:"source,omitempty" dynamodbav:"source,omitempty" firestore:"source,omitempty"`
	Url    string `json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
	Type   string `json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
}

type UploadModel struct {
	UserId   string       `json:"userId,omitempty" gorm:"column:userid;primary_key" bson:"_id,omitempty" dynamodbav:"userId,omitempty" firestore:"userId,omitempty" validate:"required,max=40"`
	ImageURL *string      `json:"imageURL,omitempty" gorm:"column:imageURL" bson:"imageURL,omitempty" dynamodbav:"imageURL,omitempty" firestore:"imageURL,omitempty"`
	CoverURL *string      `json:"coverURL,omitempty" gorm:"column:coverUrl" bson:"coverURL,omitempty" dynamodbav:"coverURL,omitempty" firestore:"coverURL,omitempty" `
	Gallery  []UploadInfo `json:"gallery,omitempty" gorm:"column:gallery" bson:"gallery,omitempty" dynamodbav:"gallery,omitempty" firestore:"gallery,omitempty"`
}
type UploadData struct {
	Name string `json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Data []byte `json:"data,omitempty" gorm:"column:data" bson:"data,omitempty" dynamodbav:"data,omitempty" firestore:"data,omitempty"`
}

func (u UploadInfo) Value() (driver.Value, error) {
	return json.Marshal(u)
}
func (u *UploadInfo) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &u)
}

func NewUploadService(DB *sql.DB, table string,
	toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}, service StorageService, provider string, generalDirectory string,
	keyFile string, directory string,
	sizesCover []int,
	sizesImage []int,
	IdCol string,
	CoverCol string,
	ImageCol string,
	Gallery string) UploadService {

	columns := UploadFieldColumn{
		Image:   &ImageCol,
		Cover:   &CoverCol,
		Gallery: &Gallery,
		Id:      "id",
	}
	repository := CreateStorageRepository(DB, table, columns, toArray)
	var sizesI = []int{40, 400}
	var sizesC = []int{576, 768}

	if len(sizesCover) != 0 {
		sizesC = sizesCover
	}
	if len(sizesImage) != 0 {
		sizesI = sizesCover
	}
	return &uploadService{Service: service, Provider: provider, GeneralDirectory: generalDirectory,
		KeyFile: keyFile, Directory: directory,
		SizesImage: sizesI, SizesCover: sizesC, repository: repository, Columns: columns}
}
func (u *uploadService) UploadCover(id string, data []UploadData, contentType string, r *http.Request) (string, error) {
	//delete
	result, err := u.repository.Load(r.Context(), id, *u.Columns.Cover)
	if err != nil {
		return "", err
	}

	if result.CoverURL != nil {
		_, err := u.DeleteFileUpload(u.SizesCover, *result.CoverURL, r)
		if err != nil {
			return "", errors.New("internal server error")
		}
	}
	//upload
	var newUrl string
	for i := range data {
		file := data[i]
		rs, errorRespone := u.UploadFile(file.Name, contentType, file.Data, r)
		if errorRespone != nil {
			return "", errorRespone
		}
		if i != 0 {
			continue
		}
		newUrl = rs
	}
	user := UploadModel{UserId: id, CoverURL: &newUrl}

	_, err1 := u.Update(r.Context(), user)
	if err1 != nil {
		return "", err1
	}
	return newUrl, nil
}

func (u *uploadService) UploadImage(id string, data []UploadData, contentType string, r *http.Request) (string, error) {
	//delete
	result, _ := u.repository.Load(r.Context(), id, *u.Columns.Image)

	if result.ImageURL != nil {
		_, err := u.DeleteFileUpload(u.SizesImage, *result.ImageURL, r)
		if err != nil {
			return "", errors.New("internal server error")
		}
	}
	//upload
	var newUrl string
	for i := range data {
		file := data[i]
		rs, errorRespone := u.UploadFile(file.Name, contentType, file.Data, r)
		if errorRespone != nil {
			return "", errorRespone
		}
		if i != 0 {
			continue
		}
		newUrl = rs
	}
	user := UploadModel{UserId: id, ImageURL: &newUrl}

	_, err1 := u.Update(r.Context(), user)
	if err1 != nil {
		return "", err1
	}
	return newUrl, nil
}

func (u *uploadService) UploadGallery(data Upload, r *http.Request) ([]UploadInfo, error) {
	result, _ := u.repository.Load(r.Context(), data.UserId, *u.Columns.Gallery)
	gallery := []UploadInfo{}
	if result.Gallery != nil {
		gallery = result.Gallery
	}

	rs, errorRespone := u.UploadFile(data.Name, data.Type, data.Data, r)
	if errorRespone != nil {
		return nil, errorRespone
	}

	gallery = append(gallery, UploadInfo{Source: data.Source, Type: strings.Split(data.Type, "/")[0], Url: rs})
	user := UploadModel{UserId: data.UserId, Gallery: gallery}

	_, err := u.Update(r.Context(), user)
	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (u *uploadService) UploadFile(fileName string, contentType string, data []byte, r *http.Request) (rs string, errorRespone error) {
	directory := u.Directory
	rs, err2 := u.Service.Upload(r.Context(), directory, fileName, data, contentType)
	if err2 != nil {
		errorRespone = err2
		return
	}
	return
}

func (u *uploadService) DeleteGalleryFile(id string, url string, r *http.Request) (int64, error) {
	rs, err := u.repository.Load(r.Context(), id, *u.Columns.Gallery)
	if err != nil {
		return 0, err
	}
	gallery := rs.Gallery
	// if find url in gallery
	idx := -1
	for i := range gallery {
		if gallery[i].Url == url {
			idx = i
		}
	}
	if idx == -1 {
		return 0, nil
	}
	_, err2 := u.DeleteFile(url, r)
	if err2 != nil {
		return 0, err2
	}
	gallery = append(gallery[:idx], gallery[idx+1:]...)
	user := UploadModel{UserId: id, Gallery: gallery}
	_, err3 := u.Update(r.Context(), user)
	if err3 != nil {
		return 0, err3
	}
	return 1, nil
}

func (u *uploadService) GetGallery(id string, r *http.Request) ([]UploadInfo, error) {
	rs, err := u.repository.Load(r.Context(), id, *u.Columns.Gallery)
	return rs.Gallery, err
}

func (u *uploadService) UpdateGallery(data map[string][]UploadInfo, id string, r *http.Request) (int64, error) {

	user := UploadModel{UserId: id, Gallery: data["data"]}
	_, err2 := u.Update(r.Context(), user)
	if err2 != nil {
		return 0, err2
	}
	return 1, err2
}
func (u *uploadService) DeleteFileUpload(sizes []int, url string, r *http.Request) (bool, error) {
	rs, err := u.DeleteFile(url, r)
	fmt.Print(rs, err)
	// if err != nil {
	// 	return false, errors.New("internal server error")
	// }
	for i := range sizes {
		oldUrl := removeExt(url) + "_" + strconv.Itoa(sizes[i]) + getExt(url)
		arr := strings.Split(oldUrl, "/")
		delUrl := arr[len(arr)-2] + "/" + arr[len(arr)-1]
		rss, err := u.DeleteFile(delUrl, r)
		fmt.Print(rss, err)
	}
	return true, nil
}

func (u *uploadService) DeleteFile(url string, r *http.Request) (bool, error) {
	arrOrigin := strings.Split(url, "/")
	delOriginUrl := arrOrigin[len(arrOrigin)-2] + "/" + arrOrigin[len(arrOrigin)-1]
	rs, err := u.Service.Delete(r.Context(), delOriginUrl)
	return rs, err
}

func (r *uploadService) Update(ctx context.Context, user UploadModel) (bool, error) {
	_, err := r.repository.Update(ctx, user)
	if err == nil {
		return false, nil
	}
	return false, err
}

func getExt(file string) string {
	ext := filepath.Ext(file)
	if strings.HasPrefix(ext, ":") {
		ext = ext[1:]
		return ext
	}
	return ext
}

func removeExt(file string) string {
	return file[:len(file)-len(filepath.Ext(file))]
}
