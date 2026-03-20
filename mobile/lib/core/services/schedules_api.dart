import 'package:dio/dio.dart';
import 'dart:convert';

/// Schedule model
class Schedule {
  final int? id;
  final String name;
  final int? playbookId;
  final String? playbookName;
  final String? category;
  final String cron;
  final String timezone;
  final List<dynamic>? targetNodes;
  final Map<String, dynamic>? variables;
  final bool enabled;
  final DateTime? nextRun;
  final DateTime? lastRun;
  final String? lastTaskId;
  final int? createdBy;
  final String? createdByEmail;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  Schedule({
    this.id,
    required this.name,
    this.playbookId,
    this.playbookName,
    this.category,
    required this.cron,
    this.timezone = 'UTC',
    this.targetNodes,
    this.variables,
    this.enabled = true,
    this.nextRun,
    this.lastRun,
    this.lastTaskId,
    this.createdBy,
    this.createdByEmail,
    this.createdAt,
    this.updatedAt,
  });

  factory Schedule.fromJson(Map<String, dynamic> json) {
    List<dynamic>? nodes;
    if (json['target_nodes'] != null) {
      if (json['target_nodes'] is String) {
        try {
          nodes = jsonDecode(json['target_nodes'] as String) as List<dynamic>;
        } catch (_) {}
      } else {
        nodes = json['target_nodes'] as List<dynamic>;
      }
    }

    Map<String, dynamic>? vars;
    if (json['variables'] != null) {
      if (json['variables'] is String) {
        try {
          vars = jsonDecode(json['variables'] as String) as Map<String, dynamic>;
        } catch (_) {}
      } else {
        vars = json['variables'] as Map<String, dynamic>?;
      }
    }

    return Schedule(
      id: json['id'] as int?,
      name: json['name'] as String? ?? '',
      playbookId: json['playbook_id'] as int?,
      playbookName: json['playbook_name'] as String?,
      category: json['category'] as String?,
      cron: json['cron'] as String? ?? '0 * * * *',
      timezone: json['timezone'] as String? ?? 'UTC',
      targetNodes: nodes,
      variables: vars,
      enabled: (json['enabled'] as int?) == 1 || json['enabled'] == true,
      nextRun: json['next_run'] != null
          ? DateTime.tryParse(json['next_run'] as String)
          : null,
      lastRun: json['last_run'] != null
          ? DateTime.tryParse(json['last_run'] as String)
          : null,
      lastTaskId: json['last_task_id'] as String?,
      createdBy: json['created_by'] as int?,
      createdByEmail: json['created_by_email'] as String?,
      createdAt: json['created_at'] != null
          ? DateTime.tryParse(json['created_at'] as String)
          : null,
      updatedAt: json['updated_at'] != null
          ? DateTime.tryParse(json['updated_at'] as String)
          : null,
    );
  }

  String get cronDescription {
    final parts = cron.split(' ');
    if (parts.length != 5) return cron;

    // Every N minutes
    if (parts[0].startsWith('*/')) {
      final min = parts[0].substring(2);
      return 'Every $min minutes';
    }

    // Hourly
    if (parts[0] == '0' && parts[1] == '*') {
      return 'Hourly';
    }

    // Daily at specific time
    if (parts[1] != '*' && !parts[1].contains('/')) {
      final hour = int.tryParse(parts[1]) ?? 0;
      final minute = int.tryParse(parts[0]) ?? 0;
      final period = hour >= 12 ? 'PM' : 'AM';
      final displayHour = hour > 12 ? hour - 12 : (hour == 0 ? 12 : hour);
      return 'Daily at ${displayHour}:${minute.toString().padLeft(2, '0')} $period';
    }

    return cron;
  }
}

/// Schedules API endpoints
class SchedulesApi {
  final Dio _dio;

  SchedulesApi(this._dio);

  /// Get all schedules
  Future<List<Schedule>> getSchedules() async {
    final response = await _dio.get('/schedules');
    final data = response.data['data'] as List<dynamic>;
    return data.map((json) => Schedule.fromJson(json as Map<String, dynamic>)).toList();
  }

  /// Get single schedule
  Future<Schedule> getSchedule(int id) async {
    final response = await _dio.get('/schedules/$id');
    return Schedule.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  /// Create schedule
  Future<Schedule> createSchedule({
    required String name,
    required int playbookId,
    required String cron,
    String timezone = 'UTC',
    required List<dynamic> targetNodes,
    Map<String, dynamic>? variables,
    bool enabled = true,
  }) async {
    final response = await _dio.post('/schedules', data: {
      'name': name,
      'playbook_id': playbookId,
      'cron': cron,
      'timezone': timezone,
      'target_nodes': targetNodes,
      if (variables != null) 'variables': variables,
      'enabled': enabled,
    });
    return Schedule.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  /// Update schedule
  Future<Schedule> updateSchedule(int id, {
    String? name,
    String? cron,
    String? timezone,
    List<dynamic>? targetNodes,
    Map<String, dynamic>? variables,
    bool? enabled,
  }) async {
    final response = await _dio.put('/schedules/$id', data: {
      if (name != null) 'name': name,
      if (cron != null) 'cron': cron,
      if (timezone != null) 'timezone': timezone,
      if (targetNodes != null) 'target_nodes': targetNodes,
      if (variables != null) 'variables': variables,
      if (enabled != null) 'enabled': enabled,
    });
    return Schedule.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  /// Delete schedule
  Future<void> deleteSchedule(int id) async {
    await _dio.delete('/schedules/$id');
  }

  /// Toggle schedule enabled
  Future<bool> toggleSchedule(int id) async {
    final response = await _dio.post('/schedules/$id/toggle');
    return response.data['data']['enabled'] == 1;
  }

  /// Run schedule now
  Future<String?> runScheduleNow(int id) async {
    final response = await _dio.post('/schedules/$id/run');
    return response.data['data']['task_id'] as String?;
  }
}