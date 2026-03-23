import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { web3Api, ipfsApi } from '@/api'

export const useWeb3Store = defineStore('web3', () => {
  const walletAddress = ref(localStorage.getItem('walletAddress') || null)
  const did = ref(localStorage.getItem('did') || null)
  const isConnected = computed(() => !!walletAddress.value)
  const challenge = ref(null)
  const signature = ref(null)
  const isConnecting = ref(false)
  const error = ref(null)

  // Ethereum address validation
  function isValidAddress(address) {
    return /^0x[a-fA-F0-9]{40}$/.test(address)
  }

  // Generate DID from Ethereum address
  function generateDID(address) {
    return `did:ethr:${address.toLowerCase()}`
  }

  // Request wallet connection (MetaMask)
  async function connectWallet() {
    isConnecting.value = true
    error.value = null

    try {
      if (typeof window.ethereum === 'undefined') {
        throw new Error('Please install MetaMask to use Web3 features')
      }

      const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' })
      if (accounts.length === 0) {
        throw new Error('No accounts found')
      }

      walletAddress.value = accounts[0]
      did.value = generateDID(accounts[0])
      localStorage.setItem('walletAddress', accounts[0])
      localStorage.setItem('did', did.value)

      return { address: accounts[0], did: did.value }
    } catch (err) {
      error.value = err.message
      return null
    } finally {
      isConnecting.value = false
    }
  }

  // Get SIWE challenge from server
  async function getChallenge() {
    if (!walletAddress.value) {
      await connectWallet()
    }

    try {
      const response = await web3Api.challenge(walletAddress.value)
      challenge.value = response.data?.data || response.data
      return challenge.value
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to get challenge'
      return null
    }
  }

  // Sign message with wallet
  async function signMessage(message) {
    if (!walletAddress.value || typeof window.ethereum === 'undefined') {
      return null
    }

    try {
      const sig = await window.ethereum.request({
        method: 'personal_sign',
        params: [message, walletAddress.value]
      })
      signature.value = sig
      return sig
    } catch (err) {
      error.value = err.message
      return null
    }
  }

  // Verify signature and authenticate
  async function verifySignature() {
    if (!challenge.value || !signature.value) {
      return false
    }

    try {
      const response = await web3Api.verify(walletAddress.value, signature.value)
      return response.data?.success || response.data?.data?.success || false
    } catch (err) {
      error.value = err.response?.data?.error || 'Verification failed'
      return false
    }
  }

  // Full Web3 login flow
  async function web3Login() {
    isConnecting.value = true
    error.value = null

    try {
      // Step 1: Connect wallet
      const wallet = await connectWallet()
      if (!wallet) return { success: false, error: 'Wallet connection failed' }

      // Step 2: Get challenge
      const ch = await getChallenge()
      if (!ch) return { success: false, error: 'Failed to get challenge' }

      // Step 3: Sign challenge
      const sig = await signMessage(ch.message || ch.challenge || ch)
      if (!sig) return { success: false, error: 'Signature rejected' }

      // Step 4: Verify
      const verified = await verifySignature()
      return { success: verified, did: wallet.did }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      isConnecting.value = false
    }
  }

  // Store audit on blockchain
  async function storeAudit(auditData) {
    try {
      const response = await web3Api.audit({
        ...auditData,
        walletAddress: walletAddress.value,
        timestamp: new Date().toISOString()
      })
      return response.data?.data || response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Audit failed'
      return null
    }
  }

  // IPFS operations
  async function uploadToIPFS(data) {
    try {
      const response = await ipfsApi.upload(data)
      return response.data?.cid || response.data?.data?.cid
    } catch (err) {
      error.value = err.response?.data?.error || 'IPFS upload failed'
      return null
    }
  }

  async function getFromIPFS(cid) {
    try {
      const response = await ipfsApi.get(cid)
      return response.data?.data || response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'IPFS fetch failed'
      return null
    }
  }

  function disconnect() {
    walletAddress.value = null
    did.value = null
    challenge.value = null
    signature.value = null
    localStorage.removeItem('walletAddress')
    localStorage.removeItem('did')
  }

  return {
    walletAddress,
    did,
    isConnected,
    isConnecting,
    challenge,
    signature,
    error,
    isValidAddress,
    generateDID,
    connectWallet,
    getChallenge,
    signMessage,
    verifySignature,
    web3Login,
    storeAudit,
    uploadToIPFS,
    getFromIPFS,
    disconnect
  }
})