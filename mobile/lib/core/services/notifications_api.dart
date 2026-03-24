import 'package:dio/dio.dart';

/// Notification model
class Notification {
  final String id;
  final String title;
  final String message;
  final String type;
  final bool read;
  final DateTime createdAt;
  final Map<String, dynamic>? metadata;

  const Notification({
    required this.id,
    required this.title,
    required this.message,
    required this.type,
    this.read = false,
    required this.createdAt,
    this.metadata,
  });

  factory Notification.fromJson(Map<String, dynamic> json) {
    return Notification(
      id: json['id']?.toString() ?? '',
      title: json['title'] ?? '',
      message: json['message'] ?? '',
      type: json['type'] ?? 'info',
      read: json['read'] ?? false,
      createdAt: json['created_at'] != null
          ? DateTime.tryParse(json['created_at']) ?? DateTime.now()
          : DateTime.now(),
      metadata: json['metadata'] as Map<String, dynamic>?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'message': message,
      'type': type,
      'read': read,
      'created_at': createdAt.toIso8601String(),
      'metadata': metadata,
    };
  }

  Notification copyWith({
    String? id,
    String? title,
    String? message,
    String? type,
    bool? read,
    DateTime? createdAt,
    Map<String, dynamic>? metadata,
  }) {
    return Notification(
      id: id ?? this.id,
      title: title ?? this.title,
      message: message ?? this.message,
      type: type ?? this.type,
      read: read ?? this.read,
      createdAt: createdAt ?? this.createdAt,
      metadata: metadata ?? this.metadata,
    );
  }
}

/// Notifications API service
class NotificationsApi {
  final Dio _dio;

  NotificationsApi(this._dio);

  /// List notifications
  Future<List<Notification>> getNotifications({
    int limit = 50,
    int offset = 0,
    bool? unreadOnly,
  }) async {
    final queryParams = <String, dynamic>{
      'limit': limit,
      'offset': offset,
    };
    if (unreadOnly == true) {
      queryParams['unread_only'] = true;
    }

    final response = await _dio.get('/notifications', queryParameters: queryParams);
    if (response.data['success'] == true) {
      return (response.data['data']['items'] as List)
          .map((json) => Notification.fromJson(json))
          .toList();
    }
    throw Exception(response.data['error'] ?? 'Failed to get notifications');
  }

  /// Get unread count
  Future<int> getUnreadCount() async {
    final response = await _dio.get('/notifications/unread-count');
    if (response.data['success'] == true) {
      return response.data['data']['count'] ?? 0;
    }
    return 0;
  }

  /// Mark notification as read
  Future<void> markAsRead(String id) async {
    final response = await _dio.put('/notifications/$id/read');
    if (response.data['success'] != true) {
      throw Exception(response.data['error'] ?? 'Failed to mark as read');
    }
  }

  /// Mark all as read
  Future<void> markAllAsRead() async {
    final response = await _dio.put('/notifications/read-all');
    if (response.data['success'] != true) {
      throw Exception(response.data['error'] ?? 'Failed to mark all as read');
    }
  }

  /// Delete notification
  Future<void> deleteNotification(String id) async {
    final response = await _dio.delete('/notifications/$id');
    if (response.data['success'] != true) {
      throw Exception(response.data['error'] ?? 'Failed to delete notification');
    }
  }
}