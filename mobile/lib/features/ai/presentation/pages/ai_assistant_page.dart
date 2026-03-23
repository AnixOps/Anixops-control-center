import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:anixops_mobile/core/services/api_client.dart';

/// AI Chat message model
class ChatMessage {
  final String role;
  final String content;
  final DateTime timestamp;

  ChatMessage({
    required this.role,
    required this.content,
    required this.timestamp,
  });
}

/// Chat history provider
final chatHistoryProvider = StateProvider<List<ChatMessage>>((ref) => []);

/// Thinking state provider
final isThinkingProvider = StateProvider<bool>((ref) => false);

/// AI Assistant Page
class AIAssistantPage extends ConsumerStatefulWidget {
  const AIAssistantPage({super.key});

  @override
  ConsumerState<AIAssistantPage> createState() => _AIAssistantPageState();
}

class _AIAssistantPageState extends ConsumerState<AIAssistantPage> {
  final _controller = TextEditingController();
  final _scrollController = ScrollController();
  int _activeTab = 0;

  @override
  void dispose() {
    _controller.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  Future<void> _sendMessage() async {
    final message = _controller.text.trim();
    if (message.isEmpty) return;

    _controller.clear();
    final history = ref.read(chatHistoryProvider);

    // Add user message
    ref.read(chatHistoryProvider.notifier).state = [
      ...history,
      ChatMessage(role: 'user', content: message, timestamp: DateTime.now()),
    ];

    ref.read(isThinkingProvider.notifier).state = true;

    try {
      final response = await apiClient.ai.chat(
        message,
        history: history.map((m) => {'role': m.role, 'content': m.content}).toList(),
      );

      final assistantMessage = response['response'] ?? response['data']?['response'] ?? 'I could not generate a response.';

      ref.read(chatHistoryProvider.notifier).state = [
        ...ref.read(chatHistoryProvider),
        ChatMessage(role: 'assistant', content: assistantMessage, timestamp: DateTime.now()),
      ];
    } catch (e) {
      ref.read(chatHistoryProvider.notifier).state = [
        ...ref.read(chatHistoryProvider),
        ChatMessage(
          role: 'assistant',
          content: 'Sorry, I encountered an error: $e',
          timestamp: DateTime.now(),
        ),
      ];
    } finally {
      ref.read(isThinkingProvider.notifier).state = false;
    }

    // Scroll to bottom
    Future.delayed(const Duration(milliseconds: 100), () {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  void _quickPrompt(String prompt) {
    _controller.text = prompt;
    _sendMessage();
  }

  @override
  Widget build(BuildContext context) {
    final chatHistory = ref.watch(chatHistoryProvider);
    final isThinking = ref.watch(isThinkingProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('AI Assistant'),
            Text(
              'Powered by Workers AI',
              style: TextStyle(fontSize: 12, fontWeight: FontWeight.normal),
            ),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.delete_outline),
            onPressed: () {
              ref.read(chatHistoryProvider.notifier).state = [];
            },
            tooltip: 'Clear chat',
          ),
        ],
      ),
      body: Column(
        children: [
          // Tab bar
          Container(
            padding: const EdgeInsets.all(16),
            child: SegmentedButton<int>(
              segments: const [
                ButtonSegment(value: 0, label: Text('Chat')),
                ButtonSegment(value: 1, label: Text('Analyze')),
                ButtonSegment(value: 2, label: Text('Search')),
              ],
              selected: {_activeTab},
              onSelectionChanged: (Set<int> selection) {
                setState(() {
                  _activeTab = selection.first;
                });
              },
            ),
          ),

          // Content
          Expanded(
            child: _activeTab == 0
                ? _buildChatTab(chatHistory, isThinking)
                : _activeTab == 1
                    ? _buildAnalyzeTab()
                    : _buildSearchTab(),
          ),
        ],
      ),
    );
  }

