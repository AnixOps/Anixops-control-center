import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

// Log Model
class LogEntry {
  final DateTime time;
  final String level;
  final String source;
  final String message;

  const LogEntry({
    required this.time,
    required this.level,
    required this.source,
    required this.message,
  });
}

// Logs State
class LogsState {
  final List<LogEntry> logs;
  final bool loading;
  final bool isStreaming;
  final Set<String> levelFilters;
  final String search;

  const LogsState({
    this.logs = const [],
    this.loading = false,
    this.isStreaming = false,
    this.levelFilters = const {'INFO', 'WARN', 'ERROR'},
    this.search = '',
  });

  LogsState copyWith({
    List<LogEntry>? logs,
    bool? loading,
    bool? isStreaming,
    Set<String>? levelFilters,
    String? search,
  }) {
    return LogsState(
      logs: logs ?? this.logs,
      loading: loading ?? this.loading,
      isStreaming: isStreaming ?? this.isStreaming,
      levelFilters: levelFilters ?? this.levelFilters,
      search: search ?? this.search,
    );
  }

  List<LogEntry> get filteredLogs {
    var result = logs.where((l) => levelFilters.contains(l.level)).toList();
    if (search.isNotEmpty) {
      result = result.where((l) =>
        l.message.toLowerCase().contains(search.toLowerCase()) ||
        l.source.toLowerCase().contains(search.toLowerCase())
      ).toList();
    }
    return result;
  }

  int get errorCount => logs.where((l) => l.level == 'ERROR').length;
  int get warnCount => logs.where((l) => l.level == 'WARN').length;
  int get infoCount => logs.where((l) => l.level == 'INFO').length;
}

class LogsNotifier extends StateNotifier<LogsState> {
  LogsNotifier() : super(const LogsState()) {
    loadLogs();
  }

  Future<void> loadLogs() async {
    state = state.copyWith(loading: true);
    await Future.delayed(const Duration(milliseconds: 500));

    final now = DateTime.now();
    final logs = List.generate(100, (i) {
      final levels = ['INFO', 'INFO', 'INFO', 'WARN', 'ERROR'];
      final sources = ['api', 'node', 'auth', 'plugin', 'system'];
      final messages = [
        'Request processed successfully',
        'Node connection established',
        'User authenticated',
        'High memory usage detected',
        'Connection timeout',
        'Database query executed',
        'Configuration updated',
        'Plugin loaded',
        'Certificate renewal scheduled',
        'Backup completed',
      ];

      return LogEntry(
        time: now.subtract(Duration(minutes: i)),
        level: levels[i % levels.length],
        source: sources[i % sources.length],
        message: messages[i % messages.length],
      );
    });

    state = state.copyWith(logs: logs, loading: false);
  }

  void toggleStreaming() {
    state = state.copyWith(isStreaming: !state.isStreaming);
  }

  void toggleLevelFilter(String level) {
    final filters = Set<String>.from(state.levelFilters);
    if (filters.contains(level)) {
      filters.remove(level);
    } else {
      filters.add(level);
    }
    state = state.copyWith(levelFilters: filters);
  }

  void setSearch(String search) {
    state = state.copyWith(search: search);
  }

  void clearLogs() {
    state = state.copyWith(logs: []);
  }

  void addLog(LogEntry log) {
    state = state.copyWith(logs: [log, ...state.logs]);
  }
}

final logsProvider = StateNotifierProvider<LogsNotifier, LogsState>((ref) {
  return LogsNotifier();
});

// Logs Page
class LogsPage extends ConsumerStatefulWidget {
  LogsPage({super.key});

  @override
  ConsumerState<LogsPage> createState() => _LogsPageState();
}

