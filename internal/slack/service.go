package slack

import "io"

//go:generate mockgen -source=service.go -destination=mocks/mock_service.go -package=mocks

// Service defines the interface for all Slack API operations.
// *Client satisfies this interface.
type Service interface {
	AuthTest() (*AuthTestResult, error)

	ListChannels(params PaginationParams) (*PaginatedResult[Channel], error)
	GetChannelInfo(channelID string) (*Channel, error)
	CreateChannel(name string, isPrivate bool) (*Channel, error)
	ArchiveChannel(channelID string) error
	InviteToChannel(channelID string, userIDs ...string) error
	KickFromChannel(channelID, userID string) error
	SetChannelTopic(channelID, topic string) error
	SetChannelPurpose(channelID, purpose string) error

	ListMessages(params ListMessagesParams) (*PaginatedResult[Message], error)
	SendMessage(params SendMessageParams) (*Message, error)
	EditMessage(channelID, timestamp, text string) (*Message, error)
	DeleteMessage(channelID, timestamp string) error
	SearchMessages(params SearchParams) (*SearchResult, error)

	ListUsers(params PaginationParams) (*PaginatedResult[User], error)
	GetUserInfo(userID string) (*User, error)
	GetUserPresence(userID string) (string, error)

	AddReaction(channelID, timestamp, name string) error
	RemoveReaction(channelID, timestamp, name string) error
	ListReactions(userID string, params PaginationParams) (*PaginatedResult[ReactedItem], error)

	ListFiles(params PaginationParams, channelID, userID string) (*PaginatedResult[File], error)
	GetFileInfo(fileID string) (*File, error)
	UploadFile(channelID, filename, title string, reader io.Reader) (*File, error)
	DownloadFile(url, destPath string) error
	DeleteFile(fileID string) error
}
