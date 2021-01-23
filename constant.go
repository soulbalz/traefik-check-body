package checkbodyplugin

// MatchType defines an enum which can be used to specify the match type for the 'contains' config.
type MatchType string

const (
	//MatchAll requires all values to be matched
	MatchAll MatchType = "all"
	//MatchOne requires only one value to be matched
	MatchOne MatchType = "one"
)
