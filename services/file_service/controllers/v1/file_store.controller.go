package v1

import (
	"context"
	"cortex_api/common"
	"cortex_api/common/api_models"
	"cortex_api/common/api_utils"
	"cortex_api/database/db_gen"
	"cortex_api/services/file_service/config"
	"cortex_api/services/file_service/utils"
	"encoding/json"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type FileTag struct {
	ID  pgtype.UUID `db:"id" json:"id"`
	Tag string      `db:"tag" json:"tag"`
}

type FileMeta struct {
	Key          string `json:"key"`
	Value        string `json:"value"`
	AttachedFile *File  `json:"attachedFile"`
}

type File struct {
	ID                  pgtype.UUID `json:"id"`
	IsSystem            bool        `json:"isSystem"`
	Filename            string      `json:"filename"`
	ContentType         string      `json:"contentType"`
	FileSize            int64       `json:"fileSize"`
	Browsable           bool        `json:"browsable"`
	Secure              bool        `json:"secure"`
	CompletedProcessing bool        `json:"completedProcessing"`
	Meta                *[]FileMeta `json:"meta"`
	Tags                *[]FileTag  `json:"tags"`
}

type FileStoreModel struct {
	Queries *db_gen.Queries
	DbCtx   context.Context
}

type FileStoreInterface interface {
	UploadFile(ctx *fiber.Ctx) error
	ListFiles(ctx *fiber.Ctx) error
	GetFile(ctx *fiber.Ctx) error
	DeleteFile(ctx *fiber.Ctx) error
	DownloadFileById(ctx *fiber.Ctx) error
	DownloadSecureFileByFilename(ctx *fiber.Ctx) error
	DownloadFileByFilename(ctx *fiber.Ctx) error
}

func FileStoreController(
	queries *db_gen.Queries,
	dbCtx context.Context,
) *FileStoreModel {
	return &FileStoreModel{
		Queries: queries,
		DbCtx:   dbCtx,
	}
}

func GetMimeType(mime *mimetype.MIME) (string, string) {
	types := strings.Split(mime.String(), "/")
	return types[0], types[1]
}

func ProcessThumbnail(controller *FileStoreModel, fileId pgtype.UUID, filename string, storedFilename string) {
	err := controller.Queries.DeleteFileByFilename(controller.DbCtx, utils.GetImageFilenameWithThumbnailSuffix(filename))
	if err != nil {
		log.Printf("ProcessThumbnail error while deleting previous thumbnail, error %v\n", err)
	}

	thumbnailStoredFilename := uuid.New().String()

	err = utils.ResizeImage(utils.GetFileStoragePathname(storedFilename), utils.GetFileStoragePathname(thumbnailStoredFilename), 200)
	if err != nil {
		log.Printf("ProcessThumbnail error while resizing image, error %v\n", err)
	}

	thumbnailMime, err := mimetype.DetectFile(utils.GetFileStoragePathname(thumbnailStoredFilename))
	if err != nil {
		log.Printf("ProcessThumbnail error while detecting thumbnail mime type, error %v\n", err)
	}

	stat, err := os.Stat(utils.GetFileStoragePathname(thumbnailStoredFilename))
	if err != nil {
		log.Printf("ProcessThumbnail error while getting thumbnail file stats, error %v\n", err)
	}

	fileStore := db_gen.UploadFileParams{
		Filename:       utils.GetImageFilenameWithThumbnailSuffix(filename),
		StoredFilename: thumbnailStoredFilename,
		ContentType:    pgtype.Text{String: thumbnailMime.String(), Valid: true},
		FileSize:       stat.Size(),
		Browsable:      false,
	}

	uploadedFile, dbErr := controller.Queries.UploadFile(controller.DbCtx, fileStore)
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("ProcessThumbnail error while saving object, error %v\n", dbErr)
		}
	}
	_, err = controller.Queries.AddFileMeta(controller.DbCtx, db_gen.AddFileMetaParams{
		FileStoreID:         fileId,
		Key:                 "thumbnail_filename",
		Value:               pgtype.Text{String: utils.GetImageFilenameWithThumbnailSuffix(filename), Valid: true},
		AttachedFileStoreID: uploadedFile,
	})
	if err != nil {
		log.Printf("ProcessThumbnail error while adding image thumbnail filename, error %v\n", err)
	}

	_, err = controller.Queries.AddFileMeta(controller.DbCtx, db_gen.AddFileMetaParams{
		FileStoreID: fileId,
		Key:         "thumbnail_file_id",
		Value:       pgtype.Text{String: common.FormatIdAsString(uploadedFile), Valid: true},
	})
	if err != nil {
		log.Printf("ProcessThumbnail error while adding image thumbnail file ID, error %v\n", err)
	}

	_, err = controller.Queries.AddFileMeta(controller.DbCtx, db_gen.AddFileMetaParams{
		FileStoreID: uploadedFile,
		Key:         "thumbnail_type",
		Value: pgtype.Text{
			String: "image",
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("ProcessThumbnail error while adding file thumbnail identification, error %v\n", err)
	}

	_, err = controller.Queries.AddFileMeta(controller.DbCtx, db_gen.AddFileMetaParams{
		FileStoreID: uploadedFile,
		Key:         "mime_type",
		Value: pgtype.Text{
			String: thumbnailMime.String(),
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("ProcessThumbnail error while adding file mime_type, error %v\n", err)
	}

	_, err = controller.Queries.AddFileMeta(controller.DbCtx, db_gen.AddFileMetaParams{
		FileStoreID: uploadedFile,
		Key:         "extension",
		Value: pgtype.Text{
			String: thumbnailMime.Extension(),
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("ProcessThumbnail error while adding file extension, error %v\n", err)
	}

	imageData, err := utils.GetImageData(utils.GetFileStoragePathname(thumbnailStoredFilename))
	if err != nil {
		log.Printf("ProcessThumbnail error while getting image data, error %v\n", err)
	} else {
		_, err = controller.Queries.AddFileMeta(controller.DbCtx, db_gen.AddFileMetaParams{
			FileStoreID: uploadedFile,
			Key:         "width",
			Value:       pgtype.Text{String: strconv.Itoa(imageData.Width), Valid: true},
		})
		if err != nil {
			log.Printf("ProcessThumbnail error while adding image width, error %v\n", err)
		}

		_, err = controller.Queries.AddFileMeta(controller.DbCtx, db_gen.AddFileMetaParams{
			FileStoreID: uploadedFile,
			Key:         "height",
			Value:       pgtype.Text{String: strconv.Itoa(imageData.Height), Valid: true},
		})
		if err != nil {
			log.Printf("ProcessThumbnail error while adding image height, error %v\n", err)
		}
	}

	err = controller.Queries.ProcessingComplete(controller.DbCtx, uploadedFile)
	if err != nil {
		log.Printf("ProcessThumbnail error while setting processing complete, error %v\n", err)
	}
}

func ProcessFile(controller *FileStoreModel, fileId pgtype.UUID, filename string, storedFilename string) {
	time.Sleep(time.Second)

	mimetype.SetLimit(1024 * 1024) // Set limit to 1MB.
	mime, err := mimetype.DetectFile(utils.GetFileStoragePathname(storedFilename))
	if err != nil {
		log.Printf("ProcessFile error while detecting file mime type, error %v\n", err)
	}
	log.Printf("ProcessFile detected file mime type: %s\n", mime.String())

	attachedFiles, err := controller.Queries.ListAttachmentsByFileStoreId(controller.DbCtx, fileId)
	if err != nil {
		log.Printf("ProcessFile error while fetching attached files, error %v\n", err)
	} else {
		for _, attachedFile := range attachedFiles {
			err = controller.Queries.DeleteFile(controller.DbCtx, attachedFile)
			if err != nil {
				log.Printf("ProcessFile error while deleting attached file, error %v\n", err)
			}
		}
	}

	err = controller.Queries.DeleteFileMetaByFileStoreId(controller.DbCtx, fileId)
	if err != nil {
		log.Printf("ProcessFile error while deleting file meta, error %v\n", err)
	}

	_, err = controller.Queries.AddFileMeta(controller.DbCtx, db_gen.AddFileMetaParams{
		FileStoreID: fileId,
		Key:         "mime_type",
		Value: pgtype.Text{
			String: mime.String(),
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("ProcessFile error while adding file mime_type, error %v\n", err)
	}

	_, err = controller.Queries.AddFileMeta(controller.DbCtx, db_gen.AddFileMetaParams{
		FileStoreID: fileId,
		Key:         "extension",
		Value: pgtype.Text{
			String: mime.Extension(),
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("ProcessFile error while adding file extension, error %v\n", err)
	}

	log.Printf("ProcessFile added mime type: %s %s\n", mime.String(), mime)

	mainType, _ := GetMimeType(mime)
	if mainType == "image" {
		imageData, err := utils.GetImageData(utils.GetFileStoragePathname(storedFilename))
		if err != nil {
			log.Printf("ProcessFile error while getting image data, error %v\n", err)
		} else {
			_, err = controller.Queries.AddFileMeta(controller.DbCtx, db_gen.AddFileMetaParams{
				FileStoreID: fileId,
				Key:         "width",
				Value:       pgtype.Text{String: strconv.Itoa(imageData.Width), Valid: true},
			})
			if err != nil {
				log.Printf("ProcessFile error while adding image width, error %v\n", err)
			}

			_, err = controller.Queries.AddFileMeta(controller.DbCtx, db_gen.AddFileMetaParams{
				FileStoreID: fileId,
				Key:         "height",
				Value:       pgtype.Text{String: strconv.Itoa(imageData.Height), Valid: true},
			})
			if err != nil {
				log.Printf("ProcessFile error while adding image height, error %v\n", err)
			}
		}

		ProcessThumbnail(controller, fileId, filename, storedFilename)
	}

	err = controller.Queries.ProcessingComplete(controller.DbCtx, fileId)
	if err != nil {
		log.Printf("ProcessFile error while setting processing complete, error %v\n", err)
	}
}

func SaveNewImage(controller *FileStoreModel, ctx *fiber.Ctx, file *multipart.FileHeader, secure bool) (*pgtype.UUID, string, error) {
	fileStore := db_gen.UploadFileParams{
		Filename:       file.Filename,
		StoredFilename: uuid.New().String(),
		ContentType:    pgtype.Text{String: file.Header.Get("Content-Type"), Valid: true},
		FileSize:       file.Size,
		Browsable:      true,
		Secure:         secure,
	}

	destination := fmt.Sprintf(utils.GetFileStoragePathname(fileStore.StoredFilename))
	if err := ctx.SaveFile(file, destination); err != nil {
		log.Printf("UploadFile error while saving file to the underlying file system, error %v\n", err)
		return nil, "", err
	}

	uploadedFile, dbErr := controller.Queries.UploadFile(controller.DbCtx, fileStore)
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("UploadFile error while uploading object, error %v\n", dbErr)
			return nil, "", dbErr
		}
	}

	return &uploadedFile, fileStore.StoredFilename, nil
}

func SaveUpdatedImage(controller *FileStoreModel, ctx *fiber.Ctx, file *multipart.FileHeader, currentFile db_gen.FileStore, secure bool) (*pgtype.UUID, string, error) {
	fileStore := db_gen.UpdateFileParams{
		ID:          currentFile.ID,
		Filename:    file.Filename,
		ContentType: pgtype.Text{String: file.Header.Get("Content-Type"), Valid: true},
		FileSize:    file.Size,
		Browsable:   secure,
	}

	destination := fmt.Sprintf(utils.GetFileStoragePathname(currentFile.StoredFilename))
	if err := ctx.SaveFile(file, destination); err != nil {
		log.Printf("UploadFile error while saving file to the underlying file system, error %v\n", err)
		return nil, "", err
	}

	_, dbErr := controller.Queries.UpdateFile(controller.DbCtx, fileStore)
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("UploadFile error while updating object, error %v\n", dbErr)
			return nil, "", dbErr
		}
	}

	return &currentFile.ID, currentFile.StoredFilename, nil
}

func SaveImage(controller *FileStoreModel, ctx *fiber.Ctx, file *multipart.FileHeader, currentFile db_gen.FileStore) (*pgtype.UUID, string, bool, error) {
	secure := ctx.FormValue("secure") == "true"
	// todo verify secure is working

	currentFile, dbErr := controller.Queries.GetFileByFilename(controller.DbCtx, file.Filename)
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("UploadFile error while fetching object, error %v\n", dbErr)
			return nil, "", false, dbErr
		}
	}
	if dbErr == nil {
		// file exists, perform update if overwrite is true
		overwrite := ctx.FormValue("overwrite")
		if overwrite != "true" {
			return nil, "", true, nil
		}

		uploadedFile, storedFilename, err := SaveUpdatedImage(controller, ctx, file, currentFile, secure)
		return uploadedFile, storedFilename, false, err
	} else {
		uploadedFile, storedFilename, err := SaveNewImage(controller, ctx, file, secure)
		return uploadedFile, storedFilename, false, err
	}
}

func (controller *FileStoreModel) UploadFile(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	file, err := ctx.FormFile("file")
	if err != nil {
		log.Printf("UploadFile error while saving file from form data, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid file upload",
		})
		return ctx.Send(payload)
	}

	acceptableType := false
	for _, acceptedType := range config.AppConfig().GetAllConfig().AcceptedFileTypes {
		if acceptedType == file.Header.Get("Content-Type") {
			acceptableType = true
			break
		}
	}
	if !acceptableType {
		log.Printf("UploadFile error while validating file type\n")
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid file type, please confirm the file type is supported",
		})
		return ctx.Send(payload)
	}

	var uploadedFile = &pgtype.UUID{}
	var storedFilename = ""
	var overwriteError = false

	uploadedFile, storedFilename, overwriteError, err = SaveImage(controller, ctx, file, db_gen.FileStore{})
	if overwriteError {
		ctx.Status(http.StatusConflict)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Unable to upload file, file already exists",
		})
		return ctx.Send(payload)
	}

	go func() {
		ProcessFile(controller, *uploadedFile, file.Filename, storedFilename)
	}()

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   uploadedFile,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: marshalErr.Error(),
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *FileStoreModel) ListFiles(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	requestedPage, requestedPageSize, sortBy, sortOrder := api_utils.GetPaginationParams(ctx)
	rowCount, dbErr := controller.Queries.CountFiles(controller.DbCtx)
	if dbErr != nil {
		log.Printf("ListFiles error while retrieving count, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve file count",
		})
		return ctx.Send(payload)
	}

	query, err := api_utils.ParseQueryString(ctx)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid query",
		})
		return ctx.Send(payload)
	}

	files := make([]db_gen.FileStore, 0)

	browsable := query.Get("browsable")
	if browsable == "" {
		if sortOrder == "desc" {
			files, dbErr = controller.Queries.ListFilesDesc(controller.DbCtx, db_gen.ListFilesDescParams{
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		} else {
			files, dbErr = controller.Queries.ListFilesAsc(controller.DbCtx, db_gen.ListFilesAscParams{
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		}
	} else {
		if sortOrder == "desc" {
			files, dbErr = controller.Queries.ListBrowsableFilesDesc(controller.DbCtx, db_gen.ListBrowsableFilesDescParams{
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
				Browsable:         browsable == "true",
			})
		} else {
			files, dbErr = controller.Queries.ListBrowsableFilesAsc(controller.DbCtx, db_gen.ListBrowsableFilesAscParams{
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
				Browsable:         browsable == "true",
			})
		}
	}
	if dbErr != nil {
		log.Printf("ListFiles error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve files",
		})
		return ctx.Send(payload)
	}

	fileViews := make([]File, 0)
	for _, file := range files {
		metaData, err := controller.Queries.ListFileMetaByFileStoreId(controller.DbCtx, file.ID)
		if err != nil {
			log.Printf("ListFiles error while fetching file meta information, error %v\n", err)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to retrieve file meta information",
			})
			return ctx.Send(payload)
		}

		metaViews := make([]FileMeta, 0)
		for _, meta := range metaData {
			if meta.AttachedFileStoreID.Valid {
				attachedFile, err := controller.Queries.GetFileById(controller.DbCtx, meta.AttachedFileStoreID)
				if err != nil {
					log.Printf("ListFiles error while fetching attached file, error %v\n", err)
					payload, _ := json.Marshal(api_models.Failed{
						Status:       api_utils.FailedResponse,
						ErrorCode:    http.StatusBadGateway,
						ErrorMessage: "Unable to retrieve attached file",
					})
					return ctx.Send(payload)
				}
				metaViews = append(metaViews, FileMeta{
					Key:   meta.Key,
					Value: meta.Value.String,
					AttachedFile: &File{
						ID:                  attachedFile.ID,
						IsSystem:            attachedFile.IsSystem,
						Filename:            attachedFile.Filename,
						ContentType:         attachedFile.ContentType.String,
						Browsable:           attachedFile.Browsable,
						FileSize:            attachedFile.FileSize,
						CompletedProcessing: attachedFile.CompletedProcessing,
					},
				})
				continue
			} else {
				metaViews = append(metaViews, FileMeta{
					Key:   meta.Key,
					Value: meta.Value.String,
				})
			}
		}

		tagData, err := controller.Queries.ListFileTagByFileStoreId(controller.DbCtx, file.ID)
		if err != nil {
			log.Printf("ListFiles error while fetching file meta information, error %v\n", err)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to retrieve file meta information",
			})
			return ctx.Send(payload)
		}

		tags := make([]FileTag, 0)
		for _, tag := range tagData {
			tags = append(tags, FileTag{
				ID:  tag.ID,
				Tag: tag.Tag,
			})
		}

		fileViews = append(fileViews, File{
			ID:                  file.ID,
			IsSystem:            file.IsSystem,
			Filename:            file.Filename,
			ContentType:         file.ContentType.String,
			FileSize:            file.FileSize,
			Browsable:           file.Browsable,
			Meta:                &metaViews,
			Tags:                &tags,
			CompletedProcessing: file.CompletedProcessing,
		})
	}

	if fileViews == nil {
		fileViews = make([]File, 0)
	}
	results := api_utils.GetPaginatedResults(rowCount, requestedPage, requestedPageSize, fileViews)
	payload, marshalErr := json.Marshal(results)
	if marshalErr != nil {
		log.Printf("ListFiles error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *FileStoreModel) GetFile(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	fileUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("GetFile error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	file, dbErr := controller.Queries.GetFileById(controller.DbCtx, pgtype.UUID{Bytes: fileUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "File not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("GetFile error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve file",
		})
		return ctx.Send(payload)
	}

	metaData, err := controller.Queries.ListFileMetaByFileStoreId(controller.DbCtx, file.ID)
	if err != nil {
		log.Printf("GetFile error while fetching file meta information, error %v\n", err)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve file meta information",
		})
		return ctx.Send(payload)
	}

	metaViews := make([]FileMeta, 0)
	for _, meta := range metaData {
		if meta.AttachedFileStoreID.Valid {
			attachedFile, err := controller.Queries.GetFileById(controller.DbCtx, meta.AttachedFileStoreID)
			if err != nil {
				log.Printf("ListFiles error while fetching attached file, error %v\n", err)
				payload, _ := json.Marshal(api_models.Failed{
					Status:       api_utils.FailedResponse,
					ErrorCode:    http.StatusBadGateway,
					ErrorMessage: "Unable to retrieve attached file",
				})
				return ctx.Send(payload)
			}
			metaViews = append(metaViews, FileMeta{
				Key:   meta.Key,
				Value: meta.Value.String,
				AttachedFile: &File{
					ID:                  attachedFile.ID,
					IsSystem:            attachedFile.IsSystem,
					Filename:            attachedFile.Filename,
					ContentType:         attachedFile.ContentType.String,
					Browsable:           attachedFile.Browsable,
					FileSize:            attachedFile.FileSize,
					CompletedProcessing: attachedFile.CompletedProcessing,
				},
			})
			continue
		} else {
			metaViews = append(metaViews, FileMeta{
				Key:   meta.Key,
				Value: meta.Value.String,
			})
		}
	}

	tagData, err := controller.Queries.ListFileTagByFileStoreId(controller.DbCtx, file.ID)
	if err != nil {
		log.Printf("GetFile error while fetching file meta information, error %v\n", err)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve file meta information",
		})
		return ctx.Send(payload)
	}

	tags := make([]FileTag, 0)
	for _, tag := range tagData {
		tags = append(tags, FileTag{
			ID:  tag.ID,
			Tag: tag.Tag,
		})
	}

	fileView := File{
		ID:                  file.ID,
		IsSystem:            file.IsSystem,
		Filename:            file.Filename,
		ContentType:         file.ContentType.String,
		Browsable:           file.Browsable,
		FileSize:            file.FileSize,
		CompletedProcessing: file.CompletedProcessing,
		Secure:              file.Secure,
		Meta:                &metaViews,
		Tags:                &tags,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   fileView,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetFile error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *FileStoreModel) GetFileByFilename(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	filename := ctx.Params("filename")
	file, dbErr := controller.Queries.GetFileByFilename(controller.DbCtx, filename)
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "File not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("GetFile error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve file",
		})
		return ctx.Send(payload)
	}

	metaData, err := controller.Queries.ListFileMetaByFileStoreId(controller.DbCtx, file.ID)
	if err != nil {
		log.Printf("GetFile error while fetching file meta information, error %v\n", err)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve file meta information",
		})
		return ctx.Send(payload)
	}

	metaViews := make([]FileMeta, 0)
	for _, meta := range metaData {
		if meta.AttachedFileStoreID.Valid {
			attachedFile, err := controller.Queries.GetFileById(controller.DbCtx, meta.AttachedFileStoreID)
			if err != nil {
				log.Printf("ListFiles error while fetching attached file, error %v\n", err)
				payload, _ := json.Marshal(api_models.Failed{
					Status:       api_utils.FailedResponse,
					ErrorCode:    http.StatusBadGateway,
					ErrorMessage: "Unable to retrieve attached file",
				})
				return ctx.Send(payload)
			}
			metaViews = append(metaViews, FileMeta{
				Key:   meta.Key,
				Value: meta.Value.String,
				AttachedFile: &File{
					ID:                  attachedFile.ID,
					IsSystem:            attachedFile.IsSystem,
					Filename:            attachedFile.Filename,
					ContentType:         attachedFile.ContentType.String,
					Browsable:           attachedFile.Browsable,
					FileSize:            attachedFile.FileSize,
					CompletedProcessing: attachedFile.CompletedProcessing,
				},
			})
			continue
		} else {
			metaViews = append(metaViews, FileMeta{
				Key:   meta.Key,
				Value: meta.Value.String,
			})
		}
	}

	tagData, err := controller.Queries.ListFileTagByFileStoreId(controller.DbCtx, file.ID)
	if err != nil {
		log.Printf("GetFile error while fetching file meta information, error %v\n", err)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve file meta information",
		})
		return ctx.Send(payload)
	}

	tags := make([]FileTag, 0)
	for _, tag := range tagData {
		tags = append(tags, FileTag{
			ID:  tag.ID,
			Tag: tag.Tag,
		})
	}

	fileView := File{
		ID:                  file.ID,
		IsSystem:            file.IsSystem,
		Filename:            file.Filename,
		ContentType:         file.ContentType.String,
		Browsable:           file.Browsable,
		FileSize:            file.FileSize,
		CompletedProcessing: file.CompletedProcessing,
		Secure:              file.Secure,
		Meta:                &metaViews,
		Tags:                &tags,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   fileView,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetFile error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *FileStoreModel) DownloadSecureFileByFilename(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	filename := ctx.Params("filename")
	if filename == "" {
		log.Printf("DownloadSecureFileByFilename missing filename\n")
		ctx.Status(http.StatusBadRequest)
		return ctx.Send([]byte("Invalid filename"))
	}

	file, dbErr := controller.Queries.GetFileByFilename(controller.DbCtx, filename)
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			return ctx.Send([]byte("File not found"))
		}
		log.Printf("DownloadSecureFileByFilename error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		return ctx.Send([]byte("Unable to retrieve file"))
	}

	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	ctx.Set("Content-Type", file.ContentType.String)
	ctx.Status(http.StatusOK)
	return ctx.SendFile(fmt.Sprintf(utils.GetFileStoragePathname(file.StoredFilename)))
}

