package services

import (
	"errors"
	"time"

	"github.com/anixops/anixops-control-center/internal/security/auth"
	"github.com/anixops/anixops-control-center/internal/storage/models"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserDisabled       = errors.New("user is disabled")
	ErrEmailExists        = errors.New("email already exists")
)

// AuthService handles authentication operations
type AuthService struct {
	db        *gorm.DB
	jwtMgr    *auth.JWTManager
}

// NewAuthService creates a new auth service
func NewAuthService(db *gorm.DB, jwtMgr *auth.JWTManager) *AuthService {
	return &AuthService{
		db:     db,
		jwtMgr: jwtMgr,
	}
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(email, password, ip, userAgent string) (*models.User, string, string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", ErrInvalidCredentials
		}
		return nil, "", "", err
	}

	if !user.Enabled {
		return nil, "", "", ErrUserDisabled
	}

	if !auth.CheckPassword(password, user.PasswordHash) {
		return nil, "", "", ErrInvalidCredentials
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	s.db.Save(&user)

	// Generate tokens
	accessToken, err := s.jwtMgr.GenerateAccessToken(
		string(rune(user.ID)),
		user.Role,
		nil,
		user.AuthProvider,
	)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.jwtMgr.GenerateRefreshToken(
		string(rune(user.ID)),
		user.Role,
		user.AuthProvider,
	)
	if err != nil {
		return nil, "", "", err
	}

	return &user, accessToken, refreshToken, nil
}

// Register creates a new user
func (s *AuthService) Register(email, password, role string) (*models.User, error) {
	// Check if email exists
	var count int64
	s.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return nil, ErrEmailExists
	}

	// Hash password
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		AuthProvider: "local",
		Enabled:      true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// RefreshToken refreshes an access token
func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	return s.jwtMgr.RefreshAccessToken(refreshToken)
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(token string) (*auth.Claims, error) {
	return s.jwtMgr.ValidateToken(token)
}

