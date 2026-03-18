/**
 * Constants used throughout the application
 */

// API endpoints
export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    ME: '/auth/me',
    REGISTER: '/auth/register',
    FORGOT_PASSWORD: '/auth/forgot-password',
    RESET_PASSWORD: '/auth/reset-password'
  },
  NODES: {
    LIST: '/nodes',
    CREATE: '/nodes',
    GET: (id) => `/nodes/${id}`,
    UPDATE: (id) => `/nodes/${id}`,
    DELETE: (id) => `/nodes/${id}`,
    START: (id) => `/nodes/${id}/start`,
    STOP: (id) => `/nodes/${id}/stop`,
    RESTART: (id) => `/nodes/${id}/restart`,
    STATS: (id) => `/nodes/${id}/stats`,
    LOGS: (id) => `/nodes/${id}/logs`
  },
  USERS: {
    LIST: '/users',
    CREATE: '/users',
    GET: (id) => `/users/${id}`,
    UPDATE: (id) => `/users/${id}`,
    DELETE: (id) => `/users/${id}`,
    BAN: (id) => `/users/${id}/ban`,
    UNBAN: (id) => `/users/${id}/unban`
  },
  PLUGINS: {
    LIST: '/plugins',
    GET: (name) => `/plugins/${name}`,
    ENABLE: (name) => `/plugins/${name}/enable`,
    DISABLE: (name) => `/plugins/${name}/disable`,
    RESTART: (name) => `/plugins/${name}/restart`
  },
  PLAYBOOKS: {
    LIST: '/playbooks',
    CREATE: '/playbooks',
    GET: (id) => `/playbooks/${id}`,
    RUN: (id) => `/playbooks/${id}/run`,
    STOP: (id) => `/playbooks/${id}/stop`
  },
  LOGS: {
    LIST: '/logs',
    STREAM: '/logs/stream'
  },
  SETTINGS: {
    GET: '/settings',
    UPDATE: '/settings'
  }
}

// Status colors
export const STATUS_COLORS = {
  ONLINE: 'green',
  OFFLINE: 'red',
  STARTING: 'yellow',
  STOPPING: 'yellow',
  ERROR: 'red',
  ACTIVE: 'green',
  INACTIVE: 'gray',
  BANNED: 'red',
  SUSPENDED: 'yellow',
  PENDING: 'blue',
  RUNNING: 'yellow',
  SUCCESS: 'green',
  FAILED: 'red'
}

// Log levels
export const LOG_LEVELS = {
  DEBUG: 'DEBUG',
  INFO: 'INFO',
  WARN: 'WARN',
  ERROR: 'ERROR',
  FATAL: 'FATAL'
}

// Log level colors
export const LOG_LEVEL_COLORS = {
  DEBUG: 'gray',
  INFO: 'blue',
  WARN: 'yellow',
  ERROR: 'red',
  FATAL: 'red'
}

// User roles
export const USER_ROLES = {
  ADMIN: 'admin',
  SUPERADMIN: 'superadmin',
  USER: 'user'
}

// Plan types
export const PLAN_TYPES = {
  FREE: 'free',
  BASIC: 'basic',
  PRO: 'pro',
  ENTERPRISE: 'enterprise'
}

// Node types
export const NODE_TYPES = {
  V2RAY: 'v2ray',
  XRAY: 'xray',
  TROJAN: 'trojan',
  SHADOWSOCKS: 'shadowsocks'
}

// Traffic units
export const TRAFFIC_UNITS = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']

// Time formats
export const TIME_FORMATS = {
  SHORT: 'HH:mm',
  MEDIUM: 'HH:mm:ss',
  LONG: 'YYYY-MM-DD HH:mm:ss',
  DATE: 'YYYY-MM-DD',
  DATE_SHORT: 'MM/DD',
  RELATIVE: 'relative'
}

// Pagination defaults
export const PAGINATION = {
  DEFAULT_PAGE: 1,
  DEFAULT_PAGE_SIZE: 20,
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100]
}

// WebSocket events
export const WS_EVENTS = {
  NODE_STATUS: 'node:status',
  NODE_STATS: 'node:stats',
  USER_UPDATE: 'user:update',
  LOG: 'log',
  ALERT: 'alert',
  NOTIFICATION: 'notification',
  PING: 'ping',
  PONG: 'pong'
}

// Error codes
export const ERROR_CODES = {
  UNKNOWN: 'E0000',
  VALIDATION_ERROR: 'E0001',
  NOT_FOUND: 'E0002',
  UNAUTHORIZED: 'E0003',
  FORBIDDEN: 'E0004',
  INTERNAL_ERROR: 'E0005',
  PLUGIN_ERROR: 'E0100',
  AUTH_ERROR: 'E0200',
  DATABASE_ERROR: 'E0300',
  NODE_ERROR: 'E0400',
  TASK_ERROR: 'E0500',
  CONFIG_ERROR: 'E0600',
  TRACE_ERROR: 'E0700',
  RATELIMIT_ERROR: 'E0800'
}

// HTTP status codes
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  CONFLICT: 409,
  UNPROCESSABLE_ENTITY: 422,
  TOO_MANY_REQUESTS: 429,
  INTERNAL_SERVER_ERROR: 500,
  BAD_GATEWAY: 502,
  SERVICE_UNAVAILABLE: 503
}

// Local storage keys
export const STORAGE_KEYS = {
  TOKEN: 'auth_token',
  REFRESH_TOKEN: 'refresh_token',
  USER: 'user_data',
  THEME: 'theme',
  LANGUAGE: 'language',
  SIDEBAR_COLLAPSED: 'sidebar_collapsed',
  API_URL: 'api_url'
}

// Theme options
export const THEMES = {
  LIGHT: 'light',
  DARK: 'dark',
  SYSTEM: 'system'
}

// Languages
export const LANGUAGES = {
  EN: 'en',
  ZH: 'zh'
}

// Sort directions
export const SORT_DIRECTION = {
  ASC: 'asc',
  DESC: 'desc'
}

// Filter operators
export const FILTER_OPERATORS = {
  EQ: 'eq',
  NE: 'ne',
  GT: 'gt',
  GTE: 'gte',
  LT: 'lt',
  LTE: 'lte',
  LIKE: 'like',
  IN: 'in',
  NOT_IN: 'not_in',
  IS_NULL: 'is_null',
  IS_NOT_NULL: 'is_not_null'
}