import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'api_client.dart';

/// Web3 API service for blockchain/IPFS integration
class Web3Api {
  final Dio _dio;

  Web3Api(this._dio);

  /// Get SIWE challenge for wallet authentication
  Future<Map<String, dynamic>> getChallenge(String address) async {
    final response = await _dio.post('/web3/challenge', data: {'address': address});
    return response.data;
  }

  /// Verify wallet signature
  Future<Map<String, dynamic>> verifySignature(String address, String signature) async {
    final response = await _dio.post('/web3/verify', data: {
      'address': address,
      'signature': signature,
    });
    return response.data;
  }

  /// Store audit on blockchain
  Future<Map<String, dynamic>> storeAudit(Map<String, dynamic> auditData) async {
    final response = await _dio.post('/web3/audit', data: auditData);
    return response.data;
  }

  /// Upload file to IPFS
  Future<Map<String, dynamic>> uploadToIPFS(dynamic data) async {
    final response = await _dio.post('/ipfs/upload', data: data);
    return response.data;
  }

  /// Get file from IPFS
  Future<Map<String, dynamic>> getFromIPFS(String cid) async {
    final response = await _dio.get('/ipfs/$cid');
    return response.data;
  }
}