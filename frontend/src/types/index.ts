export interface Software {
  vendor?: string
  product?: string
  version?: string
}

export interface TLS {
  version?: string
  cipher?: string
  cert_fingerprint_sha256?: string
}

export interface Service {
  port: number
  protocol: string
  status?: number
  software?: Software
  tls?: TLS
  vulnerabilities?: string[]
}

export interface Snapshot {
  id?: number
  ip: string
  timestamp: string
  filename?: string
  services: Service[]
  service_count: number
  uploaded_at?: string
}

export type ChangeType = 'added' | 'removed' | 'modified'

export interface PortChange {
  port: number
  change_type: ChangeType
  old_service?: Service
  new_service?: Service
}

export interface ServiceChange {
  port: number
  protocol: string
  status_changed: boolean
  old_status?: number
  new_status?: number
  software_changed: boolean
  old_software?: Software
  new_software?: Software
  tls_changed: boolean
  old_tls?: TLS
  new_tls?: TLS
  vulnerabilities_added?: string[]
  vulnerabilities_fixed?: string[]
}

export interface DiffReport {
  old_snapshot: Snapshot
  new_snapshot: Snapshot
  ports_added: PortChange[]
  ports_removed: PortChange[]
  services_changed: ServiceChange[]
  has_changes: boolean
}
