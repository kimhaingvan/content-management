package metax

type MetaxType string

const (
	Datetime         MetaxType = "datetime"
	SelectOne        MetaxType = "select_one"
	SelectMany       MetaxType = "select_many"
	String           MetaxType = "string"
	Checkbox         MetaxType = "checkbox"
	File             MetaxType = "file"
	Number           MetaxType = "number"
	Float            MetaxType = "float"
	SingleEdit       MetaxType = "single_edit"
	CollectionEdit   MetaxType = "collection_edit"
	RichEditor       MetaxType = "rich_editor"
	HiddenPrimaryKey MetaxType = "hidden_primary_key"
	Hidden           MetaxType = "hidden"
	Readonly         MetaxType = "readonly"
	MediaLibrary     MetaxType = "media_library"
	MediaBox         MetaxType = "media_box"
	PublishLiveNow   MetaxType = "publish_live_now"
	Password         MetaxType = "password"
	Text             MetaxType = "text"
)

func (m MetaxType) String() string {
	return string(m)
}
