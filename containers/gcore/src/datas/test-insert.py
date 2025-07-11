import csv
import requests
import random
import os

def send_movement_data(csv_file_path, base_url, user_id):
    """
    Reads movement data from a CSV file and sends it as POST requests to /report/movement.

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

    all_movement_data = [] # To store data for later use in circle creation

    print("\n--- /report/movement エンドポイントにデータを送信中 ---")
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
                    steps = i * 5 # Steps calculated as i * 5
                    
                    all_movement_data.append({"latitude": latitude, "longitude": longitude})

                    payload = {
                        "latitude": latitude,
                        "longitude": longitude,
                        "steps": steps
                    }

                    print(f"({i+1}) /report/movement リクエストを送信中: {payload}")
                    response = requests.post(url, headers=headers, json=payload, verify=False) # verify=False for localhost with self-signed certs
                    response.raise_for_status()  # Raise an exception for HTTP errors
                    print(f"成功 ({response.status_code}): {response.text}")
                except (ValueError, IndexError) as e:
                    print(f"エラー: 行 {i+1} の処理中に問題が発生しました: {row}. スキップします. エラー: {e}")
                except requests.exceptions.RequestException as e:
                    print(f"エラー: 行 {i+1} のリクエスト送信中に問題が発生しました: {e}")
    except FileNotFoundError:
        print(f"エラー: CSVファイルが見つかりません: {csv_file_path}")
        return [] # Return empty list if file not found
    except Exception as e:
        print(f"CSVファイルの読み込み中に予期せぬエラーが発生しました: {e}")
        return [] # Return empty list on unexpected error
    print("--- /report/movement エンドポイントへの送信が完了しました ---")
    
    return all_movement_data

def create_random_circles(movement_data, base_url, user_id, num_samples=10, fixed_steps=10):
    """
    Selects random samples from provided movement data and sends them as POST requests
    to the /create/circle endpoint with a fixed steps value.

    Args:
        movement_data (list): A list of dictionaries, each containing 'latitude' and 'longitude'.
        base_url (str): The base URL of the API.
        user_id (str): The UserID to be included in the requests.
        num_samples (int): The number of random samples to select.
        fixed_steps (int): The fixed steps value to send in the request payload.
    """
    url = f"{base_url}/create/circle"
    headers = {
        "UserID": user_id,
        "Content-Type": "application/json"
    }

    if not movement_data:
        print("円の作成に使用できる移動データがありません。")
        return

    # Select random samples
    # ここで、movement_dataの件数がnum_samplesより少ない場合は、movement_dataの全件を選択します。
    # それ以外の場合は、num_samples件（デフォルトでは10件）を無作為に選択します。
    samples_to_take = min(num_samples, len(movement_data))
    selected_samples = random.sample(movement_data, samples_to_take)

    print(f"\n--- /create/circle エンドポイントに無作為に{samples_to_take}件のデータを送信します ---")
    for i, data in enumerate(selected_samples):
        payload = {
            "latitude": data["latitude"],
            "longitude": data["longitude"],
            "steps": fixed_steps  # Fixed steps value for circle creation
        }

        try:
            print(f"({i+1}/{samples_to_take}) /create/circle リクエストを送信中: {payload}")
            response = requests.post(url, headers=headers, json=payload, verify=False) # verify=False for localhost with self-signed certs
            response.raise_for_status()  # HTTPエラーがあれば例外を発生
            print(f"成功 ({response.status_code}): {response.text}")
        except requests.exceptions.RequestException as e:
            print(f"エラー: リクエストの送信中に問題が発生しました: {e}")
        except Exception as e:
            print(f"予期せぬエラーが発生しました: {e}")
    print(f"--- {samples_to_take}件の円作成リクエストの送信が完了しました ---")


if __name__ == "__main__":
    # --- 設定 ---
    CSV_FILE = "./log-1-movement.txt"  # 使用するCSVファイル名
    BASE_URL = "https://localhost:8943/gcore"
    USER_ID = "userid-79541130-3275-4b90-8677-01323045aca5"
    NUM_CIRCLE_SAMPLES = 10 # 円作成のために無作為に選択するデータの件数 (この値が使用されます)
    FIXED_CIRCLE_STEPS = 10 # 円作成のために固定で送信するステップ数
    # -------------------

    # テスト用にダミーのCSVファイルを作成 (存在しない場合のみ)
    if not os.path.exists(CSV_FILE):
        print(f"'{CSV_FILE}' が見つかりませんでした。テスト用のダミーCSVを作成します。")
        # ダミーデータは十分な行数を用意して、ランダム選択が機能することを確認
        dummy_data = [
            (34.7001, 135.5000, 6321, 1751280000),
            (34.6650, 135.4950, 4876, 1751300000),
            (34.7100, 135.5100, 9102, 1751320000),
            (34.7050, 135.5050, 5500, 1751340000),
            (34.6900, 135.4800, 7200, 1751360000),
            (34.7150, 135.5200, 8100, 1751380000),
            (34.6800, 135.4900, 4200, 1751400000),
            (34.7200, 135.5300, 9900, 1751420000),
            (34.6700, 135.4700, 3100, 1751440000),
            (34.7005, 135.5010, 6800, 1751460000),
            (34.6950, 135.4980, 5900, 1751480000),
            (34.7080, 135.5120, 8500, 1751500000),
            (34.6600, 135.4900, 3900, 1751520000),
            (34.7120, 135.5150, 7700, 1751540000),
            (34.6850, 135.4850, 5000, 1751560000),
            (34.7030, 135.5070, 6100, 1751580000),
            (34.6980, 135.4990, 5300, 1751600000),
            (34.7110, 135.5110, 8900, 1751620000),
            (34.6670, 135.4960, 4600, 1751640000),
            (34.7090, 135.5130, 9500, 1751660000),
        ]
        with open(CSV_FILE, 'w', encoding='utf-8', newline='') as f:
            writer = csv.writer(f)
            for row_data in dummy_data:
                writer.writerow(row_data)
        print(f"サンプルCSVファイル '{CSV_FILE}' を作成しました。")
    else:
        print(f"既存のCSVファイル '{CSV_FILE}' を使用します。")

    # 1. /report/movement エンドポイントにデータを送信
    movement_data_for_circles = send_movement_data(CSV_FILE, BASE_URL, USER_ID)

    # 2. /create/circle エンドポイントに無作為に選択したデータを送信
    create_random_circles(movement_data_for_circles, BASE_URL, USER_ID, NUM_CIRCLE_SAMPLES, FIXED_CIRCLE_STEPS)