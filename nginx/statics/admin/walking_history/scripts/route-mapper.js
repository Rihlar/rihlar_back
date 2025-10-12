// scripts/route-mapper.js

// アイコン設定
const startIcon = L.icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-green.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41]
});
const endIcon = L.icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41]
});
const routeColor = '#1a73e8';

/**
 * 距離から歩数を計算する関数（平均歩幅：約70cm）
 */
function calculateSteps(distanceMeters) {
    const averageStrideLength = 0.7; // メートル
    return Math.round(distanceMeters / averageStrideLength);
}

/**
 * 2つの座標間の距離（メートル）を計算する関数 (Haversine formula)
 */
function calculateDistance(lat1, lon1, lat2, lon2) {
    const R = 6371000; // 地球の半径（メートル）
    const dLat = (lat2 - lat1) * Math.PI / 180;
    const dLon = (lon2 - lon1) * Math.PI / 180;
    const a = 
        Math.sin(dLat / 2) * Math.sin(dLat / 2) +
        Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) * Math.sin(dLon / 2) * Math.sin(dLon / 2);
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    return R * c;
}


/**
 * OSRMのエンコードされたポリラインをデコードする関数
 */
function decodePolyline(encoded) {
    const inv = 1.0 / 1e5;
    const decoded = [];
    let previous = [0, 0];
    let i = 0;

    while (i < encoded.length) {
        const ll = [0, 0];
        for (let j = 0; j < 2; j++) {
            let shift = 0;
            let result = 0;
            let byte;
            do {
                byte = encoded.charCodeAt(i++) - 63;
                result |= (byte & 0x1f) << shift;
                shift += 5;
            } while (byte >= 0x20);
            ll[j] = previous[j] + (result & 1 ? ~(result >> 1) : result >> 1);
            previous[j] = ll[j];
        }
        decoded.push([ll[0] * inv, ll[1] * inv]);
    }
    return decoded;
}

export class RouteMapper {
    constructor(mapInstance) {
        this.map = mapInstance;
        this.selectedPoints = [];
        this.markers = [];
        this.displayedRoute = null;
        this.calculatedRouteData = null; // {routeCoordinates, duration, distance, steps}
        this.isActive = false; // モードがアクティブかどうかを判定するフラグ

        // DOM要素
        this.routeInstruction = document.getElementById('route-instruction');
        this.routeTimeSpan = document.getElementById('route-time');
        this.routeDistSpan = document.getElementById('route-dist');
        this.routeStepsSpan = document.getElementById('route-steps');
        this.addLogBtn = document.getElementById('add-log-btn');
        this.clearRouteBtn = document.getElementById('clear-route-btn');

        // イベントリスナー設定
        this.clearRouteBtn.addEventListener('click', () => this.clearRoute());
        this.addLogBtn.addEventListener('click', () => this.generateLogJSON());
        
        // Leafletのダブルクリックイベントをマップインスタンスに直接追加
        this.map.on('dblclick', (e) => this.handleMapDoubleClick(e.latlng));
    }
    
    /**
     * ルートマッピング機能を有効/無効にする
     */
    setActive(state) {
        this.isActive = state;
        if (!state) {
            this.clearRoute();
        }
    }

    /**
     * ルート情報をクリアし、UIをリセットする
     */
    clearRoute() {
        this.selectedPoints = [];
        this.markers.forEach(marker => this.map.removeLayer(marker));
        this.markers = [];
        if (this.displayedRoute) {
            this.map.removeLayer(this.displayedRoute);
            this.displayedRoute = null;
        }
        this.calculatedRouteData = null;
        this.routeTimeSpan.textContent = '';
        this.routeDistSpan.textContent = '';
        this.routeStepsSpan.textContent = '';
        this.routeInstruction.textContent = '地図上の2点をダブルクリックしてルートを検索してください。';
        this.addLogBtn.disabled = true;
        this.clearRouteBtn.disabled = true;
    }

    /**
     * マップダブルクリック時の処理（ログ追加モードでのみ実行）
     */
    handleMapDoubleClick(latlng) {
        if (!this.isActive) return;
        
        if (this.selectedPoints.length >= 2) {
            this.clearRoute();
        }

        this.selectedPoints.push(latlng);

        // マーカーを仮追加
        const icon = this.selectedPoints.length === 1 ? startIcon : endIcon;
        const label = this.selectedPoints.length === 1 ? '出発地' : '目的地';
        const tempMarker = L.marker(latlng, { icon: icon }).bindPopup(label).addTo(this.map).openPopup();
        this.markers.push(tempMarker);

        // 2点が選択されたらルートを計算
        if (this.selectedPoints.length === 2) {
            // 仮マーカーはクリアされるため、一旦削除
            this.markers.forEach(marker => this.map.removeLayer(marker));
            this.markers = [];

            this.calculateRoute(this.selectedPoints[0], this.selectedPoints[1]);
        }
    }

