package services

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/smoothweb/backend/internal/models"
)

type MediaService struct {
	db             *gorm.DB
	maxFileSize    int64    // Maximum file size in bytes
	allowedTypes   []string // Allowed MIME types
	uploadPath     string   // Base upload directory
	thumbnailSizes []ThumbnailSize
}

type ThumbnailSize struct {
	Name   string
	Width  int
	Height int
}

type MediaUploadRequest struct {
	UserID uint
	File   multipart.File
	Header *multipart.FileHeader
}

type MediaUploadResponse struct {
	MediaID       uint   `json:"media_id"`
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	FilePath      string `json:"file_path"`
	ThumbnailPath string `json:"thumbnail_path,omitempty"`
	URL           string `json:"url"`
	ThumbnailURL  string `json:"thumbnail_url,omitempty"`
	MimeType      string `json:"mime_type"`
	FileSize      int64  `json:"file_size"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	FileType      string `json:"file_type"`
}

type CropRequest struct {
	MediaID uint
	UserID  uint
	X       int
	Y       int
	Width   int
	Height  int
}

// NewMediaService creates a new media service
func NewMediaService(db *gorm.DB) *MediaService {
	return &MediaService{
		db:          db,
		maxFileSize: 5 * 1024 * 1024, // 5MB default
		allowedTypes: []string{
			"image/jpeg",
			"image/jpg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		uploadPath: "./uploads",
		thumbnailSizes: []ThumbnailSize{
			{Name: "small", Width: 32, Height: 32},
			{Name: "medium", Width: 150, Height: 150},
			{Name: "large", Width: 300, Height: 300},
		},
	}
}

// UploadAvatar handles avatar upload with validation, processing, and storage
func (s *MediaService) UploadAvatar(req *MediaUploadRequest) (*MediaUploadResponse, error) {
	// Set avatar-specific limits
	savedMaxSize := s.maxFileSize
	s.maxFileSize = 5 * 1024 * 1024 // 5MB for avatars
	defer func() { s.maxFileSize = savedMaxSize }()

	return s.uploadImage(req, models.MediaTypeAvatar)
}

// UploadCoverPhoto handles cover photo upload
func (s *MediaService) UploadCoverPhoto(req *MediaUploadRequest) (*MediaUploadResponse, error) {
	// Set cover photo-specific limits
	savedMaxSize := s.maxFileSize
	s.maxFileSize = 10 * 1024 * 1024 // 10MB for cover photos
	defer func() { s.maxFileSize = savedMaxSize }()

	return s.uploadImage(req, models.MediaTypeCoverPhoto)
}

// uploadImage is the core image upload function
func (s *MediaService) uploadImage(req *MediaUploadRequest, fileType models.MediaType) (*MediaUploadResponse, error) {
	// 1. Validate file size
	if req.Header.Size > s.maxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum limit of %d bytes", s.maxFileSize)
	}

	// 2. Validate file type
	contentType := req.Header.Header.Get("Content-Type")
	if !s.isAllowedType(contentType) {
		return nil, fmt.Errorf("file type %s is not allowed", contentType)
	}

	// 3. Read file content
	fileBytes, err := io.ReadAll(req.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 4. Decode and validate image
	img, format, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, fmt.Errorf("invalid image format: %w", err)
	}

	// 5. Get image dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 6. Validate image dimensions
	if width < 100 || height < 100 {
		return nil, errors.New("image dimensions too small (minimum 100x100)")
	}
	if width > 4096 || height > 4096 {
		return nil, errors.New("image dimensions too large (maximum 4096x4096)")
	}

	// 7. Generate unique filename
	fileExtension := s.getFileExtension(contentType)
	fileName := fmt.Sprintf("%s_%s%s", fileType, uuid.New().String(), fileExtension)

	// 8. Create user-specific directory
	var userDir string
	if fileType == models.MediaTypeAvatar {
		userDir = filepath.Join(s.uploadPath, "avatars", fmt.Sprintf("%d", req.UserID))
	} else if fileType == models.MediaTypeCoverPhoto {
		userDir = filepath.Join(s.uploadPath, "covers", fmt.Sprintf("%d", req.UserID))
	} else {
		userDir = filepath.Join(s.uploadPath, "media", fmt.Sprintf("%d", req.UserID))
	}

	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 9. Process and save the image
	filePath := filepath.Join(userDir, fileName)
	if err := s.saveImage(img, format, filePath); err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	// 10. Generate thumbnails
	thumbnailPath := ""
	if len(s.thumbnailSizes) > 0 {
		thumbnailPath, err = s.generateThumbnails(img, format, userDir, fileName)
		if err != nil {
			// Cleanup uploaded file if thumbnail generation fails
			os.Remove(filePath)
			return nil, fmt.Errorf("failed to generate thumbnails: %w", err)
		}
	}

	// 11. Delete old avatar/cover photo if this is a replacement
	if fileType == models.MediaTypeAvatar || fileType == models.MediaTypeCoverPhoto {
		s.deleteOldMedia(req.UserID, fileType)
	}

	// 12. Save media record to database
	media := models.MediaFile{
		UserID:          req.UserID,
		FileName:        fileName,
		OriginalName:    req.Header.Filename,
		FilePath:        filePath,
		ThumbnailPath:   thumbnailPath,
		MimeType:        contentType,
		FileSize:        req.Header.Size,
		Width:           &width,
		Height:          &height,
		FileType:        fileType,
		IsPublic:        true,
		StorageProvider: "local",
	}

	if err := s.db.Create(&media).Error; err != nil {
		// Cleanup uploaded file if database save fails
		os.Remove(filePath)
		if thumbnailPath != "" {
			os.Remove(thumbnailPath)
		}
		return nil, fmt.Errorf("failed to save media record: %w", err)
	}

	// 13. Build response
	response := &MediaUploadResponse{
		MediaID:       media.ID,
		FileName:      media.FileName,
		OriginalName:  media.OriginalName,
		FilePath:      media.FilePath,
		ThumbnailPath: media.ThumbnailPath,
		URL:           s.getPublicURL(media.FilePath),
		ThumbnailURL:  s.getPublicURL(media.ThumbnailPath),
		MimeType:      media.MimeType,
		FileSize:      media.FileSize,
		Width:         width,
		Height:        height,
		FileType:      string(media.FileType),
	}

	// 14. Update profile with avatar/cover photo reference
	if fileType == models.MediaTypeAvatar {
		if err := s.updateProfileAvatar(req.UserID, media.ID); err != nil {
			// Non-blocking error: log but don't fail the upload
			fmt.Printf("Warning: failed to update profile avatar: %v\n", err)
		}
	} else if fileType == models.MediaTypeCoverPhoto {
		if err := s.updateProfileCoverPhoto(req.UserID, media.ID); err != nil {
			// Non-blocking error: log but don't fail the upload
			fmt.Printf("Warning: failed to update profile cover photo: %v\n", err)
		}
	}

	return response, nil
}

// CropMedia crops an image to the specified coordinates
func (s *MediaService) CropMedia(req *CropRequest) (*MediaUploadResponse, error) {
	var media models.MediaFile

	// Find the media file
	if err := s.db.Where("id = ? AND user_id = ?", req.MediaID, req.UserID).First(&media).Error; err != nil {
		return nil, err
	}

	// Open the image file
	file, err := os.Open(media.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the image
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Validate crop coordinates
	bounds := img.Bounds()
	if req.X < 0 || req.Y < 0 || req.X+req.Width > bounds.Dx() || req.Y+req.Height > bounds.Dy() {
		return nil, errors.New("invalid crop coordinates")
	}

	// Crop the image
	croppedImg := imaging.Crop(img, image.Rect(req.X, req.Y, req.X+req.Width, req.Y+req.Height))

	// Resize to standard size (400x400 for avatars, keep aspect for covers)
	var finalImg image.Image
	if media.FileType == models.MediaTypeAvatar {
		finalImg = imaging.Resize(croppedImg, 400, 400, imaging.Lanczos)
	} else {
		// For cover photos, maintain aspect ratio but limit width to 1200px
		if croppedImg.Bounds().Dx() > 1200 {
			finalImg = imaging.Resize(croppedImg, 1200, 0, imaging.Lanczos)
		} else {
			finalImg = croppedImg
		}
	}

	// Generate new filename
	fileExtension := s.getFileExtension(media.MimeType)
	newFileName := fmt.Sprintf("%s_cropped_%s%s", media.FileType, uuid.New().String(), fileExtension)

	// Save the cropped image
	userDir := filepath.Dir(media.FilePath)
	newFilePath := filepath.Join(userDir, newFileName)

	if err := s.saveImage(finalImg, format, newFilePath); err != nil {
		return nil, err
	}

	// Generate thumbnails
	thumbnailPath, err := s.generateThumbnails(finalImg, format, userDir, newFileName)
	if err != nil {
		os.Remove(newFilePath)
		return nil, err
	}

	// Delete old files
	os.Remove(media.FilePath)
	if media.ThumbnailPath != "" {
		os.Remove(media.ThumbnailPath)
	}
	s.deleteGeneratedThumbnails(media.FilePath)

	// Update database record
	bounds = finalImg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	media.FileName = newFileName
	media.FilePath = newFilePath
	media.ThumbnailPath = thumbnailPath
	media.Width = &width
	media.Height = &height
	media.UpdatedAt = time.Now()

	s.db.Save(&media)

	// Build response
	response := &MediaUploadResponse{
		MediaID:       media.ID,
		FileName:      media.FileName,
		OriginalName:  media.OriginalName,
		FilePath:      media.FilePath,
		ThumbnailPath: media.ThumbnailPath,
		URL:           s.getPublicURL(media.FilePath),
		ThumbnailURL:  s.getPublicURL(media.ThumbnailPath),
		MimeType:      media.MimeType,
		FileSize:      media.FileSize,
		Width:         width,
		Height:        height,
		FileType:      string(media.FileType),
	}

	return response, nil
}

// DeleteMedia removes a media file and all its thumbnails
func (s *MediaService) DeleteMedia(mediaID uint, userID uint) error {
	var media models.MediaFile

	// Find the media file
	if err := s.db.Where("id = ? AND user_id = ?", mediaID, userID).First(&media).Error; err != nil {
		return err
	}

	// Delete main file
	if media.FilePath != "" {
		os.Remove(media.FilePath)
	}

	// Delete thumbnail
	if media.ThumbnailPath != "" {
		os.Remove(media.ThumbnailPath)
	}

	// Delete all generated thumbnails
	s.deleteGeneratedThumbnails(media.FilePath)

	// Delete database record
	if err := s.db.Delete(&media).Error; err != nil {
		return err
	}

	return nil
}

// GetMedia retrieves a media file by ID
func (s *MediaService) GetMedia(mediaID uint, userID *uint) (*models.MediaFile, error) {
	var media models.MediaFile
	query := s.db.Where("id = ?", mediaID)

	if userID != nil {
		query = query.Where("user_id = ? OR is_public = ?", *userID, true)
	} else {
		query = query.Where("is_public = ?", true)
	}

	if err := query.First(&media).Error; err != nil {
		return nil, err
	}

	return &media, nil
}

// GetMediaByFilename retrieves a media file by filename
func (s *MediaService) GetMediaByFilename(filename string) (*models.MediaFile, error) {
	var media models.MediaFile
	if err := s.db.Where("file_name = ?", filename).First(&media).Error; err != nil {
		return nil, err
	}
	return &media, nil
}

// Helper functions

// saveImage saves an image to disk with the specified format
func (s *MediaService) saveImage(img image.Image, format string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch format {
	case "jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 85})
	case "png":
		return png.Encode(file, img)
	default:
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 85})
	}
}

// generateThumbnails creates multiple thumbnail sizes
func (s *MediaService) generateThumbnails(img image.Image, format string, dir string, filename string) (string, error) {
	var primaryThumbnail string

	for _, size := range s.thumbnailSizes {
		// Resize image to thumbnail size (maintain aspect ratio, fill center)
		thumbnail := imaging.Fill(img, size.Width, size.Height, imaging.Center, imaging.Lanczos)

		// Create thumbnail filename
		ext := filepath.Ext(filename)
		base := strings.TrimSuffix(filename, ext)
		thumbnailFilename := fmt.Sprintf("%s_%s%s", base, size.Name, ext)
		thumbnailPath := filepath.Join(dir, thumbnailFilename)

		// Save thumbnail
		if err := s.saveImage(thumbnail, format, thumbnailPath); err != nil {
			return "", err
		}

		// Set primary thumbnail (usually medium size)
		if size.Name == "medium" {
			primaryThumbnail = thumbnailPath
		}
	}

	return primaryThumbnail, nil
}

// deleteGeneratedThumbnails removes all generated thumbnails for a file
func (s *MediaService) deleteGeneratedThumbnails(filePath string) {
	if filePath == "" {
		return
	}

	dir := filepath.Dir(filePath)
	base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	for _, size := range s.thumbnailSizes {
		thumbnailPath := filepath.Join(dir, fmt.Sprintf("%s_%s%s", base, size.Name, filepath.Ext(filePath)))
		os.Remove(thumbnailPath)
	}
}

// deleteOldMedia removes old avatar/cover photos for a user
func (s *MediaService) deleteOldMedia(userID uint, fileType models.MediaType) {
	var oldMedia []models.MediaFile
	s.db.Where("user_id = ? AND file_type = ?", userID, fileType).Find(&oldMedia)

	for _, media := range oldMedia {
		// Delete files
		if media.FilePath != "" {
			os.Remove(media.FilePath)
		}
		if media.ThumbnailPath != "" {
			os.Remove(media.ThumbnailPath)
		}
		s.deleteGeneratedThumbnails(media.FilePath)

		// Delete database record
		s.db.Delete(&media)
	}
}

// isAllowedType checks if the MIME type is allowed
func (s *MediaService) isAllowedType(contentType string) bool {
	for _, allowedType := range s.allowedTypes {
		if strings.EqualFold(contentType, allowedType) {
			return true
		}
	}
	return false
}

// getFileExtension returns the file extension for a MIME type
func (s *MediaService) getFileExtension(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

// getPublicURL returns the public URL for a file path
func (s *MediaService) getPublicURL(filePath string) string {
	// Return the relative path for serving via static file handler
	return strings.TrimPrefix(filePath, "./")
}

// EnsureUploadDirectories creates the necessary upload directories
func (s *MediaService) EnsureUploadDirectories() error {
	dirs := []string{
		filepath.Join(s.uploadPath, "avatars"),
		filepath.Join(s.uploadPath, "covers"),
		filepath.Join(s.uploadPath, "media"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// updateProfileAvatar updates user's profile with new avatar reference
func (s *MediaService) updateProfileAvatar(userID uint, avatarID uint) error {
	return s.db.Model(&models.UserProfile{}).Where("user_id = ?", userID).Update("avatar_id", avatarID).Error
}

// updateProfileCoverPhoto updates user's profile with new cover photo reference
func (s *MediaService) updateProfileCoverPhoto(userID uint, coverPhotoID uint) error {
	return s.db.Model(&models.UserProfile{}).Where("user_id = ?", userID).Update("cover_photo_id", coverPhotoID).Error
}
