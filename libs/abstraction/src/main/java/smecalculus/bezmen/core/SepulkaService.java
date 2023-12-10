package smecalculus.bezmen.core;

import java.util.List;
import smecalculus.bezmen.core.MessageDm.PreviewRequest;
import smecalculus.bezmen.core.MessageDm.PreviewResponse;
import smecalculus.bezmen.core.MessageDm.RegistrationRequest;
import smecalculus.bezmen.core.MessageDm.RegistrationResponse;

public interface SepulkaService {
    RegistrationResponse register(RegistrationRequest request);

    List<PreviewResponse> view(PreviewRequest request);
}
