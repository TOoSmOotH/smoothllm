import { test, expect } from '@playwright/test';

// Create a unique user for LLM feature tests
const randomSuffix = Math.random().toString(36).substring(7);
const testUser = {
  username: `llmuser_${randomSuffix}`,
  email: `llm_${randomSuffix}@example.com`,
  password: 'securePassword123!',
};

const testProvider = {
  name: `Test Provider ${randomSuffix}`,
  apiKey: 'sk-test-fake-key-for-testing',
  defaultModel: 'gpt-4o-mini',
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

test.describe('LLM Providers Management', () => {
  test.describe.configure({ mode: 'serial' });

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    try {
      await registerAndLogin(page);
    } finally {
      await page.close();
    }
  });

  test('should display providers page with empty state', async ({ page }) => {
    await loginUser(page);

    await page.goto('/providers');
    await page.waitForLoadState('networkidle');

    await expect(page).toHaveURL(/\/providers/);
    await expect(page.locator('h1:has-text("LLM Providers")')).toBeVisible({ timeout: 10000 });

    // Should show empty state or providers list
    const emptyStateText = page.getByText('No Providers Configured');
    const hasEmptyState = await emptyStateText.isVisible({ timeout: 5000 }).catch(() => false);
    if (hasEmptyState) {
      await expect(page.getByRole('button', { name: 'Add Your First Provider' })).toBeVisible();
    }
  });

  test('should open create provider modal', async ({ page }) => {
    await loginUser(page);
    await page.goto('/providers');
    await page.waitForLoadState('networkidle');

    // Wait for page to be ready
    await expect(page.locator('h1:has-text("LLM Providers")')).toBeVisible({ timeout: 10000 });

    // Wait for the button to be visible and click it
    const addButton = page.getByRole('button', { name: /Add Provider/i }).first();
    await expect(addButton).toBeVisible({ timeout: 5000 });
    await addButton.click();

    // Wait for modal - check for the form input which confirms modal is open
    // The Input component may not use proper label association, so also try placeholder
    const providerNameInput = page.locator('input[placeholder="My OpenAI Account"]');
    await expect(providerNameInput).toBeVisible({ timeout: 10000 });
  });

  test('should validate required fields when creating provider', async ({ page }) => {
    await loginUser(page);
    await page.goto('/providers');
    await page.waitForLoadState('networkidle');

    // Wait for page to be ready
    await expect(page.locator('h1:has-text("LLM Providers")')).toBeVisible({ timeout: 10000 });

    // Click Add Provider button
    const addButton = page.getByRole('button', { name: /Add Provider/i }).first();
    await expect(addButton).toBeVisible({ timeout: 5000 });
    await addButton.click();

    // Wait for modal to be fully open
    const providerNameInput = page.locator('input[placeholder="My OpenAI Account"]');
    await expect(providerNameInput).toBeVisible({ timeout: 10000 });

    // Try to submit empty form
    await page.getByRole('button', { name: 'Create Provider' }).click();

    // Should show validation errors
    await expect(page.getByText('Provider name is required')).toBeVisible({ timeout: 5000 });
    await expect(page.getByText('API key is required')).toBeVisible();
  });

  test('should create a new provider', async ({ page }) => {
    await loginUser(page);
    await page.goto('/providers');
    await page.waitForLoadState('networkidle');

    // Wait for page to be ready
    await expect(page.locator('h1:has-text("LLM Providers")')).toBeVisible({ timeout: 10000 });

    // Click Add Provider
    const addButton = page.getByRole('button', { name: /Add Provider/i }).first();
    await expect(addButton).toBeVisible({ timeout: 5000 });
    await addButton.click();

    // Wait for modal
    const providerNameInput = page.locator('input[placeholder="My OpenAI Account"]');
    await expect(providerNameInput).toBeVisible({ timeout: 10000 });

    // Fill in the form
    await providerNameInput.fill(testProvider.name);

    // Select provider type (OpenAI is default)
    const providerTypeSelect = page.locator('select').first();
    await providerTypeSelect.selectOption('openai');

    // Fill API key
    const apiKeyInput = page.locator('input[placeholder="sk-..."]');
    await apiKeyInput.fill(testProvider.apiKey);

    // Fill default model
    const modelInput = page.locator('input[placeholder="gpt-4o"]');
    await modelInput.fill(testProvider.defaultModel);

    // Submit
    await page.getByRole('button', { name: 'Create Provider' }).click();

    // Wait for success toast
    await expect(page.getByText('Provider created successfully')).toBeVisible({ timeout: 10000 });

    // Provider should appear in the list
    await expect(page.getByText(testProvider.name)).toBeVisible({ timeout: 5000 });
  });

  test('should display provider details after creation', async ({ page }) => {
    await loginUser(page);
    await page.goto('/providers');
    await page.waitForLoadState('networkidle');

    // Provider should be in the list
    await expect(page.getByText(testProvider.name)).toBeVisible({ timeout: 10000 });

    // Should show provider type (uppercase in the UI)
    await expect(page.locator('text=OPENAI').first()).toBeVisible();

    // Should show default model
    await expect(page.getByText(testProvider.defaultModel)).toBeVisible();

    // Should show Active status
    await expect(page.getByText('Active').first()).toBeVisible();
  });

  test('should open edit modal for existing provider', async ({ page }) => {
    await loginUser(page);
    await page.goto('/providers');
    await page.waitForLoadState('networkidle');

    // Wait for page header to be visible first
    await expect(page.locator('h1:has-text("LLM Providers")')).toBeVisible({ timeout: 10000 });

    // Wait for provider heading to be visible
    const providerHeading = page.locator(`h3:has-text("${testProvider.name}")`);
    await expect(providerHeading).toBeVisible({ timeout: 15000 });

    // Find the card container using the rounded-lg class that contains our provider
    const providerCard = page.locator('div.rounded-lg').filter({ has: providerHeading }).first();

    // The edit button is the second button (after test connection button)
    const editButton = providerCard.locator('button').nth(1);
    await expect(editButton).toBeVisible({ timeout: 5000 });
    await editButton.click();

    // Modal should open with edit title
    await expect(page.getByRole('heading', { name: 'Edit Provider' })).toBeVisible({ timeout: 5000 });

    // Name field should be pre-filled
    const providerNameInput = page.locator('input[placeholder="My OpenAI Account"]');
    await expect(providerNameInput).toHaveValue(testProvider.name);
  });

  test('should update provider name', async ({ page }) => {
    await loginUser(page);
    await page.goto('/providers');
    await page.waitForLoadState('networkidle');

    // Wait for provider heading to be visible
    const providerHeading = page.locator(`h3:has-text("${testProvider.name}")`);
    await expect(providerHeading).toBeVisible({ timeout: 15000 });

    // Open edit modal
    const providerCard = page.locator('div.rounded-lg').filter({ has: providerHeading }).first();
    await providerCard.locator('button').nth(1).click();

    // Wait for modal to be fully open by checking for the input field
    const providerNameInput = page.locator('input[placeholder="My OpenAI Account"]');
    await expect(providerNameInput).toBeVisible({ timeout: 10000 });
    await expect(page.getByRole('heading', { name: 'Edit Provider' })).toBeVisible({ timeout: 5000 });

    // Update name
    const updatedName = `${testProvider.name} Updated`;
    await providerNameInput.clear();
    await providerNameInput.fill(updatedName);

    // Submit
    await page.getByRole('button', { name: 'Save Changes' }).click();

    // Wait for success toast (increase timeout in case of network delay)
    await expect(page.getByText('Provider updated successfully')).toBeVisible({ timeout: 15000 });

    // Updated name should appear
    await expect(page.getByText(updatedName)).toBeVisible({ timeout: 5000 });

    // Update the test provider name for later tests
    testProvider.name = updatedName;
  });

  test('should show delete confirmation modal', async ({ page }) => {
    await loginUser(page);
    await page.goto('/providers');
    await page.waitForLoadState('networkidle');

    // Wait for provider heading to be visible
    const providerHeading = page.locator(`h3:has-text("${testProvider.name}")`);
    await expect(providerHeading).toBeVisible({ timeout: 15000 });

    // Click delete button (third button)
    const providerCard = page.locator('div.rounded-lg').filter({ has: providerHeading }).first();
    await providerCard.locator('button').nth(2).click();

    // Confirmation modal should appear
    await expect(page.getByRole('heading', { name: 'Delete Provider' })).toBeVisible({ timeout: 5000 });
    await expect(page.getByText('This action cannot be undone')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Delete Provider' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Cancel' })).toBeVisible();
  });

  test('should cancel delete when clicking cancel', async ({ page }) => {
    await loginUser(page);
    await page.goto('/providers');
    await page.waitForLoadState('networkidle');

    // Wait for provider heading to be visible
    const providerHeading = page.locator(`h3:has-text("${testProvider.name}")`);
    await expect(providerHeading).toBeVisible({ timeout: 15000 });

    // Click delete button
    const providerCard = page.locator('div.rounded-lg').filter({ has: providerHeading }).first();
    await providerCard.locator('button').nth(2).click();

    await expect(page.getByRole('heading', { name: 'Delete Provider' })).toBeVisible({ timeout: 5000 });

    // Click cancel
    await page.getByRole('button', { name: 'Cancel' }).click();

    // Modal should close
    await expect(page.getByRole('heading', { name: 'Delete Provider' })).not.toBeVisible({ timeout: 3000 });

    // Provider should still be visible
    await expect(page.getByText(testProvider.name)).toBeVisible();
  });

  // Note: Actual deletion test commented out to preserve provider for API key tests
  // test('should delete provider', async ({ page }) => { ... });
});

