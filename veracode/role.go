package veracode

type Role struct {
	IsApi               bool   `json:"is_api,omitempty"`
	IsScanType          bool   `json:"is_scan_type,omitempty"`
	TeamAdminManageable bool   `json:"team_admin_manageable,omitempty"`
	RoleDescription     string `json:"role_description,omitempty"`
	RoleId              string `json:"role_id,omitempty"`
	RoleName            string `json:"role_name,omitempty"`
	RoleLegacyId        int    `json:"role_legacy_id,omitempty"`
	// BACKLOG: Add remaining fields for model as required.
}
