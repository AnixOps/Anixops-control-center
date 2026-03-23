import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../core/services/backup_api.dart';
import '../../core/providers/api_providers.dart';

/// Backup state
class BackupState {
  final List<Backup> backups;
  final BackupStatus? status;
  final bool isLoading;
  final String? error;

  const BackupState({
    this.backups = const [],
    this.status,
    this.isLoading = false,
    this.error,
  });

  BackupState copyWith({
    List<Backup>? backups,
    BackupStatus? status,
    bool? isLoading,
    String? error,
  }) {
    return BackupState(
      backups: backups ?? this.backups,
      status: status ?? this.status,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }

  List<Backup> get completedBackups =>
      backups.where((b) => b.isCompleted).toList();

  List<Backup> get pendingBackups =>
      backups.where((b) => b.isPending).toList();

  List<Backup> get failedBackups =>
      backups.where((b) => b.isFailed).toList();
}

/// Backup notifier
class BackupNotifier extends StateNotifier<BackupState> {
  final BackupApi _api;

  BackupNotifier(this._api) : super(const BackupState());

  Future<void> loadBackups() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final backups = await _api.listBackups();
      final status = await _api.getStatus();
      state = state.copyWith(
        backups: backups,
        status: status,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  Future<Backup?> createBackup({String? name, String? description}) async {
    try {
      final backup = await _api.createBackup(
        name: name,
        description: description,
      );
      await loadBackups();
      return backup;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return null;
    }
  }

  Future<bool> restoreBackup(String id) async {
    try {
      await _api.restoreBackup(id);
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  Future<bool> deleteBackup(String id) async {
    try {
      await _api.deleteBackup(id);
      final backups = state.backups.where((b) => b.id != id).toList();
      state = state.copyWith(backups: backups);
      return true;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return false;
    }
  }

  Future<int> cleanupBackups({int keepLast = 10}) async {
    try {
      final deletedCount = await _api.cleanupBackups(keepLast: keepLast);
      await loadBackups();
      return deletedCount;
    } catch (e) {
      state = state.copyWith(error: e.toString());
      return 0;
    }
  }

  void clearError() {
    state = state.copyWith(error: null);
  }
}

/// Provider for backup state
final backupProvider =
    StateNotifierProvider<BackupNotifier, BackupState>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return BackupNotifier(apiClient.backup);
});

/// Provider for backup status
final backupStatusProvider = Provider<BackupStatus?>((ref) {
  return ref.watch(backupProvider).status;
});