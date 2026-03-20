import 'package:dio/dio.dart';

/// Nodes API endpoints
class NodesApi {
  final Dio _dio;

  NodesApi(this._dio);

  /// List all nodes with optional filters
  Future<Response> list({
    String? search,
    String? status,
    String? type,
    int page = 1,
    int limit = 20,
  }) async {
    return _dio.get('/nodes', queryParameters: {
      if (search != null) 'search': search,
      if (status != null) 'status': status,
      if (type != null) 'type': type,
      'page': page,
      'limit': limit,
    });
  }

  /// Get a single node by ID
  Future<Response> get(String id) async {
    return _dio.get('/nodes/$id');
  }

  /// Create a new node
  Future<Response> create(Map<String, dynamic> data) async {
    return _dio.post('/nodes', data: data);
  }

  /// Update a node
  Future<Response> update(String id, Map<String, dynamic> data) async {
    return _dio.put('/nodes/$id', data: data);
  }

  /// Delete a node
  Future<Response> delete(String id) async {
    return _dio.delete('/nodes/$id');
  }

  /// Start a node
  Future<Response> start(String id) async {
    return _dio.post('/nodes/$id/start');
  }

  /// Stop a node
  Future<Response> stop(String id) async {
    return _dio.post('/nodes/$id/stop');
  }

  /// Restart a node
  Future<Response> restart(String id) async {
    return _dio.post('/nodes/$id/restart');
  }

  /// Get node statistics
  Future<Response> stats(String id) async {
    return _dio.get('/nodes/$id/stats');
  }

  /// Get node logs
  Future<Response> logs(String id, {int limit = 100, String? level}) async {
    return _dio.get('/nodes/$id/logs', queryParameters: {
      'limit': limit,
      if (level != null) 'level': level,
    });
  }

  /// Test node connection
  Future<Response> testConnection(String id) async {
    return _dio.post('/nodes/$id/test');
  }

  /// Sync node configuration
  Future<Response> sync(String id) async {
    return _dio.post('/nodes/$id/sync');
  }

  /// Bulk operations on nodes
  Future<Response> bulkAction({
    required List<String> nodeIds,
    required String action, // start, stop, restart, delete
  }) async {
    return _dio.post('/nodes/bulk', data: {
      'node_ids': nodeIds,
      'action': action,
    });
  }
}