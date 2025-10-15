import time
from playwright.sync_api import sync_playwright, expect

def run_verification(playwright):
    browser = playwright.chromium.launch(headless=True)
    page = browser.new_page()

    try:
        # Navigate to the homepage
        page.goto("http://localhost:8080")

        # Fill in the form with the example data
        page.click(".example-item:has-text('贵州茅台')")

        # Click the analyze button
        page.click("#analyzeBtn")

        # Wait for the status text to indicate that the analysis is running
        status_text_locator = page.locator("#statusText")
        expect(status_text_locator).to_contain_text("正在启动分析...", timeout=10000)

        # Wait for a status update from the callback handler
        # expect(status_text_locator).to_contain_text("开始执行任务", timeout=20000)
        time.sleep(10) # wait for 10s
        page.pause()
        # Take a screenshot of the progress
        page.screenshot(path="jules-scratch/verification/verification.png")

    finally:
        browser.close()

with sync_playwright() as playwright:
    run_verification(playwright)