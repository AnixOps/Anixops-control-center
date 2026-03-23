import 'package:dio/dio.dart';

/// Backup model
class Backup {
  final String id;
  final String name;
  final String status;
  final int size;
  final String? description;
  final DateTime createdAt;
  final DateTime? completedAt;
  final String? error;

  const Backup({
    required this.id,
    required this.name,
    required this.status,
    this.size = 0,
    this.description,
    required this.createdAt,
    this.completedAt,
    this.error,
  });

  factory Backup.fromJson(Map<String, dynamic> json) {
    return Backup(
      id: json['id']?.toString() ?? '',
      name: json['name'] ?? '',
      status: json['status'] ?? 'pending',
      size: json['size'] ?? 0,
      description: json['description'],
      createdAt: json['created_at'] != null
          ? DateTime.tryParse(json['created_at']) ?? DateTime.now()
          : DateTime.now(),
      completedAt: json['completed_at'] != null
          ? DateTime.tryParse(json['completed_at'])
          : null,
      error: json['error'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'status': status,
      'size': size,
      'description': description,
      'created_at': createdAt.toIso8601String(),
      'completed_at': completedAt?.toIso8601String(),
      'error': error,
    };
  }

  bool get isCompleted => status == 'completed';
  bool get isPending => status == 'pending';
  bool get isFailed => status == 'failed';

  String get formattedSize {
    if (size < 1024) return '$size B';
    if (size < 1024 * 1024) return '${(size / 1024).toStringAsFixed(1)} KB';
    if (size < 1024 * 1024 * 1024) {
      return '${(size / (1024 * 1024)).toStringAsFixed(1)} MB';
    }
    return '${(size / (1024 * 1024 * 1024)).toStringAsFixed(1)} GB';
  }
}

/// Backup status model
class BackupStatus {
  final bool isRunning;
  final DateTime? lastBackupAt;
  final int totalBackups;
  final int totalSize;

  const BackupStatus({
    this.isRunning = false,
    this.lastBackupAt,
    this.totalBackups = 0,
    this.totalSize = 0,
  });

  factory BackupStatus.fromJson(Map<String, dynamic> json) {
    return BackupStatus(
      isRunning: json['is_running'] ?? false,
      lastBackupAt: json['last_backup_at'] != null
          ? DateTime.tryParse(json['last_backup_at'])
          : null,
      totalBackups: json['total_backups'] ?? 0,
      totalSize: json['total_size'] ?? 0,
    );
  }
}

/// Backup API service
class BackupApi {
  final Dio _dio;

  BackupApi(this._dio);

  /// List backups
  Future<List<Backup>> listBackups({
    int limit = 50,
    int offset = 0,
  }) async {
    final response = await _dio.get('/backups', queryParameters: {
      'limit': limit,
      'offset': offset,
    });
    if (response.data['success'] == true) {
      return (response.data['data']['items'] as List)
          .map((json) => Backup.fromJson(json))
          .toList();
    }
    throw Exception(response.data['error'] ?? 'Failed to list backups');
  }

  /// Get backup status
  Future<BackupStatus> getStatus() async {
    final response = await _dio.get('/backups/status');
    if (response.data['success'] == true) {
      return BackupStatus.fromJson(response.data['data']);
    }
    throw Exception(response.data['error'] ?? 'Failed to get backup status');
  }

  /// Create backup
  Future<Backup> createBackup({
    String? name,
    String? description,
  }) async {
    final response = await _dio.post('/backups', data: {
      if (name != null) 'name': name,
      if (description != null) 'description': description,
    });
    if (response.data['success'] == true) {
      return Backup.fromJson(response.data['data']);
    }
    throw Exception(response.data['error'] ?? 'Failed to create backup');
  }

  /// Get backup details
  Future<Backup> getBackup(String id) async {
    final response = await _dio.get('/backups/$id');
    if (response.data['success'] == true) {
      return Backup.fromJson(response.data['data']);
    }
    throw Exception(response.data['error'] ?? 'Failed to get backup');
  }

  /// Restore from backup
  Future<void> restoreBackup(String id) async {
    final response = await _dio.post('/backups/$id/restore');
    if (response.data['success'] != true) {
      throw Exception(response.data['error'] ?? 'Failed to restore backup');
    }
  }

  /// Delete backup
  Future<void> deleteBackup(String id) async {
    final response = await _dio.delete('/backups/$id');
    if (response.data['success'] != true) {
      throw Exception(response.data['error'] ?? 'Failed to delete backup');
    }
  }

  /// Cleanup old backups
  Future<int> cleanupBackups({int keepLast = 10}) async {
    final response = await _dio.post('/backups/cleanup', data: {
      'keep_last': keepLast,
    });
    if (response.data['success'] == true) {
      return response.data['data']['deleted_count'] ?? 0;
    }
    throw Exception(response.data['error'] ?? 'Failed to cleanup backups');
  }
}