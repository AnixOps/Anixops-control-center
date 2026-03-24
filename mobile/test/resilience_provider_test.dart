import 'package:flutter_test/flutter_test.dart';

void main() {
  group('CircuitBreaker', () {
    test('creates circuit breaker in closed state', () {
      final breaker = CircuitBreaker(name: 'api-service', failureThreshold: 5);
      expect(breaker.state, CircuitBreakerState.closed);
      expect(breaker.failureCount, 0);
    });

    test('opens after failure threshold', () {
      final breaker = CircuitBreaker(name: 'test', failureThreshold: 3);
      breaker.recordFailure();
      expect(breaker.state, CircuitBreakerState.closed);
      breaker.recordFailure();
      expect(breaker.state, CircuitBreakerState.closed);
      breaker.recordFailure();
      expect(breaker.state, CircuitBreakerState.open);
    });

    test('resets on success', () {
      final breaker = CircuitBreaker(name: 'test', failureThreshold: 3);
      breaker.recordFailure();
      breaker.recordFailure();
      breaker.recordSuccess();
      expect(breaker.failureCount, 0);
    });

    test('opens from half-open on failure', () {
      final breaker = CircuitBreaker(name: 'test', failureThreshold: 3);
      breaker.state = CircuitBreakerState.halfOpen;
      breaker.recordFailure();
      expect(breaker.state, CircuitBreakerState.open);
    });
  });

  group('RateLimiter', () {
    test('creates rate limiter with tokens', () {
      final limiter = RateLimiter(name: 'api', maxTokens: 100, refillRate: 10);
      expect(limiter.tokens, 100);
    });

    test('consumes tokens', () {
      final limiter = RateLimiter(name: 'api', maxTokens: 100, refillRate: 10);
      final allowed = limiter.tryConsume();
      expect(allowed, isTrue);
      expect(limiter.tokens, 99);
    });

    test('rejects when no tokens', () {
      final limiter = RateLimiter(name: 'api', maxTokens: 100, refillRate: 10);
      limiter.tokens = 0;
      final allowed = limiter.tryConsume();
      expect(allowed, isFalse);
    });
  });

  group('RetryConfig', () {
    test('creates retry config', () {
      final config = RetryConfig(name: 'api', maxRetries: 3, backoffMultiplier: 2, initialDelay: 100);
      expect(config.maxRetries, 3);
      expect(config.backoffMultiplier, 2);
    });

    test('checks if can retry', () {
      final config = RetryConfig(name: 'api', maxRetries: 3, backoffMultiplier: 2, initialDelay: 100);
      expect(config.canRetry(1), isTrue);
      expect(config.canRetry(3), isTrue);
      expect(config.canRetry(4), isFalse);
    });
  });

  group('ResilienceStats', () {
    test('calculates open breakers', () {
      final breakers = <CircuitBreaker>[
        CircuitBreaker(name: 'a', failureThreshold: 5)..state = CircuitBreakerState.open,
        CircuitBreaker(name: 'b', failureThreshold: 5)..state = CircuitBreakerState.closed,
      ];
      final openCount = breakers.where((b) => b.state == CircuitBreakerState.open).length;
      expect(openCount, 1);
    });

    test('calculates total available tokens', () {
      final limiters = <RateLimiter>[
        RateLimiter(name: 'a', maxTokens: 100, refillRate: 10)..tokens = 80,
        RateLimiter(name: 'b', maxTokens: 200, refillRate: 20)..tokens = 150,
      ];
      final total = limiters.fold<int>(0, (sum, l) => sum + l.tokens);
      expect(total, 230);
    });
  });
}

enum CircuitBreakerState { closed, open, halfOpen }

class CircuitBreaker {
  final String name;
  final int failureThreshold;
  final int successThreshold;
  int failureCount = 0;
  int successCount = 0;
  CircuitBreakerState state = CircuitBreakerState.closed;

  CircuitBreaker({required this.name, required this.failureThreshold, this.successThreshold = 3});

  void recordFailure() {
    failureCount++;
    if (state == CircuitBreakerState.halfOpen) {
      state = CircuitBreakerState.open;
    } else if (failureCount >= failureThreshold) {
      state = CircuitBreakerState.open;
    }
  }

  void recordSuccess() {
    successCount++;
    failureCount = 0;
    if (state == CircuitBreakerState.halfOpen && successCount >= successThreshold) {
      state = CircuitBreakerState.closed;
      successCount = 0;
    }
  }
}

class RateLimiter {
  final String name;
  final int maxTokens;
  final int refillRate;
  int tokens;

  RateLimiter({required this.name, required this.maxTokens, required this.refillRate}) : tokens = maxTokens;

  bool tryConsume() {
    if (tokens > 0) {
      tokens--;
      return true;
    }
    return false;
  }
}

class RetryConfig {
  final String name;
  final int maxRetries;
  final int backoffMultiplier;
  final int initialDelay;
  final int maxDelay;

  RetryConfig({
    required this.name,
    required this.maxRetries,
    required this.backoffMultiplier,
    required this.initialDelay,
    this.maxDelay = 30000,
  });

  bool canRetry(int attempt) => attempt <= maxRetries;
}