import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:anixops_mobile/core/services/schedules_api.dart';

void main() {
  group('Schedule model', () {
    test('is created correctly from JSON', () {
      final json = {
        'id': 1,
        'name': 'Daily Backup',
        'playbook_id': 5,
        'playbook_name': 'backup-system',
        'category': 'maintenance',
        'cron': '0 2 * * *',
        'timezone': 'UTC',
        'target_nodes': ['node-1', 'node-2'],
        'variables': {'backup_type': 'full'},
        'enabled': true,
        'next_run': '2026-03-21T02:00:00Z',
        'last_run': '2026-03-20T02:00:00Z',
        'last_task_id': 'task-123',
        'created_by': 1,
        'created_by_email': 'admin@example.com',
        'created_at': '2026-03-15T10:00:00Z',
        'updated_at': '2026-03-20T08:00:00Z',
      };

      final schedule = Schedule.fromJson(json);

      expect(schedule.id, 1);
      expect(schedule.name, 'Daily Backup');
      expect(schedule.playbookId, 5);
      expect(schedule.playbookName, 'backup-system');
      expect(schedule.category, 'maintenance');
      expect(schedule.cron, '0 2 * * *');
      expect(schedule.timezone, 'UTC');
      expect(schedule.targetNodes, ['node-1', 'node-2']);
      expect(schedule.variables, {'backup_type': 'full'});
      expect(schedule.enabled, true);
      expect(schedule.nextRun, isNotNull);
      expect(schedule.lastRun, isNotNull);
      expect(schedule.lastTaskId, 'task-123');
      expect(schedule.createdBy, 1);
      expect(schedule.createdByEmail, 'admin@example.com');
    });

    test('handles missing optional fields', () {
      final json = {
        'id': 2,
        'name': 'Test Schedule',
        'playbook_id': 1,
        'cron': '0 * * * *',
      };

      final schedule = Schedule.fromJson(json);

      expect(schedule.id, 2);
      expect(schedule.name, 'Test Schedule');
      expect(schedule.playbookId, 1);
      expect(schedule.playbookName, isNull);
      expect(schedule.category, isNull);
      expect(schedule.timezone, 'UTC');
      expect(schedule.targetNodes, isNull);
      expect(schedule.variables, isNull);
      // enabled defaults to false when not provided in JSON
      expect(schedule.enabled, false);
      expect(schedule.nextRun, isNull);
      expect(schedule.lastRun, isNull);
    });

    test('handles enabled as integer (0 or 1)', () {
      final jsonEnabled = {
        'id': 1,
        'name': 'Enabled Schedule',
        'playbook_id': 1,
        'cron': '0 * * * *',
        'enabled': 1,
      };

      final jsonDisabled = {
        'id': 2,
        'name': 'Disabled Schedule',
        'playbook_id': 1,
        'cron': '0 * * * *',
        'enabled': 0,
      };

      final enabledSchedule = Schedule.fromJson(jsonEnabled);
      final disabledSchedule = Schedule.fromJson(jsonDisabled);

      expect(enabledSchedule.enabled, true);
      expect(disabledSchedule.enabled, false);
    });

    test('cronDescription parses hourly cron', () {
      final schedule = Schedule(
        id: 1,
        name: 'Test',
        playbookId: 1,
        cron: '0 * * * *',
      );

      expect(schedule.cronDescription, 'Hourly');
    });

    test('cronDescription parses daily cron', () {
      final schedule = Schedule(
        id: 1,
        name: 'Test',
        playbookId: 1,
        cron: '30 14 * * *',
      );

      expect(schedule.cronDescription, contains('2:30'));
    });

    test('cronDescription parses every N minutes', () {
      final schedule = Schedule(
        id: 1,
        name: 'Test',
        playbookId: 1,
        cron: '*/15 * * * *',
      );

      expect(schedule.cronDescription, 'Every 15 minutes');
    });

    test('cronDescription returns original for invalid cron', () {
      final schedule = Schedule(
        id: 1,
        name: 'Test',
        playbookId: 1,
        cron: 'invalid',
      );

      expect(schedule.cronDescription, 'invalid');
    });
  });
}