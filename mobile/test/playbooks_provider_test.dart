import 'package:flutter_test/flutter_test.dart';
import 'package:anixops_mobile/core/services/playbooks_api.dart';
import 'package:anixops_mobile/features/playbooks/presentation/providers/playbooks_provider.dart';

void main() {
  group('Playbook model', () {
    test('is created correctly from JSON', () {
      final json = {
        'name': 'install-docker',
        'display_name': 'Install Docker',
        'description': 'Install Docker on target nodes',
        'category': 'software',
        'source': 'built-in',
        'variables': {'version': 'latest'},
        'tags': ['docker', 'container'],
        'author': 'anixops',
        'version': '1.0.0',
        'created_at': '2026-03-20T10:00:00Z',
        'updated_at': '2026-03-20T12:00:00Z',
      };

      final playbook = Playbook.fromJson(json);

      expect(playbook.name, 'install-docker');
      expect(playbook.displayName, 'Install Docker');
      expect(playbook.description, 'Install Docker on target nodes');
      expect(playbook.category, 'software');
      expect(playbook.source, 'built-in');
      expect(playbook.variables, {'version': 'latest'});
      expect(playbook.tags, ['docker', 'container']);
      expect(playbook.author, 'anixops');
      expect(playbook.version, '1.0.0');
    });

    test('title getter formats name correctly', () {
      final playbook = Playbook(name: 'install-docker-compose');

      expect(playbook.title, 'Install Docker Compose');
    });

    test('title uses displayName if available', () {
      final playbook = Playbook(
        name: 'install-docker',
        displayName: 'Custom Docker Name',
      );

      expect(playbook.title, 'Custom Docker Name');
    });

    test('handles missing optional fields', () {
      final json = {
        'name': 'minimal-playbook',
      };

      final playbook = Playbook.fromJson(json);

      expect(playbook.name, 'minimal-playbook');
      expect(playbook.displayName, isNull);
      expect(playbook.description, isNull);
      expect(playbook.category, isNull);
      expect(playbook.variables, isNull);
      expect(playbook.tags, isNull);
    });
  });

  group('PlaybookCategory model', () {
    test('is created correctly from JSON', () {
      final json = {
        'id': 'security',
        'name': 'Security',
        'icon': 'shield',
        'description': 'Security-related playbooks',
      };

      final category = PlaybookCategory.fromJson(json);

      expect(category.id, 'security');
      expect(category.name, 'Security');
      expect(category.icon, 'shield');
      expect(category.description, 'Security-related playbooks');
    });
  });

  group('PlaybooksState', () {
    test('allPlaybooks combines built-in and custom playbooks', () {
      final state = PlaybooksState(
        builtInPlaybooks: [
          Playbook(name: 'built-in-1'),
          Playbook(name: 'built-in-2'),
        ],
        playbooks: [
          Playbook(name: 'custom-1'),
        ],
      );

      expect(state.allPlaybooks.length, 3);
    });

    test('filteredPlaybooks returns all when no category selected', () {
      final state = PlaybooksState(
        playbooks: [
          Playbook(name: 'playbook-1', category: 'security'),
          Playbook(name: 'playbook-2', category: 'maintenance'),
        ],
        builtInPlaybooks: [
          Playbook(name: 'built-in-1', category: 'software'),
        ],
      );

      expect(state.filteredPlaybooks.length, 3);
    });

    test('filteredPlaybooks filters by selected category', () {
      final state = PlaybooksState(
        playbooks: [
          Playbook(name: 'playbook-1', category: 'security'),
          Playbook(name: 'playbook-2', category: 'maintenance'),
        ],
        builtInPlaybooks: [
          Playbook(name: 'built-in-1', category: 'security'),
        ],
        selectedCategory: 'security',
      );

      expect(state.filteredPlaybooks.length, 2);
      expect(state.filteredPlaybooks.every((p) => p.category == 'security'), true);
    });

    test('filteredPlaybooks handles "all" category', () {
      final state = PlaybooksState(
        playbooks: [
          Playbook(name: 'playbook-1', category: 'security'),
        ],
        selectedCategory: 'all',
      );

      expect(state.filteredPlaybooks.length, 1);
    });

    test('copyWith works correctly', () {
      const state = PlaybooksState(
        playbooks: [],
        isLoading: false,
      );

      final newState = state.copyWith(
        isLoading: true,
        error: 'Test error',
        selectedCategory: 'security',
      );

      expect(newState.isLoading, true);
      expect(newState.error, 'Test error');
      expect(newState.selectedCategory, 'security');
      expect(newState.playbooks, isEmpty);
    });
  });
}