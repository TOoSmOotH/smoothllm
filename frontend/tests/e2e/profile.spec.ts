import { test, expect } from '@playwright/test';

// Create a unique user for profile tests
const randomSuffix = Math.random().toString(36).substring(7);
const testUser = {
  username: `profileuser_${randomSuffix}`,
  email: `profile_${randomSuffix}@example.com`,
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

  // Handle either dashboard redirect or login page
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

test.describe('Profile Management', () => {
  test.describe.configure({ mode: 'serial' });

  test.beforeAll(async ({ browser }) => {
    // Register a new user for profile tests
    const page = await browser.newPage();
    try {
      await registerAndLogin(page);
    } finally {
      await page.close();
    }
  });

  test('should display user profile page', async ({ page }) => {
    // Login
    await loginUser(page);

    // Navigate to profile
    await page.goto('/profile');
    await page.waitForLoadState('networkidle');
    await expect(page).toHaveURL(/\/profile/);

    // Verify profile page elements
    await expect(page.getByText('About')).toBeVisible({ timeout: 10000 });
    await expect(page.getByText('Contact')).toBeVisible();
    await expect(page.getByText('Professional')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Edit Profile' })).toBeVisible();
  });

  test('should navigate to edit profile from profile page', async ({ page }) => {
    // Login
    await loginUser(page);

    // Navigate to profile
    await page.goto('/profile');
    await page.waitForLoadState('networkidle');

    // Click edit profile button
    await page.getByRole('button', { name: 'Edit Profile' }).first().click();

    // Verify navigation to edit page
    await expect(page).toHaveURL(/\/profile\/edit/, { timeout: 10000 });
    await expect(page.getByText('Edit Profile')).toBeVisible();
  });

  test('should save profile changes successfully', async ({ page }) => {
    // Login
    await loginUser(page);

    // Navigate to edit profile
    await page.goto('/profile/edit');
    await page.waitForLoadState('networkidle');
    await expect(page).toHaveURL(/\/profile\/edit/, { timeout: 10000 });

    // Wait for form to be ready
    await expect(page.getByText('Basic Information')).toBeVisible({ timeout: 10000 });

    // Generate unique test data
    const uniqueSuffix = Math.random().toString(36).substring(7);
    const bio = `E2E test bio ${uniqueSuffix}`;

    // Fill bio - use placeholder to find the textarea
    const bioTextarea = page.locator('textarea[placeholder="Tell us about yourself..."]');
    await expect(bioTextarea).toBeVisible({ timeout: 5000 });
    await bioTextarea.clear();
    await bioTextarea.fill(bio);

    // Click save button
    await page.getByRole('button', { name: 'Save Changes' }).click();

    // Wait for success toast to appear (indicates save completed)
    await expect(page.getByText('Profile updated successfully')).toBeVisible({ timeout: 10000 });

    // Wait for network to settle after save
    await page.waitForLoadState('networkidle');

    // Navigate to profile view to verify changes persisted
    await page.goto('/profile');
    await page.waitForLoadState('networkidle');

    // Wait for the About section to load
    await expect(page.getByText('About')).toBeVisible({ timeout: 10000 });

    // Verify the bio text appears on the page
    // The bio should be in a paragraph element under the About section
    await expect(page.getByText(bio, { exact: false })).toBeVisible({ timeout: 15000 });
  });

  test('should navigate from dashboard to profile via quick action', async ({ page }) => {
    // Login
    await loginUser(page);

    // Click View Profile quick action
    await page.getByText('View Profile').click();

    // Verify navigation
    await expect(page).toHaveURL(/\/profile/, { timeout: 10000 });
  });

  test('should navigate from dashboard to edit profile via quick action', async ({ page }) => {
    // Login
    await loginUser(page);

    // Click Edit Profile quick action
    await page.getByText('Edit Profile').click();

    // Verify navigation
    await expect(page).toHaveURL(/\/profile\/edit/, { timeout: 10000 });
  });
});
