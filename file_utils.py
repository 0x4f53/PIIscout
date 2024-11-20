RAG_FILE_PII_SCAN = "ragdata.txt"

def readRagData (ragFile: str) -> str:
    with open (ragFile, "r") as ragFileData:
        ragData = ragFileData.read()
        ragFileData.close()
        return ragData
