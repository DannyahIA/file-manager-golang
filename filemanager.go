package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
)

const dirName = "./drive"

type File struct {
	Name         string `json:"name,omitempty"`
	Path         string `json:"path,omitempty"`
	IsFolder     bool   `json:"is_folder"`
	Size         string `json:"size"`
	DataModified string `json:"data_modified,omitempty"`
	Items        []File `json:"items,omitempty"`
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

func GetRootItems() ([]File, error) {
	var folders []File
	var files []File

	err := filepath.WalkDir(dirName, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == dirName {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		parentDir := filepath.ToSlash(filepath.Dir(path))
		if !d.IsDir() {
			found := false
			for i := range folders {
				if folders[i].Path == parentDir {
					folders[i].Items = append(folders[i].Items, File{
						Name:         d.Name(),
						Path:         filepath.ToSlash(path),
						IsFolder:     false,
						Items:        nil,
						Size:         convertSizeToMB(fileInfo.Size()),
						DataModified: fileInfo.ModTime().String(),
					})
					found = true
					break
				}
			}
			if !found {
				files = append(files, File{
					Name:         d.Name(),
					Path:         filepath.ToSlash(path),
					IsFolder:     false,
					Items:        nil,
					Size:         convertSizeToMB(fileInfo.Size()),
					DataModified: fileInfo.ModTime().String(),
				})
			}
		} else {
			found := false
			for i := range folders {
				if folders[i].Path == parentDir {
					folders[i].Items = append(folders[i].Items, File{
						Name:         d.Name(),
						Path:         filepath.ToSlash(path),
						IsFolder:     true,
						Items:        nil,
						Size:         convertSizeToMB(fileInfo.Size()),
						DataModified: fileInfo.ModTime().String(),
					})
					found = true
					break
				}
			}
			if !found {
				folders = append(folders, File{
					Name:         d.Name(),
					Path:         filepath.ToSlash(path),
					IsFolder:     true,
					Items:        nil,
					Size:         convertSizeToMB(fileInfo.Size()),
					DataModified: fileInfo.ModTime().String(),
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for i := range folders {
		for j := range folders[i].Items {
			if folders[i].Items[j].IsFolder {
				folders[i].Items = nil
				break
			}
		}
	}

	var result []File
	result = append(result, folders...)
	result = append(result, files...)

	return result, nil
}

func CreateFolder(folderName string) error {
	return os.Mkdir(filepath.Join(dirName, folderName), 0755)
}

func DeleteItem(files []File, path string) error {
	exeDir, err := os.Executable()
	for i := range files {
		if files[i].Path == path {
			if files[i].IsFolder {
				err := os.RemoveAll(path)
				if err != nil {
					return err
				}
			} else {
				if err != nil {
					return err
				}

				err = os.Remove(filepath.Join(filepath.Base(exeDir), path))
				if err != nil {
					return err
				}
			}
			return nil
		}

		for j := range files[i].Items {
			if files[i].Items[j].Path == path {
				if files[i].IsFolder {
					err := os.RemoveAll(path)
					if err != nil {
						return err
					}
				} else {
					if err != nil {
						return err
					}

					err = os.Remove(filepath.Join(filepath.Base(exeDir), path))
					if err != nil {
						return err
					}
				}
				return nil
			}
		}
	}
	return fmt.Errorf("file or directory not found")
}
