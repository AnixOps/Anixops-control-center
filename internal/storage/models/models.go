package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the control center
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;size:255" json:"email"`
	PasswordHash string         `gorm:"size:255" json:"-"`
	Role         string         `gorm:"size:50;default:viewer" json:"role"`
	AuthProvider string         `gorm:"size:50;default:local" json:"auth_provider"`
	Enabled      bool           `gorm:"default:true" json:"enabled"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name
func (User) TableName() string {
	return "users"
}

// APIToken represents an API token
type APIToken struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index" json:"user_id"`
	Name      string         `gorm:"size:100" json:"name"`
	Token     string         `gorm:"uniqueIndex;size:255" json:"token"`
	ExpiresAt *time.Time     `json:"expires_at,omitempty"`
	LastUsed  *time.Time     `json:"last_used,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name
func (APIToken) TableName() string {
	return "api_tokens"
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Action    string    `gorm:"size:100;index" json:"action"`
	Resource  string    `gorm:"size:100;index" json:"resource"`
	ResourceID string   `gorm:"size:100;index" json:"resource_id"`
	IP        string    `gorm:"size:50" json:"ip"`
	UserAgent string    `gorm:"size:255" json:"user_agent"`
	Details   string    `gorm:"type:text" json:"details"`
	Status    string    `gorm:"size:20;default:success" json:"status"` // success, failed
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

// TableName returns the table name
func (AuditLog) TableName() string {
	return "audit_logs"
}

// Node represents a proxy node in the system
type Node struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"uniqueIndex;size:100" json:"name"`
	Host           string         `gorm:"size:255" json:"host"`
	Port           int            `gorm:"default:443" json:"port"`
	Type           string         `gorm:"size:50;default:'v2ray'" json:"type"` // v2ray, xray, sing-box, etc.
	Status         string         `gorm:"size:20;default:'unknown'" json:"status"` // online, offline, unknown, maintenance
	Region         string         `gorm:"size:100" json:"region"`
	ServerID       uint           `gorm:"index" json:"server_id"`
	TrafficUp      int64          `gorm:"default:0" json:"traffic_up"`       // bytes
	TrafficDown    int64          `gorm:"default:0" json:"traffic_down"`     // bytes
	UserCount      int            `gorm:"default:0" json:"user_count"`
	MaxUsers       int            `gorm:"default:0" json:"max_users"`
	Enabled        bool           `gorm:"default:true" json:"enabled"`
	LastCheckedAt  *time.Time     `json:"last_checked_at,omitempty"`
	Tags           string         `gorm:"size:255" json:"tags"` // JSON array of tags
	Config         string         `gorm:"type:text" json:"config"` // JSON config
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name
func (Node) TableName() string {
	return "nodes"
}

// Server represents a physical/virtual server
type Server struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"uniqueIndex;size:100" json:"name"`
	IP           string         `gorm:"size:50" json:"ip"`
	Provider     string         `gorm:"size:100" json:"provider"` // AWS, GCP, DigitalOcean, etc.
	Region       string         `gorm:"size:100" json:"region"`
	CPU          int            `gorm:"default:0" json:"cpu"`         // cores
	Memory       int            `gorm:"default:0" json:"memory"`      // MB
	Disk         int            `gorm:"default:0" json:"disk"`        // GB
	Bandwidth    int64          `gorm:"default:0" json:"bandwidth"`   // bytes/month
	Status       string         `gorm:"size:20;default:'unknown'" json:"status"` // active, inactive, maintenance
	NodeCount    int            `gorm:"default:0" json:"node_count"`
	SSHHost      string         `gorm:"size:255" json:"ssh_host"`
	SSHPort      int            `gorm:"default:22" json:"ssh_port"`
	SSHUser      string         `gorm:"size:100" json:"ssh_user"`
	SSHKey       string         `gorm:"type:text" json:"ssh_key"`     // encrypted SSH key
	Notes        string         `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name
func (Server) TableName() string {
	return "servers"
}

// Subscription represents a user subscription
type Subscription struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index" json:"user_id"`
	PlanID        uint           `gorm:"index" json:"plan_id"`
	Status        string         `gorm:"size:20;default:'active'" json:"status"` // active, expired, cancelled, suspended
	StartDate     time.Time      `json:"start_date"`
	EndDate       time.Time      `json:"end_date"`
	TrafficUsed   int64          `gorm:"default:0" json:"traffic_used"`    // bytes
	TrafficLimit  int64          `gorm:"default:0" json:"traffic_limit"`   // bytes, 0 = unlimited
	NodeIDs       string         `gorm:"type:text" json:"node_ids"`        // JSON array of allowed node IDs
	AutoRenew     bool           `gorm:"default:false" json:"auto_renew"`
	PaymentID     string         `gorm:"size:100" json:"payment_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name
func (Subscription) TableName() string {
	return "subscriptions"
}

// Plan represents a subscription plan
type Plan struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"uniqueIndex;size:100" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Price        float64        `json:"price"`
	Currency     string         `gorm:"size:10;default:'USD'" json:"currency"`
	Duration     int            `json:"duration"` // days
	TrafficLimit int64          `json:"traffic_limit"` // bytes, 0 = unlimited
	NodeLimit    int            `gorm:"default:0" json:"node_limit"` // max nodes, 0 = unlimited
	Features     string         `gorm:"type:text" json:"features"` // JSON array of features
	SortOrder    int            `gorm:"default:0" json:"sort_order"`
	Enabled      bool           `gorm:"default:true" json:"enabled"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name
func (Plan) TableName() string {
	return "plans"
}

// Order represents a purchase order
type Order struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"index" json:"user_id"`
	PlanID         uint           `gorm:"index" json:"plan_id"`
	SubscriptionID uint           `gorm:"index" json:"subscription_id"`
	Amount         float64        `json:"amount"`
	Currency       string         `gorm:"size:10;default:'USD'" json:"currency"`
	Status         string         `gorm:"size:20;default:'pending'" json:"status"` // pending, paid, cancelled, refunded
	PaymentMethod  string         `gorm:"size:50" json:"payment_method"`
	PaymentID      string         `gorm:"size:100" json:"payment_id"`
	PaidAt         *time.Time     `json:"paid_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name
func (Order) TableName() string {
	return "orders"
}
