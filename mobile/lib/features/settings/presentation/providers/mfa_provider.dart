import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/services/mfa_api.dart';
import '../../../../core/providers/api_providers.dart';

/// MFA state
class MFAState {
  final MFAStatus? status;
  final MFASetupResult? setupResult;
  final bool isLoading;
  final String? error;

  const MFAState({
    this.status,
    this.setupResult,
    this.isLoading = false,
    this.error,
  });

  MFAState copyWith({
    MFAStatus? status,
    MFASetupResult? setupResult,
    bool? isLoading,
    String? error,
  }) {
    return MFAState(
      status: status ?? this.status,
      setupResult: setupResult ?? this.setupResult,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
    );
  }
}

/// MFA notifier
class MFANotifier extends StateNotifier<MFAState> {
  final MFAApi _api;

  MFANotifier(this._api) : super(const MFAState()) {
    loadStatus();
  }

  /// Load MFA status
  Future<void> loadStatus() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final status = await _api.getStatus();
      state = state.copyWith(
        status: status,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// Start MFA setup
  Future<MFASetupResult?> setup() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final result = await _api.setup();
      state = state.copyWith(
        setupResult: result,
        isLoading: false,
      );
      return result;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      return null;
    }
  }

  /// Enable MFA with verification code
  Future<bool> enable(String code) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final success = await _api.enable(code);
      if (success) {
        await loadStatus();
        state = state.copyWith(setupResult: null);
      } else {
        state = state.copyWith(
          isLoading: false,
          error: 'Invalid verification code',
        );
      }
      return success;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      return false;
    }
  }

  /// Disable MFA with verification code
  Future<bool> disable(String code) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final success = await _api.disable(code);
      if (success) {
        await loadStatus();
      } else {
        state = state.copyWith(
          isLoading: false,
          error: 'Invalid verification code',
        );
      }
      return success;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      return false;
    }
  }

  /// Regenerate recovery codes
  Future<List<String>?> regenerateRecoveryCodes(String code) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final codes = await _api.regenerateRecoveryCodes(code);
      state = state.copyWith(isLoading: false);
      return codes;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      return null;
    }
  }

  /// Clear error
  void clearError() {
    state = state.copyWith(error: null);
  }

  /// Clear setup result
  void clearSetupResult() {
    state = state.copyWith(setupResult: null);
  }
}

/// Provider for MFAState
final mfaProvider = StateNotifierProvider<MFANotifier, MFAState>((ref) {
  final client = ref.watch(apiClientProvider);
  return MFANotifier(client.mfa);
});