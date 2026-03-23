import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/services/playbooks_api.dart';
import '../../../../core/providers/api_providers.dart';

/// Playbooks state
class PlaybooksState {
  final List<Playbook> playbooks;
  final List<Playbook> builtInPlaybooks;
  final List<PlaybookCategory> categories;
  final bool isLoading;
  final String? error;
  final String? selectedCategory;

  const PlaybooksState({
    this.playbooks = const [],
    this.builtInPlaybooks = const [],
    this.categories = const [],
    this.isLoading = false,
    this.error,
    this.selectedCategory,
  });

  PlaybooksState copyWith({
    List<Playbook>? playbooks,
    List<Playbook>? builtInPlaybooks,
    List<PlaybookCategory>? categories,
    bool? isLoading,
    String? error,
    String? selectedCategory,
  }) {
    return PlaybooksState(
      playbooks: playbooks ?? this.playbooks,
      builtInPlaybooks: builtInPlaybooks ?? this.builtInPlaybooks,
      categories: categories ?? this.categories,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
      selectedCategory: selectedCategory ?? this.selectedCategory,
    );
  }

  /// Get all playbooks (custom + built-in)
  List<Playbook> get allPlaybooks => [...builtInPlaybooks, ...playbooks];

  /// Get playbooks filtered by category
  List<Playbook> get filteredPlaybooks {
    if (selectedCategory == null || selectedCategory == 'all') {
      return allPlaybooks;
    }
    return allPlaybooks.where((p) => p.category == selectedCategory).toList();
  }
}

/// Playbooks notifier
class PlaybooksNotifier extends StateNotifier<PlaybooksState> {
  final PlaybooksApi _api;

  PlaybooksNotifier(this._api) : super(const PlaybooksState()) {
    loadAll();
  }

  /// Load all playbooks data
  Future<void> loadAll() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final results = await Future.wait([
        _api.getPlaybooks(),
        _api.getBuiltInPlaybooks(),
        _api.getCategories(),
      ]);

      state = state.copyWith(
        playbooks: results[0] as List<Playbook>,
        builtInPlaybooks: results[1] as List<Playbook>,
        categories: results[2] as List<PlaybookCategory>,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// Set selected category filter
  void setCategory(String? category) {
    state = state.copyWith(selectedCategory: category);
  }

  /// Get single playbook
  Future<Playbook> getPlaybook(String name) async {
    return await _api.getPlaybook(name);
  }

  /// Upload playbook
  Future<bool> uploadPlaybook({
    required String name,
    required String content,
    String? description,
    String? category,
  }) async {
    try {
      final playbook = await _api.uploadPlaybook(
        name: name,
        content: content,
        description: description,
        category: category ?? 'custom',
      );
      state = state.copyWith(
        playbooks: [...state.playbooks, playbook],
      );
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// Delete playbook
  Future<bool> deletePlaybook(String name) async {
    try {
      await _api.deletePlaybook(name);
      state = state.copyWith(
        playbooks: state.playbooks.where((p) => p.name != name).toList(),
      );
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// Sync built-in playbooks
  Future<bool> syncBuiltIn() async {
    try {
      await _api.syncBuiltIn();
      await loadAll();
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  /// Clear error
  void clearError() {
    state = state.copyWith(error: null);
  }
}

/// Provider for PlaybooksState
final playbooksProvider = StateNotifierProvider<PlaybooksNotifier, PlaybooksState>((ref) {
  final client = ref.watch(apiClientProvider);
  return PlaybooksNotifier(client.playbooks);
});