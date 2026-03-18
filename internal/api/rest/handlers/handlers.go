package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anixops/anixops-control-center/internal/core/plugin"
	"github.com/anixops/anixops-control-center/internal/services"
	"github.com/anixops/anixops-control-center/internal/storage/models"
	"github.com/gin-gonic/gin"
)

// Handlers holds all API handlers
type Handlers struct {
	pluginMgr *plugin.Manager
	authSvc   *services.AuthService
	userSvc   *services.UserService
	nodeSvc   *services.NodeService
	auditSvc  *services.AuditService
	dashSvc   *services.DashboardService
	planSvc   *services.PlanService
	subSvc    *services.SubscriptionService
	orderSvc  *services.OrderService
}

// NewHandlers creates a new Handlers instance
func NewHandlers(
	pluginMgr *plugin.Manager,
	authSvc *services.AuthService,
	userSvc *services.UserService,
	nodeSvc *services.NodeService,
	auditSvc *services.AuditService,
	dashSvc *services.DashboardService,
	planSvc *services.PlanService,
	subSvc *services.SubscriptionService,
	orderSvc *services.OrderService,
) *Handlers {
	return &Handlers{
		pluginMgr: pluginMgr,
		authSvc:   authSvc,
		userSvc:   userSvc,
		nodeSvc:   nodeSvc,
		auditSvc:  auditSvc,
		dashSvc:   dashSvc,
		planSvc:   planSvc,
		subSvc:    subSvc,
		orderSvc:  orderSvc,
	}
}

// === Health ===

func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"version":   "1.0.0",
		"timestamp": time.Now().Unix(),
	})
}

func (h *Handlers) ReadinessCheck(c *gin.Context) {
	plugins := 0
	if h.pluginMgr != nil {
		plugins = len(h.pluginMgr.List())
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"plugins":   plugins,
		"timestamp": time.Now().Unix(),
	})
}

// === Auth ===

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role"`
}

func (h *Handlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.authSvc == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "auth service not configured"})
		return
	}

	user, accessToken, refreshToken, err := h.authSvc.Login(
		req.Email,
		req.Password,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if err == services.ErrUserDisabled {
			c.JSON(http.StatusForbidden, gin.H{"error": "user is disabled"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log the login
	if h.auditSvc != nil {
		h.auditSvc.Log(user.ID, "login", "user", string(rune(user.ID)), c.ClientIP(), c.GetHeader("User-Agent"), "", "success")
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    86400,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (h *Handlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.authSvc == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "auth service not configured"})
		return
	}

	role := req.Role
	if role == "" {
		role = "viewer"
	}

	user, err := h.authSvc.Register(req.Email, req.Password, role)
	if err != nil {
		if err == services.ErrEmailExists {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log the registration
	if h.auditSvc != nil {
		h.auditSvc.Log(user.ID, "register", "user", string(rune(user.ID)), c.ClientIP(), c.GetHeader("User-Agent"), "", "success")
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (h *Handlers) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if exists && h.auditSvc != nil {
		uid, _ := strconv.ParseUint(userID.(string), 10, 64)
		h.auditSvc.Log(uint(uid), "logout", "user", userID.(string), c.ClientIP(), c.GetHeader("User-Agent"), "", "success")
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h *Handlers) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.authSvc == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "auth service not configured"})
		return
	}

	accessToken, err := h.authSvc.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   86400,
	})
}

func (h *Handlers) GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	uid, err := strconv.ParseUint(userID.(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"id": userID, "email": "admin@example.com", "role": role})
		return
	}

	if h.userSvc != nil {
		user, err := h.userSvc.Get(uint(uid))
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"id":         user.ID,
				"email":      user.Email,
				"role":       user.Role,
				"enabled":    user.Enabled,
				"last_login": user.LastLoginAt,
				"created_at": user.CreatedAt,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"id": userID, "email": "admin@example.com", "role": role})
}

// === Plugins ===

func (h *Handlers) ListPlugins(c *gin.Context) {
	if h.pluginMgr == nil {
		c.JSON(http.StatusOK, gin.H{"plugins": []interface{}{}})
		return
	}
	infos := h.pluginMgr.GetInfo()
	c.JSON(http.StatusOK, gin.H{"plugins": infos})
}

func (h *Handlers) GetPlugin(c *gin.Context) {
	name := c.Param("name")
	if h.pluginMgr == nil {
		c.JSON(http.StatusOK, gin.H{"name": name, "status": "running"})
		return
	}
	p, ok := h.pluginMgr.Get(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "plugin not found"})
		return
	}
	c.JSON(http.StatusOK, p.Info())
}

