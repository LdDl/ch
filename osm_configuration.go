package ch

// OsmConfiguration Allows to filter ways by certain tags from OSM data
type OsmConfiguration struct {
	EntityName string // Currrently we support 'highway' only
	Tags       []string
}

// CheckTag Checks if incoming tag is represented in configuration
func (cfg *OsmConfiguration) CheckTag(tag string) bool {
	for i := range cfg.Tags {
		if cfg.Tags[i] == tag {
			return true
		}
	}
	return false
}
