import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/services/schedules_api.dart';
import '../../../../core/providers/api_providers.dart';

/// Schedules state
class SchedulesState {
  final List<Schedule> schedules;
  final bool isLoading;
  final String? error;

  const SchedulesState({
    this.schedules = const [],
    this.isLoading = false,
    this.error,
  });

  SchedulesState copyWith({
    List<Schedule>? schedules,
    bool? isLoading,
    String? error,
  }) {
    return SchedulesState(
      schedules: schedules ?? this.schedules,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
    );
  }

  List<Schedule> get enabledSchedules => schedules.where((s) => s.enabled).toList();
  List<Schedule> get disabledSchedules => schedules.where((s) => !s.enabled).toList();
}

/// Schedules notifier
class SchedulesNotifier extends StateNotifier<SchedulesState> {
  final SchedulesApi _api;

  SchedulesNotifier(this._api) : super(const SchedulesState()) {
    loadSchedules();
  }

  /// Load all schedules
  Future<void> loadSchedules() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final schedules = await _api.getSchedules();
      state = state.copyWith(
        schedules: schedules,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// Create schedule
  Future<Schedule?> createSchedule({
    required String name,
    required int playbookId,
    required String cron,
    String timezone = 'UTC',
    required List<dynamic> targetNodes,
    Map<String, dynamic>? variables,
  }) async {
    try {
      final schedule = await _api.createSchedule(
        name: name,
        playbookId: playbookId,
        cron: cron,
        timezone: timezone,
        targetNodes: targetNodes,
        variables: variables,
      );
      await loadSchedules();
      return schedule;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return null;
    }
  }

  /// Update schedule
  Future<bool> updateSchedule(int id, {
    String? name,
    String? cron,
    String? timezone,
    List<dynamic>? targetNodes,
    Map<String, dynamic>? variables,
    bool? enabled,
  }) async {
    try {
      await _api.updateSchedule(
        id,
        name: name,
        cron: cron,
        timezone: timezone,
        targetNodes: targetNodes,
        variables: variables,
        enabled: enabled,
      );
      await loadSchedules();
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// Delete schedule
  Future<bool> deleteSchedule(int id) async {
    try {
      await _api.deleteSchedule(id);
      state = state.copyWith(
        schedules: state.schedules.where((s) => s.id != id).toList(),
      );
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// Toggle schedule enabled
  Future<bool> toggleSchedule(int id) async {
    try {
      final newEnabled = await _api.toggleSchedule(id);
      state = state.copyWith(
        schedules: state.schedules.map((s) {
          if (s.id == id) {
            return Schedule(
              id: s.id,
              name: s.name,
              playbookId: s.playbookId,
              playbookName: s.playbookName,
              category: s.category,
              cron: s.cron,
              timezone: s.timezone,
              targetNodes: s.targetNodes,
              variables: s.variables,
              enabled: newEnabled,
              nextRun: s.nextRun,
              lastRun: s.lastRun,
              lastTaskId: s.lastTaskId,
              createdBy: s.createdBy,
              createdByEmail: s.createdByEmail,
              createdAt: s.createdAt,
              updatedAt: s.updatedAt,
            );
          }
          return s;
        }).toList(),
      );
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// Run schedule now
  Future<String?> runScheduleNow(int id) async {
    try {
      final taskId = await _api.runScheduleNow(id);
      await loadSchedules();
      return taskId;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return null;
    }
  }

  /// Clear error
  void clearError() {
    state = state.copyWith(error: null);
  }
}

/// Provider for SchedulesState
final schedulesProvider = StateNotifierProvider<SchedulesNotifier, SchedulesState>((ref) {
  final client = ref.watch(apiClientProvider);
  return SchedulesNotifier(client.schedules);
});