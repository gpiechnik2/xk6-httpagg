import { check } from 'k6';
import http from 'k6/http';
import httpagg from 'k6/x/httpagg';


export default function () {
  const response = http.get('http://httpbin.test.k6.io/endpointThatWillReturn404Error');
  const status = check(
    response,
    {
      'response code was 200': (res) => res.status == 200
    }
  ); // the status variable will be false because the assertion inside does not match

  httpagg.checkRequest(
    response,
    status,
    {
        fileName: "myFilenameWithRequestsAggregated.json",
        aggregateLevel: "onSuccess" // response with the request above will not be 
        // aggregated because we set the  aggregation level to "onSuccess". 
        // The default level is "onError", which is when any of the assertions from
        // the k6 "check" function fails and the entire function returns false
    }
  );

  // or (without the optional fields)
  httpagg.checkRequest(response, status); // this request & response will be aggregated because 
  // we have not set the aggregation level and the default "onError" will be used. Additionally, 
  // a file will be created with the default name "httpagg.json"

  // or
  // IMPORTANT: We can use the "all" aggregation level to aggregate all requests regardless of 
  // the check result
  httpagg.checkRequest(response, status, {
      aggregateLevel: "all"
  });
}

export function teardown(data) {
    httpagg.generateRaport("myFilenameWithRequestsAggregated.json", "myHtmlReport.html")

    // or (without the optional fields)
    httpagg.generateRaport() // the default name of the html report that will be created 
    // is "httpaggReport.html". In turn, the name of the request results file that will 
    // be checked is "httpagg.json"
}