    /**
     * OSRM APIを使用してルートを計算し、地図に表示する
     */
    async calculateRoute(point1, point2) {
        this.routeInstruction.textContent = 'ルートを計算中...';
        const osrmBaseUrl = 'https://routing.openstreetmap.de/routed-foot/route/v1/driving/';
        const coordsString = `${point1.lng},${point1.lat};${point2.lng},${point2.lat}`;
        const url = `${osrmBaseUrl}${coordsString}?overview=full&geometries=polyline`;

        try {
            const response = await fetch(url);
            const data = await response.json();

            if (data.code === 'Ok' && data.routes && data.routes.length > 0) {
                // 成功したら以前の情報をクリア
                this.clearRoute();
                
                // マーカーを再追加
                this.markers.push(L.marker(point1, { icon: startIcon }).addTo(this.map).bindPopup('出発地').openPopup());
                this.markers.push(L.marker(point2, { icon: endIcon }).addTo(this.map).bindPopup('目的地'));
                this.selectedPoints.push(point1, point2); // clearRouteでリセットされたため再設定

                const primaryRoute = data.routes[0];
                const routeCoordinates = decodePolyline(primaryRoute.geometry);

                this.displayedRoute = L.polyline(routeCoordinates, {
                    color: routeColor,
                    weight: 5,
                    opacity: 0.8
                }).addTo(this.map);

                const durationInSeconds = primaryRoute.duration;
                const distanceInMeters = primaryRoute.distance;
                const steps = calculateSteps(distanceInMeters);

                // 計算結果を保存
                this.calculatedRouteData = {
                    routeCoordinates,
                    duration: durationInSeconds,
                    distance: distanceInMeters,
                    steps: steps
                };

                // UIの更新
                const durationInMinutes = Math.round(durationInSeconds / 60);
                const km = (distanceInMeters / 1000).toFixed(1);

                this.routeTimeSpan.textContent = `${durationInMinutes} 分`;
                this.routeDistSpan.textContent = `${km} km`;
                this.routeStepsSpan.textContent = `${steps.toLocaleString()} 歩`;
                this.routeInstruction.textContent = 'ルートが表示されました。「このルートをログに追加」ボタンでJSONを生成します。';
                this.addLogBtn.disabled = false;
                this.clearRouteBtn.disabled = false;

                this.map.fitBounds(this.displayedRoute.getBounds());
            } else {
                this.routeInstruction.textContent = 'ルートが見つかりませんでした。再度2点をダブルクリックしてください。';
                alert('ルートが見つかりませんでした');
                this.clearRoute();
            }
        } catch (error) {
            console.error('ルート計算中にエラーが発生しました:', error);
            this.routeInstruction.textContent = 'ルート計算中にエラーが発生しました。再度2点をダブルクリックしてください。';
            alert('ルート計算中にエラーが発生しました');
            this.clearRoute();
        }
    }

    /**
     * 計算されたルートデータから歩行ログJSONを生成する
     */
    generateLogJSON() {
        const gameSelect = document.getElementById('game-select');
        const userId = document.getElementById('user-id').textContent;

        if (!this.calculatedRouteData || !gameSelect.value) {
            alert('ルート情報またはゲームが選択されていません。');
            return;
        }

        const selectedGameId = gameSelect.value;
        const { routeCoordinates, duration, distance, steps } = this.calculatedRouteData;
        const currentTimeStamp = Math.floor(Date.now() / 1000); // 現在のUNIXタイムスタンプ（秒）
        const totalDistance = distance;
        
        let logData = [];
        let accumulatedSteps = 0;
        let accumulatedDistance = 0;

        // ログの開始時刻を現在の推定時刻（currentTimeStamp）からduration分前（約-duration秒）に設定
        const startTime = currentTimeStamp - Math.floor(duration);

        for (let i = 0; i < routeCoordinates.length; i++) {
            const coord = routeCoordinates[i];
            const isFirst = i === 0;
            const isLast = i === routeCoordinates.length - 1;
            
            let segmentDistance = 0;
            let segmentTimeOffset = 0;
            
            if (!isFirst) {
                const prevCoord = routeCoordinates[i - 1];
                // 座標間の距離を計算
                segmentDistance = calculateDistance(prevCoord[0], prevCoord[1], coord[0], coord[1]);
                accumulatedDistance += segmentDistance;
                
                // 距離の割合に基づいて歩数を割り当て
                const stepRatio = segmentDistance / totalDistance;
                let stepsForSegment = Math.round(steps * stepRatio);

                // 最後のログで残りの歩数を調整し、合計が合うようにする
                if (isLast) {
                    stepsForSegment = steps - accumulatedSteps;
                }
                accumulatedSteps += stepsForSegment;
                
                // タイムスタンプのオフセット: 総時間に対する距離の割合で時間を配分
                const timeRatio = accumulatedDistance / totalDistance;
                segmentTimeOffset = Math.floor(duration * timeRatio);
            }

            // タイムスタンプを計算 (開始時刻 + 経過時間)
            const timeStamp = isFirst ? startTime : startTime + segmentTimeOffset;

            logData.push({
                timeStamp: Math.floor(timeStamp),
                latitude: coord[0],
                longitude: coord[1],
                // 歩数は、その座標までの移動（つまり前の座標からこの座標への移動）として割り当てる
                steps: isFirst ? 0 : Math.max(0, accumulatedSteps - logData[i-1].steps), // このセグメントで稼いだ歩数
                gameID: selectedGameId,
                userID: userId
            });
            
            // NOTE: Leafletのデコードされた座標配列では、各座標がほぼ等間隔ではないため、距離に基づいて時間を配分するのがより正確です。
            // stepsは、この地点に到達するまでに稼いだ歩数として、前地点から移動した分の歩数を割り当てます。
        }

        // ログデータをJSON形式で表示
        console.log('--- 生成された歩行ログJSON ---');
        console.log(`// ユーザーID: ${userId}`);
        console.log(`// ゲームID: ${selectedGameId}`);
        console.log(`// 総歩数（推定）：${steps.toLocaleString()} 歩`);
        console.log(JSON.stringify(logData, null, 2));
        alert('生成された歩行ログJSONをコンソールに出力しました。\n(APIへの追加処理は、このコードでは実装されていません。)');

        this.clearRoute(); // ログ生成後、ルートをクリア

        return true;
    }
}
