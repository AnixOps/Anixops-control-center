import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:anixops_mobile/core/theme/app_theme.dart';

/// Desktop sidebar navigation
class DesktopSidebar extends StatelessWidget {
  final int selectedIndex;
  final ValueChanged<int>? onDestinationSelected;

  const DesktopSidebar({
    super.key,
    required this.selectedIndex,
    this.onDestinationSelected,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 260,
      decoration: BoxDecoration(
        color: AppTheme.darkSurface,
        border: Border(
          right: BorderSide(color: AppTheme.darkBorder),
        ),
      ),
      child: Column(
        children: [
          // Logo
          _buildLogo(context),

          // Navigation items
          Expanded(
            child: _buildNavItems(context),
          ),

          // User profile section
          _buildUserSection(context),
        ],
      ),
    );
  }

  Widget _buildLogo(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 24),
      child: Row(
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: AppTheme.primaryColor,
              borderRadius: BorderRadius.circular(10),
            ),
            child: const Icon(
              Icons.dashboard_rounded,
              color: Colors.white,
              size: 24,
            ),
          ),
          const SizedBox(width: 12),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'AnixOps',
                style: Theme.of(context).textTheme.titleLarge?.copyWith(
                      fontWeight: FontWeight.bold,
                      color: AppTheme.darkText,
                    ),
              ),
              Text(
                'Control Center',
                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                      color: AppTheme.darkTextSecondary,
                    ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildNavItems(BuildContext context) {
    final items = _getNavItems();

    return ListView.builder(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      itemCount: items.length,
      itemBuilder: (context, index) {
        final item = items[index];
        final isSelected = index == selectedIndex;

        return _NavItem(
          icon: item.icon,
          label: item.label,
          isSelected: isSelected,
          onTap: () {
            onDestinationSelected?.call(index);
            context.go(item.route);
          },
        );
      },
    );
  }

  Widget _buildUserSection(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        border: Border(
          top: BorderSide(color: AppTheme.darkBorder),
        ),
      ),
      child: Row(
        children: [
          CircleAvatar(
            radius: 18,
            backgroundColor: AppTheme.primaryColor,
            child: const Text(
              'A',
              style: TextStyle(
                color: Colors.white,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Admin',
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                        color: AppTheme.darkText,
                        fontWeight: FontWeight.w500,
                      ),
                ),
                Text(
                  'admin@anixops.dev',
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: AppTheme.darkTextSecondary,
                      ),
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(Icons.logout_rounded, color: AppTheme.darkTextSecondary),
            onPressed: () {
              // Handle logout
              context.go('/login');
            },
            tooltip: 'Logout',
          ),
        ],
      ),
    );
  }

  List<_NavItemData> _getNavItems() {
    return [
      _NavItemData(Icons.dashboard_rounded, 'Dashboard', '/dashboard'),
      _NavItemData(Icons.dns_rounded, 'Nodes', '/nodes'),
      _NavItemData(Icons.play_circle_rounded, 'Playbooks', '/playbooks'),
      _NavItemData(Icons.extension_rounded, 'Plugins', '/plugins'),
      _NavItemData(Icons.people_rounded, 'Users', '/users'),
      _NavItemData(Icons.description_rounded, 'Logs', '/logs'),
      _NavItemData(Icons.settings_rounded, 'Settings', '/settings'),
    ];
  }
}

class _NavItemData {
  final IconData icon;
  final String label;
  final String route;

  _NavItemData(this.icon, this.label, this.route);
}

class _NavItem extends StatelessWidget {
  final IconData icon;
  final String label;
  final bool isSelected;
  final VoidCallback? onTap;

  const _NavItem({
    required this.icon,
    required this.label,
    this.isSelected = false,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(8),
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 200),
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          decoration: BoxDecoration(
            color: isSelected ? AppTheme.primaryColor.withOpacity(0.15) : Colors.transparent,
            borderRadius: BorderRadius.circular(8),
          ),
          child: Row(
            children: [
              Icon(
                icon,
                size: 22,
                color: isSelected ? AppTheme.primaryColor : AppTheme.darkTextSecondary,
              ),
              const SizedBox(width: 14),
              Text(
                label,
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                  color: isSelected ? AppTheme.primaryColor : AppTheme.darkText,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}