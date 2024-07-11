package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var DefaultRoot = "./drive"

type File struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	IsFolder     bool   `json:"is_folder"`
	Size         string `json:"size"`
	Extension    string `json:"extension"`
	DataModified string `json:"data_modified"`
}

func convertSizeToMB(size int64) string {
	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB
		GB
		TB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	}

	return fmt.Sprintf("%d bytes", size)
}

func GetRootFolders() ([]File, error) {
	var folders []File

	err := filepath.WalkDir(DefaultRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == DefaultRoot {
			return nil
		}

		// Check if the path is a subdirectory of DefaultRoot
		relPath, err := filepath.Rel(DefaultRoot, path)
		if err != nil {
			return err
		}
		if strings.Contains(relPath, string(filepath.Separator)) {
			return filepath.SkipDir
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		if d.IsDir() {
			folders = append(folders, File{
				Name:         d.Name(),
				Path:         filepath.ToSlash(path),
				IsFolder:     true,
				Size:         convertSizeToMB(fileInfo.Size()),
				DataModified: fileInfo.ModTime().Format(time.RFC3339),
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return folders, nil
}

func GetFolderItems(folderPath string) ([]File, error) {
	var items []File
	var folders []File
	var files []File

	err := filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == folderPath {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		if d.IsDir() {
			folders = append(folders, File{
				Name:         d.Name(),
				Path:         filepath.ToSlash(path),
				IsFolder:     true,
				Size:         convertSizeToMB(fileInfo.Size()),
				DataModified: fileInfo.ModTime().Format(time.RFC3339),
			})
			return filepath.SkipDir
		}

		files = append(files, File{
			Name:         d.Name(),
			Path:         filepath.ToSlash(path),
			IsFolder:     false,
			Size:         convertSizeToMB(fileInfo.Size()),
			DataModified: fileInfo.ModTime().Format(time.RFC3339),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	items = append(folders, files...)
	return items, nil
}

func CreateFolder(folderName string) error {
	return os.Mkdir(filepath.Join(DefaultRoot, folderName), 0755)
}

func DeleteItem(isFolder bool, path string) error {
	if isFolder {
		err := os.RemoveAll(filepath.Join(".", path))
		if err != nil {
			return err
		}
	} else {
		err := os.Remove(filepath.Join(".", path))
		if err != nil {
			return err
		}
	}
	return nil
}

func RenameItem(isFolder bool, path, newName string) error {
	dir := filepath.Dir(path)
	var newPath string
	if !isFolder {
		ext := filepath.Ext(path)
		newPath = filepath.Join(dir, newName+ext)
	} else {
		newPath = filepath.Join(dir, newName)
	}

	return os.Rename(path, newPath)
}
