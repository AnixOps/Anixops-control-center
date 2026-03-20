import 'package:dio/dio.dart';
import 'dart:convert';

/// Task model
class Task {
  final String taskId;
  final int? playbookId;
  final String? playbookName;
  final String status;
  final String triggerType;
  final int? triggeredBy;
  final String? triggeredByEmail;
  final List<TargetNode>? targetNodes;
  final Map<String, dynamic>? variables;
  final Map<String, dynamic>? result;
  final String? error;
  final DateTime? createdAt;
  final DateTime? startedAt;
  final DateTime? completedAt;
  final String? category;

  Task({
    required this.taskId,
    this.playbookId,
    this.playbookName,
    required this.status,
    required this.triggerType,
    this.triggeredBy,
    this.triggeredByEmail,
    this.targetNodes,
    this.variables,
    this.result,
    this.error,
    this.createdAt,
    this.startedAt,
    this.completedAt,
    this.category,
  });

  factory Task.fromJson(Map<String, dynamic> json) {
    List<TargetNode>? nodes;
    if (json['target_nodes'] != null) {
      if (json['target_nodes'] is String) {
        try {
          final decoded = jsonDecode(json['target_nodes'] as String);
          if (decoded is List) {
            nodes = decoded.map((n) => TargetNode.fromJson(n as Map<String, dynamic>)).toList();
          }
        } catch (_) {
          nodes = null;
        }
      } else if (json['target_nodes'] is List) {
        nodes = (json['target_nodes'] as List)
            .map((n) => TargetNode.fromJson(n as Map<String, dynamic>))
            .toList();
      }
    }

    Map<String, dynamic>? vars;
    if (json['variables'] != null) {
      if (json['variables'] is String) {
        try {
          vars = jsonDecode(json['variables'] as String) as Map<String, dynamic>;
        } catch (_) {
          vars = null;
        }
      } else {
        vars = json['variables'] as Map<String, dynamic>?;
      }
    }

    return Task(
      taskId: json['task_id'] as String? ?? json['id'].toString(),
      playbookId: json['playbook_id'] as int?,
      playbookName: json['playbook_name'] as String?,
      status: json['status'] as String? ?? 'pending',
      triggerType: json['trigger_type'] as String? ?? 'manual',
      triggeredBy: json['triggered_by'] as int?,
      triggeredByEmail: json['triggered_by_email'] as String?,
      targetNodes: nodes,
      variables: vars,
      result: json['result'] as Map<String, dynamic>?,
      error: json['error'] as String?,
      createdAt: json['created_at'] != null
          ? DateTime.tryParse(json['created_at'] as String)
          : null,
      startedAt: json['started_at'] != null
          ? DateTime.tryParse(json['started_at'] as String)
          : null,
      completedAt: json['completed_at'] != null
          ? DateTime.tryParse(json['completed_at'] as String)
          : null,
      category: json['category'] as String?,
    );
  }

  String get title => playbookName ?? 'Task $taskId';
}

/// Target node model
class TargetNode {
  final int id;
  final String name;
  final String? host;

  TargetNode({
    required this.id,
    required this.name,
    this.host,
  });

  factory TargetNode.fromJson(Map<String, dynamic> json) {
    return TargetNode(
      id: json['id'] as int? ?? 0,
      name: json['name'] as String? ?? 'Unknown',
      host: json['host'] as String?,
    );
  }
}

/// Task log model
class TaskLog {
  final int? id;
  final String taskId;
  final int? nodeId;
  final String? nodeName;
  final String level;
  final String message;
  final Map<String, dynamic>? metadata;
  final DateTime? createdAt;

  TaskLog({
    this.id,
    required this.taskId,
    this.nodeId,
    this.nodeName,
    required this.level,
    required this.message,
    this.metadata,
    this.createdAt,
  });

  factory TaskLog.fromJson(Map<String, dynamic> json) {
    return TaskLog(
      id: json['id'] as int?,
      taskId: json['task_id'] as String? ?? '',
      nodeId: json['node_id'] as int?,
      nodeName: json['node_name'] as String?,
      level: json['level'] as String? ?? 'info',
      message: json['message'] as String? ?? '',
      metadata: json['metadata'] as Map<String, dynamic>?,
      createdAt: json['created_at'] != null
          ? DateTime.tryParse(json['created_at'] as String)
          : null,
    );
  }
}

/// Tasks API endpoints
class TasksApi {
  final Dio _dio;

  TasksApi(this._dio);

  /// Get all tasks
  Future<List<Task>> getTasks({
    String? status,
    int? playbookId,
    int page = 1,
    int perPage = 20,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'per_page': perPage,
    };
    if (status != null) queryParams['status'] = status;
    if (playbookId != null) queryParams['playbook_id'] = playbookId;

    final response = await _dio.get('/tasks', queryParameters: queryParams);
    final data = response.data['data'] as Map<String, dynamic>;
    final items = data['items'] as List<dynamic>;
    return items.map((json) => Task.fromJson(json as Map<String, dynamic>)).toList();
  }

  /// Get single task
  Future<Task> getTask(String taskId) async {
    final response = await _dio.get('/tasks/$taskId');
    return Task.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  /// Create task (run playbook)
  Future<Task> createTask({
    int? playbookId,
    String? playbookName,
    required List<dynamic> targetNodes,
    Map<String, dynamic>? variables,
  }) async {
    final response = await _dio.post('/tasks', data: {
      if (playbookId != null) 'playbook_id': playbookId,
      if (playbookName != null) 'playbook_name': playbookName,
      'target_nodes': targetNodes,
      if (variables != null) 'variables': variables,
    });
    return Task.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  /// Cancel task
  Future<void> cancelTask(String taskId) async {
    await _dio.post('/tasks/$taskId/cancel');
  }

  /// Retry task
  Future<Task> retryTask(String taskId) async {
    final response = await _dio.post('/tasks/$taskId/retry');
    return Task.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  /// Get task logs
  Future<List<TaskLog>> getTaskLogs(
    String taskId, {
    String level = 'all',
    int limit = 1000,
    int offset = 0,
  }) async {
    final response = await _dio.get('/tasks/$taskId/logs', queryParameters: {
      'level': level,
      'limit': limit,
      'offset': offset,
    });
    final data = response.data['data'] as List<dynamic>;
    return data.map((json) => TaskLog.fromJson(json as Map<String, dynamic>)).toList();
  }
}