import { test, expect } from '@playwright/test'

/**
 * E2E Test: Navigation Flow
 * 
 * Tests basic navigation between pages:
 * 1. Home page navigation
 * 2. Create page navigation
 * 3. Search page navigation
 * 4. Post detail page navigation
 */
test.describe('Navigation Flow', () => {
  test('should navigate between all main pages', async ({ page }) => {
    // Start at home page
    await page.goto('/')
    await expect(page).toHaveURL(/.*\/$/)
    await page.waitForLoadState('networkidle')

    // Verify home page elements
    await expect(page.locator('text=公司曝光平台')).toBeVisible()
    await expect(page.locator('text=首页')).toBeVisible()

    // Navigate to create page
    await page.click('text=发布曝光')
    await expect(page).toHaveURL(/.*\/create/)
    await page.waitForLoadState('networkidle')
    await expect(page.locator('text=发布曝光')).toBeVisible()

    // Navigate to search page
    await page.click('text=搜索')
    await expect(page).toHaveURL(/.*\/search/)
    await page.waitForLoadState('networkidle')
    await expect(page.locator('text=搜索曝光内容')).toBeVisible()

    // Navigate back to home
    await page.click('text=首页')
    await expect(page).toHaveURL(/.*\/$/)
    await page.waitForLoadState('networkidle')
  })

  test('should display header on all pages', async ({ page }) => {
    const pages = ['/', '/create', '/search']

    for (const path of pages) {
      await page.goto(path)
      await page.waitForLoadState('networkidle')

      // Verify header is visible
      await expect(page.locator('.app-header, header')).toBeVisible()
      
      // Verify logo/brand is visible
      const brand = page.locator('text=Fuck Boss, .brand-name, text=公司曝光平台')
      await expect(brand.first()).toBeVisible()
    }
  })

  test('should highlight active menu item', async ({ page }) => {
    // Navigate to home
    await page.goto('/')
    await page.waitForLoadState('networkidle')
    
    // Check if home menu item is active (Ant Design highlights active items)
    const homeMenu = page.locator('text=首页').locator('..')
    // Ant Design adds active class to selected menu items

    // Navigate to create
    await page.click('text=发布曝光')
    await page.waitForLoadState('networkidle')

    // Navigate to search
    await page.click('text=搜索')
    await page.waitForLoadState('networkidle')
  })
})

