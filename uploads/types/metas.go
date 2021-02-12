package types

type SwingUploadMeta struct {
	TimestampSeconds int    `json:"timestamp"`
	Frames           int    `json:"frames"`
	Swing            int    `json:"swing"`
	Clip             int    `json:"clip"`
	UploadKey        string `json:"uploadKey"`
}
