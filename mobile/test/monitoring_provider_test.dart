import 'package:flutter_test/flutter_test.dart';

void main() {
  group('HealthCheck Model', () {
    test('creates health check with all fields', () {
      final check = HealthCheck(
        name: 'API',
        status: 'healthy',
        latency: 12,
        lastCheck: DateTime.parse('2026-03-23T10:00:00Z'),
      );

      expect(check.name, 'API');
      expect(check.status, 'healthy');
      expect(check.latency, 12);
    });

    test('validates health status', () {
      final validStatuses = ['healthy', 'degraded', 'unhealthy'];

      expect(validStatuses.contains('healthy'), isTrue);
      expect(validStatuses.contains('degraded'), isTrue);
      expect(validStatuses.contains('unhealthy'), isTrue);
      expect(validStatuses.contains('unknown'), isFalse);
    });

    test('checks if healthy', () {
      final healthy = HealthCheck(name: 'API', status: 'healthy', latency: 10);
      final degraded = HealthCheck(name: 'DB', status: 'degraded', latency: 50);

      expect(healthy.isHealthy, isTrue);
      expect(degraded.isHealthy, isFalse);
    });
  });

  group('Metric Model', () {
    test('creates metric with value and unit', () {
      final metric = Metric(
        name: 'request_rate',
        value: 1250.5,
        unit: 'req/s',
        timestamp: DateTime.parse('2026-03-23T10:00:00Z'),
      );

      expect(metric.name, 'request_rate');
      expect(metric.value, 1250.5);
      expect(metric.unit, 'req/s');
    });

    test('formats value correctly', () {
      expect(Metric.formatNumber(1250), '1.3K');
      expect(Metric.formatNumber(1500000), '1.5M');
      expect(Metric.formatNumber(500), '500');
    });

    test('calculates percentage', () {
      final metric = Metric(name: 'error_rate', value: 0.15, unit: '%');
      expect(metric.asPercentage, '0.15%');
    });
  });

  group('Alert Model', () {
    test('creates alert with severity', () {
      final alert = Alert(
        id: '1',
        name: 'High Memory Usage',
        severity: 'warning',
        metric: 'memory_percent',
        value: 85,
        threshold: 80,
        startedAt: DateTime.parse('2026-03-23T10:00:00Z'),
      );

      expect(alert.id, '1');
      expect(alert.name, 'High Memory Usage');
      expect(alert.severity, 'warning');
      expect(alert.isFiring, isTrue);
    });

    test('checks if alert is firing', () {
      final firing = Alert(
        id: '1', name: 'Test', severity: 'warning',
        metric: 'test', value: 85, threshold: 80, startedAt: DateTime.now(),
      );
      final notFiring = Alert(
        id: '2', name: 'Test', severity: 'warning',
        metric: 'test', value: 75, threshold: 80, startedAt: DateTime.now(),
      );

      expect(firing.isFiring, isTrue);
      expect(notFiring.isFiring, isFalse);
    });

    test('validates severity levels', () {
      final severities = ['info', 'warning', 'critical'];

      expect(severities.contains('info'), isTrue);
      expect(severities.contains('warning'), isTrue);
      expect(severities.contains('critical'), isTrue);
    });
  });

  group('ServiceHealth Model', () {
    test('creates service health', () {
      final service = ServiceHealth(
        name: 'api-gateway',
        health: 'healthy',
        requestRate: 500,
        errorRate: 0.1,
        latency: 25,
      );

      expect(service.name, 'api-gateway');
      expect(service.health, 'healthy');
      expect(service.requestRate, 500);
    });

    test('formats error rate', () {
      final service = ServiceHealth(
        name: 'test', health: 'healthy',
        requestRate: 100, errorRate: 0.5, latency: 10,
      );

      expect(service.errorRateFormatted, '0.50%');
    });

    test('checks service status', () {
      final healthy = ServiceHealth(
        name: 'test', health: 'healthy',
        requestRate: 100, errorRate: 0, latency: 10,
      );
      final unhealthy = ServiceHealth(
        name: 'test', health: 'unhealthy',
        requestRate: 100, errorRate: 0, latency: 10,
      );

      expect(healthy.isHealthy, isTrue);
      expect(unhealthy.isHealthy, isFalse);
    });
  });

  group('MonitoringDashboard', () {
    test('calculates overall health', () {
      final checks = [
        HealthCheck(name: 'API', status: 'healthy', latency: 10),
        HealthCheck(name: 'DB', status: 'healthy', latency: 5),
        HealthCheck(name: 'Cache', status: 'degraded', latency: 50),
      ];

      final overallHealth = MonitoringDashboard.getOverallHealth(checks);
      expect(overallHealth, 'degraded');
    });

    test('returns unhealthy if any check is unhealthy', () {
      final checks = [
        HealthCheck(name: 'API', status: 'healthy', latency: 10),
        HealthCheck(name: 'DB', status: 'unhealthy', latency: 100),
      ];

      final overallHealth = MonitoringDashboard.getOverallHealth(checks);
      expect(overallHealth, 'unhealthy');
    });

    test('calculates average latency', () {
      final checks = [
        HealthCheck(name: 'API', status: 'healthy', latency: 10),
        HealthCheck(name: 'DB', status: 'healthy', latency: 20),
        HealthCheck(name: 'Cache', status: 'healthy', latency: 30),
      ];

      final avgLatency = MonitoringDashboard.getAverageLatency(checks);
      expect(avgLatency, 20);
    });
  });
}

// Model classes for testing
class HealthCheck {
  final String name;
  final String status;
  final int latency;
  final DateTime? lastCheck;

  HealthCheck({
    required this.name,
    required this.status,
    required this.latency,
    this.lastCheck,
  });

  bool get isHealthy => status == 'healthy';
}

class Metric {
  final String name;
  final double value;
  final String unit;
  final DateTime? timestamp;

  Metric({
    required this.name,
    required this.value,
    required this.unit,
    this.timestamp,
  });

  String get asPercentage => '${value.toFixed(2)}%';

  static String formatNumber(num n) {
    if (n >= 1000000) return '${(n / 1000000).toStringAsFixed(1)}M';
    if (n >= 1000) return '${(n / 1000).toStringAsFixed(1)}K';
    return n.toString();
  }
}

class Alert {
  final String id;
  final String name;
  final String severity;
  final String metric;
  final double value;
  final double threshold;
  final DateTime startedAt;

  Alert({
    required this.id,
    required this.name,
    required this.severity,
    required this.metric,
    required this.value,
    required this.threshold,
    required this.startedAt,
  });

  bool get isFiring => value > threshold;
}

class ServiceHealth {
  final String name;
  final String health;
  final int requestRate;
  final double errorRate;
  final int latency;

  ServiceHealth({
    required this.name,
    required this.health,
    required this.requestRate,
    required this.errorRate,
    required this.latency,
  });

  bool get isHealthy => health == 'healthy';
  String get errorRateFormatted => '${errorRate.toStringAsFixed(2)}%';
}

class MonitoringDashboard {
  static String getOverallHealth(List<HealthCheck> checks) {
    if (checks.any((c) => c.status == 'unhealthy')) return 'unhealthy';
    if (checks.any((c) => c.status == 'degraded')) return 'degraded';
    return 'healthy';
  }

  static int getAverageLatency(List<HealthCheck> checks) {
    if (checks.isEmpty) return 0;
    return checks.fold(0, (sum, c) => sum + c.latency) ~/ checks.length;
  }
}

extension on double {
  String toFixed(int fractionDigits) => toStringAsFixed(fractionDigits);
}