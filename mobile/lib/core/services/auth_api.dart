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

  /// Update current user profile
  Future<Response> updateProfile({
    String? name,
    String? email,
  }) async {
    return _dio.put('/users/me', data: {
      if (name != null) 'name': name,
      if (email != null) 'email': email,
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

/// API Tokens API endpoints
class TokensApi {
  final Dio _dio;

  TokensApi(this._dio);

  /// List API tokens
  Future<Response> list() async {
    return _dio.get('/users/me/tokens');
  }

  /// Create API token
  Future<Response> create(String name, {int? expiresInDays}) async {
    return _dio.post('/users/me/tokens', data: {
      'name': name,
      if (expiresInDays != null) 'expires_in_days': expiresInDays,
    });
  }

  /// Delete API token
  Future<Response> delete(String id) async {
    return _dio.delete('/users/me/tokens/$id');
  }
}

/// Sessions API endpoints
class SessionsApi {
  final Dio _dio;

  SessionsApi(this._dio);

  /// List active sessions
  Future<Response> list() async {
    return _dio.get('/users/me/sessions');
  }

  /// Delete other sessions
  Future<Response> deleteOthers() async {
    return _dio.delete('/users/me/sessions/others');
  }
}