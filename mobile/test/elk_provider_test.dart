import 'package:flutter_test/flutter_test.dart';

void main() {
  group('Elasticsearch Index', () {
    test('creates index with valid name', () {
      final index = ElasticsearchIndex(
        name: 'logs-app-2026.03.23',
        health: 'green',
        docs: 500000,
        size: 1073741824,
      );
      expect(index.name, 'logs-app-2026.03.23');
      expect(index.health, 'green');
      expect(index.docs, 500000);
    });

    test('formats size in GB', () {
      final index = ElasticsearchIndex(name: 'test', health: 'green', docs: 100, size: 2147483648);
      expect(index.formattedSize, '2.0GB');
    });

    test('formats size in MB', () {
      final index = ElasticsearchIndex(name: 'test', health: 'green', docs: 100, size: 524288000);
      expect(index.formattedSize, '500.0MB');
    });

    test('checks if index is healthy', () {
      final greenIndex = ElasticsearchIndex(name: 'test', health: 'green', docs: 0, size: 0);
      final redIndex = ElasticsearchIndex(name: 'test', health: 'red', docs: 0, size: 0);
      expect(greenIndex.isHealthy, isTrue);
      expect(redIndex.isHealthy, isFalse);
    });
  });

  group('Cluster Health', () {
    test('creates cluster health', () {
      final health = ClusterHealth(
        clusterName: 'anixops-logs',
        status: 'green',
        numberOfNodes: 3,
        activeShards: 90,
        unassignedShards: 0,
      );
      expect(health.clusterName, 'anixops-logs');
      expect(health.status, 'green');
    });

    test('calculates shards per node', () {
      final health = ClusterHealth(clusterName: 'test', status: 'green', numberOfNodes: 3, activeShards: 90, unassignedShards: 0);
      expect(health.shardsPerNode, 30);
    });

    test('checks cluster status', () {
      final greenCluster = ClusterHealth(clusterName: 'test', status: 'green', numberOfNodes: 1, activeShards: 1, unassignedShards: 0);
      expect(greenCluster.isHealthy, isTrue);
    });
  });

  group('ILM Policy', () {
    test('creates policy with phases', () {
      final policy = ILMPolicy(name: 'logs-policy', phases: ['hot', 'warm', 'cold', 'delete']);
      expect(policy.name, 'logs-policy');
      expect(policy.phases.length, 4);
    });

    test('checks required phases', () {
      final policy = ILMPolicy(name: 'test', phases: ['hot', 'delete']);
      expect(policy.hasHotPhase, isTrue);
      expect(policy.hasDeletePhase, isTrue);
    });
  });

  group('Index Template', () {
    test('creates template with patterns', () {
      final template = IndexTemplate(name: 'logs-app', indexPatterns: ['logs-app-*'], numberOfShards: 3, numberOfReplicas: 1);
      expect(template.name, 'logs-app');
      expect(template.indexPatterns.length, 1);
    });

    test('matches index to template', () {
      final template = IndexTemplate(name: 'logs-app', indexPatterns: ['logs-app-*'], numberOfShards: 3, numberOfReplicas: 1);
      expect(template.matches('logs-app-2026.03.23'), isTrue);
      expect(template.matches('metrics-app-2026.03.23'), isFalse);
    });
  });

  group('Log Entry', () {
    test('creates log entry', () {
      final log = LogEntry(timestamp: DateTime.parse('2026-03-23T10:00:00Z'), message: 'Request processed', level: 'INFO', service: 'api-gateway');
      expect(log.message, 'Request processed');
      expect(log.level, 'INFO');
    });

    test('checks if error log', () {
      final errorLog = LogEntry(timestamp: DateTime.now(), message: 'Error', level: 'ERROR', service: 'test');
      final infoLog = LogEntry(timestamp: DateTime.now(), message: 'Info', level: 'INFO', service: 'test');
      expect(errorLog.isError, isTrue);
      expect(infoLog.isError, isFalse);
    });
  });

  group('Search Query', () {
    test('builds term query', () {
      final query = SearchQuery.term('level', 'ERROR');
      expect(query['query']['term']['level'], 'ERROR');
    });

    test('builds range query', () {
      final query = SearchQuery.range('@timestamp', 'now-1h', null);
      expect(query['query']['range']['@timestamp']['gte'], 'now-1h');
    });
  });
}

class ElasticsearchIndex {
  final String name;
  final String health;
  final int docs;
  final int size;

  ElasticsearchIndex({required this.name, required this.health, required this.docs, required this.size});

  bool get isHealthy => health == 'green';

  String get formattedSize {
    if (size >= 1073741824) return '${(size / 1073741824).toStringAsFixed(1)}GB';
    if (size >= 1048576) return '${(size / 1048576).toStringAsFixed(1)}MB';
    if (size >= 1024) return '${(size / 1024).toStringAsFixed(1)}KB';
    return '${size}B';
  }
}

class ClusterHealth {
  final String clusterName;
  final String status;
  final int numberOfNodes;
  final int activeShards;
  final int unassignedShards;

  ClusterHealth({required this.clusterName, required this.status, required this.numberOfNodes, required this.activeShards, required this.unassignedShards});

  bool get isHealthy => status == 'green';
  int get shardsPerNode => activeShards ~/ numberOfNodes;
}

class ILMPolicy {
  final String name;
  final List<String> phases;

  ILMPolicy({required this.name, required this.phases});

  bool get hasHotPhase => phases.contains('hot');
  bool get hasDeletePhase => phases.contains('delete');
}

class IndexTemplate {
  final String name;
  final List<String> indexPatterns;
  final int numberOfShards;
  final int numberOfReplicas;

  IndexTemplate({required this.name, required this.indexPatterns, required this.numberOfShards, required this.numberOfReplicas});

  bool matches(String indexName) {
    for (final pattern in indexPatterns) {
      final regex = pattern.replaceAll('*', '.*');
      if (RegExp('^$regex\$').hasMatch(indexName)) return true;
    }
    return false;
  }
}

class LogEntry {
  final DateTime timestamp;
  final String message;
  final String level;
  final String service;

  LogEntry({required this.timestamp, required this.message, required this.level, required this.service});

  bool get isError => level == 'ERROR';
}

class SearchQuery {
  static Map<String, dynamic> term(String field, String value) {
    return {'query': {'term': {field: value}}};
  }

  static Map<String, dynamic> range(String field, String? gte, String? lte) {
    final range = <String, dynamic>{};
    if (gte != null) range['gte'] = gte;
    if (lte != null) range['lte'] = lte;
    return {'query': {'range': {field: range}}};
  }
}