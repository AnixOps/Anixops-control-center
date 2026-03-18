import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:anixops_mobile/core/providers/api_providers.dart';

/// Log entry model
class LogEntry {
  final String id;
  final String level;
  final String source;
  final String message;
  final DateTime timestamp;
  final Map<String, dynamic>? metadata;

  const LogEntry({
    required this.id,
    required this.level,
    required this.source,
    required this.message,
    required this.timestamp,
    this.metadata,
  });

  factory LogEntry.fromJson(Map<String, dynamic> json) {
    return LogEntry(
      id: json['id']?.toString() ?? '',
      level: json['level'] ?? 'info',
      source: json['source'] ?? '',
      message: json['message'] ?? '',
      timestamp: json['timestamp'] != null
          ? DateTime.tryParse(json['timestamp']) ?? DateTime.now()
          : DateTime.now(),
      metadata: json['metadata'],
    );
  }
}

/// Logs state
class LogsState {
  final List<LogEntry> logs;
  final bool loading;
  final String? error;
  final String search;
  final String levelFilter;
  final String sourceFilter;
  final bool streaming;
  final int page;
  final bool hasMore;

  const LogsState({
    this.logs = const [],
    this.loading = false,
    this.error,
    this.search = '',
    this.levelFilter = '',
    this.sourceFilter = '',
    this.streaming = false,
    this.page = 1,
    this.hasMore = true,
  });

  LogsState copyWith({
    List<LogEntry>? logs,
    bool? loading,
    String? error,
    String? search,
    String? levelFilter,
    String? sourceFilter,
    bool? streaming,
    int? page,
    bool? hasMore,
  }) {
    return LogsState(
      logs: logs ?? this.logs,
      loading: loading ?? this.loading,
      error: error,
      search: search ?? this.search,
      levelFilter: levelFilter ?? this.levelFilter,
      sourceFilter: sourceFilter ?? this.sourceFilter,
      streaming: streaming ?? this.streaming,
      page: page ?? this.page,
      hasMore: hasMore ?? this.hasMore,
    );
  }

  List<LogEntry> get filteredLogs {
    var result = logs;
    if (search.isNotEmpty) {
      result = result
          .where((l) =>
              l.message.toLowerCase().contains(search.toLowerCase()) ||
              l.source.toLowerCase().contains(search.toLowerCase()))
          .toList();
    }
    if (levelFilter.isNotEmpty) {
      result = result.where((l) => l.level == levelFilter).toList();
    }
    if (sourceFilter.isNotEmpty) {
      result = result.where((l) => l.source == sourceFilter).toList();
    }
    return result;
  }

  int get errorCount => logs.where((l) => l.level == 'error').length;
  int get warningCount => logs.where((l) => l.level == 'warning').length;
  List<String> get sources => logs.map((l) => l.source).toSet().toList();
}

/// Logs notifier
class LogsNotifier extends StateNotifier<LogsState> {
  final Ref _ref;

  LogsNotifier(this._ref) : super(const LogsState());

  Future<void> fetchLogs({int? page, bool refresh = false}) async {
    if (state.loading) return;

    state = state.copyWith(loading: true, error: null);

    try {
      final api = _ref.read(apiClientProvider);
      final response = await api.dio.get('/logs', queryParameters: {
        if (state.levelFilter.isNotEmpty) 'level': state.levelFilter,
        if (state.sourceFilter.isNotEmpty) 'source': state.sourceFilter,
        'page': page ?? state.page,
        'limit': 50,
      });

      final data = response.data;
      final List<LogEntry> logs = (data['data'] ?? data)
          .map<LogEntry>((json) => LogEntry.fromJson(json))
          .toList();

      state = state.copyWith(
        logs: refresh ? logs : [...logs, ...state.logs],
        loading: false,
        page: page ?? state.page,
        hasMore: logs.length >= 50,
      );
    } catch (e) {
      state = state.copyWith(
        error: e.toString(),
        loading: false,
      );
    }
  }

  Future<void> refresh() => fetchLogs(page: 1, refresh: true);

  Future<void> loadMore() => fetchLogs(page: state.page + 1);

  void setSearch(String search) {
    state = state.copyWith(search: search);
  }

  void setLevelFilter(String filter) {
    state = state.copyWith(levelFilter: filter);
    fetchLogs(page: 1, refresh: true);
  }

  void setSourceFilter(String filter) {
    state = state.copyWith(sourceFilter: filter);
    fetchLogs(page: 1, refresh: true);
  }

  void addLog(LogEntry log) {
    state = state.copyWith(logs: [log, ...state.logs]);
  }

  void startStreaming() {
    state = state.copyWith(streaming: true);
  }

  void stopStreaming() {
    state = state.copyWith(streaming: false);
  }

  void clearLogs() {
    state = state.copyWith(logs: []);
  }
}

/// Provider for logs state
final logsProvider = StateNotifierProvider<LogsNotifier, LogsState>((ref) {
  return LogsNotifier(ref);
});