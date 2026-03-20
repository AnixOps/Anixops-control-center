import 'package:dio/dio.dart';

/// Playbook model
class Playbook {
  final String name;
  final String? displayName;
  final String? description;
  final String? category;
  final String? source;
  final String? content;
  final Map<String, dynamic>? variables;
  final List<String>? tags;
  final String? author;
  final String? version;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  Playbook({
    required this.name,
    this.displayName,
    this.description,
    this.category,
    this.source,
    this.content,
    this.variables,
    this.tags,
    this.author,
    this.version,
    this.createdAt,
    this.updatedAt,
  });

  factory Playbook.fromJson(Map<String, dynamic> json) {
    return Playbook(
      name: json['name'] as String,
      displayName: json['display_name'] as String?,
      description: json['description'] as String?,
      category: json['category'] as String?,
      source: json['source'] as String?,
      content: json['content'] as String?,
      variables: json['variables'] as Map<String, dynamic>?,
      tags: (json['tags'] as List<dynamic>?)?.cast<String>(),
      author: json['author'] as String?,
      version: json['version'] as String?,
      createdAt: json['created_at'] != null
          ? DateTime.parse(json['created_at'] as String)
          : null,
      updatedAt: json['updated_at'] != null
          ? DateTime.parse(json['updated_at'] as String)
          : null,
    );
  }

  String get title => displayName ?? name.replaceAll('-', ' ').split(' ').map((word) => word[0].toUpperCase() + word.substring(1)).join(' ');
}

/// Playbook category model
class PlaybookCategory {
  final String id;
  final String name;
  final String icon;
  final String description;

  PlaybookCategory({
    required this.id,
    required this.name,
    required this.icon,
    required this.description,
  });

  factory PlaybookCategory.fromJson(Map<String, dynamic> json) {
    return PlaybookCategory(
      id: json['id'] as String,
      name: json['name'] as String,
      icon: json['icon'] as String,
      description: json['description'] as String,
    );
  }
}

/// Playbooks API endpoints
class PlaybooksApi {
  final Dio _dio;

  PlaybooksApi(this._dio);

  /// Get all playbooks
  Future<List<Playbook>> getPlaybooks({String? category, String? source}) async {
    final queryParams = <String, String>{};
    if (category != null) queryParams['category'] = category;
    if (source != null) queryParams['source'] = source;

    final response = await _dio.get('/playbooks', queryParameters: queryParams);
    final data = response.data['data'] as List<dynamic>;
    return data.map((json) => Playbook.fromJson(json as Map<String, dynamic>)).toList();
  }

  /// Get built-in playbooks
  Future<List<Playbook>> getBuiltInPlaybooks() async {
    final response = await _dio.get('/playbooks/built-in');
    final data = response.data['data'] as List<dynamic>;
    return data.map((json) => Playbook.fromJson(json as Map<String, dynamic>)).toList();
  }

  /// Get playbook categories
  Future<List<PlaybookCategory>> getCategories() async {
    final response = await _dio.get('/playbooks/categories');
    final data = response.data['data'] as List<dynamic>;
    return data.map((json) => PlaybookCategory.fromJson(json as Map<String, dynamic>)).toList();
  }

  /// Get single playbook
  Future<Playbook> getPlaybook(String name) async {
    final response = await _dio.get('/playbooks/$name');
    return Playbook.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  /// Upload playbook
  Future<Playbook> uploadPlaybook({
    required String name,
    required String content,
    String? description,
    String? category,
    Map<String, dynamic>? variables,
    List<String>? tags,
  }) async {
    final response = await _dio.post('/playbooks', data: {
      'name': name,
      'content': content,
      if (description != null) 'description': description,
      if (category != null) 'category': category,
      if (variables != null) 'variables': variables,
      if (tags != null) 'tags': tags,
    });
    return Playbook.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  /// Delete playbook
  Future<void> deletePlaybook(String name) async {
    await _dio.delete('/playbooks/$name');
  }

  /// Sync built-in playbooks
  Future<List<String>> syncBuiltIn() async {
    final response = await _dio.post('/playbooks/sync-builtin');
    return (response.data['data'] as List<dynamic>).cast<String>();
  }
}