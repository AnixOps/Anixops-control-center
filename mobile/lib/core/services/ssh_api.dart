import 'package:dio/dio.dart';

/// SSH API endpoints for server import
class SshApi {
  final Dio _dio;

  SshApi(this._dio);

  /// Test SSH connection with password
  Future<Response> testConnectionWithPassword({
    required String host,
    required int port,
    required String username,
    required String password,
  }) async {
    return _dio.post('/ssh/test', data: {
      'host': host,
      'port': port,
      'username': username,
      'auth_type': 'password',
      'password': password,
    });
  }

  /// Test SSH connection with private key
  Future<Response> testConnectionWithKey({
    required String host,
    required int port,
    required String username,
    required String privateKey,
    String? passphrase,
  }) async {
    return _dio.post('/ssh/test', data: {
      'host': host,
      'port': port,
      'username': username,
      'auth_type': 'key',
      'private_key': privateKey,
      'passphrase': passphrase,
    });
  }

  /// Import server via SSH
  Future<Response> importServer({
    required String host,
    required int port,
    required String username,
    required String authType,
    String? password,
    String? privateKey,
    String? passphrase,
    String? name,
    String? group,
    List<String>? tags,
  }) async {
    return _dio.post('/ssh/import', data: {
      'host': host,
      'port': port,
      'username': username,
      'auth_type': authType,
      if (password != null) 'password': password,
      if (privateKey != null) 'private_key': privateKey,
      if (passphrase != null) 'passphrase': passphrase,
      if (name != null) 'name': name,
      if (group != null) 'group': group,
      if (tags != null) 'tags': tags,
    });
  }

  /// Detect server type via SSH
  Future<Response> detectServerType({
    required String host,
    required int port,
    required String username,
    required String authType,
    String? password,
    String? privateKey,
    String? passphrase,
  }) async {
    return _dio.post('/ssh/detect', data: {
      'host': host,
      'port': port,
      'username': username,
      'auth_type': authType,
      if (password != null) 'password': password,
      if (privateKey != null) 'private_key': privateKey,
      if (passphrase != null) 'passphrase': passphrase,
    });
  }

  /// Execute command via SSH
  Future<Response> executeCommand({
    required String nodeId,
    required String command,
  }) async {
    return _dio.post('/nodes/$nodeId/execute', data: {
      'command': command,
    });
  }
}