import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'privacy_policy_page.dart';
import 'terms_of_service_page.dart';

// Settings State
class SettingsState {
  final bool twoFactorEnabled;
  final bool autoBackup;
  final String theme;
  final String language;
  final bool notificationsEnabled;
  final int sessionTimeout;

  const SettingsState({
    this.twoFactorEnabled = false,
    this.autoBackup = true,
    this.theme = 'dark',
    this.language = 'en',
    this.notificationsEnabled = true,
    this.sessionTimeout = 30,
  });

  SettingsState copyWith({
    bool? twoFactorEnabled,
    bool? autoBackup,
    String? theme,
    String? language,
    bool? notificationsEnabled,
    int? sessionTimeout,
  }) {
    return SettingsState(
      twoFactorEnabled: twoFactorEnabled ?? this.twoFactorEnabled,
      autoBackup: autoBackup ?? this.autoBackup,
      theme: theme ?? this.theme,
      language: language ?? this.language,
      notificationsEnabled: notificationsEnabled ?? this.notificationsEnabled,
      sessionTimeout: sessionTimeout ?? this.sessionTimeout,
    );
  }
}

class SettingsNotifier extends StateNotifier<SettingsState> {
  SettingsNotifier() : super(const SettingsState());

  void setTwoFactor(bool value) => state = state.copyWith(twoFactorEnabled: value);
  void setAutoBackup(bool value) => state = state.copyWith(autoBackup: value);
  void setTheme(String value) => state = state.copyWith(theme: value);
  void setLanguage(String value) => state = state.copyWith(language: value);
  void setNotifications(bool value) => state = state.copyWith(notificationsEnabled: value);
  void setSessionTimeout(int value) => state = state.copyWith(sessionTimeout: value);
}

final settingsProvider = StateNotifierProvider<SettingsNotifier, SettingsState>((ref) {
  return SettingsNotifier();
});

