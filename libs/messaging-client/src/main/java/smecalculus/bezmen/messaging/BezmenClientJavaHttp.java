package smecalculus.bezmen.messaging;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpRequest.BodyPublishers;
import java.net.http.HttpResponse.BodyHandlers;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import smecalculus.bezmen.messaging.SepulkaMessageEm.RegistrationRequest;
import smecalculus.bezmen.messaging.SepulkaMessageEm.RegistrationResponse;

@RequiredArgsConstructor
public class BezmenClientJavaHttp implements BezmenClient {

    @NonNull
    private ObjectMapper mapper;

    @NonNull
    private HttpClient client;

    @Override
    public RegistrationResponse register(RegistrationRequest request) {
        try {
            var requestJson = mapper.writeValueAsString(request);
            var httpRequest = HttpRequest.newBuilder()
                    .uri(URI.create("http://localhost:8080/sepulkas"))
                    .POST(BodyPublishers.ofString(requestJson))
                    .header("Content-Type", "application/json")
                    .header("Accept", "application/json")
                    .build();
            var httpResponse = client.send(httpRequest, BodyHandlers.ofString());
            return mapper.readValue(httpResponse.body(), SepulkaMessageEm.RegistrationResponse.class);
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        } catch (IOException | InterruptedException e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public boolean isReady() {
        try {
            var httpRequest = HttpRequest.newBuilder()
                    .uri(URI.create("http://localhost:8080/actuator/health"))
                    .GET()
                    .build();
            var httpResponse = client.send(httpRequest, BodyHandlers.discarding());
            return httpResponse.statusCode() == 200;
        } catch (IOException | InterruptedException e) {
            throw new RuntimeException(e);
        }
    }
}
