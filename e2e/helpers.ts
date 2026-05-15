import { type Page } from '@playwright/test';

// Creates a game room and clicks through to the waiting state (WS connected, waiting for P2).
// Returns the 4-character room code.
export async function createAndEnterGame(page: Page): Promise<string> {
  await page.goto('/');
  await page.click('#createBtn');
  await page.waitForSelector('#displayCode');
  const code = (await page.textContent('#displayCode'))!.trim();
  await page.click('#enterBtn');
  return code;
}

// Joins an existing game room (WS connects immediately after the POST succeeds).
export async function joinGame(page: Page, code: string): Promise<void> {
  await page.goto('/');
  await page.fill('#codeInput', code);
  await page.click('#joinBtn');
}

// Waits until the game is active: no "Waiting for opponent" overlay and HUD shows Gold.
export async function waitForGameActive(page: Page): Promise<void> {
  await page.waitForFunction(
    () => {
      const hud = document.getElementById('hud');
      return hud !== null && (hud.textContent?.includes('Gold') ?? false);
    },
    { timeout: 10_000 },
  );
}
