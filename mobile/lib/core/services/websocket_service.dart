import 'dart:async';
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:web_socket_channel/status.dart' as status;

/// WebSocket service for real-time communication
class WebSocketService {
  WebSocketChannel? _channel;
  final Map<String, List<Function(dynamic)>> _handlers = {};
  Timer? _reconnectTimer;
  Timer? _heartbeatTimer;
  String? _url;
  String? _token;
  bool _isConnecting = false;
  int _reconnectAttempts = 0;
  static const int _maxReconnectAttempts = 5;
  static const Duration _reconnectDelay = Duration(seconds: 3);
  static const Duration _heartbeatInterval = Duration(seconds: 30);

  /// Whether the WebSocket is connected
  bool get isConnected => _channel != null && !_isConnecting;

  /// Stream of connection state changes
  final StreamController<bool> _connectionStateController = StreamController<bool>.broadcast();
  Stream<bool> get connectionState => _connectionStateController.stream;

  /// Connect to WebSocket server
  Future<void> connect(String url, {String? token}) async {
    if (_isConnecting || isConnected) return;

    _url = url;
    _token = token;
    _isConnecting = true;

    try {
      final uri = Uri.parse(url);
      final wsUrl = uri.replace(scheme: uri.scheme == 'https' ? 'wss' : 'ws');

      _channel = WebSocketChannel.connect(
        Uri.parse('$wsUrl?token=$token'),
      );

      await _channel!.ready;
      _isConnecting = false;
      _reconnectAttempts = 0;
      _connectionStateController.add(true);

      // Start heartbeat
      _startHeartbeat();

      // Listen for messages
      _channel!.stream.listen(
        _onMessage,
        onError: _onError,
        onDone: _onDone,
      );
    } catch (e) {
      _isConnecting = false;
      _scheduleReconnect();
    }
  }

  /// Disconnect from WebSocket server
  void disconnect() {
    _reconnectTimer?.cancel();
    _heartbeatTimer?.cancel();
    _channel?.sink.close(status.goingAway);
    _channel = null;
    _connectionStateController.add(false);
  }

  /// Subscribe to an event
  void on(String event, Function(dynamic) handler) {
    _handlers.putIfAbsent(event, () => []).add(handler);
  }

  /// Unsubscribe from an event
  void off(String event, [Function(dynamic)? handler]) {
    if (handler == null) {
      _handlers.remove(event);
    } else {
      _handlers[event]?.remove(handler);
    }
  }

  /// Emit an event to the server
  void emit(String event, dynamic data) {
    if (!isConnected) return;

    final message = jsonEncode({
      'event': event,
      'data': data,
    });

    _channel?.sink.add(message);
  }

  void _onMessage(dynamic message) {
    try {
      final decoded = jsonDecode(message);
      final event = decoded['event'] as String?;
      final data = decoded['data'];

      if (event != null && _handlers.containsKey(event)) {
        for (final handler in _handlers[event]!) {
          handler(data);
        }
      }
    } catch (e) {
      // Handle parse error
    }
  }

  void _onError(dynamic error) {
    _connectionStateController.add(false);
    _scheduleReconnect();
  }

  void _onDone() {
    _connectionStateController.add(false);
    _scheduleReconnect();
  }

  void _scheduleReconnect() {
    if (_reconnectAttempts >= _maxReconnectAttempts) return;

    _reconnectTimer?.cancel();
    _reconnectTimer = Timer(_reconnectDelay, () {
      _reconnectAttempts++;
      if (_url != null) {
        connect(_url!, token: _token);
      }
    });
  }

  void _startHeartbeat() {
    _heartbeatTimer?.cancel();
    _heartbeatTimer = Timer.periodic(_heartbeatInterval, (_) {
      emit('ping', {});
    });
  }

  /// Dispose resources
  void dispose() {
    disconnect();
    _connectionStateController.close();
  }
}

/// Global WebSocket service instance
final WebSocketService webSocketService = WebSocketService();