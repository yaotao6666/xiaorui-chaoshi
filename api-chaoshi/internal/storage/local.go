package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"chaoshi_api/internal/config"
)

type LocalStorage struct {
	rootDir       string
	publicBaseURL string
	urlPrefix     string
}

var localStorage *LocalStorage

func InitLocalStorage() error {
	cfg := config.Config.Storage
	rootDir := strings.TrimSpace(cfg.RootDir)
	if rootDir == "" {
		rootDir = "./uploads"
	}

	absoluteRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return fmt.Errorf("解析上传目录失败: %w", err)
	}
	if err := os.MkdirAll(absoluteRoot, 0o755); err != nil {
		return fmt.Errorf("创建上传目录失败: %w", err)
	}

	urlPrefix := normalizeURLPrefix(cfg.URLPrefix)
	publicBaseURL := strings.TrimRight(strings.TrimSpace(cfg.PublicBaseURL), "/")

	localStorage = &LocalStorage{
		rootDir:       absoluteRoot,
		publicBaseURL: publicBaseURL,
		urlPrefix:     urlPrefix,
	}
	return nil
}

func GetService() *LocalStorage {
	return localStorage
}

func (s *LocalStorage) RootDir() string {
	if s == nil {
		return ""
	}
	return s.rootDir
}

func (s *LocalStorage) URLPrefix() string {
	if s == nil {
		return "/uploads"
	}
	return s.urlPrefix
}

func (s *LocalStorage) SaveUploadedFile(fileHeader *multipart.FileHeader, scope string) (string, string, error) {
	if s == nil {
		return "", "", fmt.Errorf("本地存储未初始化")
	}
	if fileHeader == nil {
		return "", "", fmt.Errorf("上传文件不能为空")
	}

	relativeDir := buildRelativeDir(scope)
	targetDir := filepath.Join(s.rootDir, filepath.FromSlash(relativeDir))
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", "", fmt.Errorf("创建上传子目录失败: %w", err)
	}

	filename := buildFilename(fileHeader.Filename)
	relativePath := filepath.ToSlash(filepath.Join(relativeDir, filename))
	targetPath := filepath.Join(s.rootDir, filepath.FromSlash(relativePath))

	src, err := fileHeader.Open()
	if err != nil {
		return "", "", fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		return "", "", fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", "", fmt.Errorf("保存上传文件失败: %w", err)
	}

	return relativePath, s.BuildURL(relativePath), nil
}

func (s *LocalStorage) BuildURL(resource string) string {
	if s == nil {
		return resource
	}

	trimmed := strings.TrimSpace(resource)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return trimmed
	}

	normalizedPath := filepath.ToSlash(strings.TrimLeft(trimmed, "/"))
	if s.publicBaseURL != "" {
		return s.publicBaseURL + s.urlPrefix + "/" + normalizedPath
	}
	return s.urlPrefix + "/" + normalizedPath
}

func buildRelativeDir(scope string) string {
	now := time.Now()
	scopePath := strings.Trim(scope, "/")
	if scopePath == "" {
		scopePath = "common"
	}
	return filepath.ToSlash(filepath.Join(scopePath, now.Format("2006"), now.Format("01"), now.Format("02")))
}

func buildFilename(original string) string {
	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(original)))
	if ext == "" {
		ext = ".bin"
	}

	base := strings.TrimSuffix(filepath.Base(strings.TrimSpace(original)), ext)
	base = sanitizeSegment(base)
	if base == "" {
		base = "file"
	}

	return fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), base, ext)
}

func sanitizeSegment(value string) string {
	replacer := strings.NewReplacer(
		"\\", "_",
		"/", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		" ", "_",
	)
	return strings.Trim(replacer.Replace(value), "._")
}

func normalizeURLPrefix(prefix string) string {
	trimmed := strings.TrimSpace(prefix)
	if trimmed == "" {
		return "/uploads"
	}
	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}
	return strings.TrimRight(trimmed, "/")
}
