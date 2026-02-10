package slack

import (
	"io"
	"os"

	slackapi "github.com/slack-go/slack"
)

type File struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Mimetype   string `json:"mimetype"`
	Filetype   string `json:"filetype"`
	Size       int    `json:"size"`
	User       string `json:"user"`
	Created    int64  `json:"created"`
	URLPrivate string `json:"url_private"`
	Permalink  string `json:"permalink"`
}

func fileFromAPI(f slackapi.File) File {
	return File{
		ID:         f.ID,
		Name:       f.Name,
		Title:      f.Title,
		Mimetype:   f.Mimetype,
		Filetype:   f.Filetype,
		Size:       f.Size,
		User:       f.User,
		Created:    int64(f.Created),
		URLPrivate: f.URLPrivateDownload,
		Permalink:  f.Permalink,
	}
}

func (c *Client) ListFiles(params PaginationParams, channelID, userID string) (*PaginatedResult[File], error) {
	if params.All {
		return c.listAllFiles(params, channelID, userID)
	}
	return c.listFilesPage(params, channelID, userID)
}

func (c *Client) listFilesPage(params PaginationParams, channelID, userID string) (*PaginatedResult[File], error) {
	listParams := slackapi.ListFilesParameters{
		Channel: channelID,
		User:    userID,
		Limit:   params.EffectiveLimit(),
		Cursor:  params.Cursor,
	}

	type listResult struct {
		files      []slackapi.File
		nextParams *slackapi.ListFilesParameters
	}

	r, err := retry(func() (listResult, error) {
		files, nextParams, err := c.api.ListFiles(listParams)
		return listResult{files, nextParams}, err
	})
	if err != nil {
		return nil, classifyError(err)
	}

	items := make([]File, len(r.files))
	for i, f := range r.files {
		items[i] = fileFromAPI(f)
	}

	var nextCursor string
	if r.nextParams != nil {
		nextCursor = r.nextParams.Cursor
	}

	return &PaginatedResult[File]{
		Items:      items,
		NextCursor: nextCursor,
		HasMore:    nextCursor != "",
	}, nil
}

func (c *Client) listAllFiles(params PaginationParams, channelID, userID string) (*PaginatedResult[File], error) {
	var allItems []File
	cursor := params.Cursor
	for {
		page, err := c.listFilesPage(PaginationParams{
			Cursor: cursor,
			Limit:  params.EffectiveLimit(),
		}, channelID, userID)
		if err != nil {
			return nil, err
		}
		allItems = append(allItems, page.Items...)
		if !page.HasMore {
			break
		}
		cursor = page.NextCursor
	}
	return &PaginatedResult[File]{
		Items:   allItems,
		HasMore: false,
	}, nil
}

func (c *Client) GetFileInfo(fileID string) (*File, error) {
	type fileInfoResult struct {
		file *slackapi.File
	}

	r, err := retry(func() (fileInfoResult, error) {
		f, _, _, err := c.api.GetFileInfo(fileID, 0, 0)
		return fileInfoResult{f}, err
	})
	if err != nil {
		return nil, classifyError(err)
	}
	result := fileFromAPI(*r.file)
	return &result, nil
}

func (c *Client) UploadFile(channelID, filename, title string, reader io.Reader) (*File, error) {
	params := slackapi.UploadFileV2Parameters{
		Channel:  channelID,
		Filename: filename,
		Title:    title,
		Reader:   reader,
	}

	f, err := retry(func() (*slackapi.FileSummary, error) {
		return c.api.UploadFileV2(params)
	})
	if err != nil {
		return nil, classifyError(err)
	}
	return &File{
		ID:    f.ID,
		Title: f.Title,
	}, nil
}

func (c *Client) DownloadFile(url, destPath string) error {
	outFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func() { _ = outFile.Close() }()

	err = c.api.GetFile(url, outFile)
	if err != nil {
		return classifyError(err)
	}
	return nil
}

func (c *Client) DeleteFile(fileID string) error {
	_, err := retry(func() (struct{}, error) {
		return struct{}{}, c.api.DeleteFile(fileID)
	})
	if err != nil {
		return classifyError(err)
	}
	return nil
}
