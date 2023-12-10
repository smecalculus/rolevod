package smecalculus.bezmen.messaging;

import smecalculus.bezmen.messaging.SepulkaMessageEm.RegistrationRequest;
import smecalculus.bezmen.messaging.SepulkaMessageEm.RegistrationResponse;

/**
 * Port: client side
 */
public interface SepulkaClient {
    RegistrationResponse register(RegistrationRequest request);
}
