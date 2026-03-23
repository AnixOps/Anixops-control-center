import 'package:flutter_test/flutter_test.dart';

void main() {
  group('Webhook', () {
    test('creates webhook', () {
      final webhook = Webhook(id: 'w1', name: 'Slack', url: 'https://hooks.slack.com/test', events: ['alert.created'], enabled: true);
      expect(webhook.name, 'Slack');
      expect(webhook.enabled, isTrue);
    });

    test('toggles webhook', () {
      final webhook = Webhook(id: 'w1', name: 'test', url: 'https://example.com', events: [], enabled: true);
      webhook.enabled = false;
      expect(webhook.enabled, isFalse);
    });

    test('validates URL', () {
      final webhook = Webhook(id: 'w1', name: 'test', url: 'https://example.com', events: [], enabled: true);
      expect(webhook.isSecure, isTrue);
    });
  });

  group('Webhook Delivery', () {
    test('creates delivery', () {
      final delivery = WebhookDelivery(webhookId: 'w1', status: 'success', attempts: 1);
      expect(delivery.status, 'success');
      expect(delivery.attempts, 1);
    });

    test('checks if can retry', () {
      final delivery = WebhookDelivery(webhookId: 'w1', status: 'failed', attempts: 2);
      expect(delivery.canRetry(3), isTrue);
      expect(delivery.canRetry(2), isFalse);
    });
  });
}

class Webhook {
  final String id;
  final String name;
  final String url;
  final List<String> events;
  bool enabled;

  Webhook({required this.id, required this.name, required this.url, required this.events, required this.enabled});

  bool get isSecure => url.startsWith('https://');
}

class WebhookDelivery {
  final String webhookId;
  final String status;
  final int attempts;

  WebhookDelivery({required this.webhookId, required this.status, required this.attempts});

  bool canRetry(int maxAttempts) => attempts < maxAttempts;
}