import { test, expect } from '@playwright/test';

// Use a random suffix to avoid user conflicts in persistent DBs during testing
const randomSuffix = Math.random().toString(36).substring(7);
const testUser = {
  username: `testuser_${randomSuffix}`,
  email: `test_${randomSuffix}@example.com`,
  password: 'securePassword123!',
};

test.describe('Authentication Flow', () => {
  // Run tests serially because they depend on the same user state
  test.describe.configure({ mode: 'serial' });

  test('should allow a user to register', async ({ page }) => {
    await page.goto('/register');

    // Fill in the registration form
    await page.fill('input#username', testUser.username);
    await page.fill('input#email', testUser.email);
    await page.fill('input#password', testUser.password);
    await page.fill('input#confirmPassword', testUser.password);

    // Check the terms agreement checkbox
    const termsCheckbox = page.locator('input[type="checkbox"]');
    await termsCheckbox.check();

    // Click the submit button
    await page.click('button[type="submit"]');

    // Wait for navigation to dashboard or login depending on approval logic
    // The backend seems to default to active status now, so likely dashboard
    await expect(page).toHaveURL(/\/dashboard/);

    // Verify successful login/registration state
    await expect(page.getByText(/Welcome|Dashboard/i)).toBeVisible();
  });

  test('should allow a user to login', async ({ page }) => {
    // Navigate to login page
    await page.goto('/login');

    // Fill in credentials
    await page.fill('input#email', testUser.email);
    await page.fill('input#password', testUser.password);

    // Submit
    await page.click('button[type="submit"]');

    // Verify redirect to dashboard
    await expect(page).toHaveURL(/\/dashboard/);

    // Verify user info is potentially visible or at least the dashboard loaded
    await expect(page.locator('#app').first()).toBeVisible();
  });

  test('should show validation errors on invalid login', async ({ page }) => {
    await page.goto('/login');

    await page.fill('input#email', 'wrong@example.com');
    await page.fill('input#password', 'wrongpassword');
    await page.click('button[type="submit"]');

    // Expect an error message to appear
    await expect(page.locator('.text-error-500')).toBeVisible();
  });

  test('should allow a user to logout', async ({ page }) => {
    // First login
    await page.goto('/login');
    await page.fill('input#email', testUser.email);
    await page.fill('input#password', testUser.password);
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/\/dashboard/, { timeout: 10000 });

    // Wait for the page to fully load
    await page.waitForLoadState('networkidle');

    // Click the Sign Out button
    await page.getByRole('button', { name: 'Sign Out' }).click();

    // Verify redirect to login page
    await expect(page).toHaveURL(/\/login/, { timeout: 10000 });

    // Verify user is no longer authenticated (Sign In link visible in header)
    await expect(page.getByRole('link', { name: 'Sign In' })).toBeVisible({ timeout: 5000 });
  });
});

test.describe('Route Protection', () => {
  test('should redirect unauthenticated user to login when accessing protected route', async ({ page }) => {
    // Try to access dashboard without being logged in
    await page.goto('/dashboard');

    // Should be redirected to login with redirect query param
    await expect(page).toHaveURL(/\/login/);
    await expect(page).toHaveURL(/redirect/);
  });

  test('should redirect authenticated user away from guest routes', async ({ page }) => {
    // First register/login a new user
    const suffix = Math.random().toString(36).substring(7);
    const user = {
      username: `guesttest_${suffix}`,
      email: `guesttest_${suffix}@example.com`,
      password: 'securePassword123!',
    };

    await page.goto('/register');
    await page.fill('input#username', user.username);
    await page.fill('input#email', user.email);
    await page.fill('input#password', user.password);
    await page.fill('input#confirmPassword', user.password);
    await page.locator('input[type="checkbox"]').check();
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/\/dashboard/);

    // Now try to access login page while authenticated
    await page.goto('/login');

    // Should be redirected to dashboard
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('should redirect non-admin user away from admin routes', async ({ page }) => {
    // Create a fresh non-admin user for this test
    const suffix = Math.random().toString(36).substring(7);
    const nonAdminUser = {
      username: `nonadmin_${suffix}`,
      email: `nonadmin_${suffix}@example.com`,
      password: 'securePassword123!',
    };

    // Register new user (won't be admin since not first user)
    await page.goto('/register');
    await page.fill('input#username', nonAdminUser.username);
    await page.fill('input#email', nonAdminUser.email);
    await page.fill('input#password', nonAdminUser.password);
    await page.fill('input#confirmPassword', nonAdminUser.password);
    await page.locator('input[type="checkbox"]').check();
    await page.click('button[type="submit"]');

    // Wait for registration to complete
    await page.waitForURL(/\/(dashboard|login)/, { timeout: 10000 });

    // If redirected to login (pending approval), login
    if (page.url().includes('/login')) {
      await page.fill('input#email', nonAdminUser.email);
      await page.fill('input#password', nonAdminUser.password);
      await page.click('button[type="submit"]');
      await expect(page).toHaveURL(/\/dashboard/, { timeout: 10000 });
    }

    // Try to access admin dashboard
    await page.goto('/admin');

    // Should be redirected to user dashboard (not admin)
    await expect(page).toHaveURL(/\/dashboard/, { timeout: 5000 });
  });

  test('should persist session after page refresh', async ({ page }) => {
    // Register a fresh user for this test
    const suffix = Math.random().toString(36).substring(7);
    const sessionUser = {
      username: `session_${suffix}`,
      email: `session_${suffix}@example.com`,
      password: 'securePassword123!',
    };

    await page.goto('/register');
    await page.fill('input#username', sessionUser.username);
    await page.fill('input#email', sessionUser.email);
    await page.fill('input#password', sessionUser.password);
    await page.fill('input#confirmPassword', sessionUser.password);
    await page.locator('input[type="checkbox"]').check();
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/\/dashboard/, { timeout: 10000 });

    // Verify we're logged in (Sign Out button visible)
    await expect(page.getByRole('button', { name: 'Sign Out' })).toBeVisible({ timeout: 5000 });

    // Refresh the page
    await page.reload();
    await page.waitForLoadState('networkidle');

    // Should still be on dashboard and logged in (not redirected to login)
    await expect(page).toHaveURL(/\/dashboard/);
    // Sign Out should still be visible, confirming session persisted
    await expect(page.getByRole('button', { name: 'Sign Out' })).toBeVisible({ timeout: 5000 });
  });
});
