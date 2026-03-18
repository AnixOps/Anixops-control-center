export default {
  common: {
    appName: 'AnixOps Control Center',
    loading: 'Loading...',
    // ... 保持原有内容
    offlineMode: 'Offline Mode',
    onlineMode: 'Online Mode',
    syncPending: 'Sync Pending ({count})',
    lastSync: 'Last Sync: {time}',
    exportData: 'Export Data',
    importData: 'Import Data',
    downloadBackup: 'Download Backup',
    restoreBackup: 'Restore Backup'
  },

  offline: {
    title: 'Offline Mode',
    subtitle: 'Work without internet connection',
    enabled: 'Offline mode enabled',
    disabled: 'Online mode',
    pendingActions: 'Pending Actions',
    pendingCount: '{count} actions pending sync',
    syncNow: 'Sync Now',
    syncSuccess: 'Synced {success} actions successfully',
    syncFailed: 'Failed to sync {failed} actions',
    lastSyncTime: 'Last synced: {time}',
    autoSync: 'Auto-sync when online',
    saveForOffline: 'Save for Offline',
    clearCache: 'Clear Offline Cache',
    cacheSize: 'Cache Size: {size}',
    availableOffline: 'Available Offline',
    dataExpiry: 'Offline data expires after {days} days'
  },

  notifications: {
    title: 'Notifications',
    subtitle: 'Manage notification preferences',
    pushEnabled: 'Push Notifications',
    pushEnabledDesc: 'Receive push notifications on your device',
    emailEnabled: 'Email Notifications',
    emailEnabledDesc: 'Receive notifications via email',
    slackEnabled: 'Slack Notifications',
    slackEnabledDesc: 'Send notifications to Slack channel',
    webhookEnabled: 'Webhook Notifications',
    webhookEnabledDesc: 'Send notifications to custom webhook URL',

    types: 'Notification Types',
    nodeAlerts: 'Node Alerts',
    nodeAlertsDesc: 'Alerts when nodes go offline or have issues',
    userAlerts: 'User Alerts',
    userAlertsDesc: 'Alerts for user registration, bans, etc.',
    systemAlerts: 'System Alerts',
    systemAlertsDesc: 'Critical system notifications',
    trafficAlerts: 'Traffic Alerts',
    trafficAlertsDesc: 'Alerts when traffic thresholds are exceeded',
    securityAlerts: 'Security Alerts',
    securityAlertsDesc: 'Login attempts, permission changes',
    deploymentAlerts: 'Deployment Alerts',
    deploymentAlertsDesc: 'Playbook execution results',

    quietHours: 'Quiet Hours',
    quietHoursEnabled: 'Enable Quiet Hours',
    quietHoursDesc: 'No notifications during specified hours',
    startTime: 'Start Time',
    endTime: 'End Time',

    sound: 'Notification Sound',
    soundEnabled: 'Enable Sound',
    soundVolume: 'Volume',
    testSound: 'Test Sound',

    digest: 'Notification Digest',
    digestEnabled: 'Enable Daily Digest',
    digestTime: 'Digest Time',
    digestDesc: 'Receive a daily summary of all notifications'
  },

  export: {
    title: 'Data Export',
    subtitle: 'Export and import your data',
    exportNow: 'Export Now',
    importFromFile: 'Import from File',
    selectFormat: 'Select Format',
    formatJSON: 'JSON',
    formatCSV: 'CSV',
    formatYAML: 'YAML',

    includeData: 'Include Data',
    includeNodes: 'Include Nodes',
    includeUsers: 'Include Users',
    includePlugins: 'Include Plugins',
    includeSettings: 'Include Settings',
    includePlaybooks: 'Include Playbooks',
    includeSensitive: 'Include Sensitive Data',
    includeSensitiveDesc: 'Include API keys, tokens, and other sensitive information',

    importMode: 'Import Mode',
    importMerge: 'Merge',
    importMergeDesc: 'Merge with existing data',
    importReplace: 'Replace',
    importReplaceDesc: 'Replace all existing data',
    importAppend: 'Append',
    importAppendDesc: 'Add to existing data',

    lastExport: 'Last Export',
    lastImport: 'Last Import',
    exportHistory: 'Export History',
    noHistory: 'No export history',

    validating: 'Validating file...',
    importing: 'Importing data...',
    exporting: 'Exporting data...',

    importSuccess: 'Import completed successfully',
    importStats: 'Imported: {nodes} nodes, {users} users, {plugins} plugins',
    exportSuccess: 'Export completed successfully',
    downloadReady: 'Download ready',

    invalidFormat: 'Invalid file format',
    incompatibleVersion: 'Incompatible version',
    importFailed: 'Import failed: {error}'
  }
}