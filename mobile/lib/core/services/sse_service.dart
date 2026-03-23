import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;

/// SSE Service for real-time communication with Workers API
/// Uses Server-Sent Events instead of WebSocket due to Cloudflare Durable Object issues
class SSEService {
  http.Client? _client;
  final Map<String, List<Function(dynamic)>> _handlers = {};
  Timer? _reconnectTimer;
  String? _url;
  String? _token;
  bool _isConnecting = false;
  bool _shouldReconnect = true;
  int _reconnectAttempts = 0;
  static const int _maxReconnectAttempts = 10;
  static const Duration _baseReconnectDelay = Duration(seconds: 1);
  static const Duration _maxReconnectDelay = Duration(seconds: 30);

  /// Whether the SSE is connected
  bool get isConnected => _client != null && !_isConnecting;

  /// Stream of connection state changes
  final StreamController<bool> _connectionStateController =
      StreamController<bool>.broadcast();
  Stream<bool> get connectionState => _connectionStateController.stream;

  /// Current subscribed channels
  final Set<String> _subscribedChannels = {};
  Set<String> get subscribedChannels => Set.unmodifiable(_subscribedChannels);

  /// Connect to SSE endpoint
  Future<void> connect(String url, {String? token}) async {
    if (_isConnecting || isConnected) return;

    _url = url;
    _token = token;
    _isConnecting = true;
    _shouldReconnect = true;

    try {
      _client = http.Client();

      final request = http.Request('GET', Uri.parse(url));
      request.headers['Accept'] = 'text/event-stream';
      request.headers['Cache-Control'] = 'no-cache';
      if (token != null && token.isNotEmpty) {
        request.headers['Authorization'] = 'Bearer $token';
      }

      final response = await _client!.send(request);

      if (response.statusCode != 200) {
        throw Exception('SSE connection failed: ${response.statusCode}');
      }

      _isConnecting = false;
      _reconnectAttempts = 0;
      _connectionStateController.add(true);

      // Listen to the stream
      final stream = response.stream.transform(utf8.decoder);
      String buffer = '';

      await for (final chunk in stream) {
        if (!_shouldReconnect) break;

        buffer += chunk;
        _processBuffer(buffer);
        buffer = _clearProcessedEvents(buffer);
      }
    } catch (e) {
      _isConnecting = false;
      _client?.close();
      _client = null;
      _connectionStateController.add(false);

      if (_shouldReconnect) {
        _scheduleReconnect();
      }
    }
  }

  /// Process the buffer for complete events
  void _processBuffer(String buffer) {
    // Split by double newlines (event separator)
    final events = buffer.split('\n\n');

    for (int i = 0; i < events.length - 1; i++) {
      final event = events[i].trim();
      if (event.isNotEmpty) {
        _parseAndDispatchEvent(event);
      }
    }
  }

  /// Clear processed events from buffer
  String _clearProcessedEvents(String buffer) {
    final lastSeparator = buffer.lastIndexOf('\n\n');
    if (lastSeparator >= 0) {
      return buffer.substring(lastSeparator + 2);
    }
    return buffer;
  }

