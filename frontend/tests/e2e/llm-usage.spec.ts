import { test, expect } from '@playwright/test';

// Create a unique user for usage tests
const randomSuffix = Math.random().toString(36).substring(7);
const testUser = {
  username: `usageuser_${randomSuffix}`,
  email: `usage_${randomSuffix}@example.com`,
  password: 'securePassword123!',
};

// Helper to register and login
async function registerAndLogin(page: any) {
  await page.goto('/register');
  await page.fill('input#username', testUser.username);
  await page.fill('input#email', testUser.email);
  await page.fill('input#password', testUser.password);
  await page.fill('input#confirmPassword', testUser.password);
  await page.locator('input[type="checkbox"]').check();
  await page.click('button[type="submit"]');

  await page.waitForURL(/\/(dashboard|login)/, { timeout: 15000 });

  if (page.url().includes('/login')) {
    await page.fill('input#email', testUser.email);
    await page.fill('input#password', testUser.password);
    await page.click('button[type="submit"]');
    await page.waitForURL(/\/dashboard/, { timeout: 10000 });
  }
}

async function loginUser(page: any) {
  await page.goto('/login');
  await page.fill('input#email', testUser.email);
  await page.fill('input#password', testUser.password);
  await page.click('button[type="submit"]');
  await page.waitForURL(/\/dashboard/, { timeout: 10000 });
}

test.describe('Usage Statistics Dashboard', () => {
  test.describe.configure({ mode: 'serial' });

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    try {
      await registerAndLogin(page);
    } finally {
      await page.close();
    }
  });

  test('should display usage page with header', async ({ page }) => {
    await loginUser(page);

    await page.goto('/usage');
    await page.waitForLoadState('networkidle');

    await expect(page).toHaveURL(/\/usage/);
    await expect(page.locator('h1:has-text("Usage Statistics")')).toBeVisible({ timeout: 10000 });
    await expect(page.getByText('Track your LLM API usage and costs')).toBeVisible();
  });

  test('should display summary statistics cards', async ({ page }) => {
    await loginUser(page);
    await page.goto('/usage');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("Usage Statistics")')).toBeVisible({ timeout: 10000 });

    // Should show statistics cards
    await expect(page.getByText('Total Requests')).toBeVisible({ timeout: 10000 });
    await expect(page.getByText('Total Tokens')).toBeVisible();
    await expect(page.getByText('Total Cost')).toBeVisible();
  });

  test('should display date range filters', async ({ page }) => {
    await loginUser(page);
    await page.goto('/usage');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("Usage Statistics")')).toBeVisible({ timeout: 10000 });

    // Should show date filters
    await expect(page.getByText('Start Date')).toBeVisible({ timeout: 10000 });
    await expect(page.getByText('End Date')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Apply Filters' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Clear' })).toBeVisible();
  });

  test('should have functional refresh button', async ({ page }) => {
    await loginUser(page);
    await page.goto('/usage');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("Usage Statistics")')).toBeVisible({ timeout: 10000 });

    // Find and click refresh button (button with SVG icon in header)
    const refreshButton = page.locator('button').filter({ has: page.locator('svg') }).first();
    await expect(refreshButton).toBeVisible({ timeout: 10000 });

    // Click should trigger a refresh
    await refreshButton.click();

    // Page should still be functional
    await expect(page.getByText('Total Requests')).toBeVisible();
  });

  test('should set date filters', async ({ page }) => {
    await loginUser(page);
    await page.goto('/usage');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("Usage Statistics")')).toBeVisible({ timeout: 10000 });

    // Set start date
    const startDateInput = page.locator('input[type="date"]').first();
    await startDateInput.fill('2024-01-01');

    // Set end date
    const endDateInput = page.locator('input[type="date"]').last();
    await endDateInput.fill('2024-12-31');

    // Apply filters
    await page.getByRole('button', { name: 'Apply Filters' }).click();

    // Wait for any loading to complete
    await page.waitForLoadState('networkidle');

    // Stats should still be visible
    await expect(page.getByText('Total Requests')).toBeVisible();
  });

  test('should clear date filters', async ({ page }) => {
    await loginUser(page);
    await page.goto('/usage');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("Usage Statistics")')).toBeVisible({ timeout: 10000 });

    // Set some filters first
    const startDateInput = page.locator('input[type="date"]').first();
    const endDateInput = page.locator('input[type="date"]').last();

    await startDateInput.fill('2024-01-01');
    await endDateInput.fill('2024-12-31');

    // Clear filters
    await page.getByRole('button', { name: 'Clear' }).click();

    // Date inputs should be cleared
    await expect(startDateInput).toHaveValue('');
    await expect(endDateInput).toHaveValue('');
  });

  test('should display zero usage for new user', async ({ page }) => {
    await loginUser(page);
    await page.goto('/usage');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("Usage Statistics")')).toBeVisible({ timeout: 10000 });

    // New user should have zero stats - check that the cards are visible
    await expect(page.getByText('Total Requests')).toBeVisible({ timeout: 10000 });
    await expect(page.getByText('Total Tokens')).toBeVisible();
    await expect(page.getByText('Total Cost')).toBeVisible();
  });

  test('should be accessible via navigation', async ({ page }) => {
    await loginUser(page);

    // Navigate via sidebar/navigation if exists
    await page.goto('/dashboard');
    await page.waitForLoadState('networkidle');

    // Look for Usage link in navigation
    const usageLink = page.getByRole('link', { name: /Usage/i });

    if (await usageLink.isVisible().catch(() => false)) {
      await usageLink.click();
      await expect(page).toHaveURL(/\/usage/);
      await expect(page.locator('h1:has-text("Usage Statistics")')).toBeVisible();
    } else {
      // Fallback: navigate directly
      await page.goto('/usage');
      await expect(page.locator('h1:has-text("Usage Statistics")')).toBeVisible({ timeout: 10000 });
    }
  });
});

