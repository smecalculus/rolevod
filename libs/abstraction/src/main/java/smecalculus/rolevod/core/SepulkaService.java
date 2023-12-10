package smecalculus.rolevod.core;

import java.util.List;
import smecalculus.rolevod.core.SepulkaMessageDm.PreviewRequest;
import smecalculus.rolevod.core.SepulkaMessageDm.PreviewResponse;
import smecalculus.rolevod.core.SepulkaMessageDm.RegistrationRequest;
import smecalculus.rolevod.core.SepulkaMessageDm.RegistrationResponse;

public interface SepulkaService {
    RegistrationResponse register(RegistrationRequest request);

    List<PreviewResponse> view(PreviewRequest request);
}
