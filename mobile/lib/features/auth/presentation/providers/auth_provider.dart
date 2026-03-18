import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:dio/dio.dart';
import 'package:anixops_mobile/core/services/api_client.dart';

class AuthState {
  final bool isAuthenticated;
  final bool isLoading;
  final String? token;
  final String? refreshToken;
  final int? userId;
  final String? email;
  final String? role;
  final String? error;

  const AuthState({
    this.isAuthenticated = false,
    this.isLoading = false,
    this.token,
    this.refreshToken,
    this.userId,
    this.email,
    this.role,
    this.error,
  });

  AuthState copyWith({
    bool? isAuthenticated,
    bool? isLoading,
    String? token,
    String? refreshToken,
    int? userId,
    String? email,
    String? role,
    String? error,
  }) {
    return AuthState(
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      isLoading: isLoading ?? this.isLoading,
      token: token ?? this.token,
      refreshToken: refreshToken ?? this.refreshToken,
      userId: userId ?? this.userId,
      email: email ?? this.email,
      role: role ?? this.role,
      error: error,
    );
  }
}

class AuthNotifier extends StateNotifier<AuthState> {
  final ApiClient _apiClient;
  final SharedPreferences _prefs;

  static const String _tokenKey = 'auth_token';
  static const String _refreshTokenKey = 'refresh_token';
  static const String _userIdKey = 'user_id';
  static const String _emailKey = 'user_email';
  static const String _roleKey = 'user_role';

  AuthNotifier(this._apiClient, this._prefs) : super(const AuthState()) {
    _loadStoredAuth();
  }

  Future<void> _loadStoredAuth() async {
    final token = _prefs.getString(_tokenKey);
    if (token != null) {
      _apiClient.setAuthToken(token);
      state = AuthState(
        isAuthenticated: true,
        token: token,
        refreshToken: _prefs.getString(_refreshTokenKey),
        userId: _prefs.getInt(_userIdKey),
        email: _prefs.getString(_emailKey),
        role: _prefs.getString(_roleKey),
      );
    }
  }

  Future<bool> login(String email, String password) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final response = await _apiClient.auth.login(email, password);

      if (response.data['success'] == true) {
        final data = response.data['data'];
        final token = data['access_token'] as String;
        final refreshToken = data['refresh_token'] as String;
        final user = data['user'] as Map<String, dynamic>;

        // Store tokens
        await _prefs.setString(_tokenKey, token);
        await _prefs.setString(_refreshTokenKey, refreshToken);
        await _prefs.setInt(_userIdKey, user['id'] as int);
        await _prefs.setString(_emailKey, user['email'] as String);
        await _prefs.setString(_roleKey, user['role'] as String);

        // Set auth header
        _apiClient.setAuthToken(token);

        state = AuthState(
          isAuthenticated: true,
          isLoading: false,
          token: token,
          refreshToken: refreshToken,
          userId: user['id'] as int,
          email: user['email'] as String,
          role: user['role'] as String,
        );

        return true;
      } else {
        state = state.copyWith(
          isLoading: false,
          error: response.data['error'] ?? 'Login failed',
        );
        return false;
      }
    } on DioException catch (e) {
      final error = e.response?.data?['error'] ?? 'Network error';
      state = state.copyWith(isLoading: false, error: error);
      return false;
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
      return false;
    }
  }

  Future<void> logout() async {
    try {
      await _apiClient.auth.logout();
    } catch (_) {
      // Ignore logout API errors
    }

    // Clear stored tokens
    await _prefs.remove(_tokenKey);
    await _prefs.remove(_refreshTokenKey);
    await _prefs.remove(_userIdKey);
    await _prefs.remove(_emailKey);
    await _prefs.remove(_roleKey);

    // Clear API client auth
    _apiClient.clearAuthToken();

    state = const AuthState();
  }

  Future<bool> refreshAuthToken() async {
    final storedRefreshToken = _prefs.getString(_refreshTokenKey);
    if (storedRefreshToken == null) return false;

    try {
      final response = await _apiClient.auth.refresh(storedRefreshToken);

      if (response.data['success'] == true) {
        final data = response.data['data'];
        final token = data['access_token'] as String;

        await _prefs.setString(_tokenKey, token);
        _apiClient.setAuthToken(token);

        state = state.copyWith(token: token);
        return true;
      }
    } catch (_) {}

    // Refresh failed, logout
    await logout();
    return false;
  }

  Future<bool> checkAuth() async {
    if (!state.isAuthenticated) return false;

    try {
      final response = await _apiClient.auth.me();
      return response.data['success'] == true;
    } catch (_) {
      return false;
    }
  }
}

// Provider for SharedPreferences
final sharedPreferencesProvider = Provider<SharedPreferences>((ref) {
  throw UnimplementedError('SharedPreferences must be overridden');
});

// Provider for ApiClient
final apiClientProvider = Provider<ApiClient>((ref) {
  return ApiClient();
});

// Auth state provider
final authStateProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  final prefs = ref.watch(sharedPreferencesProvider);
  return AuthNotifier(apiClient, prefs);
});