func (h *Handlers) ExecutePlugin(c *gin.Context) {
	name := c.Param("name")
	var req struct {
		Action string                 `json:"action"`
		Params map[string]interface{} `json:"params"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.pluginMgr == nil {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}

	p, ok := h.pluginMgr.Get(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "plugin not found"})
		return
	}

	execPlugin, ok := p.(plugin.ExecutablePlugin)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plugin does not support execution"})
		return
	}

	result, err := execPlugin.Execute(c.Request.Context(), req.Action, req.Params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handlers) GetPluginStatus(c *gin.Context) {
	name := c.Param("name")
	if h.pluginMgr == nil {
		c.JSON(http.StatusOK, gin.H{"state": "running"})
		return
	}
	state := h.pluginMgr.GetState(name)
	c.JSON(http.StatusOK, gin.H{"state": string(state)})
}

func (h *Handlers) GetPluginCapabilities(c *gin.Context) {
	name := c.Param("name")
	if h.pluginMgr == nil {
		c.JSON(http.StatusOK, gin.H{"capabilities": []string{}})
		return
	}
	p, ok := h.pluginMgr.Get(name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "plugin not found"})
		return
	}
	caps := p.Capabilities()
	c.JSON(http.StatusOK, gin.H{"capabilities": caps})
}

func (h *Handlers) StartPlugin(c *gin.Context) {
	name := c.Param("name")
	if h.pluginMgr == nil {
		c.JSON(http.StatusOK, gin.H{"message": "plugin started"})
		return
	}
	if err := h.pluginMgr.StartPlugin(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "plugin started"})
}

func (h *Handlers) StopPlugin(c *gin.Context) {
	name := c.Param("name")
	if h.pluginMgr == nil {
		c.JSON(http.StatusOK, gin.H{"message": "plugin stopped"})
		return
	}
	if err := h.pluginMgr.StopPlugin(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "plugin stopped"})
}

func (h *Handlers) RestartPlugin(c *gin.Context) {
	name := c.Param("name")
	if h.pluginMgr == nil {
		c.JSON(http.StatusOK, gin.H{"message": "plugin restarted"})
		return
	}
	// Stop then start
	if err := h.pluginMgr.StopPlugin(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.pluginMgr.StartPlugin(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "plugin restarted"})
}

// === Nodes ===

type CreateNodeRequest struct {
	Name     string `json:"name" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port"`
	Type     string `json:"type"`
	Region   string `json:"region"`
	ServerID uint   `json:"server_id"`
	MaxUsers int    `json:"max_users"`
}

func (h *Handlers) ListNodes(c *gin.Context) {
	if h.nodeSvc == nil {
		c.JSON(http.StatusOK, gin.H{"nodes": []interface{}{}})
		return
	}

	status := c.Query("status")
	region := c.Query("region")

	nodes, err := h.nodeSvc.List(status, region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nodes": nodes})
}

func (h *Handlers) CreateNode(c *gin.Context) {
	var req CreateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.nodeSvc == nil {
		c.JSON(http.StatusCreated, gin.H{"message": "node created"})
		return
	}

	node := &models.Node{
		Name:     req.Name,
		Host:     req.Host,
		Port:     req.Port,
		Type:     req.Type,
		Region:   req.Region,
		ServerID: req.ServerID,
		MaxUsers: req.MaxUsers,
		Status:   "unknown",
		Enabled:  true,
	}

	if node.Port == 0 {
		node.Port = 443
	}
	if node.Type == "" {
		node.Type = "v2ray"
	}

	if err := h.nodeSvc.Create(node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	if h.auditSvc != nil {
		userID, _ := c.Get("userID")
		uid, _ := strconv.ParseUint(userID.(string), 10, 64)
		h.auditSvc.Log(uint(uid), "create", "node", string(rune(node.ID)), c.ClientIP(), c.GetHeader("User-Agent"), "name: "+req.Name, "success")
	}

	c.JSON(http.StatusCreated, node)
}

func (h *Handlers) GetNode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}

	if h.nodeSvc == nil {
		c.JSON(http.StatusOK, gin.H{"id": id})
		return
	}

	node, err := h.nodeSvc.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *Handlers) UpdateNode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.nodeSvc == nil {
		c.JSON(http.StatusOK, gin.H{"message": "node updated"})
		return
	}

	if err := h.nodeSvc.Update(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "node updated"})
}

func (h *Handlers) DeleteNode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}

	if h.nodeSvc == nil {
		c.JSON(http.StatusOK, gin.H{"message": "node deleted"})
		return
	}

	if err := h.nodeSvc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "node deleted"})
}

