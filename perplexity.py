import time, file_utils, string_utils, argparse, os, json

parser = argparse.ArgumentParser(description="A script to demonstrate CLI arguments.")
parser.add_argument("file_name", help="The filename to scan (supported: PDF, JPG, BMP, PNG, DOCX, TXT, XLS)")
args = parser.parse_args()

from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager
import undetected_chromedriver as uc
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.common.by import By
from selenium.common.exceptions import NoSuchElementException
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.action_chains import ActionChains

from seleniumbase import Driver

url = "https://perplexity.ai"

# configure driver
headless_mode = True
if file_utils.existsSignInRequiredFile():
    headless_mode = False

options = webdriver.ChromeOptions()
options.add_argument("--lang=en")
driver = Driver(uc=True, headless=headless_mode, user_data_dir="./.chrome-configs/")

start_time = time.time()
elapsed_time = time.time() - start_time

# 0. start page
driver.get(url)

# 0a. check if the user needs to sign in
elements = driver.find_elements(By.XPATH, "//div[text()='Log in']")
if elements:
    elements[0].click()
    print ("You need to sign in. Please restart this program and sign in using your credentials...")
    
    if not file_utils.existsSignInRequiredFile():
        file_utils.makeSignInRequiredFile()
        driver.quit()
        exit(-1)

    try:
        while True:
            try:
                element = driver.find_element(By.XPATH, '//img[@alt="User avatar"]')
                if element: break
            except NoSuchElementException: time.sleep(1)

    finally: 
        time.sleep(3)
        file_utils.removeSignInRequiredFile()
        driver.quit()
        print ("Signed in successfully! Please restart this program and sign in using your credentials...")
        exit(-1)
            

# 1. find the upload file button
file_input = driver.find_element(By.XPATH, "//input[@type='file']")

file_path = file_utils.fullPath(args.file_name)
file_input.send_keys(file_path)

# 2. wait for file upload to be completed and for the button to be activated
button = driver.find_element(By.XPATH, "//button[@aria-label='Submit']")

def is_button_disabled(button):
    return button.get_attribute("disabled") is not None

while is_button_disabled(button):
    elapsed_time = time.time() - start_time
    formatted_time = time.strftime("%H:%M:%S", time.gmtime(elapsed_time))
    print(f"Processing file... (elapsed time: {formatted_time})")
    time.sleep(1)

    textarea = driver.find_element(By.TAG_NAME, "textarea")
    textarea.clear()
    textarea.send_keys(file_utils.readRagData(file_utils.RAG_FILE_PII_SCAN))

    # 3. Click the button
    button = driver.find_element(By.XPATH, "//button[@aria-label='Submit']")  # Re-locate the button if necessary
    button.click()

# 4. wait for text to be generated completely
text = driver.execute_script("return document.body.innerText;")

while "Share\nRewrite" not in text:
    time.sleep(1)
    text = driver.execute_script("return document.body.innerText;")

# 5. delete conversation
button = driver.find_element(By.XPATH, "//button[@data-testid='thread-dropdown-menu']")  # Re-locate the button if necessary
button.click()

button = driver.find_element(By.XPATH, "//div[@data-testid='thread-delete']")  # Re-locate the button if necessary
button.click()

element = driver.find_element(By.XPATH, "//div[text()='Confirm']")
element.click()

# 6. kill driver
driver.quit()

text = string_utils.extract_substring(text)

data_dict = json.loads(text)
data_dict["file_metadata"] = {
    "file_path" : file_path
    "file_url" : file_url
}
file_utils.write_output(f"output/{string_utils.print_current_timestamp()}.json", json.dumps(data_dict))

print(text)

