package transform

type OutputTransform struct {
	Type    string `json:"type"`
	Convert string `json:"convert"`
}

type OutputCriteria struct {
	Type  string `json:"type"`
	Match string `json:"match"`
}

type ConvertElement struct {
	Criteria  OutputCriteria  `json:"criteria"`
	Transform OutputTransform `json:"transform"`
}
