package smecalculus.rolevod.messaging;

import smecalculus.rolevod.messaging.SepulkaMessageEm.RegistrationRequest;
import smecalculus.rolevod.messaging.SepulkaMessageEm.RegistrationResponse;

/**
 * Port: client side
 */
public interface SepulkaClient {
    RegistrationResponse register(RegistrationRequest request);
}
