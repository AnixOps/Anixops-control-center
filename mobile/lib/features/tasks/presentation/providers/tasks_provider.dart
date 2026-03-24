import 'dart:async';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/services/tasks_api.dart';
import '../../../../core/providers/api_providers.dart';

/// Tasks state
class TasksState {
  final List<Task> tasks;
  final Task? selectedTask;
  final List<TaskLog> taskLogs;
  final bool isLoading;
  final bool isLoadingLogs;
  final String? error;
  final String? statusFilter;
  final int currentPage;
  final int totalPages;
  final int total;

  const TasksState({
    this.tasks = const [],
    this.selectedTask,
    this.taskLogs = const [],
    this.isLoading = false,
    this.isLoadingLogs = false,
    this.error,
    this.statusFilter,
    this.currentPage = 1,
    this.totalPages = 1,
    this.total = 0,
  });

  TasksState copyWith({
    List<Task>? tasks,
    Task? selectedTask,
    List<TaskLog>? taskLogs,
    bool? isLoading,
    bool? isLoadingLogs,
    String? error,
    String? statusFilter,
    int? currentPage,
    int? totalPages,
    int? total,
  }) {
    return TasksState(
      tasks: tasks ?? this.tasks,
      selectedTask: selectedTask ?? this.selectedTask,
      taskLogs: taskLogs ?? this.taskLogs,
      isLoading: isLoading ?? this.isLoading,
      isLoadingLogs: isLoadingLogs ?? this.isLoadingLogs,
      error: error ?? this.error,
      statusFilter: statusFilter ?? this.statusFilter,
      currentPage: currentPage ?? this.currentPage,
      totalPages: totalPages ?? this.totalPages,
      total: total ?? this.total,
    );
  }

  List<Task> get filteredTasks {
    if (statusFilter == null || statusFilter == 'all') {
      return tasks;
    }
    return tasks.where((t) => t.status == statusFilter).toList();
  }
}

/// Tasks notifier
class TasksNotifier extends StateNotifier<TasksState> {
  final TasksApi _api;
  final Ref _ref;
  StreamSubscription<bool>? _connectionSubscription;

  TasksNotifier(this._api, this._ref) : super(const TasksState()) {
    _bindRealtimeUpdates();
    loadTasks();
  }

  void _bindRealtimeUpdates() {
    final sse = _ref.read(sseServiceProvider);

    sse.on('task_update', _handleTaskUpdate);
    sse.on('log', _handleTaskLog);
    _connectionSubscription = sse.connectionState.listen((connected) {
      if (!connected) return;
      if (!sse.subscribedChannels.contains('tasks')) {
        unawaited(sse.subscribe('tasks'));
      }
      if (!sse.subscribedChannels.contains('logs')) {
        unawaited(sse.subscribe('logs'));
      }
    });
  }

  void _handleTaskUpdate(dynamic payload) {
    if (payload is! Map) return;

    final taskId = payload['task_id']?.toString();
    if (taskId == null || taskId.isEmpty) return;

    final updatedTasks = state.tasks.map((task) {
      if (task.taskId != taskId) return task;

      return Task(
        taskId: task.taskId,
        playbookId: task.playbookId,
        playbookName: task.playbookName,
        status: payload['status']?.toString() ?? task.status,
        triggerType: task.triggerType,
        triggeredBy: task.triggeredBy,
        triggeredByEmail: task.triggeredByEmail,
        targetNodes: task.targetNodes,
        variables: task.variables,
        result: payload['result'] is Map<String, dynamic>
            ? payload['result'] as Map<String, dynamic>
            : task.result,
        error: payload['error']?.toString() ?? task.error,
        createdAt: task.createdAt,
        startedAt: task.startedAt,
        completedAt: task.completedAt,
        category: task.category,
      );
    }).toList();

    Task? updatedSelectedTask = state.selectedTask;
    if (updatedSelectedTask?.taskId == taskId) {
      updatedSelectedTask = updatedTasks.firstWhere(
        (task) => task.taskId == taskId,
        orElse: () => updatedSelectedTask!,
      );
    }

    state = state.copyWith(tasks: updatedTasks, selectedTask: updatedSelectedTask);
  }

  void _handleTaskLog(dynamic payload) {
    if (payload is! Map) return;

    final taskId = payload['task_id']?.toString();
    if (taskId == null || taskId.isEmpty) return;
    if (state.selectedTask?.taskId != taskId) return;

    final log = TaskLog.fromJson(Map<String, dynamic>.from(payload));
    state = state.copyWith(taskLogs: [...state.taskLogs, log]);
  }

  /// Load tasks
  Future<void> loadTasks({int page = 1, String? status}) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final tasks = await _api.getTasks(
        page: page,
        status: status ?? state.statusFilter,
      );

      state = state.copyWith(
        tasks: tasks,
        isLoading: false,
        currentPage: page,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// Set status filter
  void setStatusFilter(String? status) {
    state = state.copyWith(statusFilter: status);
    loadTasks(page: 1, status: status);
  }

  /// Load task detail
  Future<void> loadTask(String taskId) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final task = await _api.getTask(taskId);
      state = state.copyWith(
        selectedTask: task,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// Load task logs
  Future<void> loadTaskLogs(String taskId) async {
    state = state.copyWith(isLoadingLogs: true);

    try {
      final logs = await _api.getTaskLogs(taskId);
      state = state.copyWith(
        taskLogs: logs,
        isLoadingLogs: false,
      );
    } catch (e) {
      state = state.copyWith(isLoadingLogs: false);
    }
  }

  /// Create task
  Future<Task?> createTask({
    int? playbookId,
    String? playbookName,
    required List<dynamic> targetNodes,
    Map<String, dynamic>? variables,
  }) async {
    try {
      final task = await _api.createTask(
        playbookId: playbookId,
        playbookName: playbookName,
        targetNodes: targetNodes,
        variables: variables,
      );
      await loadTasks();
      return task;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return null;
    }
  }

  /// Cancel task
  Future<bool> cancelTask(String taskId) async {
    try {
      await _api.cancelTask(taskId);
      await loadTasks();
      if (state.selectedTask?.taskId == taskId) {
        await loadTask(taskId);
      }
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// Retry task
  Future<Task?> retryTask(String taskId) async {
    try {
      final task = await _api.retryTask(taskId);
      await loadTasks();
      return task;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return null;
    }
  }

  /// Clear selected task
  void clearSelectedTask() {
    state = state.copyWith(
      selectedTask: null,
      taskLogs: const [],
    );
  }

  @override
  void dispose() {
    final sse = _ref.read(sseServiceProvider);
    sse.off('task_update', _handleTaskUpdate);
    sse.off('log', _handleTaskLog);
    _connectionSubscription?.cancel();
    super.dispose();
  }
}

/// Provider for TasksState
final tasksProvider = StateNotifierProvider<TasksNotifier, TasksState>((ref) {
  final client = ref.watch(apiClientProvider);
  return TasksNotifier(client.tasks, ref);
});