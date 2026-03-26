import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:anixops_mobile/core/services/api_client.dart';

/// Wallet connection state providers
final walletAddressProvider = NotifierProvider<WalletAddressNotifier, String?>(WalletAddressNotifier.new);
final didProvider = NotifierProvider<DIDNotifier, String?>(DIDNotifier.new);
final isConnectingProvider = NotifierProvider<IsConnectingNotifier, bool>(IsConnectingNotifier.new);

class WalletAddressNotifier extends Notifier<String?> {
  @override
  String? build() => null;

  void setAddress(String? address) {
    state = address;
  }
}

class DIDNotifier extends Notifier<String?> {
  @override
  String? build() => null;

  void setDID(String? did) {
    state = did;
  }
}

class IsConnectingNotifier extends Notifier<bool> {
  @override
  bool build() => false;

  void setConnecting(bool value) {
    state = value;
  }
}

/// Web3 Dashboard Page
class Web3DashboardPage extends ConsumerStatefulWidget {
  const Web3DashboardPage({super.key});

  @override
  ConsumerState<Web3DashboardPage> createState() => _Web3DashboardPageState();
}

class _Web3DashboardPageState extends ConsumerState<Web3DashboardPage> {
  int _statsIpfsFiles = 0;
  int _statsOnChainAudits = 0;
  String? _ipfsResult;
  String? _auditResult;
  bool _uploading = false;
  bool _storing = false;

  final _auditDetailsController = TextEditingController();
  String _auditAction = 'node.restart';

  final List<Map<String, dynamic>> _recentAudits = [];

  @override
  void dispose() {
    _auditDetailsController.dispose();
    super.dispose();
  }

  bool get _isValidEthereumAddress {
    final address = ref.read(walletAddressProvider);
    if (address == null) return false;
    return RegExp(r'^0x[a-fA-F0-9]{40}$').hasMatch(address);
  }

  String _generateDID(String address) {
    return 'did:ethr:${address.toLowerCase()}';
  }

