import axios from 'axios'
import { Snapshot, DiffReport } from '../types'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

const api = axios.create({
  baseURL: API_BASE_URL,
})

export const uploadSnapshot = async (file: File): Promise<{ message: string; ip: string; timestamp: string }> => {
  const formData = new FormData()
  formData.append('file', file)

  const response = await api.post('/snapshots/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })

  return response.data
}

export const getSnapshot = async (id: number): Promise<Snapshot> => {
  const response = await api.get(`/snapshots/${id}`)
  return response.data
}

export const listHosts = async (): Promise<string[]> => {
  const response = await api.get('/hosts')
  return response.data.hosts || []
}

export const listSnapshotsByIP = async (ip: string): Promise<Snapshot[]> => {
  const response = await api.get(`/hosts/${ip}/snapshots`)
  return response.data.snapshots || []
}

export const compareSnapshots = async (oldId: number, newId: number): Promise<DiffReport> => {
  const response = await api.get(`/diff/${oldId}/${newId}`)
  return response.data
}
