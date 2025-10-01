import { test, expect } from '@playwright/test'

test.describe('ScanDiff Application', () => {
  test('should display the homepage with upload functionality', async ({ page }) => {
    await page.goto('/')
    
    // Check if the main heading is visible
    await expect(page.getByRole('heading', { name: 'Upload Snapshot' })).toBeVisible()
    
    // Check if the upload button is present but disabled initially
    const uploadButton = page.getByRole('button', { name: 'Upload Snapshot' })
    await expect(uploadButton).toBeVisible()
    await expect(uploadButton).toBeDisabled()
  })

  test('should navigate to hosts page', async ({ page }) => {
    await page.goto('/')
    
    // Click on Hosts link
    await page.getByRole('link', { name: 'Hosts' }).click()
    
    // Check if we're on the hosts page
    await expect(page.getByRole('heading', { name: 'Monitored Hosts' })).toBeVisible()
  })

  test('should show file selection in upload page', async ({ page }) => {
    await page.goto('/')
    
    // Create a test file
    const fileContent = JSON.stringify({
      timestamp: '2025-09-10T03:00:00Z',
      ip: '192.168.1.1',
      services: [],
      service_count: 0
    })
    
    // Set up the file chooser handler
    const fileChooserPromise = page.waitForEvent('filechooser')
    await page.click('label[for="file-input"]')
    const fileChooser = await fileChooserPromise
    
    // Create a buffer from the file content
    await fileChooser.setFiles({
      name: 'test_snapshot.json',
      mimeType: 'application/json',
      buffer: Buffer.from(fileContent)
    })
    
    // Verify the file is displayed
    await expect(page.getByText('test_snapshot.json')).toBeVisible()
    
    // Verify upload button is now enabled
    const uploadButton = page.getByRole('button', { name: 'Upload Snapshot' })
    await expect(uploadButton).toBeEnabled()
  })

  test('should have proper navigation links', async ({ page }) => {
    await page.goto('/')
    
    // Check navigation links exist
    await expect(page.getByRole('link', { name: 'Upload' })).toBeVisible()
    await expect(page.getByRole('link', { name: 'Hosts' })).toBeVisible()
    
    // Check ScanDiff branding
    await expect(page.getByRole('heading', { name: 'ScanDiff' })).toBeVisible()
  })

  test('should display empty state on hosts page when no hosts exist', async ({ page }) => {
    // Mock the API response
    await page.route('**/api/hosts', async (route) => {
      await route.fulfill({
        status: 200,
        body: JSON.stringify({ hosts: [] }),
      })
    })
    
    await page.goto('/hosts')
    
    // Check for empty state
    await expect(page.getByText('No Hosts Found')).toBeVisible()
    await expect(page.getByText('Upload your first snapshot to start monitoring hosts.')).toBeVisible()
  })

  test('should handle drag and drop in upload area', async ({ page }) => {
    await page.goto('/')
    
    // The drag area should be present
    const dropZone = page.locator('div').filter({ hasText: 'Click to upload or drag and drop' }).first()
    await expect(dropZone).toBeVisible()
  })
})
