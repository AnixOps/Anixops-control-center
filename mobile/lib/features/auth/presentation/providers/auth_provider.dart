import 'package:flutter_riverpod/flutter_riverpod.dart';

class AuthState {
  final bool isAuthenticated;
  final String? token;
  final String? userId;
  final String? email;
  final String? role;

  const AuthState({
    this.isAuthenticated = false,
    this.token,
    this.userId,
    this.email,
    this.role,
  });

  AuthState copyWith({
    bool? isAuthenticated,
    String? token,
    String? userId,
    String? email,
    String? role,
  }) {
    return AuthState(
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      token: token ?? this.token,
      userId: userId ?? this.userId,
      email: email ?? this.email,
      role: role ?? this.role,
    );
  }
}

class AuthNotifier extends StateNotifier<AuthState> {
  AuthNotifier() : super(const AuthState());

  Future<bool> login(String email, String password) async {
    // TODO: Implement actual login
    state = AuthState(
      isAuthenticated: true,
      token: 'demo-token',
      userId: '1',
      email: email,
      role: 'admin',
    );
    return true;
  }

  Future<void> logout() async {
    // TODO: Implement logout
    state = const AuthState();
  }

  Future<bool> refreshToken() async {
    // TODO: Implement token refresh
    return true;
  }

  Future<bool> checkAuth() async {
    // TODO: Check if token is valid
    return state.isAuthenticated;
  }
}

final authStateProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier();
});