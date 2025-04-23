package entity

import (
	"fmt"
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", ct.Format("2006-01-02"))
	return []byte(formatted), nil
}

// Restriction struct corresponds to the Restriction class in Java, used to store restriction information
type Restriction struct {
	Plan            string `json:"plan"`              // Default Go JSON serialization will use field name as JSON key
	ActiveUserCount int    `json:"active_user_count"` // Using JsonProperty tag to map JSON field
}

// LicenseInfo struct corresponds to the LicenseInfo class in Java, used to store license information
type LicenseInfo struct {
	Name    string `json:"Name"`
	Company string `json:"Company"`
	Email   string `json:"Email"`
}

// License represents a license, corresponding to the License class in Java
type License struct {
	Version                      int         `json:"version"`
	License                      LicenseInfo `json:"licensee"`
	StartsAt                     CustomTime  `json:"issued_at"`
	ExpiresAt                    CustomTime  `json:"expires_at"`
	NotifyAdminsAt               CustomTime  `json:"notify_admins_at"`
	NotifyUsersAt                CustomTime  `json:"notify_users_at"`
	BlockChangesAt               CustomTime  `json:"block_changes_at"`
	CloudLicensingEnabled        bool        `json:"cloud_licensing_enabled"`
	OfflineCloudLicensingEnabled bool        `json:"offline_cloud_licensing_enabled"`
	AutoRenewEnabled             bool        `json:"auto_renew_enabled"`
	SeatReconciliationEnabled    bool        `json:"seat_reconciliation_enabled"`
	OperationalMetricsEnabled    bool        `json:"operational_metrics_enabled"`
	GeneratedFromCustomersDot    bool        `json:"generated_from_customers_dot"`
	Restrictions                 Restriction `json:"restrictions"`
}
