import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:anixops_mobile/core/theme/app_theme.dart';

import '../providers/nodes_provider.dart';

// Nodes Page
class NodesPage extends ConsumerStatefulWidget {
  const NodesPage({super.key});

  @override
  ConsumerState<NodesPage> createState() => _NodesPageState();
}

class _NodesPageState extends ConsumerState<NodesPage> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() => ref.read(nodesProvider.notifier).fetchNodes());
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(nodesProvider);
    final theme = Theme.of(context);

    return Scaffold(
      body: RefreshIndicator(
        onRefresh: () => ref.read(nodesProvider.notifier).refresh(),
        child: CustomScrollView(
          slivers: [
            // App Bar
            SliverAppBar(
              floating: true,
              snap: true,
              title: const Text('Nodes'),
              actions: [
                IconButton(
                  icon: const Icon(Icons.add),
                  onPressed: () {
                    // TODO: Add node
                  },
                ),
              ],
              bottom: PreferredSize(
                preferredSize: const Size.fromHeight(100),
                child: Column(
                  children: [
                    // Search Bar
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                      child: TextField(
                        onChanged: (value) => ref.read(nodesProvider.notifier).setSearch(value),
                        decoration: InputDecoration(
                          hintText: 'Search nodes...',
                          prefixIcon: const Icon(Icons.search),
                          filled: true,
                          fillColor: theme.colorScheme.surface,
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(12),
                            borderSide: BorderSide.none,
                          ),
                          contentPadding: EdgeInsets.zero,
                        ),
                      ),
                    ),
                    // Stats Row
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                      child: Row(
                        children: [
                          _buildStatCard('Total', state.nodes.length.toString(), Colors.blue, theme),
                          const SizedBox(width: 8),
                          _buildStatCard('Online', state.onlineCount.toString(), Colors.green, theme),
                          const SizedBox(width: 8),
                          _buildStatCard('Offline', state.offlineCount.toString(), Colors.red, theme),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),

            // Filter Chips
            SliverToBoxAdapter(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: SingleChildScrollView(
                  scrollDirection: Axis.horizontal,
                  child: Row(
                    children: [
                      FilterChip(
                        label: const Text('All'),
                        selected: state.statusFilter.isEmpty,
                        onSelected: (_) => ref.read(nodesProvider.notifier).setStatusFilter(''),
                      ),
                      const SizedBox(width: 8),
                      FilterChip(
                        label: const Text('Online'),
                        selected: state.statusFilter == 'online',
                        onSelected: (_) => ref.read(nodesProvider.notifier).setStatusFilter('online'),
                      ),
                      const SizedBox(width: 8),
                      FilterChip(
                        label: const Text('Offline'),
                        selected: state.statusFilter == 'offline',
                        onSelected: (_) => ref.read(nodesProvider.notifier).setStatusFilter('offline'),
                      ),
                    ],
                  ),
                ),
              ),
            ),

            // Nodes List
            state.loading
                ? const SliverFillRemaining(
                    child: Center(child: CircularProgressIndicator()),
                  )
                : state.filteredNodes.isEmpty
                    ? SliverFillRemaining(
                        child: Center(
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Icon(
                                Icons.dns_outlined,
                                size: 64,
                                color: theme.colorScheme.onSurfaceVariant,
                              ),
                              const SizedBox(height: 16),
                              Text(
                                'No nodes found',
                                style: theme.textTheme.bodyLarge?.copyWith(
                                  color: theme.colorScheme.onSurfaceVariant,
                                ),
                              ),
                            ],
                          ),
                        ),
                      )
                    : SliverPadding(
                        padding: const EdgeInsets.all(16),
                        sliver: SliverList(
                          delegate: SliverChildBuilderDelegate(
                            (context, index) {
                              final node = state.filteredNodes[index];
                              return _NodeCard(
                                node: node,
                                onStart: () => ref.read(nodesProvider.notifier).startNode(node.id),
                                onStop: () => ref.read(nodesProvider.notifier).stopNode(node.id),
                                onRestart: () => ref.read(nodesProvider.notifier).restartNode(node.id),
                              );
                            },
                            childCount: state.filteredNodes.length,
                          ),
                        ),
                      ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatCard(String label, String value, Color color, ThemeData theme) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: color.withOpacity(0.1),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: color.withOpacity(0.3)),
        ),
        child: Column(
          children: [
            Text(
              value,
              style: theme.textTheme.headlineSmall?.copyWith(
                fontWeight: FontWeight.bold,
                color: color,
              ),
            ),
            Text(
              label,
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _NodeCard extends StatelessWidget {
  final Node node;
  final VoidCallback? onStart;
  final VoidCallback? onStop;
  final VoidCallback? onRestart;

  const _NodeCard({
    required this.node,
    this.onStart,
    this.onStop,
    this.onRestart,
  });

  String _formatBytes(int bytes) {
    if (bytes == 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    final i = (bytes.bitLength - 1) ~/ 10;
    return '${(bytes / (1 << (i * 10))).toStringAsFixed(1)} ${sizes[i]}';
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isOnline = node.status == 'online';

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: () {
          // TODO: Navigate to node detail
        },
        borderRadius: BorderRadius.circular(16),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    width: 48,
                    height: 48,
                    decoration: BoxDecoration(
                      color: theme.colorScheme.primary.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(
                      Icons.dns,
                      color: theme.colorScheme.primary,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          children: [
                            Expanded(
                              child: Text(
                                node.name,
                                style: theme.textTheme.titleMedium?.copyWith(
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                            ),
                            Container(
                              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                              decoration: BoxDecoration(
                                color: isOnline
                                    ? Colors.green.withOpacity(0.1)
                                    : Colors.red.withOpacity(0.1),
                                borderRadius: BorderRadius.circular(12),
                              ),
                              child: Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  Container(
                                    width: 6,
                                    height: 6,
                                    decoration: BoxDecoration(
                                      color: isOnline ? Colors.green : Colors.red,
                                      shape: BoxShape.circle,
                                    ),
                                  ),
                                  const SizedBox(width: 4),
                                  Text(
                                    node.status,
                                    style: TextStyle(
                                      fontSize: 12,
                                      color: isOnline ? Colors.green : Colors.red,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 4),
                        Text(
                          '${node.host}:${node.port}',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: theme.colorScheme.onSurfaceVariant,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 12),
              Row(
                children: [
                  _buildInfoChip(
                    icon: Icons.people_outline,
                    label: '${node.users} users',
                    theme: theme,
                  ),
                  const SizedBox(width: 8),
                  _buildInfoChip(
                    icon: Icons.swap_vert,
                    label: _formatBytes(node.traffic),
                    theme: theme,
                  ),
                  const SizedBox(width: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: theme.colorScheme.surface,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Text(
                      node.type.toUpperCase(),
                      style: theme.textTheme.labelSmall,
                    ),
                  ),
                ],
              ),
              if (onStart != null || onStop != null || onRestart != null) ...[
                const SizedBox(height: 12),
                Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    if (isOnline && onStop != null)
                      TextButton.icon(
                        onPressed: onStop,
                        icon: const Icon(Icons.stop, size: 18),
                        label: const Text('Stop'),
                        style: TextButton.styleFrom(foregroundColor: Colors.red),
                      ),
                    if (!isOnline && onStart != null)
                      TextButton.icon(
                        onPressed: onStart,
                        icon: const Icon(Icons.play_arrow, size: 18),
                        label: const Text('Start'),
                        style: TextButton.styleFrom(foregroundColor: Colors.green),
                      ),
                    if (onRestart != null)
                      TextButton.icon(
                        onPressed: onRestart,
                        icon: const Icon(Icons.refresh, size: 18),
                        label: const Text('Restart'),
                      ),
                  ],
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildInfoChip({
    required IconData icon,
    required String label,
    required ThemeData theme,
  }) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(
          icon,
          size: 16,
          color: theme.colorScheme.onSurfaceVariant,
        ),
        const SizedBox(width: 4),
        Text(
          label,
          style: theme.textTheme.bodySmall?.copyWith(
            color: theme.colorScheme.onSurfaceVariant,
          ),
        ),
      ],
    );
  }
}