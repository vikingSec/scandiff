import { useEffect, useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { listSnapshotsByIP } from '../api/api'
import { Snapshot } from '../types'

function SnapshotsPage() {
  const { ip } = useParams<{ ip: string }>()
  const navigate = useNavigate()
  const [snapshots, setSnapshots] = useState<Snapshot[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedSnapshots, setSelectedSnapshots] = useState<number[]>([])

  useEffect(() => {
    if (ip) {
      fetchSnapshots()
    }
  }, [ip])

  const fetchSnapshots = async () => {
    try {
      const data = await listSnapshotsByIP(ip!)
      setSnapshots(data)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch snapshots')
    } finally {
      setLoading(false)
    }
  }

  const handleSnapshotSelect = (id: number) => {
    if (selectedSnapshots.includes(id)) {
      setSelectedSnapshots(selectedSnapshots.filter((s) => s !== id))
    } else if (selectedSnapshots.length < 2) {
      setSelectedSnapshots([...selectedSnapshots, id])
    } else {
      setSelectedSnapshots([selectedSnapshots[1], id])
    }
  }

  const handleCompare = () => {
    if (selectedSnapshots.length === 2) {
      const [id1, id2] = selectedSnapshots.sort((a, b) => a - b)
      navigate(`/diff/${id1}/${id2}`)
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

  return (
    <div>
      <div className="mb-8">
        <Link to="/hosts" className="text-primary-600 hover:text-primary-700 font-medium mb-4 inline-flex items-center">
          <svg className="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Back to Hosts
        </Link>
        <h2 className="text-3xl font-bold text-gray-900 mb-2">Snapshots for {ip}</h2>
        <p className="text-gray-600">
          Select two snapshots to compare their configurations and identify changes.
        </p>
      </div>

      {selectedSnapshots.length === 2 && (
        <div className="mb-6 bg-primary-50 border border-primary-200 rounded-lg p-4 flex items-center justify-between">
          <div className="flex items-center">
            <svg className="w-5 h-5 text-primary-600 mr-2" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
            </svg>
            <span className="text-primary-800 font-medium">2 snapshots selected for comparison</span>
          </div>
          <button
            onClick={handleCompare}
            className="bg-primary-600 hover:bg-primary-700 text-white font-semibold py-2 px-6 rounded-lg transition-colors"
          >
            Compare Snapshots
          </button>
        </div>
      )}

      <div className="bg-white rounded-xl shadow-lg overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Select
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Timestamp
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Services
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Filename
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {snapshots.map((snapshot) => (
                <tr
                  key={snapshot.id}
                  className={`hover:bg-gray-50 ${
                    selectedSnapshots.includes(snapshot.id!) ? 'bg-primary-50' : ''
                  }`}
                >
                  <td className="px-6 py-4 whitespace-nowrap">
                    <input
                      type="checkbox"
                      checked={selectedSnapshots.includes(snapshot.id!)}
                      onChange={() => handleSnapshotSelect(snapshot.id!)}
                      className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded cursor-pointer"
                    />
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900">
                      {formatDate(snapshot.timestamp)}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className="px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                      {snapshot.service_count} services
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {snapshot.filename}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <Link
                      to={`/hosts/${ip}/snapshots/${snapshot.id}`}
                      className="text-primary-600 hover:text-primary-900 inline-flex items-center"
                    >
                      <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                      </svg>
                      View Details
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {snapshots.length === 0 && (
        <div className="text-center py-12">
          <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
          </svg>
          <h3 className="mt-2 text-sm font-medium text-gray-900">No snapshots</h3>
          <p className="mt-1 text-sm text-gray-500">
            Upload snapshots for this host to get started.
          </p>
        </div>
      )}
    </div>
  )
}

export default SnapshotsPage
