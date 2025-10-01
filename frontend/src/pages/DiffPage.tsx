import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { compareSnapshots } from '../api/api'
import { DiffReport, ServiceChange, Software } from '../types'

function DiffPage() {
  const { oldId, newId } = useParams<{ oldId: string; newId: string }>()
  const [diff, setDiff] = useState<DiffReport | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (oldId && newId) {
      fetchDiff()
    }
  }, [oldId, newId])

  const fetchDiff = async () => {
    try {
      const data = await compareSnapshots(Number(oldId), Number(newId))
      setDiff(data)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch comparison')
    } finally {
      setLoading(false)
    }
  }

  const formatDate = (timestamp: string) => {
    return new Date(timestamp).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const formatSoftware = (software?: Software) => {
    if (!software) return 'N/A'
    const parts = []
    if (software.vendor) parts.push(software.vendor)
    if (software.product) parts.push(software.product)
    if (software.version) parts.push(`v${software.version}`)
    return parts.join(' ') || 'N/A'
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="bg-red-50 text-red-800 p-4 rounded-lg">
        <p>{error}</p>
      </div>
    )
  }

  if (!diff) return null

  return (
    <div>
      <div className="mb-8">
        <Link
          to={`/hosts/${diff.old_snapshot.ip}/snapshots`}
          className="text-primary-600 hover:text-primary-700 font-medium mb-4 inline-flex items-center"
        >
          <svg className="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Back to Snapshots
        </Link>
        <h2 className="text-3xl font-bold text-gray-900 mb-2">Snapshot Comparison</h2>
        <p className="text-gray-600">Host: {diff.old_snapshot.ip}</p>
      </div>

      {/* Timeline Header */}
      <div className="bg-white rounded-xl shadow-lg p-6 mb-6">
        <div className="flex justify-between items-center">
          <div className="flex-1">
            <div className="flex items-center space-x-2 mb-2">
              <div className="w-3 h-3 bg-red-500 rounded-full"></div>
              <span className="text-sm font-medium text-gray-500">OLD SNAPSHOT</span>
            </div>
            <p className="text-lg font-semibold text-gray-900">
              {formatDate(diff.old_snapshot.timestamp)}
            </p>
            <p className="text-sm text-gray-600">{diff.old_snapshot.service_count} services</p>
          </div>
          
          <div className="px-8">
            <svg className="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
            </svg>
          </div>

          <div className="flex-1 text-right">
            <div className="flex items-center justify-end space-x-2 mb-2">
              <span className="text-sm font-medium text-gray-500">NEW SNAPSHOT</span>
              <div className="w-3 h-3 bg-green-500 rounded-full"></div>
            </div>
            <p className="text-lg font-semibold text-gray-900">
              {formatDate(diff.new_snapshot.timestamp)}
            </p>
            <p className="text-sm text-gray-600">{diff.new_snapshot.service_count} services</p>
          </div>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white rounded-xl shadow-md p-6 border-l-4 border-green-500">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Ports Added</p>
              <p className="text-3xl font-bold text-gray-900">{diff.ports_added.length}</p>
            </div>
            <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
              <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-md p-6 border-l-4 border-red-500">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Ports Removed</p>
              <p className="text-3xl font-bold text-gray-900">{diff.ports_removed.length}</p>
            </div>
            <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
              <svg className="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-md p-6 border-l-4 border-yellow-500">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Services Modified</p>
              <p className="text-3xl font-bold text-gray-900">{diff.services_changed.length}</p>
            </div>
            <div className="w-12 h-12 bg-yellow-100 rounded-full flex items-center justify-center">
              <svg className="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      {!diff.has_changes && (
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-6 text-center">
          <svg className="mx-auto h-12 w-12 text-blue-500 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <h3 className="text-lg font-semibold text-blue-900 mb-1">No Changes Detected</h3>
          <p className="text-blue-700">The snapshots are identical.</p>
        </div>
      )}

      {/* Added Ports */}
      {diff.ports_added.length > 0 && (
        <div className="bg-white rounded-xl shadow-lg p-6 mb-6">
          <h3 className="text-xl font-bold text-gray-900 mb-4 flex items-center">
            <span className="w-2 h-2 bg-green-500 rounded-full mr-3"></span>
            Ports Added
          </h3>
          <div className="space-y-4">
            {diff.ports_added.map((change) => (
              <div key={change.port} className="border border-green-200 bg-green-50 rounded-lg p-4">
                <div className="flex items-start justify-between">
                  <div>
                    <p className="font-semibold text-gray-900">
                      Port {change.port} ({change.new_service?.protocol})
                    </p>
                    {change.new_service?.software && (
                      <p className="text-sm text-gray-600 mt-1">
                        Software: {formatSoftware(change.new_service.software)}
                      </p>
                    )}
                    {change.new_service?.vulnerabilities && change.new_service.vulnerabilities.length > 0 && (
                      <div className="mt-2">
                        <span className="text-sm font-medium text-red-600">Vulnerabilities:</span>
                        <div className="flex flex-wrap gap-2 mt-1">
                          {change.new_service.vulnerabilities.map((vuln) => (
                            <span key={vuln} className="px-2 py-1 bg-red-100 text-red-800 text-xs rounded">
                              {vuln}
                            </span>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                  <span className="px-3 py-1 bg-green-600 text-white text-xs font-semibold rounded-full">
                    NEW
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Removed Ports */}
      {diff.ports_removed.length > 0 && (
        <div className="bg-white rounded-xl shadow-lg p-6 mb-6">
          <h3 className="text-xl font-bold text-gray-900 mb-4 flex items-center">
            <span className="w-2 h-2 bg-red-500 rounded-full mr-3"></span>
            Ports Removed
          </h3>
          <div className="space-y-4">
            {diff.ports_removed.map((change) => (
              <div key={change.port} className="border border-red-200 bg-red-50 rounded-lg p-4">
                <div className="flex items-start justify-between">
                  <div>
                    <p className="font-semibold text-gray-900">
                      Port {change.port} ({change.old_service?.protocol})
                    </p>
                    {change.old_service?.software && (
                      <p className="text-sm text-gray-600 mt-1">
                        Software: {formatSoftware(change.old_service.software)}
                      </p>
                    )}
                  </div>
                  <span className="px-3 py-1 bg-red-600 text-white text-xs font-semibold rounded-full">
                    REMOVED
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Modified Services */}
      {diff.services_changed.length > 0 && (
        <div className="bg-white rounded-xl shadow-lg p-6">
          <h3 className="text-xl font-bold text-gray-900 mb-4 flex items-center">
            <span className="w-2 h-2 bg-yellow-500 rounded-full mr-3"></span>
            Services Modified
          </h3>
          <div className="space-y-6">
            {diff.services_changed.map((change: ServiceChange, idx: number) => (
              <div key={idx} className="border border-yellow-200 bg-yellow-50 rounded-lg p-4">
                <p className="font-semibold text-gray-900 mb-3">
                  Port {change.port} ({change.protocol})
                </p>

                <div className="space-y-3">
                  {change.status_changed && (
                    <div className="flex items-center justify-between bg-white rounded p-3">
                      <span className="text-sm font-medium text-gray-700">Status:</span>
                      <div className="flex items-center space-x-2">
                        <span className="px-2 py-1 bg-gray-100 text-gray-800 text-sm rounded">
                          {change.old_status}
                        </span>
                        <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                        </svg>
                        <span className="px-2 py-1 bg-yellow-100 text-yellow-800 text-sm rounded font-semibold">
                          {change.new_status}
                        </span>
                      </div>
                    </div>
                  )}

                  {change.software_changed && (
                    <div className="flex items-center justify-between bg-white rounded p-3">
                      <span className="text-sm font-medium text-gray-700">Software:</span>
                      <div className="flex items-center space-x-2">
                        <span className="px-2 py-1 bg-gray-100 text-gray-800 text-sm rounded">
                          {formatSoftware(change.old_software)}
                        </span>
                        <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                        </svg>
                        <span className="px-2 py-1 bg-yellow-100 text-yellow-800 text-sm rounded font-semibold">
                          {formatSoftware(change.new_software)}
                        </span>
                      </div>
                    </div>
                  )}

                  {change.tls_changed && (
                    <div className="bg-white rounded p-3">
                      <span className="text-sm font-medium text-gray-700">TLS Configuration Changed</span>
                    </div>
                  )}

                  {change.vulnerabilities_added && change.vulnerabilities_added.length > 0 && (
                    <div className="bg-white rounded p-3">
                      <span className="text-sm font-medium text-red-600">⚠️ New Vulnerabilities:</span>
                      <div className="flex flex-wrap gap-2 mt-2">
                        {change.vulnerabilities_added.map((vuln) => (
                          <span key={vuln} className="px-2 py-1 bg-red-100 text-red-800 text-xs rounded font-semibold">
                            {vuln}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}

                  {change.vulnerabilities_fixed && change.vulnerabilities_fixed.length > 0 && (
                    <div className="bg-white rounded p-3">
                      <span className="text-sm font-medium text-green-600">✓ Fixed Vulnerabilities:</span>
                      <div className="flex flex-wrap gap-2 mt-2">
                        {change.vulnerabilities_fixed.map((vuln) => (
                          <span key={vuln} className="px-2 py-1 bg-green-100 text-green-800 text-xs rounded">
                            {vuln}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export default DiffPage
