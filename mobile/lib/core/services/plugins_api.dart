import 'package:dio/dio.dart';

/// Plugins API endpoints
class PluginsApi {
  final Dio _dio;

  PluginsApi(this._dio);

  /// List all plugins
  Future<Response> list() async {
    return _dio.get('/plugins');
  }

  /// Get a single plugin by name
  Future<Response> get(String name) async {
    return _dio.get('/plugins/$name');
  }

  /// Execute a plugin action
  Future<Response> execute(String name, String action, {Map<String, dynamic>? params}) async {
    return _dio.post('/plugins/$name/execute', data: {
      'action': action,
      'params': params ?? {},
    });
  }

  /// Get plugin status
  Future<Response> status(String name) async {
    return _dio.get('/plugins/$name/status');
  }

  /// Get plugin configuration
  Future<Response> config(String name) async {
    return _dio.get('/plugins/$name/config');
  }

  /// Update plugin configuration
  Future<Response> updateConfig(String name, Map<String, dynamic> config) async {
    return _dio.put('/plugins/$name/config', data: config);
  }

  /// Enable a plugin
  Future<Response> enable(String name) async {
    return _dio.post('/plugins/$name/enable');
  }

  /// Disable a plugin
  Future<Response> disable(String name) async {
    return _dio.post('/plugins/$name/disable');
  }

  /// Start a plugin
  Future<Response> start(String name) async {
    return _dio.post('/plugins/$name/start');
  }

  /// Stop a plugin
  Future<Response> stop(String name) async {
    return _dio.post('/plugins/$name/stop');
  }

  /// Restart a plugin
  Future<Response> restart(String name) async {
    return _dio.post('/plugins/$name/restart');
  }

  /// Get plugin logs
  Future<Response> logs(String name, {int limit = 100}) async {
    return _dio.get('/plugins/$name/logs', queryParameters: {'limit': limit});
  }
}