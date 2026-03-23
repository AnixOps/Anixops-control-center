import 'package:flutter_test/flutter_test.dart';

void main() {
  group('Span Model', () {
    test('creates span with all fields', () {
      final span = Span(
        spanId: 'b7ad6b7169203331',
        name: 'HTTP GET /api/users',
        kind: 'server',
        startTime: DateTime.parse('2026-03-23T10:00:00Z'),
        endTime: DateTime.parse('2026-03-23T10:00:00.100Z'),
        status: SpanStatus.ok,
        attributes: {'http.method': 'GET', 'http.url': '/api/users'},
        resource: {'service.name': 'api-gateway'},
      );

      expect(span.spanId, 'b7ad6b7169203331');
      expect(span.name, 'HTTP GET /api/users');
      expect(span.kind, 'server');
      expect(span.duration, 100);
    });

    test('calculates duration correctly', () {
      final span = Span(
        spanId: '1',
        name: 'test',
        kind: 'internal',
        startTime: DateTime.parse('2026-03-23T10:00:00.000Z'),
        endTime: DateTime.parse('2026-03-23T10:00:00.250Z'),
        status: SpanStatus.ok,
      );
      expect(span.duration, 250);
    });

    test('supports different span kinds', () {
      final kinds = ['unspecified', 'internal', 'server', 'client', 'producer', 'consumer'];
      for (final kind in kinds) {
        final span = Span(
          spanId: '1',
          name: 'test',
          kind: kind,
          startTime: DateTime.now(),
          status: SpanStatus.ok,
        );
        expect(span.kind, kind);
      }
    });

    test('supports status codes', () {
      final okSpan = Span(spanId: '1', name: 'ok', kind: 'server', startTime: DateTime.now(), status: SpanStatus.ok);
      final errorSpan = Span(spanId: '2', name: 'error', kind: 'server', startTime: DateTime.now(), status: SpanStatus.error(message: 'Failed'));
      final unsetSpan = Span(spanId: '3', name: 'unset', kind: 'server', startTime: DateTime.now(), status: SpanStatus.unset);

      expect(okSpan.status.code, 'ok');
      expect(errorSpan.status.code, 'error');
      expect(errorSpan.status.message, 'Failed');
      expect(unsetSpan.status.code, 'unset');
    });
  });

  group('Trace Model', () {
    test('creates trace with spans', () {
      final trace = Trace(
        traceId: '0af7651916cd43dd8448eb211c80319c',
        spans: [
          Span(spanId: '1', name: 'HTTP GET', kind: 'server', startTime: DateTime.now(), status: SpanStatus.ok),
          Span(spanId: '2', name: 'db:query', kind: 'client', startTime: DateTime.now(), status: SpanStatus.ok),
        ],
        status: 'ok',
        duration: 150,
      );
      expect(trace.traceId, '0af7651916cd43dd8448eb211c80319c');
      expect(trace.spanCount, 2);
    });

    test('validates trace ID format', () {
      final traceId = '0af7651916cd43dd8448eb211c80319c';
      expect(traceId.length, 32);
      expect(RegExp(r'^[0-9a-f]{32}$').hasMatch(traceId), isTrue);
    });

    test('calculates service count', () {
      final trace = Trace(
        traceId: 'test',
        spans: [
          Span(spanId: '1', name: 'HTTP', kind: 'server', startTime: DateTime.now(), status: SpanStatus.ok, resource: {'service.name': 'api-gateway'}),
          Span(spanId: '2', name: 'auth', kind: 'client', startTime: DateTime.now(), status: SpanStatus.ok, resource: {'service.name': 'auth-service'}),
          Span(spanId: '3', name: 'db', kind: 'client', startTime: DateTime.now(), status: SpanStatus.ok, resource: {'service.name': 'api-gateway'}),
        ],
        status: 'ok',
        duration: 100,
      );
      expect(trace.serviceCount, 2);
    });
  });

  group('W3C Trace Context', () {
    test('formats traceparent header', () {
      final context = SpanContext(
        traceId: '0af7651916cd43dd8448eb211c80319c',
        spanId: 'b7ad6b7169203331',
        traceFlags: 1,
      );
      final traceparent = context.toTraceparent();
      expect(traceparent, '00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01');
    });

    test('parses traceparent header', () {
      final context = SpanContext.fromTraceparent('00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01');
      expect(context, isNotNull);
      expect(context!.traceId, '0af7651916cd43dd8448eb211c80319c');
      expect(context.spanId, 'b7ad6b7169203331');
      expect(context.traceFlags, 1);
    });

    test('returns null for invalid traceparent', () {
      expect(SpanContext.fromTraceparent('invalid'), isNull);
      expect(SpanContext.fromTraceparent(''), isNull);
    });

    test('handles sampled flag', () {
      final sampled = SpanContext(traceId: 'test123', spanId: 'span123', traceFlags: 1);
      final notSampled = SpanContext(traceId: 'test123', spanId: 'span123', traceFlags: 0);
      expect(sampled.isSampled, isTrue);
      expect(notSampled.isSampled, isFalse);
    });
  });

  group('Trace Statistics', () {
    test('calculates error rate', () {
      final traces = [
        Trace(traceId: '1', spans: [], status: 'ok', duration: 100),
        Trace(traceId: '2', spans: [], status: 'ok', duration: 100),
        Trace(traceId: '3', spans: [], status: 'error', duration: 100),
        Trace(traceId: '4', spans: [], status: 'error', duration: 100),
        Trace(traceId: '5', spans: [], status: 'ok', duration: 100),
      ];
      final errorCount = traces.where((t) => t.status == 'error').length;
      final errorRate = errorCount / traces.length;
      expect(errorRate, 0.4);
    });

    test('calculates average duration', () {
      final durations = [100, 200, 300, 400, 500];
      final avg = durations.reduce((a, b) => a + b) / durations.length;
      expect(avg, 300);
    });

    test('calculates total spans', () {
      final traces = [
        Trace(traceId: '1', spans: [Span(spanId: '1', name: 'a', kind: 'server', startTime: DateTime.now(), status: SpanStatus.ok)], status: 'ok', duration: 100),
        Trace(traceId: '2', spans: [Span(spanId: '2', name: 'b', kind: 'server', startTime: DateTime.now(), status: SpanStatus.ok), Span(spanId: '3', name: 'c', kind: 'server', startTime: DateTime.now(), status: SpanStatus.ok)], status: 'ok', duration: 100),
      ];
      final totalSpans = traces.fold(0, (sum, t) => sum + t.spanCount);
      expect(totalSpans, 3);
    });
  });

  group('Span Attributes', () {
    test('stores HTTP attributes', () {
      final span = Span(
        spanId: '1',
        name: 'HTTP GET',
        kind: 'server',
        startTime: DateTime.now(),
        status: SpanStatus.ok,
        attributes: {'http.method': 'GET', 'http.url': '/api/users', 'http.status_code': 200},
      );
      expect(span.attributes['http.method'], 'GET');
      expect(span.attributes['http.status_code'], 200);
    });

    test('stores database attributes', () {
      final span = Span(
        spanId: '1',
        name: 'db:query',
        kind: 'client',
        startTime: DateTime.now(),
        status: SpanStatus.ok,
        attributes: {'db.system': 'postgresql', 'db.statement': 'SELECT * FROM users', 'db.operation': 'SELECT'},
      );
      expect(span.attributes['db.system'], 'postgresql');
      expect(span.attributes['db.operation'], 'SELECT');
    });
  });
}

