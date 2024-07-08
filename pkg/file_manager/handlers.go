package filemanager

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

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

func ListFilesHandler(c *gin.Context) {
	files, err := GetRootItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	Files = files

	c.JSON(http.StatusOK, Files)
}

func CreateFileHandler(c *gin.Context) {
	folder := c.Param("folder")
	if folder == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "folder parameter is required"})
		return
	}

	if err := CreateFolder(folder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func DeleteFileHandler(c *gin.Context) {
	var requestBody struct {
		Path string `json:"path"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path := requestBody.Path
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path query parameter is required"})
		return
	}

	if err := DeleteItem(Files, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
