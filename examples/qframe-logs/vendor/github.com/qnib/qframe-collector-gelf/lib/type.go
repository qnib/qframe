package qframe_collector_gelf


type GelfMsg struct {
  Version string `json:"version"`
  Host string `json:"host"`
  ShortMsg string `json:"short_message"`
  TimeNano float64 `json:"timestamp"`
  Level int `json:"level"`
  Command string `json:"_command"`
  ContainerID string `json:"_container_id"`
  ContainerName string `json:"_container_name"`
  ImageID string `json:"_image_id"` // Not the Digest!
  ImageName string `json:"_image_name"`
  Tag string `json:"_tag"`
  SourceAddr string
}