  Future<void> _connectWallet() async {
    ref.read(isConnectingProvider.notifier).setConnecting(true);

    try {
      // Simulate wallet connection
      // In production, this would use walletconnect or web3dart
      const mockAddress = '0x742d35Cc6634C0532925a3b844Bc9e7595f2bD18';

      ref.read(walletAddressProvider.notifier).setAddress(mockAddress);
      ref.read(didProvider.notifier).setDID(_generateDID(mockAddress));

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Wallet connected successfully!')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to connect wallet: $e')),
        );
      }
    } finally {
      ref.read(isConnectingProvider.notifier).setConnecting(false);
    }
  }

  Future<void> _uploadToIPFS() async {
    setState(() {
      _uploading = true;
      _ipfsResult = null;
    });

    try {
      final response = await apiClient.web3.uploadToIPFS({'test': 'data'});
      setState(() {
        _ipfsResult = response['cid'] ?? 'QmTest123...';
        _statsIpfsFiles++;
      });
    } catch (e) {
      setState(() {
        _ipfsResult = 'Error: $e';
      });
    } finally {
      setState(() {
        _uploading = false;
      });
    }
  }

  Future<void> _storeAudit() async {
    if (!_isValidEthereumAddress) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please connect wallet first')),
      );
      return;
    }

    setState(() {
      _storing = true;
      _auditResult = null;
    });

    try {
      final response = await apiClient.web3.storeAudit({
        'action': _auditAction,
        'details': _auditDetailsController.text,
        'timestamp': DateTime.now().toIso8601String(),
      });

      setState(() {
        _auditResult = response['txHash'] ?? '0x123...abc';
        _statsOnChainAudits++;
        _recentAudits.insert(0, {
          'action': _auditAction,
          'details': _auditDetailsController.text,
          'timestamp': DateTime.now().toIso8601String(),
          'txHash': _auditResult,
        });
        _auditDetailsController.clear();
      });
    } catch (e) {
      setState(() {
        _auditResult = 'Error: $e';
      });
    } finally {
      setState(() {
        _storing = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final walletAddress = ref.watch(walletAddressProvider);
    final did = ref.watch(didProvider);
    final isConnecting = ref.watch(isConnectingProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Web3 Dashboard'),
            Text(
              'IPFS & Blockchain Integration',
              style: TextStyle(fontSize: 12, fontWeight: FontWeight.normal),
            ),
          ],
        ),
        actions: [
          if (walletAddress != null)
            Container(
              margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
              decoration: BoxDecoration(
                color: Colors.green.withAlpha(50),
                borderRadius: BorderRadius.circular(16),
                border: Border.all(color: Colors.green),
              ),
              child: Row(
                children: [
                  Container(
                    width: 8,
                    height: 8,
                    decoration: const BoxDecoration(
                      color: Colors.green,
                      shape: BoxShape.circle,
                    ),
                  ),
                  const SizedBox(width: 8),
                  const Text('Connected', style: TextStyle(color: Colors.green)),
                ],
              ),
            )
          else
            TextButton.icon(
              icon: const Icon(Icons.link),
              label: const Text('Connect Wallet'),
              onPressed: isConnecting ? null : _connectWallet,
            ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Stats cards
            SizedBox(
              height: 120,
              child: ListView(
                scrollDirection: Axis.horizontal,
                children: [
                  _buildStatCard(
                    icon: Icons.folder,
                    iconColor: Colors.blue,
                    title: 'IPFS Files',
                    value: _statsIpfsFiles.toString(),
                  ),
                  _buildStatCard(
                    icon: Icons.security,
                    iconColor: Colors.purple,
                    title: 'On-Chain Audits',
                    value: _statsOnChainAudits.toString(),
                  ),
                  _buildStatCard(
                    icon: Icons.account_balance_wallet,
                    iconColor: Colors.orange,
                    title: 'DID',
                    value: did ?? 'Not set',
                    isLongValue: true,
                  ),
                  _buildStatCard(
                    icon: Icons.network_check,
                    iconColor: Colors.green,
                    title: 'Network',
                    value: 'Ethereum',
                  ),
                ],
              ),
            ),

            const SizedBox(height: 24),

            // Main content
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // IPFS Card
                Expanded(
                  child: _buildIpfsCard(),
                ),
                const SizedBox(width: 16),
                // Audit Card
                Expanded(
                  child: _buildAuditCard(),
                ),
              ],
            ),

            const SizedBox(height: 24),

            // Recent audits
            Text(
              'Recent Audit Records',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Card(
              child: _recentAudits.isEmpty
                  ? const Padding(
                      padding: EdgeInsets.all(24),
                      child: Center(child: Text('No audit records yet')),
                    )
                  : ListView.builder(
                      shrinkWrap: true,
                      itemCount: _recentAudits.length,
                      itemBuilder: (context, index) {
                        final audit = _recentAudits[index];
                        return ListTile(
                          title: Text(audit['action']),
                          subtitle: Text(audit['details']),
                          trailing: Text(
                            audit['txHash']?.toString().substring(0, 10) ?? '',
                            style: Theme.of(context).textTheme.bodySmall,
                          ),
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatCard({
    required IconData icon,
    required Color iconColor,
    required String title,
    required String value,
    bool isLongValue = false,
  }) {
    return Container(
      width: isLongValue ? 200 : 150,
      margin: const EdgeInsets.only(right: 12),
      child: Card(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: iconColor.withAlpha(50),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Icon(icon, color: iconColor),
              ),
              const SizedBox(height: 12),
              Text(
                title,
                style: Theme.of(context).textTheme.bodySmall,
              ),
              const SizedBox(height: 4),
              Text(
                value,
                style: Theme.of(context).textTheme.titleMedium,
                overflow: isLongValue ? TextOverflow.ellipsis : null,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildIpfsCard() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'IPFS Storage',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            const Text(
              'Files are stored on IPFS decentralized network',
              style: TextStyle(color: Colors.grey),
            ),
            const SizedBox(height: 16),
            SizedBox(
              width: double.infinity,
              child: OutlinedButton.icon(
                icon: const Icon(Icons.upload_file),
                label: const Text('Select File to Upload'),
                onPressed: () {
                  // File picker would go here
                  _uploadToIPFS();
                },
              ),
            ),
            const SizedBox(height: 16),
            FilledButton.icon(
              icon: _uploading
                  ? const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                    )
                  : const Icon(Icons.cloud_upload),
              label: Text(_uploading ? 'Uploading...' : 'Upload to IPFS'),
              onPressed: _uploading ? null : _uploadToIPFS,
            ),
            if (_ipfsResult != null) ...[
              const SizedBox(height: 16),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.green.withAlpha(50),
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(color: Colors.green),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Uploaded Successfully!',
                      style: TextStyle(color: Colors.green, fontWeight: FontWeight.bold),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      'CID: $_ipfsResult',
                      style: Theme.of(context).textTheme.bodySmall,
                    ),
                  ],
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildAuditCard() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'On-Chain Audit Log',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 16),
            DropdownButtonFormField<String>(
              value: _auditAction,
              decoration: const InputDecoration(
                labelText: 'Action',
                border: OutlineInputBorder(),
              ),
              items: const [
                DropdownMenuItem(value: 'node.restart', child: Text('Node Restart')),
                DropdownMenuItem(value: 'node.create', child: Text('Node Create')),
                DropdownMenuItem(value: 'node.delete', child: Text('Node Delete')),
                DropdownMenuItem(value: 'playbook.run', child: Text('Playbook Run')),
                DropdownMenuItem(value: 'user.login', child: Text('User Login')),
                DropdownMenuItem(value: 'settings.change', child: Text('Settings Change')),
              ],
              onChanged: (value) {
                setState(() {
                  _auditAction = value ?? 'node.restart';
                });
              },
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _auditDetailsController,
              maxLines: 3,
              decoration: const InputDecoration(
                labelText: 'Details',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 16),
            FilledButton.icon(
              icon: _storing
                  ? const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                    )
                  : const Icon(Icons.security),
              label: Text(_storing ? 'Recording...' : 'Record on Blockchain'),
              onPressed: _storing ? null : _storeAudit,
            ),
            if (_auditResult != null) ...[
              const SizedBox(height: 16),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.green.withAlpha(50),
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(color: Colors.green),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Audit Recorded!',
                      style: TextStyle(color: Colors.green, fontWeight: FontWeight.bold),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      'Tx: $_auditResult',
                      style: Theme.of(context).textTheme.bodySmall,
                    ),
                  ],
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}