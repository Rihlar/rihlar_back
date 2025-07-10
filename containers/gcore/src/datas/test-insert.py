import csv
import requests
import datetime

def send_movement_data(csv_file_path, base_url, user_id):
    """
    Reads movement data from a CSV file and sends it as POST requests.

    Args:
        csv_file_path (str): The path to the CSV file.
        base_url (str): The base URL of the API (e.g., "https://localhost:8943/gcore").
        user_id (str): The UserID to be included in the requests.
    """
    url = f"{base_url}/report/movement"
    headers = {
        "UserID": user_id,
        "Content-Type": "application/json"
    }

    try:
        with open(csv_file_path, 'r', encoding='utf-8') as csvfile:
            reader = csv.reader(csvfile)
            for i, row in enumerate(reader):
                if not row:  # Skip empty rows
                    continue
                try:
                    # Assuming the CSV format is: latitude,longitude,steps,timestamp
                    latitude = float(row[0])
                    longitude = float(row[1])
                    steps = i * 5 #int(row[2])
                    # The timestamp from your example seems to be a Unix timestamp
                    # For this example, we'll just use the current time, as the
                    # provided API request body doesn't include a timestamp.
                    # If your API expects a timestamp, you'll need to add it to the payload.

                    payload = {
                        "latitude": latitude,
                        "longitude": longitude,
                        "steps": steps
                    }

                    print(f"Sending request for row {i+1}: {payload}")
                    response = requests.post(url, headers=headers, json=payload, verify=False) # verify=False for localhost with self-signed certs
                    response.raise_for_status()  # Raise an exception for HTTP errors
                    print(f"Successfully sent data for row {i+1}. Status Code: {response.status_code}")
                except (ValueError, IndexError) as e:
                    print(f"Error processing row {i+1}: {row}. Skipping. Error: {e}")
                except requests.exceptions.RequestException as e:
                    print(f"Error sending request for row {i+1}: {e}")

    except FileNotFoundError:
        print(f"Error: CSV file not found at {csv_file_path}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")

if __name__ == "__main__":
    # --- Configuration ---
    CSV_FILE = "./log-1-movement.txt"  # Name of your CSV file
    BASE_URL = "https://localhost:8943/gcore"
    USER_ID = "userid-79541130-3275-4b90-8677-01323045aca5"
    # -------------------

    # Create a dummy CSV file for testing if it doesn't exist
    try:
        with open(CSV_FILE, 'x', encoding='utf-8') as f:
            f.write("34.7001,135.5000,6321,1751280000\n")
            f.write("34.6650,135.4950,4876,1751300000\n")
            f.write("34.7100,135.5100,9102,1751320000\n")
        print(f"Created a sample CSV file: {CSV_FILE}")
    except FileExistsError:
        print(f"Using existing CSV file: {CSV_FILE}")

    send_movement_data(CSV_FILE, BASE_URL, USER_ID)