func (controller *FileStoreModel) DownloadFileByFilename(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	filename := ctx.Params("filename")
	if filename == "" {
		log.Printf("DownloadFileByFilename missing filename\n")
		ctx.Status(http.StatusBadRequest)
		return ctx.Send([]byte("Invalid filename"))
	}

	file, dbErr := controller.Queries.GetPublicFileByFilename(controller.DbCtx, filename)
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			return ctx.Send([]byte("File not found"))
		}
		log.Printf("DownloadFileByFilename error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		return ctx.Send([]byte("Unable to retrieve file"))
	}

	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	ctx.Set("Content-Type", file.ContentType.String)
	ctx.Status(http.StatusOK)
	return ctx.SendFile(fmt.Sprintf(utils.GetFileStoragePathname(file.StoredFilename)))
}

func (controller *FileStoreModel) DownloadFileById(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	fileUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("DownloadFileById error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	file, dbErr := controller.Queries.GetFileById(controller.DbCtx, pgtype.UUID{Bytes: fileUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "File not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("DownloadFileById error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve file",
		})
		return ctx.Send(payload)
	}

	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	ctx.Set("Content-Type", file.ContentType.String)
	ctx.Status(http.StatusOK)
	return ctx.SendFile(fmt.Sprintf(utils.GetFileStoragePathname(file.StoredFilename)))
}

