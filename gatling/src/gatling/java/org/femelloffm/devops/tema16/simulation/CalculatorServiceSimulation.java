package org.femelloffm.devops.tema16.simulation;

import io.gatling.javaapi.core.*;

import static io.gatling.javaapi.core.CoreDsl.*;

import io.gatling.javaapi.http.*;

import static io.gatling.javaapi.http.HttpDsl.*;

public class CalculatorServiceSimulation extends Simulation {
    String calculatorProtocol = System.getProperty("protocol", "http");
    String calculatorHost = System.getProperty("host", "localhost");
    int calculatorPort = Integer.getInteger("port", 8080);

    HttpProtocolBuilder httpProtocol = http
            .baseUrl(calculatorProtocol + "://" + calculatorHost + ":" + calculatorPort + "")
            .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
            .doNotTrackHeader("1")
            .acceptLanguageHeader("en-US,en;q=0.5")
            .acceptEncodingHeader("gzip, deflate")
            .userAgentHeader("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:16.0) Gecko/20100101 Firefox/16.0");

    ScenarioBuilder scn = scenario("Go calculator stress test")
            .exec(http("History request")
                    .get("/calc/history")
                    .check(status().is(200),
                            header("Content-Type").is("application/json; charset=utf-8"),
                            bodyString().exists()));

    {
        setUp(scn.injectOpen(
                rampUsersPerSec(1).to(900).during(60),
                rampUsersPerSec(1).to(900).during(60).randomized()
        ).protocols(httpProtocol))
                .assertions(global().responseTime().max().lt(1000));
    }
}