class _LogsPageState extends ConsumerState<LogsPage> {
  final ScrollController _scrollController = ScrollController();

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(logsProvider);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Logs'),
        actions: [
          IconButton(
            icon: Icon(
              state.isStreaming ? Icons.pause : Icons.play_arrow,
            ),
            onPressed: () => ref.read(logsProvider.notifier).toggleStreaming(),
            tooltip: state.isStreaming ? 'Stop Stream' : 'Start Stream',
          ),
          IconButton(
            icon: const Icon(Icons.delete_outline),
            onPressed: () => ref.read(logsProvider.notifier).clearLogs(),
            tooltip: 'Clear Logs',
          ),
        ],
      ),
      body: Column(
        children: [
          // Stats Row
          Container(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                _buildStatChip('Errors', state.errorCount, Colors.red),
                const SizedBox(width: 8),
                _buildStatChip('Warnings', state.warnCount, Colors.orange),
                const SizedBox(width: 8),
                _buildStatChip('Info', state.infoCount, Colors.blue),
                if (state.isStreaming) ...[
                  const SizedBox(width: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: Colors.green.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Container(
                          width: 8,
                          height: 8,
                          decoration: const BoxDecoration(
                            color: Colors.green,
                            shape: BoxShape.circle,
                          ),
                        ),
                        const SizedBox(width: 6),
                        const Text(
                          'Live',
                          style: TextStyle(
                            color: Colors.green,
                            fontSize: 12,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ],
            ),
          ),

          // Filter Chips
          SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Row(
              children: ['ERROR', 'WARN', 'INFO', 'DEBUG'].map((level) {
                final isSelected = state.levelFilters.contains(level);
                return Padding(
                  padding: const EdgeInsets.only(right: 8),
                  child: FilterChip(
                    label: Text(level),
                    selected: isSelected,
                    onSelected: (_) => ref.read(logsProvider.notifier).toggleLevelFilter(level),
                    selectedColor: _getLevelColor(level).withOpacity(0.2),
                    checkmarkColor: _getLevelColor(level),
                  ),
                );
              }).toList(),
            ),
          ),
          const SizedBox(height: 8),

          // Search
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: TextField(
              onChanged: (value) => ref.read(logsProvider.notifier).setSearch(value),
              decoration: InputDecoration(
                hintText: 'Search logs...',
                prefixIcon: const Icon(Icons.search, size: 20),
                filled: true,
                fillColor: theme.colorScheme.surface,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: BorderSide.none,
                ),
                contentPadding: EdgeInsets.zero,
                isDense: true,
              ),
            ),
          ),
          const SizedBox(height: 8),

          // Logs List
          Expanded(
            child: state.loading
                ? const Center(child: CircularProgressIndicator())
                : state.filteredLogs.isEmpty
                    ? Center(
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Icon(
                              Icons.article_outlined,
                              size: 64,
                              color: theme.colorScheme.onSurfaceVariant,
                            ),
                            const SizedBox(height: 16),
                            Text(
                              'No logs found',
                              style: theme.textTheme.bodyLarge?.copyWith(
                                color: theme.colorScheme.onSurfaceVariant,
                              ),
                            ),
                          ],
                        ),
                      )
                    : ListView.builder(
                        controller: _scrollController,
                        padding: const EdgeInsets.symmetric(horizontal: 16),
                        itemCount: state.filteredLogs.length,
                        itemBuilder: (context, index) {
                          final log = state.filteredLogs[index];
                          return _LogEntryCard(log: log);
                        },
                      ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatChip(String label, int count, Color color) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(
        '$count $label',
        style: TextStyle(
          color: color,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }

  Color _getLevelColor(String level) {
    switch (level) {
      case 'ERROR':
        return Colors.red;
      case 'WARN':
        return Colors.orange;
      case 'INFO':
        return Colors.blue;
      case 'DEBUG':
        return Colors.grey;
      default:
        return Colors.grey;
    }
  }
}

class _LogEntryCard extends StatelessWidget {
  final LogEntry log;

  const _LogEntryCard({required this.log});

  Color _getLevelColor() {
    switch (log.level) {
      case 'ERROR':
        return Colors.red;
      case 'WARN':
        return Colors.orange;
      case 'INFO':
        return Colors.blue;
      case 'DEBUG':
        return Colors.grey;
      default:
        return Colors.grey;
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final levelColor = _getLevelColor();

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: InkWell(
        onTap: () {
          // Show log detail
        },
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Text(
                    log.time.toString().substring(11, 19),
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: theme.colorScheme.onSurfaceVariant,
                      fontFamily: 'monospace',
                    ),
                  ),
                  const SizedBox(width: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                    decoration: BoxDecoration(
                      color: levelColor.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(
                      log.level,
                      style: TextStyle(
                        color: levelColor,
                        fontSize: 10,
                        fontWeight: FontWeight.bold,
                        fontFamily: 'monospace',
                      ),
                    ),
                  ),
                  const SizedBox(width: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                    decoration: BoxDecoration(
                      color: theme.colorScheme.surfaceVariant,
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(
                      log.source,
                      style: theme.textTheme.labelSmall,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              Text(
                log.message,
                style: theme.textTheme.bodyMedium,
              ),
            ],
          ),
        ),
      ),
    );
  }
}