// Settings Page
class SettingsPage extends ConsumerWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(settingsProvider);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Settings'),
      ),
      body: ListView(
        children: [
          // Account Section
          _buildSectionHeader('Account', context),
          ListTile(
            leading: CircleAvatar(
              backgroundColor: theme.colorScheme.primary,
              child: const Icon(Icons.person, color: Colors.white),
            ),
            title: const Text('Admin User'),
            subtitle: const Text('admin@example.com'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              // Navigate to profile
            },
          ),
          ListTile(
            leading: const Icon(Icons.lock_outline),
            title: const Text('Change Password'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              // Navigate to change password
            },
          ),
          SwitchListTile(
            secondary: const Icon(Icons.security),
            title: const Text('Two-Factor Authentication'),
            subtitle: const Text('Add extra security to your account'),
            value: state.twoFactorEnabled,
            onChanged: (value) => ref.read(settingsProvider.notifier).setTwoFactor(value),
          ),

          const Divider(),

          // Server Section
          _buildSectionHeader('Server Configuration', context),
          ListTile(
            leading: const Icon(Icons.dns_outlined),
            title: const Text('Server Settings'),
            subtitle: const Text('Host, port, SSL configuration'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              // Navigate to server settings
            },
          ),
          ListTile(
            leading: const Icon(Icons.storage_outlined),
            title: const Text('Database'),
            subtitle: const Text('SQLite'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              // Navigate to database settings
            },
          ),

          const Divider(),

          // Plugins Section
          _buildSectionHeader('Plugins', context),
          ListTile(
            leading: const Icon(Icons.extension_outlined),
            title: const Text('Manage Plugins'),
            subtitle: const Text('4 plugins installed'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              // Navigate to plugins
            },
          ),

          const Divider(),

          // Security Section
          _buildSectionHeader('Security', context),
          ListTile(
            leading: const Icon(Icons.admin_panel_settings_outlined),
            title: const Text('IP Whitelist'),
            subtitle: const Text('Restrict access by IP'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              // Navigate to IP whitelist
            },
          ),
          ListTile(
            leading: const Icon(Icons.timer_outlined),
            title: const Text('Session Timeout'),
            subtitle: Text('${state.sessionTimeout} minutes'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              _showTimeoutDialog(context, ref, state.sessionTimeout);
            },
          ),

          const Divider(),

          // Notifications Section
          _buildSectionHeader('Notifications', context),
          SwitchListTile(
            secondary: const Icon(Icons.notifications_outlined),
            title: const Text('Push Notifications'),
            subtitle: const Text('Receive alerts on your device'),
            value: state.notificationsEnabled,
            onChanged: (value) => ref.read(settingsProvider.notifier).setNotifications(value),
          ),
          ListTile(
            leading: const Icon(Icons.email_outlined),
            title: const Text('Email Settings'),
            subtitle: const Text('Configure email notifications'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              // Navigate to email settings
            },
          ),

          const Divider(),

          // Backup Section
          _buildSectionHeader('Backup & Restore', context),
          SwitchListTile(
            secondary: const Icon(Icons.backup_outlined),
            title: const Text('Automatic Backup'),
            subtitle: const Text('Daily at midnight'),
            value: state.autoBackup,
            onChanged: (value) => ref.read(settingsProvider.notifier).setAutoBackup(value),
          ),
          ListTile(
            leading: const Icon(Icons.cloud_download_outlined),
            title: const Text('Create Backup'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              _showBackupDialog(context);
            },
          ),
          ListTile(
            leading: const Icon(Icons.restore_outlined),
            title: const Text('Restore from Backup'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              // Show restore dialog
            },
          ),

          const Divider(),

          // Appearance Section
          _buildSectionHeader('Appearance', context),
          ListTile(
            leading: const Icon(Icons.palette_outlined),
            title: const Text('Theme'),
            subtitle: Text(state.theme == 'dark' ? 'Dark' : 'Light'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              _showThemeDialog(context, ref, state.theme);
            },
          ),
          ListTile(
            leading: const Icon(Icons.language),
            title: const Text('Language'),
            subtitle: Text(state.language == 'en' ? 'English' : '中文'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              _showLanguageDialog(context, ref, state.language);
            },
          ),

          const Divider(),

          // About Section
          _buildSectionHeader('About', context),
          ListTile(
            leading: const Icon(Icons.info_outline),
            title: const Text('Version'),
            subtitle: const Text('0.9.9'),
          ),
          ListTile(
            leading: const Icon(Icons.update),
            title: const Text('Check for Updates'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('You are running the latest version')),
              );
            },
          ),
          ListTile(
            leading: const Icon(Icons.description_outlined),
            title: const Text('Documentation'),
            trailing: const Icon(Icons.open_in_new),
            onTap: () {
              // Open documentation
            },
          ),
          ListTile(
            leading: const Icon(Icons.bug_report_outlined),
            title: const Text('Report an Issue'),
            trailing: const Icon(Icons.open_in_new),
            onTap: () {
              // Open issue tracker
            },
          ),
          ListTile(
            leading: const Icon(Icons.privacy_tip_outlined),
            title: const Text('Privacy Policy'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => const PrivacyPolicyPage()),
              );
            },
          ),
          ListTile(
            leading: const Icon(Icons.gavel_outlined),
            title: const Text('Terms of Service'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => const TermsOfServicePage()),
              );
            },
          ),

          const SizedBox(height: 24),

          // Logout Button
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: OutlinedButton.icon(
              onPressed: () {
                _showLogoutDialog(context);
              },
              icon: const Icon(Icons.logout),
              label: const Text('Sign Out'),
              style: OutlinedButton.styleFrom(
                foregroundColor: Colors.red,
                side: const BorderSide(color: Colors.red),
                padding: const EdgeInsets.symmetric(vertical: 12),
              ),
            ),
          ),

          const SizedBox(height: 32),
        ],
      ),
    );
  }

  Widget _buildSectionHeader(String title, BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      child: Text(
        title,
        style: Theme.of(context).textTheme.titleSmall?.copyWith(
          color: Theme.of(context).colorScheme.primary,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }

  void _showTimeoutDialog(BuildContext context, WidgetRef ref, int currentValue) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Session Timeout'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            RadioListTile<int>(
              title: const Text('15 minutes'),
              value: 15,
              groupValue: currentValue,
              onChanged: (value) {
                ref.read(settingsProvider.notifier).setSessionTimeout(value!);
                Navigator.pop(context);
              },
            ),
            RadioListTile<int>(
              title: const Text('30 minutes'),
              value: 30,
              groupValue: currentValue,
              onChanged: (value) {
                ref.read(settingsProvider.notifier).setSessionTimeout(value!);
                Navigator.pop(context);
              },
            ),
            RadioListTile<int>(
              title: const Text('1 hour'),
              value: 60,
              groupValue: currentValue,
              onChanged: (value) {
                ref.read(settingsProvider.notifier).setSessionTimeout(value!);
                Navigator.pop(context);
              },
            ),
            RadioListTile<int>(
              title: const Text('Never'),
              value: 0,
              groupValue: currentValue,
              onChanged: (value) {
                ref.read(settingsProvider.notifier).setSessionTimeout(value!);
                Navigator.pop(context);
              },
            ),
          ],
        ),
      ),
    );
  }

  void _showThemeDialog(BuildContext context, WidgetRef ref, String currentValue) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Theme'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            RadioListTile<String>(
              title: const Text('Light'),
              value: 'light',
              groupValue: currentValue,
              onChanged: (value) {
                ref.read(settingsProvider.notifier).setTheme(value!);
                Navigator.pop(context);
              },
            ),
            RadioListTile<String>(
              title: const Text('Dark'),
              value: 'dark',
              groupValue: currentValue,
              onChanged: (value) {
                ref.read(settingsProvider.notifier).setTheme(value!);
                Navigator.pop(context);
              },
            ),
            RadioListTile<String>(
              title: const Text('System'),
              value: 'system',
              groupValue: currentValue,
              onChanged: (value) {
                ref.read(settingsProvider.notifier).setTheme(value!);
                Navigator.pop(context);
              },
            ),
          ],
        ),
      ),
    );
  }

  void _showLanguageDialog(BuildContext context, WidgetRef ref, String currentValue) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Language'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            RadioListTile<String>(
              title: const Text('English'),
              value: 'en',
              groupValue: currentValue,
              onChanged: (value) {
                ref.read(settingsProvider.notifier).setLanguage(value!);
                Navigator.pop(context);
              },
            ),
            RadioListTile<String>(
              title: const Text('中文'),
              value: 'zh',
              groupValue: currentValue,
              onChanged: (value) {
                ref.read(settingsProvider.notifier).setLanguage(value!);
                Navigator.pop(context);
              },
            ),
          ],
        ),
      ),
    );
  }

  void _showBackupDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Create Backup'),
        content: const Text('This will create a backup of your database and configuration. Continue?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('Backup created successfully')),
              );
            },
            child: const Text('Create'),
          ),
        ],
      ),
    );
  }

  void _showLogoutDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Sign Out'),
        content: const Text('Are you sure you want to sign out?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              // Navigate to login
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.red,
            ),
            child: const Text('Sign Out'),
          ),
        ],
      ),
    );
  }
}