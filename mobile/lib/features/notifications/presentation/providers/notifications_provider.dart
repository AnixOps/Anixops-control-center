import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/services/notifications_api.dart';
import '../../core/providers/api_providers.dart';

/// Notifications state
class NotificationsState {
  final List<Notification> notifications;
  final bool isLoading;
  final String? error;
  final int unreadCount;

  const NotificationsState({
    this.notifications = const [],
    this.isLoading = false,
    this.error,
    this.unreadCount = 0,
  });

  NotificationsState copyWith({
    List<Notification>? notifications,
    bool? isLoading,
    String? error,
    int? unreadCount,
  }) {
    return NotificationsState(
      notifications: notifications ?? this.notifications,
      isLoading: isLoading ?? this.isLoading,
      error: error,
      unreadCount: unreadCount ?? this.unreadCount,
    );
  }

  List<Notification> get unreadNotifications =>
      notifications.where((n) => !n.read).toList();

  List<Notification> get readNotifications =>
      notifications.where((n) => n.read).toList();
}

/// Notifications notifier
class NotificationsNotifier extends StateNotifier<NotificationsState> {
  final NotificationsApi _api;

  NotificationsNotifier(this._api) : super(const NotificationsState());

  Future<void> loadNotifications() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final notifications = await _api.getNotifications();
      state = state.copyWith(
        notifications: notifications,
        isLoading: false,
        unreadCount: notifications.where((n) => !n.read).length,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  Future<void> markAsRead(String id) async {
    try {
      await _api.markAsRead(id);
      final notifications = state.notifications.map((n) {
        if (n.id == id) {
          return Notification(
            id: n.id,
            title: n.title,
            message: n.message,
            type: n.type,
            read: true,
            createdAt: n.createdAt,
            metadata: n.metadata,
          );
        }
        return n;
      }).toList();
      state = state.copyWith(
        notifications: notifications,
        unreadCount: notifications.where((n) => !n.read).length,
      );
    } catch (e) {
      // Ignore error, still update UI
    }
  }

  Future<void> markAllAsRead() async {
    try {
      await _api.markAllAsRead();
      final notifications = state.notifications.map((n) {
        return Notification(
          id: n.id,
          title: n.title,
          message: n.message,
          type: n.type,
          read: true,
          createdAt: n.createdAt,
          metadata: n.metadata,
        );
      }).toList();
      state = state.copyWith(
        notifications: notifications,
        unreadCount: 0,
      );
    } catch (e) {
      // Ignore error, still update UI
    }
  }

  Future<void> deleteNotification(String id) async {
    try {
      await _api.deleteNotification(id);
      final notifications = state.notifications.where((n) => n.id != id).toList();
      state = state.copyWith(
        notifications: notifications,
        unreadCount: notifications.where((n) => !n.read).length,
      );
    } catch (e) {
      // Still update UI
      final notifications = state.notifications.where((n) => n.id != id).toList();
      state = state.copyWith(
        notifications: notifications,
        unreadCount: notifications.where((n) => !n.read).length,
      );
    }
  }

  Future<int> fetchUnreadCount() async {
    try {
      final count = await _api.getUnreadCount();
      state = state.copyWith(unreadCount: count);
      return count;
    } catch (e) {
      return state.unreadCount;
    }
  }

  void clearError() {
    state = state.copyWith(error: null);
  }
}

/// Provider for notifications
final notificationsProvider =
    StateNotifierProvider<NotificationsNotifier, NotificationsState>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return NotificationsNotifier(apiClient.notifications);
});

/// Provider for unread count (can be watched separately)
final unreadCountProvider = Provider<int>((ref) {
  return ref.watch(notificationsProvider).unreadCount;
});