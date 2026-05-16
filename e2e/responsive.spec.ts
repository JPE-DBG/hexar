import { test, expect, type Browser } from '@playwright/test';
import { createAndEnterGame, joinGame, waitForGameActive } from './helpers';

async function setupGame(browser: Browser, viewport: { width: number; height: number }) {
  const ctx1 = await browser.newContext({ viewport });
  const ctx2 = await browser.newContext({ viewport });
  const p1 = await ctx1.newPage();
  const p2 = await ctx2.newPage();
  const code = await createAndEnterGame(p1);
  await joinGame(p2, code);
  await waitForGameActive(p1);
  await waitForGameActive(p2);
  return { ctx1, ctx2, p1, p2 };
}

test.describe('Desktop layout (1280×800)', () => {
  test('sidebar is visible on the left edge', async ({ browser }) => {
    const { ctx1, ctx2, p1 } = await setupGame(browser, { width: 1280, height: 800 });
    try {
      const sidebar = p1.locator('#sidebar');
      await expect(sidebar).toBeVisible();

      // Sidebar is on the left — its left edge should be at/near 0
      const box = await sidebar.boundingBox();
      expect(box).not.toBeNull();
      expect(box!.x).toBeLessThan(10);
    } finally {
      await ctx1.close();
      await ctx2.close();
    }
  });

  test('destructive panel is not shown at desktop resolution', async ({ browser }) => {
    const { ctx1, ctx2, p1 } = await setupGame(browser, { width: 1280, height: 800 });
    try {
      const destructivePanel = p1.locator('#destructive-panel');
      // Should exist in DOM but not be visible (CSS hides it on desktop)
      await expect(destructivePanel).toBeHidden();
    } finally {
      await ctx1.close();
      await ctx2.close();
    }
  });
});

test.describe('Mobile portrait layout (390×844)', () => {
  test('sidebar moves to bottom bar', async ({ browser }) => {
    const { ctx1, ctx2, p1 } = await setupGame(browser, { width: 390, height: 844 });
    try {
      const sidebar = p1.locator('#sidebar');
      await expect(sidebar).toBeVisible();

      const box = await sidebar.boundingBox();
      expect(box).not.toBeNull();

      // In mobile portrait, sidebar is at the bottom — its bottom edge should be near screen bottom
      const bottomEdge = box!.y + box!.height;
      expect(bottomEdge).toBeGreaterThan(800); // near bottom of 844px screen
    } finally {
      await ctx1.close();
      await ctx2.close();
    }
  });

  test('destructive panel appears at top-left when context is active', async ({ browser }) => {
    const { ctx1, ctx2, p1 } = await setupGame(browser, { width: 390, height: 844 });
    try {
      // Click the canvas to try to select an owned hex
      // Player 1 starts at top-left spawn — click near center to find a hex
      const canvas = p1.locator('canvas#game');
      await expect(canvas).toBeVisible();

      // Click in the middle of the screen to trigger hex selection
      // We need to trigger a hex click that results in an owned hex with context
      await canvas.click({ position: { x: 195, y: 422 } });
      await p1.waitForTimeout(200);

      // If a hex with context is selected, #destructive-panel.visible appears at top-left
      // Check that if it's visible, it's positioned at the top-left (not bottom)
      const destructivePanel = p1.locator('#destructive-panel');
      const isVisible = await destructivePanel.isVisible();

      if (isVisible) {
        const box = await destructivePanel.boundingBox();
        expect(box).not.toBeNull();
        // Top-left: x near 0, y near 0
        expect(box!.x).toBeLessThan(50);
        expect(box!.y).toBeLessThan(150);
      }
    } finally {
      await ctx1.close();
      await ctx2.close();
    }
  });
});

test.describe('Mobile landscape layout (896×414)', () => {
  test('sidebar is visible in landscape mode', async ({ browser }) => {
    const { ctx1, ctx2, p1 } = await setupGame(browser, { width: 896, height: 414 });
    try {
      await expect(p1.locator('#sidebar')).toBeVisible();
    } finally {
      await ctx1.close();
      await ctx2.close();
    }
  });

  test('HUD is visible in landscape mode', async ({ browser }) => {
    const { ctx1, ctx2, p1 } = await setupGame(browser, { width: 896, height: 414 });
    try {
      await expect(p1.locator('#hud')).toBeVisible();
    } finally {
      await ctx1.close();
      await ctx2.close();
    }
  });
});
