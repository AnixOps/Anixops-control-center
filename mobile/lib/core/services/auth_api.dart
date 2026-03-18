import 'package:dio/dio.dart';

/// Authentication API endpoints
class AuthApi {
  final Dio _dio;

  AuthApi(this._dio);

  /// Login with email and password
  Future<Response> login(String email, String password) async {
    return _dio.post('/auth/login', data: {
      'email': email,
      'password': password,
    });
  }

  /// Register new user
  Future<Response> register(String email, String password, {String role = 'viewer'}) async {
    return _dio.post('/auth/register', data: {
      'email': email,
      'password': password,
      'role': role,
    });
  }

  /// Logout current user
  Future<Response> logout() async {
    return _dio.post('/auth/logout');
  }

  /// Refresh authentication token
  Future<Response> refresh(String refreshToken) async {
    return _dio.post('/auth/refresh', data: {
      'refresh_token': refreshToken,
    });
  }

  /// Get current user info
  Future<Response> me() async {
    return _dio.get('/users/me');
  }

  /// Update password
  Future<Response> updatePassword({
    required String currentPassword,
    required String newPassword,
  }) async {
    return _dio.put('/auth/password', data: {
      'current_password': currentPassword,
      'new_password': newPassword,
    });
  }

  /// Request password reset email
  Future<Response> forgotPassword(String email) async {
    return _dio.post('/auth/forgot-password', data: {
      'email': email,
    });
  }

  /// Reset password with token
  Future<Response> resetPassword({
    required String token,
    required String password,
  }) async {
    return _dio.post('/auth/reset-password', data: {
      'token': token,
      'password': password,
    });
  }
}