func (controller *FileStoreModel) DeleteFile(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	fileUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("DeleteFile error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	_, dbErr := controller.Queries.GetFileById(controller.DbCtx, pgtype.UUID{Bytes: fileUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "File not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("DeleteFile error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete file",
		})
		return ctx.Send(payload)
	}

	// we only need to delete the file from the database, the file deletion is handled by the housekeeping process
	dbErr = controller.Queries.DeleteFile(controller.DbCtx, pgtype.UUID{Bytes: fileUuid, Valid: true})
	if dbErr != nil {
		log.Printf("DeleteFile error while deleting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete file",
		})
		return ctx.Send(payload)
	}

	attachments, dbErr := controller.Queries.ListAttachmentsByFileStoreId(controller.DbCtx, pgtype.UUID{Bytes: fileUuid, Valid: true})
	if dbErr != nil {
		log.Printf("DeleteFile error while fetching attachments, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete file(s)",
		})
		return ctx.Send(payload)
	}

	err = controller.Queries.DeleteFileMetaByFileStoreId(controller.DbCtx, pgtype.UUID{Bytes: fileUuid, Valid: true})
	if err != nil {
		log.Printf("DeleteFile error while deleting file meta, error %v\n", err)
	}

	for _, attachment := range attachments {
		dbErr = controller.Queries.DeleteFile(controller.DbCtx, attachment)
		if dbErr != nil {
			log.Printf("DeleteFile error while deleting attachment, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to delete file(s)",
			})
			return ctx.Send(payload)
		}
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   nil,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("DeleteFile error serialising response, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}
