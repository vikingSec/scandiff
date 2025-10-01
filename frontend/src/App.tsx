import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom'
import UploadPage from './pages/UploadPage'
import HostsPage from './pages/HostsPage'
import SnapshotsPage from './pages/SnapshotsPage'
import SnapshotDetailPage from './pages/SnapshotDetailPage'
import DiffPage from './pages/DiffPage'

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
        <nav className="bg-white shadow-lg border-b border-gray-200">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between items-center h-16">
              <div className="flex items-center space-x-2">
                <svg className="w-8 h-8 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                <h1 className="text-2xl font-bold text-gray-900">ScanDiff</h1>
              </div>
              <div className="flex space-x-4">
                <Link
                  to="/"
                  className="px-4 py-2 rounded-lg text-gray-700 hover:bg-primary-50 hover:text-primary-700 transition-colors font-medium"
                >
                  Upload
                </Link>
                <Link
                  to="/hosts"
                  className="px-4 py-2 rounded-lg text-gray-700 hover:bg-primary-50 hover:text-primary-700 transition-colors font-medium"
                >
                  Hosts
                </Link>
              </div>
            </div>
          </div>
        </nav>

        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <Routes>
            <Route path="/" element={<UploadPage />} />
            <Route path="/hosts" element={<HostsPage />} />
            <Route path="/hosts/:ip/snapshots" element={<SnapshotsPage />} />
            <Route path="/hosts/:ip/snapshots/:id" element={<SnapshotDetailPage />} />
            <Route path="/diff/:oldId/:newId" element={<DiffPage />} />
          </Routes>
        </main>
      </div>
    </Router>
  )
}

export default App
