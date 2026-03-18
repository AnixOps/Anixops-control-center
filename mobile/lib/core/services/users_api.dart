import 'package:dio/dio.dart';

/// Users API endpoints
class UsersApi {
  final Dio _dio;

  UsersApi(this._dio);

  /// List all users with optional filters
  Future<Response> list({
    String? search,
    String? role,
    String? status,
    int page = 1,
    int limit = 20,
  }) async {
    return _dio.get('/users', queryParameters: {
      if (search != null) 'search': search,
      if (role != null) 'role': role,
      if (status != null) 'status': status,
      'page': page,
      'limit': limit,
    });
  }

  /// Get a single user by ID
  Future<Response> get(String id) async {
    return _dio.get('/users/$id');
  }

  /// Create a new user
  Future<Response> create(Map<String, dynamic> data) async {
    return _dio.post('/users', data: data);
  }

  /// Update a user
  Future<Response> update(String id, Map<String, dynamic> data) async {
    return _dio.put('/users/$id', data: data);
  }

  /// Delete a user
  Future<Response> delete(String id) async {
    return _dio.delete('/users/$id');
  }

  /// Ban a user
  Future<Response> ban(String id) async {
    return _dio.post('/users/$id/ban');
  }

  /// Unban a user
  Future<Response> unban(String id) async {
    return _dio.post('/users/$id/unban');
  }

  /// Reset user password
  Future<Response> resetPassword(String id) async {
    return _dio.post('/users/$id/reset-password');
  }

  /// Update user role
  Future<Response> updateRole(String id, String role) async {
    return _dio.put('/users/$id/role', data: {'role': role});
  }

  /// Export users to CSV
  Future<Response> export({String? status}) async {
    return _dio.get('/users/export', queryParameters: {
      if (status != null) 'status': status,
    });
  }
}