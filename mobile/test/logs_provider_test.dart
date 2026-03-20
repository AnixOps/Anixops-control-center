import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:anixops_mobile/features/logs/presentation/providers/logs_provider.dart';

void main() {
  group('LogsProvider', () {
    test('initial state is correct', () {
      final container = ProviderContainer();
      final state = container.read(logsProvider);

      expect(state.logs, isEmpty);
      expect(state.loading, false);
      expect(state.error, isNull);
      expect(state.search, isEmpty);
      expect(state.levelFilter, isEmpty);
      expect(state.sourceFilter, isEmpty);
      expect(state.streaming, false);
      expect(state.page, 1);
      expect(state.hasMore, true);
    });

    group('LogEntry model', () {
      test('is created correctly from JSON', () {
        final json = {
          'id': 'log-123',
          'level': 'error',
          'source': 'node.tokyo-01',
          'message': 'Connection failed to upstream server',
          'timestamp': '2026-03-16T14:30:00Z',
          'metadata': {'node_id': 'tokyo-01', 'retry_count': 3},
        };

        final log = LogEntry.fromJson(json);

        expect(log.id, 'log-123');
        expect(log.level, 'error');
        expect(log.source, 'node.tokyo-01');
        expect(log.message, 'Connection failed to upstream server');
        expect(log.metadata, isNotNull);
        expect(log.metadata!['node_id'], 'tokyo-01');
        expect(log.metadata!['retry_count'], 3);
      });

      test('handles missing fields with defaults', () {
        final json = <String, dynamic>{};
        final log = LogEntry.fromJson(json);

        expect(log.id, '');
        expect(log.level, 'info');
        expect(log.source, '');
        expect(log.message, '');
        expect(log.metadata, isNull);
      });

      test('converts id to string if numeric', () {
        final json = {
          'id': 12345,
          'level': 'info',
          'source': 'system',
          'message': 'Test log',
          'timestamp': '2026-03-16T14:30:00Z',
        };

        final log = LogEntry.fromJson(json);

        expect(log.id, '12345');
      });
    });

    group('LogsState', () {
      test('copyWith works correctly', () {
        const original = LogsState();
        final updated = original.copyWith(
          loading: true,
          error: 'Test error',
          streaming: true,
        );

        expect(updated.loading, true);
        expect(updated.error, 'Test error');
        expect(updated.streaming, true);
        expect(updated.logs, isEmpty); // Unchanged
      });

      test('filteredLogs filters by search', () {
        final state = LogsState(
          logs: [
            LogEntry(id: '1', level: 'info', source: 'system', message: 'Server started', timestamp: DateTime.now()),
            LogEntry(id: '2', level: 'error', source: 'node', message: 'Connection error', timestamp: DateTime.now()),
            LogEntry(id: '3', level: 'info', source: 'system', message: 'Health check passed', timestamp: DateTime.now()),
          ],
          search: 'error',
        );

        expect(state.filteredLogs.length, 1);
        expect(state.filteredLogs.first.level, 'error');
      });

      test('filteredLogs filters by level', () {
        final state = LogsState(
          logs: [
            LogEntry(id: '1', level: 'info', source: 'system', message: 'Info 1', timestamp: DateTime.now()),
            LogEntry(id: '2', level: 'error', source: 'node', message: 'Error 1', timestamp: DateTime.now()),
            LogEntry(id: '3', level: 'error', source: 'node', message: 'Error 2', timestamp: DateTime.now()),
          ],
          levelFilter: 'error',
        );

        expect(state.filteredLogs.length, 2);
      });

      test('filteredLogs filters by source', () {
        final state = LogsState(
          logs: [
            LogEntry(id: '1', level: 'info', source: 'system', message: 'System log', timestamp: DateTime.now()),
            LogEntry(id: '2', level: 'info', source: 'node.tokyo', message: 'Node log', timestamp: DateTime.now()),
            LogEntry(id: '3', level: 'info', source: 'system', message: 'Another system log', timestamp: DateTime.now()),
          ],
          sourceFilter: 'system',
        );

        expect(state.filteredLogs.length, 2);
      });

      test('filteredLogs combines multiple filters', () {
        final state = LogsState(
          logs: [
            LogEntry(id: '1', level: 'error', source: 'system', message: 'System error', timestamp: DateTime.now()),
            LogEntry(id: '2', level: 'error', source: 'node', message: 'Node error', timestamp: DateTime.now()),
            LogEntry(id: '3', level: 'info', source: 'system', message: 'System info', timestamp: DateTime.now()),
            LogEntry(id: '4', level: 'error', source: 'system', message: 'Another system error', timestamp: DateTime.now()),
          ],
          search: 'another',
          levelFilter: 'error',
          sourceFilter: 'system',
        );

        expect(state.filteredLogs.length, 1);
        expect(state.filteredLogs.first.message, 'Another system error');
      });

      test('errorCount computes correctly', () {
        final state = LogsState(
          logs: [
            LogEntry(id: '1', level: 'info', source: 's1', message: 'm1', timestamp: DateTime.now()),
            LogEntry(id: '2', level: 'error', source: 's2', message: 'm2', timestamp: DateTime.now()),
            LogEntry(id: '3', level: 'warning', source: 's3', message: 'm3', timestamp: DateTime.now()),
            LogEntry(id: '4', level: 'error', source: 's4', message: 'm4', timestamp: DateTime.now()),
          ],
        );

        expect(state.errorCount, 2);
      });

      test('warningCount computes correctly', () {
        final state = LogsState(
          logs: [
            LogEntry(id: '1', level: 'info', source: 's1', message: 'm1', timestamp: DateTime.now()),
            LogEntry(id: '2', level: 'warning', source: 's2', message: 'm2', timestamp: DateTime.now()),
            LogEntry(id: '3', level: 'warning', source: 's3', message: 'm3', timestamp: DateTime.now()),
            LogEntry(id: '4', level: 'error', source: 's4', message: 'm4', timestamp: DateTime.now()),
          ],
        );

        expect(state.warningCount, 2);
      });

      test('sources returns unique list', () {
        final state = LogsState(
          logs: [
            LogEntry(id: '1', level: 'info', source: 'system', message: 'm1', timestamp: DateTime.now()),
            LogEntry(id: '2', level: 'info', source: 'node.tokyo', message: 'm2', timestamp: DateTime.now()),
            LogEntry(id: '3', level: 'info', source: 'system', message: 'm3', timestamp: DateTime.now()),
            LogEntry(id: '4', level: 'info', source: 'node.london', message: 'm4', timestamp: DateTime.now()),
          ],
        );

        expect(state.sources.length, 3);
        expect(state.sources, containsAll(['system', 'node.tokyo', 'node.london']));
      });

      test('hasMore defaults to true', () {
        const state = LogsState();
        expect(state.hasMore, true);
      });

      test('page defaults to 1', () {
        const state = LogsState();
        expect(state.page, 1);
      });
    });

    group('LogsNotifier', () {
      test('setSearch updates state', () {
        final container = ProviderContainer();
        final notifier = container.read(logsProvider.notifier);

        notifier.setSearch('error');

        expect(container.read(logsProvider).search, 'error');
      });

      test('setLevelFilter updates state', () {
        final container = ProviderContainer();
        final notifier = container.read(logsProvider.notifier);

        notifier.state = const LogsState(); // Reset state
        notifier.setLevelFilter('error');

        expect(container.read(logsProvider).levelFilter, 'error');
      });

      test('setSourceFilter updates state', () {
        final container = ProviderContainer();
        final notifier = container.read(logsProvider.notifier);

        notifier.setSourceFilter('system');

        expect(container.read(logsProvider).sourceFilter, 'system');
      });

      test('addLog prepends log to list', () {
        final container = ProviderContainer();
        final notifier = container.read(logsProvider.notifier);

        final log1 = LogEntry(
          id: '1',
          level: 'info',
          source: 'system',
          message: 'First log',
          timestamp: DateTime.now(),
        );
        final log2 = LogEntry(
          id: '2',
          level: 'error',
          source: 'node',
          message: 'Second log',
          timestamp: DateTime.now(),
        );

        notifier.addLog(log1);
        notifier.addLog(log2);

        final state = container.read(logsProvider);
        expect(state.logs.length, 2);
        expect(state.logs.first.id, '2'); // log2 was added last, should be first
      });

      test('startStreaming sets streaming to true', () {
        final container = ProviderContainer();
        final notifier = container.read(logsProvider.notifier);

        notifier.startStreaming();

        expect(container.read(logsProvider).streaming, true);
      });

      test('stopStreaming sets streaming to false', () {
        final container = ProviderContainer();
        final notifier = container.read(logsProvider.notifier);

        notifier.startStreaming();
        expect(container.read(logsProvider).streaming, true);

        notifier.stopStreaming();
        expect(container.read(logsProvider).streaming, false);
      });

      test('clearLogs removes all logs', () {
        final container = ProviderContainer();
        final notifier = container.read(logsProvider.notifier);

        notifier.addLog(LogEntry(
          id: '1',
          level: 'info',
          source: 'system',
          message: 'Test',
          timestamp: DateTime.now(),
        ));

        expect(container.read(logsProvider).logs.length, 1);

        notifier.clearLogs();

        expect(container.read(logsProvider).logs.length, 0);
      });
    });
  });
}