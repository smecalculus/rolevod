package smecalculus.bezmen.core;

import java.util.List;
import smecalculus.bezmen.core.SepulkaMessageDm.PreviewRequest;
import smecalculus.bezmen.core.SepulkaMessageDm.PreviewResponse;
import smecalculus.bezmen.core.SepulkaMessageDm.RegistrationRequest;
import smecalculus.bezmen.core.SepulkaMessageDm.RegistrationResponse;

public interface SepulkaService {
    RegistrationResponse register(RegistrationRequest request);

    List<PreviewResponse> view(PreviewRequest request);
}
