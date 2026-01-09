import { useState, useEffect } from 'react';
import { Upload, Download, Trash2, Lock, LogIn, LogOut, AlertCircle, CheckCircle, Loader2, FileText, RefreshCw } from 'lucide-react';

export default function App() {
  const [config, setConfig] = useState(null);
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  
  // Auth states
  const [pinVerified, setPinVerified] = useState(false);
  const [adminAuthenticated, setAdminAuthenticated] = useState(false);
  const [pin, setPin] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [showAdminLogin, setShowAdminLogin] = useState(false);
  
  // Upload state
  const [selectedFile, setSelectedFile] = useState(null);
  const [uploading, setUploading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);

  const API_BASE = '/api';

  useEffect(() => {
    fetchConfig();
  }, []);

  useEffect(() => {
    if (!config) return;
    if (config.pinProtected && !pinVerified) return;
    fetchFiles();
  }, [config, pinVerified]);

  // Auto-clear messages after 5 seconds
  useEffect(() => {
    if (error || success) {
      const timer = setTimeout(() => {
        setError('');
        setSuccess('');
      }, 5000);
      return () => clearTimeout(timer);
    }
  }, [error, success]);

  const fetchConfig = async () => {
    try {
      const res = await fetch(`${API_BASE}/config`);
      const data = await res.json();
      setConfig(data);
      
      if (!data.pinProtected) {
        setPinVerified(true);
      }
    } catch (err) {
      setError('Failed to connect to server');
    } finally {
      setLoading(false);
    }
  };

  const fetchFiles = async () => {
    setRefreshing(true);
    try {
      const res = await fetch(`${API_BASE}/files`, {
        credentials: 'include'
      });
      
      if (res.status === 401) {
        setPinVerified(false);
        return;
      }
      
      const data = await res.json();
      setFiles(data.files || []);
      setError('');
    } catch (err) {
      setError('Failed to fetch files');
    } finally {
      setRefreshing(false);
    }
  };

  const verifyPIN = async () => {
    setError('');
    
    if (!pin) {
      setError('Please enter a PIN');
      return;
    }

    try {
      const res = await fetch(`${API_BASE}/verify-pin`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ pin })
      });
      
      if (res.ok) {
        setPinVerified(true);
        setPin('');
      } else {
        setError('Invalid PIN');
      }
    } catch (err) {
      setError('Failed to verify PIN');
    }
  };

  const adminLogin = async () => {
    setError('');
    
    if (!username || !password) {
      setError('Please enter username and password');
      return;
    }

    try {
      const res = await fetch(`${API_BASE}/admin/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ username, password })
      });
      
      if (res.ok) {
        setAdminAuthenticated(true);
        setShowAdminLogin(false);
        setUsername('');
        setPassword('');
        setSuccess('Admin authenticated successfully');
      } else {
        setError('Invalid credentials');
      }
    } catch (err) {
      setError('Failed to authenticate');
    }
  };

  const adminLogout = async () => {
    try {
      await fetch(`${API_BASE}/admin/logout`, {
        method: 'POST',
        credentials: 'include'
      });
      setAdminAuthenticated(false);
      setSuccess('Logged out successfully');
    } catch (err) {
      setError('Failed to logout');
    }
  };

  const handleFileSelect = (e) => {
    const file = e.target.files[0];
    if (file) {
      if (config && file.size > config.maxFileSize) {
        setError(`File too large. Max size: ${(config.maxFileSize / (1024 * 1024)).toFixed(0)} MB`);
        return;
      }
      setSelectedFile(file);
      setError('');
    }
  };

  const uploadFile = async () => {
    if (!selectedFile) return;
    
    setUploading(true);
    setError('');
    setSuccess('');
    
    const formData = new FormData();
    formData.append('file', selectedFile);
    
    try {
      const res = await fetch(`${API_BASE}/files/upload`, {
        method: 'POST',
        credentials: 'include',
        body: formData
      });
      
      if (res.status === 401) {
        setError('Admin authentication required');
        setShowAdminLogin(true);
      } else if (res.ok) {
        setSuccess('File uploaded successfully');
        setSelectedFile(null);
        document.getElementById('fileInput').value = '';
        fetchFiles();
      } else {
        const data = await res.json();
        setError(data.error || 'Failed to upload file');
      }
    } catch (err) {
      setError('Failed to upload file');
    } finally {
      setUploading(false);
    }
  };

  const downloadFile = async (filename) => {
    try {
      const res = await fetch(`${API_BASE}/files/download/${filename}`, {
        credentials: 'include'
      });
      
      if (res.ok) {
        const blob = await res.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
        setSuccess(`Downloaded ${filename}`);
      } else {
        setError('Failed to download file');
      }
    } catch (err) {
      setError('Failed to download file');
    }
  };

  const deleteFile = async (filename) => {
    if (!confirm(`Delete ${filename}?`)) return;
    
    try {
      const res = await fetch(`${API_BASE}/files/${filename}`, {
        method: 'DELETE',
        credentials: 'include'
      });
      
      if (res.status === 401) {
        setError('Admin authentication required');
        setShowAdminLogin(true);
      } else if (res.ok) {
        setSuccess(`Deleted ${filename}`);
        fetchFiles();
      } else {
        setError('Failed to delete file');
      }
    } catch (err) {
      setError('Failed to delete file');
    }
  };

  const formatBytes = (bytes) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
  };

  const handlePinKeyPress = (e) => {
    if (e.key === 'Enter') {
      verifyPIN();
    }
  };

  const handleAdminKeyPress = (e) => {
    if (e.key === 'Enter') {
      adminLogin();
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="w-12 h-12 animate-spin text-blue-600 mx-auto mb-4" />
          <p className="text-gray-600">Connecting to server...</p>
        </div>
      </div>
    );
  }

  // PIN Entry Screen
  if (config?.pinProtected && !pinVerified) {
    return (
      <div className="min-h-screen flex items-center justify-center p-4">
        <div className="bg-white rounded-2xl shadow-xl p-8 w-full max-w-md animate-fade-in">
          <div className="flex items-center justify-center mb-6">
            <div className="bg-blue-100 p-4 rounded-full">
              <Lock className="w-8 h-8 text-blue-600" />
            </div>
          </div>
          <h2 className="text-2xl font-bold text-center mb-2">Protected Access</h2>
          <p className="text-gray-600 text-center mb-6">Enter PIN to continue</p>
          
          <div>
            <input
              type="password"
              value={pin}
              onChange={(e) => setPin(e.target.value)}
              onKeyPress={handlePinKeyPress}
              placeholder="Enter PIN"
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent mb-4 text-center text-2xl tracking-widest"
              autoFocus
              maxLength={6}
            />
            
            {error && (
              <div className="flex items-center gap-2 text-red-600 text-sm mb-4 bg-red-50 p-3 rounded-lg">
                <AlertCircle className="w-4 h-4 flex-shrink-0" />
                <span>{error}</span>
              </div>
            )}
            
            <button
              onClick={verifyPIN}
              className="w-full bg-blue-600 text-white py-3 rounded-lg hover:bg-blue-700 transition font-medium"
            >
              Verify PIN
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <div className="max-w-6xl mx-auto p-6">
        {/* Header */}
        <div className="bg-white rounded-2xl shadow-lg p-6 mb-6 animate-fade-in">
          <div className="flex items-center justify-between flex-wrap gap-4">
            <div className="flex items-center gap-3 min-w-0">
              <img src="/src/assets/localshare.png" alt="LocalShare" className="h-34 w-15 w-auto" />
              <div className='leading-tight'>
                <h1 className='text-lg font-semibold text-slate-800'>LocalShare</h1>
                <p className="text-sm text-slate-500">Share files on your local network</p>
              </div>
            </div>
            
            <div className="flex items-center gap-3">
              {config?.pinProtected && (
                <span className="flex items-center gap-2 px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm">
                  <Lock className="w-4 h-4" />
                  PIN Verified
                </span>
              )}
              
              {config?.adminRequired && (
                adminAuthenticated ? (
                  <button
                    onClick={adminLogout}
                    className="flex items-center gap-2 px-4 py-2 bg-gray-100 hover:bg-gray-200 rounded-lg transition text-sm"
                  >
                    <LogOut className="w-4 h-4" />
                    Logout
                  </button>
                ) : (
                  <button
                    onClick={() => setShowAdminLogin(true)}
                    className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white hover:bg-blue-700 rounded-lg transition text-sm"
                  >
                    <LogIn className="w-4 h-4" />
                    Admin Login
                  </button>
                )
              )}
            </div>
          </div>
        </div>

        {/* Alerts */}
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6 flex items-start gap-2 text-red-700 animate-slide-up">
            <AlertCircle className="w-5 h-5 flex-shrink-0 mt-0.5" />
            <span className="flex-1">{error}</span>
            <button onClick={() => setError('')} className="text-red-400 hover:text-red-600">
              ×
            </button>
          </div>
        )}
        
        {success && (
          <div className="bg-green-50 border border-green-200 rounded-lg p-4 mb-6 flex items-start gap-2 text-green-700 animate-slide-up">
            <CheckCircle className="w-5 h-5 flex-shrink-0 mt-0.5" />
            <span className="flex-1">{success}</span>
            <button onClick={() => setSuccess('')} className="text-green-400 hover:text-green-600">
              ×
            </button>
          </div>
        )}

        {/* Admin Login Modal */}
        {showAdminLogin && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50 animate-fade-in">
            <div className="bg-white rounded-2xl p-8 w-full max-w-md animate-slide-up">
              <h3 className="text-xl font-bold mb-4">Admin Login</h3>
              <div>
                <input
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="Username"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg mb-3 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  autoFocus
                />
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  onKeyPress={handleAdminKeyPress}
                  placeholder="Password"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg mb-4 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                <div className="flex gap-3">
                  <button
                    onClick={adminLogin}
                    className="flex-1 bg-blue-600 text-white py-3 rounded-lg hover:bg-blue-700 transition"
                  >
                    Login
                  </button>
                  <button
                    onClick={() => {
                      setShowAdminLogin(false);
                      setUsername('');
                      setPassword('');
                    }}
                    className="flex-1 bg-gray-200 text-gray-800 py-3 rounded-lg hover:bg-gray-300 transition"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Upload Section */}
        <div className="bg-white rounded-2xl shadow-lg p-6 mb-6 animate-fade-in">
          <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
            <Upload className="w-5 h-5" />
            Upload Files
          </h2>
          
          <div className="flex flex-col sm:flex-row gap-3">
            <input
              id="fileInput"
              type="file"
              onChange={handleFileSelect}
              className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
            />
            <button
              onClick={uploadFile}
              disabled={!selectedFile || uploading}
              className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition flex items-center justify-center gap-2"
            >
              {uploading ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  Uploading...
                </>
              ) : (
                <>
                  <Upload className="w-4 h-4" />
                  Upload
                </>
              )}
            </button>
          </div>
          
          {selectedFile && (
            <div className="mt-3 p-3 bg-blue-50 rounded-lg flex items-center gap-2">
              <FileText className="w-4 h-4 text-blue-600" />
              <div className="flex-1">
                <p className="text-sm font-medium text-gray-900">{selectedFile.name}</p>
                <p className="text-xs text-gray-600">{formatBytes(selectedFile.size)}</p>
              </div>
            </div>
          )}
        </div>

        {/* Files List */}
        <div className="bg-white rounded-2xl shadow-lg p-6 animate-fade-in">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold">Files ({files.length})</h2>
            <button
              onClick={fetchFiles}
              disabled={refreshing}
              className="flex items-center gap-2 px-3 py-2 text-sm text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition"
            >
              <RefreshCw className={`w-4 h-4 ${refreshing ? 'animate-spin' : ''}`} />
              Refresh
            </button>
          </div>
          
          {files.length === 0 ? (
            <div className="text-center py-12">
              <div className="inline-flex items-center justify-center w-16 h-16 bg-gray-100 rounded-full mb-4">
                <FileText className="w-8 h-8 text-gray-400" />
              </div>
              <p className="text-gray-500">No files yet</p>
              <p className="text-sm text-gray-400 mt-1">Upload some files to get started</p>
            </div>
          ) : (
            <div className="space-y-2 max-h-[500px] overflow-y-auto custom-scrollbar pr-2">
              {files.map((file) => (
                <div
                  key={file.name}
                  className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition group"
                >
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <FileText className="w-4 h-4 text-gray-400 flex-shrink-0" />
                      <p className="font-medium text-gray-900 truncate">{file.name}</p>
                    </div>
                    <div className="flex items-center gap-3 text-xs text-gray-500">
                      <span>{formatBytes(file.size)}</span>
                      <span>•</span>
                      <span>{formatDate(file.modifiedTime)}</span>
                    </div>
                  </div>
                  
                  <div className="flex gap-2 ml-4">
                    <button
                      onClick={() => downloadFile(file.name)}
                      className="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition"
                      title="Download"
                    >
                      <Download className="w-5 h-5" />
                    </button>
                    <button
                      onClick={() => deleteFile(file.name)}
                      className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition opacity-0 group-hover:opacity-100"
                      title="Delete"
                    >
                      <Trash2 className="w-5 h-5" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="mt-6 text-center text-sm text-gray-500">
          <p>LocalShare v1.0 • Files are stored locally on the server</p>
        </div>
      </div>
    </div>
  );
}