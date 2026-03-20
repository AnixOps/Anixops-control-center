import 'package:dio/dio.dart';

/// MFA status model
class MFAStatus {
  final bool enabled;
  final bool hasRecoveryCodes;
  final DateTime? enabledAt;
  final String? method;

  MFAStatus({
    required this.enabled,
    this.hasRecoveryCodes = false,
    this.enabledAt,
    this.method,
  });

  factory MFAStatus.fromJson(Map<String, dynamic> json) {
    return MFAStatus(
      enabled: json['enabled'] as bool? ?? false,
      hasRecoveryCodes: json['has_recovery_codes'] as bool? ?? false,
      enabledAt: json['enabled_at'] != null
          ? DateTime.tryParse(json['enabled_at'] as String)
          : null,
      method: json['method'] as String?,
    );
  }
}

/// MFA setup result model
class MFASetupResult {
  final String secret;
  final String otpauthUrl;
  final List<String> recoveryCodes;

  MFASetupResult({
    required this.secret,
    required this.otpauthUrl,
    required this.recoveryCodes,
  });

  factory MFASetupResult.fromJson(Map<String, dynamic> json) {
    return MFASetupResult(
      secret: json['secret'] as String? ?? '',
      otpauthUrl: json['otpauth_url'] as String? ?? '',
      recoveryCodes: (json['recovery_codes'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList() ?? [],
    );
  }
}

/// MFA API endpoints
class MFAApi {
  final Dio _dio;

  MFAApi(this._dio);

  /// Get MFA status
  Future<MFAStatus> getStatus() async {
    final response = await _dio.get('/mfa/status');
    return MFAStatus.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  /// Setup MFA (generate secret and recovery codes)
  Future<MFASetupResult> setup() async {
    final response = await _dio.post('/mfa/setup');
    return MFASetupResult.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  /// Enable MFA (verify setup and enable)
  Future<bool> enable(String code) async {
    try {
      final response = await _dio.post('/mfa/enable', data: {'code': code});
      return response.data['success'] as bool? ?? false;
    } catch (e) {
      return false;
    }
  }

  /// Disable MFA
  Future<bool> disable(String code) async {
    try {
      final response = await _dio.post('/mfa/disable', data: {'code': code});
      return response.data['success'] as bool? ?? false;
    } catch (e) {
      return false;
    }
  }

  /// Verify MFA code
  Future<bool> verify(String code) async {
    try {
      final response = await _dio.post('/mfa/verify', data: {'code': code});
      return response.data['success'] as bool? ?? false;
    } catch (e) {
      return false;
    }
  }

  /// Regenerate recovery codes
  Future<List<String>?> regenerateRecoveryCodes(String code) async {
    try {
      final response = await _dio.post('/mfa/recovery-codes/regenerate', data: {'code': code});
      final data = response.data['data'] as Map<String, dynamic>?;
      return (data?['recovery_codes'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList();
    } catch (e) {
      return null;
    }
  }
}