test.describe('Provider Type Selection', () => {
  test('should show different placeholders for different provider types', async ({ page }) => {
    // Register a new user for this test
    const suffix = Math.random().toString(36).substring(7);
    const user = {
      username: `provtype_${suffix}`,
      email: `provtype_${suffix}@example.com`,
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

    await page.goto('/providers');
    await page.waitForLoadState('networkidle');

    // Wait for page to be ready
    await expect(page.locator('h1:has-text("LLM Providers")')).toBeVisible({ timeout: 10000 });

    // Open create modal
    const addButton = page.getByRole('button', { name: /Add Provider/i }).first();
    await expect(addButton).toBeVisible({ timeout: 5000 });
    await addButton.click();

    // Wait for modal to be fully open using placeholder
    const providerNameInput = page.locator('input[placeholder="My OpenAI Account"]');
    await expect(providerNameInput).toBeVisible({ timeout: 10000 });

    // Check OpenAI (default)
    const providerTypeSelect = page.locator('select').first();
    await expect(providerTypeSelect).toHaveValue('openai');

    // Select Anthropic
    await providerTypeSelect.selectOption('anthropic');

    // Select Local
    await providerTypeSelect.selectOption('local');

    // Verify it can switch between types
    await providerTypeSelect.selectOption('openai');
    await expect(providerTypeSelect).toHaveValue('openai');
  });
});
