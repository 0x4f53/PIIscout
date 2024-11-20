from bs4 import BeautifulSoup

# Your HTML content
html_content = '''
<code style="font-size: inherit; font-family: inherit; line-height: 1.66667; padding: 8px; white-space: pre-wrap;">
    <span style="opacity: 1; font-size: inherit; line-height: 1.42857; color: rgb(197, 200, 198); background-color: transparent; flex-shrink: 0; padding: 8px; text-align: right; user-select: none;">
        <span class="token token punctuation">{</span><span>
        </span></span><span style="opacity: 1; font-size: inherit; line-height: 1.42857; color: rgb(197, 200, 198); background-color: transparent; flex-shrink: 0; padding: 8px; text-align: right; user-select: none;">
        <span>    </span><span class="token token property">"type_of_document"</span><span class="token token operator">:</span><span> </span><span class="token token" style="color: rgb(181, 189, 104);">"Bank Response Letter"</span><span class="token token punctuation">,</span><span>
        </span></span>
    </span>
</code>
'''

# Parse the HTML content
soup = BeautifulSoup(html_content, 'html.parser')

# Prettify the HTML
formatted_html = soup.prettify()

print(formatted_html)
