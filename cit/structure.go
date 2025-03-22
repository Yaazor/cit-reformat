package cit

type SourceCIT struct {
	Type       string `properties:"type,default=item"`
	Item       string `properties:"items,default=wooden_sword"`
	CustomData string `properties:"custom_data,default=none"`
	Pattern    string `properties:"components.custom_name,default=none"`
}
