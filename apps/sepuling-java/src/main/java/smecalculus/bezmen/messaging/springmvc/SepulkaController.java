package smecalculus.bezmen.messaging.springmvc;

import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import smecalculus.bezmen.messaging.SepulkaClient;
import smecalculus.bezmen.messaging.SepulkaMessageEm.RegistrationRequest;
import smecalculus.bezmen.messaging.SepulkaMessageEm.RegistrationResponse;

@RestController
@RequestMapping("sepulkas")
@RequiredArgsConstructor
public class SepulkaController {

    @NonNull
    private SepulkaClient client;

    @PostMapping
    ResponseEntity<RegistrationResponse> register(@RequestBody RegistrationRequest requestEdge) {
        var responseEdge = client.register(requestEdge);
        return ResponseEntity.status(HttpStatus.CREATED).body(responseEdge);
    }
}