  Widget _buildChatTab(List<ChatMessage> chatHistory, bool isThinking) {
    return Column(
      children: [
        // Chat messages
        Expanded(
          child: chatHistory.isEmpty
              ? _buildEmptyState()
              : ListView.builder(
                  controller: _scrollController,
                  padding: const EdgeInsets.all(16),
                  itemCount: chatHistory.length + (isThinking ? 1 : 0),
                  itemBuilder: (context, index) {
                    if (index == chatHistory.length && isThinking) {
                      return _buildThinkingIndicator();
                    }
                    final message = chatHistory[index];
                    return _buildMessageBubble(message);
                  },
                ),
        ),

        // Input
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: Theme.of(context).colorScheme.surface,
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.1),
                blurRadius: 4,
                offset: const Offset(0, -2),
              ),
            ],
          ),
          child: Row(
            children: [
              Expanded(
                child: TextField(
                  controller: _controller,
                  decoration: const InputDecoration(
                    hintText: 'Ask about your infrastructure...',
                    border: OutlineInputBorder(),
                  ),
                  onSubmitted: (_) => _sendMessage(),
                ),
              ),
              const SizedBox(width: 8),
              IconButton.filled(
                icon: const Icon(Icons.send),
                onPressed: isThinking ? null : _sendMessage,
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 80,
            height: 80,
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [Colors.blue, Colors.purple],
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
              borderRadius: BorderRadius.circular(24),
            ),
            child: const Icon(Icons.psychology, size: 40, color: Colors.white),
          ),
          const SizedBox(height: 16),
          Text(
            'AI DevOps Assistant',
            style: Theme.of(context).textTheme.headlineSmall,
          ),
          const SizedBox(height: 8),
          const Text('Ask me anything about your infrastructure'),
          const SizedBox(height: 24),
          Wrap(
            spacing: 8,
            children: [
              ActionChip(
                label: const Text('Show node status'),
                onPressed: () => _quickPrompt('Show node status'),
              ),
              ActionChip(
                label: const Text('Analyze errors'),
                onPressed: () => _quickPrompt('Analyze recent errors'),
              ),
              ActionChip(
                label: const Text('Security tips'),
                onPressed: () => _quickPrompt('Best practices for security'),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildMessageBubble(ChatMessage message) {
    final isUser = message.role == 'user';
    return Align(
      alignment: isUser ? Alignment.centerRight : Alignment.centerLeft,
      child: Container(
        margin: const EdgeInsets.symmetric(vertical: 4),
        padding: const EdgeInsets.all(12),
        constraints: BoxConstraints(
          maxWidth: MediaQuery.of(context).size.width * 0.8,
        ),
        decoration: BoxDecoration(
          color: isUser
              ? Theme.of(context).colorScheme.primary
              : Theme.of(context).colorScheme.surfaceContainerHighest,
          borderRadius: BorderRadius.circular(16),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              message.content,
              style: TextStyle(
                color: isUser ? Colors.white : null,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              '${message.timestamp.hour}:${message.timestamp.minute.toString().padLeft(2, '0')}',
              style: TextStyle(
                fontSize: 10,
                color: isUser ? Colors.white70 : null,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildThinkingIndicator() {
    return Container(
      margin: const EdgeInsets.symmetric(vertical: 4),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surfaceContainerHighest,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          SizedBox(
            width: 16,
            height: 16,
            child: CircularProgressIndicator(
              strokeWidth: 2,
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
          const SizedBox(width: 12),
          const Text('Thinking...'),
        ],
      ),
    );
  }

  Widget _buildAnalyzeTab() {
    final logController = TextEditingController();
    String? result;

    return StatefulBuilder(
      builder: (context, setState) {
        return Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(
                'Log Analysis',
                style: Theme.of(context).textTheme.titleLarge,
              ),
              const SizedBox(height: 8),
              const Text('Paste your logs for AI-powered analysis'),
              const SizedBox(height: 16),
              TextField(
                controller: logController,
                maxLines: 10,
                decoration: const InputDecoration(
                  hintText: 'Paste log content here...',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),
              FilledButton.icon(
                icon: const Icon(Icons.analytics),
                label: const Text('Analyze Logs'),
                onPressed: () async {
                  if (logController.text.isEmpty) return;

                  setState(() => result = null);

                  try {
                    final response = await apiClient.ai.analyzeLog(logController.text);
                    setState(() {
                      result = response['analysis'] ?? response.toString();
                    });
                  } catch (e) {
                    setState(() {
                      result = 'Error: $e';
                    });
                  }
                },
              ),
              if (result != null) ...[
                const SizedBox(height: 16),
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Analysis Result',
                          style: Theme.of(context).textTheme.titleMedium,
                        ),
                        const SizedBox(height: 8),
                        Text(result!),
                      ],
                    ),
                  ),
                ),
              ],
            ],
          ),
        );
      },
    );
  }

  Widget _buildSearchTab() {
    final searchController = TextEditingController();
    List<Map<String, dynamic>> results = [];

    return StatefulBuilder(
      builder: (context, setState) {
        return Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(
                'Semantic Search',
                style: Theme.of(context).textTheme.titleLarge,
              ),
              const SizedBox(height: 8),
              const Text('Search logs and tasks using natural language'),
              const SizedBox(height: 16),
              TextField(
                controller: searchController,
                decoration: const InputDecoration(
                  hintText: 'e.g., "connection timeout errors"',
                  border: OutlineInputBorder(),
                  suffixIcon: Icon(Icons.search),
                ),
                onSubmitted: (_) async {
                  if (searchController.text.isEmpty) return;

                  try {
                    final searchResults = await apiClient.ai.semanticSearch(searchController.text);
                    setState(() {
                      results = searchResults;
                    });
                  } catch (e) {
                    // Handle error
                  }
                },
              ),
              const SizedBox(height: 16),
              Expanded(
                child: results.isEmpty
                    ? const Center(child: Text('No results'))
                    : ListView.builder(
                        itemCount: results.length,
                        itemBuilder: (context, index) {
                          final result = results[index];
                          return Card(
                            child: ListTile(
                              title: Text(result['id'] ?? 'Unknown'),
                              subtitle: Text(result['metadata']?.toString() ?? ''),
                              trailing: Text(
                                result['score']?.toStringAsFixed(3) ?? '',
                                style: Theme.of(context).textTheme.bodySmall,
                              ),
                            ),
                          );
                        },
                      ),
              ),
            ],
          ),
        );
      },
    );
  }
}