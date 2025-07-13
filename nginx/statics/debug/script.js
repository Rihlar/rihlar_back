document.addEventListener('DOMContentLoaded', () => {
    // マップの初期化
    const map = L.map('map').setView([34.706, 135.501], 10); // 初期表示の中心座標とズームレベルを調整

    // OpenStreetMapのタイルレイヤーを追加
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    // 各順位・チームに対応する色を定義
    const teamColors = {
        'Top1': 'red',     // 1位のチーム
        'Top2': 'blue',    // 2位のチーム
        'Top3': 'green',   // 3位のチーム
        'Self': 'purple',  // 自分のチーム
        'Other': 'gray'    // その他のチーム (データには含まれないが念のため)
    };

    // --- 歩いたデータを取得する関数 ---
    /**
     * 指定されたUserIDとGameIDで移動データをAPIから取得します。
     * @param {string} userId - ユーザーID
     * @param {string} gameId - ゲームID
     * @returns {Promise<Array<Object>>} 取得した移動データの配列
     */
    async function fetchMovementData(userId, gameId) {
        const url = '/gcore/get/movement'; // 歩いたデータのパス
        const headers = {
            'UserID': userId,
            'GameID': gameId, // 歩いたデータのエンドポイントにはGameIDヘッダーが必要
            'Content-Type': 'application/json'
        };

        try {
            const response = await fetch(url, {
                method: 'GET',
                headers: headers
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            if (data && data.data) {
                return data.data;
            } else {
                console.warn("Movement API response does not contain a 'data' array:", data);
                return [];
            }
        } catch (error) {
            console.error("Failed to fetch movement data:", error);
            // alert("移動データの取得に失敗しました。"); // 開発中はコメントアウトするとデバッグしやすい
            return [];
        }
    }

    // --- 円データを取得する関数 ---
    /**
     * 指定されたUserIDとGameIDでサークルデータをAPIから取得します。
     * @param {string} userId - ユーザーID
     * @param {string} gameId - ゲームID
     * @returns {Promise<Object>} 取得したチームごとのサークルデータ
     */
    async function fetchCircleData(userId, gameId) {
        const url = `/game/ranking/top/${gameId}`; // 円データのパス
        const headers = {
            'UserID': userId,
            'Content-Type': 'application/json'
        };

        try {
            const response = await fetch(url, {
                method: 'GET',
                headers: headers
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const apiResponse = await response.json();
            if (apiResponse && apiResponse.Data) {
                return apiResponse.Data;
            } else {
                console.warn("Circle API response does not contain a 'Data' object:", apiResponse);
                return {};
            }
        } catch (error) {
            console.error("Failed to fetch circle data:", error);
            // alert("サークルデータの取得に失敗しました。"); // 開発中はコメントアウトするとデバッグしやすい
            return {};
        }
    }

    // --- 歩いたデータを地図に表示する関数 ---
    /**
     * 取得した移動データをLeaflet地図上に線で表示します。
     * @param {Array<Object>} movementData - 移動データの配列
     */
    function displayMovementData(movementData) {
        if (!movementData || movementData.length === 0) {
            console.warn("表示する移動データがありません。");
            return []; // fitBoundsのために空の配列を返す
        }

        // タイムスタンプでデータをソート
        const sortedData = [...movementData].sort((a, b) => a.timeStamp - b.timeStamp);

        // ポリラインの座標配列を作成
        const polylineCoords = sortedData.map(point => [point.latitude, point.longitude]);

        // ポリラインを地図に追加
        if (polylineCoords.length > 1) {
            // 歩いたデータは青色の実線で表示
            L.polyline(polylineCoords, { color: 'blue', weight: 4, dashArray: '0' }).addTo(map);
        } else if (polylineCoords.length === 1) {
            // 1点のみの場合、その地点に円マーカーを表示
            L.circleMarker(polylineCoords[0], {
                radius: 6,
                fillColor: "blue",
                color: "white",
                weight: 1,
                opacity: 1,
                fillOpacity: 0.8
            }).addTo(map)
              .bindPopup(`
                <b>移動データ</b><br>
                <b>Time:</b> ${new Date(sortedData[0].timeStamp * 1000).toLocaleString()}<br>
                <b>Steps:</b> ${sortedData[0].steps}<br>
                <b>Lat:</b> ${sortedData[0].latitude.toFixed(6)}<br>
                <b>Lng:</b> ${sortedData[0].longitude.toFixed(6)}
            `);
        }
        return polylineCoords; // fitBoundsのために座標を返す
    }

    // --- 取得した円データを地図に表示する関数 ---
    /**
     * 取得したサークルデータをLeaflet地図上に表示します。
     * 各チームのサークルを異なる色で描画し、ポリラインも描画します。
     * @param {Object} allTeamData - 全てのチームのサークルデータを含むオブジェクト
     */
    function displayCircleData(allTeamData) {
        if (!allTeamData || Object.keys(allTeamData).length === 0) {
            console.warn("表示するサークルデータがありません。");
            return []; // fitBoundsのために空の配列を返す
        }

        let currentCircleCoords = []; // fitBoundsのために収集する座標

        // 各チームのデータをループ処理
        for (const teamKey in allTeamData) {
            if (allTeamData.hasOwnProperty(teamKey)) {
                const team = allTeamData[teamKey];
                const circles = team.Circles;
                const teamColor = teamColors[teamKey] || 'black'; // teamColorsに定義された色、なければ黒

                if (circles && circles.length > 0) {
                    // タイムスタンプでソート
                    const sortedCircles = [...circles].sort((a, b) => a.TimeStamp - b.TimeStamp);

                    // 各サークルを円として描画（Sizeプロパティを半径として使用）
                    sortedCircles.forEach(circle => {
                        const circleOptions = {
                            radius: circle.Size, // Sizeを半径として使用 (メートル単位)
                            color: teamColor,
                            fillColor: teamColor,
                            fillOpacity: 0.15, // 円は少し薄めに表示
                            weight: 1.5
                        };
                        L.circle([circle.Latitude, circle.Longitude], circleOptions)
                            .addTo(map)
                            .bindPopup(`
                                <b>チーム:</b> ${teamKey}<br>
                                <b>CircleID:</b> ${circle.CircleID}<br>
                                <b>サイズ:</b> ${circle.Size}m<br>
                                <b>レベル:</b> ${circle.Level}<br>
                                <b>タイム:</b> ${new Date(circle.TimeStamp * 1000).toLocaleString()}<br>
                                <b>緯度:</b> ${circle.Latitude.toFixed(6)}<br>
                                <b>経度:</b> ${circle.Longitude.toFixed(6)}
                            `);
                        currentCircleCoords.push([circle.Latitude, circle.Longitude]); // fitBounds用に座標を収集
                    });

                    // ★★★ 変更点: 円同士を結ぶポリラインの描画を削除 ★★★
                    // チームの移動経路を点線で表示したい場合は、以下のコードをコメントアウトせずに残します。
                    // もし、単に円の場所を示したいだけで、円と円を結ぶ線も不要であれば、このifブロック全体をコメントアウトまたは削除してください。
                    // const polylineCoords = sortedCircles.map(circle => [circle.Latitude, circle.Longitude]);
                    // if (polylineCoords.length > 1) {
                    //     // 円の**中心**の移動経路を点線で表示
                    //     const polyline = L.polyline(polylineCoords, { color: teamColor, weight: 2, dashArray: '5, 5' }).addTo(map);
                    //     polylineCoords.forEach(coord => currentCircleCoords.push(coord)); // fitBounds用に座標を収集
                    // }
                }
            }
        }
        return currentCircleCoords; // fitBoundsのために座標を返す
    }

    // --- メイン処理 ---
    async function main() {
        const userId = "userid-79541130-3275-4b90-8677-01323045aca5"; // 指定のユーザーID
        const gameId = "gameid-413a287b-213c-414f-a287-c1397db8f9bf"; // 指定のゲームID

        // 1. 歩いたデータを取得・表示
        const movementData = await fetchMovementData(userId, gameId);
        const movementCoords = displayMovementData(movementData);

        // 2. 円データを取得・表示
        const allTeamData = await fetchCircleData(userId, gameId);
        const circleAndPolylineCoords = displayCircleData(allTeamData);

        // 3. 全ての表示要素が地図に収まるようにビューを調整
        const allCombinedCoords = [...movementCoords, ...circleAndPolylineCoords];
        if (allCombinedCoords.length > 0) {
            map.fitBounds(L.latLngBounds(allCombinedCoords), { padding: [50, 50] }); // パディングを追加して少し余裕を持たせる
        } else {
            // データが全くない場合のデフォルトビュー
            map.setView([34.706, 135.501], 10);
        }
    }

    // メイン処理の実行
    main();
});