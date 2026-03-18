// Wails runtime bindings - generated from Go backend
// This file provides TypeScript bindings for the Go services

declare global {
  interface Window {
    go: {
      main: {
        App: {
          GetConfig: () => Promise<AppConfig>
          SetConfig: (config: AppConfig) => Promise<void>
        }
        AppServices: {
          // Auth methods
          SetAuthTokens: (token: string, refreshToken: string) => Promise<void>
          GetAuthToken: () => Promise<string>
          SetUser: (user: UserInfo) => Promise<void>
          GetUser: () => Promise<UserInfo>
          IsAuthenticated: () => Promise<boolean>
          Logout: () => Promise<void>

          // Update methods
          CheckForUpdates: () => Promise<UpdateInfo>
          DownloadUpdate: () => Promise<void>
          InstallUpdate: () => Promise<void>
          GetUpdateProgress: () => Promise<number>

          // System methods
          GetSystemTime: () => Promise<string>
          GetTimezones: () => Promise<string[]>
          GetNetworkInterfaces: () => Promise<NetworkInterface[]>
          Ping: (host: string) => Promise<PingResult>
          Traceroute: (host: string) => Promise<TracerouteHop[]>
          DNSLookup: (domain: string) => Promise<DNSResult>

          // Utility methods
          FormatBytes: (bytes: number) => Promise<string>
          GenerateUUID: () => Promise<string>
          HashString: (input: string) => Promise<string>
        }
        ConfigManager: {
          Load: () => Promise<void>
          Save: () => Promise<void>
          GetConfig: () => Promise<AppConfig>
          SetConfig: (config: AppConfig) => Promise<void>
          SetTheme: (theme: string) => Promise<void>
          SetAPIUrl: (url: string) => Promise<void>
          SetLanguage: (lang: string) => Promise<void>
        }
      }
    }
    runtime: {
      BrowserOpenURL: (url: string) => void
      EventsEmit: (eventName: string, ...data: unknown[]) => void
      EventsOn: (eventName: string, callback: (...data: unknown[]) => void) => void
      EventsOff: (eventName: string) => void
      WindowReload: () => void
      WindowSetTitle: (title: string) => void
      WindowShow: () => void
      WindowHide: () => void
      WindowMinimise: () => void
      WindowMaximise: () => void
      WindowUnmaximise: () => void
      WindowToggleMaximise: () => void
      WindowFullscreen: () => void
      WindowUnfullscreen: () => void
      WindowClose: () => void
    }
  }
}

// Type definitions
export interface AppConfig {
  api_url: string
  theme: string
  language: string
  auto_update: boolean
  shortcuts: Record<string, string>
  window: WindowConfig
}

export interface WindowConfig {
  width: number
  height: number
  maximized: boolean
  x: number
  y: number
}

export interface UserInfo {
  id: string
  email: string
  name: string
  role: string
}

export interface UpdateInfo {
  update_available: boolean
  current_version: string
  latest_version: string
  release_notes: string
}

export interface NetworkInterface {
  name: string
  ip: string
  mac: string
  is_up: boolean
  is_loopback: boolean
}

export interface PingResult {
  host: string
  success: boolean
  latency: number
  time: string
}

export interface TracerouteHop {
  hop: number
  host: string
  latency: number
}

export interface DNSResult {
  domain: string
  records: DNSRecord[]
}

export interface DNSRecord {
  type: string
  value: string
}

// API wrapper class
export class WailsAPI {
  private static instance: WailsAPI

  static getInstance(): WailsAPI {
    if (!WailsAPI.instance) {
      WailsAPI.instance = new WailsAPI()
    }
    return WailsAPI.instance
  }

  // Auth methods
  async setAuthTokens(token: string, refreshToken: string): Promise<void> {
    return window.go.main.AppServices.SetAuthTokens(token, refreshToken)
  }

  async getAuthToken(): Promise<string> {
    return window.go.main.AppServices.GetAuthToken()
  }

  async setUser(user: UserInfo): Promise<void> {
    return window.go.main.AppServices.SetUser(user)
  }

  async getUser(): Promise<UserInfo | null> {
    return window.go.main.AppServices.GetUser()
  }

  async isAuthenticated(): Promise<boolean> {
    return window.go.main.AppServices.IsAuthenticated()
  }

  async logout(): Promise<void> {
    return window.go.main.AppServices.Logout()
  }

  // Config methods
  async getConfig(): Promise<AppConfig> {
    return window.go.main.ConfigManager.GetConfig()
  }

  async setConfig(config: AppConfig): Promise<void> {
    return window.go.main.ConfigManager.SetConfig(config)
  }

  async setTheme(theme: string): Promise<void> {
    return window.go.main.ConfigManager.SetTheme(theme)
  }

  async setAPIUrl(url: string): Promise<void> {
    return window.go.main.ConfigManager.SetAPIUrl(url)
  }

  async setLanguage(lang: string): Promise<void> {
    return window.go.main.ConfigManager.SetLanguage(lang)
  }

  // Update methods
  async checkForUpdates(): Promise<UpdateInfo> {
    return window.go.main.AppServices.CheckForUpdates()
  }

  async downloadUpdate(): Promise<void> {
    return window.go.main.AppServices.DownloadUpdate()
  }

  async installUpdate(): Promise<void> {
    return window.go.main.AppServices.InstallUpdate()
  }

  async getUpdateProgress(): Promise<number> {
    return window.go.main.AppServices.GetUpdateProgress()
  }

  // System methods
  async getSystemTime(): Promise<string> {
    return window.go.main.AppServices.GetSystemTime()
  }

  async getTimezones(): Promise<string[]> {
    return window.go.main.AppServices.GetTimezones()
  }

  async getNetworkInterfaces(): Promise<NetworkInterface[]> {
    return window.go.main.AppServices.GetNetworkInterfaces()
  }

  async ping(host: string): Promise<PingResult> {
    return window.go.main.AppServices.Ping(host)
  }

  async traceroute(host: string): Promise<TracerouteHop[]> {
    return window.go.main.AppServices.Traceroute(host)
  }

  async dnsLookup(domain: string): Promise<DNSResult> {
    return window.go.main.AppServices.DNSLookup(domain)
  }

  // Utility methods
  async formatBytes(bytes: number): Promise<string> {
    return window.go.main.AppServices.FormatBytes(bytes)
  }

  async generateUUID(): Promise<string> {
    return window.go.main.AppServices.GenerateUUID()
  }

  async hashString(input: string): Promise<string> {
    return window.go.main.AppServices.HashString(input)
  }

  // Window methods
  openURL(url: string): void {
    window.runtime.BrowserOpenURL(url)
  }

  reload(): void {
    window.runtime.WindowReload()
  }

  setTitle(title: string): void {
    window.runtime.WindowSetTitle(title)
  }

  minimize(): void {
    window.runtime.WindowMinimise()
  }

  maximize(): void {
    window.runtime.WindowMaximise()
  }

  toggleMaximize(): void {
    window.runtime.WindowToggleMaximise()
  }

  close(): void {
    window.runtime.WindowClose()
  }

  // Events
  emit(eventName: string, ...data: unknown[]): void {
    window.runtime.EventsEmit(eventName, ...data)
  }

  on(eventName: string, callback: (...data: unknown[]) => void): void {
    window.runtime.EventsOn(eventName, callback)
  }

  off(eventName: string): void {
    window.runtime.EventsOff(eventName)
  }
}

export const wailsAPI = WailsAPI.getInstance()