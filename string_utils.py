import re

def extract_substring(text: str) -> str:
    start_marker = "Answer\njson"
    end_marker = "Share\nRewrite"
    
    start_index = text.find(start_marker)
    end_index = text.find(end_marker)
    
    if start_index == -1 or end_index == -1 or start_index > end_index:
        return None  # Return None if markers are not found or out of order
    
    start_index += len(start_marker)  # Move past the start marker
    substring = text[start_index:end_index].strip()  # Extract and strip any surrounding whitespace
    
    return substring.strip()