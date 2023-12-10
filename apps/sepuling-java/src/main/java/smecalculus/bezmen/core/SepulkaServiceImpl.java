package smecalculus.bezmen.core;

import static java.util.UUID.randomUUID;

import java.time.LocalDateTime;
import java.util.Collections;
import java.util.List;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import smecalculus.bezmen.core.SepulkaMessageDm.PreviewRequest;
import smecalculus.bezmen.core.SepulkaMessageDm.PreviewResponse;
import smecalculus.bezmen.core.SepulkaMessageDm.RegistrationRequest;
import smecalculus.bezmen.core.SepulkaMessageDm.RegistrationResponse;
import smecalculus.bezmen.storage.SepulkaDao;

@RequiredArgsConstructor
public class SepulkaServiceImpl implements SepulkaService {

    @NonNull
    private SepulkaMapper mapper;

    @NonNull
    private SepulkaDao dao;

    @Override
    public RegistrationResponse register(RegistrationRequest request) {
        var now = LocalDateTime.now();
        var sepulkaCreated = mapper.toState(request)
                .internalId(randomUUID())
                .revision(0)
                .createdAt(now)
                .updatedAt(now)
                .build();
        var sepulkaSaved = dao.add(sepulkaCreated);
        return mapper.toMessage(sepulkaSaved).build();
    }

    @Override
    public List<PreviewResponse> view(PreviewRequest request) {
        return Collections.emptyList();
    }
}
