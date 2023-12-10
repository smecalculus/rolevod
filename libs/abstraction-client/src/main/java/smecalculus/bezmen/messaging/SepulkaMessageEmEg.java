package smecalculus.bezmen.messaging;

import java.util.UUID;
import smecalculus.bezmen.messaging.SepulkaMessageEm.RegistrationRequest;
import smecalculus.bezmen.messaging.SepulkaMessageEm.RegistrationResponse;

public abstract class SepulkaMessageEmEg {
    public static RegistrationRequest registrationRequest() {
        var requestEdge = new RegistrationRequest();
        requestEdge.setExternalId(UUID.randomUUID().toString());
        return requestEdge;
    }

    public static RegistrationRequest registrationRequest(String id) {
        var requestEdge = registrationRequest();
        requestEdge.setExternalId(id);
        return requestEdge;
    }

    public static RegistrationResponse registrationResponse() {
        var responseEdge = new RegistrationResponse();
        responseEdge.setExternalId(UUID.randomUUID().toString());
        return responseEdge;
    }

    public static RegistrationResponse registrationResponse(String externalId) {
        var responseEdge = registrationResponse();
        responseEdge.setExternalId(externalId);
        return responseEdge;
    }
}
