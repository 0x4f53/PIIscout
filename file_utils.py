from pathlib import Path
import os

RAG_FILE_PII_SCAN = "ragdata.txt"

def readRagData (ragFile: str) -> str:
    with open (ragFile, "r") as ragFileData:
        ragData = ragFileData.read()
        ragFileData.close()
        return ragData
    
def fullPath (file: str) -> str:
    return os.path.abspath(file)

def read_file (file: str) -> bool:
    file_path = Path(file)
    if file_path.exists():
        return True
    return False

SIGN_IN_REQUIRED=".signInRequired"

def existsSignInRequiredFile():
    return os.path.isfile(SIGN_IN_REQUIRED)

def makeSignInRequiredFile():
    with open(SIGN_IN_REQUIRED, "w") as signInFile:
        signInFile.write("")
        signInFile.close()

def removeSignInRequiredFile():
    os.remove(SIGN_IN_REQUIRED)