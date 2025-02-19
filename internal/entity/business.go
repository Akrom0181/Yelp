package entity

type ContactInfo struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

type HoursOfOperation struct {
	Monday    string `json:"monday"`
	Tuesday   string `json:"tuesday"`
	Wednesday string `json:"wednesday"`
	Thursday  string `json:"thursday"`
	Friday    string `json:"friday"`
	Saturday  string `json:"saturday"`
	Sunday    string `json:"sunday"`
}

type Business struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	Description      string               `json:"description"`
	CategoryID       string               `json:"category_id"`
	Address          string               `json:"address"`
	Attachments      []BusinessAttachment `json:"attachments"`
	Latitude         float64              `json:"latitude"`
	Longitude        float64              `json:"longitude"`
	ContactInfo      ContactInfo          `json:"contact_info"`
	HoursOfOperation HoursOfOperation     `json:"hours_of_operation"`
	OwnerID          string               `json:"owner_id"`
	CreatedAt        string               `json:"created_at"`
	UpdatedAt        string               `json:"updated_at"`
}

type BusinessList struct {
	Items []Business `json:"businesses"`
	Count int        `json:"count"`
}

type BusinessSingleRequest struct {
	ID         string `json:"id"`
	OwnerID    string `json:"owner_id"`
	CategoryID string `json:"category_id"`
}

type BusinessCategory struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BusinessCategoryList struct {
	Items []BusinessCategory `json:"businesses_categories"`
	Count int                `json:"count"`
}

type BusinessCategorySingleRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BusinessAttachment struct {
	Id          string `json:"id"`
	BusinessId  string `json:"-"`
	FilePath    string `json:"filepath"`
	ContentType string `json:"content_type"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type BusinessAttachmentList struct {
	Items []BusinessAttachment `json:"items"`
	Count int64                `json:"count"`
}

type BusinessAttachmentMultipleInsertRequest struct {
	BusinessId  string               `json:"business_id"`
	Attachments []BusinessAttachment `json:"attachments"`
}
