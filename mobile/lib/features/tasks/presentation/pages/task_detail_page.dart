import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/services/tasks_api.dart';
import '../providers/tasks_provider.dart';

class TaskDetailPage extends ConsumerStatefulWidget {
  final String taskId;

  const TaskDetailPage({super.key, required this.taskId});

  @override
  ConsumerState<TaskDetailPage> createState() => _TaskDetailPageState();
}

class _TaskDetailPageState extends ConsumerState<TaskDetailPage> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() {
      ref.read(tasksProvider.notifier).loadTask(widget.taskId);
      ref.read(tasksProvider.notifier).loadTaskLogs(widget.taskId);
    });
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(tasksProvider);
    final task = state.selectedTask;
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Task Details'),
        actions: [
          if (task != null && (task.status == 'pending' || task.status == 'running'))
            IconButton(
              icon: const Icon(Icons.cancel),
              tooltip: 'Cancel',
              onPressed: () => _cancelTask(context, task.taskId),
            ),
          if (task != null && (task.status == 'failed' || task.status == 'cancelled'))
            IconButton(
              icon: const Icon(Icons.refresh),
              tooltip: 'Retry',
              onPressed: () => _retryTask(context, task.taskId),
            ),
        ],
      ),
      body: state.isLoading
          ? const Center(child: CircularProgressIndicator())
          : task == null
              ? _buildNotFound()
              : _buildContent(context, task, state, theme),
    );
  }

  Widget _buildNotFound() {
    return const Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.search_off, size: 48, color: Colors.grey),
          SizedBox(height: 16),
          Text('Task not found'),
        ],
      ),
    );
  }

  Widget _buildContent(BuildContext context, Task task, TasksState state, ThemeData theme) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Status card
          _buildStatusCard(task, theme),
          const SizedBox(height: 16),

          // Info card
          _buildInfoCard(task, theme),
          const SizedBox(height: 16),

          // Target nodes card
          if (task.targetNodes != null && task.targetNodes!.isNotEmpty) ...[
            _buildNodesCard(task, theme),
            const SizedBox(height: 16),
          ],

          // Error card
          if (task.error != null && task.error!.isNotEmpty) ...[
            _buildErrorCard(task, theme),
            const SizedBox(height: 16),
          ],

          // Logs card
          _buildLogsCard(state, theme),
        ],
      ),
    );
  }

  Widget _buildStatusCard(Task task, ThemeData theme) {
    Color statusColor;
    IconData statusIcon;

    switch (task.status) {
      case 'pending':
        statusColor = Colors.grey;
        statusIcon = Icons.schedule;
        break;
      case 'running':
        statusColor = Colors.blue;
        statusIcon = Icons.play_circle;
        break;
      case 'success':
        statusColor = Colors.green;
        statusIcon = Icons.check_circle;
        break;
      case 'failed':
        statusColor = Colors.red;
        statusIcon = Icons.error;
        break;
      case 'cancelled':
        statusColor = Colors.orange;
        statusIcon = Icons.cancel;
        break;
      default:
        statusColor = Colors.grey;
        statusIcon = Icons.help;
    }

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              width: 64,
              height: 64,
              decoration: BoxDecoration(
                color: statusColor.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(
                statusIcon,
                size: 32,
                color: statusColor,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    task.title,
                    style: theme.textTheme.titleLarge?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    task.playbookName ?? 'Unknown playbook',
                    style: theme.textTheme.bodyMedium,
                  ),
                  const SizedBox(height: 4),
                  Text(
                    'ID: ${task.taskId.substring(0, 8)}...',
                    style: theme.textTheme.bodySmall?.copyWith(color: Colors.grey),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoCard(Task task, ThemeData theme) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Task Information', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 12),
            _buildInfoRow('Status', task.status.toUpperCase(), theme),
            _buildInfoRow('Trigger', task.triggerType, theme),
            if (task.triggeredByEmail != null)
              _buildInfoRow('Triggered by', task.triggeredByEmail!, theme),
            if (task.createdAt != null)
              _buildInfoRow('Created', _formatDateTime(task.createdAt!), theme),
            if (task.startedAt != null)
              _buildInfoRow('Started', _formatDateTime(task.startedAt!), theme),
            if (task.completedAt != null)
              _buildInfoRow('Completed', _formatDateTime(task.completedAt!), theme),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow(String label, String value, ThemeData theme) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: theme.textTheme.bodySmall?.copyWith(color: Colors.grey)),
          Text(value, style: theme.textTheme.bodyMedium),
        ],
      ),
    );
  }

  Widget _buildNodesCard(Task task, ThemeData theme) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text('Target Nodes', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: theme.colorScheme.primary.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    '${task.targetNodes?.length ?? 0} nodes',
                    style: TextStyle(fontSize: 12, color: theme.colorScheme.primary),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: task.targetNodes!.map((node) {
                return Container(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                  decoration: BoxDecoration(
                    color: Colors.grey.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      const Icon(Icons.computer, size: 16),
                      const SizedBox(width: 4),
                      Text(node.name),
                    ],
                  ),
                );
              }).toList(),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildErrorCard(Task task, ThemeData theme) {
    return Card(
      color: Colors.red.withValues(alpha: 0.1),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.error, color: Colors.red),
                const SizedBox(width: 8),
                Text('Error', style: theme.textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: Colors.red,
                )),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              task.error!,
              style: theme.textTheme.bodyMedium?.copyWith(color: Colors.red),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildLogsCard(TasksState state, ThemeData theme) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text('Execution Logs', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                if (state.isLoadingLogs)
                  const SizedBox(
                    width: 16,
                    height: 16,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  ),
              ],
            ),
            const SizedBox(height: 12),
            if (state.taskLogs.isEmpty)
              const Center(
                child: Padding(
                  padding: EdgeInsets.all(16),
                  child: Text('No logs available'),
                ),
              )
            else
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: theme.brightness == Brightness.dark
                      ? Colors.grey[900]
                      : Colors.grey[100],
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: state.taskLogs.take(50).map((log) {
                    Color logColor;
                    switch (log.level) {
                      case 'error':
                        logColor = Colors.red;
                        break;
                      case 'warning':
                        logColor = Colors.orange;
                        break;
                      default:
                        logColor = theme.brightness == Brightness.dark
                            ? Colors.grey[300]!
                            : Colors.grey[700]!;
                    }

                    return Padding(
                      padding: const EdgeInsets.only(bottom: 4),
                      child: Text(
                        '${log.nodeName != null ? "[${log.nodeName}] " : ""}${log.message}',
                        style: TextStyle(
                          fontFamily: 'monospace',
                          fontSize: 11,
                          color: logColor,
                        ),
                      ),
                    );
                  }).toList(),
                ),
              ),
          ],
        ),
      ),
    );
  }

  String _formatDateTime(DateTime dt) {
    return '${dt.year}-${dt.month.toString().padLeft(2, '0')}-${dt.day.toString().padLeft(2, '0')} '
           '${dt.hour.toString().padLeft(2, '0')}:${dt.minute.toString().padLeft(2, '0')}';
  }

  Future<void> _cancelTask(BuildContext context, String taskId) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Cancel Task'),
        content: const Text('Are you sure you want to cancel this task?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('No'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context, true),
            style: ElevatedButton.styleFrom(backgroundColor: Colors.orange),
            child: const Text('Yes, Cancel'),
          ),
        ],
      ),
    );

    if (confirmed == true && context.mounted) {
      final success = await ref.read(tasksProvider.notifier).cancelTask(taskId);
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(success ? 'Task cancelled' : 'Cancel failed')),
        );
      }
    }
  }

  Future<void> _retryTask(BuildContext context, String taskId) async {
    final task = await ref.read(tasksProvider.notifier).retryTask(taskId);
    if (context.mounted) {
      if (task != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Task retry created: ${task.taskId}')),
        );
        // Navigate to new task
        Navigator.pop(context);
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Retry failed')),
        );
      }
    }
  }
}