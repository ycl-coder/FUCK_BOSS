import { test, expect } from '@playwright/test'

/**
 * E2E Test: Search Flow
 * 
 * Tests the search functionality:
 * 1. Navigate to search page
 * 2. Enter search keyword
 * 3. Optionally select city filter
 * 4. Verify search results are displayed
 * 5. Click on a search result to view details
 */
test.describe('Search Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to home page first
    await page.goto('/')
    await page.waitForLoadState('networkidle')
  })

  test('should navigate to search page and display search bar', async ({ page }) => {
    // Navigate to search page
    await page.click('text=搜索')
    await expect(page).toHaveURL(/.*\/search/)

    // Verify search page elements
    await expect(page.locator('text=搜索曝光内容')).toBeVisible()
    await expect(page.locator('input[placeholder*="关键词"]')).toBeVisible()
  })

  test('should perform search from home page', async ({ page }) => {
    // Find search input on home page
    const searchInput = page.locator('input[placeholder*="搜索公司名称"]')
    
    // Check if search input exists on home page
    const searchExists = await searchInput.isVisible().catch(() => false)
    
    if (searchExists) {
      // Enter search keyword
      const keyword = '测试'
      await searchInput.fill(keyword)
      await searchInput.press('Enter')
      
      // Wait for navigation or results
      await page.waitForTimeout(1000)
      
      // Verify we're on search page or results are shown
      const currentUrl = page.url()
      expect(currentUrl).toMatch(/.*\/search.*/)
    } else {
      // If search input doesn't exist on home page, navigate to search page
      await page.click('text=搜索')
      await expect(page).toHaveURL(/.*\/search/)
      
      // Enter search keyword
      const searchInput = page.locator('input[placeholder*="关键词"]')
      const keyword = '测试'
      await searchInput.fill(keyword)
      
      // Click search button or press Enter
      const searchButton = page.locator('button:has-text("搜索")')
      if (await searchButton.isVisible().catch(() => false)) {
        await searchButton.click()
      } else {
        await searchInput.press('Enter')
      }
      
      // Wait for results
      await page.waitForTimeout(2000)
    }
  })

  test('should filter search results by city', async ({ page }) => {
    // Navigate to search page
    await page.click('text=搜索')
    await expect(page).toHaveURL(/.*\/search/)

    // Enter search keyword
    const searchInput = page.locator('input[placeholder*="关键词"]')
    await searchInput.fill('测试')

    // Select city filter
    const citySelect = page.locator('div.ant-select-selector').first()
    if (await citySelect.isVisible().catch(() => false)) {
      await citySelect.click()
      await page.click('text=北京')
    }

    // Perform search
    const searchButton = page.locator('button:has-text("搜索")')
    if (await searchButton.isVisible().catch(() => false)) {
      await searchButton.click()
    } else {
      await searchInput.press('Enter')
    }

    // Wait for results
    await page.waitForTimeout(2000)

    // Verify results are displayed (or empty state)
    const resultsContainer = page.locator('.search-results-container, .ant-empty')
    await expect(resultsContainer).toBeVisible()
  })

  test('should navigate to post detail from search results', async ({ page }) => {
    // Navigate to search page
    await page.click('text=搜索')
    await expect(page).toHaveURL(/.*\/search/)

    // Perform a search
    const searchInput = page.locator('input[placeholder*="关键词"]')
    await searchInput.fill('测试')
    
    const searchButton = page.locator('button:has-text("搜索")')
    if (await searchButton.isVisible().catch(() => false)) {
      await searchButton.click()
    } else {
      await searchInput.press('Enter')
    }

    // Wait for results
    await page.waitForTimeout(2000)

    // Try to click on first result if available
    const firstResult = page.locator('.search-results-item-card, .post-list-card').first()
    const resultExists = await firstResult.isVisible().catch(() => false)

    if (resultExists) {
      await firstResult.click()
      await page.waitForLoadState('networkidle')

      // Verify we're on detail page
      await expect(page).toHaveURL(/.*\/post\/.*/)

      // Verify post details are displayed
      await expect(page.locator('text=曝光内容')).toBeVisible()
    }
  })

  test('should show empty state when no results found', async ({ page }) => {
    // Navigate to search page
    await page.click('text=搜索')
    await expect(page).toHaveURL(/.*\/search/)

    // Search for a keyword that likely doesn't exist
    const searchInput = page.locator('input[placeholder*="关键词"]')
    await searchInput.fill('这是一个非常不可能存在的关键词123456789')
    
    const searchButton = page.locator('button:has-text("搜索")')
    if (await searchButton.isVisible().catch(() => false)) {
      await searchButton.click()
    } else {
      await searchInput.press('Enter')
    }

    // Wait for results
    await page.waitForTimeout(2000)

    // Verify empty state is shown
    const emptyState = page.locator('.ant-empty, text=未找到')
    await expect(emptyState).toBeVisible()
  })
})