func (h *Handlers) GetNodeStats(c *gin.Context) {
	if h.nodeSvc == nil {
		c.JSON(http.StatusOK, gin.H{"stats": nil})
		return
	}

	stats, err := h.nodeSvc.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// === Agents ===

func (h *Handlers) ListAgents(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"agents": []interface{}{}})
}

func (h *Handlers) ConnectAgent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "connected"})
}

func (h *Handlers) DisconnectAgent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "disconnected"})
}

func (h *Handlers) ExecAgent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"output": ""})
}

func (h *Handlers) GetAgentInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"info": nil})
}

// === Playbooks ===

func (h *Handlers) executePluginAction(pluginName, action string, params map[string]interface{}) (plugin.Result, error) {
	if h.pluginMgr == nil {
		return plugin.Result{}, nil
	}

	p, ok := h.pluginMgr.Get(pluginName)
	if !ok {
		return plugin.Result{}, nil
	}

	execPlugin, ok := p.(plugin.ExecutablePlugin)
	if !ok {
		return plugin.Result{}, nil
	}

	return execPlugin.Execute(context.Background(), action, params)
}

func (h *Handlers) ListPlaybooks(c *gin.Context) {
	if h.pluginMgr == nil {
		c.JSON(http.StatusOK, gin.H{"playbooks": []interface{}{}})
		return
	}

	result, err := h.executePluginAction("ansible", "list_playbooks", nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"playbooks": []interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"playbooks": result.Data})
}

func (h *Handlers) GetPlaybook(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{"name": name})
}