// UserService handles user operations
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// List retrieves users with pagination
func (s *UserService) List(page, perPage int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := s.db.Model(&models.User{})
	if search != "" {
		query = query.Where("email LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Get retrieves a user by ID
func (s *UserService) Get(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Create creates a new user
func (s *UserService) Create(email, password, role string) (*models.User, error) {
	// Check if email exists
	var count int64
	s.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return nil, ErrEmailExists
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		AuthProvider: "local",
		Enabled:      true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates a user
func (s *UserService) Update(id uint, updates map[string]interface{}) (*models.User, error) {
	user, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if err := s.db.Model(user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Delete soft deletes a user
func (s *UserService) Delete(id uint) error {
	result := s.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// SetEnabled enables or disables a user
func (s *UserService) SetEnabled(id uint, enabled bool) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("enabled", enabled).Error
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.Get(id)
	if err != nil {
		return err
	}

	if !auth.CheckPassword(oldPassword, user.PasswordHash) {
		return ErrInvalidCredentials
	}

	passwordHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.db.Model(user).Update("password_hash", passwordHash).Error
}

// NodeService handles node operations
type NodeService struct {
	db *gorm.DB
}

// NewNodeService creates a new node service
func NewNodeService(db *gorm.DB) *NodeService {
	return &NodeService{db: db}
}

// List retrieves nodes with optional filtering
func (s *NodeService) List(status, region string) ([]models.Node, error) {
	var nodes []models.Node
	query := s.db.Model(&models.Node{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if region != "" {
		query = query.Where("region = ?", region)
	}

	if err := query.Find(&nodes).Error; err != nil {
		return nil, err
	}

	return nodes, nil
}

// Get retrieves a node by ID
func (s *NodeService) Get(id uint) (*models.Node, error) {
	var node models.Node
	if err := s.db.First(&node, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("node not found")
		}
		return nil, err
	}
	return &node, nil
}

// Create creates a new node
func (s *NodeService) Create(node *models.Node) error {
	return s.db.Create(node).Error
}

// Update updates a node
func (s *NodeService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.Node{}).Where("id = ?", id).Updates(updates).Error
}

// Delete soft deletes a node
func (s *NodeService) Delete(id uint) error {
	return s.db.Delete(&models.Node{}, id).Error
}

// UpdateStats updates node statistics
func (s *NodeService) UpdateStats(id uint, trafficUp, trafficDown int64, userCount int) error {
	return s.db.Model(&models.Node{}).Where("id = ?", id).Updates(map[string]interface{}{
		"traffic_up":   trafficUp,
		"traffic_down": trafficDown,
		"user_count":   userCount,
	}).Error
}

// UpdateStatus updates node status
func (s *NodeService) UpdateStatus(id uint, status string) error {
	now := time.Now()
	return s.db.Model(&models.Node{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":          status,
		"last_checked_at": now,
	}).Error
}

// GetStats returns aggregate node statistics
func (s *NodeService) GetStats() (map[string]interface{}, error) {
	var totalNodes, onlineNodes, offlineNodes int64
	var totalTrafficUp, totalTrafficDown int64
	var totalUsers int

	s.db.Model(&models.Node{}).Count(&totalNodes)
	s.db.Model(&models.Node{}).Where("status = ?", "online").Count(&onlineNodes)
	s.db.Model(&models.Node{}).Where("status = ?", "offline").Count(&offlineNodes)
	s.db.Model(&models.Node{}).Select("COALESCE(SUM(traffic_up), 0)").Scan(&totalTrafficUp)
	s.db.Model(&models.Node{}).Select("COALESCE(SUM(traffic_down), 0)").Scan(&totalTrafficDown)
	s.db.Model(&models.Node{}).Select("COALESCE(SUM(user_count), 0)").Scan(&totalUsers)

	return map[string]interface{}{
		"total_nodes":     totalNodes,
		"online_nodes":    onlineNodes,
		"offline_nodes":   offlineNodes,
		"total_traffic":   totalTrafficUp + totalTrafficDown,
		"traffic_up":      totalTrafficUp,
		"traffic_down":    totalTrafficDown,
		"total_users":     totalUsers,
	}, nil
}

// AuditService handles audit log operations
type AuditService struct {
	db *gorm.DB
}

// NewAuditService creates a new audit service
func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

// Log creates an audit log entry
func (s *AuditService) Log(userID uint, action, resource, resourceID, ip, userAgent, details, status string) error {
	log := &models.AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IP:         ip,
		UserAgent:  userAgent,
		Details:    details,
		Status:     status,
	}
	return s.db.Create(log).Error
}

// List retrieves audit logs with pagination and filtering
func (s *AuditService) List(page, perPage int, userID uint, action, resource string) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := s.db.Model(&models.AuditLog{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resource != "" {
		query = query.Where("resource = ?", resource)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// DashboardService handles dashboard operations
type DashboardService struct {
	db *gorm.DB
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

// GetStats returns dashboard statistics
func (s *DashboardService) GetStats() (map[string]interface{}, error) {
	var userCount, nodeCount, activeSubs int64
	var totalTrafficUp, totalTrafficDown int64

	s.db.Model(&models.User{}).Count(&userCount)
	s.db.Model(&models.Node{}).Count(&nodeCount)
	s.db.Model(&models.Subscription{}).Where("status = ?", "active").Count(&activeSubs)
	s.db.Model(&models.Node{}).Select("COALESCE(SUM(traffic_up), 0)").Scan(&totalTrafficUp)
	s.db.Model(&models.Node{}).Select("COALESCE(SUM(traffic_down), 0)").Scan(&totalTrafficDown)

	// Get online nodes count
	var onlineNodes int64
	s.db.Model(&models.Node{}).Where("status = ?", "online").Count(&onlineNodes)

	// Get recent audit logs
	var recentLogs []models.AuditLog
	s.db.Order("created_at DESC").Limit(10).Find(&recentLogs)

	return map[string]interface{}{
		"users":           userCount,
		"nodes":           nodeCount,
		"online_nodes":    onlineNodes,
		"active_subs":     activeSubs,
		"total_traffic":   totalTrafficUp + totalTrafficDown,
		"traffic_up":      totalTrafficUp,
		"traffic_down":    totalTrafficDown,
		"recent_activity": recentLogs,
	}, nil
}

// PlanService handles plan operations
type PlanService struct {
	db *gorm.DB
}

// NewPlanService creates a new plan service
func NewPlanService(db *gorm.DB) *PlanService {
	return &PlanService{db: db}
}

// List retrieves all plans
func (s *PlanService) List(enabledOnly bool) ([]models.Plan, error) {
	var plans []models.Plan
	query := s.db.Model(&models.Plan{})
	if enabledOnly {
		query = query.Where("enabled = ?", true)
	}
	if err := query.Order("sort_order ASC").Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

// Get retrieves a plan by ID
func (s *PlanService) Get(id uint) (*models.Plan, error) {
	var plan models.Plan
	if err := s.db.First(&plan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("plan not found")
		}
		return nil, err
	}
	return &plan, nil
}

// SubscriptionService handles subscription operations
type SubscriptionService struct {
	db *gorm.DB
}

// NewSubscriptionService creates a new subscription service
func NewSubscriptionService(db *gorm.DB) *SubscriptionService {
	return &SubscriptionService{db: db}
}

// GetByUserID retrieves subscriptions for a user
func (s *SubscriptionService) GetByUserID(userID uint) ([]models.Subscription, error) {
	var subs []models.Subscription
	if err := s.db.Where("user_id = ?", userID).Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}

// OrderService handles order operations
type OrderService struct {
	db *gorm.DB
}

// NewOrderService creates a new order service
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

// List retrieves orders with pagination
func (s *OrderService) List(page, perPage int, userID uint, status string) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := s.db.Model(&models.Order{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Get retrieves an order by ID
func (s *OrderService) Get(id uint) (*models.Order, error) {
	var order models.Order
	if err := s.db.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}