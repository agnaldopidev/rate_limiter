import requests
import time

for i in range(1, 101):
    try:
        r = requests.get(
            'http://localhost:8080',
            headers={'API_KEY': 'abc123'}
        )
        print(f"Req {i}: HTTP {r.status_code} - {r.text.strip()}")
    except Exception as e:
        print(f"Req {i}: ERROR - {str(e)}")
    time.sleep(0.1)