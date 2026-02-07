package api

import "time"

// DeviceCodeResponse is returned when requesting a device code
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// TokenResponse is returned when device code is exchanged for a token
type TokenResponse struct {
	AccessToken  string             `json:"access_token"`
	TokenType    string             `json:"token_type"`
	Scope        string             `json:"scope"`
	Organization OrganizationInfo   `json:"organization"`
	User         UserInfo           `json:"user"`
}

// TokenErrorResponse is returned when polling for token fails
type TokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// OrganizationInfo contains basic org info from token response
type OrganizationInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// UserInfo contains basic user info from token response
type UserInfo struct {
	Email string `json:"email"`
}

// BrandDNA represents the brand DNA from the API
type BrandDNA struct {
	ID          int                    `json:"id"`
	Version     int                    `json:"version"`
	Active      bool                   `json:"active"`
	Status      string                 `json:"status"`
	SourceURL   string                 `json:"source_url"`
	SourceType  string                 `json:"source_type"`
	CompanyName string                 `json:"company_name"`
	Palette     map[string]interface{} `json:"palette"`
	Fonts       map[string]interface{} `json:"fonts"`
	Voice       map[string]interface{} `json:"voice"`
	LogoURL     string                 `json:"logo_url,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// BrandVersionsResponse is returned from /api/v1/brand/versions
type BrandVersionsResponse struct {
	Versions []BrandDNA `json:"versions"`
}

// Task represents a marketing task from the API
type Task struct {
	ID                int        `json:"id"`
	Title             string     `json:"title"`
	Description       string     `json:"description,omitempty"`
	Status            string     `json:"status"`
	Position          int        `json:"position"`
	Platform          string     `json:"platform,omitempty"`
	AssetType         string     `json:"asset_type,omitempty"`
	SuggestedPostDate *string    `json:"suggested_post_date,omitempty"`
	MarketingPlanID   int        `json:"marketing_plan_id"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Detailed fields (only in single task response)
	ImagePrompt        string          `json:"image_prompt,omitempty"`
	CopyPrompt         string          `json:"copy_prompt,omitempty"`
	VideoPrompt        string          `json:"video_prompt,omitempty"`
	SelectedImagesCount int            `json:"selected_images_count,omitempty"`
	SelectedCopy       *SelectedCopy   `json:"selected_copy,omitempty"`
	ReadyForCompletion bool            `json:"ready_for_completion,omitempty"`
	GithubEventID      *int            `json:"github_event_id,omitempty"`
}

// SelectedCopy represents selected copy for a task
type SelectedCopy struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

// TasksResponse is returned from /api/v1/tasks
type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

// ContextResponse is returned from /api/v1/context
type ContextResponse struct {
	Organization  ContextOrganization   `json:"organization"`
	Brand         *ContextBrand         `json:"brand,omitempty"`
	MarketingPlan *ContextMarketingPlan `json:"marketing_plan,omitempty"`
	Tasks         ContextTasks          `json:"tasks"`
	GeneratedAt   time.Time             `json:"generated_at"`
}

// ContextOrganization contains org info for context
type ContextOrganization struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ContextBrand contains brand info for context
type ContextBrand struct {
	CompanyName string                 `json:"company_name"`
	SourceURL   string                 `json:"source_url"`
	Palette     map[string]interface{} `json:"palette"`
	Fonts       map[string]interface{} `json:"fonts"`
	Voice       map[string]interface{} `json:"voice"`
	Version     int                    `json:"version"`
}

// ContextMarketingPlan contains marketing plan info for context
type ContextMarketingPlan struct {
	ID               int       `json:"id"`
	Content          string    `json:"content"`
	GenerationStatus string    `json:"generation_status"`
	TaskCount        int       `json:"task_count"`
	PendingTasks     int       `json:"pending_tasks"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ContextTasks contains task summary for context
type ContextTasks struct {
	Total      int               `json:"total"`
	ByStatus   map[string]int    `json:"by_status"`
	Pending    []TaskSummary     `json:"pending"`
	InProgress []TaskSummary     `json:"in_progress"`
}

// TaskSummary is a brief task summary for context
type TaskSummary struct {
	ID                int     `json:"id"`
	Title             string  `json:"title"`
	Description       string  `json:"description,omitempty"`
	Platform          string  `json:"platform,omitempty"`
	SuggestedPostDate *string `json:"suggested_post_date,omitempty"`
	ImagePrompt       string  `json:"image_prompt,omitempty"`
	CopyPrompt        string  `json:"copy_prompt,omitempty"`
}

// UpdateTaskRequest represents the request body for updating a task
type UpdateTaskRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

// DiscardResponse represents the response from discarding a resource
type DiscardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error string `json:"error"`
}
