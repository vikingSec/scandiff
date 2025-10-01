import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { getSnapshot } from '../api/api'
import { Snapshot, Software } from '../types'

function SnapshotDetailPage() {
  const { ip, id } = useParams<{ ip: string; id: string }>()
  const [snapshot, setSnapshot] = useState<Snapshot | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (id) {
      fetchSnapshot()
    }
  }, [id])

  const fetchSnapshot = async () => {
    try {
      const data = await getSnapshot(Number(id))
      setSnapshot(data)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch snapshot')
    } finally {
      setLoading(false)
    }
  }

  const formatDate = (timestamp: string) => {
    return new Date(timestamp).toLocaleString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      timeZoneName: 'short',
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

  if (error || !snapshot) {
    return (
      <div className="bg-red-50 text-red-800 p-4 rounded-lg">
        <p>{error || 'Snapshot not found'}</p>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-8">
        <Link
          to={`/hosts/${ip}/snapshots`}
          className="text-primary-600 hover:text-primary-700 font-medium mb-4 inline-flex items-center"
        >
          <svg className="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Back to Snapshots
        </Link>
        <h2 className="text-3xl font-bold text-gray-900 mb-2">Snapshot Details</h2>
        <p className="text-gray-600">Complete view of host state at this point in time</p>
      </div>

      {/* Snapshot Header */}
      <div className="bg-white rounded-xl shadow-lg p-6 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div>
            <p className="text-sm font-medium text-gray-500 mb-1">Host IP</p>
            <p className="text-lg font-bold text-gray-900">{snapshot.ip}</p>
          </div>
          <div>
            <p className="text-sm font-medium text-gray-500 mb-1">Snapshot Time</p>
            <p className="text-lg font-semibold text-gray-900">{formatDate(snapshot.timestamp)}</p>
          </div>
          <div>
            <p className="text-sm font-medium text-gray-500 mb-1">Total Services</p>
            <p className="text-lg font-bold text-primary-600">{snapshot.service_count}</p>
          </div>
          <div>
            <p className="text-sm font-medium text-gray-500 mb-1">Filename</p>
            <p className="text-sm text-gray-700 truncate">{snapshot.filename}</p>
          </div>
        </div>
      </div>

      {/* Services Summary */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white rounded-xl shadow-md p-6 border-l-4 border-blue-500">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Open Ports</p>
              <p className="text-3xl font-bold text-gray-900">{snapshot.services.length}</p>
            </div>
            <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
              <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-md p-6 border-l-4 border-yellow-500">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Vulnerabilities</p>
              <p className="text-3xl font-bold text-gray-900">
                {snapshot.services.reduce((count, service) => 
                  count + (service.vulnerabilities?.length || 0), 0
                )}
              </p>
            </div>
            <div className="w-12 h-12 bg-yellow-100 rounded-full flex items-center justify-center">
              <svg className="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-md p-6 border-l-4 border-green-500">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Secure Services</p>
              <p className="text-3xl font-bold text-gray-900">
                {snapshot.services.filter(s => s.tls).length}
              </p>
            </div>
            <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
              <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      {/* Services Details */}
      <div className="bg-white rounded-xl shadow-lg p-6">
        <h3 className="text-xl font-bold text-gray-900 mb-6 flex items-center">
          <svg className="w-6 h-6 mr-2 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 10h16M4 14h16M4 18h16" />
          </svg>
          Services Running
        </h3>

        {snapshot.services.length === 0 ? (
          <div className="text-center py-12">
            <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
            </svg>
            <h3 className="mt-2 text-sm font-medium text-gray-900">No services found</h3>
            <p className="mt-1 text-sm text-gray-500">This snapshot has no running services.</p>
          </div>
        ) : (
          <div className="space-y-4">
            {snapshot.services.map((service, idx) => (
              <div
                key={idx}
                className="border border-gray-200 rounded-lg p-5 hover:shadow-md transition-shadow"
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center space-x-4">
                    <div className="flex items-center justify-center w-12 h-12 bg-primary-100 rounded-lg">
                      <span className="text-lg font-bold text-primary-700">{service.port}</span>
                    </div>
                    <div>
                      <h4 className="text-lg font-semibold text-gray-900">
                        Port {service.port} - {service.protocol}
                      </h4>
                      {service.status && (
                        <p className="text-sm text-gray-600">Status: {service.status}</p>
                      )}
                    </div>
                  </div>
                  {service.tls && (
                    <span className="px-3 py-1 bg-green-100 text-green-800 text-xs font-semibold rounded-full flex items-center">
                      <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clipRule="evenodd" />
                      </svg>
                      TLS Enabled
                    </span>
                  )}
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {/* Software Information */}
                  {service.software && (
                    <div className="bg-gray-50 rounded-lg p-4">
                      <h5 className="text-sm font-semibold text-gray-700 mb-2 flex items-center">
                        <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
                        </svg>
                        Software
                      </h5>
                      <p className="text-sm text-gray-900">{formatSoftware(service.software)}</p>
                      {service.software.vendor && (
                        <p className="text-xs text-gray-500 mt-1">Vendor: {service.software.vendor}</p>
                      )}
                    </div>
                  )}

                  {/* TLS Information */}
                  {service.tls && (
                    <div className="bg-green-50 rounded-lg p-4">
                      <h5 className="text-sm font-semibold text-green-900 mb-2 flex items-center">
                        <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clipRule="evenodd" />
                        </svg>
                        TLS Configuration
                      </h5>
                      <p className="text-sm text-gray-900">Version: {service.tls.version || 'N/A'}</p>
                      <p className="text-xs text-gray-600 mt-1">Cipher: {service.tls.cipher || 'N/A'}</p>
                      {service.tls.cert_fingerprint_sha256 && (
                        <p className="text-xs text-gray-600 mt-1 break-all">
                          Cert: {service.tls.cert_fingerprint_sha256.substring(0, 16)}...
                        </p>
                      )}
                    </div>
                  )}
                </div>

                {/* Vulnerabilities */}
                {service.vulnerabilities && service.vulnerabilities.length > 0 && (
                  <div className="mt-4 bg-red-50 border border-red-200 rounded-lg p-4">
                    <h5 className="text-sm font-semibold text-red-900 mb-2 flex items-center">
                      <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                      </svg>
                      Vulnerabilities ({service.vulnerabilities.length})
                    </h5>
                    <div className="flex flex-wrap gap-2">
                      {service.vulnerabilities.map((vuln, vIdx) => (
                        <span
                          key={vIdx}
                          className="px-3 py-1 bg-red-100 text-red-800 text-xs font-semibold rounded-full"
                        >
                          {vuln}
                        </span>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default SnapshotDetailPage
