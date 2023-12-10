package smecalculus.bezmen.core;

import java.util.UUID;
import smecalculus.bezmen.core.MessageDm.RegistrationResponse;

public class MessageDmEg {
    public static RegistrationResponse.Builder registrationResponse() {
        return RegistrationResponse.builder().externalId(UUID.randomUUID().toString());
    }

    public static RegistrationResponse.Builder registrationResponse(String externalId) {
        return registrationResponse().externalId(externalId);
    }
}
