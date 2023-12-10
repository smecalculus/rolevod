package smecalculus.rolevod.messaging;

import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import smecalculus.rolevod.core.SepulkaService;
import smecalculus.rolevod.messaging.SepulkaMessageEm.RegistrationRequest;
import smecalculus.rolevod.messaging.SepulkaMessageEm.RegistrationResponse;
import smecalculus.rolevod.validation.EdgeValidator;

@RequiredArgsConstructor
public class SepulkaClientImpl implements SepulkaClient {

    @NonNull
    private EdgeValidator validator;

    @NonNull
    private SepulkaMessageMapper mapper;

    @NonNull
    private SepulkaService service;

    @Override
    public RegistrationResponse register(RegistrationRequest requestEdge) {
        validator.validate(requestEdge);
        var request = mapper.toDomain(requestEdge);
        var response = service.register(request);
        return mapper.toEdge(response);
    }
}
