import csv
import requests
import random
import os

## Function to Parse CSV
def parse_csv_data(csv_file_path):
    """
    Parses a CSV file and extracts latitude and longitude for each valid row.

    Args:
        csv_file_path (str): The path to the CSV file.

    Returns:
        list: A list of dictionaries, each containing 'latitude' and 'longitude'.
              Returns an empty list if the file is not found or no valid data is parsed.
    """
    parsed_data = []
    try:
        with open(csv_file_path, 'r', encoding='utf-8') as csvfile:
            reader = csv.reader(csvfile)
            for i, row in enumerate(reader):
                if not row:  # Skip empty rows
                    continue
                try:
                    # Assuming CSV format: latitude,longitude,steps,timestamp
                    latitude = float(row[0])
                    longitude = float(row[1])
                    parsed_data.append({"latitude": latitude, "longitude": longitude})
                except (ValueError, IndexError) as e:
                    print(f"Error parsing row {i+1} from CSV: {row}. Skipping. Error: {e}")
    except FileNotFoundError:
        print(f"Error: CSV file not found at {csv_file_path}")
    except Exception as e:
        print(f"An unexpected error occurred while parsing CSV: {e}")
    return parsed_data

## Function to Send Movement Data
def send_movement_data(data_to_send, base_url, user_id):
    """
    Sends movement data as POST requests to the /report/movement endpoint.

    Args:
        data_to_send (list): A list of dictionaries, each containing 'latitude' and 'longitude'.
        base_url (str): The base URL of the API (e.g., "https://localhost:8943/gcore").
        user_id (str): The UserID to be included in the requests.
    """
    url = f"{base_url}/report/movement"
    headers = {
        "UserID": user_id,
        "Content-Type": "application/json"
    }

    print("\n--- Sending data to /report/movement endpoint ---")
    for i, data in enumerate(data_to_send):
        try:
            # Steps are calculated based on the index (i)
            steps = i * 5 

            payload = {
                "latitude": data["latitude"],
                "longitude": data["longitude"],
                "steps": steps
            }

            print(f"({i+1}/{len(data_to_send)}) Sending /report/movement request: {payload}")
            response = requests.post(url, headers=headers, json=payload, verify=False) # verify=False for localhost with self-signed certs
            response.raise_for_status()  # Raise an exception for HTTP errors
            print(f"Success ({response.status_code}): {response.text}")
        except requests.exceptions.RequestException as e:
            print(f"Error sending request for item {i+1}: {e}")
        except Exception as e:
            print(f"An unexpected error occurred during request for item {i+1}: {e}")
    print("--- Finished sending data to /report/movement endpoint ---")


## Function to Get 10 Random Data Points
def get_random_samples(data_list, num_samples=10):
    """
    Selects a specified number of random samples from a list of data.

    Args:
        data_list (list): The list of data from which to sample.
        num_samples (int): The number of random samples to select.

    Returns:
        list: A list of randomly selected samples.
    """
    if not data_list:
        print("No data available to select samples from.")
        return []

    samples_to_take = min(num_samples, len(data_list))
    print(f"\n--- Selecting {samples_to_take} random samples ---")
    return random.sample(data_list, samples_to_take)


## Function to Create Circles from Data
def create_circles_from_data(data_for_circles, base_url, user_id, fixed_steps=10):
    """
    Sends circle creation requests to the /create/circle endpoint.

    Args:
        data_for_circles (list): A list of dictionaries (latitude, longitude) for circle creation.
        base_url (str): The base URL of the API.
        user_id (str): The UserID to be included in the requests.
        fixed_steps (int): The fixed steps value to send in the request payload.
    """
    url = f"{base_url}/create/circle"
    headers = {
        "UserID": user_id,
        "Content-Type": "application/json"
    }

    if not data_for_circles:
        print("No data provided for circle creation.")
        return

    print(f"\n--- Sending {len(data_for_circles)} circle creation requests to /create/circle endpoint ---")
    for i, data in enumerate(data_for_circles):
        payload = {
            "latitude": data["latitude"] + 0.007,
            "longitude": data["longitude"] - 0.007,
            "steps": 2000  # Fixed steps value for circle creation
        }

        try:
            print(f"({i+1}/{len(data_for_circles)}) Sending /create/circle request: {payload}")
            response = requests.post(url, headers=headers, json=payload, verify=False) # verify=False for localhost with self-signed certs
            response.raise_for_status()  # Raise an exception for HTTP errors
            print(f"Success ({response.status_code}): {response.text}")
        except requests.exceptions.RequestException as e:
            print(f"Error sending circle creation request for item {i+1}: {e}")
        except Exception as e:
            print(f"An unexpected error occurred during circle creation request for item {i+1}: {e}")
    print(f"--- Finished sending {len(data_for_circles)} circle creation requests ---")

