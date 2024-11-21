from pathlib import Path

RAG_FILE_PII_SCAN = "ragdata.txt"

def readRagData (ragFile: str) -> str:
    with open (ragFile, "r") as ragFileData:
        ragData = ragFileData.read()
        ragFileData.close()
        return ragData

def read_file (file: str) -> bool:
    file_path = Path(file)
    if file_path.exists():
        return True
    return False
