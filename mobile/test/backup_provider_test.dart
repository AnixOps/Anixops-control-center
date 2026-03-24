import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:anixops_mobile/core/services/backup_api.dart';

void main() {
  group('Backup model', () {
    test('is created correctly from JSON', () {
      final json = {
        'id': 'backup-123',
        'name': 'Daily Backup',
        'status': 'completed',
        'size': 1048576,
        'description': 'System backup',
        'created_at': '2026-03-20T10:00:00Z',
        'completed_at': '2026-03-20T10:05:00Z',
      };

      final backup = Backup.fromJson(json);

      expect(backup.id, 'backup-123');
      expect(backup.name, 'Daily Backup');
      expect(backup.status, 'completed');
      expect(backup.size, 1048576);
      expect(backup.description, 'System backup');
      expect(backup.isCompleted, true);
      expect(backup.isPending, false);
      expect(backup.isFailed, false);
    });

    test('handles missing optional fields', () {
      final json = <String, dynamic>{
        'id': '2',
        'name': 'Test Backup',
        'status': 'pending',
      };

      final backup = Backup.fromJson(json);

      expect(backup.id, '2');
      expect(backup.name, 'Test Backup');
      expect(backup.status, 'pending');
      expect(backup.size, 0);
      expect(backup.description, isNull);
      expect(backup.completedAt, isNull);
      expect(backup.error, isNull);
      expect(backup.isPending, true);
    });

    test('formattedSize formats bytes correctly', () {
      final now = DateTime.now();
      expect(Backup(id: '', name: '', status: '', size: 500, createdAt: now).formattedSize, '500 B');
      expect(Backup(id: '', name: '', status: '', size: 1024, createdAt: now).formattedSize, '1.0 KB');
      expect(Backup(id: '', name: '', status: '', size: 1048576, createdAt: now).formattedSize, '1.0 MB');
      expect(Backup(id: '', name: '', status: '', size: 1073741824, createdAt: now).formattedSize, '1.0 GB');
    });

    test('isCompleted/isPending/isFailed work correctly', () {
      final now = DateTime.now();
      final completed = Backup(id: '', name: '', status: 'completed', createdAt: now);
      final pending = Backup(id: '', name: '', status: 'pending', createdAt: now);
      final failed = Backup(id: '', name: '', status: 'failed', createdAt: now);
      final running = Backup(id: '', name: '', status: 'running', createdAt: now);

      expect(completed.isCompleted, true);
      expect(pending.isPending, true);
      expect(failed.isFailed, true);
      expect(running.isCompleted, false);
      expect(running.isPending, false);
      expect(running.isFailed, false);
    });
  });

  group('BackupStatus model', () {
    test('is created correctly from JSON', () {
      final json = {
        'is_running': true,
        'last_backup_at': '2026-03-20T10:00:00Z',
        'total_backups': 10,
        'total_size': 10737418240,
      };

      final status = BackupStatus.fromJson(json);

      expect(status.isRunning, true);
      expect(status.lastBackupAt, isNotNull);
      expect(status.totalBackups, 10);
      expect(status.totalSize, 10737418240);
    });

    test('handles missing fields with defaults', () {
      final json = <String, dynamic>{};

      final status = BackupStatus.fromJson(json);

      expect(status.isRunning, false);
      expect(status.lastBackupAt, isNull);
      expect(status.totalBackups, 0);
      expect(status.totalSize, 0);
    });
  });
}