## Main Execution Block
if __name__ == "__main__":
    # --- Configuration ---
    CSV_FILE = "./log-1-movement.txt"  # Your CSV file name
    BASE_URL = "https://localhost:8943/gcore"
    USER_ID = "userid-79541130-3275-4b90-8677-01323045aca5"
    NUM_CIRCLE_SAMPLES = 10 # Number of random samples for circle creation
    FIXED_CIRCLE_STEPS = 10 # Fixed steps value for circle creation
    # -------------------

    # Create a dummy CSV file for testing if it doesn't exist
    if not os.path.exists(CSV_FILE):
        print(f"'{CSV_FILE}' not found. Creating a dummy CSV for testing.")
        # Ensure enough dummy data for 10 random samples
        dummy_data = [
            (34.7001, 135.5000, 6321, 1751280000), (34.6650, 135.4950, 4876, 1751300000),
            (34.7100, 135.5100, 9102, 1751320000), (34.7050, 135.5050, 5500, 1751340000),
            (34.6900, 135.4800, 7200, 1751360000), (34.7150, 135.5200, 8100, 1751380000),
            (34.6800, 135.4900, 4200, 1751400000), (34.7200, 135.5300, 9900, 1751420000),
            (34.6700, 135.4700, 3100, 1751440000), (34.7005, 135.5010, 6800, 1751460000),
            (34.6950, 135.4980, 5900, 1751480000), (34.7080, 135.5120, 8500, 1751500000),
            (34.6600, 135.4900, 3900, 1751520000), (34.7120, 135.5150, 7700, 1751540000),
            (34.6850, 135.4850, 5000, 1751560000), (34.7030, 135.5070, 6100, 1751580000),
            (34.6980, 135.4990, 5300, 1751600000), (34.7110, 135.5110, 8900, 1751620000),
            (34.6670, 135.4960, 4600, 1751640000), (34.7090, 135.5130, 9500, 1751660000),
        ]
        with open(CSV_FILE, 'w', encoding='utf-8', newline='') as f:
            writer = csv.writer(f)
            for row_data in dummy_data:
                writer.writerow(row_data)
        print(f"Sample CSV file '{CSV_FILE}' created.")
    else:
        print(f"Using existing CSV file: '{CSV_FILE}'.")

    # 1. Parse the CSV file
    parsed_movement_data = parse_csv_data(CSV_FILE)

    if not parsed_movement_data:
        print("No valid data parsed from CSV. Exiting.")
    else:
        # 2. Send movement data
        send_movement_data(parsed_movement_data, BASE_URL, USER_ID)

        # 3. Get 10 random data points for circle creation
        random_circle_data = get_random_samples(parsed_movement_data, NUM_CIRCLE_SAMPLES)

        # 4. Create circles from the 10 random data points
        if random_circle_data:
            random_circle_data = get_random_samples(parsed_movement_data, NUM_CIRCLE_SAMPLES)
            create_circles_from_data(random_circle_data, BASE_URL, "userid-c3d4e5f6-a7b8-9012-3456-7890abcdef01", FIXED_CIRCLE_STEPS)

            random_circle_data = get_random_samples(parsed_movement_data, NUM_CIRCLE_SAMPLES)
            create_circles_from_data(random_circle_data, BASE_URL, "userid-b2c3d4e5-f6a7-8901-2345-67890abcdef0", FIXED_CIRCLE_STEPS)
        else:
            print("Could not select enough random samples to create circles.")