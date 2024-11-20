import time, file_utils, string_utils

from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager
import undetected_chromedriver as uc
from selenium.webdriver.common.by import By

from seleniumbase import Driver

options = webdriver.ChromeOptions()
options.add_argument("--lang=en")

#driver = uc.Chrome(service=Service(ChromeDriverManager().install()), options=options) # Create a driver instance using undetected_chromedriver
driver = Driver(uc=True, headless=False, user_data_dir="./.chrome-configs/")

url = "https://perplexity.ai"

driver.get(url)

file_input = driver.find_element(By.XPATH, "//input[@type='file']")

file_path = "/home/owais/Desktop/research/PIIscout/donotcommit/icici.jpeg"  # Replace with the actual file path
file_input.send_keys(file_path)

button = driver.find_element(By.XPATH, "//button[@aria-label='Submit']")

def is_button_disabled(button):
    return button.get_attribute("disabled") is not None

while is_button_disabled(button):
    print("Button is disabled, waiting for it to be enabled...")
    time.sleep(1)
    textarea = driver.find_element(By.TAG_NAME, "textarea")
    textarea.clear()
    textarea.send_keys(file_utils.readRagData(file_utils.RAG_FILE_PII_SCAN))
    button = driver.find_element(By.XPATH, "//button[@aria-label='Submit']")  # Re-locate the button if necessary
    button.click()

text = driver.execute_script("return document.body.innerText;")

while "Share\nRewrite" not in text:
    # print("Waiting for 'Share\\nRewrite' to appear in the text...")
    time.sleep(1)
    text = driver.execute_script("return document.body.innerText;")

# print("Found 'Share\\nRewrite' in the text. Quitting the driver.")
driver.quit()

print(string_utils.extract_substring(text))
