import { test, expect } from '@playwright/test';

test.describe('Current Week App', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('has correct title', async ({ page }) => {
    await expect(page).toHaveTitle(/Current week is/);
  });

  test('navigates to next week', async ({ page }) => {
    const weekElement = page.locator('#week');
    const initialWeekText = await weekElement.innerText();

    // Click next week
    await page.locator('#nextWeekButton').click();

    // Verify week changed
    await expect(weekElement).not.toHaveText(initialWeekText);

    // Verify URL change if applicable, but the app uses fetch/AJAX.
    // The week number should change.
  });

  test('navigates to previous week', async ({ page }) => {
    const weekElement = page.locator('#week');
    const initialWeekText = await weekElement.innerText();

    // Click prev week
    await page.locator('#prevWeekButton').click();

    // Verify week changed
    await expect(weekElement).not.toHaveText(initialWeekText);
  });

  test('resets to current week', async ({ page }) => {
    const weekElement = page.locator('#week');
    const initialWeekText = await weekElement.innerText();

    // Navigate away first
    await page.locator('#nextWeekButton').click();
    await expect(weekElement).not.toHaveText(initialWeekText);

    // Click week number to reset
    await weekElement.click();

    // Verify it's back to initial
    await expect(weekElement).toHaveText(initialWeekText);
  });

  test('resets to current week on ArrowUp', async ({ page }) => {
    const weekElement = page.locator('#week');
    const initialWeekText = await weekElement.innerText();

    // Navigate away first
    await page.locator('#nextWeekButton').click();
    await expect(weekElement).not.toHaveText(initialWeekText);

    // Press ArrowUp to reset
    await page.keyboard.press('ArrowUp');

    // Verify it's back to initial
    await expect(weekElement).toHaveText(initialWeekText);
  });
});