  /// Parse and dispatch an SSE event
  void _parseAndDispatchEvent(String eventText) {
    String? eventType;
    String? data;
    String? id;

    for (final line in eventText.split('\n')) {
      if (line.startsWith('event:')) {
        eventType = line.substring(6).trim();
      } else if (line.startsWith('data:')) {
        data = line.substring(5).trim();
      } else if (line.startsWith('id:')) {
        id = line.substring(3).trim();
      } else if (line.startsWith(':')) {
        // Comment line (heartbeat), ignore
      }
    }

    // Handle heartbeat
    if (eventText.contains(': heartbeat')) {
      return;
    }

    // Parse data as JSON if present
    dynamic parsedData;
    if (data != null && data.isNotEmpty) {
      try {
        parsedData = jsonDecode(data);
      } catch (e) {
        parsedData = data;
      }
    }

    // Dispatch to handlers
    if (eventType != null && _handlers.containsKey(eventType)) {
      for (final handler in List.from(_handlers[eventType]!)) {
        try {
          handler(parsedData);
        } catch (e) {
          // Handler error, continue
        }
      }
    }

    // Also dispatch to generic 'message' handlers if no specific event type
    if (eventType == null && _handlers.containsKey('message')) {
      for (final handler in List.from(_handlers['message']!)) {
        try {
          handler(parsedData);
        } catch (e) {
          // Handler error, continue
        }
      }
    }

    // Dispatch by type field in data (Workers API format)
    if (parsedData is Map && parsedData.containsKey('type')) {
      final type = parsedData['type'] as String;
      if (_handlers.containsKey(type)) {
        for (final handler in List.from(_handlers[type]!)) {
          try {
            handler(parsedData['payload'] ?? parsedData);
          } catch (e) {
            // Handler error, continue
          }
        }
      }
    }
  }

  /// Subscribe to a channel via REST API
  Future<bool> subscribe(String channel) async {
    if (_url == null || _token == null) return false;

    try {
      // Extract base URL from SSE URL
      final baseUri = Uri.parse(_url!);
      final baseUrl = '${baseUri.scheme}://${baseUri.host}';
      final subscribeUrl = '$baseUrl/api/v1/sse/subscribe';

      final response = await http.post(
        Uri.parse(subscribeUrl),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $_token',
        },
        body: jsonEncode({'channel': channel}),
      );

      if (response.statusCode == 200) {
        _subscribedChannels.add(channel);
        return true;
      }
    } catch (e) {
      // Subscribe failed
    }
    return false;
  }

  /// Unsubscribe from a channel
  Future<bool> unsubscribe(String channel) async {
    if (_url == null || _token == null) return false;

    try {
      final baseUri = Uri.parse(_url!);
      final baseUrl = '${baseUri.scheme}://${baseUri.host}';
      final unsubscribeUrl = '$baseUrl/api/v1/sse/unsubscribe';

      final response = await http.post(
        Uri.parse(unsubscribeUrl),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $_token',
        },
        body: jsonEncode({'channel': channel}),
      );

      if (response.statusCode == 200) {
        _subscribedChannels.remove(channel);
        return true;
      }
    } catch (e) {
      // Unsubscribe failed
    }
    return false;
  }

  /// Register a handler for an event type
  void on(String eventType, Function(dynamic) handler) {
    _handlers.putIfAbsent(eventType, () => []).add(handler);
  }

  /// Remove a handler for an event type
  void off(String eventType, [Function(dynamic)? handler]) {
    if (handler == null) {
      _handlers.remove(eventType);
    } else {
      _handlers[eventType]?.remove(handler);
    }
  }

  /// Disconnect from SSE endpoint
  void disconnect() {
    _shouldReconnect = false;
    _reconnectTimer?.cancel();
    _client?.close();
    _client = null;
    _subscribedChannels.clear();
    _connectionStateController.add(false);
  }

  void _scheduleReconnect() {
    if (!_shouldReconnect) return;
    if (_reconnectAttempts >= _maxReconnectAttempts) return;

    _reconnectTimer?.cancel();

    // Exponential backoff
    final delay = Duration(
      milliseconds: (_baseReconnectDelay.inMilliseconds *
              (1 << _reconnectAttempts.clamp(0, 5)))
          .clamp(_baseReconnectDelay.inMilliseconds,
              _maxReconnectDelay.inMilliseconds),
    );

    _reconnectTimer = Timer(delay, () {
      _reconnectAttempts++;
      if (_url != null) {
        connect(_url!, token: _token);
      }
    });
  }

  /// Reset reconnection attempts (call on successful operation)
  void resetReconnectAttempts() {
    _reconnectAttempts = 0;
  }

  /// Dispose resources
  void dispose() {
    disconnect();
    _connectionStateController.close();
  }
}

/// Global SSE service instance
final SSEService sseService = SSEService();