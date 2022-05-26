import http from 'k6/http';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";
import { jUnit, textSummary } from 'https://jslib.k6.io/k6-summary/0.0.1/index.js';
import { check } from 'k6';
import { httpagg } from 'k6/x/httpagg';

export const options = { vus: 5, iterations: 30 };

export default function () {

  // Send the results to some remote server or trigger a hook
  const res = http.get('http://httpbin.test.k6.io');
  check(
    res,
    {
      'response code was 200': (res) => res.status == 200,
      'body size was 1234 bytes': (res) => res.body.length == 1234,
    },
    { myTag: "I'm a tag" }
  );
  httpagg.checkRequest(JSON.stringify(res), "errors");
}

export function handleSummary(data) {
  return {
    'stdout': textSummary(data, { indent: ' ', enableColors: true }), // Show the text summary to stdout...
    './junit.xml': jUnit(data), // but also transform it and save it as a JUnit XML...
    './summary.json': JSON.stringify(data), // and a JSON with all the details...
    "./summary.html": htmlReport(data),
    // And any other JS transformation of the data you can think of,
    // you can write your own JS helpers to transform the summary data however you like!
  };
}
