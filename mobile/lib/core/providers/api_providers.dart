import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:anixops_mobile/core/services/api_client.dart';
import 'package:anixops_mobile/core/services/websocket_service.dart';

/// Provider for the API client
final apiClientProvider = Provider<ApiClient>((ref) {
  return apiClient;
});

/// Provider for the WebSocket service
final webSocketServiceProvider = Provider<WebSocketService>((ref) {
  return webSocketService;
});

/// Provider for WebSocket connection state
final wsConnectionStateProvider = StreamProvider<bool>((ref) {
  final ws = ref.watch(webSocketServiceProvider);
  return ws.connectionState;
});