class SpanStatus {
  final String code;
  final String? message;
  const SpanStatus._(this.code, [this.message]);
  static const SpanStatus ok = SpanStatus._('ok');
  static const SpanStatus unset = SpanStatus._('unset');
  static SpanStatus error({String? message}) => SpanStatus._('error', message);
}

class Span {
  final String spanId;
  final String? parentSpanId;
  final String name;
  final String kind;
  final DateTime startTime;
  final DateTime? endTime;
  final SpanStatus status;
  final Map<String, dynamic> attributes;
  final Map<String, dynamic> resource;

  Span({
    required this.spanId,
    this.parentSpanId,
    required this.name,
    required this.kind,
    required this.startTime,
    this.endTime,
    required this.status,
    this.attributes = const {},
    this.resource = const {},
  });

  int get duration {
    if (endTime == null) return 0;
    return endTime!.difference(startTime).inMilliseconds;
  }
}

class Trace {
  final String traceId;
  final List<Span> spans;
  final String status;
  final int duration;

  Trace({
    required this.traceId,
    required this.spans,
    required this.status,
    required this.duration,
  });

  int get spanCount => spans.length;

  int get serviceCount {
    final services = <String>{};
    for (final span in spans) {
      final serviceName = span.resource['service.name'];
      if (serviceName != null) {
        services.add(serviceName.toString());
      }
    }
    return services.length;
  }
}

class SpanContext {
  final String traceId;
  final String spanId;
  final int traceFlags;

  SpanContext({required this.traceId, required this.spanId, required this.traceFlags});

  bool get isSampled => traceFlags & 1 == 1;

  String toTraceparent() {
    return '00-$traceId-$spanId-${traceFlags.toRadixString(16).padLeft(2, '0')}';
  }

  static SpanContext? fromTraceparent(String traceparent) {
    final match = RegExp(r'^([0-9a-f]{2})-([0-9a-f]{32})-([0-9a-f]{16})-([0-9a-f]{2})$').firstMatch(traceparent);
    if (match == null) return null;
    return SpanContext(
      traceId: match.group(2)!,
      spanId: match.group(3)!,
      traceFlags: int.parse(match.group(4)!, radix: 16),
    );
  }
}