package app

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"collector/internal/service"
)

type request struct {
	method string
	path   string
}

type want struct {
	status int
	body   string
}

func Test_handleConnection(t *testing.T) {
	a := App{
		Url:     "127.0.0.1:8380",
		Router:  mux.NewRouter(),
		Service: service.Mock{},
	}
	go a.Run()

	tests := map[string]struct {
		request request
		want    want
	}{
		"testCollector": {
			request{
				"GET",
				"http://127.0.0.1:8380/",
			},
			want{200, `{"status":true,"data":{"sms":[[{"Country":"Saint Barthélemy","Bandwidth":"68","ResponseTime":"1594","Provider":"Kildy"},{"Country":"United States of America (the)","Bandwidth":"36","ResponseTime":"1576","Provider":"Rond"}],[{"Country":"Saint Barthélemy","Bandwidth":"68","ResponseTime":"1594","Provider":"Kildy"},{"Country":"United States of America (the)","Bandwidth":"36","ResponseTime":"1576","Provider":"Rond"}]],"mms":[[{"country":"Russian Federation (the)","provider":"Kildy","bandwidth":"3","response_time":"511"}],[{"country":"Russian Federation (the)","provider":"Kildy","bandwidth":"3","response_time":"511"}]],"voice_call":[{"Country":"BG","Bandwidth":"40","ResponseTime":"609","Provider":"E-Voice","ConnectionStability":0.86,"TTFB":160,"VoicePurity":36,"MedianOfCallsTime":5},{"Country":"DK","Bandwidth":"11","ResponseTime":"743","Provider":"JustPhone","ConnectionStability":0.67,"TTFB":82,"VoicePurity":74,"MedianOfCallsTime":41}],"email":{"RU":[[{"Country":"RU","Provider":"Yahoo","DeliveryTime":124},{"Country":"RU","Provider":"Gmail","DeliveryTime":428},{"Country":"RU","Provider":"MSN","DeliveryTime":463}],[{"Country":"RU","Provider":"Gmail","DeliveryTime":428},{"Country":"RU","Provider":"MSN","DeliveryTime":463},{"Country":"RU","Provider":"Hotmail","DeliveryTime":592}]],"US":[[{"Country":"US","Provider":"Orange","DeliveryTime":45},{"Country":"US","Provider":"MSN","DeliveryTime":124},{"Country":"US","Provider":"Yahoo","DeliveryTime":305}],[{"Country":"US","Provider":"MSN","DeliveryTime":124},{"Country":"US","Provider":"Yahoo","DeliveryTime":305},{"Country":"US","Provider":"Hotmail","DeliveryTime":391}]]},"billing":{"CreateCustomer":true,"Purchase":true,"Payout":false,"Recurring":false,"FraudControl":true,"CheckoutPage":false},"support":[2,36],"incident":[{"topic":"Wrong SMS delivery time","status":"active"},{"topic":"Support overloaded because of EU affect","status":"active"},{"topic":"Billing isn’t allowed in US","status":"closed"}]},"error":""}`},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &http.Client{}
			r := strings.NewReader("")
			req, err := http.NewRequest(tc.request.method, tc.request.path, r)
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			b, _ := io.ReadAll(res.Body)
			res.Body.Close()

			assert.Equal(t, tc.want.status, res.StatusCode)
			assert.Equal(t, tc.want.body, string(b))
		})
	}
}
