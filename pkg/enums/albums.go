package enums

type SwingVideoStatus string

const (
	SwingVideoStatusProcessing SwingVideoStatus = "Processing"
	SwingVideoStatusCreated    SwingVideoStatus = "Created"
	SwingVideoStatusDeleted    SwingVideoStatus = "Deleted"
)

type AlbumStatus string

const (
	AlbumStatusProcessing AlbumStatus = "Processing"
	AlbumStatusClipped    AlbumStatus = "Clipped"
	AlbumStatusCreated    AlbumStatus = "Created"
	AlbumStatusDeleted    AlbumStatus = "Deleted"
)

type SwingUploadStatus string

const (
	SwingUploadStatusOriginal SwingUploadStatus = "Original"
	SwingUploadStatusClipped  SwingUploadStatus = "Clipped"
	SwingUploadStatusFinished SwingUploadStatus = "Finished"
)
