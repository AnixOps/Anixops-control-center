import 'package:dio/dio.dart';

/// AI Models configuration
const aiModels = {
  'textGeneration': '@cf/meta/llama-3.1-8b-instruct',
  'textEmbeddings': '@cf/baai/bge-base-en-v1.5',
};

/// AI API service for Workers AI integration
class AiApi {
  final Dio _dio;

  AiApi(this._dio);

  /// Send chat message to AI assistant
  Future<Map<String, dynamic>> chat(String message, {List<Map<String, String>>? history}) async {
    final response = await _dio.post('/ai/chat', data: {
      'message': message,
      'history': history ?? [],
    });
    return response.data;
  }

  /// Analyze log content with AI
  Future<Map<String, dynamic>> analyzeLog(String logContent) async {
    final response = await _dio.post('/ai/analyze-log', data: {
      'logContent': logContent,
    });
    return response.data;
  }

  /// Get operations advice
  Future<Map<String, dynamic>> getOpsAdvice(Map<String, dynamic> context) async {
    final response = await _dio.post('/ai/ops-advice', data: context);
    return response.data;
  }

  /// Generate text embedding
  Future<List<double>> generateEmbedding(String text) async {
    final response = await _dio.post('/ai/embedding', data: {'text': text});
    return List<double>.from(response.data['embedding'] ?? []);
  }

  /// Natural language query
  Future<Map<String, dynamic>> query(String query) async {
    final response = await _dio.post('/ai/query', data: {'query': query});
    return response.data;
  }

  /// Semantic search using Vectorize
  Future<List<Map<String, dynamic>>> semanticSearch(String query) async {
    // First get embedding
    final embedding = await generateEmbedding(query);
    if (embedding.isEmpty) return [];

    // Then search vectors
    final response = await _dio.post('/vectors/search', data: {
      'embedding': embedding,
      'topK': 10,
    });
    return List<Map<String, dynamic>>.from(response.data['results'] ?? []);
  }
}