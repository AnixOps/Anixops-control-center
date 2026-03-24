import 'package:flutter_test/flutter_test.dart';
import 'package:anixops_mobile/core/services/tasks_api.dart';
import 'package:anixops_mobile/features/tasks/presentation/providers/tasks_provider.dart';

void main() {
  group('Task model', () {
    test('is created correctly from JSON', () {
      final json = {
        'task_id': 'task-123',
        'playbook_id': 1,
        'playbook_name': 'install-docker',
        'status': 'completed',
        'trigger_type': 'manual',
        'triggered_by': 1,
        'triggered_by_email': 'admin@test.com',
        'target_nodes': [
          {'id': 1, 'name': 'node-1', 'host': '10.0.0.1'}
        ],
        'variables': {'version': 'latest'},
        'result': {'success': true},
        'error': null,
        'created_at': '2026-03-20T10:00:00Z',
        'started_at': '2026-03-20T10:01:00Z',
        'completed_at': '2026-03-20T10:05:00Z',
        'category': 'software',
      };

      final task = Task.fromJson(json);

      expect(task.taskId, 'task-123');
      expect(task.playbookId, 1);
      expect(task.playbookName, 'install-docker');
      expect(task.status, 'completed');
      expect(task.triggerType, 'manual');
      expect(task.triggeredBy, 1);
      expect(task.triggeredByEmail, 'admin@test.com');
      expect(task.targetNodes?.length, 1);
      expect(task.targetNodes?[0].name, 'node-1');
      expect(task.variables, {'version': 'latest'});
      expect(task.result, {'success': true});
    });

    test('title uses playbookName if available', () {
      final task = Task(
        taskId: 'task-123',
        playbookName: 'install-docker',
        status: 'completed',
        triggerType: 'manual',
      );

      expect(task.title, 'install-docker');
    });

    test('title falls back to taskId', () {
      final task = Task(
        taskId: 'task-456',
        status: 'pending',
        triggerType: 'manual',
      );

      expect(task.title, 'Task task-456');
    });

    test('handles missing optional fields', () {
      final json = {
        'task_id': 'task-789',
        'status': 'pending',
      };

      final task = Task.fromJson(json);

      expect(task.taskId, 'task-789');
      expect(task.status, 'pending');
      expect(task.playbookId, isNull);
      expect(task.playbookName, isNull);
      expect(task.targetNodes, isNull);
      expect(task.variables, isNull);
    });
  });

  group('TargetNode model', () {
    test('is created correctly from JSON', () {
      final json = {
        'id': 1,
        'name': 'web-server',
        'host': '192.168.1.1',
      };

      final node = TargetNode.fromJson(json);

      expect(node.id, 1);
      expect(node.name, 'web-server');
      expect(node.host, '192.168.1.1');
    });

    test('handles missing fields with defaults', () {
      final json = <String, dynamic>{};

      final node = TargetNode.fromJson(json);

      expect(node.id, 0);
      expect(node.name, 'Unknown');
      expect(node.host, isNull);
    });
  });

  group('TaskLog model', () {
    test('is created correctly from JSON', () {
      final json = {
        'id': 1,
        'task_id': 'task-123',
        'node_id': 1,
        'node_name': 'node-1',
        'level': 'info',
        'message': 'Task started successfully',
        'metadata': {'duration': 5000},
        'created_at': '2026-03-20T10:00:00Z',
      };

      final log = TaskLog.fromJson(json);

      expect(log.id, 1);
      expect(log.taskId, 'task-123');
      expect(log.nodeId, 1);
      expect(log.nodeName, 'node-1');
      expect(log.level, 'info');
      expect(log.message, 'Task started successfully');
      expect(log.metadata, {'duration': 5000});
    });

    test('handles missing fields with defaults', () {
      final json = <String, dynamic>{};

      final log = TaskLog.fromJson(json);

      expect(log.id, isNull);
      expect(log.taskId, '');
      expect(log.nodeId, isNull);
      expect(log.nodeName, isNull);
      expect(log.level, 'info');
      expect(log.message, '');
    });
  });

  group('TasksState', () {
    test('filteredTasks returns all when no filter', () {
      final state = TasksState(
        tasks: [
          Task(taskId: '1', status: 'completed', triggerType: 'manual'),
          Task(taskId: '2', status: 'running', triggerType: 'manual'),
          Task(taskId: '3', status: 'pending', triggerType: 'manual'),
        ],
      );

      expect(state.filteredTasks.length, 3);
    });

    test('filteredTasks filters by status', () {
      final state = TasksState(
        tasks: [
          Task(taskId: '1', status: 'completed', triggerType: 'manual'),
          Task(taskId: '2', status: 'running', triggerType: 'manual'),
          Task(taskId: '3', status: 'completed', triggerType: 'manual'),
        ],
        statusFilter: 'completed',
      );

      expect(state.filteredTasks.length, 2);
      expect(state.filteredTasks.every((t) => t.status == 'completed'), true);
    });

    test('filteredTasks handles "all" filter', () {
      final state = TasksState(
        tasks: [
          Task(taskId: '1', status: 'completed', triggerType: 'manual'),
          Task(taskId: '2', status: 'running', triggerType: 'manual'),
        ],
        statusFilter: 'all',
      );

      expect(state.filteredTasks.length, 2);
    });

    test('copyWith works correctly', () {
      const state = TasksState(
        tasks: [],
        isLoading: false,
      );

      final newState = state.copyWith(
        isLoading: true,
        error: 'Test error',
      );

      expect(newState.isLoading, true);
      expect(newState.error, 'Test error');
      expect(newState.tasks, isEmpty);
    });
  });

  group('RC3 realtime helpers', () {
    test('task update can merge status for existing task', () {
      final original = Task(
        taskId: 'task-123',
        playbookName: 'deploy',
        status: 'pending',
        triggerType: 'manual',
      );

      final updated = Task(
        taskId: original.taskId,
        playbookId: original.playbookId,
        playbookName: original.playbookName,
        status: 'running',
        triggerType: original.triggerType,
        triggeredBy: original.triggeredBy,
        triggeredByEmail: original.triggeredByEmail,
        targetNodes: original.targetNodes,
        variables: original.variables,
        result: original.result,
        error: original.error,
        createdAt: original.createdAt,
        startedAt: original.startedAt,
        completedAt: original.completedAt,
        category: original.category,
      );

      expect(updated.taskId, 'task-123');
      expect(updated.status, 'running');
      expect(updated.playbookName, 'deploy');
    });

    test('task logs can append for selected task', () {
      final state = TasksState(
        selectedTask: Task(taskId: 'task-123', status: 'running', triggerType: 'manual'),
        taskLogs: const [],
      );

      final nextLogs = [
        ...state.taskLogs,
        TaskLog.fromJson({
          'task_id': 'task-123',
          'message': 'playbook step started',
          'level': 'info',
        }),
      ];

      expect(nextLogs.length, 1);
      expect(nextLogs.first.taskId, 'task-123');
      expect(nextLogs.first.message, 'playbook step started');
    });
  });
}
