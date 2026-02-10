package files

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jackchuka/slackcli/internal/cmdutil"
	"github.com/jackchuka/slackcli/internal/slack"
)

func NewFilesCmd() *cobra.Command {
	filesCmd := &cobra.Command{
		Use:   "files",
		Short: "Manage files",
	}
	filesCmd.AddCommand(newListCmd())
	filesCmd.AddCommand(newInfoCmd())
	filesCmd.AddCommand(newUploadCmd())
	filesCmd.AddCommand(newDownloadCmd())
	filesCmd.AddCommand(newDeleteCmd())
	return filesCmd
}

func newListCmd() *cobra.Command {
	var channelID string
	var userID string
	var cursor string
	var limit int
	var all bool

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List files",
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			result, err := rc.Client.ListFiles(slack.PaginationParams{
				Cursor: cursor,
				Limit:  limit,
				All:    all,
			}, channelID, userID)
			if err != nil {
				return err
			}
			return rc.Formatter.Format(result)
		},
	}
	listCmd.Flags().StringVar(&channelID, "channel", "", "Filter by channel ID")
	listCmd.Flags().StringVar(&userID, "user", "", "Filter by user ID")
	listCmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")
	listCmd.Flags().IntVar(&limit, "limit", 100, "Number of files per page")
	listCmd.Flags().BoolVar(&all, "all", false, "Fetch all files (auto-paginate)")
	return listCmd
}

func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info <file-id>",
		Short: "Get file info",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			file, err := rc.Client.GetFileInfo(args[0])
			if err != nil {
				return err
			}
			return rc.Formatter.Format(file)
		},
	}
}

func newUploadCmd() *cobra.Command {
	var channelID string
	var title string
	var filePath string

	uploadCmd := &cobra.Command{
		Use:         "upload",
		Short:       "Upload a file",
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			f, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer func() { _ = f.Close() }()

			filename := filePath
			if len(args) > 0 {
				filename = args[0]
			}

			file, err := rc.Client.UploadFile(channelID, filename, title, f)
			if err != nil {
				return err
			}
			return rc.Formatter.Format(file)
		},
	}
	uploadCmd.Flags().StringVar(&channelID, "channel", "", "Channel ID (required)")
	_ = uploadCmd.MarkFlagRequired("channel")
	uploadCmd.Flags().StringVar(&filePath, "file", "", "File path (required)")
	_ = uploadCmd.MarkFlagRequired("file")
	uploadCmd.Flags().StringVar(&title, "title", "", "File title")
	return uploadCmd
}

func newDownloadCmd() *cobra.Command {
	var dest string

	downloadCmd := &cobra.Command{
		Use:   "download <file-id>",
		Short: "Download a file",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			file, err := rc.Client.GetFileInfo(args[0])
			if err != nil {
				return err
			}
			destPath := dest
			if destPath == "" {
				destPath = file.Name
			}
			if err := rc.Client.DownloadFile(file.URLPrivate, destPath); err != nil {
				return err
			}
			return rc.Formatter.Format(map[string]string{
				"status": "downloaded",
				"file":   args[0],
				"path":   destPath,
			})
		},
	}
	downloadCmd.Flags().StringVarP(&dest, "dest", "d", "", "Destination path")
	return downloadCmd
}

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "delete <file-id>",
		Short:       "Delete a file",
		Args:        cobra.ExactArgs(1),
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			if err := rc.Client.DeleteFile(args[0]); err != nil {
				return err
			}
			return rc.Formatter.Format(map[string]string{
				"status": "deleted",
				"file":   args[0],
			})
		},
	}
}
