package humstatPrinter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go-team-room/humstat"
	"go-team-room/models/context"
	//"go-team-room/controllers"
)

//HandlerStat - responces to req GET stat
func HandlerStat(w http.ResponseWriter, r *http.Request) {

	currentUserID := int(context.GetIdFromContext(r))

	if currentUserID < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "401 Can't find your ID")
		return
	}
	currentUserRole := string(context.GetRoleFromContext(r))
	if !strings.EqualFold(currentUserRole, "admin") {

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "401 You are not welcome")
		return

	}

	keys, ok := r.URL.Query()["pretty"]
	if !ok || len(keys) < 1 {
		//we going to serve json
		humstat.OutputStat <- map[string]int{"stat": 0}
		totalStat := make(map[string]int)
		totalStatTemp := <-humstat.OutputStat
		for i := range totalStatTemp {
			totalStat["total_"+i] = totalStatTemp[i]
		}
		humstat.OutputStat <- map[string]int{"stat": 1}
		minuteStat := <-humstat.OutputStat
		//jsonString, err := json.NewEncoder().Encode(totalStat)
		jsonStringTotalStat, err := json.Marshal(totalStat)
		if err != nil {
			//something wrong with total stat
		}
		jsonStringMinuteStat, err := json.Marshal(minuteStat)
		if err != nil {
			//something wrong with minute stat
		}

		//we going to serve json
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(jsonStringTotalStat), ",", string(jsonStringMinuteStat))

	} else {
		//we going to serve html
		w.WriteHeader(http.StatusOK)
		//scriptToChart := ""

		// humstat.OutputStat <- map[string]int{"stat": 0}
		// totalStat := make(map[string]int)
		// totalStatTemp := <-humstat.OutputStat
		// for i := range totalStatTemp {
		// 	totalStat["total_"+i] = totalStatTemp[i]
		// }

		fmt.Fprint(w, `<html><head><meta charset="utf-8"></head>
<body>
<script src="https://sdk.amazonaws.com/js/aws-sdk-2.195.0.min.js"></script>
<script src="/dist/chart.js"></script>
<script>AWS.config.update({	region: "eu-central-1",	accessKeyId: "AKIAIKYJK67CAJ36L4CA",secretAccessKey: "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"})</script>`)
		fmt.Fprint(w, `<script src="/dist/chart-dashboard.js"></script>
<canvas id="canvas" width="600px" height="500px"  style="width: 600px; height: 500px"></canvas>
</body></html>`)
	}
}
