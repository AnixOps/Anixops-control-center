import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:anixops_mobile/features/notifications/presentation/providers/notifications_provider.dart';
import 'package:anixops_mobile/core/services/notifications_api.dart';

void main() {
  group('NotificationsProvider', () {
    test('initial state is correct', () {
      final container = ProviderContainer();
      final state = container.read(notificationsProvider);

      expect(state.notifications, isEmpty);
      expect(state.isLoading, false);
      expect(state.error, isNull);
      expect(state.unreadCount, 0);
    });

    group('Notification model', () {
      test('is created correctly from JSON', () {
        final json = {
          'id': 'notif-123',
          'title': 'Test Notification',
          'message': 'This is a test message',
          'type': 'warning',
          'read': false,
          'created_at': '2026-03-20T10:00:00Z',
          'metadata': {'node_id': 'node-1'},
        };

        final notification = Notification.fromJson(json);

        expect(notification.id, 'notif-123');
        expect(notification.title, 'Test Notification');
        expect(notification.message, 'This is a test message');
        expect(notification.type, 'warning');
        expect(notification.read, false);
        expect(notification.createdAt, isNotNull);
        expect(notification.metadata, isNotNull);
        expect(notification.metadata!['node_id'], 'node-1');
      });

      test('handles missing optional fields', () {
        final json = <String, dynamic>{
          'id': '2',
          'title': 'Simple',
          'message': 'Message',
        };

        final notification = Notification.fromJson(json);

        expect(notification.id, '2');
        expect(notification.title, 'Simple');
        expect(notification.message, 'Message');
        expect(notification.type, 'info');
        expect(notification.read, false);
        expect(notification.metadata, isNull);
      });
    });

    group('NotificationsState', () {
      test('copyWith works correctly', () {
        const original = NotificationsState();
        final updated = original.copyWith(
          isLoading: true,
          error: 'Test error',
          unreadCount: 5,
        );

        expect(updated.isLoading, true);
        expect(updated.error, 'Test error');
        expect(updated.unreadCount, 5);
        expect(updated.notifications, isEmpty);
      });

      test('unreadNotifications filters correctly', () {
        final state = NotificationsState(
          notifications: [
            Notification(id: '1', title: '1', message: 'm1', type: 'info', read: false, createdAt: DateTime.now()),
            Notification(id: '2', title: '2', message: 'm2', type: 'info', read: true, createdAt: DateTime.now()),
            Notification(id: '3', title: '3', message: 'm3', type: 'info', read: false, createdAt: DateTime.now()),
          ],
        );

        expect(state.unreadNotifications.length, 2);
        expect(state.readNotifications.length, 1);
      });
    });
  });
}