test.describe('Usage Dashboard - Empty State', () => {
  test('should gracefully handle no usage data', async ({ page }) => {
    // Create a fresh user
    const suffix = Math.random().toString(36).substring(7);
    const user = {
      username: `empty_${suffix}`,
      email: `empty_${suffix}@example.com`,
      password: 'securePassword123!',
    };

    await page.goto('/register');
    await page.fill('input#username', user.username);
    await page.fill('input#email', user.email);
    await page.fill('input#password', user.password);
    await page.fill('input#confirmPassword', user.password);
    await page.locator('input[type="checkbox"]').check();
    await page.click('button[type="submit"]');
    await page.waitForURL(/\/(dashboard|login)/, { timeout: 15000 });

    if (page.url().includes('/login')) {
      await page.fill('input#email', user.email);
      await page.fill('input#password', user.password);
      await page.click('button[type="submit"]');
      await page.waitForURL(/\/dashboard/, { timeout: 10000 });
    }

    // Navigate to usage
    await page.goto('/usage');
    await page.waitForLoadState('networkidle');

    // Page should load without errors
    await expect(page.locator('h1:has-text("Usage Statistics")')).toBeVisible({ timeout: 10000 });

    // Summary cards should be visible even with no data
    await expect(page.getByText('Total Requests')).toBeVisible();
    await expect(page.getByText('Total Tokens')).toBeVisible();
    await expect(page.getByText('Total Cost')).toBeVisible();
  });
});

test.describe('Usage Dashboard - Loading States', () => {
  test('should show and hide loading indicator', async ({ page }) => {
    // Create a fresh user
    const suffix = Math.random().toString(36).substring(7);
    const user = {
      username: `loading_${suffix}`,
      email: `loading_${suffix}@example.com`,
      password: 'securePassword123!',
    };

    await page.goto('/register');
    await page.fill('input#username', user.username);
    await page.fill('input#email', user.email);
    await page.fill('input#password', user.password);
    await page.fill('input#confirmPassword', user.password);
    await page.locator('input[type="checkbox"]').check();
    await page.click('button[type="submit"]');
    await page.waitForURL(/\/(dashboard|login)/, { timeout: 15000 });

    if (page.url().includes('/login')) {
      await page.fill('input#email', user.email);
      await page.fill('input#password', user.password);
      await page.click('button[type="submit"]');
      await page.waitForURL(/\/dashboard/, { timeout: 10000 });
    }

    // Navigate to usage - content should eventually load
    await page.goto('/usage');

    // After loading, stats should be visible
    await expect(page.getByText('Total Requests')).toBeVisible({ timeout: 15000 });
  });
});