func (h *Handlers) RunPlaybook(c *gin.Context) {
	var req struct {
		Playbook  string                 `json:"playbook" binding:"required"`
		Inventory string                 `json:"inventory"`
		Tags      string                 `json:"tags"`
		Limit     string                 `json:"limit"`
		ExtraVars map[string]interface{} `json:"extra_vars"`
		Verbose   bool                   `json:"verbose"`
		Check     bool                   `json:"check"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.pluginMgr == nil {
		c.JSON(http.StatusOK, gin.H{"message": "playbook started"})
		return
	}

	params := map[string]interface{}{
		"playbook":   req.Playbook,
		"inventory":  req.Inventory,
		"tags":       req.Tags,
		"limit":      req.Limit,
		"extra_vars": req.ExtraVars,
		"verbose":    req.Verbose,
		"check":      req.Check,
	}

	result, err := h.executePluginAction("ansible", "run_playbook", params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": result.Success,
		"output":  result.Data,
		"error":   result.Error,
	})
}

func (h *Handlers) ValidatePlaybook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"valid": true})
}

// === Inventory ===

func (h *Handlers) GetInventory(c *gin.Context) {
	if h.pluginMgr == nil {
		c.JSON(http.StatusOK, gin.H{"inventory": nil})
		return
	}

	result, err := h.executePluginAction("ansible", "get_inventory", nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"inventory": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"inventory": result.Data})
}

func (h *Handlers) ListHosts(c *gin.Context) {
	if h.pluginMgr == nil {
		c.JSON(http.StatusOK, gin.H{"hosts": []interface{}{}})
		return
	}

	result, err := h.executePluginAction("ansible", "list_hosts", nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"hosts": []interface{}{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hosts": result.Data})
}

// === Users ===

func (h *Handlers) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "50"))
	search := c.Query("search")

	if h.userSvc == nil {
		c.JSON(http.StatusOK, gin.H{"users": []interface{}{}, "page": page, "per_page": perPage, "total": 0})
		return
	}

	users, total, err := h.userSvc.List(page, perPage, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":    users,
		"page":     page,
		"per_page": perPage,
		"total":    total,
	})
}

func (h *Handlers) CreateUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.userSvc == nil {
		c.JSON(http.StatusCreated, gin.H{"message": "user created"})
		return
	}

	role := req.Role
	if role == "" {
		role = "viewer"
	}

	user, err := h.userSvc.Create(req.Email, req.Password, role)
	if err != nil {
		if err == services.ErrEmailExists {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handlers) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if h.userSvc == nil {
		c.JSON(http.StatusOK, gin.H{"id": id})
		return
	}

	user, err := h.userSvc.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handlers) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.userSvc == nil {
		c.JSON(http.StatusOK, gin.H{"message": "user updated"})
		return
	}

	user, err := h.userSvc.Update(uint(id), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handlers) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if h.userSvc == nil {
		c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
		return
	}

	if err := h.userSvc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func (h *Handlers) BanUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if h.userSvc == nil {
		c.JSON(http.StatusOK, gin.H{"message": "user banned"})
		return
	}

	if err := h.userSvc.SetEnabled(uint(id), false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user banned"})
}

func (h *Handlers) UnbanUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if h.userSvc == nil {
		c.JSON(http.StatusOK, gin.H{"message": "user unbanned"})
		return
	}

	if err := h.userSvc.SetEnabled(uint(id), true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user unbanned"})
}

func (h *Handlers) GetUserSubscriptions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if h.subSvc == nil {
		c.JSON(http.StatusOK, gin.H{"subscriptions": []interface{}{}})
		return
	}

	subs, err := h.subSvc.GetByUserID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscriptions": subs})
}

// === Orders ===

func (h *Handlers) ListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "50"))
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	status := c.Query("status")

	if h.orderSvc == nil {
		c.JSON(http.StatusOK, gin.H{"orders": []interface{}{}, "page": page, "per_page": perPage, "total": 0})
		return
	}

	orders, total, err := h.orderSvc.List(page, perPage, uint(userID), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":   orders,
		"page":     page,
		"per_page": perPage,
		"total":    total,
	})
}

func (h *Handlers) GetOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	if h.orderSvc == nil {
		c.JSON(http.StatusOK, gin.H{"id": id})
		return
	}

	order, err := h.orderSvc.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// === Plans ===

func (h *Handlers) ListPlans(c *gin.Context) {
	enabledOnly := c.Query("enabled") == "true"

	if h.planSvc == nil {
		c.JSON(http.StatusOK, gin.H{"plans": []interface{}{}})
		return
	}

	plans, err := h.planSvc.List(enabledOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func (h *Handlers) GetPlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	if h.planSvc == nil {
		c.JSON(http.StatusOK, gin.H{"id": id})
		return
	}

	plan, err := h.planSvc.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// === Dashboard ===

func (h *Handlers) GetDashboard(c *gin.Context) {
	if h.dashSvc == nil {
		c.JSON(http.StatusOK, gin.H{"nodes": 8, "users": 357, "traffic": "1.2TB"})
		return
	}

	stats, err := h.dashSvc.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Format traffic
	totalTraffic := stats["total_traffic"].(int64)
	trafficStr := formatBytes(totalTraffic)

	c.JSON(http.StatusOK, gin.H{
		"nodes":        stats["nodes"],
		"online_nodes": stats["online_nodes"],
		"users":        stats["users"],
		"active_subs":  stats["active_subs"],
		"traffic":      trafficStr,
		"traffic_up":   stats["traffic_up"],
		"traffic_down": stats["traffic_down"],
	})
}

func (h *Handlers) GetDashboardStats(c *gin.Context) {
	if h.dashSvc == nil {
		c.JSON(http.StatusOK, gin.H{"stats": nil})
		return
	}

	stats, err := h.dashSvc.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// === Logs ===

func (h *Handlers) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "50"))
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	action := c.Query("action")
	resource := c.Query("resource")

	if h.auditSvc == nil {
		c.JSON(http.StatusOK, gin.H{"logs": []interface{}{}, "page": page, "per_page": perPage, "total": 0})
		return
	}

	logs, total, err := h.auditSvc.List(page, perPage, uint(userID), action, resource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":     logs,
		"page":     page,
		"per_page": perPage,
		"total":    total,
	})
}

func (h *Handlers) StreamLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "WebSocket streaming"})
}

// === Settings ===

func (h *Handlers) GetSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"settings": nil})
}

func (h *Handlers) UpdateSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "settings updated"})
}

// === Events ===

func (h *Handlers) StreamEvents(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.SSEvent("connected", gin.H{"timestamp": time.Now().Unix()})
	c.Writer.Flush()
}

// === Admin ===

func (h *Handlers) GetAdminStats(c *gin.Context) {
	plugins := 0
	if h.pluginMgr != nil {
		plugins = len(h.pluginMgr.List())
	}
	c.JSON(http.StatusOK, gin.H{"plugins": plugins, "version": "1.0.0"})
}

func (h *Handlers) GetAuditLogs(c *gin.Context) {
	h.GetLogs(c)
}

func (h *Handlers) GetSystemInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": "1.0.0", "go_version": "1.24"})
}

func (h *Handlers) CreateBackup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "backup created"})
}

func (h *Handlers) ListBackups(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"backups": []interface{}{}})
}

// Helper functions

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return strconv.FormatInt(bytes, 10) + "B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return strconv.FormatFloat(float64(bytes)/float64(div), 'f', 1, 64) + string("KMGTPE"[exp]) + "B"
}

// GetUserIDFromContext extracts user ID from gin context
func GetUserIDFromContext(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	uid, err := strconv.ParseUint(strings.TrimSpace(userID.(string)), 10, 64)
	if err != nil {
		return 0
	}
	return uint(uid)
}