package smecalculus.bezmen.messaging;

import smecalculus.bezmen.messaging.MessageEm.RegistrationRequest;
import smecalculus.bezmen.messaging.MessageEm.RegistrationResponse;

/**
 * Port: client side
 */
public interface SepulkaClient {
    RegistrationResponse register(RegistrationRequest request);
}
